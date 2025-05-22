package main

import (
	"embed"
	"os"
	"path/filepath"
)

const LIB_DIR = "library"

type location struct {
	fs          embed.FS
	source      string
	destination string
	file        string
}

func copyEmbedFile(l location) error {
	data, err := l.fs.ReadFile(l.source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(l.destination), 0o755); err != nil {
		return err
	}
	return os.WriteFile(l.destination, data, 0o644)
}

func addComponent(l location) error {
	l.source = filepath.Join(LIB_DIR, l.source, l.file+".templ")
	l.destination = filepath.Join(l.destination, l.file+".templ")
	if err := copyEmbedFile(l); err != nil {
		return err
	}
	return nil
}

func copyDir(l location) error {
	entries, err := l.fs.ReadDir(l.source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(l.destination, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(l.source, entry.Name())
		dstPath := filepath.Join(l.destination, entry.Name())
		if entry.IsDir() {
			if err := copyDir(location{fs: l.fs, source: srcPath, destination: dstPath}); err != nil {
				return err
			}
		} else {
			data, err := l.fs.ReadFile(srcPath)
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
