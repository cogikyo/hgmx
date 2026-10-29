package render

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func HTML(w http.ResponseWriter, r *http.Request, status int, c templ.Component) error {
	var buf bytes.Buffer
	if err := c.Render(r.Context(), &buf); err != nil {
		return err
	}
	body := buf.Bytes()

	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	if h.Get("Cache-Control") == "" {
		h.Set("Cache-Control", "no-cache")
	}
	read := r.Method == http.MethodGet || r.Method == http.MethodHead
	if status == http.StatusOK && read {
		if h.Get("ETag") == "" {
			sum := fnv.New64a()
			sum.Write(body)
			h.Set("ETag", fmt.Sprintf(`"%016x"`, sum.Sum64()))
		}
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(body))
		return nil
	}

	h.Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		w.Write(body)
	}
	return nil
}
