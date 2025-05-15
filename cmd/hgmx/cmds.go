package main

import (
	"embed"
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

	// Copy static directory
	if err := copyEmbedDir(hgmx.LibraryFS, "static", filepath.Join(appDir, "static")); err != nil {
		l.Error("Failed to copy static directory", slog.String("error", err.Error()))
		return 1
	}
	l.Info("Copied static directory", slog.String("to", filepath.Join(appDir, "static")))

	// Copy all components
	if err := copyEmbedDir(hgmx.LibraryFS, "components", filepath.Join(appDir, "components")); err != nil {
		l.Error("Failed to copy components directory", slog.String("error", err.Error()))
		return 1
	}
	l.Info("Copied components directory", slog.String("to", filepath.Join(appDir, "components")))

	// Copy pages/home/home.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "pages/home/home.templ", filepath.Join(appDir, "pages", "home", "home.templ")); err != nil {
		l.Error("Failed to copy home.templ", slog.String("error", err.Error()))
		return 1
	}

	// Copy blocks/layouts/base.templ and auth.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/layouts/base.templ", filepath.Join(appDir, "blocks", "layouts", "base.templ")); err != nil {
		l.Error("Failed to copy base.templ", slog.String("error", err.Error()))
		return 1
	}
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/layouts/auth.templ", filepath.Join(appDir, "blocks", "layouts", "auth.templ")); err != nil {
		l.Error("Failed to copy auth.templ", slog.String("error", err.Error()))
		return 1
	}

	// Copy blocks/forms/settings.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/forms/settings.templ", filepath.Join(appDir, "blocks", "forms", "settings.templ")); err != nil {
		l.Error("Failed to copy settings.templ", slog.String("error", err.Error()))
		return 1
	}

	// Copy blocks/content/hero.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/content/hero.templ", filepath.Join(appDir, "blocks", "content", "hero.templ")); err != nil {
		l.Error("Failed to copy hero.templ", slog.String("error", err.Error()))
		return 1
	}

	// Copy blocks/navigation/navbar.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/navigation/navbar.templ", filepath.Join(appDir, "blocks", "navigation", "navbar.templ")); err != nil {
		l.Error("Failed to copy navbar.templ", slog.String("error", err.Error()))
		return 1
	}

	// Copy blocks/partials/alert.templ
	if err := copyEmbedFile(hgmx.LibraryFS, "blocks/partials/alert.templ", filepath.Join(appDir, "blocks", "partials", "alert.templ")); err != nil {
		l.Error("Failed to copy alert.templ", slog.String("error", err.Error()))
		return 1
	}

	// Ensure CSS directories exist
	cssDir := filepath.Join(appDir, "static", "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		l.Error("Failed to create css directory", slog.String("error", err.Error()))
		return 1
	}

	// Write import statements for CSS files
	blocksCSS := filepath.Join(cssDir, "blocks.css")
	componentsCSS := filepath.Join(cssDir, "components.css")
	pagesCSS := filepath.Join(cssDir, "pages.css")

	writeCSSImports(blocksCSS, []string{
		"../../blocks/layouts/base.css",
		"../../blocks/layouts/auth.css",
		"../../blocks/forms/settings.css",
		"../../blocks/content/hero.css",
		"../../blocks/navigation/navbar.css",
		"../../blocks/partials/alert.css",
	})
	writeCSSImports(componentsCSS, collectComponentCSSImports("../../components"))
	writeCSSImports(pagesCSS, []string{
		"../../pages/home/home.css",
	})

	l.Info("hgmx project initialized successfully in ./app")
	return 0
}

func copyEmbedFile(efs embed.FS, src, dst string) error {
	data, err := efs.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func writeCSSImports(target string, imports []string) {
	f, err := os.Create(target)
	if err != nil {
		return
	}
	defer f.Close()
	for _, imp := range imports {
		fmt.Fprintf(f, "@import '%s';\n", imp)
	}
}

func collectComponentCSSImports(base string) []string {
	var imports []string
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".css" {
			imports = append(imports, path)
		}
		return nil
	})
	return imports
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

// Helper to recursively copy a directory from embed.FS to disk
func copyEmbedDir(efs embed.FS, src, dst string) error {
	entries, err := efs.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyEmbedDir(efs, srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := efs.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}
