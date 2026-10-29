package theme

import (
	"math"

	"github.com/cogikyo/hgmx/palette"
)

type mock struct {
	Base, Inset, Surface, Raised, Line     palette.Color
	Strong, Fg, Muted, Disabled, Border    palette.Color
	Hover, Pressed, Selected, OnSelected   palette.Color
	Focus, Link, LinkHover                 palette.Color
	Primary, OnPrimary, Interface, White   palette.Color
	Secondary, OnSecondary, Error, OnError palette.Color
	Chips                                  []chip
}

type chip struct {
	Name, Glyph    string
	Bg, Edge, Text palette.Color
}

type pair struct {
	Label  string
	Fg, Bg palette.Color
	Min    float64
	Scale  string
	Step   int
}

var (
	white = palette.Color{R: 255, G: 255, B: 255}
	black = palette.Color{}
)

const (
	hoverAlpha   = 0.05
	pressedAlpha = 0.10
)

func (v View) mock() mock {
	pri, sec, acc, errs, neu := v.scale("primary"), v.scale("secondary"), v.scale("accent"), v.scale("error"), v.scale("neutral")
	surface, ink := v.role("surface"), v.role("ink")
	m := mock{
		Base: v.role("base"), Inset: v.role("inset"), Surface: surface, Raised: v.role("raised"), Line: v.role("line"),
		Strong: v.role("strong"), Fg: v.role("fg"), Muted: v.role("muted"), Disabled: neu.Steps[4], Border: neu.Steps[3],
		Hover: blend(ink, surface, hoverAlpha), Pressed: blend(ink, surface, pressedAlpha),
		Selected: acc.Steps[1], OnSelected: acc.Steps[7], Focus: acc.Steps[palette.FillStep],
		Link: pri.Steps[palette.TextStep], LinkHover: pri.Beta[palette.TextStep],
		Primary: pri.Fill, OnPrimary: pri.On, Interface: darken(pri.Fill), White: white,
		Secondary: sec.Fill, OnSecondary: sec.On, Error: errs.Fill, OnError: errs.On,
	}
	for _, st := range []struct{ name, glyph string }{{"info", "i"}, {"success", "✓"}, {"warning", "!"}, {"error", "✕"}} {
		s := v.scale(st.name)
		m.Chips = append(m.Chips, chip{st.name, st.glyph, s.Steps[1], s.Steps[3], s.Steps[7]})
	}
	return m
}

func (m mock) pairs() []pair {
	out := []pair{
		{"title on inset", m.Strong, m.Inset, 4.5, "neutral", 9},
		{"title on surface", m.Strong, m.Surface, 4.5, "neutral", 9},
		{"body on surface", m.Fg, m.Surface, 4.5, "neutral", 8},
		{"muted on surface", m.Muted, m.Surface, 4.5, "neutral", palette.TextStep},
		{"placeholder on base", m.Muted, m.Base, 4.5, "neutral", palette.TextStep},
		{"rest on surface", m.Fg, m.Surface, 4.5, "neutral", 8},
		{"hover", m.Fg, m.Hover, 4.5, "neutral", 8},
		{"pressed", m.Fg, m.Pressed, 4.5, "neutral", 8},
		{"selected", m.OnSelected, m.Selected, 4.5, "accent", 7},
		{"focus ring on surface", m.Focus, m.Surface, 3, "accent", palette.FillStep},
		{"link on surface", m.Link, m.Surface, 4.5, "primary", palette.TextStep},
		{"link hover on surface", m.LinkHover, m.Surface, 4.5, "primary", palette.TextStep},
		{"body on raised", m.Fg, m.Raised, 4.5, "neutral", 8},
		{"primary button", m.OnPrimary, m.Primary, 4.5, "", 0},
		{"interface fill", m.White, m.Interface, 4.5, "", 0},
		{"secondary button", m.OnSecondary, m.Secondary, 4.5, "", 0},
		{"destructive button", m.OnError, m.Error, 4.5, "", 0},
	}
	for _, c := range m.Chips {
		out = append(out, pair{c.Name + " chip", c.Text, c.Bg, 4.5, c.Name, 7})
	}
	return out
}

func blend(top, bottom palette.Color, alpha float64) palette.Color {
	mix := func(a, b uint8) uint8 {
		return uint8(math.Round(float64(alpha*float64(a)) + float64((1-alpha)*float64(b))))
	}
	return palette.Color{R: mix(top.R, bottom.R), G: mix(top.G, bottom.G), B: mix(top.B, bottom.B)}
}

func darken(c palette.Color) palette.Color {
	l := c.LCh()
	for i := range 101 {
		k := 1 - float64(i)/100
		x := palette.LCh{L: l.L * k, C: l.C * k, H: l.H}.Color()
		if palette.Ratio(white, x) >= 4.5 {
			return x
		}
	}
	return black
}
