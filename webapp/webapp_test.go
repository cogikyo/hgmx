package webapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"path"
	"regexp"
	"testing"
)

func TestInstallable(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for i := range src.Pix {
		src.Pix[i] = 0xff
	}
	src.Set(256, 256, color.RGBA{R: 0xff, A: 0xff})
	var icon bytes.Buffer
	if err := png.Encode(&icon, src); err != nil {
		t.Fatal(err)
	}
	h, err := New(App{Name: "test", Theme: "#222536", Background: "#222536", Icon: icon.Bytes()})
	if err != nil {
		t.Fatal(err)
	}

	link := regexp.MustCompile(`rel="manifest" href="([^"]+)"`).FindStringSubmatch(string(h.Head()))
	if link == nil {
		t.Fatalf("Head has no manifest link: %s", h.Head())
	}
	res := serve(h, link[1])
	if res.Code != http.StatusOK {
		t.Fatalf("manifest %s: %d", link[1], res.Code)
	}
	var m struct {
		Name     string `json:"name"`
		StartURL string `json:"start_url"`
		Display  string `json:"display"`
		Icons    []struct{ Src, Sizes, Purpose string }
	}
	if err := json.Unmarshal(res.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	if m.Name == "" || m.StartURL == "" || m.Display != "standalone" {
		t.Fatalf("manifest misses an install field: %s", res.Body)
	}

	sizes := map[string]bool{}
	for _, ic := range m.Icons {
		url := path.Join(path.Dir(link[1]), ic.Src)
		res := serve(h, url)
		img, err := png.Decode(res.Body)
		if err != nil {
			t.Fatalf("icon %s: %d %v", url, res.Code, err)
		}
		if got := fmt.Sprintf("%dx%d", img.Bounds().Dx(), img.Bounds().Dy()); got != ic.Sizes {
			t.Fatalf("icon %s is %s, manifest says %s", url, got, ic.Sizes)
		}
		sizes[ic.Sizes] = true
	}
	if !sizes["192x192"] || !sizes["512x512"] {
		t.Fatalf("manifest icons cover %v, install needs 192x192 and 512x512", sizes)
	}
}

func serve(h http.Handler, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, url, nil))
	return w
}
