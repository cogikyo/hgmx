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

	"github.com/cogikyo/hgmx/icons/glyph"
)

const root = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">`

type source struct {
	Package string            `json:"package"`
	Version string            `json:"version"`
	URL     string            `json:"url"`
	SHA256  string            `json:"sha256"`
	Icons   map[string]string `json:"icons,omitempty"`
}

type lock struct {
	Lucide   source `json:"lucide"`
	Lab      source `json:"lab"`
	Metadata source `json:"metadata"`
}

func (l *lock) catalogs() map[string]*source {
	return map[string]*source{"lucide": &l.Lucide, "lab": &l.Lab}
}

func (l *lock) owners() (map[string]string, error) {
	owners := map[string]string{}
	catalogs := l.catalogs()
	for _, src := range slices.Sorted(maps.Keys(catalogs)) {
		for name := range catalogs[src].Icons {
			if prev, ok := owners[name]; ok {
				return nil, fmt.Errorf("%s: locked by both %s and %s", name, prev, src)
			}
			owners[name] = src
		}
	}
	return owners, nil
}

type meta struct {
	Categories []category       `json:"categories"`
	Icons      map[string]entry `json:"icons"`
}

type category struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type entry struct {
	Source     string   `json:"source"`
	Categories []string `json:"categories"`
	Aliases    []string `json:"aliases,omitempty"`
}

type upstream struct {
	svgs       map[string]map[string][]byte
	tags       map[string][]string
	categories map[string][]string
	titles     map[string]string
	aliases    map[string]string
	license    []byte
}

var originals = map[string][]string{
	"github":       {"development", "social"},
	"question":     {"accessibility", "text", "notifications"},
	"user-contact": {"account", "communication", "social"},
}

var labCategories = map[string]string{"maths": "math"}

