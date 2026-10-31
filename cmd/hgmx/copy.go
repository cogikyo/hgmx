package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func copyFS(fsys fs.FS, src, dst string, edit func(name string, data []byte) ([]byte, bool)) error {
	return fs.WalkDir(fsys, src, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.FromSlash(src), filepath.FromSlash(name))
		if err != nil {
			return err
		}
		out := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(out, 0o755)
		}
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err
		}
		data, keep := edit(name, data)
		if !keep {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		return os.WriteFile(out, data, 0o644)
	})
}
