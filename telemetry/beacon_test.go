package telemetry

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const chrome = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0 Safari/537.36"

func post(b *Beacon, body string) int {
	r := httptest.NewRequest("POST", Path, strings.NewReader(body))
	r.Host = "hgmx.dev"
	r.RemoteAddr = "198.51.100.9:41234"
	r.Header.Set("User-Agent", chrome)
	r.Header.Set("CF-Connecting-IP", "203.0.113.77")
	r.Header.Set("CF-IPCountry", "NZ")
	w := httptest.NewRecorder()
	b.ServeHTTP(w, r)
	return w.Code
}

func view(size int) string {
	body := `{"type":"view","path":"/icons?q=secret","referrer":"https://news.ycombinator.com/item?id=1","traceparent":"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01","x":"`
	return body + strings.Repeat("a", size-len(body)-2) + `"}`
}

func TestBeaconCaps(t *testing.T) {
	b := offline(nil).Beacon()
	if code := post(b, view(maxBeacon+1)); code != http.StatusRequestEntityTooLarge {
		t.Errorf("2 KiB + 1: %d, want 413", code)
	}
	if code := post(b, view(maxBeacon)); code != http.StatusBadRequest {
		t.Errorf("2 KiB with an unknown field: %d, want 400", code)
	}

	b.limit = newBucket(beaconsPS)
	codes := map[int]int{}
	for range 200 {
		codes[post(b, `{"type":"view","path":"/"}`)]++
	}
	if codes[http.StatusNoContent] < beaconsPS || codes[http.StatusNoContent] > beaconsPS+5 || codes[http.StatusTooManyRequests] == 0 {
		t.Errorf("200 posts at once: %v; want about %d accepted", codes, beaconsPS)
	}
}

func TestBeaconPrivacy(t *testing.T) {
	p := offline(nil)
	b := p.Beacon()
	if code := post(b, `{"type":"view","path":"/icons?q=secret#top","referrer":"https://news.ycombinator.com/item?id=1","traceparent":"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"}`); code != http.StatusNoContent {
		t.Fatalf("view: %d", code)
	}
	if code := post(b, `{"type":"vitals","path":"/icons","ttfb":12,"lcp":420,"cls":0.02,"inp":1e9}`); code != http.StatusNoContent {
		t.Fatalf("vitals: %d", code)
	}
	for _, ref := range []string{"http://192.0.2.9/private", "http://[2001:db8::1]:8080/"} {
		if code := post(b, `{"type":"view","path":"/","referrer":"`+ref+`"}`); code != http.StatusNoContent {
			t.Fatalf("view from %s: %d", ref, code)
		}
	}
	recs := p.recs.take()
	for _, r := range recs[2:] {
		if h := attr(r.attrs, "referrer.host"); h != "" {
			t.Errorf("IP-literal referrer exported as %q", h)
		}
	}
	recs = recs[:2]
	if len(recs) != 2 {
		t.Fatalf("%d records", len(recs))
	}
	out := encodeLogs(nil, nil, recs)
	for _, leak := range []string{"203.0.113.77", "198.51.100.9", "secret", "item?id", chrome} {
		if bytes.Contains(out, []byte(leak)) {
			t.Errorf("exported logs contain %q", leak)
		}
	}
	v := recs[0]
	if v.event != "browser.page_view" || attr(v.attrs, "url.path") != "/icons" || attr(v.attrs, "referrer.host") != "news.ycombinator.com" ||
		len(attr(v.attrs, "visitor.hash")) != 16 || attr(v.attrs, "geo.country.iso_code") != "NZ" || v.trace == (traceID{}) {
		t.Errorf("page view: %+v", v)
	}
	if attr(recs[1].attrs, "browser.vitals.inp") != "" || attr(recs[1].attrs, "browser.vitals.lcp") != "420" {
		t.Errorf("vitals: %+v", recs[1].attrs)
	}
}

func TestBeaconCrossOrigin(t *testing.T) {
	for _, p := range []*Pipeline{offline(nil), Disabled(slog.DiscardHandler)} {
		b := p.Beacon()
		if (b.Head() == "") != p.off() {
			t.Errorf("off=%v: Head %q", p.off(), b.Head())
		}
		for site, want := range map[string]int{"same-origin": http.StatusNoContent, "cross-site": http.StatusForbidden} {
			r := httptest.NewRequest("POST", Path, strings.NewReader(`{"type":"view","path":"/"}`))
			r.Header.Set("User-Agent", chrome)
			r.Header.Set("Sec-Fetch-Site", site)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, r)
			if w.Code != want {
				t.Errorf("off=%v, %s: %d, want %d", p.off(), site, w.Code, want)
			}
		}
		if p.recs != nil {
			if n := len(p.recs.take()); n != 1 {
				t.Errorf("%d events accepted, want the same-origin one", n)
			}
		}
	}
}
