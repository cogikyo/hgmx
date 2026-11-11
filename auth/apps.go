package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/cogikyo/hgmx/internal/appid"
	"github.com/cogikyo/hgmx/internal/webhost"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
)

// App registers one platform; set exactly one of IOS and Android.
type App struct {
	ID      string // Public ID; never reassign to another platform or host.
	IOS     *IOS
	Android *Android
	Build   int // Minimum integer CFBundleVersion or versionCode.
	Store   string
	Google  []GoogleToken
}

type IOS struct {
	Prefix string // Provisioning app ID prefix, not necessarily the Team ID.
	Bundle string
}

type Android struct {
	Package string
	Signers [][]byte // SHA-256 of installed builds' signing certificate DER.
}

type GoogleToken struct {
	Audience  string // Exact aud claim.
	Presenter string // Exact azp claim for this audience.
}

type app struct {
	id      string
	host    *host
	kind    authv1.ClientKind
	origins []string
	build   int
	store   string
	google  []GoogleToken
	wa      *webauthn.WebAuthn
}

var (
	iosPrefix      = regexp.MustCompile(`^[A-Z0-9]{10}$`)
	googleClient   = regexp.MustCompile(`^[A-Za-z0-9-]+\.apps\.googleusercontent\.com$`)
	stores         = map[authv1.ClientKind]string{authv1.ClientKind_CLIENT_KIND_IOS: "apps.apple.com", authv1.ClientKind_CLIENT_KIND_ANDROID: "play.google.com"}
	associationAge = 24 * time.Hour
)

func parseApps(hs []Host, hosts map[string]*host, rp string) (map[string]*app, error) {
	apps := map[string]*app{}
	platforms := map[string]string{}
	for _, h := range hs {
		origin, _, _ := webhost.Parse(h.Origin)
		for _, c := range h.Apps {
			a, ids, err := parseApp(c, hosts[origin], rp)
			if err != nil {
				return nil, fmt.Errorf("auth: host %s: app %q: %w", origin, c.ID, err)
			}
			if apps[a.id] != nil {
				return nil, fmt.Errorf("auth: app %q is listed twice", a.id)
			}
			for _, id := range ids {
				if other, ok := platforms[id]; ok {
					return nil, fmt.Errorf("auth: apps %q and %q share %s", other, a.id, id)
				}
				platforms[id] = a.id
			}
			apps[a.id] = a
		}
	}
	if len(apps) > 0 && hosts["https://"+rp] == nil {
		return nil, fmt.Errorf("auth: apps need https://%s in Config.Hosts to serve app associations", rp)
	}
	return apps, nil
}

func parseApp(c App, h *host, rp string) (*app, []string, error) {
	switch {
	case !strings.HasPrefix(h.origin, "https://"):
		return nil, nil, fmt.Errorf("apps need an https origin")
	case !appid.Valid(c.ID):
		return nil, nil, fmt.Errorf("ID must be lowercase letters, digits, dots and dashes")
	case (c.IOS == nil) == (c.Android == nil):
		return nil, nil, fmt.Errorf("set exactly one of IOS and Android")
	case c.Build < 1:
		return nil, nil, fmt.Errorf("Build must be positive")
	}
	a := &app{id: c.ID, host: h, build: c.Build, google: slices.Clone(c.Google)}
	var ids []string
	if c.IOS != nil {
		if !iosPrefix.MatchString(c.IOS.Prefix) || !appid.ValidBundle(c.IOS.Bundle) {
			return nil, nil, fmt.Errorf("IOS needs a 10-character app ID prefix and a bundle ID")
		}
		a.kind = authv1.ClientKind_CLIENT_KIND_IOS
		a.origins = []string{"https://" + rp}
		ids = []string{c.IOS.id()}
	} else {
		if !appid.ValidPackage(c.Android.Package) || len(c.Android.Signers) == 0 {
			return nil, nil, fmt.Errorf("Android needs a package name and signing certificate digests")
		}
		a.kind = authv1.ClientKind_CLIENT_KIND_ANDROID
		for _, s := range c.Android.Signers {
			if len(s) != sha256.Size {
				return nil, nil, fmt.Errorf("Android signer digests must be SHA-256")
			}
			o := "android:apk-key-hash:" + base64.RawURLEncoding.EncodeToString(s)
			if slices.Contains(a.origins, o) {
				return nil, nil, fmt.Errorf("Android signer %s is listed twice", fingerprint(s))
			}
			a.origins = append(a.origins, o)
		}
		ids = []string{"android:" + c.Android.Package}
	}
	u, err := url.Parse(c.Store)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Host != stores[a.kind] {
		return nil, nil, fmt.Errorf("Store must be an https://%s URL", stores[a.kind])
	}
	a.store = u.String()
	for i, g := range c.Google {
		if !googleClient.MatchString(g.Audience) || !googleClient.MatchString(g.Presenter) {
			return nil, nil, fmt.Errorf("Google needs OAuth client IDs for Audience and Presenter")
		}
		if slices.Contains(c.Google[:i], g) {
			return nil, nil, fmt.Errorf("Google lists %s/%s twice", g.Audience, g.Presenter)
		}
	}
	return a, ids, nil
}

