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

	log := newLogger(*logLevelFlag, stderr)

	log.Info("Environment:",
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

type templateLocation struct {
	Dir  string
	Name string
}

var components = map[string]templateLocation{
	"button": {"library/components/action", "button"},
	"avatar": {"library/components/display", "avatar"},
	"text":   {"library/components/display", "text"},
	"loader": {"library/components/feedback", "loader"},
	"input":  {"library/components/input", "input"},
}

var blocks = map[string]templateLocation{
	"base":     {"library/blocks/layouts", "base"},
	"settings": {"library/blocks/forms", "settings"},
	"hero":     {"library/blocks/content", "hero"},
	"navbar":   {"library/blocks/navigation", "navbar"},
	"alert":    {"library/blocks/partials", "alert"},
}

var pages = map[string]templateLocation{
	"home": {"library/pages/home", "home"},
}

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

	log := newLogger(*logLevelFlag, stderr)

	if *baseFlag {
		log.Info("Base/essentials-only mode is not yet implemented")
		return 1
	}

	appDir := "app"
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		log.Error("Failed to create app directory", slog.String("error", err.Error()))
		return 1
	}

	cssDir := filepath.Join(appDir, "static", "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		log.Error("Failed to create css directory", slog.String("error", err.Error()))
		return 1
	}

	if err := copyDirDirect(hgmx.LibraryFS, "library/static", filepath.Join(appDir, "static")); err != nil {
		log.Error("Failed to copy static directory", slog.String("error", err.Error()))
		return 1
	}

	for _, c := range components {
		dstDir := filepath.Join(appDir, "components", filepath.Base(c.Dir))
		if err := addComponent(hgmx.LibraryFS, c.Dir, dstDir, c.Name, "components", cssDir); err != nil {
			log.Error("Failed to copy component", slog.String("dir", c.Dir), slog.String("name", c.Name), slog.String("error", err.Error()))
			return 1
		}
	}

	for _, b := range blocks {
		dstDir := filepath.Join(appDir, "blocks", filepath.Base(b.Dir))
		if err := addComponent(hgmx.LibraryFS, b.Dir, dstDir, b.Name, "blocks", cssDir); err != nil {
			log.Error("Failed to copy block", slog.String("dir", b.Dir), slog.String("name", b.Name), slog.String("error", err.Error()))
			return 1
		}
	}

	for _, p := range pages {
		dstDir := filepath.Join(appDir, "pages", filepath.Base(p.Dir))
		if err := addComponent(hgmx.LibraryFS, p.Dir, dstDir, p.Name, "pages", cssDir); err != nil {
			log.Error("Failed to copy page", slog.String("dir", p.Dir), slog.String("name", p.Name), slog.String("error", err.Error()))
			return 1
		}
	}

	log.Info("hgmx project initialized successfully in ./app")
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

	log := newLogger(*logLevelFlag, stderr)

	remainingArgs := cmd.Args()
	if len(remainingArgs) != 1 {
		log.Error("Missing or too many arguments: expected exactly one (hex) color argument.")
		fmt.Fprint(stderr, paletteUsageText)
		return 64
	}

	hexColor := remainingArgs[0]

	// TODO: more than hex
	if len(hexColor) != 7 || hexColor[0] != '#' {
		log.Error("Invalid hex color format. Expected #RRGGBB.", slog.String("color", hexColor))
		return 1
	}

	log.Info("Generating palette for color:", slog.String("hex", hexColor))

	generatedPalette := palette.Generate(hexColor)

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

// --- link command ---

const linkUsageText = `usage: hgmx link <dir>

Symlinks files from the given directory into the local app directory, replacing any duplicate files with symlinks.

Args:
  dir   Path to the directory to link from (e.g., '.')
  -l    Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func linkCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("link", flag.ExitOnError)
	logLevelFlag := cmd.String("l", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, linkUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	log := newLogger(*logLevelFlag, stderr)
	remainingArgs := cmd.Args()
	if len(remainingArgs) != 1 {
		log.Error("Missing or too many arguments: expected exactly one directory argument.")
		fmt.Fprint(stderr, linkUsageText)
		return 64
	}

	srcDir := remainingArgs[0]
	appDir := "app"
	if _, err := os.Stat(appDir); err != nil {
		log.Error("Could not find 'app' directory in current project.", slog.String("error", err.Error()))
		return 1
	}

	err = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		appPath := filepath.Join(appDir, rel)
		if _, err := os.Stat(appPath); err == nil {
			// Remove the file in appDir
			if err := os.Remove(appPath); err != nil {
				log.Error("Failed to remove file in appDir", slog.String("file", appPath), slog.String("error", err.Error()))
				return err
			}
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(appPath), 0o755); err != nil {
				log.Error("Failed to create parent directory in appDir", slog.String("dir", filepath.Dir(appPath)), slog.String("error", err.Error()))
				return err
			}
			// Create symlink
			absPath, absErr := filepath.Abs(path)
			if absErr != nil {
				log.Error("Failed to get absolute path", slog.String("src", path), slog.String("error", absErr.Error()))
				return absErr
			}
			if err := os.Symlink(absPath, appPath); err != nil {
				log.Error("Failed to create symlink", slog.String("src", absPath), slog.String("dst", appPath), slog.String("error", err.Error()))
				return err
			}
			log.Info("Symlinked", slog.String("src", path), slog.String("dst", appPath))
		}
		return nil
	})
	if err != nil {
		log.Error("Error during linking", slog.String("error", err.Error()))
		return 1
	}

	log.Info("Linking complete.")
	return 0
}
