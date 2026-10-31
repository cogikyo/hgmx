package theme

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cogikyo/hgmx/palette"
)

type Preset struct {
	Name string
	Seed palette.Seed
}

var Default = palette.Seed{Base: palette.Color{R: 0x22, G: 0x25, B: 0x36}, Primary: palette.Color{R: 0x01, G: 0x84, B: 0xfe}}

const exportName = "custom"

// View holds an explorer's seed, solved palette, and input error.
type View struct {
	Seed    palette.Seed
	Palette palette.Palette
	Action  string
	Presets []Preset
	Err     error
	assets  assetURLs
}

// View prepares [Explorer] using Default unless base or primary is present.
// Invalid seed input sets Err; use v.Status() for the HTTP response status.
// Action is the explorer page's URL, used for forms and preset links.
func (h *Handler) View(action string, q url.Values, presets ...Preset) View {
	v := View{Seed: Default, Action: action, Presets: presets, assets: h.explorer}
	if q.Has("base") || q.Has("primary") {
		seed, err := palette.ParseSeed(q.Get("base"), q.Get("primary"), q.Get("accent"))
		if err == nil {
			v.Seed = seed
		}
		v.Err = err
	}
	v.Palette = palette.Generate(v.Seed)
	return v
}

// Status returns 422 for invalid seed input, otherwise 200.
func (v View) Status() int {
	if v.Err != nil {
		return http.StatusUnprocessableEntity
	}
	return http.StatusOK
}

func (v View) Link(s palette.Seed) string {
	q := url.Values{"base": {s.Base.Hex()}, "primary": {s.Primary.Hex()}}
	if s.Accent != nil {
		q.Set("accent", s.Accent.Hex())
	}
	return v.Action + "?" + q.Encode()
}

func (v View) Accent() string {
	if v.Seed.Accent == nil {
		return ""
	}
	return v.Seed.Accent.Hex()
}

// Live returns the seed query for a data-theme-live attribute.
func (v View) Live() string {
	return strings.TrimPrefix(Sheet("live", v.Seed), sheet+"?")
}

func (v View) download(format string) string {
	return Sheet(exportName, v.Seed) + "&download=" + format
}

func (v View) style() string {
	primary := v.scale("primary")
	return fmt.Sprintf(
		"--tx-base:%s;--tx-fg:%s;--tx-muted:%s;--tx-line:%s;--tx-surface:%s;--tx-accent:%s;--tx-on-accent:%s;--tx-focus:%s;--tx-error:%s",
		v.Seed.Base.Hex(), v.role("fg").Hex(), v.role("muted").Hex(),
		v.role("line").Hex(), v.role("surface").Hex(), primary.Fill.Hex(),
		primary.On.Hex(), v.scale("accent").Steps[palette.FillStep].Hex(), v.scale("error").Steps[7].Hex(),
	)
}

func (v View) scale(name string) palette.Scale {
	s, _ := v.Palette.Scale(name)
	return s
}

func (v View) role(name string) palette.Color {
	for _, t := range v.Palette.Roles {
		if t.Name == name {
			return t.Color
		}
	}
	return palette.Color{}
}

type detail struct {
	Token   string `json:"token"`
	Hex     string `json:"hex"`
	OKLCH   string `json:"oklch"`
	Base    string `json:"base"`
	Surface string `json:"surface"`
	Inset   string `json:"inset"`
	DL      string `json:"dl"`
	Note    string `json:"note"`
}

func (d detail) JSON() string {
	b, _ := json.Marshal(d)
	return string(b)
}

func (v View) detail(token string, c palette.Color, note string) detail {
	l := c.LCh()
	ratio := func(g palette.Color) string { return fmt.Sprintf("%.2f:1", palette.Ratio(c, g)) }
	return detail{
		Token:   token,
		Hex:     c.Hex(),
		OKLCH:   fmt.Sprintf("%.3f %.3f %.1f", l.L, l.C, l.H),
		Base:    ratio(v.Seed.Base),
		Surface: ratio(v.role("surface")),
		Inset:   ratio(v.role("inset")),
		DL:      fmt.Sprintf("%.3f", math.Abs(l.L-v.Seed.Base.Lab().L)),
		Note:    note,
	}
}

type swatch struct {
	Token  string
	Family string
	Color  palette.Color
	Detail detail
	Short  bool
}

type row struct {
	Label  string
	Family string
	Cells  []swatch
}

func family(s palette.Scale, beta bool) string {
	if beta {
		return s.Name + "-beta"
	}
	return s.Name
}

func (v View) rows() []row {
	aliases := map[string][]string{}
	for _, s := range v.Palette.Scales {
		if s.Alias != "" {
			aliases[s.Alias] = append(aliases[s.Alias], s.Name)
		}
	}
	var out []row
	add := func(s palette.Scale, label string, steps []palette.Color, beta bool) {
		r := row{Label: label, Family: family(s, beta)}
		for i, c := range steps {
			token := s.Token(i)
			if beta {
				token = s.BetaToken(i)
			}
			note := "meets its tier"
			short := i >= s.Reach
			if short {
				note = "cannot reach its tier on this base; the CSS still emits this colour"
			}
			r.Cells = append(r.Cells, swatch{token, r.Family, c, v.detail(token, c, note), short})
		}
		out = append(out, r)
	}
	for _, s := range v.Palette.Scales {
		if s.Alias != "" {
			continue
		}
		label := s.Name
		if a := aliases[s.Name]; len(a) > 0 {
			label += " · " + strings.Join(a, ", ")
		}
		add(s, label, s.Steps, false)
		if s.Beta != nil {
			add(s, s.Name+" β", s.Beta, true)
		}
	}
	return out
}

