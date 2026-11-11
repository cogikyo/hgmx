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
	"log/slog"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cogikyo/hgmx/skills"
)

const (
	skillsDir    = ".agents/skills"
	skillsSum    = skillsDir + "/hgmx.sum"
	claudeDir    = ".claude"
	claudeSkills = claudeDir + "/skills"
	siteSkill    = "SITE.md"
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

type skillItem struct {
	path  string
	skill string
	data  []byte
	link  string
	seen  string
	drift drift
}

type changed struct{ it skillItem }

func (c changed) Error() string {
	return c.it.path + " changed during the sync"
}

func (it skillItem) blocks() bool {
	return it.drift == edited || it.drift == foreign
}

func (it skillItem) finding() finding {
	f := finding{file: it.path, rule: ruleHgmxSkills}
	edit := fmt.Sprintf("move site-specific changes into %s/%s/%s, delete this file, then run go tool hgmx skills", skillsDir, it.skill, siteSkill)
	switch {
	case it.link != "" && it.drift == missing:
		f.msg = "missing link to " + it.link
	case it.link != "" && it.drift == dropped:
		f.msg = "links a skill hgmx no longer ships"
	case it.link != "":
		f.msg = "not hgmx's link to " + it.link
		f.fix = "move it aside, then run go tool hgmx skills"
	case it.drift == missing:
		f.msg = "missing"
	case it.drift == stale:
		f.msg = "stale; hgmx ships a different version"
	case it.drift == dropped:
		f.msg = "no longer shipped by hgmx"
	case it.drift == edited:
		f.msg = "edited since hgmx wrote it"
		f.fix = edit
	case it.drift == foreign:
		f.msg = "not written by hgmx and differs from its copy"
		f.fix = edit
	}
	return f
}

var skillsMissing = finding{file: skillsDir + "/", rule: ruleHgmxSkills, msg: "hgmx skills not installed"}

type skillSync struct {
	root      *os.Root
	items     []skillItem
	prune     []string
	sum, want []byte
	installed bool
}

func openSkills(s site) (*os.Root, error) {
	if err := os.MkdirAll(s.path(skillsDir), 0o755); err != nil {
		return nil, err
	}
	return os.OpenRoot(s.path(skillsDir))
}

func planSkills(s site, root *os.Root, src fs.FS) (*skillSync, error) {
	want := map[string][]byte{}
	err := fs.WalkDir(src, ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Base(rel) == siteSkill {
			return err
		}
		want[rel], err = fs.ReadFile(src, rel)
		return err
	})
	if err != nil {
		return nil, err
	}
	sum, err := root.ReadFile(path.Base(skillsSum))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	ps := &skillSync{root: root, sum: sum, want: formatSum(want), installed: err == nil}
	recorded, err := parseSum(sum)
	if err != nil {
		return nil, fmt.Errorf("%s: %w; delete it, then run go tool hgmx skills", skillsSum, err)
	}

	rels := slices.Collect(maps.Keys(want))
	for rel := range recorded {
		if _, ok := want[rel]; !ok {
			rels = append(rels, rel)
			ps.prune = append(ps.prune, path.Dir(rel))
		}
	}
	slices.Sort(rels)
	for _, rel := range rels {
		data, ok := want[rel]
		cur, err := root.ReadFile(rel)
		exists := err == nil
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s/%s: %w", skillsDir, rel, err)
		}
		hash, has := recorded[rel]
		it := skillItem{path: skillsDir + "/" + rel, skill: skillName(rel), data: data}
		if exists {
			it.seen = hashOf(cur)
		}
		switch {
		case ok && !exists:
			it.drift = missing
		case ok && bytes.Equal(cur, data), !exists:
			it.drift = current
		case !has:
			it.drift = foreign
		case it.seen != hash:
			it.drift = edited
		case ok:
			it.drift = stale
		default:
			it.drift = dropped
		}
		ps.items = append(ps.items, it)
	}

	fi, err := os.Stat(s.path(claudeDir))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	if err != nil || !fi.IsDir() {
		return ps, nil
	}
	// Resolve the link's physical parent; .claude/skills may itself be a symlink.
	agents, err := filepath.EvalSymlinks(s.path(skillsDir))
	if err != nil {
		return nil, err
	}
	links, err := filepath.EvalSymlinks(s.path(claudeSkills))
	if errors.Is(err, fs.ErrNotExist) {
		links, err = filepath.EvalSymlinks(s.path(claudeDir))
		links = filepath.Join(links, "skills")
	}
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for rel := range recorded {
		names[skillName(rel)] = false
	}
	for rel := range want {
		names[skillName(rel)] = true
	}
	for _, name := range slices.Sorted(maps.Keys(names)) {
		it, err := linkItem(s, filepath.Join(agents, name), links, name, names[name])
		if err != nil {
			return nil, err
		}
		ps.items = append(ps.items, it)
	}
	return ps, nil
}

