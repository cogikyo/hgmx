// Package theme serves palette stylesheets, theme switching, and an explorer.
// Mount the handler with its full prefix:
//
//	h := theme.New()
//	root.Handle(theme.Path, h)
//
// Include h.Head() in the layout and link a stylesheet from [Sheet].
package theme

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"time"

	"github.com/cogikyo/hgmx/palette"
	"github.com/cogikyo/hgmx/static"
)

//go:embed assets
var assets embed.FS

// Path is the asset and palette endpoint prefix; mount without StripPrefix.
const Path = "/theme/"

const sheet = Path + "palette.css"

// Handler serves theme assets and /theme/palette.css.
// The query requires name, base, and primary; accent is optional.
// Colours use #RRGGBB, and download=css or download=tailwind exports a file.
type Handler struct {
	*static.Server
	head     template.HTML
	explorer assetURLs
}

type assetURLs struct{ css, js string }

func New() *Handler {
	sub, _ := fs.Sub(assets, "assets")
	srv, err := static.New(Path, sub)
	if err != nil {
		panic(err)
	}
	h := &Handler{Server: srv}
	switcher, _ := h.URL("theme.js")
	h.explorer.css, _ = h.URL("explorer.css")
	h.explorer.js, _ = h.URL("explorer.js")
	h.head = template.HTML(fmt.Sprintf("<script type=\"module\" src=\"%s\"></script>", switcher))
	return h
}

// Head returns the layout's theme-switching module tag.
// Click data-theme-set="name" to select a loaded theme on the document root.
// A select or input with data-theme-set selects its value on change.
// Click data-theme-live="query" to load a palette.css seed query as theme live.
func (h *Handler) Head() template.HTML { return h.head }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != sheet {
		h.Server.ServeHTTP(w, r)
		return
	}
	q := r.URL.Query()
	name := q.Get("name")
	sel, err := palette.Selector(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	seed, err := palette.ParseSeed(q.Get("base"), q.Get("primary"), q.Get("accent"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	p := palette.Generate(seed)
	switch q.Get("download") {
	case "":
		p.CSS(&buf, sel)
	case "css", "tailwind":
		p.CSS(&buf, ":root, "+sel)
		file := name + ".css"
		if q.Get("download") == "tailwind" {
			buf.WriteByte('\n')
			p.Theme(&buf)
			file = name + ".tailwind.css"
		}
		w.Header().Set("Content-Disposition", `attachment; filename="`+file+`"`)
	default:
		http.Error(w, "theme: download must be css or tailwind", http.StatusBadRequest)
		return
	}
	sum := sha256.Sum256(buf.Bytes())
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("ETag", `"`+hex.EncodeToString(sum[:8])+`"`)
	http.ServeContent(w, r, "palette.css", time.Time{}, bytes.NewReader(buf.Bytes()))
}

// Sheet returns a palette.css URL; the endpoint validates the name and seed.
func Sheet(name string, s palette.Seed) string {
	q := url.Values{"name": {name}, "base": {s.Base.Hex()}, "primary": {s.Primary.Hex()}}
	if s.Accent != nil {
		q.Set("accent", s.Accent.Hex())
	}
	return sheet + "?" + q.Encode()
}
