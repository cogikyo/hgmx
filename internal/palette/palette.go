package palette

import (
	"fmt"
	"io"
	"log"
	"math"

	"github.com/alltom/oklab"
	"github.com/nosvagor/hgmx/views/builder"
)

func Generate(seed string) Palette {
	p := make(Palette)
	for code, color := range colors {
		if code == "bgc" {
			color.Seed = Hex(seed)
		}
		code.valid()
		oklch := color.Seed.toOklch()
		details := &ColorDetails{Color: color, Base: *oklch}
		p[code] = details
		details.generate()
	}
	return p
}

// === Models ==================================================================

var shadeValues = []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950}

var orderedCodes = []Code{
	"bgc", "fgc", "tes",
	"rse", "bry", "chy", "rby", "red", "crl", "pmk", "orn", "sun", "gld", "hny", "yel",
	"lem", "acd", "lim", "spr", "grn", "emr", "jde", "frs", "lea", "tea", "cyn",
	"aqu", "rbn", "azr", "sky", "blu", "cbt", "sph", "ind", "lav", "prp", "vio", "pnk", "mag",
	"blk", "gry", "wht",
}

type Code string

func (c Code) valid() {
	if len(c) != 3 {
		panic(fmt.Sprintf("invalid code: %s", c))
	}
}

type Hex string

func (h Hex) toOklch() *oklab.Oklch {
	c, err := HexToOklch(string(h))
	if err != nil {
		panic(err)
	}
	return &c
}

type Colors map[Code]Color

type Color struct {
	Seed Hex
	Name string
}

type Details struct {
	oklab.Oklch
	RL float64
	CR float64
}

type ColorDetails struct {
	Color
	Base   oklab.Oklch
	Shades map[int]Details
}

type Palette map[Code]*ColorDetails

var colors = Colors{
	"bgc": Color{"", "Background"},
	"fgc": Color{"#aeb9f8", "Foreground"},
	// "tes": Color{"#000000", "Test"},
	"rse": Color{"#fc0086", "Rose"},
	"bry": Color{"#fd016f", "Berry"},
	"chy": Color{"#ff0457", "Cherry"},
	"rby": Color{"#f9043a", "Ruby"},
	"red": Color{"#fd181a", "Red"},
	"crl": Color{"#fb3d03", "Coral"},
	"pmk": Color{"#fd5802", "Pumpkin"},
	"orn": Color{"#ff7220", "Orange"},
	"sun": Color{"#ff9004", "Sun"},
	"gld": Color{"#fead05", "Gold"},
	"hny": Color{"#ffcc00", "Honey"},
	"yel": Color{"#fddf00", "Yellow"},
	"lem": Color{"#ecec00", "Lemon"},
	"acd": Color{"#cdf118", "Acid"},
	"lim": Color{"#aae801", "Lime"},
	"spr": Color{"#86e401", "Spring"},
	"grn": Color{"#58d300", "Green"},
	"emr": Color{"#28c624", "Emerald"},
	"jde": Color{"#01b947", "Jade"},
	"frs": Color{"#03bb65", "Forest"},
	"lea": Color{"#01c37e", "Leaf"},
	"tea": Color{"#0ed39a", "Teal"},
	"cyn": Color{"#00e7cb", "Cyan"},
	"aqu": Color{"#02eeef", "Aqua"},
	"rbn": Color{"#07e3fe", "Robin"},
	"azr": Color{"#0acbff", "Azure"},
	"sky": Color{"#0aafff", "Sky"},
	"blu": Color{"#0184fe", "Blue"},
	"cbt": Color{"#256eff", "Cobalt"},
	"sph": Color{"#4158fa", "Sapphire"},
	"ind": Color{"#5a4aff", "Indigo"},
	"lav": Color{"#6e40ff", "Lavender"},
	"prp": Color{"#972eff", "Purple"},
	"vio": Color{"#c602fe", "Violet"},
	"pnk": Color{"#ea0aeb", "Pink"},
	"mag": Color{"#fd01b9", "Magenta"},
	// "blk": Color{"#1d1d21", "Black"},
	// "gry": Color{"#4f5163", "Gray"},
	// "wht": Color{"#ddddf6", "White"},
}

