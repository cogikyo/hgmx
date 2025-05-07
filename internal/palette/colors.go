package palette

import (
	"github.com/alltom/oklab"
)

func hexMust(hex string) *oklab.Oklch {
	c, err := HexToOklch(hex)
	if err != nil {
		panic(err)
	}
	return &c
}

var shadesMap = map[int]int{
	50:  0,
	100: 1,
	200: 2,
	300: 3,
	400: 4,
	500: 5,
}

type shadeGenParams struct {
	l50 float64
	c50 float64
}

// generateShades dynamically creates shades 50-500 based on 500
// interpolating L and C linearly between the target 50 values and the 500 values.
func (c *ColorScale) generateShades(p *shadeGenParams) {
	c600 := c.shade[600]
	hue := c600.H

	numIntervals := 6.0

	if p == nil {
		p = &shadeGenParams{
			l50: 0.95,
			c50: 0.02,
		}
	}

	if p.l50 == 0 {
		p.l50 = 0.95
	}
	if p.c50 == 0 {
		p.c50 = 0.02
	}

	for shadeValue, index := range shadesMap {
		t := float64(index) / numIntervals
		l := p.l50 + (c600.L-p.l50)*t
		chroma := p.c50 + (c600.C-p.c50)*t
		c.shade[shadeValue] = &oklab.Oklch{L: l, C: chroma, H: hue}
	}
}

func Background(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("bgc")
	c.shade[50] = hexMust("#5a61aa")
	c.shade[100] = hexMust("#4e5492")
	c.shade[200] = hexMust("#3f4578")
	c.shade[300] = hexMust("#30345a")
	c.shade[400] = hexMust("#282b48")
	c.shade[500] = &base
	c.shade[600] = hexMust("#252841")
	c.shade[700] = hexMust("#1e2133")
	c.shade[800] = hexMust("#181a2c")
	c.shade[900] = hexMust("#131626")
	c.shade[950] = hexMust("#0d0f1b")
	return
}

func Foreground(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("fgc")
	c.shade[700] = hexMust("#aeb9f8")
	c.shade[600] = hexMust("#b6c0f7")
	c.shade[500] = hexMust("#bec6f8")
	c.shade[400] = hexMust("#cad1fb")
	c.shade[300] = hexMust("#d1d8ff")
	return
}

func Berry(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("bry")
	c.shade[600] = hexMust("#fd016f")
	c.generateShades(nil)
	return
}

func Cherry(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("chy")
	c.shade[600] = hexMust("#ff0457")
	c.generateShades(nil)
	return
}

func Ruby(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("rby")
	c.shade[600] = hexMust("#f9043a")
	c.generateShades(nil)
	return
}

func Red(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("red")
	c.shade[600] = hexMust("#fd181a")
	c.generateShades(nil)
	return
}

func Coral(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("crl")
	c.shade[600] = hexMust("#fd3a00")
	c.generateShades(nil)
	return
}

func Pumpkin(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("pmk")
	c.shade[600] = hexMust("#fd5802")
	c.generateShades(nil)
	return
}

func Orange(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("orn")
	c.shade[600] = hexMust("#ff7220")
	c.generateShades(nil)
	return
}

func Sun(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sun")
	c.shade[600] = hexMust("#ff9004")
	c.generateShades(nil)
	return
}

func Gold(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("gld")
	c.shade[600] = hexMust("#fead05")
	c.generateShades(nil)
	return
}

func Yellow(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("yel")
	c.shade[600] = hexMust("#ffcc00")
	c.generateShades(nil)
	return
}

func Lemon(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lem")
	c.shade[600] = hexMust("#f2e303")
	c.generateShades(nil)
	return
}

func Acid(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("acd")
	c.shade[600] = hexMust("#cdf118")
	c.generateShades(nil)
	return
}

func Lime(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lim")
	c.shade[600] = hexMust("#a7e902")
	c.generateShades(nil)
	return
}

func Kiwi(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("kwi")
	c.shade[600] = hexMust("#8ee300")
	c.generateShades(nil)
	return
}

func Green(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("grn")
	c.shade[600] = hexMust("#75dd03")
	c.generateShades(nil)
	return
}

func Spring(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("spr")
	c.shade[600] = hexMust("#53d102")
	c.generateShades(nil)
	return
}

func Emerald(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("emr")
	c.shade[600] = hexMust("#25c626")
	c.generateShades(nil)
	return
}

func Jade(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("jde")
	c.shade[600] = hexMust("#00bf4a")
	c.generateShades(nil)
	return
}

func Forest(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("frs")
	c.shade[600] = hexMust("#08bb65")
	c.generateShades(nil)
	return
}

func Leaf(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lea")
	c.shade[600] = hexMust("#01c37e")
	c.generateShades(nil)
	return
}

func Teal(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("tea")
	c.shade[600] = hexMust("#0ed39a")
	c.generateShades(nil)
	return
}

func Cyan(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("cyn")
	c.shade[600] = hexMust("#00e7cb")
	c.generateShades(nil)
	return
}

func Aqua(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("aqu")
	c.shade[600] = hexMust("#02eeef")
	c.generateShades(nil)
	return
}

func Robin(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("rbn")
	c.shade[600] = hexMust("#07e3fe")
	c.generateShades(nil)
	return
}

func Azure(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("azr")
	c.shade[600] = hexMust("#0acbff")
	c.generateShades(nil)
	return
}

