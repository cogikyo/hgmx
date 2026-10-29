package palette

import (
	"fmt"
	"math"
	"strings"
)

var Steps = [...]int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950}

type target struct {
	dl    float64
	ratio float64
	floor float64
	tint  float64
}

const (
	FillStep = 5
	TextStep = 6
	hueSteps = len(Steps) - 1
)

var targets = [len(Steps)]target{
	{ratio: 1.08, tint: 0.35},
	{ratio: 1.18, tint: 0.45},
	{ratio: 1.35, tint: 0.55},
	{dl: 0.20, ratio: 1.6, tint: 0.7},
	{dl: 0.32, ratio: 2.1, tint: 0.85},
	{dl: 0.42, ratio: 3, floor: 3, tint: 1},
	{dl: 0.51, ratio: 4.5, floor: 4.5, tint: 1},
	{dl: 0.59, ratio: 7, floor: 4.5, tint: 0.9},
	{dl: 0.66, ratio: 10, floor: 4.5, tint: 0.8},
	{dl: 0.70, ratio: 14, floor: 4.5, tint: 1},
	{dl: math.Inf(1), ratio: math.Inf(1), tint: 0.3},
}

const (
	insetRatio = 1.1
	onRatio    = 4.5
	accentTurn = 150
	betaTurn   = 15
	secondary  = 0.4
	toneGain   = 2
	clashTurn  = 25
	clashRel   = 0.2
)

var spectrum = []struct {
	name  string
	angle float64
	rel   float64
}{
	{"red", 20, 0.85},
	{"orange", 50, 0.85},
	{"yellow", 85, 0.85},
	{"lime", 125, 0.7},
	{"green", 150, 0.85},
	{"cyan", 190, 0.9},
	{"sky", 230, 0.85},
	{"blue", 265, 0.8},
	{"purple", 300, 0.65},
	{"pink", 345, 0.65},
}

var status = []struct{ name, hue string }{
	{"info", "sky"},
	{"success", "green"},
	{"warning", "yellow"},
	{"error", "red"},
}

type Seed struct {
	Base    Color
	Primary Color
	Accent  *Color
}

type Palette struct {
	Seed   Seed
	Dark   bool
	Scales []Scale
	Roles  []Token
}

type Scale struct {
	Name   string
	Alias  string
	Hue    float64
	Steps  []Color
	Beta   []Color
	Fill   Color
	Seeded bool
	On     Color
	Reach  int
}

func ParseSeed(base, primary, accent string) (Seed, error) {
	var s Seed
	var err error
	if s.Base, err = ParseHex(base); err != nil {
		return s, err
	}
	if s.Primary, err = ParseHex(primary); err != nil {
		return s, err
	}
	if accent != "" {
		a, err := ParseHex(accent)
		if err != nil {
			return s, err
		}
		s.Accent = &a
	}
	return s, nil
}

func Generate(s Seed) Palette {
	g := newGen(s.Base)
	draft := g.ramp(g.lch.H, len(Steps), g.tone)
	g.grounds = []Color{s.Base, g.inset(draft[0]), draft[1]}
	neutral := Scale{Name: "neutral", Hue: g.lch.H, Steps: g.ramp(g.lch.H, len(Steps), g.tone)}
	g.toned = []Color{s.Base, neutral.Steps[hueSteps-1]}
	neutral = g.finish(neutral, nil)
	p := Palette{Seed: s, Dark: g.dark, Scales: []Scale{neutral}}

	pri := s.Primary.LCh()
	rel := relChroma(pri)
	p.Scales = append(p.Scales, g.brand("primary", pri.H, rel, &s.Primary), g.brand("secondary", pri.H, rel*secondary, nil))
	if s.Accent != nil {
		acc := s.Accent.LCh()
		p.Scales = append(p.Scales, g.brand("accent", acc.H, relChroma(acc), s.Accent))
	} else {
		p.Scales = append(p.Scales, g.brand("accent", math.Mod(pri.H+accentTurn, 360), rel, nil))
	}

	hues := map[string]Scale{}
	for _, h := range spectrum {
		hues[h.name] = g.finish(Scale{Name: h.name, Hue: h.angle, Steps: g.ramp(h.angle, hueSteps, vivid(h.angle, h.rel))}, nil)
		p.Scales = append(p.Scales, hues[h.name])
	}
	for _, st := range status {
		sc := hues[st.hue]
		sc.Name, sc.Alias = st.name, st.hue
		p.Scales = append(p.Scales, sc)
	}

	role := func(name string, i int) Token { return Token{name, neutral.Steps[i], neutral.Token(i)} }
	p.Roles = []Token{
		{Name: "base", Color: s.Base},
		{Name: "inset", Color: g.grounds[1]},
		role("surface", 1),
		role("raised", 2),
		role("line", 2),
		role("muted", TextStep),
		role("fg", 8),
		role("strong", 9),
		role("ink", 10),
	}
	return p
}

func (s Scale) Token(i int) string { return fmt.Sprintf("%s-%d", s.Name, Steps[i]) }

func (s Scale) BetaToken(i int) string { return fmt.Sprintf("%s-%d", s.Name, Steps[i]+1) }

func (p Palette) Scale(name string) (Scale, bool) {
	for _, s := range p.Scales {
		if s.Name == name {
			return s, true
		}
	}
	return Scale{}, false
}

type gen struct {
	base    Color
	lch     LCh
	dark    bool
	grounds []Color
	toned   []Color
}

func newGen(base Color) gen {
	return gen{base: base, lch: base.LCh(), dark: Ratio(white, base) >= Ratio(black, base), grounds: []Color{base}}
}

