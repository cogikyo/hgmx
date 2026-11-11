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
			wantErr, _ := os.ReadFile(filepath.Join(dir, "error"))
			switch {
			case len(wantErr) == 0 && err != nil:
				t.Fatal(err)
			case len(wantErr) > 0 && (err == nil || !strings.Contains(err.Error(), strings.TrimSpace(string(wantErr)))):
				t.Fatalf("error: %v; want one containing %q", err, strings.TrimSpace(string(wantErr)))
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
