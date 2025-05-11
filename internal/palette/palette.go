package palette

import (
	"fmt"
	"io"
	"log"
	"math"

	"github.com/alltom/oklab"
	colorsPage "github.com/nosvagor/hgmx/app/views/pages/colors"
)

// === Models ==================================================================

var shadesMap = map[int]int{50: 0, 100: 1, 200: 2, 300: 3, 400: 4, 500: 5, 600: 6, 700: 7, 800: 8, 900: 9, 950: 10}
var shadeValues = []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950}

var orderedCodes = []Code{
	"bgc", "fgc",
	"rse", "bry", "chy", "rby", "red", "crl", "pmk", "orn", "sun", "gld", "hny", "yel",
	"lem", "acd", "lim", "spr", "grn", "emr", "jde", "frs", "lea", "tea", "cyn",
	"aqu", "rbn", "azr", "sky", "blu", "cbt", "sph", "ind", "lav", "prp", "vio", "pnk", "mag",
	"brk", "rst", "bej", "olv", "mss", "znc", "gry", "stn", "slt", "ash",
	"wht", "blk",
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
	"bgc": Color{"", "Base"},
	"fgc": Color{"", "Surface"},
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

	"brk": Color{"#847e80", "Brick"},
	"rst": Color{"#847f7e", "Rust"},
	"bej": Color{"#82807c", "Beige"},
	"olv": Color{"#7f817c", "Olive"},
	"mss": Color{"#7d817e", "Moss"},
	"znc": Color{"#7c8280", "Zinc"},
	"gry": Color{"#7c8182", "Gray"},
	"slt": Color{"#7d8184", "Slate"},
	"stn": Color{"#7e8084", "Stone"},
	"ash": Color{"#817f84", "Ash"},

	"wht": Color{"#cfcfcf", "White"},
	"blk": Color{"#0d0d0d", "Black"},
}

// === Handlers ================================================================

func Generate(seed string) Palette {
	p := make(Palette)
	// seed = "#ffffff"
	for _, code := range orderedCodes {
		color, ok := colors[code]
		if !ok {
			log.Println("[Warning]", code, "not found in palette for seed.")
			continue
		}
		if code == "bgc" || code == "fgc" {
			color.Seed = Hex(seed)
		}
		code.valid()
		oklch := color.Seed.toOklch()
		details := &ColorDetails{Color: color, Base: *oklch, Shades: make(map[int]Details, len(shadesMap))}
		p[code] = details

		switch code {
		case "bgc":
			details.generateBg()
		case "fgc":
			details.generateFg(p["bgc"].Shades[50].Oklch)
		case "brk", "rst", "olv", "mss", "znc", "slt", "gry", "stn", "ash", "bej":
			details.generateGrey()
		case "wht", "blk":
			details.generateBW()
		default:
			details.generateColor()
		}
	}
	return p
}

func (p *Palette) ToCSS(w io.Writer) {
	seed := (*p)["bgc"].Base
	fmt.Fprintln(w, ":root {")
	for _, code := range orderedCodes {
		colorDetails, ok := (*p)[code]
		if ok {
			colorDetails.ToCSS(w, code, seed)
		}
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "@theme {")
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
	for _, shadeKey := range shadeValues {
		color, ok := c.Shades[shadeKey]
		if !ok {
			continue
		}
		css := OklchToString(&color.Oklch)
		fmt.Fprintf(w, "  --%s-%d: %s;\n", string(code), shadeKey, css)
	}
	fmt.Fprintln(w, " ")
}

func (c ColorDetails) ToTheme(w io.Writer, code Code, seed oklab.Oklch) {
	for _, shadeKey := range shadeValues {
		color, ok := c.Shades[shadeKey]
		if !ok {
			continue
		}
		css := OklchToString(&color.Oklch)
		fmt.Fprintf(w, "  --color-%s-%d: %s;\n", string(code), shadeKey, css)
	}
	fmt.Fprintln(w, " ")
}

func (p Palette) ToView() []colorsPage.ColorScaleView {
	var views []colorsPage.ColorScaleView

	bgcDetails, ok := p["bgc"]
	if !ok {
		return views
	}
	seed := bgcDetails.Base

	convertScale := func(code Code, scale *ColorDetails) colorsPage.ColorScaleView {
		view := colorsPage.ColorScaleView{
			Name:   scale.Name,
			Code:   string(code),
			Value:  scale.Shades[600].Oklch,
			Shades: make([]colorsPage.Shade, len(shadeValues)),
		}

		for i, shadeVal := range shadeValues {
			colorDetails, ok := scale.Shades[shadeVal]
			if !ok {
				continue
			}

			color := colorDetails.Oklch
			rl, cr := OklchCompare(seed, color)
			hue := toDegree(color.H)

			radius := 37.0

			scaledRadiusC := radius * math.Tanh(6.0*color.C)
			totalDistanceC := scaledRadiusC

			scaledRadiusL := radius * math.Pow(color.L, 1.5)
			totalDistanceL := scaledRadiusL

			angle := -color.H
			view.Shades[i] = colorsPage.Shade{
				Code:  string(code),
				Value: shadeVal,
				RL:    rl,
				CR:    cr,
				L:     fmt.Sprintf("%0.1f%%", color.L*100),
				C:     fmt.Sprintf("%0.3f", color.C),
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
		details, ok := p[code]
		if ok {
			views = append(views, convertScale(code, details))
		}
	}

	return views
}
