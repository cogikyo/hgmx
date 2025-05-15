package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nosvagor/hgmx"
	"github.com/nosvagor/hgmx/internal/palette"
	l "github.com/nosvagor/hgmx/internal/slog"
)

// --- info command ---

const infoUsageText = `usage: hgmx info [<args>]

Displays information about the hgmx environment.

Args:
  -l	Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func infoCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("info", flag.ExitOnError)
	logLevelFlag := cmd.String("l", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, infoUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	l := l.NewLogger(*logLevelFlag, stderr)

	l.Info("Environment:",
		slog.Group("versions",
			slog.String("hgmx", hgmx.Version()),
			slog.String("go", runtime.Version()),
		),
	)

	return 0
}

// --- init command ---

const initUsageText = `usage: hgmx init [options]

Initializes a new hgmx project in ./app.

Options:
  -b, --base   Only copy essential files (not yet implemented)
  -l, --log    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func initCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("init", flag.ExitOnError)
	baseFlag := cmd.Bool("b", false, "only copy essentials")
	cmd.BoolVar(baseFlag, "base", false, "only copy essentials")
	logLevelFlag := cmd.String("l", "info", "")
	cmd.StringVar(logLevelFlag, "log", "info", "set log verbosity level")

	cmd.Usage = func() {
		fmt.Fprint(stderr, initUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	l := l.NewLogger(*logLevelFlag, stderr)

	if *baseFlag {
		l.Info("Base/essentials-only mode is not yet implemented")
		return 1
	}

	appDir := "app"
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		l.Error("Failed to create app directory", slog.String("error", err.Error()))
		return 1
	}

	cssDir := filepath.Join(appDir, "static", "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		l.Error("Failed to create css directory", slog.String("error", err.Error()))
		return 1
	}

	copyDirs := []struct {
		src, dst string
	}{
		{"library/static", filepath.Join(appDir, "static")},
		{"library/components", filepath.Join(appDir, "components")},
	}
	for _, d := range copyDirs {
		if err := copyEmbedDir(hgmx.LibraryFS, d.src, d.dst); err != nil {
			l.Error("Failed to copy directory", slog.String("src", d.src), slog.String("error", err.Error()))
			return 1
		}
	}

	targets := []copyTarget{
		{"library/blocks/layouts", filepath.Join(appDir, "blocks", "layouts"), "base", "blocks"},
		{"library/blocks/layouts", filepath.Join(appDir, "blocks", "layouts"), "auth", "blocks"},
		{"library/blocks/forms", filepath.Join(appDir, "blocks", "forms"), "settings", "blocks"},
		{"library/blocks/content", filepath.Join(appDir, "blocks", "content"), "hero", "blocks"},
		{"library/blocks/navigation", filepath.Join(appDir, "blocks", "navigation"), "navbar", "blocks"},
		{"library/blocks/partials", filepath.Join(appDir, "blocks", "partials"), "alert", "blocks"},
		{"library/pages/home", filepath.Join(appDir, "pages", "home"), "home", "pages"},
	}
	for _, t := range targets {
		if err := copyTemplAndCSS(hgmx.LibraryFS, t.srcDir, t.dstDir, t.name, t.cssGroup, cssDir); err != nil {
			l.Error("Failed to copy templ/css", slog.String("srcDir", t.srcDir), slog.String("name", t.name), slog.String("error", err.Error()))
			return 1
		}
	}

	l.Info("hgmx project initialized successfully in ./app")
	return 0
}

// --- palette command ---

const paletteUsageText = `usage: hgmx palette <hex_color>

Generates a color palette based on the input hex color using OKLCH.

Args:
  color    The base background color in hex format (e.g., "#RRGGBB").
  -l       Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func paletteCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("palette", flag.ExitOnError)
	logLevelFlag := cmd.String("log-level", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, paletteUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	lg := l.NewLogger(*logLevelFlag, stderr)

	remainingArgs := cmd.Args()
	if len(remainingArgs) != 1 {
		lg.Error("Missing or too many arguments: expected exactly one (hex) color argument.")
		fmt.Fprint(stderr, paletteUsageText)
		return 64
	}

	hexColor := remainingArgs[0]

	// TODO: more than hex
	if len(hexColor) != 7 || hexColor[0] != '#' {
		lg.Error("Invalid hex color format. Expected #RRGGBB.", slog.String("color", hexColor))
		return 1
	}

	lg.Info("Generating palette for color:", slog.String("hex", hexColor))

	generatedPalette := palette.Generate(hexColor)

	outputFile := "static/css/colors.css"
	f, err := os.Create(outputFile)
	if err != nil {
		lg.Error("Failed to open output file for writing", slog.String("file", outputFile), slog.String("error", err.Error()))
		return 1
	}
	defer f.Close()
	generatedPalette.ToCSS(f)

	lg.Info("Palette successfully generated and written", slog.String("file", outputFile))
	return 0
}
