// Package icons provides SVG templ components and hover animations.
// Render a name directly in templ, for example @icons.Check.
// To enable animations, mount the assets and include [Handler.Head] in the layout:
//
//	h, err := icons.New()
//	if err != nil {
//		return err
//	}
//	root.Handle(icons.Path, h)
package icons

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/a-h/templ"

	"github.com/cogikyo/hgmx/static"
)

//go:generate go run ./fetch

// Path is the animation asset prefix; mount without StripPrefix.
const Path = "/icons/"

//go:embed svg
var svgs embed.FS

//go:embed assets
var assets embed.FS

//go:embed tags.json
var tagsJSON []byte

// Hover is spread onto a hover target to animate its first SVG descendant.
// It requires _hyperscript and also handles mouse presses and touch starts.
var Hover = templ.Attributes{
	"_": "on mouseenter or mousedown or touchstart queue none send hovered to first <svg/> in me then wait 420ms",
}

// Name is a 1em SVG templ component; the generated constants name bundled icons.
// An unknown name returns a render error without writing anything.
// Only icons with animations emit data-icon; static icons need no scripts.
type Name string

type glyph struct {
	markup   string
	tags     []string
	animated bool
}

var (
	glyphs = map[Name]glyph{}
	names  []Name
)

func init() {
	var tags map[Name][]string
	if err := json.Unmarshal(tagsJSON, &tags); err != nil {
		panic(fmt.Errorf("icons: tags.json: %w", err))
	}
	files, err := fs.Glob(svgs, "svg/*.svg")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := Name(strings.TrimSuffix(path.Base(file), ".svg"))
		data, err := svgs.ReadFile(file)
		if err != nil {
			panic(err)
		}
		_, err = fs.Stat(assets, "assets/anim/"+string(name)+".js")
		animated := err == nil
		markup, err := render(name, data, animated)
		if err != nil {
			panic(fmt.Errorf("icons: %s: %w", file, err))
		}
		glyphs[name] = glyph{markup: markup, tags: tags[name], animated: animated}
		names = append(names, name)
	}
}

func render(name Name, data []byte, animated bool) (string, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	tok, err := d.Token()
	if err != nil {
		return "", err
	}
	root, ok := tok.(xml.StartElement)
	if !ok || root.Name.Local != "svg" {
		return "", fmt.Errorf("want an <svg> root, got %T", tok)
	}
	body, ok := bytes.CutSuffix(bytes.TrimSpace(data[d.InputOffset():]), []byte("</svg>"))
	if !ok {
		return "", fmt.Errorf("want </svg> at the end")
	}

	var b strings.Builder
	b.WriteString("<svg")
	if animated {
		fmt.Fprintf(&b, ` data-icon="%s"`, html.EscapeString(string(name)))
	}
	b.WriteString(` width="1em" height="1em"`)
	for _, a := range root.Attr {
		key := a.Name.Local
		if a.Name.Space != "" {
			key = a.Name.Space + ":" + key
		}
		fmt.Fprintf(&b, ` %s="%s"`, key, html.EscapeString(a.Value))
	}
	b.WriteString(">")
	sep := ""
	for line := range bytes.Lines(body) {
		line = bytes.TrimSpace(line)
		b.WriteString(sep)
		b.Write(line)
		sep = ""
		if !bytes.HasSuffix(line, []byte(">")) {
			sep = " "
		}
	}
	b.WriteString("</svg>")
	return b.String(), nil
}

func (n Name) Render(ctx context.Context, w io.Writer) error {
	g, ok := glyphs[n]
	if !ok {
		return fmt.Errorf("icons: no icon %q", n)
	}
	_, err := io.WriteString(w, g.markup)
	return err
}

// Animated reports whether n has a bundled hover animation.
func (n Name) Animated() bool {
	return glyphs[n].animated
}

// Tags returns a copy of n's search tags.
func (n Name) Tags() []string {
	return slices.Clone(glyphs[n].tags)
}

// All returns a copy of the bundled names for browsing or selection controls.
func All() []Name {
	return slices.Clone(names)
}

// Match finds case-insensitive substrings in names and tags.
// An empty or whitespace-only query returns all icons.
func Match(q string) []Name {
	q = strings.ToLower(strings.TrimSpace(q))
	var found []Name
	for _, n := range names {
		if strings.Contains(string(n), q) || slices.ContainsFunc(glyphs[n].tags, func(tag string) bool {
			return strings.Contains(strings.ToLower(tag), q)
		}) {
			found = append(found, n)
		}
	}
	return found
}

// Handler serves the styles and per-icon animation modules.
type Handler struct {
	*static.Server
	head template.HTML
}

func New() (*Handler, error) {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	s, err := static.New(Path, sub)
	if err != nil {
		return nil, fmt.Errorf("icons: %w", err)
	}
	h := &Handler{Server: s}
	css, err := h.URL("icons.css")
	if err != nil {
		return nil, err
	}
	js, err := h.URL("icons.js")
	if err != nil {
		return nil, err
	}
	h.head = template.HTML(fmt.Sprintf(
		"<link rel=\"stylesheet\" href=\"%s\">\n<script type=\"module\" src=\"%s\"></script>",
		css, js,
	))
	return h, nil
}

// Head returns layout tags; place them after the script defining global Motion.
// Animations load per icon on first hover; reduced motion disables all playback.
func (h *Handler) Head() template.HTML {
	return h.head
}
