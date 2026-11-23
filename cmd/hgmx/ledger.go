package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type drift int

const (
	current drift = iota
	missing
	stale
	dropped
	edited
	foreign
)

func judge(want, have, rec string) drift {
	switch {
	case have == want:
		return current
	case have == "":
		return missing
	case rec == "":
		return foreign
	case have != rec:
		return edited
	case want == "":
		return dropped
	}
	return stale
}

func (d drift) blocks() bool { return d == edited || d == foreign }

type ledger map[string]string

func readLedger(root *os.Root, name, form string, valid func(rel string) bool) (l ledger, found bool, err error) {
	data, err := root.ReadFile(name)
	if errors.Is(err, fs.ErrNotExist) {
		return ledger{}, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	l = ledger{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		hash, rel, ok := strings.Cut(sc.Text(), "  ")
		if _, err := hex.DecodeString(hash); !ok || err != nil || len(hash) != sha256.Size*2 || !valid(rel) {
			return nil, true, fmt.Errorf("line %d: want <sha256>  %s", n, form)
		}
		l[rel] = hash
	}
	return l, true, sc.Err()
}

func (l ledger) format() []byte {
	var b bytes.Buffer
	for _, rel := range slices.Sorted(maps.Keys(l)) {
		fmt.Fprintf(&b, "%s  %s\n", l[rel], rel)
	}
	return b.Bytes()
}

func hashOf(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hashes(files map[string][]byte) map[string]string {
	out := make(map[string]string, len(files))
	for rel, data := range files {
		out[rel] = hashOf(data)
	}
	return out
}

type copyItem struct {
	rel   string
	data  []byte
	seen  string
	drift drift
}

type copyPlan struct {
	root    *os.Root
	sum     string
	files   []copyItem
	prune   []string
	pending ledger
}

func planCopy(root *os.Root, sum string, rec ledger, src map[string][]byte) (*copyPlan, error) {
	want := ledger(hashes(src))
	c := &copyPlan{root: root, sum: sum}
	if !maps.Equal(rec, want) {
		c.pending = want
	}
	rels := slices.Collect(maps.Keys(src))
	for rel := range rec {
		if _, ok := src[rel]; !ok {
			rels = append(rels, rel)
			c.prune = append(c.prune, path.Dir(rel))
		}
	}
	slices.Sort(rels)
	for _, rel := range rels {
		f := copyItem{rel: rel, data: src[rel]}
		data, err := root.ReadFile(rel)
		switch {
		case err == nil:
			f.seen = hashOf(data)
		case !errors.Is(err, fs.ErrNotExist):
			return nil, err
		}
		f.drift = judge(want[rel], f.seen, rec[rel])
		c.files = append(c.files, f)
	}
	return c, nil
}

func (c *copyPlan) recorded() bool { return c.pending == nil }

type changed struct{ rel string }

func (c changed) Error() string { return c.rel + " changed during the sync" }

func (c *copyPlan) apply() (wrote, removed int, err error) {
	for _, f := range c.files {
		if f.drift != missing && f.drift != stale && f.drift != dropped {
			continue
		}
		seen := ""
		data, err := c.root.ReadFile(f.rel)
		switch {
		case err == nil:
			seen = hashOf(data)
		case !errors.Is(err, fs.ErrNotExist):
			return wrote, removed, err
		}
		if seen != f.seen {
			return wrote, removed, changed{f.rel}
		}
		if f.drift == dropped {
			if err := c.root.Remove(f.rel); err != nil {
				return wrote, removed, err
			}
			removed++
			continue
		}
		if err := writeRoot(c.root, f.rel, f.data); err != nil {
			return wrote, removed, err
		}
		wrote++
	}
	for _, dir := range c.prune {
		for ; dir != "."; dir = path.Dir(dir) {
			entries, err := fs.ReadDir(c.root.FS(), dir)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return wrote, removed, err
			}
			if len(entries) > 0 {
				break
			}
			if err := c.root.Remove(dir); err != nil {
				return wrote, removed, err
			}
		}
	}
	if c.recorded() {
		return wrote, removed, nil
	}
	if err := writeRoot(c.root, c.sum, c.pending.format()); err != nil {
		return wrote, removed, err
	}
	return wrote + 1, removed, nil
}

func tempPrefix(rel string) string {
	base := path.Base(rel)
	if len(base) > 32 {
		cut := 0
		for i := range base {
			if i > 32 {
				break
			}
			cut = i
		}
		base = base[:cut]
	}
	return path.Join(path.Dir(rel), "."+base+".hgmx-")
}

func tempName(rel string) string { return tempPrefix(rel) + rand.Text() }

func tempGlob(rel string) string { return tempPrefix(rel) + "*" }

func temp(name string) bool {
	return strings.HasPrefix(name, ".") && strings.Contains(name, ".hgmx-")
}

func writeRoot(root *os.Root, rel string, data []byte) (err error) {
	if err := root.MkdirAll(path.Dir(rel), 0o755); err != nil {
		return err
	}
	tmp := tempName(rel)
	f, err := root.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			root.Remove(tmp)
		}
	}()
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return root.Rename(tmp, rel)
}

func writeAtomic(name string, data []byte) error {
	root, err := os.OpenRoot(filepath.Dir(name))
	if err != nil {
		return err
	}
	defer root.Close()
	return writeRoot(root, filepath.Base(name), data)
}
