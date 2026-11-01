package main

import (
	"context"
	"errors"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const hgmxRefLines = "294690d89e1aea47982c8be93bfb4a011b69b363\trefs/heads/master\n" +
	"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee\trefs/tags/v0.1.0\n" +
	"294690d89e1aea47982c8be93bfb4a011b69b363\trefs/tags/v0.1.0^{}\n" +
	"294690d89e1aea47982c8be93bfb4a011b69b363\trefs/tags/v0.1.1\n" +
	"edcd391fa47d0000000000000000000000000000\trefs/tags/v0.0.33\n" +
	"294690d89e1aea47982c8be93bfb4a011b69b363\trefs/tags/v0.2.0\n" +
	"edcd391fa47d0000000000000000000000000000\trefs/tags/v0.2.0^{}\n"

func TestParseHgmxRefs(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(hgmxRefLines), "\n")
	slices.Reverse(lines)
	for _, input := range []string{hgmxRefLines, strings.Join(lines, "\n")} {
		refs, err := parseHgmxRefs(input)
		if err != nil {
			t.Fatal(err)
		}
		if refs.master != "294690d89e1aea47982c8be93bfb4a011b69b363" {
			t.Fatalf("master: %s", refs.master)
		}
		want := map[string]string{
			"v0.1.0":  refs.master,
			"v0.1.1":  refs.master,
			"v0.0.33": "edcd391fa47d0000000000000000000000000000",
			"v0.2.0":  "edcd391fa47d0000000000000000000000000000",
		}
		if !maps.Equal(refs.tags, want) {
			t.Fatalf("tags: %v; want %v", refs.tags, want)
		}
	}
	for _, input := range []string{"", "not-a-ref", strings.ReplaceAll(hgmxRefLines, "294690", "zzzzzz"), strings.ReplaceAll(hgmxRefLines, "refs/heads/master", "refs/heads/other")} {
		if _, err := parseHgmxRefs(input); err == nil {
			t.Fatalf("invalid refs accepted: %q", input)
		}
	}
}

func TestHgmxLatest(t *testing.T) {
	const latest = "294690d89e1a"
	const current = "v0.0.36-0.20261101201201-" + latest
	refs, err := parseHgmxRefs(hgmxRefLines)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name    string
		version string
		vendor  string
		want    string
	}{
		{"current", current, "# " + hgmxPath + " " + current + "\n## explicit; go 1.27.0\n" + hgmxTool + "\n", ""},
		{"same-commit", "v0.0.0-20261101201201-" + latest, "", ""},
		{"annotated-tag", "v0.1.0", "# " + hgmxPath + " v0.1.0\n", ""},
		{"lightweight-tag", "v0.1.1", "# " + hgmxPath + " v0.1.1\n", ""},
		{"old-annotated-tag", "v0.2.0", "", "go.mod [hgmx-latest] site requires hgmx v0.2.0; master resolves to " + latest + "\n"},
		{"unknown-tag", "v0.3.0", "", "go.mod [hgmx-latest] site requires hgmx v0.3.0; master resolves to " + latest + "\n"},
		{"no-vendor", current, "", ""},
		{"old-module", "v0.0.33", "", "go.mod [hgmx-latest] site requires hgmx v0.0.33; master resolves to " + latest + "\n"},
		{"old-commit", "v0.0.36-0.20261101193023-edcd391fa47d", "", "go.mod [hgmx-latest] site requires hgmx v0.0.36-0.20261101193023-edcd391fa47d; master resolves to " + latest + "\n"},
		{"old-vendor", current, "# " + hgmxPath + " v0.0.33\n", "vendor/modules.txt [hgmx-latest] vendor contains hgmx v0.0.33; master resolves to " + latest + "\n"},
		{"replacement-footer", current, "# " + hgmxPath + " " + current + "\n# " + hgmxPath + " => ../hgmx\n", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			mod := "module example.com/site\n\ngo 1.27.0\n\nrequire " + hgmxPath + " " + tt.version + "\n"
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
				t.Fatal(err)
			}
			if tt.vendor != "" {
				if err := os.Mkdir(filepath.Join(dir, "vendor"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "vendor/modules.txt"), []byte(tt.vendor), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			s, err := loadSite(dir)
			if err != nil {
				t.Fatal(err)
			}
			a := &audit{site: s}
			if err := a.hgmxLatest(refs); err != nil {
				t.Fatal(err)
			}
			var got strings.Builder
			for _, f := range a.findings {
				got.WriteString(f.String() + "\n")
			}
			if got.String() != tt.want {
				t.Fatalf("findings:\n%swant:\n%s", got.String(), tt.want)
			}
		})
	}
}

