package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHgmxLatest(t *testing.T) {
	const latest = "v0.0.36-0.20261101193023-edcd391fa47d"
	for _, tt := range []struct {
		name    string
		version string
		vendor  string
		want    string
	}{
		{"current", latest, "# " + hgmxPath + " " + latest + "\n## explicit; go 1.27.0\n" + hgmxTool + "\n", ""},
		{"no-vendor", latest, "", ""},
		{"old-module", "v0.0.33", "", "go.mod [hgmx-latest] site requires hgmx v0.0.33; master resolves to " + latest + "\n"},
		{"old-vendor", latest, "# " + hgmxPath + " v0.0.33\n", "vendor/modules.txt [hgmx-latest] vendor contains hgmx v0.0.33; master resolves to " + latest + "\n"},
		{"replacement-footer", latest, "# " + hgmxPath + " " + latest + "\n# " + hgmxPath + " => ../hgmx\n", ""},
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
			if err := a.hgmxLatest(latest); err != nil {
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
	t.Setenv("GOPROXY", "off")
	t.Setenv("GONOPROXY", "none")
	t.Setenv("GOMODCACHE", t.TempDir())
	t.Setenv("GOTOOLCHAIN", "local")
	if _, err := latestHgmx(context.Background()); err == nil {
		t.Fatal("offline resolution succeeded")
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
	t.Setenv("GOPROXY", "off")
	t.Setenv("GONOPROXY", "none")
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
