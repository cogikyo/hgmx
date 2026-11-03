package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

// Host shares accounts and passkeys, but binds sessions to its origin.
type Host struct {
	// Origin is the public HTTPS origin; HTTP is allowed only on loopback hosts.
	Origin string
	Mode   Mode
	// Allow gates sign-in, CLI links, account changes, User and Require.
	// A nil Allow admits every user.
	Allow func(User) bool
	// Layout wraps auth's page root in a full document; required by New.
	Layout func(r *http.Request, title string, body templ.Component) templ.Component
	Views  Views
	Apps   []App // Native clients bound to this host's API origin.
}

type host struct {
	origin   string
	mode     Mode
	allow    func(User) bool
	layout   func(*http.Request, string, templ.Component) templ.Component
	views    Views
	sends    *limiter
	verifies *limiter
	begins   *limiter
}

func (h *host) admits(u User) bool { return h.allow == nil || h.allow(u) }

func parseHosts(hs []Host) (map[string]*host, string, error) {
	if len(hs) == 0 {
		return nil, "", errors.New("auth: Config.Hosts needs at least one host")
	}
	hosts := make(map[string]*host, len(hs))
	var rp string
	for _, h := range hs {
		origin, name, err := parseOrigin(h.Origin)
		if err != nil {
			return nil, "", err
		}
		if h.Mode < Operator || h.Mode > Open {
			return nil, "", fmt.Errorf("auth: host %s: Mode must be Operator, Invite or Open", origin)
		}
		if hosts[origin] != nil {
			return nil, "", fmt.Errorf("auth: host %s is listed twice", origin)
		}
		if _, err := netip.ParseAddr(name); err == nil && len(hs) > 1 {
			return nil, "", fmt.Errorf("auth: IP host %s must be the only host", origin)
		}
		site, err := registrable(name, len(hs) == 1)
		if err != nil {
			return nil, "", fmt.Errorf("auth: host %s: %w", origin, err)
		}
		if rp != "" && site != rp {
			return nil, "", fmt.Errorf("auth: host %s is not on %s", origin, rp)
		}
		rp = site
		hosts[origin] = &host{
			origin:   origin,
			mode:     h.Mode,
			allow:    h.Allow,
			layout:   h.Layout,
			views:    h.Views.withDefaults(),
			sends:    newLimiter(20, time.Hour),
			verifies: newLimiter(20, time.Hour),
			begins:   newLimiter(30, 10*time.Minute),
		}
	}
	return hosts, rp, nil
}

func registrable(name string, alone bool) (string, error) {
	switch {
	case loopback(name) && alone:
		return name, nil
	case name == "localhost" || strings.HasSuffix(name, ".localhost"):
		return "localhost", nil
	}
	return publicsuffix.EffectiveTLDPlusOne(name)
}

func parseOrigin(raw string) (origin, name string, err error) {
	bad := fmt.Errorf("auth: host origin %q must be https://host, or http:// on localhost", raw)
	u, err := url.Parse(raw)
	if err != nil || u.User != nil || strings.TrimSuffix(u.Path, "/") != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", "", bad
	}
	key, name, ok := authority(u.Scheme, u.Host)
	_, err = netip.ParseAddr(name)
	domain := err != nil
	if !ok || !(u.Scheme == "https" && (domain || loopback(name)) || u.Scheme == "http" && loopback(name)) {
		return "", "", bad
	}
	return u.Scheme + "://" + key, name, nil
}

func authority(scheme, hostport string) (key, name string, ok bool) {
	u, err := url.Parse(scheme + "://" + hostport)
	if err != nil || hostport == "" || u.Host != hostport || strings.HasSuffix(hostport, ":") {
		return "", "", false
	}
	name, port := u.Hostname(), u.Port()
	if ip, err := netip.ParseAddr(name); err == nil {
		if ip.Zone() != "" {
			return "", "", false
		}
		name = ip.String()
	} else if name, err = idna.Lookup.ToASCII(name); err != nil || name == "" || strings.HasSuffix(name, ".") {
		return "", "", false
	}
	key = name
	if strings.Contains(name, ":") {
		key = "[" + name + "]"
	}
	if port != "" {
		n, err := strconv.ParseUint(port, 10, 16)
		if err != nil || n == 0 {
			return "", "", false
		}
		if p := strconv.FormatUint(n, 10); p != defaultPort[scheme] {
			key += ":" + p
		}
	}
	return key, name, true
}

var defaultPort = map[string]string{"https": "443", "http": "80"}

func (a *Auth) host(r *http.Request) (*host, bool) {
	for _, scheme := range [...]string{"https", "http"} {
		if key, _, ok := authority(scheme, r.Host); ok {
			if h, ok := a.hosts[scheme+"://"+key]; ok {
				return h, true
			}
		}
	}
	return nil, false
}

func misdirected(w http.ResponseWriter) {
	http.Error(w, "unknown host", http.StatusMisdirectedRequest)
}
