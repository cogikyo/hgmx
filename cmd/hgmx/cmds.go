package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nosvagor/hgmx"
	"github.com/nosvagor/hgmx/internal/palette"
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

var components = map[string]string{
	"button": "action",
	"avatar": "display",
	"text":   "display",
	"loader": "feedback",
	"input":  "input",
}

var blocks = map[string]string{
	"settings": "forms",
	"hero":     "content",
	"navbar":   "navigation",
	"alert":    "partials",
}

var pages = map[string]string{
	"home":     "home",
	"settings": "settings",
	"notfound": "notfound",
}

func initCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)

	viewsDir := "views"
	if err := os.MkdirAll(viewsDir, 0o755); err != nil {
		log.Error("Failed to create views directory", slog.String("error", err.Error()))
		return 1
	}

	if err := copyDir(location{fs: hgmx.LibraryFS, source: LIB_DIR + "/static", destination: filepath.Join(viewsDir, "static")}); err != nil {
		log.Error("Failed to copy static directory", slog.String("error", err.Error()))
		return 1
	}

	for file, dir := range components {
		dir := filepath.Join("components", dir)
		if err := addComponent(location{fs: hgmx.LibraryFS, source: dir, destination: filepath.Join(viewsDir, dir), file: file}); err != nil {
			log.Error("Failed to copy component", slog.String("src", dir), slog.String("file", file), slog.String("error", err.Error()))
			return 1
		}
	}

	for file, dir := range blocks {
		dir := filepath.Join("blocks", dir)
		if err := addComponent(location{fs: hgmx.LibraryFS, source: dir, destination: filepath.Join(viewsDir, dir), file: file}); err != nil {
			log.Error("Failed to copy block", slog.String("src", dir), slog.String("file", file), slog.String("error", err.Error()))
			return 1
		}
	}

	for file, dir := range pages {
		dir := filepath.Join("pages", dir)
		if err := addComponent(location{fs: hgmx.LibraryFS, source: dir, destination: filepath.Join(viewsDir, dir), file: file}); err != nil {
			log.Error("Failed to copy page", slog.String("src", dir), slog.String("file", file), slog.String("error", err.Error()))
			return 1
		}
	}

	if err := copyEmbedFile(location{fs: hgmx.LibraryFS, source: LIB_DIR + "/views.templ", destination: filepath.Join(viewsDir, "views.templ")}); err != nil {
		log.Error("Failed to copy views.templ", slog.String("error", err.Error()))
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

func linkCmd(inputGlob, outputGlob string) (code int) {
	l := newLogger(logLevel, os.Stderr)

	inputDirs, err := filepath.Glob(inputGlob)
	if err != nil || len(inputDirs) == 0 {
		l.Error("No input directories found", slog.String("pattern", inputGlob))
		return 1
	}
	outputDirs, err := filepath.Glob(outputGlob)
	if err != nil || len(outputDirs) == 0 {
		l.Error("No output directories found", slog.String("pattern", outputGlob))
		return 1
	}

	l.Warn("Destination (output): ", slog.Any("searching in", inputDirs))
	l.Warn("Source (input): ", slog.Any("linking to", outputDirs))

	for i, srcDir := range inputDirs {
		var dstDir string
		if i < len(outputDirs) {
			dstDir = outputDirs[i]
		} else {
			dstDir = outputDirs[len(outputDirs)-1]
		}

		if _, err := os.Stat(srcDir); err != nil {
			l.Error("Input directory does not exist", slog.String("dir", srcDir), slog.String("error", err.Error()))
			continue
		}
		if _, err := os.Stat(dstDir); err != nil {
			if err := os.MkdirAll(dstDir, 0o755); err != nil {
				l.Error("Failed to create output directory", slog.String("dir", dstDir), slog.String("error", err.Error()))
				continue
			}
		}

		err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, "_templ.go") {
				l.Debug("Skipping generated file", slog.String("file", path))
				return nil
			}
			rel, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}
			dstPath := filepath.Join(dstDir, rel)
			if _, err := os.Stat(dstPath); err == nil {
				if err := os.Remove(dstPath); err != nil {
					l.Error("Failed to remove file in output dir", slog.String("file", dstPath), slog.String("error", err.Error()))
					return err
				}
			}
			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				l.Error("Failed to create parent directory in output dir", slog.String("dir", filepath.Dir(dstPath)), slog.String("error", err.Error()))
				return err
			}
			absPath, absErr := filepath.Abs(path)
			if absErr != nil {
				l.Error("Failed to get absolute path", slog.String("src", path), slog.String("error", absErr.Error()))
				return absErr
			}
			if err := os.Symlink(absPath, dstPath); err != nil {
				l.Error("Failed to create symlink", slog.String("src", absPath), slog.String("dst", dstPath), slog.String("error", err.Error()))
				return err
			}
			l.Info("Symlinked", slog.String("src", path), slog.String("dst", dstPath))
			return nil
		})
		if err != nil {
			l.Error("Error during linking", slog.String("input", srcDir), slog.String("output", dstDir), slog.String("error", err.Error()))
			continue
		}
	}

	l.Info("Linking complete.")
	return 0
}
