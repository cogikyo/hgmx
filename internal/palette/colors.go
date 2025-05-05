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
	600: 6,
	700: 7,
	800: 8,
	900: 9,
	950: 10,
}

var numIntervals = 10.0

// generateShades dynamically creates shades 50-900 based on 950
// interpolating L and C linearly between the target 50 values and the 950 values.
func generateShades(c ColorScale, l50, c50 float64) {
	c950 := c.shade[950]
	hue := c950.H // Keep hue constant for now

	for shadeValue, index := range shadesMap {
		if shadeValue == 950 {
			continue
		}
		t := float64(index) / numIntervals

		l := l50 + (c950.L-l50)*t
		chroma := c50 + (c950.C-c50)*t

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
	c.shade[700] = hexMust("#aeb9f8") // brt-0
	c.shade[600] = hexMust("#b6c0f7") // brt-1
	c.shade[500] = hexMust("#bec6f8") // brt-2
	c.shade[400] = hexMust("#cad1fb") // brt-3
	c.shade[300] = hexMust("#d1d8ff") // brt-4
	return
}

func Ruby(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("rby")
	c.shade[950] = hexMust("#ff0457")
	generateShades(c, 0.95, 0.02)
	return
}

func Red(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("red")
	c.shade[950] = hexMust("#ff1701")

	// c.shade[950] = hexMust("#ff1701") // red-0
	// c.shade[500] = hexMust("#f34658") // red-0
	// c.shade[400] = hexMust("#f36978") // red-1
	// c.shade[300] = hexMust("#f07a88") // red-2
	// c.shade[200] = hexMust("#f08898") // red-3
	// c.shade[100] = hexMust("#f29ca9") // red-4

	generateShades(c, 0.95, 0.02)

	return
}

func Orange(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("orn")
	c.shade[950] = hexMust("#fd6a00") // orn-0

	// c.shade[500] = hexMust("#e56b2c") // orn-0
	// c.shade[400] = hexMust("#ea834b") // orn-1
	// c.shade[300] = hexMust("#eb905d") // orn-2
	// c.shade[200] = hexMust("#f2a170") // orn-3
	// c.shade[100] = hexMust("#f8b486") // orn-4

	generateShades(c, 0.95, 0.02)

	return
}

func Sun(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sun")
	c.shade[950] = hexMust("#ff9004") // sun-0

	// c.shade[500] = hexMust("#f3a338") // sun-0
	// c.shade[400] = hexMust("#f5b855") // sun-1
	// c.shade[300] = hexMust("#f5c069") // sun-2
	// c.shade[200] = hexMust("#f4ce88") // sun-3
	// c.shade[100] = hexMust("#f5d599") // sun-4

	generateShades(c, 0.95, 0.02)

	return
}

func Gold(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("gld")
	c.shade[950] = hexMust("#fead05") // gld-0

	generateShades(c, 0.95, 0.02)
	return
}

func Yellow(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("yel")
	c.shade[950] = hexMust("#ffcc00") // yel-0

	generateShades(c, 0.95, 0.02)
	return
}

func Lemon(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lem")
	c.shade[950] = hexMust("#fcee29")

	generateShades(c, 0.95, 0.02)
	return
}

func Lime(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lim")
	c.shade[950] = hexMust("#d6fb12")

	generateShades(c, 0.95, 0.02)
	return
}

func Green(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("grn")
	c.shade[950] = hexMust("#74dc01")

	generateShades(c, 0.95, 0.02)
	return
}

func Emerald(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("emr")
	c.shade[950] = hexMust("#02c059")

	// c.shade[500] = hexMust("#2d9a43") // emr-0
	// c.shade[400] = hexMust("#48a95b") // emr-1
	// c.shade[300] = hexMust("#5aba6d") // emr-2
	// c.shade[200] = hexMust("#5fc976") // emr-3
	// c.shade[100] = hexMust("#76d78b") // emr-4

	generateShades(c, 0.95, 0.02)
	return
}

func Teal(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("tea")
	c.shade[950] = hexMust("#02fdb9")

	generateShades(c, 0.95, 0.02)
	return
}

func Cyan(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("cyn")
	c.shade[950] = hexMust("#36ffe2") // cyn-0

	// c.shade[500] = hexMust("#2bb198") // cyn-0
	// c.shade[400] = hexMust("#30c9b0") // cyn-1
	// c.shade[300] = hexMust("#38d2ba") // cyn-2
	// c.shade[200] = hexMust("#50dec8") // cyn-3
	// c.shade[100] = hexMust("#75e6d5") // cyn-4

	generateShades(c, 0.95, 0.02)
	return
}

func Aqua(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("aqu")
	c.shade[950] = hexMust("#50ffff")

	generateShades(c, 0.95, 0.02)
	return
}

func Cerulean(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("cer")
	c.shade[950] = hexMust("#07e3fe")

	generateShades(c, 0.95, 0.02)
	return
}

func Azure(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("azr")
	c.shade[950] = hexMust("#0acbff")

	generateShades(c, 0.95, 0.02)
	return
}

func Sky(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sky")
	c.shade[950] = hexMust("#0aafff")

	// 	c.shade[500] = hexMust("#369fd7") // sky-0
	// 	c.shade[400] = hexMust("#54b0e2") // sky-1
	// 	c.shade[300] = hexMust("#6bbdec") // sky-2
	// 	c.shade[200] = hexMust("#7cc5ef") // sky-3
	// 	c.shade[100] = hexMust("#90d1f5") // sky-4

	generateShades(c, 0.95, 0.02)
	return
}

func Blue(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("blu")
	c.shade[950] = hexMust("#0184fe")

	generateShades(c, 0.95, 0.02)
	return
}

func Sapphire(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("sap")
	c.shade[950] = hexMust("#4158fa") // sap-0

	// c.shade[500] = hexMust("#4a6be3") // blu-0
	// c.shade[400] = hexMust("#6380ec") // blu-1
	// c.shade[300] = hexMust("#7492ef") // blu-2
	// c.shade[200] = hexMust("#8aa4f3") // blu-3
	// c.shade[100] = hexMust("#9db2f4") // blu-4

	generateShades(c, 0.95, 0.02)
	return
}

func Lavender(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("lav")
	c.shade[950] = hexMust("#6e40ff")

	// c.shade[500] = hexMust("#7f61cd")
	// c.shade[400] = hexMust("#9376d8") // lav-1
	// c.shade[300] = hexMust("#a188df") // lav-2
	// c.shade[200] = hexMust("#b29ae8") // lav-3
	// c.shade[100] = hexMust("#bdaaeb") // lav-4

	generateShades(c, 0.95, 0.02)
	return
}

func Purple(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("prp")
	c.shade[950] = hexMust("#9201fd")

	generateShades(c, 0.95, 0.02)
	return
}

func Violet(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("vio")
	c.shade[950] = hexMust("#c602fe")

	generateShades(c, 0.95, 0.02)
	return
}

func Pink(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("pnk")
	c.shade[950] = hexMust("#ff13f7")

	generateShades(c, 0.95, 0.02)
	return
}

func Magenta(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("mag")
	c.shade[950] = hexMust("#fd01b9")

	generateShades(c, 0.95, 0.02)
	return
}

func Rose(seed oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("ros")
	c.shade[950] = hexMust("#f00a80") // pnk-0

	// c.shade[500] = hexMust("#d15da6") // pnk-0
	// c.shade[400] = hexMust("#e36cb8") // pnk-1
	// c.shade[300] = hexMust("#ea76c0") // pnk-2
	// c.shade[200] = hexMust("#e887c3") // pnk-3
	// c.shade[100] = hexMust("#ed9acd") // pnk-4

	generateShades(c, 0.95, 0.02)
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

func Forest(base oklab.Oklch) (c ColorScale) {
	c = ColorScale{}.New("frt")
	c.shade[400] = hexMust("#375c47") // frt-0
	c.shade[300] = hexMust("#4b8163") // frt-1
	c.shade[200] = hexMust("#5a9c78") // frt-2
	c.shade[100] = hexMust("#72b08e") // frt-3
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