func (v View) roles() []swatch {
	var out []swatch
	for _, t := range v.Palette.Roles {
		note := "role"
		if t.Ref != "" {
			note = "alias of " + t.Ref
		}
		out = append(out, swatch{Token: t.Name, Family: "neutral", Color: t.Color, Detail: v.detail(t.Name, t.Color, note)})
	}
	for _, name := range []string{"primary", "secondary", "accent"} {
		s := v.scale(name)
		note := "fill: the " + strconv.Itoa(palette.Steps[palette.FillStep]) + " step"
		if s.Seeded {
			note = "fill: the seed colour"
		}
		out = append(out, swatch{Token: name, Family: name, Color: s.Fill, Detail: v.detail(name, s.Fill, note)})
	}
	return out
}

func (v View) selected() swatch {
	s := v.scale("primary")
	return swatch{Token: "primary", Family: "primary", Color: s.Fill, Detail: v.detail("primary", s.Fill, "fill")}
}

type point struct {
	Hex      string   `json:"hex"`
	Families []string `json:"families"`
	Title    string   `json:"-"`
}

func (v View) cloud() []point {
	var out []point
	seen := map[palette.Color]int{}
	add := func(c palette.Color, fam, token string) {
		if i, ok := seen[c]; ok {
			if !strings.Contains(" "+strings.Join(out[i].Families, " ")+" ", " "+fam+" ") {
				out[i].Families = append(out[i].Families, fam)
			}
			out[i].Title += ", " + token
			return
		}
		seen[c] = len(out)
		out = append(out, point{c.Hex(), []string{fam}, c.Hex() + " " + token})
	}
	for _, t := range v.Palette.Roles {
		add(t.Color, "neutral", t.Name)
	}
	for _, s := range v.Palette.Scales {
		if s.Alias != "" {
			continue
		}
		if s.Seeded {
			add(s.Fill, s.Name, s.Name)
		}
		for i, c := range s.Steps {
			if i < s.Reach {
				add(c, s.Name, s.Token(i))
			}
		}
		for i, c := range s.Beta {
			if i < s.Reach {
				add(c, family(s, true), s.BetaToken(i))
			}
		}
	}
	return out
}

type dot struct {
	X, Y     float64
	Color    palette.Color
	Families string
	Title    string
}

const (
	hexRadius  = 92
	hexChroma  = 0.33
	dotRadius  = 4
	ringChroma = 0.1
)

func (v View) dots() []dot {
	var out []dot
	for _, p := range v.cloud() {
		c, _ := palette.ParseHex(p.Hex)
		l := c.LCh()
		r := min(l.C/hexChroma, 1) * hexRadius
		a := l.H * math.Pi / 180
		out = append(out, dot{float64(r * math.Cos(a)), -float64(r * math.Sin(a)), c, strings.Join(p.Families, " "), p.Title})
	}
	return out
}

type ring struct {
	R, X, Y float64
	Label   string
}

func rings(dots []dot) []ring {
	spoke, fewest := 0.0, len(dots)+1
	for k := range 6 {
		a := float64(k) * math.Pi / 3
		n := 0
		for _, d := range dots {
			if math.Hypot(d.X, d.Y) > ringChroma/hexChroma*hexRadius/2 && math.Abs(math.Remainder(math.Atan2(d.Y, d.X)-a, 2*math.Pi)) < math.Pi/6 {
				n++
			}
		}
		if n < fewest {
			spoke, fewest = a, n
		}
	}
	var out []ring
	for c := ringChroma; c <= hexChroma; c += ringChroma {
		r := c / hexChroma * hexRadius
		out = append(out, ring{r, r * math.Cos(spoke), r * math.Sin(spoke), fmt.Sprintf("C %.1f", c)})
	}
	return out
}

func hexagon(x, y, r float64) string {
	var b strings.Builder
	for i := range 6 {
		a := float64(i) * math.Pi / 3
		fmt.Fprintf(&b, "%.2f,%.2f ", x+float64(r*math.Cos(a)), y+float64(r*math.Sin(a)))
	}
	return strings.TrimSpace(b.String())
}

func (v View) space() string {
	b, _ := json.Marshal(map[string]any{"base": v.Seed.Base.Hex(), "ink": v.role("muted").Hex(), "points": v.cloud()})
	return string(b)
}

const (
	labelW  = 116.0
	cellW   = 40.0
	cellH   = 28.0
	cellGap = 4.0
	rowH    = 32
	gridTop = 26
)

var gridW = int(labelW + float64(len(palette.Steps))*(cellW+cellGap) + 8)

func colX(i int) float64 { return labelW + float64(i)*(cellW+cellGap) }

func rowY(r int) float64 { return float64(gridTop + r*rowH) }

func gridH(rows int) int { return gridTop + rows*rowH + 8 }

func num(f float64) string { return strconv.FormatFloat(f, 'f', 2, 64) }

func roleY(i int) float64 { return float64(8 + i*30) }

func rolesH(n int) int { return 8 + n*30 + 4 }
