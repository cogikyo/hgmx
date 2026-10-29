package telemetry

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cogikyo/hgmx/static"
)

const (
	Path      = "/_v/" // Path is the beacon and asset mount prefix.
	maxBeacon = 2 << 10
	beaconsPS = 50
)

//go:embed vitals.js
var vitals []byte

// Beacon serves browser page-view and web-vitals collection and its script.
// Mount with root.Handle(telemetry.Path, beacon), without a method or StripPrefix.
type Beacon struct {
	*static.Server
	head  template.HTML
	p     *Pipeline
	limit *bucket

	mu   sync.Mutex
	day  int64
	salt [32]byte
}

// Beacon prepares the browser collector, including when the pipeline is nil.
func (p *Pipeline) Beacon() (*Beacon, error) {
	b := &Beacon{Server: static.Map(Path, map[string][]byte{"vitals.js": vitals}), p: p, limit: newBucket(beaconsPS)}
	if p == nil {
		return b, nil
	}
	src, err := b.URL("vitals.js")
	if err != nil {
		return nil, err
	}
	b.head = template.HTML(fmt.Sprintf(`<script type="module" src="%s"></script>`, src))
	return b, nil
}

// Head returns the layout's script tag, or empty HTML for a nil pipeline.
// The CSP must allow script-src 'self' and connect-src 'self'.
func (b *Beacon) Head() template.HTML {
	return b.head
}

func (b *Beacon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		b.Server.ServeHTTP(w, r)
		return
	}
	if r.URL.Path != Path {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(b.accept(r))
}

type event struct {
	Type        string   `json:"type"`
	Path        string   `json:"path"`
	Referrer    string   `json:"referrer"`
	Traceparent string   `json:"traceparent"`
	TTFB        *float64 `json:"ttfb"`
	FCP         *float64 `json:"fcp"`
	LCP         *float64 `json:"lcp"`
	INP         *float64 `json:"inp"`
	CLS         *float64 `json:"cls"`
}

func (b *Beacon) accept(r *http.Request) int {
	if b.p.off() {
		return http.StatusNoContent
	}
	if r.ContentLength > maxBeacon {
		return http.StatusRequestEntityTooLarge
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBeacon+1))
	if err != nil {
		return http.StatusBadRequest
	}
	if len(body) > maxBeacon {
		return http.StatusRequestEntityTooLarge
	}
	if !b.limit.take() {
		return http.StatusTooManyRequests
	}
	ua := r.UserAgent()
	if bot(ua) {
		return http.StatusNoContent
	}

	var e event
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&e); err != nil {
		return http.StatusBadRequest
	}
	path, _, _ := strings.Cut(e.Path, "?")
	path, _, _ = strings.Cut(path, "#")
	if !strings.HasPrefix(path, "/") || len(path) > 512 {
		return http.StatusBadRequest
	}

	now := time.Now().UnixNano()
	rec := record{time: now, observed: now, level: slog.LevelInfo}
	rec.attrs = append(rec.attrs, slog.String("url.path", clip(path)))
	switch e.Type {
	case "view":
		rec.event = "browser.page_view"
		if host := referrer(e.Referrer); host != "" && host != strings.ToLower(r.Host) {
			rec.attrs = append(rec.attrs, slog.String("referrer.host", host))
		}
		rec.attrs = append(rec.attrs, slog.String("visitor.hash", b.visitor(r, ua)))
	case "vitals":
		rec.event = "browser.vitals"
		for _, v := range []struct {
			key string
			val *float64
			max float64
		}{
			{"browser.vitals.ttfb", e.TTFB, 60_000},
			{"browser.vitals.fcp", e.FCP, 60_000},
			{"browser.vitals.lcp", e.LCP, 60_000},
			{"browser.vitals.inp", e.INP, 60_000},
			{"browser.vitals.cls", e.CLS, 10},
		} {
			if v.val != nil && *v.val >= 0 && *v.val <= v.max && !math.IsNaN(*v.val) {
				rec.attrs = append(rec.attrs, slog.Float64(v.key, *v.val))
			}
		}
	default:
		return http.StatusBadRequest
	}
	rec.msg = rec.event
	rec.attrs = append(rec.attrs,
		slog.String("user_agent.name", family(ua)),
		slog.Bool("browser.mobile", strings.Contains(ua, "Mobi")),
	)
	if c := r.Header.Get("CF-IPCountry"); len(c) == 2 && c != "XX" && c != "T1" && upper(c) {
		rec.attrs = append(rec.attrs, slog.String("geo.country.iso_code", strings.Clone(c)))
	}
	if l, ok := parseTraceparent(e.Traceparent); ok {
		rec.trace, rec.span = l.trace, l.id
	}
	b.p.recs.add(rec)
	return http.StatusNoContent
}

func (b *Beacon) visitor(r *http.Request, ua string) string {
	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	day := time.Now().UTC().Unix() / 86400
	b.mu.Lock()
	if day != b.day {
		b.day = day
		rand.Read(b.salt[:])
	}
	h := sha256.New()
	h.Write(b.salt[:])
	b.mu.Unlock()
	fmt.Fprintf(h, "%s\x00%s\x00%s", strings.ToLower(r.Host), ip, ua)
	return hex.EncodeToString(h.Sum(nil)[:8])
}

func referrer(ref string) string {
	u, err := url.Parse(ref)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || len(u.Host) > 253 || net.ParseIP(u.Hostname()) != nil {
		return ""
	}
	return strings.Clone(strings.ToLower(u.Host))
}

func bot(ua string) bool {
	if ua == "" {
		return true
	}
	ua = strings.ToLower(ua)
	for _, s := range []string{"bot", "crawl", "spider", "headless", "lighthouse"} {
		if strings.Contains(ua, s) {
			return true
		}
	}
	return false
}

func family(ua string) string {
	for _, f := range []struct{ token, name string }{
		{"Edg/", "Edge"},
		{"OPR/", "Opera"},
		{"Firefox/", "Firefox"},
		{"Chrome/", "Chrome"},
		{"CriOS/", "Chrome"},
		{"Safari/", "Safari"},
	} {
		if strings.Contains(ua, f.token) {
			return f.name
		}
	}
	return "other"
}

func upper(s string) bool {
	for i := range len(s) {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}
