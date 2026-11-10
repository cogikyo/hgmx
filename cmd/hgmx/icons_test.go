package main

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

const testSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
  <path d="%s"/>
</svg>
`

func testCatalog(t *testing.T, files map[string]string) iconCatalog {
	t.Helper()
	fsys := fstest.MapFS{
		"tags.json": {Data: []byte(`{"plus":["add"],"x":["close"],"circle-user":["person"]}`)},
		"meta.json": {Data: []byte(`{"categories":[],"icons":{"plus":{},"x":{},"circle-user":{"aliases":["user-circle"]}}}`)},
	}
	for rel, data := range files {
		fsys[rel] = &fstest.MapFile{Data: []byte(data)}
	}
	cat, err := readCatalog(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return cat
}

func svg(d string) string { return strings.Replace(testSVG, "%s", d, 1) }

func iconsRoot(t *testing.T) (*os.Root, string) {
	t.Helper()
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	return root, dir
}

func runIcons(t *testing.T, root *os.Root, cat iconCatalog, mode iconMode, add ...string) (iconResult, map[string]iconVerdict) {
	t.Helper()
	var out strings.Builder
	res, err := syncIcons(root, cat, mode, add, &out)
	if err != nil {
		t.Fatalf("sync: %v", err)
	}
	verdicts := map[string]iconVerdict{}
	for line := range strings.Lines(out.String()) {
		verb, rest, _ := strings.Cut(line, " ")
		name, _, _ := strings.Cut(strings.TrimSpace(rest), ":")
		for v, d := range iconVerdicts {
			if d.verb == verb {
				verdicts[name] = iconVerdict(v)
			}
		}
	}
	return res, verdicts
}

func read(t *testing.T, dir, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func write(t *testing.T, dir, rel, data string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, rel), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func idle(t *testing.T, root *os.Root, cat iconCatalog, mode iconMode) {
	t.Helper()
	res, _ := runIcons(t, root, cat, mode)
	if res.wrote+res.removed > 0 {
		t.Fatalf("rerun wrote %d and removed %d; want nothing", res.wrote, res.removed)
	}
	check, _ := runIcons(t, root, cat, iconMode{refresh: true, check: true})
	if check.drift || check.blocked > 0 {
		t.Fatalf("check after rerun: %+v", check)
	}
}

func TestIconsAdd(t *testing.T) {
	root, dir := iconsRoot(t)
	cat := testCatalog(t, map[string]string{
		"svg/plus.svg": svg("M5 12h14"), "anim/plus.js": "export default 1\n",
		"svg/x.svg": svg("M18 6 6 18"),
	})
	res, v := runIcons(t, root, cat, iconMode{}, "plus", "x")
	if res.blocked > 0 || v["plus"] != iconCopy || v["x"] != iconCopy {
		t.Fatalf("add: %+v %v", res, v)
	}
	if got := read(t, dir, iconsList); got != "[\n\t\"plus\",\n\t\"x\"\n]\n" {
		t.Errorf("icons.json = %q", got)
	}
	code := read(t, dir, iconsGo)
	for _, want := range []string{"//go:embed svg anim tags.json", `Plus Name = "plus"`, `X    Name = "x"`} {
		if !strings.Contains(code, want) {
			t.Errorf("icons.go lacks %q:\n%s", want, code)
		}
	}
	if sum := read(t, dir, iconsSum); strings.Count(sum, "\n") != 3 || !strings.Contains(sum, "  anim/plus.js\n") {
		t.Errorf("icons.sum = %q", sum)
	}
	idle(t, root, cat, iconMode{})

	var out strings.Builder
	for _, name := range []string{"nope", "user-circle", "../x"} {
		if _, err := syncIcons(root, cat, iconMode{}, []string{name}, &out); err == nil {
			t.Errorf("add %s accepted", name)
		}
	}
}

func TestIconsUpdate(t *testing.T) {
	root, dir := iconsRoot(t)
	old := map[string]string{
		"svg/plus.svg": svg("M5 12h14"), "anim/plus.js": "export default 1\n",
		"svg/x.svg": svg("M18 6 6 18"), "anim/x.js": "export default 2\n",
		"svg/circle-user.svg": svg("M12 2a10 10 0 1 0 0 20"),
	}
	runIcons(t, root, testCatalog(t, old), iconMode{}, "plus", "x", "circle-user")

	next := map[string]string{
		"svg/plus.svg": svg("M5 12h14M12 5v14"), "anim/plus.js": "export default 1\n",
		"svg/x.svg": svg("M18 6 6 18M6 6l12 12"), "anim/x.js": "export default 2\n",
		"svg/circle-user.svg": svg("M12 2a10 10 0 1 0 0 20"), "anim/circle-user.js": "export default 3\n",
	}
	cat := testCatalog(t, next)
	edited := svg("M18 6 6 18 Z")
	write(t, dir, "svg/x.svg", edited)

	_, v := runIcons(t, root, cat, iconMode{})
	if v["plus"] != iconStale || v["circle-user"] != iconStale || v["x"] != iconConflict {
		t.Fatalf("sync verdicts: %v", v)
	}
	if read(t, dir, "svg/plus.svg") != old["svg/plus.svg"] {
		t.Error("sync refreshed a stale icon")
	}

	res, v := runIcons(t, root, cat, iconMode{refresh: true})
	if v["plus"] != iconRefresh || v["circle-user"] != iconRefresh || v["x"] != iconConflict || res.blocked != 1 {
		t.Fatalf("update verdicts: %v %+v", v, res)
	}
	if read(t, dir, "svg/plus.svg") != next["svg/plus.svg"] || read(t, dir, "anim/circle-user.js") != next["anim/circle-user.js"] {
		t.Error("unedited icons not refreshed")
	}
	if read(t, dir, "svg/x.svg") != edited {
		t.Error("conflicting edit overwritten")
	}

	os.Remove(filepath.Join(dir, "svg/x.svg"))
	os.Remove(filepath.Join(dir, "anim/x.js"))
	if _, v := runIcons(t, root, cat, iconMode{refresh: true}); v["x"] != iconCopy {
		t.Fatalf("deleted conflict not restored: %v", v)
	}
	idle(t, root, cat, iconMode{refresh: true})

	mine := svg("M5 12h14M12 5v14 Z")
	write(t, dir, "svg/plus.svg", mine)
	if res, v := runIcons(t, root, cat, iconMode{refresh: true}); v["plus"] != iconEdited || res.blocked > 0 || res.wrote > 0 {
		t.Fatalf("edited: %v %+v", v, res)
	}
	if read(t, dir, "svg/plus.svg") != mine {
		t.Error("edited copy overwritten")
	}
	idle(t, root, cat, iconMode{refresh: true})
}

func TestIconsRemove(t *testing.T) {
	root, dir := iconsRoot(t)
	cat := testCatalog(t, map[string]string{
		"svg/plus.svg": svg("M5 12h14"), "anim/plus.js": "export default 1\n",
		"svg/x.svg": svg("M18 6 6 18"), "anim/x.js": "export default 2\n",
	})
	runIcons(t, root, cat, iconMode{}, "plus", "x")
	write(t, dir, iconsList, `["x"]`)
	res, v := runIcons(t, root, cat, iconMode{})
	if v["plus"] != iconRemove || res.removed != 2 {
		t.Fatalf("remove: %v %+v", v, res)
	}
	for _, rel := range []string{"svg/plus.svg", "anim/plus.js"} {
		if _, err := os.Stat(filepath.Join(dir, rel)); !os.IsNotExist(err) {
			t.Errorf("%s survived removal", rel)
		}
	}
	if code := read(t, dir, iconsGo); strings.Contains(code, "Plus") {
		t.Errorf("icons.go still names Plus:\n%s", code)
	}
	if strings.Contains(read(t, dir, iconsSum), "plus") || strings.Contains(read(t, dir, iconsTags), "plus") {
		t.Error("plus survived in icons.sum or tags.json")
	}
	idle(t, root, cat, iconMode{})

	write(t, dir, "svg/x.svg", svg("M0 0"))
	write(t, dir, iconsList, `[]`)
	if res, v := runIcons(t, root, cat, iconMode{}); v["x"] != iconUnlisted || res.blocked != 1 {
		t.Fatalf("edited unlisted: %v %+v", v, res)
	}
	if read(t, dir, "svg/x.svg") != svg("M0 0") {
		t.Error("edited unlisted copy deleted")
	}
	if code := read(t, dir, iconsGo); !strings.Contains(code, "//go:embed tags.json\n") {
		t.Errorf("empty set embeds icon dirs:\n%s", code)
	}
}

func TestIconsRename(t *testing.T) {
	root, dir := iconsRoot(t)
	old := testCatalog(t, map[string]string{"svg/plus.svg": svg("M5 12h14")})
	old.renamed = map[string]string{}
	old.fsys.(fstest.MapFS)["svg/user-circle.svg"] = &fstest.MapFile{Data: []byte(svg("M12 2"))}
	runIcons(t, root, old, iconMode{}, "user-circle")
	mine := svg("M12 2 Z")
	write(t, dir, "svg/user-circle.svg", mine)

	cat := testCatalog(t, map[string]string{"svg/circle-user.svg": svg("M12 2")})
	res, v := runIcons(t, root, cat, iconMode{refresh: true})
	if v["user-circle"] != iconRenamed || res.wrote > 0 || res.blocked > 0 {
		t.Fatalf("renamed without --rename: %v %+v", v, res)
	}
	if check, _ := runIcons(t, root, cat, iconMode{refresh: true, check: true}); !check.drift {
		t.Error("check passed with a pending rename")
	}

	_, v = runIcons(t, root, cat, iconMode{refresh: true, rename: true})
	if v["user-circle"] != iconMoved || v["circle-user"] != iconEdited {
		t.Fatalf("rename: %v", v)
	}
	if read(t, dir, "svg/circle-user.svg") != mine {
		t.Error("rename lost the local edit")
	}
	if _, err := os.Stat(filepath.Join(dir, "svg/user-circle.svg")); !os.IsNotExist(err) {
		t.Error("old name survived the rename")
	}
	var list []string
	for line := range strings.Lines(read(t, dir, iconsList)) {
		if s := strings.Trim(strings.TrimSpace(line), `",`); s != "[" && s != "]" {
			list = append(list, s)
		}
	}
	if !slices.Equal(list, []string{"circle-user"}) {
		t.Errorf("icons.json = %v", list)
	}
	idle(t, root, cat, iconMode{refresh: true})
}

