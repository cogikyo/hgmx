package palette

import (
	"bytes"
	"flag"
	"math"
	"os"
	"testing"
)

var update = flag.Bool("update", false, "rewrite testdata golden files")

func TestGuarantees(t *testing.T) {
	levels := []uint8{0x00, 0x22, 0x55, 0x80, 0x99, 0xcc, 0xff}
	primaries := []Color{{0x01, 0x84, 0xfe}, {0x71, 0x32, 0x4b}, {0xe0, 0xa0, 0x30}, {0x80, 0x80, 0x80}, {0x00, 0xff, 0x00}}
	bases := []Color{{0x00, 0x11, 0x88}}
	for _, r := range levels {
		for _, g := range levels {
			for _, b := range levels {
				bases = append(bases, Color{r, g, b})
			}
		}
	}
	for _, base := range bases {
		for _, pri := range primaries {
			check(t, Generate(Seed{Base: base, Primary: pri}))
		}
	}
}

func check(t *testing.T, p Palette) {
	t.Helper()
	base := p.Seed.Base
	grounds := []Color{base, role(p, "inset"), role(p, "surface")}
	dark := Ratio(white, base) >= Ratio(black, base)
	worst := func(c Color) float64 {
		w := Ratio(c, base)
		for _, g := range grounds {
			w = min(w, Ratio(c, g))
		}
		return w
	}
	dl := func(c Color) float64 { return math.Abs(c.Lab().L - base.Lab().L) }
	meets := func(c Color, want target) bool {
		return Ratio(c, base) >= want.ratio && worst(c) >= want.floor && (!dark || want.dl == 0 || dl(c) >= want.dl)
	}
	pole := black
	if dark {
		pole = white
	}
	unreachable := func(want target) bool { return !meets(pole, want) }
	for _, s := range p.Scales {
		ramps := [][]Color{s.Steps}
		if s.Beta != nil {
			ramps = append(ramps, s.Beta)
		}
		for _, ramp := range ramps {
			if s.Name != "neutral" && len(ramp) != len(Steps)-1 {
				t.Fatalf("%s %s: %d steps, coloured scales end at 900", base, s.Name, len(ramp))
			}
			for i, c := range ramp {
				want := targets[i]
				if math.IsInf(want.ratio, 1) {
					if c == ramp[i-1] && i < s.Reach {
						t.Fatalf("%s %s-%d %s repeats the step before it, not reported", base, s.Name, Steps[i], c)
					}
				} else if !meets(c, want) && (!unreachable(want) || i < s.Reach) {
					t.Fatalf("%s %s-%d %s: %.3f:1 (tier %.2f:1), worst ground %.3f:1 (floor %.1f:1), ΔL %.3f (tier %.2f); reachable or not reported", base, s.Name, Steps[i], c, Ratio(c, base), want.ratio, worst(c), want.floor, dl(c), want.dl)
				}
				if i > 0 && Ratio(c, base) < Ratio(ramp[i-1], base) {
					t.Fatalf("%s %s-%d %s: %.4f:1 after %.4f:1, steps must not lose contrast", base, s.Name, Steps[i], c, Ratio(c, base), Ratio(ramp[i-1], base))
				}
			}
		}
		if w := worst(s.Fill); w < targets[FillStep].floor && !unreachable(targets[FillStep]) {
			t.Fatalf("%s %s fill %s: %.3f:1 on its worst ground, below %.1f:1", base, s.Name, s.Fill, w, targets[FillStep].floor)
		}
		if r := Ratio(s.On, s.Fill); r < onRatio {
			t.Fatalf("%s on-%s %s over %s: %.3f:1 below %.1f:1", base, s.Name, s.On, s.Fill, r, onRatio)
		}
		for _, toned := range []Color{base, role(p, "strong")} {
			if Ratio(toned, s.Fill) >= onRatio && s.On != base && s.On != role(p, "strong") {
				t.Fatalf("%s on-%s is %s although theme-toned %s reaches %.2f:1", base, s.Name, s.On, toned, Ratio(toned, s.Fill))
			}
		}
	}
}

func role(p Palette, name string) Color {
	for _, r := range p.Roles {
		if r.Name == name {
			return r.Color
		}
	}
	panic("palette: no role " + name)
}

func TestGolden(t *testing.T) {
	var buf bytes.Buffer
	for _, th := range []struct{ name, base, primary, accent string }{
		{"dusk", "#222536", "#0184fe", ""},
		{"paper", "#f8ecc9", "#71324b", "#2f6f8f"},
	} {
		s, err := ParseSeed(th.base, th.primary, th.accent)
		if err != nil {
			t.Fatal(err)
		}
		sel, err := Selector(th.name)
		if err != nil {
			t.Fatal(err)
		}
		if err := Generate(s).CSS(&buf, sel); err != nil {
			t.Fatal(err)
		}
	}
	const golden = "testdata/themes.css"
	if *update {
		if err := os.WriteFile(golden, buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("CSS output differs from %s; rerun with -update if the change is intended", golden)
	}
}
