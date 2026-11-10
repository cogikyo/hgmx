// Package glyph renders site-owned icon sets and serves hover animations.
// Load reads svg/<name>.svg, optional anim/<name>.js and tags.json from an fs.FS.
// Sites copy icons with go tool hgmx icons add <names>.
package glyph

import (
	"bytes"
	_ "embed"
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

// Path is the animation asset prefix; mount without StripPrefix.
const Path = "/icons/"

var (
	//go:embed icons.css
	css []byte
	//go:embed icons.js
	js []byte
	//go:embed motion.js
	motion []byte
)

// Hover is spread onto a hover target to animate its first SVG descendant.
// It requires _hyperscript and also handles mouse presses and touch starts.
var Hover = templ.Attributes{
	"_": "on mouseenter or mousedown or touchstart queue none send hovered to first <svg/> in me then wait 420ms",
}

type glyph struct {
	markup string
	tags   []string
}

type Set struct {
	glyphs map[string]glyph
	names  []string
	anims  map[string][]byte
}

// Load parses every svg/<name>.svg in fsys.
// Only icons with anim/<name>.js emit data-icon; static icons need no scripts.
func Load(fsys fs.FS) (*Set, error) {
	data, err := fs.ReadFile(fsys, "tags.json")
	if err != nil {
		return nil, err
	}
	var tags map[string][]string
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("tags.json: %w", err)
	}
	files, err := fs.Glob(fsys, "svg/*.svg")
	if err != nil {
		return nil, err
	}
	s := &Set{glyphs: map[string]glyph{}, anims: map[string][]byte{}}
	for _, file := range files {
		name := strings.TrimSuffix(path.Base(file), ".svg")
		data, err := fs.ReadFile(fsys, file)
		if err != nil {
			return nil, err
		}
		anim, err := fs.ReadFile(fsys, "anim/"+name+".js")
		animated := err == nil
		if animated {
			s.anims[name] = anim
		}
		markup, err := render(name, data, animated)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		s.glyphs[name] = glyph{markup: markup, tags: tags[name]}
		s.names = append(s.names, name)
	}
	return s, nil
}

// MustLoad is [Load] for package-level sets; it panics on error.
func MustLoad(fsys fs.FS) *Set {
	s, err := Load(fsys)
	if err != nil {
		panic(fmt.Errorf("glyph: %w", err))
	}
	return s
}

func render(name string, data []byte, animated bool) (string, error) {
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
		fmt.Fprintf(&b, ` data-icon="%s"`, html.EscapeString(name))
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

// Ident converts a hyphenated icon name to its Go constant name.
func Ident(name string) string {
	var b strings.Builder
	for part := range strings.SplitSeq(name, "-") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	return b.String()
}

// Render writes the named icon as a 1em SVG.
// An unknown name returns an error without writing anything.
func (s *Set) Render(w io.Writer, name string) error {
	g, ok := s.glyphs[name]
	if !ok {
		return fmt.Errorf("glyph: no icon %q", name)
	}
	_, err := io.WriteString(w, g.markup)
	return err
}

// Names returns a copy of the set's names in sorted order.
func (s *Set) Names() []string {
	return slices.Clone(s.names)
}

// Tags returns a copy of the named icon's search tags.
func (s *Set) Tags(name string) []string {
	return slices.Clone(s.glyphs[name].tags)
}

// Match finds case-insensitive substrings in names and tags.
// An empty or whitespace-only query returns all icons.
func (s *Set) Match(q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	var found []string
	for _, n := range s.names {
		if strings.Contains(n, q) || slices.ContainsFunc(s.glyphs[n].tags, func(tag string) bool {
			return strings.Contains(strings.ToLower(tag), q)
		}) {
			found = append(found, n)
		}
	}
	return found
}

// Handler serves the styles and the set's per-icon animation modules.
type Handler struct {
	*static.Server
	head template.HTML
}

// New serves the shared runtime and only set's animation modules.
func New(set *Set) *Handler {
	files := map[string][]byte{"icons.css": css, "icons.js": js, "motion.js": motion}
	for name, anim := range set.anims {
		files["anim/"+name+".js"] = anim
	}
	h := &Handler{Server: static.Map(Path, files)}
	cssURL, _ := h.URL("icons.css")
	jsURL, _ := h.URL("icons.js")
	h.head = template.HTML(fmt.Sprintf(
		"<link rel=\"stylesheet\" href=\"%s\">\n<script type=\"module\" src=\"%s\"></script>",
		cssURL, jsURL,
	))
	return h
}

// Head returns tags to place after scripts.Motion in the document head.
// Reduced motion skips animation downloads and playback.
func (h *Handler) Head() template.HTML {
	return h.head
}