func linkItem(s site, skill, links, name string, wanted bool) (skillItem, error) {
	target, err := filepath.Rel(links, skill)
	if err != nil {
		return skillItem{}, err
	}
	it := skillItem{path: claudeSkills + "/" + name, link: filepath.ToSlash(target)}
	fi, err := os.Lstat(s.path(it.path))
	if errors.Is(err, fs.ErrNotExist) {
		if wanted {
			it.drift = missing
		}
		return it, nil
	}
	if err != nil {
		return it, err
	}
	ours := false
	if fi.Mode()&fs.ModeSymlink != 0 {
		if it.seen, err = os.Readlink(s.path(it.path)); err != nil {
			return it, err
		}
		dir, err := os.Stat(skill)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			ours = filepath.ToSlash(it.seen) == it.link
		case err != nil:
			return it, err
		default:
			got, err := os.Stat(s.path(it.path))
			ours = err == nil && os.SameFile(dir, got)
		}
	}
	switch {
	case wanted && !ours:
		it.drift = foreign
	case !wanted && ours:
		it.drift = dropped
	}
	return it, nil
}

func skillName(rel string) string {
	name, _, _ := strings.Cut(rel, "/")
	return name
}

func hashOf(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func formatSum(files map[string][]byte) []byte {
	var b bytes.Buffer
	for _, rel := range slices.Sorted(maps.Keys(files)) {
		fmt.Fprintf(&b, "%s  %s\n", hashOf(files[rel]), rel)
	}
	return b.Bytes()
}

func parseSum(data []byte) (map[string]string, error) {
	out := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		hash, rel, ok := strings.Cut(sc.Text(), "  ")
		if _, err := hex.DecodeString(hash); !ok || err != nil || len(hash) != sha256.Size*2 || !fs.ValidPath(rel) || !strings.Contains(rel, "/") || path.Base(rel) == siteSkill {
			return nil, fmt.Errorf("line %d: want <sha256>  <skill>/<file>", n)
		}
		out[rel] = hash
	}
	return out, sc.Err()
}

func (ps *skillSync) findings() []finding {
	if !ps.installed {
		return []finding{skillsMissing}
	}
	var out []finding
	for _, it := range ps.items {
		if it.drift != current {
			out = append(out, it.finding())
		}
	}
	if len(out) == 0 && !bytes.Equal(ps.sum, ps.want) {
		out = append(out, finding{file: skillsSum, rule: ruleHgmxSkills, msg: "fingerprints out of date"})
	}
	return out
}

