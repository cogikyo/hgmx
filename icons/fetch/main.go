package main

import (
	"archive/tar"
	"bytes"
	"cmp"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"html"
	"io"
	"maps"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const root = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">`

type lock struct {
	Package string            `json:"package"`
	Version string            `json:"version"`
	URL     string            `json:"url"`
	SHA256  string            `json:"sha256"`
	Icons   map[string]string `json:"icons"`
}

type upstream struct {
	svgs    map[string][]byte
	tags    map[string][]string
	license []byte
}

func main() {
	check := flag.Bool("check", false, "verify icon files and names.go against lock.json without downloading")
	add := flag.String("add", "", "comma-separated upstream icons to import and lock")
	flag.Parse()
	dir := cmp.Or(flag.Arg(0), ".")
	if err := run(dir, *check, *add); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string, check bool, add string) error {
	var l lock
	if err := readJSON(filepath.Join(dir, "lock.json"), &l); err != nil {
		return err
	}
	if l.Icons == nil {
		l.Icons = map[string]string{}
	}
	tags := map[string][]string{}
	if err := readJSON(filepath.Join(dir, "tags.json"), &tags); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	var up *upstream
	load := func() (*upstream, error) {
		if up != nil {
			return up, nil
		}
		var err error
		up, err = fetch(l)
		return up, err
	}

	var added []string
	if add != "" {
		added = strings.Split(add, ",")
	}
	for _, name := range added {
		if _, ok := l.Icons[name]; ok {
			continue
		}
		if _, err := os.Stat(svgPath(dir, name)); err == nil {
			return fmt.Errorf("%s: an unlocked icon already owns svg/%s.svg", name, name)
		}
		u, err := load()
		if err != nil {
			return err
		}
		src, ok := u.svgs[name]
		if !ok {
			return fmt.Errorf("%s: not in %s %s", name, l.Package, l.Version)
		}
		out, err := normalize(src)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		l.Icons[name] = digest(out)
	}

	failed := false
	for _, name := range slices.Sorted(maps.Keys(l.Icons)) {
		status, err := sync(dir, name, l.Icons[name], check, load)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL   %s: %v\n", name, err)
			failed = true
			continue
		}
		if status != "ok" {
			fmt.Printf("%-6s %s\n", status, name)
		}
	}
	if failed {
		return errors.New("icons out of sync with lock.json")
	}

	files, err := filepath.Glob(filepath.Join(dir, "svg", "*.svg"))
	if err != nil {
		return err
	}
	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".svg"))
	}
	if up != nil {
		for name := range l.Icons {
			tags[name] = up.tags[name]
		}
	}
	for name := range tags {
		if !slices.Contains(names, name) {
			delete(tags, name)
		}
	}
	var untagged []string
	for _, name := range names {
		if _, ok := tags[name]; !ok {
			untagged = append(untagged, name)
		}
	}
	if len(untagged) > 0 {
		return fmt.Errorf("tags.json has no entry for %s", strings.Join(untagged, ", "))
	}

	code, err := generate(names)
	if err != nil {
		return err
	}
	tagData, err := json.MarshalIndent(tags, "", "\t")
	if err != nil {
		return err
	}
	lockData, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		return err
	}
	outputs := []struct {
		file string
		data []byte
	}{
		{"names.go", code},
		{"tags.json", append(tagData, '\n')},
		{"lock.json", append(lockData, '\n')},
	}
	if _, err := os.Stat(filepath.Join(dir, "LICENSE")); err != nil && !check {
		if _, err := load(); err != nil {
			return err
		}
	}
	if up != nil {
		outputs = append(outputs, struct {
			file string
			data []byte
		}{"LICENSE", up.license})
	} else if _, err := os.Stat(filepath.Join(dir, "LICENSE")); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL   LICENSE: %v\n", err)
		failed = true
	}
	for _, o := range outputs {
		p := filepath.Join(dir, o.file)
		current, _ := os.ReadFile(p)
		if bytes.Equal(current, o.data) {
			continue
		}
		if check {
			fmt.Fprintf(os.Stderr, "FAIL   %s is stale; run go generate ./icons\n", o.file)
			failed = true
			continue
		}
		if err := write(p, o.data); err != nil {
			return err
		}
		fmt.Printf("write  %s\n", o.file)
	}
	if failed {
		return errors.New("generated files out of date")
	}
	fmt.Printf("%d icons, %d from %s %s\n", len(names), len(l.Icons), l.Package, l.Version)
	return nil
}

