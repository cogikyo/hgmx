package palette

import (
	"fmt"
	"math"
	"strconv"
)

type Color struct{ R, G, B uint8 }

type Lab struct{ L, A, B float64 }

type LCh struct{ L, C, H float64 }

func ParseHex(s string) (Color, error) {
	if len(s) != 7 || s[0] != '#' {
		return Color{}, fmt.Errorf("palette: %q is not a #RRGGBB colour", s)
	}
	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return Color{}, fmt.Errorf("palette: %q is not a #RRGGBB colour", s)
	}
	return Color{uint8(v >> 16), uint8(v >> 8), uint8(v)}, nil
}

func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func (c Color) String() string { return c.Hex() }

func (c Color) MarshalText() ([]byte, error) { return []byte(c.Hex()), nil }

func (c *Color) UnmarshalText(b []byte) error {
	v, err := ParseHex(string(b))
	if err != nil {
		return err
	}
	*c = v
	return nil
}

func (c Color) Lab() Lab {
	r, g, b := linear(c.R), linear(c.G), linear(c.B)
	l := math.Cbrt(dot(0.4122214708, 0.5363325363, 0.0514459929, r, g, b))
	m := math.Cbrt(dot(0.2119034982, 0.6806995451, 0.1073969566, r, g, b))
	s := math.Cbrt(dot(0.0883024619, 0.2817188376, 0.6299787005, r, g, b))
	return Lab{
		L: dot(0.2104542553, 0.7936177850, -0.0040720468, l, m, s),
		A: dot(1.9779984951, -2.4285922050, 0.4505937099, l, m, s),
		B: dot(0.0259040371, 0.7827717662, -0.8086757660, l, m, s),
	}
}

func (c Color) LCh() LCh { return c.Lab().LCh() }

func (p Lab) LCh() LCh {
	h := math.Atan2(p.B, p.A) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return LCh{L: p.L, C: math.Hypot(p.A, p.B), H: h}
}

func (p LCh) Lab() Lab {
	rad := p.H * math.Pi / 180
	return Lab{L: p.L, A: p.C * math.Cos(rad), B: p.C * math.Sin(rad)}
}

func (p Lab) rgb() (r, g, b float64) {
	l := dot(1, 0.3963377774, 0.2158037573, p.L, p.A, p.B)
	m := dot(1, -0.1055613458, -0.0638541728, p.L, p.A, p.B)
	s := dot(1, -0.0894841775, -1.2914855480, p.L, p.A, p.B)
	l, m, s = l*l*l, m*m*m, s*s*s
	return dot(4.0767416621, -3.3077115913, 0.2309699292, l, m, s),
		dot(-1.2684380046, 2.6097574011, -0.3413193965, l, m, s),
		dot(-0.0041960863, -0.7034186147, 1.7076147010, l, m, s)
}

func dot(a, b, c, x, y, z float64) float64 {
	return float64(a*x) + float64(b*y) + float64(c*z)
}

const (
	gamutEps   = 1e-9
	achromatic = 1e-3
)

func (p Lab) inGamut() bool {
	r, g, b := p.rgb()
	return r >= -gamutEps && r <= 1+gamutEps &&
		g >= -gamutEps && g <= 1+gamutEps &&
		b >= -gamutEps && b <= 1+gamutEps
}

func (p Lab) Color() Color {
	if math.Hypot(p.A, p.B) < achromatic {
		v := encode(p.L * p.L * p.L)
		return Color{v, v, v}
	}
	r, g, b := p.rgb()
	return Color{encode(r), encode(g), encode(b)}
}

func (p LCh) Color() Color { return p.Lab().Color() }

func maxChroma(l, h float64) float64 {
	if l <= 0 || l >= 1 {
		return 0
	}
	rad := h * math.Pi / 180
	cos, sin := math.Cos(rad), math.Sin(rad)
	lo, hi := 0.0, 0.4
	for range 20 {
		mid := (lo + hi) / 2
		if (Lab{l, mid * cos, mid * sin}).inGamut() {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

func linear(v uint8) float64 {
	c := float64(v) / 255
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func encode(c float64) uint8 {
	c = min(max(c, 0), 1)
	if c <= 0.0031308 {
		c *= 12.92
	} else {
		c = float64(1.055*math.Pow(c, 1/2.4)) - 0.055
	}
	return uint8(math.Round(c * 255))
}
