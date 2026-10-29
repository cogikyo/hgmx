package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func copyFS(fsys fs.FS, src, dst string) error {
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
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		return os.WriteFile(out, data, 0o644)
	})
}

func scaffold(fsys fs.FS, root string, names []string, dst string) error {
	for _, name := range names {
		if err := copyFS(fsys, path.Join(root, name), filepath.Join(dst, name)); err != nil {
			return err
		}
	}
	return nil
}
