package palette

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/alltom/oklab"
)

// HexToOklch converts a hex color string (e.g., "#RRGGBB") to its OKLCH representation.
func HexToOklch(hexColor string) (oklchColor oklab.Oklch, err error) {
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return oklchColor, fmt.Errorf("invalid hex color format: expected #RRGGBB, got %s", hexColor)
	}

	r, errR := strconv.ParseUint(hexColor[1:3], 16, 8)
	g, errG := strconv.ParseUint(hexColor[3:5], 16, 8)
	b, errB := strconv.ParseUint(hexColor[5:7], 16, 8)

	if errR != nil || errG != nil || errB != nil {
		return
	}

	rgbaColor := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	oklchColor = oklab.OklchModel.Convert(rgbaColor).(oklab.Oklch)

	return
}

// OklchToString converts an OKLCH color to a css string.
func OklchToString(oklchColor *oklab.Oklch) string {
	// Format hue with leading zeros (3 digits before, 2 after decimal)
	return fmt.Sprintf("oklch(%.2f %.3f %06.2f)", oklchColor.L, oklchColor.C, to360(oklchColor.H))
}

func to360(hue float64) float64 {
	hueDegrees := hue * 180 / math.Pi
	if hueDegrees < 0 {
		hueDegrees += 360
	}
	return hueDegrees
}

// sRGBToLinear converts an sRGB color component (0.0-1.0) to its linear representation.
func sRGBToLinear(c float64) float64 {
	if c <= 0.04045 { // WCAG 2.2 uses 0.04045, slightly different from 0.03928 in some older specs
		return c / 12.92
	} else {
		return math.Pow((c+0.055)/1.055, 2.4)
	}
}

// RelativeLuminance calculates the relative luminance of a color according to WCAG standards.
// Input color is OKLCH.
func RelativeLuminance(c oklab.Oklch) float64 {
	// Convert OKLCH -> OKLAB -> RGBA (returns values in 0-65535 range)
	rInt, gInt, bInt, _ := c.Oklab().RGBA()

	// Convert 0-65535 range to 0.0-1.0
	r := float64(rInt) / 65535.0
	g := float64(gInt) / 65535.0
	b := float64(bInt) / 65535.0

	// Convert sRGB to Linear RGB
	rLin := sRGBToLinear(r)
	gLin := sRGBToLinear(g)
	bLin := sRGBToLinear(b)

	// Calculate relative luminance
	// Coefficients from WCAG
	return 0.2126*rLin + 0.7152*gLin + 0.0722*bLin
}

// ContrastRatio calculates the contrast ratio between two OKLCH colors (c1, c2).
func ContrastRatio(c1, c2 oklab.Oklch) float64 {
	l1 := RelativeLuminance(c1)
	l2 := RelativeLuminance(c2)

	// Ensure l1 is the lighter color
	if l2 > l1 {
		l1, l2 = l2, l1
	}

	// Formula from WCAG
	return (l1 + 0.05) / (l2 + 0.05)
}

func OklchCompare(c1, c2 oklab.Oklch) (rl float64, cr float64) {
	rl = RelativeLuminance(c2)
	cr = ContrastRatio(c1, c2)
	return
}
