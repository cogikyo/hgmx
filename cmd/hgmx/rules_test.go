package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRules(t *testing.T) {
	wants, err := filepath.Glob("testdata/*/want")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range wants {
		dir := filepath.Dir(want)
		t.Run(filepath.Base(dir), func(t *testing.T) {
			t.Parallel()
			s, err := loadSite(dir)
			if err != nil {
				t.Fatal(err)
			}
			findings, err := runRules(s)
			if err != nil {
				t.Fatal(err)
			}
			var got strings.Builder
			for _, f := range findings {
				got.WriteString(f.String() + "\n")
			}
			data, err := os.ReadFile(want)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != string(data) {
				t.Errorf("findings:\n%s\nwant:\n%s", got.String(), data)
			}
		})
	}
}
