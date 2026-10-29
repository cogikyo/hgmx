package main

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"runtime"

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

	if len(args) != 1 {
		log.Error("Missing or too many arguments: expected exactly one (hex) color argument.")
		return 64
	}

	hexColor := args[0]

	log.Info("Generating palette for color:", slog.String("hex", hexColor))

	generatedPalette, err := palette.Generate(hexColor)
	if err != nil {
		log.Error("Invalid color", slog.String("error", err.Error()))
		return 1
	}

	if paletteLua {
		generatedPalette.ToLua(os.Stdout)
		return 0
	}

	outputFile := "library/static/css/colors.css"
	f, err := os.Create(outputFile)
	if err != nil {
		log.Error("Failed to open output file for writing", slog.String("file", outputFile), slog.String("error", err.Error()))
		return 1
	}
	defer f.Close()
	generatedPalette.ToCSS(f)

	log.Info("Palette successfully generated and written", slog.String("file", outputFile))
	return 0
}
