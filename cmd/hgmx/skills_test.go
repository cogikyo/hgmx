package main

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestSkillSync(t *testing.T) {
	src := fstest.MapFS{
		"a/SKILL.md": {Data: []byte("new")},
		"a/SITE.md":  {Data: []byte("shipped")},
	}
	for _, tt := range []struct {
		name  string
		disk  map[string]string
		sum   map[string]string
		path  string
		drift drift
		dirs  []string
		links [][2]string
	}{
		{"no-op", map[string]string{"a/SKILL.md": "new"}, map[string]string{"a/SKILL.md": "new"}, ".agents/skills/a/SKILL.md", current, nil, nil},
		{"missing", nil, map[string]string{"a/SKILL.md": "new"}, ".agents/skills/a/SKILL.md", missing, nil, nil},
		{"stale", map[string]string{"a/SKILL.md": "old"}, map[string]string{"a/SKILL.md": "old"}, ".agents/skills/a/SKILL.md", stale, nil, nil},
		{"edited", map[string]string{"a/SKILL.md": "mine"}, map[string]string{"a/SKILL.md": "old"}, ".agents/skills/a/SKILL.md", edited, nil, nil},
		{"foreign", map[string]string{"a/SKILL.md": "mine"}, map[string]string{}, ".agents/skills/a/SKILL.md", foreign, nil, nil},
		{"dropped", map[string]string{"b/SKILL.md": "old", "b/SITE.md": "site"}, map[string]string{"b/SKILL.md": "old"}, ".agents/skills/b/SKILL.md", dropped, nil, nil},
		{"dropped-edited", map[string]string{"b/SKILL.md": "mine"}, map[string]string{"b/SKILL.md": "old"}, ".agents/skills/b/SKILL.md", edited, nil, nil},
		{"link", map[string]string{"a/SKILL.md": "new"}, map[string]string{"a/SKILL.md": "new"}, ".claude/skills/a", missing, []string{".claude"}, nil},
		{"linked-parent", map[string]string{"a/SKILL.md": "new"}, map[string]string{"a/SKILL.md": "new"}, ".claude/skills/a", foreign,
			[]string{".claude", "config/claude/skills"}, [][2]string{{".claude/skills", "../config/claude/skills"}, {".claude/skills/a", "../../.agents/skills/a"}}},
		{"escape", map[string]string{"b/SKILL.md": "old"}, map[string]string{"b/SKILL.md": "old"}, "", -1,
			[]string{".agents/skills", "outside"}, [][2]string{{".agents/skills/b", "../../outside"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := site{dir: t.TempDir()}
			write := func(rel, data string) {
				if err := os.MkdirAll(filepath.Dir(s.path(rel)), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(s.path(rel), []byte(data), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			for _, dir := range tt.dirs {
				if err := os.MkdirAll(s.path(dir), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			for _, l := range tt.links {
				if err := os.Symlink(l[1], s.path(l[0])); err != nil {
					t.Fatal(err)
				}
			}
			for rel, data := range tt.disk {
				write(skillsDir+"/"+rel, data)
			}
			sum := map[string][]byte{}
			for rel, data := range tt.sum {
				sum[rel] = []byte(data)
			}
			write(skillsSum, string(formatSum(sum)))
			root, err := openSkills(s)
			if err != nil {
				t.Fatal(err)
			}
			defer root.Close()
			ps, err := planSkills(s, root, src)
			if tt.drift == -1 {
				if err == nil {
					t.Fatal("plan followed a symlink out of " + skillsDir)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			got := current
			for _, it := range ps.items {
				if it.path == tt.path {
					got = it.drift
				}
			}
			if got != tt.drift {
				t.Fatalf("drift: %d; want %d", got, tt.drift)
			}
			if got == edited || got == foreign {
				return
			}
			if _, _, err := ps.apply(s); err != nil {
				t.Fatal(err)
			}
			if ps, err = planSkills(s, root, src); err != nil {
				t.Fatal(err)
			}
			if f := ps.findings(); len(f) != 0 {
				t.Fatalf("rerun findings: %v", f)
			}
			if _, err := os.Stat(s.path(skillsDir + "/a/SITE.md")); !os.IsNotExist(err) {
				t.Fatalf("shipped SITE.md written: %v", err)
			}
			if data, ok := tt.disk["b/SITE.md"]; ok {
				if b, err := os.ReadFile(s.path(skillsDir + "/b/SITE.md")); err != nil || string(b) != data {
					t.Fatalf("SITE.md changed: %q %v", b, err)
				}
			}
		})
	}
}