func (ps *skillSync) apply(s site) (wrote, removed int, err error) {
	for _, it := range ps.items {
		if it.drift != missing && it.drift != stale && it.drift != dropped {
			continue
		}
		if it.link != "" {
			if it.drift == missing {
				if err := os.MkdirAll(filepath.Dir(s.path(it.path)), 0o755); err != nil {
					return wrote, removed, err
				}
				if err := os.Symlink(filepath.FromSlash(it.link), s.path(it.path)); err != nil {
					return wrote, removed, err
				}
				wrote++
				continue
			}
			if got, err := os.Readlink(s.path(it.path)); err != nil || got != it.seen {
				return wrote, removed, changed{it}
			}
			if err := os.Remove(s.path(it.path)); err != nil {
				return wrote, removed, err
			}
			removed++
			continue
		}
		rel := strings.TrimPrefix(it.path, skillsDir+"/")
		seen := ""
		data, err := ps.root.ReadFile(rel)
		switch {
		case err == nil:
			seen = hashOf(data)
		case !errors.Is(err, fs.ErrNotExist):
			return wrote, removed, err
		}
		if seen != it.seen {
			it.drift = edited
			return wrote, removed, changed{it}
		}
		if it.drift == dropped {
			if err := ps.root.Remove(rel); err != nil {
				return wrote, removed, err
			}
			removed++
			continue
		}
		if err := writeRoot(ps.root, rel, it.data); err != nil {
			return wrote, removed, err
		}
		wrote++
	}
	for _, dir := range ps.prune {
		for ; dir != "."; dir = path.Dir(dir) {
			entries, err := fs.ReadDir(ps.root.FS(), dir)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return wrote, removed, err
			}
			if len(entries) > 0 {
				break
			}
			if err := ps.root.Remove(dir); err != nil {
				return wrote, removed, err
			}
		}
	}
	if ps.installed && bytes.Equal(ps.sum, ps.want) {
		return wrote, removed, nil
	}
	return wrote, removed, writeRoot(ps.root, path.Base(skillsSum), ps.want)
}

func writeRoot(root *os.Root, rel string, data []byte) (err error) {
	if err := root.MkdirAll(path.Dir(rel), 0o755); err != nil {
		return err
	}
	tmp := path.Join(path.Dir(rel), "."+path.Base(rel)+".hgmx-"+rand.Text())
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

func (a *audit) skills() error {
	if a.site.mod.Module.Mod.Path == hgmxPath {
		return nil
	}
	root, err := os.OpenRoot(a.site.path(skillsDir))
	if errors.Is(err, fs.ErrNotExist) {
		a.report(skillsMissing)
		return nil
	}
	if err != nil {
		return err
	}
	defer root.Close()
	ps, err := planSkills(a.site, root, skills.FS)
	if err != nil {
		return err
	}
	for _, f := range ps.findings() {
		a.report(f)
	}
	return nil
}

func skillsCmd() int {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod; run skills from the site's module root", slog.String("error", err.Error()))
		return 1
	}
	if s.mod.Module.Mod.Path == hgmxPath {
		log.Error("Refusing to install skills inside the hgmx module; run skills from a site's module root")
		return 1
	}
	root, err := openSkills(s)
	if err != nil {
		log.Error("Failed to open "+skillsDir, slog.String("error", err.Error()))
		return 1
	}
	defer root.Close()
	ps, err := planSkills(s, root, skills.FS)
	if err != nil {
		log.Error("Failed to read installed skills; no files were written", slog.String("error", err.Error()))
		return 1
	}
	if len(ps.findings()) == 0 {
		log.Info("Skills up to date", slog.String("dir", skillsDir))
		return 0
	}
	var blocked int
	for _, it := range ps.items {
		if it.blocks() {
			it.finding().print(os.Stdout)
			blocked++
		}
	}
	if blocked > 0 {
		log.Error("Skills not synced; no files were written", slog.Int("blocked", blocked))
		return 1
	}
	wrote, removed, err := ps.apply(s)
	var c changed
	if errors.As(err, &c) {
		c.it.finding().print(os.Stdout)
		log.Error("Skill sync stopped; resolve the file above, then rerun go tool hgmx skills", slog.Int("wrote", wrote), slog.Int("removed", removed))
		return 1
	}
	if err != nil {
		log.Error("Skill sync interrupted; rerun go tool hgmx skills", slog.String("error", err.Error()))
		return 1
	}
	log.Info("Skills synced", slog.String("dir", skillsDir), slog.Int("wrote", wrote), slog.Int("removed", removed))
	return 0
}
