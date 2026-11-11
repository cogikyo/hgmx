package main

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	cssInput  = "views/static/css/main.css"
	cssOutput = "views/static/css/main.min.css"
	templPath = "github.com/a-h/templ"
	readonly  = "-mod=readonly"
	serverPkg = "./cmd/server"
)

type site struct {
	dir string
	mod *modfile.File
}

func loadSite(dir string) (site, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return site{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return site{}, err
	}
	mod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return site{}, err
	}
	if mod.Module == nil {
		return site{}, errors.New("go.mod declares no module path")
	}
	return site{dir: dir, mod: mod}, nil
}

func (s site) path(rel string) string {
	return filepath.Join(s.dir, filepath.FromSlash(rel))
}

func (s site) templ() *modfile.Require {
	for _, r := range s.mod.Require {
		if r.Mod.Path == templPath {
			return r
		}
	}
	return nil
}

func (s site) css() bool {
	_, err := os.Stat(s.path(cssInput))
	return err == nil
}

func (s site) env() []string {
	return append(os.Environ(), "GOWORK=off")
}

func (s site) gocmd(sub string, args ...string) *exec.Cmd {
	cmd := exec.Command("go", slices.Concat([]string{sub, readonly}, args)...)
	cmd.Dir = s.dir
	cmd.Env = s.env()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (s site) walk(fn func(rel string, d fs.DirEntry) error) error {
	return filepath.WalkDir(s.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(s.dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() && rel != "." && skipDir(d.Name()) {
			return filepath.SkipDir
		}
		return fn(rel, d)
	})
}

func skipDir(name string) bool {
	switch name {
	case "vendor", "node_modules", "testdata":
		return true
	}
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}
