package main

import (
	"bytes"
	"io/fs"
	"os/exec"
	"testing"

	"github.com/cogikyo/hgmx"
)

func TestLibraryCredentialsCurrent(t *testing.T) {
	want, err := exec.Command("go", "run", hgmxPath+"/library", "credentials", "manifest").Output()
	if err != nil {
		t.Fatalf("library credentials manifest: %v", err)
	}
	have, err := fs.ReadFile(hgmx.LibraryFS, "library/"+credentialsFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(have, want) {
		t.Errorf("library/%s is stale; regenerate it with go run ./library credentials manifest:\n%s", credentialsFile, want)
	}
}