func TestIconsCrash(t *testing.T) {
	root, dir := iconsRoot(t)
	cat := testCatalog(t, map[string]string{"svg/plus.svg": svg("M5 12h14"), "anim/plus.js": "export default 1\n"})
	runIcons(t, root, cat, iconMode{}, "plus")
	os.Remove(filepath.Join(dir, iconsSum))
	os.Remove(filepath.Join(dir, "anim/plus.js"))
	write(t, dir, "svg/.plus.svg.hgmx-abc", "partial")
	if _, v := runIcons(t, root, cat, iconMode{}); v["plus"] != iconCopy {
		t.Fatalf("interrupted copy: %v", v)
	}
	if _, err := os.Stat(filepath.Join(dir, "svg/.plus.svg.hgmx-abc")); !os.IsNotExist(err) {
		t.Error("temp file survived")
	}
	idle(t, root, cat, iconMode{})
}

func TestIconRulesBroken(t *testing.T) {
	s, err := loadSite("testdata/icons-broken")
	if err != nil {
		t.Fatal(err)
	}
	findings, err := runRules(s)
	if err == nil {
		t.Fatal("broken site loaded without error")
	}
	if len(findings) != 1 || !strings.Contains(findings[0].msg, "icons.BadgeExclamationPoint") {
		t.Fatalf("findings = %v", findings)
	}
}

