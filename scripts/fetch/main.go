// Package main verifies or fetches scripts listed in lock.json.
package main

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type entry struct {
	// File is relative to vendor/.
	File string `json:"file"`
	// Version labels status output; it does not select the download.
	Version string `json:"version"`
	// URL supplies missing or mismatched files.
	URL string `json:"url"`
	// SHA256 is the required lowercase hex digest of the downloaded bytes.
	SHA256 string `json:"sha256"`
}

func main() {
	check := flag.Bool("check", false, "verify vendored files against lock.json without downloading")
	flag.Parse()
	dir := cmp.Or(flag.Arg(0), ".")

	data, err := os.ReadFile(filepath.Join(dir, "lock.json"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var lock []entry
	if err := json.Unmarshal(data, &lock); err != nil {
		fmt.Fprintf(os.Stderr, "lock.json: %v\n", err)
		os.Exit(1)
	}

	failed := false
	for _, e := range lock {
		status, err := sync(filepath.Join(dir, "vendor", e.File), e, *check)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL  %s %s: %v\n", e.File, e.Version, err)
			failed = true
			continue
		}
		fmt.Printf("%-5s %s %s\n", status, e.File, e.Version)
	}
	if failed {
		os.Exit(1)
	}
}

func sync(path string, e entry, check bool) (string, error) {
	local, err := os.ReadFile(path)
	if err == nil {
		err = verify(local, e.SHA256)
	}
	if err == nil {
		return "ok", nil
	}
	if check {
		return "", err
	}
	data, err := download(e.URL)
	if err != nil {
		return "", err
	}
	if err := verify(data, e.SHA256); err != nil {
		return "", fmt.Errorf("%s: %w", e.URL, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return "", err
	}
	return "fetch", os.Rename(tmp.Name(), path)
}

func download(url string) ([]byte, error) {
	client := http.Client{Timeout: time.Minute}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, res.Status)
	}
	return io.ReadAll(res.Body)
}

func verify(data []byte, want string) error {
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != want {
		return fmt.Errorf("sha256 %s, lock wants %s", got, want)
	}
	return nil
}
