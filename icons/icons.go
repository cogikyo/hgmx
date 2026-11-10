// Package icons is the full catalog copied by go tool hgmx icons.
// Importing it embeds every icon; sites import their own copied set instead.
// Import it only to preview the whole catalog:
//
//	h := glyph.New(icons.Set)
//	root.Handle(glyph.Path, h)
package icons

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cogikyo/hgmx/icons/glyph"
)

//go:generate go run ./fetch

// FS holds the catalog in glyph's set layout.
//
//go:embed svg anim tags.json meta.json
var FS embed.FS

// Set is the whole catalog for [glyph.New].
var Set = glyph.MustLoad(FS)

// Name is a catalog icon that renders as a templ component.
type Name string

func (n Name) Render(_ context.Context, w io.Writer) error {
	return Set.Render(w, string(n))
}

func (n Name) Tags() []string {
	return Set.Tags(string(n))
}

// Source returns where the icon comes from: lucide, lab or hgmx.
func (n Name) Source() string {
	return meta.Icons[n].Source
}

// Aliases returns the names the icon had before upstream renamed it.
func (n Name) Aliases() []string {
	return meta.Icons[n].Aliases
}

type Category struct {
	Name  string
	Title string
	Icons []Name
}

type entry struct {
	Source     string   `json:"source"`
	Categories []string `json:"categories"`
	Aliases    []string `json:"aliases"`
}

type catalogMeta struct {
	Categories []struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"categories"`
	Icons map[Name]entry `json:"icons"`
}

var meta = func() catalogMeta {
	data, err := FS.ReadFile("meta.json")
	if err != nil {
		panic(err)
	}
	m, err := parseMeta(data)
	if err != nil {
		panic(fmt.Errorf("icons: meta.json: %w", err))
	}
	return m
}()

func parseMeta(data []byte) (catalogMeta, error) {
	var m catalogMeta
	err := json.Unmarshal(data, &m)
	return m, err
}

// All returns a sorted copy of the catalog's icon names.
func All() []Name {
	return names(Set.Names())
}

// Match finds case-insensitive substrings in names and tags.
// An empty or whitespace-only query returns all icons.
func Match(q string) []Name {
	return names(Set.Match(q))
}

// Categories returns the catalog's categories in display order, each with its icons.
func Categories() []Category {
	return meta.categories(All())
}

func (m catalogMeta) categories(all []Name) []Category {
	out := make([]Category, len(m.Categories))
	index := map[string]int{}
	for i, c := range m.Categories {
		out[i] = Category{Name: c.Name, Title: c.Title}
		index[c.Name] = i
	}
	for _, n := range all {
		for _, c := range m.Icons[n].Categories {
			if i, ok := index[c]; ok {
				out[i].Icons = append(out[i].Icons, n)
			}
		}
	}
	return out
}

func names(in []string) []Name {
	out := make([]Name, len(in))
	for i, n := range in {
		out[i] = Name(n)
	}
	return out
}
