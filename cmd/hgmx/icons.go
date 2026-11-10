package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cogikyo/hgmx/icons"
	"github.com/cogikyo/hgmx/icons/glyph"
)

const (
	iconsDir  = "views/icons"
	iconsList = "icons.json"
	iconsSum  = "icons.sum"
	iconsTags = "tags.json"
	iconsGo   = "icons.go"
	glyphPath = hgmxPath + "/icons/glyph"
)

var iconName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

var iconsCobraCmd = func() *cobra.Command {
	var check, rename bool
	var keep []string
	root := &cobra.Command{
		Use:   "icons",
		Short: "Copies the site's icons from hgmx's catalog into " + iconsDir,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !check {
				cmd.Help()
				return
			}
			os.Exit(iconsCmd(iconMode{refresh: true, check: true}, nil))
		},
	}
	root.Flags().BoolVar(&check, "check", false, "Exit 1 if sync or update would change anything; write nothing")
	update := &cobra.Command{
		Use:   "update",
		Short: "Refreshes unedited icons the catalog changed and reports edited ones",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(iconsCmd(iconMode{refresh: true, rename: rename, keep: keep}, nil))
		},
	}
	update.Flags().BoolVar(&rename, "rename", false, "Rename icons the catalog renamed, in "+iconsList+" and on disk")
	update.Flags().StringSliceVar(&keep, "keep", nil, "Keep these icons' local copies over the catalog's version, resolving a conflict")
	root.AddCommand(
		&cobra.Command{
			Use:   "add <names>...",
			Short: "Lists catalog icons in " + iconsDir + "/" + iconsList + " and copies them",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(iconsCmd(iconMode{}, args))
			},
		},
		&cobra.Command{
			Use:   "sync",
			Short: "Makes " + iconsDir + " match " + iconsList + ", keeping each copy at its locked version",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(iconsCmd(iconMode{}, nil))
			},
		},
		update,
	)
	return root
}()

type iconMode struct {
	refresh bool
	rename  bool
	check   bool
	keep    []string
}

func iconsCmd(mode iconMode, add []string) int {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod; run icons from the site's module root", slog.String("error", err.Error()))
		return 1
	}
	if s.mod.Module.Mod.Path == hgmxPath {
		log.Error("Refusing to copy icons inside the hgmx module; run icons from a site's module root")
		return 1
	}
	cat, err := readCatalog(icons.FS)
	if err != nil {
		log.Error("Failed to read hgmx's icon catalog", slog.String("error", err.Error()))
		return 1
	}
	if mode.check {
		if _, err := os.Stat(s.path(iconsDir)); err != nil {
			log.Error("Icons not set up; run go tool hgmx icons add <names>", slog.String("error", err.Error()))
			return 1
		}
	} else if err := os.MkdirAll(s.path(iconsDir), 0o755); err != nil {
		log.Error("Failed to create "+iconsDir, slog.String("error", err.Error()))
		return 1
	}
	root, err := os.OpenRoot(s.path(iconsDir))
	if err != nil {
		log.Error("Failed to open "+iconsDir, slog.String("error", err.Error()))
		return 1
	}
	defer root.Close()
	res, err := syncIcons(root, cat, mode, add, os.Stdout)
	if err != nil {
		log.Error("Icons not synced; fix the error and rerun", slog.String("error", err.Error()))
		return 1
	}
	switch {
	case mode.check && (res.drift || res.blocked > 0):
		log.Error("Icons out of date; run go tool hgmx icons update", slog.String("dir", iconsDir))
		return 1
	case res.blocked > 0:
		log.Error("Some icons need a manual fix", slog.Int("blocked", res.blocked), slog.Int("wrote", res.wrote), slog.Int("removed", res.removed))
		return 1
	case mode.check:
		log.Info("Icons up to date", slog.String("dir", iconsDir))
	default:
		log.Info("Icons synced", slog.String("dir", iconsDir), slog.Int("wrote", res.wrote), slog.Int("removed", res.removed))
	}
	return 0
}

type iconCatalog struct {
	fsys    fs.FS
	tags    map[string][]string
	renamed map[string]string
}

