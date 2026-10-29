package palette

func (c Color) Luminance() float64 {
	return float64(0.2126*linear(c.R)) + float64(0.7152*linear(c.G)) + float64(0.0722*linear(c.B))
}

func Ratio(a, b Color) float64 {
	x, y := a.Luminance(), b.Luminance()
	if x < y {
		x, y = y, x
	}
	return (x + 0.05) / (y + 0.05)
}