func Sky(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sky")
	c.shade[600] = hexMust("#0aafff")
	c.generateShades(nil)
	return
}

func Blue(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("blu")
	c.shade[600] = hexMust("#0184fe")
	c.generateShades(nil)
	return
}

func Cobalt(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("cbt")
	c.shade[600] = hexMust("#256eff")
	c.generateShades(nil)
	return
}

func Sapphire(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sph")
	c.shade[600] = hexMust("#4158fa")
	c.generateShades(nil)
	return
}

func Indigo(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("ind")
	c.shade[600] = hexMust("#5a4aff")
	c.generateShades(nil)
	return
}

func Lavender(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lav")
	c.shade[600] = hexMust("#6e40ff")
	c.generateShades(nil)
	return
}

func Purple(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("prp")
	c.shade[600] = hexMust("#972eff")
	c.generateShades(nil)
	return
}

func Violet(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("vio")
	c.shade[600] = hexMust("#c602fe")
	c.generateShades(nil)
	return
}

func Pink(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("pnk")
	c.shade[600] = hexMust("#ed00e6")
	c.generateShades(nil)
	return
}

func Magenta(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("mag")
	c.shade[600] = hexMust("#fd01b9")
	c.generateShades(nil)
	return
}

func Rose(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("ros")
	c.shade[600] = hexMust("#fc0086")
	c.generateShades(nil)
	return
}

func Black(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("blk")
	c.shade[950] = hexMust("#0b0b0f") // blk-3
	c.shade[900] = hexMust("#101014") // blk-2
	c.shade[800] = hexMust("#16161a") // blk-1
	c.shade[700] = hexMust("#1d1d21") // blk-0
	return
}

func Gray(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("gry")
	c.shade[500] = hexMust("#4f5163") // gry-2 (Approx midpoint)
	c.shade[400] = hexMust("#3f414f") // gry-1
	c.shade[300] = hexMust("#373945") // gry-0
	c.shade[600] = hexMust("#5f6278") // gry-3
	c.shade[700] = hexMust("#6d7089") // gry-4
	c.shade[800] = hexMust("#7f8199") // gry-5
	return
}

func White(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("wht")
	c.shade[500] = hexMust("#ddddf6") // wht-2 (Approx midpoint)
	c.shade[400] = hexMust("#d3d3ed") // wht-1
	c.shade[300] = hexMust("#c9c9e2") // wht-0
	c.shade[600] = hexMust("#e9e9fb") // wht-3
	return
}

func Adenine(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("ade")
	c.shade[400] = hexMust("#824141") // ade-0
	c.shade[300] = hexMust("#b15e5b") // ade-1
	c.shade[200] = hexMust("#c67a79") //[] ade-2
	c.shade[100] = hexMust("#d09490") // ade-3
	return
}

func Rust(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("rst")
	c.shade[600] = hexMust("#493531") // rst-0
	c.shade[500] = hexMust("#563e39") // rst-1
	c.shade[400] = hexMust("#694b44") // rst-2 (Approx midpoint)
	c.shade[300] = hexMust("#805a52") // rst-3
	c.shade[200] = hexMust("#92675d") // rst-4
	c.shade[100] = hexMust("#a3786d") // rst-5
	return
}

func Cytosine(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("cyt")
	c.shade[400] = hexMust("#505831") // cyt-0
	c.shade[300] = hexMust("#717b45") // cyt-1
	c.shade[200] = hexMust("#8a945b") // cyt-2
	c.shade[100] = hexMust("#9ea876") // cyt-3
	return
}

func Olive(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("olv")
	c.shade[100] = hexMust("#818382") // olv-5
	c.shade[200] = hexMust("#6e7270") // olv-4
	c.shade[300] = hexMust("#5f6361") // olv-3
	c.shade[400] = hexMust("#505251") // olv-2 (Approx midpoint)
	c.shade[500] = hexMust("#414342") // olv-1
	c.shade[600] = hexMust("#383a39") // olv-0
	return
}

func Slate(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("slt")
	c.shade[100] = hexMust("#7b7fb0") // slt-5
	c.shade[200] = hexMust("#686ea1") // slt-4
	c.shade[300] = hexMust("#585f8d") // slt-3
	c.shade[400] = hexMust("#484e75") // slt-2 (Approx midpoint)
	c.shade[500] = hexMust("#3c4162") // slt-1
	c.shade[600] = hexMust("#343852") // slt-0
	return
}

func Thymine(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("thy")
	c.shade[400] = hexMust("#3b557c") // thy-0
	c.shade[300] = hexMust("#5b77a4") // thy-1
	c.shade[200] = hexMust("#7690b9") // thy-2
	c.shade[100] = hexMust("#90a4c7") // thy-3
	return
}

func Guanine(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("gau")
	c.shade[400] = hexMust("#6f447a") // gau-0
	c.shade[300] = hexMust("#9961a7") // gau-1
	c.shade[200] = hexMust("#af7dba") // gau-2
	c.shade[100] = hexMust("#c193cd") // gau-3
	return
}

func Plum(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("plm")
	c.shade[100] = hexMust("#977998") // plm-5
	c.shade[200] = hexMust("#876888") // plm-4
	c.shade[300] = hexMust("#765a77") // plm-3
	c.shade[400] = hexMust("#634a64") // plm-2 (Approx midpoint)
	c.shade[500] = hexMust("#523c52") // plm-1
	c.shade[600] = hexMust("#453445") // plm-0
	return
}