func readCatalog(fsys fs.FS) (iconCatalog, error) {
	c := iconCatalog{fsys: fsys, renamed: map[string]string{}}
	data, err := fs.ReadFile(fsys, "tags.json")
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c.tags); err != nil {
		return c, fmt.Errorf("tags.json: %w", err)
	}
	data, err = fs.ReadFile(fsys, "meta.json")
	if err != nil {
		return c, err
	}
	var meta struct {
		Icons map[string]struct {
			Aliases []string `json:"aliases"`
		} `json:"icons"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return c, fmt.Errorf("meta.json: %w", err)
	}
	for name, e := range meta.Icons {
		for _, a := range e.Aliases {
			c.renamed[a] = name
		}
	}
	return c, nil
}

func (c iconCatalog) files(name string) map[string][]byte {
	svg, err := fs.ReadFile(c.fsys, iconRel(name, "svg"))
	if err != nil {
		return nil
	}
	out := map[string][]byte{iconRel(name, "svg"): svg}
	if anim, err := fs.ReadFile(c.fsys, iconRel(name, "anim")); err == nil {
		out[iconRel(name, "anim")] = anim
	}
	return out
}

func iconRel(name, kind string) string {
	if kind == "svg" {
		return "svg/" + name + ".svg"
	}
	return "anim/" + name + ".js"
}

func relIcon(rel string) (string, bool) {
	if name, ok := strings.CutPrefix(rel, "svg/"); ok {
		return strings.CutSuffix(name, ".svg")
	}
	if name, ok := strings.CutPrefix(rel, "anim/"); ok {
		return strings.CutSuffix(name, ".js")
	}
	return "", false
}

type iconVerdict int

const (
	iconCurrent iconVerdict = iota
	iconCopy
	iconRefresh
	iconAdopt
	iconRemove
	iconMoved
	iconKept
	iconEdited
	iconOwn
	iconStale
	iconRenamed
	iconConflict
	iconForeign
	iconUnlisted
	iconUnknown
)

var iconVerdicts = [...]struct {
	verb, msg string
}{
	iconCopy:     {"copy", "copied from the catalog"},
	iconRefresh:  {"refresh", "unedited and changed in the catalog; refreshed"},
	iconAdopt:    {"adopt", "already matches the catalog; recorded in " + iconsSum},
	iconRemove:   {"remove", "no longer listed in " + iconsList + "; removed"},
	iconMoved:    {"rename", "renamed to %s, following the catalog"},
	iconKept:     {"keep", "local copy kept over the catalog's version; " + iconsSum + " now records the catalog's"},
	iconEdited:   {"edited", "edited since hgmx copied it; kept"},
	iconOwn:      {"own", "not in the catalog; kept as the site's own"},
	iconStale:    {"stale", "changed in the catalog; run go tool hgmx icons update"},
	iconRenamed:  {"renamed", "the catalog renamed it to %s; run go tool hgmx icons update --rename"},
	iconConflict: {"conflict", "edited here and changed in the catalog; merge by hand, keep yours with go tool hgmx icons update --keep <name>, or delete its files and run go tool hgmx icons sync to take the catalog's version"},
	iconForeign:  {"foreign", "differs from the catalog and hgmx did not copy it; keep it with go tool hgmx icons update --keep <name>, or delete its files and run go tool hgmx icons sync to take the catalog's version"},
	iconUnlisted: {"unlisted", "not listed in " + iconsList + " but edited or not copied by hgmx; list it, or delete its files"},
	iconUnknown:  {"unknown", "listed in " + iconsList + " but not in the catalog"},
}

func (v iconVerdict) blocks() bool { return v >= iconConflict }

func (v iconVerdict) drift() bool { return v != iconCurrent && v != iconEdited && v != iconOwn }

type iconItem struct {
	name    string
	verdict iconVerdict
	to      string
}

func (it iconItem) print(w io.Writer) {
	v := iconVerdicts[it.verdict]
	msg := v.msg
	if it.to != "" {
		msg = fmt.Sprintf(msg, it.to)
	}
	fmt.Fprintf(w, "%-9s %s: %s\n", v.verb, it.name, msg)
}

type iconResult struct {
	drift   bool
	blocked int
	wrote   int
	removed int
}

func syncIcons(root *os.Root, cat iconCatalog, mode iconMode, add []string, out io.Writer) (iconResult, error) {
	var res iconResult
	disk, temps, err := readIconFiles(root)
	if err != nil {
		return res, err
	}
	list, err := readIconList(root, add != nil)
	if err != nil {
		return res, err
	}
	for _, name := range add {
		if !iconName.MatchString(name) {
			return res, fmt.Errorf("%q is not an icon name", name)
		}
		if to, ok := cat.renamed[name]; ok {
			return res, fmt.Errorf("%s was renamed to %s in the catalog", name, to)
		}
		if cat.files(name) == nil {
			return res, fmt.Errorf("%s is not in the catalog", name)
		}
		if !slices.Contains(list, name) {
			list = append(list, name)
		}
	}
	slices.Sort(list)
	lock, err := readIconSum(root)
	if err != nil {
		return res, err
	}
	p, err := planIcons(list, lock, disk, cat, mode)
	if err != nil {
		return res, err
	}
	gen, err := p.generate(root, cat)
	if err != nil {
		return res, err
	}
	for _, it := range p.items {
		if it.verdict == iconCurrent {
			continue
		}
		it.print(out)
		res.drift = res.drift || it.verdict.drift()
		if it.verdict.blocks() {
			res.blocked++
		}
	}
	var writes, removes []string
	for _, rel := range slices.Sorted(maps.Keys(p.files)) {
		if old, ok := disk[rel]; !ok || !bytes.Equal(old, p.files[rel]) {
			writes = append(writes, rel)
		}
	}
	for _, rel := range slices.Sorted(maps.Keys(disk)) {
		if _, ok := p.files[rel]; !ok {
			removes = append(removes, rel)
		}
	}
	var generated []string
	for _, rel := range []string{iconsList, iconsTags, iconsGo, iconsSum} {
		data, ok := gen[rel]
		if !ok {
			continue
		}
		if old, err := root.ReadFile(rel); err != nil || !bytes.Equal(old, data) {
			generated = append(generated, rel)
		}
	}
	if len(writes)+len(removes)+len(generated)+len(temps) > 0 {
		res.drift = true
	}
	if mode.check {
		for _, rel := range slices.Concat(generated, temps) {
			fmt.Fprintf(out, "%-9s %s: out of date\n", "file", rel)
		}
		return res, nil
	}
	for _, rel := range temps {
		if err := root.Remove(rel); err != nil {
			return res, err
		}
	}
	for _, rel := range writes {
		if err := writeRoot(root, rel, p.files[rel]); err != nil {
			return res, err
		}
		res.wrote++
	}
	for _, rel := range removes {
		if err := root.Remove(rel); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return res, err
		}
		res.removed++
	}
	for _, rel := range generated {
		if err := writeRoot(root, rel, gen[rel]); err != nil {
			return res, err
		}
		res.wrote++
	}
	return res, nil
}

func readIconFiles(root *os.Root) (files map[string][]byte, temps []string, err error) {
	files = map[string][]byte{}
	for _, dir := range []string{"svg", "anim"} {
		entries, err := fs.ReadDir(root.FS(), dir)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		for _, e := range entries {
			rel := dir + "/" + e.Name()
			if temp(e.Name()) {
				temps = append(temps, rel)
				continue
			}
			if _, ok := relIcon(rel); !ok || !e.Type().IsRegular() {
				continue
			}
			if files[rel], err = root.ReadFile(rel); err != nil {
				return nil, nil, err
			}
		}
	}
	return files, temps, nil
}

func readIconList(root *os.Root, create bool) ([]string, error) {
	data, err := root.ReadFile(iconsList)
	if errors.Is(err, fs.ErrNotExist) && create {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s/%s: %w; run go tool hgmx icons add <names>", iconsDir, iconsList, err)
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("%s: want a JSON array of icon names: %w", iconsList, err)
	}
	seen := map[string]bool{}
	for _, name := range list {
		if !iconName.MatchString(name) || seen[name] {
			return nil, fmt.Errorf("%s: %q is not an icon name or is listed twice", iconsList, name)
		}
		seen[name] = true
	}
	return list, nil
}

func readIconSum(root *os.Root) (map[string]string, error) {
	out := map[string]string{}
	data, err := root.ReadFile(iconsSum)
	if errors.Is(err, fs.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		hash, rel, ok := strings.Cut(sc.Text(), "  ")
		name, isIcon := relIcon(rel)
		if _, err := hex.DecodeString(hash); !ok || err != nil || len(hash) != sha256.Size*2 || !isIcon || !iconName.MatchString(name) {
			return nil, fmt.Errorf("%s line %d: want <sha256>  svg/<name>.svg or anim/<name>.js", iconsSum, n)
		}
		out[rel] = hash
	}
	return out, sc.Err()
}

type iconPlan struct {
	items []iconItem
	list  []string
	files map[string][]byte
	lock  map[string]string
}

func planIcons(list []string, lock map[string]string, disk map[string][]byte, cat iconCatalog, mode iconMode) (*iconPlan, error) {
	p := &iconPlan{files: maps.Clone(disk), lock: maps.Clone(lock)}
	byName := map[string]bool{}
	for _, m := range []map[string]string{lock, hashes(disk)} {
		for rel := range m {
			name, _ := relIcon(rel)
			byName[name] = true
		}
	}
	for _, name := range list {
		to, ok := cat.renamed[name]
		if !ok || !mode.rename {
			p.list = append(p.list, name)
			continue
		}
		p.items = append(p.items, iconItem{name: name, verdict: iconMoved, to: to})
		p.list = append(p.list, to)
		p.move(name, to)
	}
	slices.Sort(p.list)
	p.list = slices.Compact(p.list)
	for _, alias := range slices.Sorted(maps.Keys(cat.renamed)) {
		to := cat.renamed[alias]
		if slices.Contains(list, alias) || !slices.Contains(p.list, to) || !p.owns(alias) {
			continue
		}
		if p.move(alias, to) {
			p.items = append(p.items, iconItem{name: alias, verdict: iconMoved, to: to})
		}
	}
	for _, name := range p.list {
		byName[name] = true
	}
	for _, name := range mode.keep {
		files := cat.files(name)
		if files == nil || !slices.Contains(p.list, name) || !p.owns(name) {
			return nil, fmt.Errorf("--keep %s: want a listed catalog icon with a local copy", name)
		}
		// Keep the edit, but compare future updates with today's catalog.
		for _, kind := range []string{"svg", "anim"} {
			rel := iconRel(name, kind)
			delete(p.lock, rel)
			if data, ok := files[rel]; ok {
				p.lock[rel] = hashOf(data)
			}
		}
	}
	for _, name := range slices.Sorted(maps.Keys(byName)) {
		it := iconItem{name: name, verdict: p.judge(name, cat, mode)}
		if it.verdict == iconRenamed {
			it.to = cat.renamed[name]
		}
		if it.verdict == iconEdited && slices.Contains(mode.keep, name) {
			it.verdict = iconKept
		}
		switch it.verdict {
		case iconCopy, iconRefresh, iconAdopt:
			p.set(name, cat.files(name))
		case iconRemove:
			p.set(name, nil)
		}
		p.items = append(p.items, it)
	}
	return p, nil
}

func (p *iconPlan) move(from, to string) bool {
	for _, kind := range []string{"svg", "anim"} {
		src, dest := iconRel(from, kind), iconRel(to, kind)
		if a, ok := p.files[src]; ok {
			if b, ok := p.files[dest]; ok && !bytes.Equal(a, b) {
				return false
			}
		}
		if a, ok := p.lock[src]; ok {
			if b, ok := p.lock[dest]; ok && a != b {
				return false
			}
		}
	}
	for _, kind := range []string{"svg", "anim"} {
		src, dest := iconRel(from, kind), iconRel(to, kind)
		if data, ok := p.files[src]; ok {
			p.files[dest] = data
			delete(p.files, src)
		}
		if hash, ok := p.lock[src]; ok {
			p.lock[dest] = hash
			delete(p.lock, src)
		}
	}
	return true
}

func (p *iconPlan) owns(name string) bool {
	for _, kind := range []string{"svg", "anim"} {
		if _, ok := p.files[iconRel(name, kind)]; ok {
			return true
		}
		if _, ok := p.lock[iconRel(name, kind)]; ok {
			return true
		}
	}
	return false
}

func (p *iconPlan) icon(name string) (local, locked map[string]string) {
	local, locked = map[string]string{}, map[string]string{}
	for _, kind := range []string{"svg", "anim"} {
		rel := iconRel(name, kind)
		if data, ok := p.files[rel]; ok {
			local[rel] = hashOf(data)
		}
		if hash, ok := p.lock[rel]; ok {
			locked[rel] = hash
		}
	}
	return local, locked
}

func (p *iconPlan) set(name string, files map[string][]byte) {
	for _, kind := range []string{"svg", "anim"} {
		rel := iconRel(name, kind)
		delete(p.lock, rel)
		if data, ok := files[rel]; ok {
			p.files[rel] = data
			p.lock[rel] = hashOf(data)
		} else {
			delete(p.files, rel)
		}
	}
}

func (p *iconPlan) judge(name string, cat iconCatalog, mode iconMode) iconVerdict {
	l, k := p.icon(name)
	files := cat.files(name)
	c := hashes(files)
	listed := slices.Contains(p.list, name)
	rels := slices.Concat(slices.Collect(maps.Keys(l)), slices.Collect(maps.Keys(k)), slices.Collect(maps.Keys(c)))

	if !listed {
		if len(l) == 0 && len(k) == 0 {
			return iconCurrent
		}
		for _, rel := range rels {
			if l[rel] != "" && l[rel] != k[rel] && l[rel] != c[rel] {
				return iconUnlisted
			}
		}
		return iconRemove
	}
	if files == nil {
		switch {
		case cat.renamed[name] != "":
			return iconRenamed
		case len(l) == 0:
			return iconUnknown
		}
		return iconOwn
	}
	if len(l) == 0 {
		return iconCopy
	}
	var edited, stale bool
	for _, rel := range rels {
		switch {
		case l[rel] == c[rel]:
		case l[rel] != k[rel]:
			edited = true
		case len(k) > 0:
			stale = true
		}
	}
	switch {
	case edited && len(k) == 0:
		return iconForeign
	case edited && maps.Equal(k, c):
		return iconEdited
	case edited:
		return iconConflict
	case stale && mode.refresh:
		return iconRefresh
	case stale:
		return iconStale
	case !maps.Equal(l, c):
		return iconCopy
	case !maps.Equal(k, c):
		return iconAdopt
	}
	return iconCurrent
}

func hashes(files map[string][]byte) map[string]string {
	out := make(map[string]string, len(files))
	for rel, data := range files {
		out[rel] = hashOf(data)
	}
	return out
}

func (p *iconPlan) generate(root *os.Root, cat iconCatalog) (map[string][]byte, error) {
	var local map[string][]string
	if data, err := root.ReadFile(iconsTags); err == nil {
		json.Unmarshal(data, &local)
	}
	tags := map[string][]string{}
	var names []string
	var svg, anim bool
	for _, name := range p.list {
		if _, ok := p.files[iconRel(name, "svg")]; !ok {
			continue
		}
		_, animated := p.files[iconRel(name, "anim")]
		svg, anim = true, anim || animated
		names = append(names, name)
		if t, ok := cat.tags[name]; ok {
			tags[name] = t
		} else {
			tags[name] = local[name]
		}
		if tags[name] == nil {
			tags[name] = []string{}
		}
	}

	list, err := json.MarshalIndent(p.list, "", "\t")
	if err != nil {
		return nil, err
	}
	if p.list == nil {
		list = []byte("[]")
	}
	tagsJSON, err := json.MarshalIndent(tags, "", "\t")
	if err != nil {
		return nil, err
	}
	code, err := iconsSource(names, svg, anim)
	if err != nil {
		return nil, err
	}
	var sum bytes.Buffer
	for _, rel := range slices.Sorted(maps.Keys(p.lock)) {
		fmt.Fprintf(&sum, "%s  %s\n", p.lock[rel], rel)
	}
	return map[string][]byte{
		iconsList: append(list, '\n'),
		iconsTags: append(tagsJSON, '\n'),
		iconsGo:   code,
		iconsSum:  sum.Bytes(),
	}, nil
}

func iconsSource(names []string, svg, anim bool) ([]byte, error) {
	embed := []string{iconsTags}
	if anim {
		embed = append([]string{"anim"}, embed...)
	}
	if svg {
		embed = append([]string{"svg"}, embed...)
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, `// Code generated by go tool hgmx icons; DO NOT EDIT.

package icons

import (
	"context"
	"embed"
	"io"

	%q
)

//go:embed %s
var files embed.FS

var Set = glyph.MustLoad(files)

type Name string

func (n Name) Render(_ context.Context, w io.Writer) error {
	return Set.Render(w, string(n))
}
`, glyphPath, strings.Join(embed, " "))
	if len(names) > 0 {
		b.WriteString("\nconst (\n")
	}
	seen := map[string]string{"Set": "Set", "Name": "Name"}
	for _, name := range names {
		id := glyph.Ident(name)
		if !token.IsIdentifier(id) {
			return nil, fmt.Errorf("%s has no Go identifier", name)
		}
		if prev, ok := seen[id]; ok {
			return nil, fmt.Errorf("%s and %s both map to %s", prev, name, id)
		}
		seen[id] = name
		fmt.Fprintf(&b, "\t%s Name = %q\n", id, name)
	}
	if len(names) > 0 {
		b.WriteString(")\n")
	}
	return format.Source(b.Bytes())
}
