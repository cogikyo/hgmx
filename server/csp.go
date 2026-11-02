package server

import (
	"fmt"
	"slices"
	"strings"
)

// CSP replaces each named directive's sources or adds it to hgmx's base.
// A nil map keeps the base; no value disables CSP.
type CSP map[string][]string

type directive struct {
	name string
	base []string
}

var directives = []directive{
	{"default-src", []string{"'none'"}},
	{"script-src", []string{"'self'", "https://static.cloudflareinsights.com"}},
	{"style-src", []string{"'self'"}},
	{"img-src", []string{"'self'"}},
	{"font-src", []string{"'self'"}},
	{"connect-src", []string{"'self'"}},
	{"manifest-src", []string{"'self'"}},
	{"form-action", []string{"'self'"}},
	{"base-uri", []string{"'none'"}},
	{"frame-ancestors", []string{"'none'"}},
	{"script-src-elem", nil},
	{"script-src-attr", nil},
	{"style-src-elem", nil},
	{"style-src-attr", nil},
	{"media-src", nil},
	{"object-src", nil},
	{"child-src", nil},
	{"frame-src", nil},
	{"fenced-frame-src", nil},
	{"worker-src", nil},
	{"webrtc", nil},
	{"sandbox", nil},
	{"upgrade-insecure-requests", nil},
	{"require-trusted-types-for", nil},
	{"trusted-types", nil},
	{"report-to", nil},
	{"report-uri", nil},
}

// Header serializes the base policy with c's overrides.
// It panics on unknown directives, empty sources or unsafe header characters.
func (c CSP) Header() string {
	for name, sources := range c {
		if !slices.ContainsFunc(directives, func(d directive) bool { return d.name == name }) {
			panic(fmt.Sprintf("server: unknown server.CSP directive %q", name))
		}
		for _, s := range sources {
			if s == "" || strings.ContainsFunc(s, unsafe) {
				panic(fmt.Sprintf("server: server.CSP %s source %q must be non-empty visible ASCII without ';' or ','", name, s))
			}
		}
	}
	var parts []string
	for _, d := range directives {
		sources, ok := c[d.name]
		if !ok {
			if d.base == nil {
				continue
			}
			sources = d.base
		}
		parts = append(parts, strings.Join(append([]string{d.name}, sources...), " "))
	}
	return strings.Join(parts, "; ")
}

func unsafe(r rune) bool {
	return r <= ' ' || r > '~' || r == ';' || r == ','
}
