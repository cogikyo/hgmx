package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/a-h/templ/cmd/templ/generatecmd"
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

	const root, prefix = "library", "github.com/cogikyo/hgmx/library"
	entries, err := fs.ReadDir(hgmx.LibraryFS, root)
	if err != nil {
		log.Error("Failed to read scaffold", slog.String("error", err.Error()))
		return 1
	}
	for _, e := range entries {
		if _, err := os.Lstat(e.Name()); err == nil {
			log.Error("Refusing to overwrite existing " + e.Name() + "; move it aside or run init in an empty project")
			return 1
		} else if !errors.Is(err, fs.ErrNotExist) {
			log.Error("Failed to check "+e.Name(), slog.String("error", err.Error()))
			return 1
		}
	}

	data, err := os.ReadFile("go.mod")
	if err != nil {
		log.Error("Failed to read go.mod; run go mod init first", slog.String("error", err.Error()))
		return 1
	}
	gomod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		log.Error("Failed to parse go.mod", slog.String("error", err.Error()))
		return 1
	}
	if gomod.Module == nil {
		log.Error("go.mod declares no module path")
		return 1
	}
	module := gomod.Module.Mod.Path
	if err := gomod.AddRequire(templPath, templ.Version()); err != nil {
		log.Error("Failed to require templ", slog.String("error", err.Error()))
		return 1
	}
	gomod.Cleanup()
	gomod.SortBlocks()
	updated, err := gomod.Format()
	if err != nil {
		log.Error("Failed to format go.mod", slog.String("error", err.Error()))
		return 1
	}

	tmp, err := os.MkdirTemp(".", "hgmx-init-*")
	if err != nil {
		log.Error("Failed to create temporary directory", slog.String("error", err.Error()))
		return 1
	}
	defer os.RemoveAll(tmp)

	rewrite := func(name string, data []byte) ([]byte, bool) {
		switch {
		case strings.HasSuffix(name, "_templ.go"):
			return nil, false
		case strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".templ"):
			return bytes.ReplaceAll(data, []byte(prefix), []byte(module)), true
		}
		return data, true
	}
	if err := copyFS(hgmx.LibraryFS, root, tmp, rewrite); err != nil {
		log.Error("Failed to copy scaffold", slog.String("error", err.Error()))
		return 1
	}
	if err := generatecmd.Run(context.Background(), io.Discard, os.Stderr, []string{"-path", tmp, "-log-level", "error"}); err != nil {
		log.Error("Failed to generate templ code", slog.String("error", err.Error()))
		return 1
	}

	var moved []string
	for _, e := range entries {
		if err := os.Rename(filepath.Join(tmp, e.Name()), e.Name()); err != nil {
			for _, name := range moved {
				os.RemoveAll(name)
			}
			log.Error("Failed to move "+e.Name()+" into place", slog.String("error", err.Error()))
			return 1
		}
		moved = append(moved, e.Name())
	}
	if err := writeAtomic("go.mod", updated); err != nil {
		log.Error("Failed to write go.mod", slog.String("error", err.Error()))
		return 1
	}

	log.Info("hgmx project initialized", slog.String("files", strings.Join(moved, " ")), slog.String("templ", templ.Version()))
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
