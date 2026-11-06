package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestCheckNoVendor(t *testing.T) {
	t.Setenv("GOFLAGS", "-mod=mod")
	t.Setenv("GOWORK", filepath.Join(t.TempDir(), "missing.work"))
	for _, missing := range []bool{false, true} {
		name := "current-module"
		if missing {
			name = "missing-requirement"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			mod := "module " + hgmxPath + "\n\ngo 1.27.0\n"
			src := "package main\n\nfunc main() {}\n"
			if missing {
				dep := filepath.Join(dir, ".dependency")
				if err := os.Mkdir(dep, 0o755); err != nil {
					t.Fatal(err)
				}
				for name, data := range map[string]string{
					"go.mod": "module example.com/dependency\n\ngo 1.27.0\n",
					"dep.go": "package dependency\n",
				} {
					if err := os.WriteFile(filepath.Join(dep, name), []byte(data), 0o644); err != nil {
						t.Fatal(err)
					}
				}
				mod += "\nreplace example.com/dependency => ./.dependency\n"
				src = "package main\n\nimport _ \"example.com/dependency\"\n\nfunc main() {}\n"
			}
			for name, data := range map[string]string{"go.mod": mod, "go.sum": "", "main.go": src} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(dir)
			s, err := loadSite(".")
			if err != nil {
				t.Fatal(err)
			}
			for _, sub := range []string{"vet", "build", "test"} {
				cmd := s.gocmd(sub, "./...")
				if !slices.Equal(cmd.Args[1:], []string{sub, "-mod=readonly", "./..."}) {
					t.Fatalf("command: %v", cmd.Args)
				}
				if !slices.Contains(cmd.Env, "GOWORK=off") {
					t.Fatal("command must disable workspaces")
				}
			}
			want := 0
			if missing {
				want = 1
			}
			if got := checkCmd(nil); got != want {
				t.Fatalf("check exit: %d; want %d", got, want)
			}
			for name, want := range map[string]string{"go.mod": mod, "go.sum": ""} {
				data, err := os.ReadFile(filepath.Join(dir, name))
				if err != nil {
					t.Fatal(err)
				}
				if string(data) != want {
					t.Fatalf("%s changed:\n%s", name, data)
				}
			}
			if _, err := os.Stat(filepath.Join(dir, "vendor")); !os.IsNotExist(err) {
				t.Fatalf("vendor directory: %v", err)
			}
			if missing {
				cmd := s.gocmd("list", "-deps", "./...")
				cmd.Stdout, cmd.Stderr = nil, nil
				out, err := cmd.CombinedOutput()
				if err == nil || !strings.Contains(string(out), "replaced but not required") {
					t.Fatalf("readonly module error: %v\n%s", err, out)
				}
			}
		})
	}
}
