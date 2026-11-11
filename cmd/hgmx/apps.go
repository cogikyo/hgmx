package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
	"unicode"

	"github.com/cogikyo/hgmx"
	"github.com/cogikyo/hgmx/internal/appid"
	"github.com/cogikyo/hgmx/internal/webhost"
	"github.com/cogikyo/hgmx/starters"
)

type identity struct {
	name, origin, appID, google string
	pkg, bundle                 string
	rp                          string
}

type starter struct {
	dir      string
	identity func(id identity) (texts, paths []string, err error)
	markers  []string
}

var appStarters = map[string]starter{
	"android": {
		dir: "apps/android",
		identity: func(id identity) (texts, paths []string, err error) {
			if !appid.ValidPackage(id.pkg) {
				return nil, nil, fmt.Errorf("--package %q: want an Android package name such as com.example.notes", id.pkg)
			}
			texts = []string{
				"dev.hgmx.starter", id.pkg,
				`"\"starter-android\""`, fmt.Sprintf(`"\"%s\""`, id.appID),
				`"\"https://starter.hgmx.dev\""`, fmt.Sprintf(`"\"%s\""`, id.origin),
				`"GOOGLE_CLIENT_ID", "\"\""`, fmt.Sprintf(`"GOOGLE_CLIENT_ID", "\"%s\""`, id.google),
				`"app_name", "Starter"`, `"app_name", "` + resText(id.name) + `"`,
				"hgmx.starter=dev", "hgmx.starter=" + revision(),
			}
			paths = []string{"/kotlin/dev/hgmx/starter/", "/kotlin/" + strings.ReplaceAll(id.pkg, ".", "/") + "/"}
			return texts, paths, nil
		},
		markers: []string{"dev.hgmx.starter", "dev/hgmx/starter", "starter-android", "starter.hgmx.dev", "hgmx.starter=dev"},
	},
	"ios": {
		dir: "apps/ios",
		identity: func(id identity) (texts, paths []string, err error) {
			if !appid.ValidBundle(id.bundle) {
				return nil, nil, fmt.Errorf("--bundle-id %q: want a reverse-DNS bundle ID such as com.example.notes", id.bundle)
			}
			if id.google != "" {
				return nil, nil, errors.New("--google-client-id: the iOS starter has no Google sign-in yet")
			}
			if strings.Contains(id.name, "$") || strings.Contains(id.name, "//") || strings.HasSuffix(id.name, ";") {
				return nil, nil, fmt.Errorf("--name %q: an Xcode config value can't hold $ or //, or end in ;", id.name)
			}
			if id.rp == "localhost" {
				return nil, nil, fmt.Errorf("--origin %q: iOS passkeys need a public domain; point the app at localhost through Info.plist instead", id.origin)
			}
			texts = []string{
				"PRODUCT_BUNDLE_IDENTIFIER = dev.hgmx.starter", "PRODUCT_BUNDLE_IDENTIFIER = " + id.bundle,
				"INFOPLIST_KEY_CFBundleDisplayName = Starter", "INFOPLIST_KEY_CFBundleDisplayName = " + id.name,
				"HGMX_API_HOST = starter.hgmx.dev", "HGMX_API_HOST = " + strings.TrimPrefix(id.origin, "https://"),
				"HGMX_RP_ID = hgmx.dev", "HGMX_RP_ID = " + id.rp,
				"HGMX_APP_ID = starter-ios", "HGMX_APP_ID = " + id.appID,
				"HGMX_STARTER = dev", "HGMX_STARTER = " + revision(),
			}
			return texts, nil, nil
		},
		markers: []string{"dev.hgmx.starter", "starter-ios", "starter.hgmx.dev", "HGMX_RP_ID = hgmx.dev", "HGMX_STARTER = dev"},
	},
}

const compatInitial = `{
  "base": "origin/master",
  "releases": []
}
`

func initAppsCmd(platform string, id identity) (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	files, err := adopt(s, platform, id)
	if err != nil {
		log.Error("Refusing to adopt the "+platform+" starter", slog.String("error", err.Error()))
		return 1
	}
	if err := install(s, appStarters[platform].dir, files); err != nil {
		log.Error("Failed to write the starter", slog.String("error", err.Error()))
		return 1
	}
	log.Info("Adopted the starter; it is site code now", slog.String("platform", platform), slog.String("dir", appStarters[platform].dir), slog.Int("files", len(files)))
	log.Info("Next: run go tool hgmx gen, then build the app")
	return 0
}