func TestIconsKeep(t *testing.T) {
	root, dir := iconsRoot(t)
	runIcons(t, root, testCatalog(t, map[string]string{"svg/x.svg": svg("M18 6 6 18")}), iconMode{}, "x")
	mine := svg("M18 6 6 18 Z")
	write(t, dir, "svg/x.svg", mine)
	cat := testCatalog(t, map[string]string{"svg/x.svg": svg("M18 6 6 18M6 6l12 12"), "anim/x.js": "export default 2\n"})

	if res, v := runIcons(t, root, cat, iconMode{refresh: true}); v["x"] != iconConflict || res.blocked != 1 {
		t.Fatalf("conflict: %v %+v", v, res)
	}
	res, v := runIcons(t, root, cat, iconMode{refresh: true, keep: []string{"x"}})
	if v["x"] != iconKept || res.blocked > 0 {
		t.Fatalf("keep: %v %+v", v, res)
	}
	if read(t, dir, "svg/x.svg") != mine {
		t.Error("kept copy overwritten")
	}
	if _, err := os.Stat(filepath.Join(dir, "anim/x.js")); !os.IsNotExist(err) {
		t.Error("keep copied the catalog's module")
	}
	if sum := read(t, dir, iconsSum); !strings.Contains(sum, hashOf([]byte(svg("M18 6 6 18M6 6l12 12")))) || !strings.Contains(sum, "  anim/x.js\n") {
		t.Errorf("icons.sum does not record the catalog's version: %q", sum)
	}
	idle(t, root, cat, iconMode{refresh: true})

	var out strings.Builder
	for _, name := range []string{"plus", "nope"} {
		if _, err := syncIcons(root, cat, iconMode{refresh: true, keep: []string{name}}, nil, &out); err == nil {
			t.Errorf("--keep %s accepted", name)
		}
	}
}

