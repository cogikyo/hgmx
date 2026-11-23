package main

import (
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

type skillItem struct {
	path  string
	skill string
	link  string
	seen  string
	drift drift
}

func fileItem(rel string, d drift) skillItem {
	return skillItem{path: skillsDir + "/" + rel, skill: skillName(rel), drift: d}
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
	plan      *copyPlan
	items     []skillItem
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
	recorded, installed, err := readLedger(root, path.Base(skillsSum), "<skill>/<file>", func(rel string) bool {
		return fs.ValidPath(rel) && strings.Contains(rel, "/") && path.Base(rel) != siteSkill
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w; delete it, then run go tool hgmx skills", skillsSum, err)
	}
	plan, err := planCopy(root, path.Base(skillsSum), recorded, want)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", skillsDir, err)
	}
	ps := &skillSync{plan: plan, installed: installed}
	for _, f := range plan.files {
		ps.items = append(ps.items, fileItem(f.rel, f.drift))
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
	if len(out) == 0 && !ps.plan.recorded() {
		out = append(out, finding{file: skillsSum, rule: ruleHgmxSkills, msg: "fingerprints out of date"})
	}
	return out
}

func (ps *skillSync) apply(s site) (wrote, removed int, err error) {
	for _, it := range ps.items {
		if it.link == "" {
			continue
		}
		switch it.drift {
		case missing:
			if err := os.MkdirAll(filepath.Dir(s.path(it.path)), 0o755); err != nil {
				return wrote, removed, err
			}
			if err := os.Symlink(filepath.FromSlash(it.link), s.path(it.path)); err != nil {
				return wrote, removed, err
			}
			wrote++
		case dropped:
			if got, err := os.Readlink(s.path(it.path)); err != nil || got != it.seen {
				return wrote, removed, fmt.Errorf("%s changed during the sync", it.path)
			}
			if err := os.Remove(s.path(it.path)); err != nil {
				return wrote, removed, err
			}
			removed++
		}
	}
	w, r, err := ps.plan.apply()
	return wrote + w, removed + r, err
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
		if it.drift.blocks() {
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
		fileItem(c.rel, edited).finding().print(os.Stdout)
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