var white, black = Color{255, 255, 255}, Color{}

func (g gen) dl(c Color) float64 { return math.Abs(c.Lab().L - g.lch.L) }

func (g gen) grounded(c Color, floor float64) bool {
	for _, gr := range g.grounds {
		if Ratio(c, gr) < floor {
			return false
		}
	}
	return true
}

func (g gen) meets(c Color, t target) bool {
	return (!g.dark || t.dl == 0 || g.dl(c) >= t.dl) && Ratio(c, g.base) >= t.ratio && g.grounded(c, t.floor)
}

func (g gen) tone(l float64, _ target) float64 {
	if g.lch.L == 0 {
		return 0
	}
	return min(g.lch.C*min(l/g.lch.L, toneGain), maxChroma(l, g.lch.H))
}

func (g gen) solve(h float64, chroma func(l float64) float64, forward bool, ok func(Color) bool) Color {
	end := 0.0
	if g.dark == forward {
		end = 1
	}
	at := func(t float64) Color {
		l := g.lch.L + float64(t*(end-g.lch.L))
		return LCh{l, chroma(l), h}.Color()
	}
	if !ok(at(1)) {
		return at(1)
	}
	lo, hi := 0.0, 1.0
	for range 24 {
		mid := (lo + hi) / 2
		if ok(at(mid)) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return at(hi)
}

func (g gen) ramp(h float64, n int, chroma func(l float64, t target) float64) []Color {
	out := make([]Color, n)
	var floor target
	for i, t := range targets[:n] {
		t.dl, t.ratio = max(t.dl, floor.dl), max(t.ratio, floor.ratio)
		ok := func(c Color) bool { return g.meets(c, t) }
		out[i] = g.solve(h, func(l float64) float64 { return chroma(l, t) }, true, ok)
		floor = target{dl: g.dl(out[i]), ratio: Ratio(out[i], g.base)}
	}
	return out
}

func (g gen) reach(steps []Color) int {
	for i, c := range steps {
		pole := math.IsInf(targets[i].ratio, 1)
		if pole && c == steps[i-1] || !pole && !g.meets(c, targets[i]) {
			return i
		}
	}
	return len(steps)
}

func (g gen) inset(fallback Color) Color {
	t := target{ratio: insetRatio}
	c := g.solve(g.lch.H, func(l float64) float64 { return g.tone(l, t) }, false, func(c Color) bool { return g.meets(c, t) })
	if g.meets(c, t) {
		return c
	}
	return fallback
}

func (g gen) brand(name string, h, rel float64, fill *Color) Scale {
	sc := Scale{Name: name, Hue: h, Steps: g.ramp(h, hueSteps, vivid(h, rel))}
	bh := math.Mod(h+betaTurn, 360)
	sc.Beta = g.ramp(bh, hueSteps, vivid(bh, rel))
	return g.finish(sc, fill)
}

func vivid(h, rel float64) func(l float64, t target) float64 {
	return func(l float64, t target) float64 { return rel * t.tint * maxChroma(l, h) }
}

func (g gen) finish(sc Scale, seed *Color) Scale {
	sc.Reach = g.reach(sc.Steps)
	sc.Fill = sc.Steps[FillStep]
	if seed != nil && Ratio(*seed, g.base) >= targets[FillStep].ratio && g.grounded(*seed, targets[FillStep].floor) {
		sc.Fill, sc.Seeded = *seed, true
	}
	sc.On = g.on(sc.Fill)
	return sc
}

func (g gen) on(fill Color) Color {
	best, score := Color{}, 0.0
	for _, c := range g.toned {
		if r := Ratio(c, fill); r >= onRatio && r > score {
			best, score = c, r
		}
	}
	if score > 0 {
		return best
	}
	if Ratio(white, fill) > Ratio(black, fill) {
		return white
	}
	return black
}

type Clash struct {
	Role    string
	Status  string
	Degrees float64
	Step    int
	Best    float64
}

func (c Clash) String() string {
	if c.Status != "" {
		return fmt.Sprintf("%s is %.0f° from %s's hue, so pair %s with an icon or label", c.Role, c.Degrees, c.Status, c.Status)
	}
	return fmt.Sprintf("the base leaves too little room: %s from %d up cannot reach their targets (best %.2f:1)", c.Role, c.Step, c.Best)
}

func (p Palette) Clashes() []Clash {
	var out []Clash
	for _, role := range []string{"primary", "accent"} {
		sc, _ := p.Scale(role)
		if relChroma(sc.Fill.LCh()) < clashRel {
			continue
		}
		for _, st := range status {
			hue, _ := p.Scale(st.hue)
			d := math.Abs(math.Mod(sc.Hue-hue.Hue+540, 360) - 180)
			if d < clashTurn {
				out = append(out, Clash{Role: role, Status: st.name, Degrees: d})
			}
		}
	}
	var short []string
	step, best := len(Steps), max(Ratio(white, p.Seed.Base), Ratio(black, p.Seed.Base))
	for _, s := range p.Scales {
		if s.Alias == "" && s.Reach < len(s.Steps) {
			short = append(short, s.Name)
			step = min(step, s.Reach)
		}
	}
	if len(short) > 0 {
		role := strings.Join(short, ", ")
		if len(short) > 3 {
			role = fmt.Sprintf("%d scales", len(short))
		}
		out = append(out, Clash{Role: role, Step: Steps[step], Best: best})
	}
	return out
}

func relChroma(c LCh) float64 {
	m := maxChroma(c.L, c.H)
	if m == 0 {
		return 0
	}
	return min(c.C/m, 1)
}
