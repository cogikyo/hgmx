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

const Path = "/theme/"

const sheet = Path + "palette.css"

type Handler struct {
	*static.Server
	head     template.HTML
	explorer assetURLs
}

type assetURLs struct{ css, js string }

func New() (*Handler, error) {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	srv, err := static.New(Path, sub)
	if err != nil {
		return nil, err
	}
	h := &Handler{Server: srv}
	switcher, err := h.URL("theme.js")
	if err != nil {
		return nil, err
	}
	if h.explorer.css, err = h.URL("explorer.css"); err != nil {
		return nil, err
	}
	if h.explorer.js, err = h.URL("explorer.js"); err != nil {
		return nil, err
	}
	h.head = template.HTML(fmt.Sprintf("<script type=\"module\" src=\"%s\"></script>", switcher))
	return h, nil
}

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

func Sheet(name string, s palette.Seed) string {
	q := url.Values{"name": {name}, "base": {s.Base.Hex()}, "primary": {s.Primary.Hex()}}
	if s.Accent != nil {
		q.Set("accent", s.Accent.Hex())
	}
	return sheet + "?" + q.Encode()
}