// generate dynamically creates all shades from 50-950.
// It uses linear interpolation from the Base color towards hardcoded targets for 50 and 950,
// with special final-step adjustments for shades 50 and 950.
func (c *ColorDetails) generate() {
	c.Shades = make(map[int]Details, len(shadeValues))
	for _, shade := range shadeValues {
		c.Shades[shade] = Details{Oklch: oklab.Oklch{}}
	}
	if _, ok := c.Shades[600]; !ok {
		panic("Shade 600 must be in shadeValues for generation")
	}
	c.Shades[600] = Details{Oklch: c.Base}
	hue := c.Base.H

	shadesMap := make(map[int]int, len(shadeValues))
	for i, v := range shadeValues {
		shadesMap[v] = i
	}

	// --- Constants for Light Shades (50-500) ---
	const targetL50 = 0.97
	const targetC50 = 0.01
	const adjustL50 = 0.25
	const adjustC50 = 0.37
	numIntervalsLight := float64(shadesMap[600] - shadesMap[50])

	// First, calculate all light shades (50-500) with linear interpolation
	for shadeValue, index := range shadesMap {
		if shadeValue > 500 {
			continue
		}
		t := float64(index) / numIntervalsLight
		l := targetL50 + (c.Base.L-targetL50)*t
		chroma := max(targetC50+(c.Base.C-targetC50)*t, 0)
		detail := c.Shades[shadeValue]
		detail.Oklch = oklab.Oklch{L: max(0, min(l, 1)), C: chroma, H: hue}
		c.Shades[shadeValue] = detail
	}

	// Special adjustment for shade 50's lightness and chroma
	details100, ok100 := c.Shades[100]
	details50, ok50 := c.Shades[50]
	if ok100 && ok50 {
		l100_linear := details100.L
		l50_initial_linear := details50.L
		l50_adjusted := l50_initial_linear + (l100_linear-l50_initial_linear)*adjustL50
		details50.L = max(0, min(l50_adjusted, 1))

		c100_linear := details100.C
		c50_initial_linear := details50.C
		c50_adjusted := c50_initial_linear + (c100_linear-c50_initial_linear)*adjustC50
		details50.C = max(c50_adjusted, 0)
		c.Shades[50] = details50
	}

	// --- Constants for Dark Shades (700-950) ---
	const targetL950 = 0.25
	const targetC950 = 0.05
	const adjustL950 = 0.37
	const adjustC950 = 0.42
	numIntervalsDark := float64(shadesMap[950] - shadesMap[600])
	baseIndexDark := float64(shadesMap[600])

	// Calculate all dark shades (700-950) with linear interpolation
	for shadeValue, index := range shadesMap {
		if shadeValue < 700 {
			continue
		}
		t := (float64(index) - baseIndexDark) / numIntervalsDark
		l := c.Base.L + (targetL950-c.Base.L)*t
		chroma := c.Base.C + (targetC950-c.Base.C)*t
		detail := c.Shades[shadeValue]
		detail.Oklch = oklab.Oklch{L: max(0, min(l, 1)), C: max(chroma, 0), H: hue}
		c.Shades[shadeValue] = detail
	}

	// Special adjustment for shade 950's lightness and chroma
	details900, ok900 := c.Shades[900]
	details950, ok950 := c.Shades[950]
	if ok900 && ok950 {
		l900 := details900.L
		l950_initial := details950.L
		l950_adjusted := l950_initial + (l900-l950_initial)*adjustL950
		details950.L = max(0, min(l950_adjusted, 1))

		c900 := details900.C
		c950_initial := details950.C
		c950_adjusted := c950_initial + (c900-c950_initial)*adjustC950
		details950.C = max(c950_adjusted, 0)
		c.Shades[950] = details950
	}
}

