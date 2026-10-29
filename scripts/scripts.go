// Package scripts serves selected vendored scripts with ordered HTML tags.
package scripts

import (
	"embed"
	"fmt"
	"html/template"
	"slices"
	"strings"

	"github.com/cogikyo/hgmx/static"
)

//go:generate go run ./fetch

//go:embed vendor
var vendor embed.FS

// Path is the mount prefix; keep it in the request path.
const Path = "/scripts/"

// Library selects a script from the vendored set.
type Library int

const (
	// HTMX selects the htmx runtime.
	HTMX Library = iota
	// Preload selects hx-preload and includes [HTMX] automatically.
	Preload
	// Hyperscript selects the _hyperscript runtime.
	Hyperscript
	// Motion provides the global Motion object, not ES module exports.
	Motion
)

var files = [...]string{
	HTMX:        "htmx.min.js",
	Preload:     "hx-preload.js",
	Hyperscript: "hyperscript.min.js",
	Motion:      "motion.min.js",
}

// Handler serves the selected scripts and their matching HTML tags.
// Mount it at [Path] without http.StripPrefix.
type Handler struct {
	*static.Server
	head template.HTML
}

// New adds dependencies and removes duplicate selections.
// Scripts load in [Library] declaration order, not argument order.
// Undeclared Library values panic.
func New(libs ...Library) (*Handler, error) {
	libs = slices.Clone(libs)
	if slices.Contains(libs, Preload) {
		libs = append(libs, HTMX)
	}
	slices.Sort(libs)
	libs = slices.Compact(libs)
	selected := make(map[string][]byte, len(libs))
	for _, lib := range libs {
		data, err := vendor.ReadFile("vendor/" + files[lib])
		if err != nil {
			return nil, fmt.Errorf("scripts: %w", err)
		}
		selected[files[lib]] = data
	}

	h := &Handler{Server: static.Map(Path, selected)}
	tags := make([]string, len(libs))
	for i, lib := range libs {
		url, err := h.URL(files[lib])
		if err != nil {
			return nil, err
		}
		tags[i] = fmt.Sprintf("<script src=\"%s\" defer></script>", url)
	}
	h.head = template.HTML(strings.Join(tags, "\n"))
	return h, nil
}

// Head returns deferred script tags for the document head.
func (h *Handler) Head() template.HTML {
	return h.head
}
