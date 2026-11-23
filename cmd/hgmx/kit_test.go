package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKitUpgrade(t *testing.T) {
	const (
		auth = managedDir + "/auth/v1/auth.proto"
		app  = managedDir + "/app/v1/app.proto"
		old  = managedDir + "/old/v1/old.proto"
	)
	dir := t.TempDir()
	k := &kit{site: site{dir: dir}}
	plan := func(t *testing.T, shipped map[string]string) (*os.Root, *inputs) {
		t.Helper()
		k.shipped = map[string][]byte{}
		for rel, data := range shipped {
			k.shipped[rel] = []byte(data)
		}
		root, err := os.OpenRoot(dir)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { root.Close() })
		in, err := k.inputs(root)
		if err != nil {
			t.Fatal(err)
		}
		return root, in
	}
	sync := func(t *testing.T, shipped map[string]string) error {
		t.Helper()
		root, in := plan(t, shipped)
		_, _, err := in.sync(root)
		return err
	}
	read := func(t *testing.T, rel string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	v1 := map[string]string{bufFile: "version: v2\n", auth: "auth v1\n", old: "old v1\n"}
	if err := sync(t, v1); err != nil {
		t.Fatalf("adopt: %v", err)
	}

	v2 := map[string]string{bufFile: "version: v2\n", auth: "auth v2\n", app: "app v2\n"}
	if err := sync(t, v2); err != nil {
		t.Fatalf("unedited upgrade refused: %v", err)
	}
	for rel, want := range v2 {
		if got := read(t, rel); got != want {
			t.Errorf("%s = %q after upgrade, want %q", rel, got, want)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.Dir(old))); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("%s survived the upgrade: %v", filepath.Dir(old), err)
	}
	if _, in := plan(t, v2); len(in.stale()) > 0 || len(in.refused) > 0 {
		t.Errorf("upgrade left stale=%q refused=%q", in.stale(), in.refused)
	}

	if err := os.WriteFile(filepath.Join(dir, auth), []byte("auth edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v3 := map[string]string{bufFile: "version: v2\n", auth: "auth v3\n", app: "app v3\n"}
	err := sync(t, v3)
	if err == nil || !strings.Contains(err.Error(), auth+": differs from hgmx's copy") {
		t.Fatalf("edited copy: err = %v, want a refusal naming %s", err, auth)
	}
	if got := read(t, auth); got != "auth edited\n" {
		t.Errorf("refused sync rewrote %s to %q", auth, got)
	}
	if got := read(t, app); got != "app v2\n" {
		t.Errorf("refused sync rewrote %s to %q", app, got)
	}
}