func (p *Palette) ToCSS(w io.Writer) {
	seed := (*p)["bgc"].Base
	fmt.Fprintln(w, ":root {")
	// Iterate using the defined order
	for _, code := range orderedCodes {
		colorDetails, ok := (*p)[code]
		if ok {
			colorDetails.ToCSS(w, code, seed)
		}
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "@theme {")
	// Iterate using the defined order
	for _, code := range orderedCodes {
		colorDetails, ok := (*p)[code]
		if ok {
			colorDetails.ToTheme(w, code, seed)
		}
	}
	fmt.Fprintln(w, "}")
}

func (c ColorDetails) ToCSS(w io.Writer, code Code, seed oklab.Oklch) {
	fmt.Fprintln(w, "\t/*", c.Name, "*/")
	// Iterate using the ordered shadeValues slice
	for _, shadeKey := range shadeValues {
		color, ok := c.Shades[shadeKey]
		if !ok {
			continue
		} // Skip if a shade is missing (shouldn't happen)
		css := OklchToString(&color.Oklch)
		fmt.Fprintf(w, "  --%s-%d: %s;\n", string(code), shadeKey, css)
	}
	fmt.Fprintln(w, " ")
}

func (c ColorDetails) ToTheme(w io.Writer, code Code, seed oklab.Oklch) {
	// Iterate using the ordered shadeValues slice
	for _, shadeKey := range shadeValues {
		color, ok := c.Shades[shadeKey]
		if !ok {
			continue
		} // Skip if a shade is missing
		css := OklchToString(&color.Oklch)
		fmt.Fprintf(w, "  --color-%s-%d: %s;\n", string(code), shadeKey, css)
	}
	fmt.Fprintln(w, " ")
}

func (p Palette) ToView() []builder.ColorScaleView {
	var views []builder.ColorScaleView

	bgcDetails, ok := p["bgc"]
	if !ok {
		log.Println("Warning: Background color 'bgc' not found in palette for ToView seed.")
		return views
	}
	seed := bgcDetails.Base

	convertScale := func(code Code, scale *ColorDetails) builder.ColorScaleView {
		view := builder.ColorScaleView{
			Name:   scale.Name,
			Code:   string(code),
			Value:  scale.Shades[600].Oklch,
			Shades: make([]builder.Shade, len(shadeValues)),
		}

		for i, shadeVal := range shadeValues {
			colorDetails, ok := scale.Shades[shadeVal]
			if !ok {
				continue
			}

			color := colorDetails.Oklch
			rl, cr := OklchCompare(seed, color)
			hue := toDegree(color.H)

			scaledRadiusC := 37.0 * math.Tanh(6.0*color.C)
			totalDistanceC := scaledRadiusC

			scaledRadiusL := 37.0 * math.Pow(color.L, 1.5)
			totalDistanceL := scaledRadiusL

			angle := -color.H
			view.Shades[i] = builder.Shade{
				Code:  string(code),
				Value: shadeVal,
				RL:    rl,
				CR:    cr,
				L:     fmt.Sprintf("%0.1f%%", color.L*100),
				C:     fmt.Sprintf("%0.2f", color.C),
				H:     fmt.Sprintf("%0.1f", hue),
				Hex:   OklchToHex(&color),
				Cx:    50.0 + totalDistanceC*math.Cos(angle),
				Cy:    50.0 + totalDistanceC*math.Sin(angle),
				Clx:   50.0 + totalDistanceL*math.Cos(angle),
				Cly:   50.0 + totalDistanceL*math.Sin(angle),
			}
		}
		return view
	}

	for _, code := range orderedCodes {
		if code == "bgc" || code == "fgc" {
			continue
		}
		details, ok := p[code]
		if ok {
			views = append(views, convertScale(code, details))
		}
	}

	return views
}
