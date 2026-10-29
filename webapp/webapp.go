// Package webapp builds an app manifest and icons from one source PNG.
// The manifest fixes id, start_url and scope at /.
// Its URL changes with the manifest or icons; the fixed id keeps app identity.
package webapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"math"
	"strconv"

	"github.com/cogikyo/hgmx/static"
	"golang.org/x/image/draw"
)

// Path is the mount prefix; keep it in the request path.
const Path = "/app/"

const (
	minIcon  = 512
	safeZone = 0.4
	touchBox = 0.4
	faint    = 8
)

// App supplies manifest metadata and the source icon, not offline support.
type App struct {
	// Name is required.
	Name string
	// ShortName is optional.
	ShortName string
	// Description is optional.
	Description string
	// Theme is a required #RRGGBB colour.
	Theme string
	// Background is a required #RRGGBB fill for maskable and touch icons.
	Background string
	// Icon is a square PNG of at least 512 px per side.
	Icon []byte
}

// Handler serves the manifest and icons under [Path].
// Mount it at [Path] without http.StripPrefix.
type Handler struct {
	*static.Server
	head template.HTML
}

type manifest struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ShortName       string `json:"short_name,omitempty"`
	Description     string `json:"description,omitempty"`
	StartURL        string `json:"start_url"`
	Scope           string `json:"scope"`
	Display         string `json:"display"`
	ThemeColor      string `json:"theme_color"`
	BackgroundColor string `json:"background_color"`
	Icons           []icon `json:"icons"`
}

type icon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
}

// New checks the name, colours and source PNG before generating files.
func New(a App) (*Handler, error) {
	if a.Name == "" {
		return nil, errors.New("webapp: name is required")
	}
	if _, err := parseHex(a.Theme); err != nil {
		return nil, fmt.Errorf("webapp: theme: %w", err)
	}
	bg, err := parseHex(a.Background)
	if err != nil {
		return nil, fmt.Errorf("webapp: background: %w", err)
	}
	src, err := decode(a.Icon)
	if err != nil {
		return nil, fmt.Errorf("webapp: icon: %w", err)
	}
	round, square := extent(src, bg)

	m := manifest{
		ID:              "/",
		Name:            a.Name,
		ShortName:       a.ShortName,
		Description:     a.Description,
		StartURL:        "/",
		Scope:           "/",
		Display:         "standalone",
		ThemeColor:      a.Theme,
		BackgroundColor: a.Background,
	}
	images := map[string]image.Image{"apple-touch-icon.png": inset(src, 180, bg, touchBox/square)}
	for _, size := range []int{192, 512} {
		plain := fmt.Sprintf("icon-%d.png", size)
		masked := fmt.Sprintf("maskable-%d.png", size)
		images[plain] = scale(src, size)
		images[masked] = inset(src, size, bg, safeZone/round)
		dim := fmt.Sprintf("%dx%d", size, size)
		m.Icons = append(m.Icons,
			icon{Src: plain, Sizes: dim, Type: "image/png", Purpose: "any"},
			icon{Src: masked, Sizes: dim, Type: "image/png", Purpose: "maskable"},
		)
	}

	files := make(map[string][]byte, len(images)+1)
	for name, img := range images {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("webapp: %s: %w", name, err)
		}
		files[name] = buf.Bytes()
	}
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("webapp: manifest: %w", err)
	}
	files["manifest.webmanifest"] = body

	h := &Handler{Server: static.Map(Path, files)}
	link, err := h.URL("manifest.webmanifest")
	if err != nil {
		return nil, err
	}
	touch, err := h.URL("apple-touch-icon.png")
	if err != nil {
		return nil, err
	}
	h.head = template.HTML(fmt.Sprintf(
		"<link rel=\"manifest\" href=\"%s\">\n<link rel=\"apple-touch-icon\" sizes=\"180x180\" href=\"%s\">\n<meta name=\"theme-color\" content=\"%s\">",
		link, touch, a.Theme,
	))
	return h, nil
}

// Head returns manifest, apple-touch-icon and theme-colour tags for the head.
func (h *Handler) Head() template.HTML {
	return h.head
}

func parseHex(s string) (color.RGBA, error) {
	if len(s) != 7 || s[0] != '#' {
		return color.RGBA{}, fmt.Errorf("%q is not a #RRGGBB colour", s)
	}
	v, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("%q is not a #RRGGBB colour", s)
	}
	return color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 255}, nil
}

func decode(b []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("not a PNG: %w", err)
	}
	r := img.Bounds()
	if r.Dx() != r.Dy() || r.Dx() < minIcon {
		return nil, fmt.Errorf("want a square PNG of at least %dpx, got %dx%d", minIcon, r.Dx(), r.Dy())
	}
	dst := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	draw.Draw(dst, dst.Bounds(), img, r.Min, draw.Src)
	return dst, nil
}

func extent(src *image.RGBA, bg color.RGBA) (round, square float64) {
	n := src.Bounds().Dx()
	c := float64(n) / 2
	reach := func(i int) float64 {
		return math.Max(math.Abs(float64(i)-c), math.Abs(float64(i+1)-c))
	}
	for y := range n {
		dy := reach(y)
		for x := range n {
			if !visible(src.Pix[y*src.Stride+x*4:][:4], bg) {
				continue
			}
			dx := reach(x)
			round = math.Max(round, math.Hypot(dx, dy))
			square = math.Max(square, math.Max(dx, dy))
		}
	}
	return round / float64(n), square / float64(n)
}

func visible(p []uint8, bg color.RGBA) bool {
	keep := 255 - int(p[3])
	for i, b := range [3]uint8{bg.R, bg.G, bg.B} {
		over := int(p[i]) + int(b)*keep/255
		if over-int(b) > faint || int(b)-over > faint {
			return true
		}
	}
	return false
}

func scale(src *image.RGBA, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return dst
}

func inset(src *image.RGBA, size int, bg color.RGBA, fill float64) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)
	side := int(float64(size) * math.Min(fill, 1))
	side -= (size - side) % 2
	off := (size - side) / 2
	draw.BiLinear.Scale(dst, image.Rect(off, off, off+side, off+side), src, src.Bounds(), draw.Over, nil)
	return dst
}
