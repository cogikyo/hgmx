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
	Berry    ColorScale
	Cherry   ColorScale
	Ruby     ColorScale
	Coral    ColorScale
	Orange   ColorScale
	Pumpkin  ColorScale
	Sun      ColorScale
	Gold     ColorScale
	Yellow   ColorScale
	Lemon    ColorScale
	Lime     ColorScale
	Acid     ColorScale
	Kiwi     ColorScale
	Teal     ColorScale
	Green    ColorScale
	Spring   ColorScale
	Emerald  ColorScale
	Jade     ColorScale
	Forest   ColorScale
	Leaf     ColorScale
	Cyan     ColorScale
	Robin    ColorScale
	Azure    ColorScale
	Aqua     ColorScale
	Sky      ColorScale
	Cobalt   ColorScale
	Sapphire ColorScale
	Indigo   ColorScale
	Blue     ColorScale
	Lavender ColorScale
	Purple   ColorScale
	Violet   ColorScale
	Magenta  ColorScale
	Rose     ColorScale
	Pink     ColorScale
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
	p.Berry = Berry(base)
	p.Cherry = Cherry(base)
	p.Ruby = Ruby(base)
	p.Coral = Coral(base)
	p.Orange = Orange(base)
	p.Pumpkin = Pumpkin(base)
	p.Sun = Sun(base)
	p.Gold = Gold(base)
	p.Yellow = Yellow(base)
	p.Lemon = Lemon(base)
	p.Lime = Lime(base)
	p.Acid = Acid(base)
	p.Kiwi = Kiwi(base)
	p.Green = Green(base)
	p.Spring = Spring(base)
	p.Emerald = Emerald(base)
	p.Jade = Jade(base)
	p.Forest = Forest(base)
	p.Leaf = Leaf(base)
	p.Spring = Spring(base)
	p.Teal = Teal(base)
	p.Aqua = Aqua(base)
	p.Cyan = Cyan(base)
	p.Robin = Robin(base)
	p.Sky = Sky(base)
	p.Azure = Azure(base)
	p.Cobalt = Cobalt(base)
	p.Sapphire = Sapphire(base)
	p.Indigo = Indigo(base)
	p.Blue = Blue(base)
	p.Lavender = Lavender(base)
	p.Purple = Purple(base)
	p.Violet = Violet(base)
	p.Magenta = Magenta(base)
	p.Rose = Rose(base)
	p.Pink = Pink(base)
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
	p.Berry.ToCSS(w, p.seed)
	p.Cherry.ToCSS(w, p.seed)
	p.Red.ToCSS(w, p.seed)
	p.Ruby.ToCSS(w, p.seed)
	p.Coral.ToCSS(w, p.seed)
	p.Pumpkin.ToCSS(w, p.seed)
	p.Orange.ToCSS(w, p.seed)
	p.Sun.ToCSS(w, p.seed)
	p.Gold.ToCSS(w, p.seed)
	p.Yellow.ToCSS(w, p.seed)
	p.Lemon.ToCSS(w, p.seed)
	p.Lime.ToCSS(w, p.seed)
	p.Acid.ToCSS(w, p.seed)
	p.Kiwi.ToCSS(w, p.seed)
	p.Green.ToCSS(w, p.seed)
	p.Spring.ToCSS(w, p.seed)
	p.Emerald.ToCSS(w, p.seed)
	p.Jade.ToCSS(w, p.seed)
	p.Forest.ToCSS(w, p.seed)
	p.Leaf.ToCSS(w, p.seed)
	p.Spring.ToCSS(w, p.seed)
	p.Teal.ToCSS(w, p.seed)
	p.Aqua.ToCSS(w, p.seed)
	p.Cyan.ToCSS(w, p.seed)
	p.Robin.ToCSS(w, p.seed)
	p.Sky.ToCSS(w, p.seed)
	p.Azure.ToCSS(w, p.seed)
	p.Cobalt.ToCSS(w, p.seed)
	p.Sapphire.ToCSS(w, p.seed)
	p.Indigo.ToCSS(w, p.seed)
	p.Blue.ToCSS(w, p.seed)
	p.Lavender.ToCSS(w, p.seed)
	p.Purple.ToCSS(w, p.seed)
	p.Violet.ToCSS(w, p.seed)
	p.Magenta.ToCSS(w, p.seed)
	p.Rose.ToCSS(w, p.seed)
	p.Pink.ToCSS(w, p.seed)
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
			Value:  *scale.shade[600],
			Shades: make([]builder.Shade, len(shades)),
		}

		for i, shade := range shades {
			color := scale.shade[shade]
			rl, cr := OklchCompare(p.seed, *color)
			hue := toDegree(color.H)

			scaledRadiusC := 37.0 * math.Tanh(6.0*color.C)
			totalDistanceC := scaledRadiusC

			scaledRadiusL := 37.0 * math.Pow(color.L, 1.5)
			totalDistanceL := scaledRadiusL

			angle := -color.H
			view.Shades[i] = builder.Shade{
				Code:  code,
				Value: fmt.Sprintf("%d", shade),
				RL:    rl,
				CR:    cr,
				L:     fmt.Sprintf("%0.1f%%", color.L*100),
				C:     fmt.Sprintf("%0.2f", color.C),
				H:     fmt.Sprintf("%0.1f", hue),
				Hex:   OklchToHex(color),
				Cx:    50.0 + totalDistanceC*math.Cos(angle),
				Cy:    50.0 + totalDistanceC*math.Sin(angle),
				Clx:   50.0 + totalDistanceL*math.Cos(angle),
				Cly:   50.0 + totalDistanceL*math.Sin(angle),
			}
		}
		return view
	}

	// Base colors
	// views = append(views, convertScale("Bg", "bgc", p.Background))
	// views = append(views, convertScale("Text", "fgc", p.Foreground))

	// Colors
	views = append(views, convertScale("Rose", "ros", p.Rose))
	views = append(views, convertScale("Berry", "bry", p.Berry))
	views = append(views, convertScale("Cherry", "chy", p.Cherry))
	views = append(views, convertScale("Ruby", "rby", p.Ruby))
	views = append(views, convertScale("Red", "red", p.Red))
	views = append(views, convertScale("Coral", "crl", p.Coral))
	views = append(views, convertScale("Pumpkin", "pmk", p.Pumpkin))
	views = append(views, convertScale("Orange", "orn", p.Orange))
	views = append(views, convertScale("Sun", "sun", p.Sun))
	views = append(views, convertScale("Gold", "gld", p.Gold))
	views = append(views, convertScale("Yellow", "yel", p.Yellow))
	views = append(views, convertScale("Lemon", "lem", p.Lemon))
	views = append(views, convertScale("Acid", "acd", p.Acid))
	views = append(views, convertScale("Lime", "lim", p.Lime))
	views = append(views, convertScale("Kiwi", "kwi", p.Kiwi))
	views = append(views, convertScale("Green", "grn", p.Green))
	views = append(views, convertScale("Spring", "spr", p.Spring))
	views = append(views, convertScale("Emerald", "emr", p.Emerald))
	views = append(views, convertScale("Jade", "jde", p.Jade))
	views = append(views, convertScale("Forest", "frs", p.Forest))
	views = append(views, convertScale("Leaf", "lea", p.Leaf))
	views = append(views, convertScale("Teal", "tea", p.Teal))
	views = append(views, convertScale("Cyan", "cyn", p.Cyan))
	views = append(views, convertScale("Aqua", "aqu", p.Aqua))
	views = append(views, convertScale("Robin", "rbn", p.Robin))
	views = append(views, convertScale("Azure", "azr", p.Azure))
	views = append(views, convertScale("Sky", "sky", p.Sky))
	views = append(views, convertScale("Blue", "blu", p.Blue))
	views = append(views, convertScale("Cobalt", "cbt", p.Cobalt))
	views = append(views, convertScale("Sapphire", "sph", p.Sapphire))
	views = append(views, convertScale("Indigo", "ind", p.Indigo))
	views = append(views, convertScale("Lavender", "lav", p.Lavender))
	views = append(views, convertScale("Purple", "prp", p.Purple))
	views = append(views, convertScale("Violet", "vio", p.Violet))
	views = append(views, convertScale("Pink", "pnk", p.Pink))
	views = append(views, convertScale("Magenta", "mag", p.Magenta))

	return views
}
