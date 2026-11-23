package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallKeepsUnrelatedStaging(t *testing.T) {
	dir := t.TempDir()
	backup := filepath.Join(dir, "apps", ".android.hgmx-backup", "OnlyCopy.txt")
	if err := os.MkdirAll(filepath.Dir(backup), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backup, []byte("only copy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files := []starterFile{{path: "app/build.gradle.kts", data: []byte("plugins {}\n"), mode: 0o644}}
	if err := install(site{dir: dir}, "apps/android", files, slog.New(slog.DiscardHandler)); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(backup); err != nil || string(data) != "only copy\n" {
		t.Errorf("adoption touched %s: %q, %v", backup, data, err)
	}
	for _, rel := range []string{"apps/android/app/build.gradle.kts", compatFile} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("adoption did not write %s: %v", rel, err)
		}
	}
	entries, err := os.ReadDir(filepath.Join(dir, "apps"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "android" && e.Name() != ".android.hgmx-backup" && e.Name() != filepath.Base(compatFile) {
			t.Errorf("adoption left %s behind", e.Name())
		}
	}
}