func main() {
	check := flag.Bool("check", false, "verify icon and generated files against lock.json without downloading")
	update := flag.Bool("update", false, "lock every icon in the pinned sources and refresh their metadata")
	flag.Parse()
	dir := cmp.Or(flag.Arg(0), ".")
	if *check && *update {
		fmt.Fprintln(os.Stderr, "-check and -update conflict")
		os.Exit(2)
	}
	if err := run(dir, *check, *update); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string, check, update bool) error {
	var l lock
	if err := readJSON(filepath.Join(dir, "lock.json"), &l); err != nil {
		return err
	}
	tags := map[string][]string{}
	if err := readJSON(filepath.Join(dir, "tags.json"), &tags); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	var m meta
	if err := readJSON(filepath.Join(dir, "meta.json"), &m); err != nil && !errors.Is(err, os.ErrNotExist) {
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

	if update {
		u, err := load()
		if err != nil {
			return err
		}
		if err := relock(dir, &l, u); err != nil {
			return err
		}
	}
	owners, err := l.owners()
	if err != nil {
		return err
	}

	failed := false
	for _, name := range slices.Sorted(maps.Keys(owners)) {
		src := owners[name]
		status, err := sync(dir, src, name, l.catalogs()[src].Icons[name], check, load)
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
	slices.Sort(names)
	_, licenseErr := os.Stat(filepath.Join(dir, "LICENSE"))
	unlisted := slices.ContainsFunc(names, func(name string) bool {
		_, ok := m.Icons[name]
		return !ok
	})
	if !check && (licenseErr != nil || unlisted) {
		if _, err := load(); err != nil {
			return err
		}
	}
	if up != nil {
		for name := range owners {
			tags[name] = up.tags[name]
		}
		m = up.meta(names)
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
	if err := m.settle(names, owners); err != nil {
		return err
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
	metaData, err := json.MarshalIndent(m, "", "\t")
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
		{"meta.json", append(metaData, '\n')},
	}
	if up != nil {
		outputs = append(outputs, struct {
			file string
			data []byte
		}{"LICENSE", up.license})
	} else if licenseErr != nil {
		fmt.Fprintf(os.Stderr, "FAIL   LICENSE: %v\n", licenseErr)
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
	counts := map[string]int{}
	for _, e := range m.Icons {
		counts[e.Source]++
	}
	fmt.Printf("%d icons: %d from %s %s, %d from %s %s, %d hgmx\n", len(names),
		counts["lucide"], l.Lucide.Package, l.Lucide.Version,
		counts["lab"], l.Lab.Package, l.Lab.Version, counts["hgmx"])
	return nil
}

func relock(dir string, l *lock, u *upstream) error {
	catalogs := l.catalogs()
	srcs := slices.Sorted(maps.Keys(catalogs))
	var errs []error
	for _, src := range srcs {
		s := catalogs[src]
		for _, name := range slices.Sorted(maps.Keys(s.Icons)) {
			if _, ok := u.svgs[src][name]; ok {
				continue
			}
			if to, ok := u.aliases[name]; ok {
				errs = append(errs, fmt.Errorf("%s: %s %s renames it %s; move svg/%[1]s.svg, anim/%[1]s.js and its lock entry", name, s.Package, s.Version, to))
			} else {
				errs = append(errs, fmt.Errorf("%s: not in %s %s", name, s.Package, s.Version))
			}
		}
	}
	for i, src := range srcs {
		for _, other := range srcs[i+1:] {
			for _, name := range slices.Sorted(maps.Keys(u.svgs[src])) {
				if _, ok := u.svgs[other][name]; ok {
					errs = append(errs, fmt.Errorf("%s: in both %s and %s", name, src, other))
				}
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	for _, src := range srcs {
		s := catalogs[src]
		if s.Icons == nil {
			s.Icons = map[string]string{}
		}
		for name, data := range u.svgs[src] {
			if _, ok := s.Icons[name]; !ok {
				if _, err := os.Stat(svgPath(dir, name)); err == nil {
					continue
				}
			}
			out, err := normalize(data)
			if err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			s.Icons[name] = digest(out)
		}
	}
	return nil
}

func sync(dir, src, name, want string, check bool, load func() (*upstream, error)) (string, error) {
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
	data, ok := u.svgs[src][name]
	if !ok {
		return "", fmt.Errorf("not in %s", src)
	}
	out, err := normalize(data)
	if err != nil {
		return "", err
	}
	if got := digest(out); got != want {
		return "", fmt.Errorf("normalized sha256 %s, lock wants %s", got, want)
	}
	return "fetch", write(p, out)
}

func (u *upstream) meta(names []string) meta {
	m := meta{Icons: map[string]entry{}}
	for _, name := range slices.Sorted(maps.Keys(u.titles)) {
		m.Categories = append(m.Categories, category{name, u.titles[name]})
	}
	aliases := map[string][]string{}
	for alias, name := range u.aliases {
		aliases[name] = append(aliases[name], alias)
	}
	for _, name := range names {
		cats, ok := u.categories[name]
		if !ok {
			cats = originals[name]
		}
		m.Icons[name] = entry{Categories: cats, Aliases: aliases[name]}
	}
	return m
}

func (m *meta) settle(names []string, owners map[string]string) error {
	slices.SortFunc(m.Categories, func(a, b category) int { return strings.Compare(a.Name, b.Name) })
	known := map[string]bool{}
	for _, c := range m.Categories {
		known[c.Name] = true
	}
	icons := map[string]entry{}
	var unlisted, unknown []string
	for _, name := range names {
		e, ok := m.Icons[name]
		if !ok {
			unlisted = append(unlisted, name)
			continue
		}
		e.Source = cmp.Or(owners[name], "hgmx")
		if e.Categories == nil {
			e.Categories = []string{}
		}
		slices.Sort(e.Aliases)
		for _, c := range e.Categories {
			if !known[c] {
				unknown = append(unknown, name+": "+c)
			}
		}
		icons[name] = e
	}
	if len(unlisted) > 0 {
		return fmt.Errorf("meta.json has no entry for %s", strings.Join(unlisted, ", "))
	}
	if len(unknown) > 0 {
		return fmt.Errorf("meta.json lists no category for %s", strings.Join(unknown, ", "))
	}
	m.Icons = icons
	return nil
}

func fetch(l lock) (*upstream, error) {
	u := &upstream{
		svgs:    map[string]map[string][]byte{},
		tags:    map[string][]string{},
		titles:  map[string]string{},
		aliases: map[string]string{},
	}
	licenses := map[string][]byte{}
	for src, s := range l.catalogs() {
		u.svgs[src] = map[string][]byte{}
		err := walk(*s, func(name string, r io.Reader) error {
			var err error
			switch {
			case path.Dir(name) == "icons" && path.Ext(name) == ".svg":
				u.svgs[src][strings.TrimSuffix(path.Base(name), ".svg")], err = io.ReadAll(r)
			case name == "LICENSE":
				licenses[src], err = io.ReadAll(r)
			case src == "lucide" && name == "tags.json":
				if err = json.NewDecoder(r).Decode(&u.tags); err != nil {
					err = fmt.Errorf("tags.json: %w", err)
				}
			}
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	if licenses["lucide"] == nil || licenses["lab"] == nil {
		return nil, errors.New("upstream LICENSE missing")
	}
	u.license = fmt.Appendf(slices.Clone(licenses["lucide"]),
		"\n---\n\nIcons with source \"lab\" in meta.json come from Lucide Lab (%s %s) under this license:\n\n%s",
		l.Lab.Package, l.Lab.Version, licenses["lab"])

	raw := map[string][]string{}
	canonical := map[string]bool{}
	err := walk(l.Metadata, func(name string, r io.Reader) error {
		dir, file := path.Split(name)
		icon, ok := strings.CutSuffix(file, ".json")
		if !ok || (dir != "icons/" && dir != "lab/" && dir != "categories/") {
			return nil
		}
		var rec struct {
			Title      string   `json:"title"`
			Tags       []string `json:"tags"`
			Categories []string `json:"categories"`
			Aliases    []struct {
				Name string `json:"name"`
			} `json:"aliases"`
		}
		if err := json.NewDecoder(r).Decode(&rec); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		for _, a := range rec.Aliases {
			u.aliases[a.Name] = icon
		}
		switch dir {
		case "icons/":
			canonical[icon] = true
			raw[icon] = rec.Categories
		case "lab/":
			for i, c := range rec.Categories {
				rec.Categories[i] = cmp.Or(labCategories[c], c)
			}
			raw[icon] = rec.Categories
			u.tags[icon] = rec.Tags
		case "categories/":
			u.titles[icon] = rec.Title
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	maps.DeleteFunc(u.svgs["lucide"], func(name string, _ []byte) bool { return !canonical[name] })
	if u.categories, err = extend(raw); err != nil {
		return nil, err
	}

	var bare []string
	for _, svgs := range u.svgs {
		for name := range svgs {
			if u.tags[name] == nil || u.categories[name] == nil {
				bare = append(bare, name)
			}
		}
	}
	if len(bare) > 0 {
		slices.Sort(bare)
		return nil, fmt.Errorf("%s %s has no tags or categories for %s", l.Metadata.Package, l.Metadata.Version, strings.Join(bare, ", "))
	}
	for _, cats := range u.categories {
		for _, c := range cats {
			if _, ok := u.titles[c]; !ok {
				u.titles[c] = strings.ToUpper(c[:1]) + c[1:]
			}
		}
	}
	return u, nil
}

func extend(raw map[string][]string) (map[string][]string, error) {
	out := map[string][]string{}
	var visit func(name string, stack []string) ([]string, error)
	visit = func(name string, stack []string) ([]string, error) {
		if vals, ok := out[name]; ok {
			return vals, nil
		}
		if slices.Contains(stack, name) {
			return nil, fmt.Errorf("categories: %s extends itself", strings.Join(append(stack, name), " -> "))
		}
		entries, ok := raw[name]
		if !ok {
			return nil, fmt.Errorf("categories: %s extends unknown icon %s", stack[len(stack)-1], name)
		}
		vals := []string{}
		for _, v := range entries {
			inherited := []string{v}
			if parent, ok := strings.CutPrefix(v, "$extends:"); ok {
				var err error
				if inherited, err = visit(parent, append(stack, name)); err != nil {
					return nil, err
				}
			}
			for _, c := range inherited {
				if !slices.Contains(vals, c) {
					vals = append(vals, c)
				}
			}
		}
		out[name] = vals
		return vals, nil
	}
	for _, name := range slices.Sorted(maps.Keys(raw)) {
		if _, err := visit(name, nil); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func walk(s source, file func(name string, r io.Reader) error) error {
	client := http.Client{Timeout: 5 * time.Minute}
	res, err := client.Get(s.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", s.URL, res.Status)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if got := digest(data); got != s.SHA256 {
		return fmt.Errorf("%s: sha256 %s, lock wants %s", s.URL, got, s.SHA256)
	}
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_, name, _ := strings.Cut(h.Name, "/")
		if err := file(name, tr); err != nil {
			return fmt.Errorf("%s: %w", s.URL, err)
		}
	}
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
		id := glyph.Ident(name)
		if prev, ok := seen[id]; ok {
			return nil, fmt.Errorf("%s and %s both map to %s", prev, name, id)
		}
		seen[id] = name
		fmt.Fprintf(&b, "\t%s Name = %q\n", id, name)
	}
	b.WriteString(")\n")
	return format.Source(b.Bytes())
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
