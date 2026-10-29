package main

import (
	"bytes"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cogikyo/hgmx"
	"github.com/cogikyo/hgmx/palette"
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

	const viewsDir = "views"
	if _, err := os.Lstat(viewsDir); err == nil {
		log.Error("Refusing to overwrite existing " + viewsDir + "; move it aside or run init in an empty project")
		return 1
	} else if !errors.Is(err, fs.ErrNotExist) {
		log.Error("Failed to check "+viewsDir, slog.String("error", err.Error()))
		return 1
	}

	tmp, err := os.MkdirTemp(".", "."+viewsDir+"-*")
	if err != nil {
		log.Error("Failed to create temporary directory", slog.String("error", err.Error()))
		return 1
	}
	defer os.RemoveAll(tmp)

	if err := scaffold(hgmx.LibraryFS, "library", []string{"static", "views.templ", "views.go", "icon.png"}, tmp); err != nil {
		log.Error("Failed to copy scaffold", slog.String("error", err.Error()))
		return 1
	}
	if err := os.Chmod(tmp, 0o755); err != nil {
		log.Error("Failed to set permissions", slog.String("dir", tmp), slog.String("error", err.Error()))
		return 1
	}
	if err := os.Rename(tmp, viewsDir); err != nil {
		log.Error("Failed to move scaffold into place", slog.String("error", err.Error()))
		return 1
	}

	log.Info("hgmx project initialized successfully in ./" + viewsDir)
	return 0
}

// --- palette command ---

func paletteCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)

	themes := make([]theme, 0, len(args))
	for _, arg := range args {
		t, err := parseTheme(arg)
		if err != nil {
			log.Error("Invalid theme", slog.String("theme", arg), slog.String("error", err.Error()))
			return 64
		}
		themes = append(themes, t)
		for _, c := range t.palette.Clashes() {
			log.Warn(c.String(), slog.String("theme", t.name))
		}
	}

	var buf bytes.Buffer
	switch paletteFormat {
	case "tailwind", "css":
		for i, t := range themes {
			sel := palette.Selector(t.name)
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
	name    string
	palette palette.Palette
}

func parseTheme(arg string) (theme, error) {
	name, spec, ok := strings.Cut(arg, "=")
	if !ok || !palette.ValidName(name) {
		return theme{}, errors.New("expected <name>=<base>,<primary>[,<accent>] with a name of a-z, 0-9 and -")
	}
	hexes := strings.Split(spec, ",")
	if len(hexes) < 2 || len(hexes) > 3 {
		return theme{}, errors.New("expected <name>=<base>,<primary>[,<accent>]")
	}
	hexes = append(hexes, "")
	seed, err := palette.ParseSeed(hexes[0], hexes[1], hexes[2])
	if err != nil {
		return theme{}, err
	}
	return theme{name, palette.Generate(seed)}, nil
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
