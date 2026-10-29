package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/cogikyo/hgmx/icons"
	"github.com/cogikyo/hgmx/scripts"
)

const css = `body { font: 14px system-ui, sans-serif; margin: 1rem; }
#icons ul { display: grid; grid-template-columns: repeat(auto-fill, minmax(6rem, 1fr)); gap: .5rem; list-style: none; padding: 0; }
#icons li { display: flex; flex-direction: column; align-items: center; gap: .4rem; padding: .75rem .25rem; font-size: 4rem; }
#icons span { font-size: .7rem; color: #666; }
`

func main() {
	addr := flag.String("addr", "127.0.0.1:18493", "listen address")
	flag.Parse()

	js, err := scripts.New(scripts.HTMX, scripts.Hyperscript, scripts.Motion)
	if err != nil {
		log.Fatal(err)
	}
	ic, err := icons.New()
	if err != nil {
		log.Fatal(err)
	}
	head := js.Head() + "\n" + ic.Head()

	mux := http.NewServeMux()
	mux.Handle(scripts.Path, js)
	mux.Handle(icons.Path, ic)
	mux.HandleFunc("GET /check.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		fmt.Fprint(w, css)
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	for _, path := range []string{"/{$}", "/a", "/b"} {
		mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			var only []icons.Name
			if v := q.Get("only"); v != "" {
				for n := range strings.SplitSeq(v, ",") {
					if !slices.Contains(icons.All(), icons.Name(n)) {
						http.Error(w, "no icon "+n, http.StatusNotFound)
						return
					}
					only = append(only, icons.Name(n))
				}
			}
			w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; connect-src 'self'")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := page(head, r.URL.Path, q.Get("q"), only).Render(r.Context(), w); err != nil {
				log.Print(err)
			}
		})
	}

	srv := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("icons gallery on http://%s", *addr)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	srv.Close()
}

func query(q string, only []icons.Name) string {
	v := url.Values{}
	if q != "" {
		v.Set("q", q)
	}
	if only != nil {
		names := make([]string, len(only))
		for i, n := range only {
			names[i] = string(n)
		}
		v.Set("only", strings.Join(names, ","))
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}
