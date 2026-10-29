package scripts

import (
	"regexp"
	"slices"
	"testing"
)

func TestLoadOrder(t *testing.T) {
	h := New(Motion, Hyperscript, Navigation)
	var got []string
	for _, m := range regexp.MustCompile(`/([^/"]+\.js)" defer>`).FindAllStringSubmatch(string(h.Head()), -1) {
		got = append(got, m[1])
	}
	want := []string{"htmx.min.js", "hx-preload.js", "hgmx.js", "hyperscript.min.js", "motion.min.js"}
	if !slices.Equal(got, want) {
		t.Fatalf("tags load %v, want %v", got, want)
	}
}