func sync(dir, name, want string, check bool, load func() (*upstream, error)) (string, error) {
	p := svgPath(dir, name)
	local, err := os.ReadFile(p)
	if err == nil {
		if digest(local) == want {
			return "ok", nil
		}
		return "edited", nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if check {
		return "", errors.New("missing")
	}
	u, err := load()
	if err != nil {
		return "", err
	}
	src, ok := u.svgs[name]
	if !ok {
		return "", errors.New("not upstream")
	}
	out, err := normalize(src)
	if err != nil {
		return "", err
	}
	if got := digest(out); got != want {
		return "", fmt.Errorf("normalized sha256 %s, lock wants %s", got, want)
	}
	return "fetch", write(p, out)
}

func fetch(l lock) (*upstream, error) {
	client := http.Client{Timeout: 5 * time.Minute}
	res, err := client.Get(l.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", l.URL, res.Status)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if got := digest(data); got != l.SHA256 {
		return nil, fmt.Errorf("%s: sha256 %s, lock wants %s", l.URL, got, l.SHA256)
	}
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	u := &upstream{svgs: map[string][]byte{}}
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		dir, file := path.Split(h.Name)
		switch {
		case dir == "package/icons/" && strings.HasSuffix(file, ".svg"):
			b, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			u.svgs[strings.TrimSuffix(file, ".svg")] = b
		case h.Name == "package/LICENSE":
			if u.license, err = io.ReadAll(tr); err != nil {
				return nil, err
			}
		case h.Name == "package/tags.json":
			if err := json.NewDecoder(tr).Decode(&u.tags); err != nil {
				return nil, fmt.Errorf("tags.json: %w", err)
			}
		}
	}
	return u, nil
}

func normalize(src []byte) ([]byte, error) {
	d := xml.NewDecoder(bytes.NewReader(src))
	var b bytes.Buffer
	depth := 0
	open := false
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				if t.Name.Local != "svg" {
					return nil, fmt.Errorf("root is <%s>", t.Name.Local)
				}
				b.WriteString(root)
				depth++
				continue
			}
			if open {
				b.WriteString(">")
			}
			b.WriteString("\n" + strings.Repeat("\t", depth) + "<" + t.Name.Local)
			for _, a := range t.Attr {
				fmt.Fprintf(&b, ` %s="%s"`, a.Name.Local, html.EscapeString(a.Value))
			}
			open = true
			depth++
		case xml.EndElement:
			depth--
			switch {
			case depth == 0:
				b.WriteString("\n</svg>\n")
			case open:
				b.WriteString("></" + t.Name.Local + ">")
			default:
				b.WriteString("\n" + strings.Repeat("\t", depth) + "</" + t.Name.Local + ">")
			}
			open = false
		}
	}
	if depth != 0 || b.Len() == 0 {
		return nil, errors.New("unbalanced svg")
	}
	return b.Bytes(), nil
}

func generate(names []string) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("// Code generated by go run ./fetch; DO NOT EDIT.\n\npackage icons\n\nconst (\n")
	seen := map[string]string{}
	for _, name := range names {
		id := ident(name)
		if prev, ok := seen[id]; ok {
			return nil, fmt.Errorf("%s and %s both map to %s", prev, name, id)
		}
		seen[id] = name
		fmt.Fprintf(&b, "\t%s Name = %q\n", id, name)
	}
	b.WriteString(")\n")
	return format.Source(b.Bytes())
}

func ident(name string) string {
	var b strings.Builder
	for part := range strings.SplitSeq(name, "-") {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	return b.String()
}

func svgPath(dir, name string) string {
	return filepath.Join(dir, "svg", name+".svg")
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func readJSON(p string, v any) error {
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s: %w", p, err)
	}
	return nil
}

func write(p string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), filepath.Base(p)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), p)
}