func TestHgmxTool(t *testing.T) {
	for _, tool := range []string{"", "tool " + hgmxTool + "\n", "tool (\n\t" + hgmxTool + "\n)\n"} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/site\n\ngo 1.27.0\n\n"+tool), 0o644); err != nil {
			t.Fatal(err)
		}
		s, err := loadSite(dir)
		if err != nil {
			t.Fatal(err)
		}
		a := &audit{site: s}
		if err := a.hgmx(); err != nil {
			t.Fatal(err)
		}
		want := 0
		if tool == "" {
			want = 1
		}
		if len(a.findings) != want || want == 1 && a.findings[0].rule != ruleHgmxTool {
			t.Fatalf("tool %q: findings %v", tool, a.findings)
		}
	}
}

func TestHgmxOffline(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := latestHgmx(context.Background()); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("missing git: got %v; want executable not found", err)
	}
	dir := t.TempDir()
	mod := "module example.com/site\n\ngo 1.27.0\n\nrequire " + hgmxPath + " v0.0.33\n\ntool " + hgmxTool + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := loadSite(dir)
	if err != nil {
		t.Fatal(err)
	}
	a := &audit{site: s}
	if err := a.hgmx(); err != nil {
		t.Fatal(err)
	}
	if len(a.findings) != 0 {
		t.Fatalf("offline findings: %v", a.findings)
	}
}

func TestHgmxReplace(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	for _, tt := range []struct {
		name    string
		replace string
		paths   []string
	}{
		{"module", "replace " + hgmxPath + " => ../hgmx\n", []string{hgmxPath}},
		{"versioned-block", "replace (\n\t" + hgmxPath + " v0.0.33 => " + hgmxPath + " v0.0.32\n)\n", []string{hgmxPath}},
		{"tool", "replace " + hgmxTool + " => ../hgmx/cmd/hgmx\n", []string{hgmxTool}},
		{"multiple", "replace (\n\t" + hgmxPath + " => ../hgmx\n\t" + hgmxTool + " => ../hgmx/cmd/hgmx\n)\n", []string{hgmxPath, hgmxTool}},
		{"unrelated", "replace example.com/dependency => ../dependency\n", nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			mod := "module example.com/site\n\ngo 1.27.0\n\ntool " + hgmxTool + "\n\nrequire " + hgmxPath + " v0.0.36-0.20261101193023-edcd391fa47d\n\n" + tt.replace
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
				t.Fatal(err)
			}
			s, err := loadSite(dir)
			if err != nil {
				t.Fatal(err)
			}
			a := &audit{site: s}
			if err := a.hgmx(); err != nil {
				t.Fatal(err)
			}
			if len(a.findings) != len(tt.paths) {
				t.Fatalf("findings: %v; want %d", a.findings, len(tt.paths))
			}
			for i, path := range tt.paths {
				f := a.findings[i]
				if f.rule != ruleHgmxLatest || f.file != "go.mod" || f.line == 0 {
					t.Fatalf("unexpected finding: %+v", f)
				}
				if !strings.Contains(strings.Split(mod, "\n")[f.line-1], path) {
					t.Fatalf("finding points to the wrong line: %+v", f)
				}
				if f.msg != "replace "+path+" overrides the pinned hgmx version" || f.fix != "remove this replace; use go.work to develop against ../hgmx" {
					t.Fatalf("unexpected finding: %+v", f)
				}
			}
		})
	}
}
