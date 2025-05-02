package palette

import (
	"fmt"
	"io"
	"strings"

	"github.com/alltom/oklab"
)

type ColorScale struct {
	name  string
	shade map[int]*oklab.Oklch
}

var shades = []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950}

func (cs ColorScale) New(name string) ColorScale {
	cs = ColorScale{name: name, shade: make(map[int]*oklab.Oklch)}
	for _, shade := range shades {
		cs.shade[shade] = &oklab.Oklch{}
	}
	return cs
}

// Palette holds the complete set of color scales for the theme.
type Palette struct {
	seed oklab.Oklch
	// Base
	Background ColorScale
	Foreground ColorScale
	// Colors
	Ruby    ColorScale
	Orange  ColorScale
	Sun     ColorScale
	Green   ColorScale
	Emerald ColorScale
	Cyan    ColorScale
	Sky     ColorScale
	Blue    ColorScale
	Purple  ColorScale
	Pink    ColorScale
	// Muted
	Adenine  ColorScale
	Rust     ColorScale
	Cytosine ColorScale
	Olive    ColorScale
	Forest   ColorScale
	Slate    ColorScale
	Thymine  ColorScale
	Glacial  ColorScale
	Guanine  ColorScale
	Plum     ColorScale
	// Grays
	Black ColorScale
	Gray  ColorScale
	White ColorScale
}

func Generate(base oklab.Oklch) (p Palette) {
	p.seed = base
	// Base
	p.Background = Background(base)
	p.Foreground = Foreground(base)
	// Colors
	p.Ruby = Ruby(base)
	p.Orange = Orange(base)
	p.Sun = Sun(base)
	p.Green = Green(base)
	p.Emerald = Emerald(base)
	p.Cyan = Cyan(base)
	p.Sky = Sky(base)
	p.Blue = Blue(base)
	p.Purple = Purple(base)
	p.Pink = Pink(base)
	// Muted
	p.Adenine = Adenine(base)
	p.Rust = Rust(base)
	p.Cytosine = Cytosine(base)
	p.Olive = Olive(base)
	p.Forest = Forest(base)
	p.Slate = Slate(base)
	p.Thymine = Thymine(base)
	p.Glacial = Glacial(base)
	p.Guanine = Guanine(base)
	p.Plum = Plum(base)
	// Grays
	p.Black = Black(base)
	p.Gray = Gray(base)
	p.White = White(base)
	return
}

func (p *Palette) ToCSS(w io.Writer) {
	fmt.Fprintln(w, ":root {")
	// Base
	fmt.Fprintln(w, "/* Base */")
	p.Background.ToCSS(w, p.seed)
	p.Foreground.ToCSS(w, p.seed)
	// Colors
	fmt.Fprintln(w, "/* Colors */")
	p.Ruby.ToCSS(w, p.seed)
	p.Orange.ToCSS(w, p.seed)
	p.Sun.ToCSS(w, p.seed)
	p.Green.ToCSS(w, p.seed)
	p.Emerald.ToCSS(w, p.seed)
	p.Cyan.ToCSS(w, p.seed)
	p.Sky.ToCSS(w, p.seed)
	p.Blue.ToCSS(w, p.seed)
	p.Purple.ToCSS(w, p.seed)
	p.Pink.ToCSS(w, p.seed)
	// Muted
	fmt.Fprintln(w, "/* Muted */")
	p.Adenine.ToCSS(w, p.seed)
	p.Rust.ToCSS(w, p.seed)
	p.Cytosine.ToCSS(w, p.seed)
	p.Olive.ToCSS(w, p.seed)
	p.Forest.ToCSS(w, p.seed)
	p.Slate.ToCSS(w, p.seed)
	p.Thymine.ToCSS(w, p.seed)
	p.Glacial.ToCSS(w, p.seed)
	p.Guanine.ToCSS(w, p.seed)
	p.Plum.ToCSS(w, p.seed)
	// Grays
	fmt.Fprintln(w, "/* Grays */")
	p.Black.ToCSS(w, p.seed)
	p.Gray.ToCSS(w, p.seed)
	p.White.ToCSS(w, p.seed)
	fmt.Fprintln(w, "}")
}

func (s ColorScale) ToCSS(w io.Writer, seed oklab.Oklch) {
	for _, shadeKey := range shades {
		color := s.shade[shadeKey]
		css := OklchToString(color)
		rl, cr := OklchCompare(seed, *color)
		if shadeKey == 50 {
			fmt.Fprintf(w, "  --%s-50:  %s; /* RL: %0.4f, CR: %05.2f */\n", strings.ToLower(s.name), css, rl, cr)
		} else {
			fmt.Fprintf(w, "  --%s-%d: %s; /* RL: %0.4f, CR: %05.2f */\n", strings.ToLower(s.name), shadeKey, css, rl, cr)
		}
	}
	fmt.Fprintln(w, " ")
}