func TestIconsRenameCrash(t *testing.T) {
	cat := testCatalog(t, map[string]string{"svg/circle-user.svg": svg("M12 2"), "anim/circle-user.js": "export default 3\n"})
	mine := svg("M12 2 Z")
	snapshot := func(dir string) map[string]string {
		out := map[string]string{}
		filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				rel, _ := filepath.Rel(dir, p)
				out[rel] = read(t, dir, rel)
			}
			return nil
		})
		return out
	}
	var want map[string]string
	for i, stage := range []string{"none", "anim-copied", "svg-copied", "old-removed", "list-written"} {
		t.Run(stage, func(t *testing.T) {
			root, dir := iconsRoot(t)
			for _, d := range []string{"svg", "anim"} {
				os.Mkdir(filepath.Join(dir, d), 0o755)
			}
			write(t, dir, iconsList, `["user-circle"]`)
			write(t, dir, "svg/user-circle.svg", mine)
			write(t, dir, "anim/user-circle.js", "export default 3\n")
			write(t, dir, iconsSum, hashOf([]byte("export default 3\n"))+"  anim/user-circle.js\n"+hashOf([]byte(svg("M12 2")))+"  svg/user-circle.svg\n")
			steps := []func(){
				func() { write(t, dir, "anim/circle-user.js", "export default 3\n") },
				func() { write(t, dir, "svg/circle-user.svg", mine) },
				func() {
					os.Remove(filepath.Join(dir, "anim/user-circle.js"))
					os.Remove(filepath.Join(dir, "svg/user-circle.svg"))
				},
				func() { write(t, dir, iconsList, `["circle-user"]`) },
			}
			for _, step := range steps[:i] {
				step()
			}
			res, v := runIcons(t, root, cat, iconMode{refresh: true, rename: true})
			if res.blocked > 0 || v["circle-user"] != iconEdited {
				t.Fatalf("resume: %v %+v", v, res)
			}
			idle(t, root, cat, iconMode{refresh: true, rename: true})
			got := snapshot(dir)
			if want == nil {
				want = got
				return
			}
			if !maps.Equal(got, want) {
				t.Errorf("state after resume differs from an uninterrupted rename:\n%v\nwant:\n%v", got, want)
			}
		})
	}
}