func install(s site, dir string, files []starterFile) error {
	root, err := os.OpenRoot(s.dir)
	if err != nil {
		return err
	}
	defer root.Close()
	tmp := path.Join(path.Dir(dir), "."+path.Base(dir)+".hgmx-"+rand.Text())
	for _, f := range files {
		name := path.Join(tmp, f.path)
		if err := root.MkdirAll(path.Dir(name), 0o755); err != nil {
			return errors.Join(err, root.RemoveAll(tmp))
		}
		if err := root.WriteFile(name, f.data, f.mode); err != nil {
			return errors.Join(err, root.RemoveAll(tmp))
		}
	}
	if err := root.Rename(tmp, dir); err != nil {
		return errors.Join(err, root.RemoveAll(tmp))
	}
	if _, err := root.Lstat(compatFile); errors.Is(err, fs.ErrNotExist) {
		return writeRoot(root, compatFile, []byte(compatInitial))
	}
	return nil
}

type starterFile struct {
	path string
	data []byte
	mode fs.FileMode
}

func adopt(s site, platform string, id identity) ([]starterFile, error) {
	st, ok := appStarters[platform]
	if !ok {
		return nil, fmt.Errorf("no %s starter; want one of: %s", platform, strings.Join(slices.Sorted(maps.Keys(appStarters)), ", "))
	}
	origin, rp, err := webhost.RelyingParty(id.origin)
	if err != nil || !strings.HasPrefix(origin, "https://") {
		return nil, fmt.Errorf("--origin %q: want the exact https API origin on a domain name, such as https://example.com", id.origin)
	}
	id.origin, id.rp = origin, rp
	if !appid.Valid(id.appID) {
		return nil, fmt.Errorf("--app-id %q: want the public app ID from the site's auth.App, such as notes-android", id.appID)
	}
	if id.name == "" || strings.TrimSpace(id.name) != id.name || strings.IndexFunc(id.name, func(r rune) bool { return unicode.IsControl(r) || unicode.In(r, unicode.Zl, unicode.Zp) }) >= 0 ||
		strings.ContainsAny(id.name, "\"\\") || strings.ContainsAny(id.name[:1], "@?") {
		return nil, fmt.Errorf("--name %q: want a one-line display name without double quotes, backslashes or surrounding spaces, not starting with @ or ?", id.name)
	}
	if strings.ContainsAny(id.google, "\"\\ \n") {
		return nil, fmt.Errorf("--google-client-id %q: want the Google web client ID", id.google)
	}
	texts, paths, err := st.identity(id)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(s.path(st.dir)); err == nil {
		return nil, fmt.Errorf("%s exists; the starter is copied once and is site code afterwards", st.dir)
	}
	replace := strings.NewReplacer(texts...)
	movePath := strings.NewReplacer(paths...)
	var files []starterFile
	var leftover []string
	err = fs.WalkDir(starters.FS, st.dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := starters.FS.ReadFile(p)
		if err != nil {
			return err
		}
		rel := movePath.Replace(strings.TrimPrefix(p, st.dir+"/"))
		mode := fs.FileMode(0o644)
		if path.Base(p) == "gradlew" {
			mode = 0o755
		}
		if !bytes.ContainsRune(data, 0) {
			data = []byte(replace.Replace(string(data)))
			for _, m := range st.markers {
				if strings.Contains(string(data), m) || strings.Contains(rel, m) {
					leftover = append(leftover, rel+": "+m)
				}
			}
		}
		files = append(files, starterFile{rel, data, mode})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(leftover) > 0 {
		return nil, fmt.Errorf("starter identity left after substitution:\n%s", strings.Join(leftover, "\n"))
	}
	return files, nil
}

func resText(s string) string {
	return strings.NewReplacer("'", `\\'`, "$", `\$`).Replace(s)
}

func revision() string {
	if v := hgmx.Version(); v != "" && v != "(devel)" {
		return v
	}
	return "devel"
}
