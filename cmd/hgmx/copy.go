package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

func cssImport(targetCSS, importPath string) error {
	data, err := os.ReadFile(targetCSS)
	if err == nil && bytes.Contains(data, []byte(importPath)) {
		return nil
	}
	f, err := os.OpenFile(targetCSS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "@import '%s';\n", importPath)
	return err
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

func addComponent(fs embed.FS, srcDir, dstDir, name, cssGroup, cssDir string) error {
	for _, ext := range []string{".templ", ".css"} {
		src := filepath.Join(srcDir, name+ext)
		dst := filepath.Join(dstDir, name+ext)
		if err := copyEmbedFile(fs, src, dst); err != nil {
			if ext == ".css" {
				continue
			}
			return err
		}
		if ext == ".css" {
			var importPath string
			switch cssGroup {
			case "blocks":
				importPath = "../../blocks/" + filepath.Base(srcDir) + "/" + name + ".css"
			case "components":
				importPath = "../../components/" + filepath.Base(srcDir) + "/" + name + ".css"
			case "pages":
				importPath = "../../pages/" + filepath.Base(srcDir) + "/" + name + ".css"
			}
			targetCSS := filepath.Join(cssDir, cssGroup+".css")
			if err := cssImport(targetCSS, importPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyDirDirect(efs embed.FS, src, dst string) error {
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
			if err := copyDirDirect(efs, srcPath, dstPath); err != nil {
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
