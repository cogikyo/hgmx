package main

import (
	"bytes"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/cogikyo/hgmx"
	"github.com/cogikyo/hgmx/palette"
	"golang.org/x/mod/modfile"
)

// --- info command ---

func infoCmd() (code int) {
	log := newLogger(logLevel, os.Stderr)

	log.Info("Environment:",
		slog.Group("versions",
			slog.String("hgmx", hgmx.Version()),
			slog.String("go", runtime.Version()),
		),
	)

	return 0
}

// --- init command ---

func initCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)

	const viewsDir, mainFile = "views", "main.go"
	for _, name := range []string{viewsDir, mainFile} {
		if _, err := os.Lstat(name); err == nil {
			log.Error("Refusing to overwrite existing " + name + "; move it aside or run init in an empty project")
			return 1
		} else if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Failed to check "+name, slog.String("error", err.Error()))
			return 1
		}
	}

	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		log.Error("Failed to read go.mod; run go mod init first", slog.String("error", err.Error()))
		return 1
	}
	module := modfile.ModulePath(gomod)
	if module == "" {
		log.Error("go.mod declares no module path")
		return 1
	}

	tmp, err := os.MkdirTemp(".", "."+viewsDir+"-*")
	if err != nil {
		log.Error("Failed to create temporary directory", slog.String("error", err.Error()))
		return 1
	}
	defer os.RemoveAll(tmp)

	if err := scaffold(hgmx.LibraryFS, "library", []string{"static", "views.templ", "views_templ.go", "views.go"}, tmp); err != nil {
		log.Error("Failed to copy scaffold", slog.String("error", err.Error()))
		return 1
	}
	src, err := hgmx.LibraryFS.ReadFile("library/main/main.go")
	if err != nil {
		log.Error("Failed to read scaffold main.go", slog.String("error", err.Error()))
		return 1
	}
	src = bytes.Replace(src, []byte(`"github.com/cogikyo/hgmx/library"`), []byte(strconv.Quote(module+"/"+viewsDir)), 1)
	staged := filepath.Join(tmp, mainFile)
	if err := os.WriteFile(staged, src, 0o644); err != nil {
		log.Error("Failed to write "+mainFile, slog.String("error", err.Error()))
		return 1
	}
	if err := os.Chmod(tmp, 0o755); err != nil {
		log.Error("Failed to set permissions", slog.String("dir", tmp), slog.String("error", err.Error()))
		return 1
	}

	if err := os.Rename(staged, mainFile); err != nil {
		log.Error("Failed to move "+mainFile+" into place", slog.String("error", err.Error()))
		return 1
	}
	if err := os.Rename(tmp, viewsDir); err != nil {
		os.Remove(mainFile)
		log.Error("Failed to move scaffold into place", slog.String("error", err.Error()))
		return 1
	}

	log.Info("hgmx project initialized successfully in ./" + viewsDir + " and ./" + mainFile)
	return 0
}

// --- palette command ---

func paletteCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)

	themes := make([]theme, 0, len(args))
	seen := map[string]bool{}
	for _, arg := range args {
		t, err := parseTheme(arg)
		if err == nil && seen[t.name] {
			err = errors.New("theme name used twice")
		}
		if err != nil {
			log.Error("Invalid theme", slog.String("theme", arg), slog.String("error", err.Error()))
			return 64
		}
		seen[t.name] = true
		themes = append(themes, t)
		for _, c := range t.palette.Clashes() {
			log.Warn(c.String(), slog.String("theme", t.name))
		}
	}

	var buf bytes.Buffer
	switch paletteFormat {
	case "tailwind", "css":
		for i, t := range themes {
			sel := t.selector
			if i == 0 {
				sel = ":root, " + sel
			} else {
				buf.WriteByte('\n')
			}
			t.palette.CSS(&buf, sel)
		}
		if paletteFormat == "tailwind" {
			buf.WriteByte('\n')
			themes[0].palette.Theme(&buf)
		}
	case "lua":
		if len(themes) != 1 {
			log.Error("Lua output takes exactly one theme")
			return 64
		}
		themes[0].palette.Lua(&buf)
	default:
		log.Error("Unknown format", slog.String("format", paletteFormat))
		return 64
	}

	if paletteOut == "" {
		if _, err := os.Stdout.Write(buf.Bytes()); err != nil {
			log.Error("Failed to write palette", slog.String("error", err.Error()))
			return 1
		}
		return 0
	}
	if err := writeAtomic(paletteOut, buf.Bytes()); err != nil {
		log.Error("Failed to write palette", slog.String("file", paletteOut), slog.String("error", err.Error()))
		return 1
	}
	log.Info("Palette written", slog.String("file", paletteOut), slog.Int("themes", len(themes)))
	return 0
}

type theme struct {
	name     string
	selector string
	palette  palette.Palette
}

func parseTheme(arg string) (theme, error) {
	name, spec, ok := strings.Cut(arg, "=")
	if !ok {
		return theme{}, errors.New("expected <name>=<base>,<primary>[,<accent>]")
	}
	sel, err := palette.Selector(name)
	if err != nil {
		return theme{}, err
	}
	hexes := strings.Split(spec, ",")
	if len(hexes) < 2 || len(hexes) > 3 || slices.Contains(hexes, "") {
		return theme{}, errors.New("expected <name>=<base>,<primary>[,<accent>]")
	}
	hexes = append(hexes, "")
	seed, err := palette.ParseSeed(hexes[0], hexes[1], hexes[2])
	if err != nil {
		return theme{}, err
	}
	return theme{name, sel, palette.Generate(seed)}, nil
}

func writeAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
