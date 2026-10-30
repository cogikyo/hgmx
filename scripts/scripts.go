// Package scripts serves selected vendored and first-party scripts with ordered HTML tags.
package scripts

import (
	_ "embed"
	"fmt"
	"html/template"
	"slices"
	"strings"

	"github.com/cogikyo/hgmx/static"
)

//go:generate go run ./fetch

var (
	//go:embed vendor/htmx.min.js
	htmx []byte
	//go:embed vendor/hx-preload.js
	preload []byte
	//go:embed hgmx.js
	hgmx []byte
	//go:embed vendor/hyperscript.min.js
	hyperscript []byte
	//go:embed vendor/motion.min.js
	motion []byte
)

// Path is the mount prefix; keep it in the request path.
const Path = "/scripts/"

// Library selects a script from the embedded set.
type Library int

const (
	// HTMX selects the htmx runtime.
	HTMX Library = iota
	// Preload selects hx-preload and includes [HTMX] automatically.
	Preload
	// Navigation includes [HTMX] and [Preload] with newest-navigation-wins swaps.
	// Late responses from removed elements or aborted requests are dropped.
	// A newer boosted navigation or history restore drops older navigation responses.
	// Other requests are unaffected by newer navigations.
	// History restores re-send the element's hx-config headers.
	// After either swap, focus moves to main's h1 or main unless already inside main.
	// Do not add hx-sync to body when loading Navigation.
	Navigation
	// Hyperscript selects the _hyperscript runtime.
	Hyperscript
	// Motion provides the global Motion object, not ES module exports.
	Motion
)

var files = [...]struct {
	name string
	data []byte
}{
	HTMX:        {"htmx.min.js", htmx},
	Preload:     {"hx-preload.js", preload},
	Navigation:  {"hgmx.js", hgmx},
	Hyperscript: {"hyperscript.min.js", hyperscript},
	Motion:      {"motion.min.js", motion},
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
func New(libs ...Library) *Handler {
	libs = slices.Clone(libs)
	if slices.Contains(libs, Navigation) {
		libs = append(libs, Preload)
	}
	if slices.Contains(libs, Preload) {
		libs = append(libs, HTMX)
	}
	slices.Sort(libs)
	libs = slices.Compact(libs)
	selected := make(map[string][]byte, len(libs))
	for _, lib := range libs {
		selected[files[lib].name] = files[lib].data
	}

	h := &Handler{Server: static.Map(Path, selected)}
	tags := make([]string, len(libs))
	for i, lib := range libs {
		url, _ := h.URL(files[lib].name)
		tags[i] = fmt.Sprintf("<script src=\"%s\" defer></script>", url)
	}
	h.head = template.HTML(strings.Join(tags, "\n"))
	return h
}

// Head returns deferred script tags for the document head.
func (h *Handler) Head() template.HTML {
	return h.head
}
