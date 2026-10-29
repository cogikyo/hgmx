package scripts

import (
	"regexp"
	"slices"
	"testing"
)

func TestLoadOrder(t *testing.T) {
	h, err := New(Motion, Hyperscript, Preload)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, m := range regexp.MustCompile(`/([^/"]+\.js)" defer>`).FindAllStringSubmatch(string(h.Head()), -1) {
		got = append(got, m[1])
	}
	want := []string{"htmx.min.js", "hx-preload.js", "hyperscript.min.js", "motion.min.js"}
	if !slices.Equal(got, want) {
		t.Fatalf("tags load %v, want %v", got, want)
	}
}
