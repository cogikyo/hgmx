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
  -v           Set log verbosity level to "debug". (default "info")
  -log-level   Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func infoCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("info", flag.ExitOnError)
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, infoUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	l := l.NewLogger(*logLevelFlag, *verboseFlag, stderr)

	l.Info("Environment:",
		slog.Group("versions",
			slog.String("hgmx", hgmx.Version()),
			slog.String("go", runtime.Version()),
		),
	)

	return 0
}

const initUsageText = `usage: hgmx init

Initializes a new hgmx project.

Args:
  -v           Set log verbosity level to "debug". (default "info")
  -log-level   Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func initCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("init", flag.ExitOnError)
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, initUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	l := l.NewLogger(*logLevelFlag, *verboseFlag, stderr)

	// Copy static directory
	srcStatic := "static"
	dstStatic := "static"
	if err := copyEmbedDir(hgmx.LibraryFS, srcStatic, dstStatic); err != nil {
		l.Error("Failed to copy static directory", slog.String("error", err.Error()))
		return 1
	}
	l.Info("Copied static directory", slog.String("from", srcStatic), slog.String("to", dstStatic))

	// Ensure views directory exists
	if err := os.MkdirAll("views", 0o755); err != nil {
		l.Error("Failed to create views directory", slog.String("error", err.Error()))
		return 1
	}

	src := "components/index/index.templ"
	dst := "views/index.templ"
	data, err := hgmx.LibraryFS.ReadFile(src)
	if err != nil {
		l.Error("Failed to read embedded index.templ", slog.String("file", src), slog.String("error", err.Error()))
		return 1
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		l.Error("Failed to write index.templ", slog.String("file", dst), slog.String("error", err.Error()))
		return 1
	}
	l.Info("Copied index.templ to views directory")

	l.Info("hgmx project initialized successfully")
	return 0
}

// --- palette command ---

const paletteUsageText = `usage: hgmx palette <hex_color>

Generates a color palette based on the input hex color using OKLCH.

Args:
  color    The base background color in hex format (e.g., "#RRGGBB").
  -v           Set log verbosity level to "debug". (default "info")
  -log-level   Set log verbosity level. (default "info", options: "debug", "info", "warn", "error")
`

func paletteCmd(stdout, stderr io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("palette", flag.ExitOnError)
	verboseFlag := cmd.Bool("v", false, "")
	logLevelFlag := cmd.String("log-level", "info", "")

	cmd.Usage = func() {
		fmt.Fprint(stderr, paletteUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return 64
	}

	lg := l.NewLogger(*logLevelFlag, *verboseFlag, stderr)

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
