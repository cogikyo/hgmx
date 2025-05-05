package palette

import (
	"fmt"
	"io"
	"math"
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
	Red      ColorScale
	Ruby     ColorScale
	Orange   ColorScale
	Sun      ColorScale
	Gold     ColorScale
	Yellow   ColorScale
	Lemon    ColorScale
	Lime     ColorScale
	Teal     ColorScale
	Green    ColorScale
	Emerald  ColorScale
	Cyan     ColorScale
	Cerulean ColorScale
	Azure    ColorScale
	Aqua     ColorScale
	Sky      ColorScale
	Sapphire ColorScale
	Blue     ColorScale
	Lavender ColorScale
	Purple   ColorScale
	Violet   ColorScale
	Magenta  ColorScale
	Rose     ColorScale
	Pink     ColorScale
	// Muted
	Adenine  ColorScale
	Rust     ColorScale
	Cytosine ColorScale
	Olive    ColorScale
	Forest   ColorScale
	Slate    ColorScale
	Thymine  ColorScale
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
	p.Red = Red(base)
	p.Ruby = Ruby(base)
	p.Orange = Orange(base)
	p.Sun = Sun(base)
	p.Gold = Gold(base)
	p.Yellow = Yellow(base)
	p.Lemon = Lemon(base)
	p.Lime = Lime(base)
	p.Green = Green(base)
	p.Emerald = Emerald(base)
	p.Teal = Teal(base)
	p.Aqua = Aqua(base)
	p.Cyan = Cyan(base)
	p.Cerulean = Cerulean(base)
	p.Sky = Sky(base)
	p.Azure = Azure(base)
	p.Sapphire = Sapphire(base)
	p.Blue = Blue(base)
	p.Lavender = Lavender(base)
	p.Magenta = Magenta(base)
	p.Violet = Violet(base)
	p.Purple = Purple(base)
	p.Rose = Rose(base)
	p.Pink = Pink(base)
	// Muted
	p.Adenine = Adenine(base)
	p.Rust = Rust(base)
	p.Cytosine = Cytosine(base)
	p.Olive = Olive(base)
	p.Forest = Forest(base)
	p.Slate = Slate(base)
	p.Thymine = Thymine(base)
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
	p.Background.ToCSS(w, p.seed)
	p.Foreground.ToCSS(w, p.seed)
	// Colors
	p.Red.ToCSS(w, p.seed)
	p.Ruby.ToCSS(w, p.seed)
	p.Orange.ToCSS(w, p.seed)
	p.Sun.ToCSS(w, p.seed)
	p.Gold.ToCSS(w, p.seed)
	p.Yellow.ToCSS(w, p.seed)
	p.Lemon.ToCSS(w, p.seed)
	p.Lime.ToCSS(w, p.seed)
	p.Aqua.ToCSS(w, p.seed)
	p.Azure.ToCSS(w, p.seed)
	p.Cerulean.ToCSS(w, p.seed)
	p.Green.ToCSS(w, p.seed)
	p.Emerald.ToCSS(w, p.seed)
	p.Teal.ToCSS(w, p.seed)
	p.Cyan.ToCSS(w, p.seed)
	p.Sky.ToCSS(w, p.seed)
	p.Sapphire.ToCSS(w, p.seed)
	p.Blue.ToCSS(w, p.seed)
	p.Lavender.ToCSS(w, p.seed)
	p.Purple.ToCSS(w, p.seed)
	p.Violet.ToCSS(w, p.seed)
	p.Magenta.ToCSS(w, p.seed)
	p.Rose.ToCSS(w, p.seed)
	p.Pink.ToCSS(w, p.seed)
	// Muted
	p.Adenine.ToCSS(w, p.seed)
	p.Rust.ToCSS(w, p.seed)
	p.Cytosine.ToCSS(w, p.seed)
	p.Olive.ToCSS(w, p.seed)
	p.Forest.ToCSS(w, p.seed)
	p.Slate.ToCSS(w, p.seed)
	p.Thymine.ToCSS(w, p.seed)
	p.Guanine.ToCSS(w, p.seed)
	p.Plum.ToCSS(w, p.seed)
	// Grays
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

	convertScale := func(name, code string, scale ColorScale) builder.ColorScaleView {
		view := builder.ColorScaleView{
			Name:   name,
			Code:   code,
			Value:  *scale.shade[500],
			Shades: make([]builder.Shade, len(shades)),
		}

		for i, shade := range shades {
			color := scale.shade[shade]
			rl, cr := OklchCompare(p.seed, *color)
			hue := toDegree(color.H)

			// Calculate radius using tanh for non-linear scaling
			maxVisualRadius := 40.0
			spreadFactor := 10.0
			scaledRadius := maxVisualRadius * math.Tanh(spreadFactor*color.C)
			totalDistance := scaledRadius

			angle := -color.H // User reverted rotation offset
			view.Shades[i] = builder.Shade{
				Code:  code,
				Value: fmt.Sprintf("%d", shade),
				RL:    rl,
				CR:    cr,
				L:     fmt.Sprintf("%0.1f%%", color.L*100),
				C:     fmt.Sprintf("%0.2f", color.C),
				H:     fmt.Sprintf("%0.1f", hue),
				Hex:   OklchToHex(color),
				Cx:    50.0 + totalDistance*math.Cos(angle),
				Cy:    50.0 + totalDistance*math.Sin(angle),
			}
		}
		return view
	}

	// Base colors
	// views = append(views, convertScale("Bg", "bgc", p.Background))
	// views = append(views, convertScale("Text", "fgc", p.Foreground))

	// Colors
	views = append(views, convertScale("Ruby", "rby", p.Ruby))
	views = append(views, convertScale("Red", "red", p.Red))
	views = append(views, convertScale("Orange", "orn", p.Orange))
	views = append(views, convertScale("Sun", "sun", p.Sun))
	views = append(views, convertScale("Gold", "gld", p.Gold))
	views = append(views, convertScale("Yellow", "yel", p.Yellow))
	views = append(views, convertScale("Lemon", "lem", p.Lemon))
	views = append(views, convertScale("Lime", "lim", p.Lime))
	views = append(views, convertScale("Green", "grn", p.Green))
	views = append(views, convertScale("Emerald", "emr", p.Emerald))
	views = append(views, convertScale("Teal", "tea", p.Teal))
	views = append(views, convertScale("Cyan", "cyn", p.Cyan))
	views = append(views, convertScale("Aqua", "aqu", p.Aqua))
	views = append(views, convertScale("Cerulean", "cer", p.Cerulean))
	views = append(views, convertScale("Azure", "azr", p.Azure))
	views = append(views, convertScale("Sky", "sky", p.Sky))
	views = append(views, convertScale("Blue", "blu", p.Blue))
	views = append(views, convertScale("Sapphire", "sap", p.Sapphire))
	views = append(views, convertScale("Lavender", "lav", p.Lavender))
	views = append(views, convertScale("Purple", "prp", p.Purple))
	views = append(views, convertScale("Violet", "vio", p.Violet))
	views = append(views, convertScale("Pink", "pnk", p.Pink))
	views = append(views, convertScale("Magenta", "mag", p.Magenta))
	views = append(views, convertScale("Rose", "ros", p.Rose))

	// Muted colors
	// views = append(views, convertScale("Forest", "frt", p.Forest))
	// views = append(views, convertScale("Cytosine", "cyt", p.Cytosine))
	// views = append(views, convertScale("Adenine", "ade", p.Adenine))
	// views = append(views, convertScale("Guanine", "gau", p.Guanine))
	// views = append(views, convertScale("Thymine", "thy", p.Thymine))
	// views = append(views, convertScale("Plum", "plm", p.Plum))
	// views = append(views, convertScale("Olive", "olv", p.Olive))
	// views = append(views, convertScale("Slate", "slt", p.Slate))
	// views = append(views, convertScale("Rust", "rst", p.Rust))
	return views
}
