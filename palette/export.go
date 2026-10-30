package palette

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Token names a CSS custom property without the -- prefix.
type Token struct {
	Name  string
	Color Color
	Ref   string // Referenced token name, or empty for a literal colour.
}

func (p Palette) Tokens() []Token {
	out := append([]Token(nil), p.Roles...)
	for _, s := range p.Scales {
		ref := func(name string) string {
			if s.Alias == "" {
				return ""
			}
			return strings.Replace(name, s.Name, s.Alias, 1)
		}
		fill := Token{Name: s.Name, Color: s.Fill, Ref: ref(s.Name)}
		if !s.Seeded && fill.Ref == "" {
			fill.Ref = s.Token(FillStep)
		}
		out = append(out, fill, Token{Name: "on-" + s.Name, Color: s.On, Ref: ref("on-" + s.Name)})
		for i, c := range s.Steps {
			out = append(out, Token{Name: s.Token(i), Color: c, Ref: ref(s.Token(i))})
			if s.Beta != nil {
				out = append(out, Token{Name: s.BetaToken(i), Color: s.Beta[i], Ref: ref(s.BetaToken(i))})
			}
		}
	}
	return out
}

func validName(s string) bool {
	if s == "" || len(s) > 64 {
		return false
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

// Selector validates name before building a data-theme CSS selector.
// Names must contain 1 to 64 characters from a-z, 0-9, and -.
func Selector(name string) (string, error) {
	if !validName(name) {
		return "", fmt.Errorf("palette: theme name %q must be 1 to 64 characters of a-z, 0-9 and -", name)
	}
	return `[data-theme="` + name + `"]`, nil
}

func (p Palette) Scheme() string {
	if p.Dark {
		return "dark"
	}
	return "light"
}

// CSS writes custom properties under selector, which is used verbatim.
// Use [Selector] to validate an untrusted theme name.
func (p Palette) CSS(w io.Writer, selector string) error {
	b := bufio.NewWriter(w)
	fmt.Fprintf(b, "%s {\n\tcolor-scheme: %s;\n", selector, p.Scheme())
	for _, t := range p.Tokens() {
		if t.Ref != "" {
			fmt.Fprintf(b, "\t--%s: var(--%s);\n", t.Name, t.Ref)
		} else {
			fmt.Fprintf(b, "\t--%s: %s;\n", t.Name, t.Color)
		}
	}
	b.WriteString("}\n")
	return b.Flush()
}

// Theme writes Tailwind's @theme inline mapping to the palette's CSS variables.
// Load the output of [Palette.CSS] too; Theme does not define their values.
func (p Palette) Theme(w io.Writer) error {
	b := bufio.NewWriter(w)
	b.WriteString("@theme inline {\n\t--color-*: initial;\n")
	for _, t := range p.Tokens() {
		fmt.Fprintf(b, "\t--color-%s: var(--%s);\n", t.Name, t.Name)
	}
	b.WriteString("}\n")
	return b.Flush()
}

// Lua writes a Lua table of literal token colours and the colour scheme.
func (p Palette) Lua(w io.Writer) error {
	b := bufio.NewWriter(w)
	fmt.Fprintf(b, "return {\n\tscheme = %q,\n", p.Scheme())
	group := ""
	for _, t := range p.Tokens() {
		g, _, _ := strings.Cut(strings.TrimPrefix(t.Name, "on-"), "-")
		if _, ok := p.Scale(g); !ok {
			g = "roles"
		}
		if g != group {
			fmt.Fprintf(b, "\n\t-- %s\n", g)
			group = g
		}
		fmt.Fprintf(b, "\t%s = %q, -- CR: %.2f\n", strings.ReplaceAll(t.Name, "-", "_"), t.Color.Hex(), Ratio(t.Color, p.Seed.Base))
	}
	b.WriteString("}\n")
	return b.Flush()
}