func (c *IOS) id() string { return c.Prefix + "." + c.Bundle }

func fingerprint(digest []byte) string {
	h := strings.ToUpper(hex.EncodeToString(digest))
	parts := make([]string, 0, len(digest))
	for i := 0; i < len(h); i += 2 {
		parts = append(parts, h[i:i+2])
	}
	return strings.Join(parts, ":")
}

func (a *Auth) verifiers(name string) error {
	for _, p := range a.apps {
		c := webauthnConfig(a.rp, name, []string{"https://" + a.rp})
		if p.kind == authv1.ClientKind_CLIENT_KIND_ANDROID {
			c.RPOpaqueOrigins = p.origins
		}
		wa, err := webauthn.New(c)
		if err != nil {
			return fmt.Errorf("auth: app %q: %w", p.id, err)
		}
		p.wa = wa
	}
	return nil
}

type association struct {
	path string
	body []byte
	etag string
}

func associations(hs []Host) []association {
	web := []string{}
	type target struct {
		Namespace string   `json:"namespace"`
		Package   string   `json:"package_name"`
		Prints    []string `json:"sha256_cert_fingerprints"`
	}
	type statement struct {
		Relation []string `json:"relation"`
		Target   target   `json:"target"`
	}
	links := []statement{}
	for _, h := range hs {
		for _, c := range h.Apps {
			if c.IOS != nil {
				web = append(web, c.IOS.id())
				continue
			}
			prints := make([]string, len(c.Android.Signers))
			for i, s := range c.Android.Signers {
				prints[i] = fingerprint(s)
			}
			links = append(links, statement{
				Relation: []string{"delegate_permission/common.get_login_creds"},
				Target:   target{Namespace: "android_app", Package: c.Android.Package, Prints: prints},
			})
		}
	}
	aasa := map[string]any{"webcredentials": map[string]any{"apps": web}}
	var out []association
	for path, v := range map[string]any{
		"/.well-known/apple-app-site-association": aasa,
		"/.well-known/assetlinks.json":            links,
	} {
		body, _ := json.Marshal(v)
		sum := sha256.Sum256(body)
		out = append(out, association{path: path, body: body, etag: `"` + base64.RawURLEncoding.EncodeToString(sum[:16]) + `"`})
	}
	return out
}

func (a *Auth) mountAssociations(mux *http.ServeMux) {
	for _, s := range a.links {
		mux.HandleFunc("GET "+s.path, func(w http.ResponseWriter, r *http.Request) {
			if h, ok := a.host(r); !ok || h.origin != "https://"+a.rp {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(associationAge/time.Second)))
			w.Header().Set("ETag", s.etag)
			http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(s.body))
		})
	}
}
