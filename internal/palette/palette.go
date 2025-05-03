package palette

import (
	"fmt"
	"io"
	"strings"

	"github.com/alltom/oklab"
	"github.com/nosvagor/hgmx/views/builder"
)

type ColorScale struct {
	name  string
	shade map[int]*oklab.Oklch
	rl    map[int]float64
	cr    map[int]float64
}

var shades = []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950}

func (cs ColorScale) New(name string) ColorScale {
	cs = ColorScale{
		name:  name,
		shade: make(map[int]*oklab.Oklch),
		rl:    make(map[int]float64),
		cr:    make(map[int]float64),
	}
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

func (p *Palette) ToView() []builder.ColorScaleView {
	var views []builder.ColorScaleView

	// Helper function to convert a ColorScale to ColorScaleView
	convertScale := func(name, code string, scale ColorScale) builder.ColorScaleView {
		view := builder.ColorScaleView{
			Name:   name,
			Code:   code,
			Shades: make([]builder.Shade, len(shades)),
		}

		for i, shade := range shades {
			color := scale.shade[shade]
			rl, cr := OklchCompare(p.seed, *color)
			view.Shades[i] = builder.Shade{
				Code:  code,
				Value: fmt.Sprintf("%d", shade),
				RL:    fmt.Sprintf("%0.4f", rl),
				CR:    fmt.Sprintf("%05.2f", cr),
			}
		}
		return view
	}

	// Base colors
	views = append(views, convertScale("Bg", "bgc", p.Background))
	views = append(views, convertScale("Text", "fgc", p.Foreground))

	// Colors
	views = append(views, convertScale("Ruby", "rby", p.Ruby))
	views = append(views, convertScale("Orange", "orn", p.Orange))
	views = append(views, convertScale("Sun", "sun", p.Sun))
	views = append(views, convertScale("Green", "grn", p.Green))
	views = append(views, convertScale("Emerald", "emr", p.Emerald))
	views = append(views, convertScale("Cyan", "cyn", p.Cyan))
	views = append(views, convertScale("Sky", "sky", p.Sky))
	views = append(views, convertScale("Blue", "blu", p.Blue))
	views = append(views, convertScale("Purple", "prp", p.Purple))
	views = append(views, convertScale("Pink", "pnk", p.Pink))

	// Muted colors
	views = append(views, convertScale("Adenine", "ade", p.Adenine))
	views = append(views, convertScale("Rust", "rst", p.Rust))
	views = append(views, convertScale("Cytosine", "cyt", p.Cytosine))
	views = append(views, convertScale("Olive", "olv", p.Olive))
	views = append(views, convertScale("Forest", "frt", p.Forest))
	views = append(views, convertScale("Slate", "slt", p.Slate))
	views = append(views, convertScale("Thymine", "thy", p.Thymine))
	views = append(views, convertScale("Glacial", "glc", p.Glacial))
	views = append(views, convertScale("Guanine", "gau", p.Guanine))
	views = append(views, convertScale("Plum", "plm", p.Plum))

	// Grays
	views = append(views, convertScale("Black", "blk", p.Black))
	views = append(views, convertScale("Gray", "gry", p.Gray))
	views = append(views, convertScale("White", "wht", p.White))

	return views
}
