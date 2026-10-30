package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type entry struct {
	Status  int
	Bytes   int64
	Aborted bool
}

func TestAccess(t *testing.T) {
	for _, tc := range []struct {
		name string
		h    http.HandlerFunc
		want entry
	}{
		{"implicit", func(w http.ResponseWriter, r *http.Request) {}, entry{Status: 200}},
		{"write", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "abc") }, entry{Status: 200, Bytes: 3}},
		{"explicit", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) }, entry{Status: 404}},
		{"first final", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusEarlyHints)
			w.WriteHeader(201)
			w.WriteHeader(500)
		}, entry{Status: 201}},
		{"switching", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusSwitchingProtocols) }, entry{Status: 101}},
		{"flush", func(w http.ResponseWriter, r *http.Request) {
			if err := http.NewResponseController(w).Flush(); err != nil {
				panic(err)
			}
		}, entry{Status: 200}},
		{"aborted after commit", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "ab")
			panic(http.ErrAbortHandler)
		}, entry{Status: 200, Bytes: 2, Aborted: true}},
		{"panic before commit", func(w http.ResponseWriter, r *http.Request) { panic("boom") }, entry{Status: 500, Aborted: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := offline(nil)
			var buf bytes.Buffer
			p.local = slog.NewJSONHandler(&buf, nil)
			func() {
				defer func() { recover() }()
				p.Handler(tc.h).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
			}()
			var got entry
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("access line %q: %v", buf.String(), err)
			}
			if got != tc.want {
				t.Fatalf("access line %+v, want %+v", got, tc.want)
			}
			spans := p.spans.take()
			if len(spans) != 1 {
				t.Fatalf("%d spans", len(spans))
			}
			s := spans[0]
			if code := attr(s.attrs, "http.response.status_code"); code != strconv.Itoa(tc.want.Status) {
				t.Errorf("span status %s, access line %d", code, tc.want.Status)
			}
			if failed := s.code == 2; failed != (tc.want.Status >= 500 || tc.want.Aborted) {
				t.Errorf("span failed = %v for %+v", failed, tc.want)
			}
		})
	}
}

var errGone = errors.New("client gone")

type failing struct{ *httptest.ResponseRecorder }

func (failing) FlushError() error { return errGone }

func TestFlushError(t *testing.T) {
	for _, tc := range []struct {
		w    http.ResponseWriter
		want error
	}{
		{failing{httptest.NewRecorder()}, errGone},
		{struct{ http.ResponseWriter }{httptest.NewRecorder()}, http.ErrNotSupported},
	} {
		var got error
		Disabled(slog.DiscardHandler).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got = http.NewResponseController(w).Flush()
		})).ServeHTTP(tc.w, httptest.NewRequest("GET", "/", nil))
		if !errors.Is(got, tc.want) {
			t.Errorf("Flush = %v, want %v", got, tc.want)
		}
	}
}
