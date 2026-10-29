package theme

import (
	"testing"

	"github.com/cogikyo/hgmx/palette"
)

func TestMockPairsPass(t *testing.T) {
	levels := []uint8{0x00, 0x22, 0x55, 0xcc, 0xf8, 0xff}
	checked, skipped := 0, 0
	primaries := []string{"#0184fe", "#71324b", "#e0a030", "#3fb950"}
	for _, r := range levels {
		for _, g := range levels {
			for _, b := range levels {
				base := palette.Color{R: r, G: g, B: b}
				for _, hex := range primaries {
					seed, err := palette.ParseSeed(base.Hex(), hex, "")
					if err != nil {
						t.Fatal(err)
					}
					v := View{Seed: seed, Palette: palette.Generate(seed)}
					for _, p := range v.mock().pairs() {
						if p.Scale != "" && p.Step >= v.scale(p.Scale).Reach {
							skipped++
							continue
						}
						checked++
						if r := palette.Ratio(p.Fg, p.Bg); r < p.Min {
							t.Fatalf("%s %s: %s on %s is %.2f:1, below %.1f:1", base, p.Label, p.Fg, p.Bg, r, p.Min)
						}
					}
				}
			}
		}
	}
	t.Logf("%d pairs checked, %d skipped because the base leaves their step out of reach", checked, skipped)
	if checked < 12000 {
		t.Fatalf("only %d pairs had room on their base", checked)
	}
}
