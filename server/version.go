package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

const versionHeader = "HGMX-Version"

func VersionConfig(version string) string {
	b, _ := json.Marshal(map[string]map[string]string{"headers": {versionHeader: version}})
	return string(b)
}

func PageVersion(version string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(versionHeader)
		read := r.Method == http.MethodGet || r.Method == http.MethodHead
		if !read || got == "" || got == version {
			h.ServeHTTP(w, r)
			return
		}
		target := r.URL.RequestURI()
		if strings.HasPrefix(target, "//") {
			target = "/"
		}
		w.Header().Set("HX-Redirect", target)
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusNoContent)
	})
}
