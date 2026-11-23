package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/cogikyo/hgmx"
	"github.com/cogikyo/hgmx/proto"
	"github.com/cogikyo/hgmx/starters"
)

const (
	bufFile     = "buf.yaml"
	managedDir  = "proto/hgmx"
	kitSum      = managedDir + "/hgmx.sum"
	compatFile  = "apps/compatibility.json"
	androidLock = "apps/android/app/gradle.lockfile"
	iosLock     = "apps/ios/Starter.xcodeproj/project.xcworkspace/xcshareddata/swiftpm/Package.resolved"
)

var mobileRuntimes = map[string]string{
	"com.google.protobuf:protobuf-javalite": starterVersion("protobuf"),
	"com.connectrpc:connect-kotlin":         connectKotlin.version,
}

func starterVersion(key string) string {
	data, err := starters.FS.ReadFile("apps/android/gradle/libs.versions.toml")
	if err != nil {
		panic(err)
	}
	for line := range strings.Lines(string(data)) {
		k, v, ok := strings.Cut(line, "=")
		if ok && strings.TrimSpace(k) == key {
			return strings.Trim(strings.TrimSpace(v), `"`)
		}
	}
	panic("the Android starter's libs.versions.toml pins no " + key)
}

type platform struct {
	name   string
	marker string
	out    string
}

var platforms = []platform{
	{name: "go", out: "proto"},
	{name: "android", marker: "apps/android", out: "apps/android/app/src/main/generated"},
	{name: "ios", marker: "apps/ios", out: "apps/ios/API/Sources/API/Generated"},
}

type kit struct {
	site
	self      bool
	shipped   map[string][]byte
	platforms []platform
	skipped   []string
	apps      bool
	log       *slog.Logger
}

type change struct {
	path string
	data []byte
	why  string
}

func loadKit(s site, skip []string, log *slog.Logger) (*kit, error) {
	for _, name := range skip {
		if name != "android" && name != "ios" {
			return nil, fmt.Errorf("--skip %s: want android or ios", name)
		}
	}
	k := &kit{site: s, self: s.mod.Module.Mod.Path == hgmxPath, log: log}
	active := exists(s.path("proto"))
	for _, p := range platforms {
		if p.marker == "" {
			continue
		}
		if !exists(s.path(p.marker)) {
			continue
		}
		active, k.apps = true, true
		if slices.Contains(skip, p.name) {
			log.Warn("Skipped by request; not run", slog.String("platform", p.name))
			k.skipped = append(k.skipped, p.name)
			continue
		}
		k.platforms = append(k.platforms, p)
	}
	if !active {
		return nil, nil
	}
	k.platforms = append([]platform{platforms[0]}, k.platforms...)
	if k.self {
		return k, nil
	}
	k.shipped = map[string][]byte{bufFile: hgmx.BufYAML}
	err := fs.WalkDir(proto.FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		k.shipped["proto/"+p], err = proto.FS.ReadFile(p)
		return err
	})
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (k *kit) incomplete() bool {
	if k == nil || len(k.skipped) == 0 {
		return false
	}
	var stages []string
	for _, name := range k.skipped {
		stages = append(stages, name+" generation")
		stages = append(stages, name+" pins")
	}
	k.log.Error("Check incomplete; --skip never passes", slog.String("not run", strings.Join(stages, ", ")))
	return true
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

type inputs struct {
	plan    *copyPlan
	temps   []string
	refused []string
}

func (k *kit) inputs(root *os.Root) (*inputs, error) {
	if k.self {
		return &inputs{plan: &copyPlan{}}, nil
	}
	rec, _, err := readLedger(root, kitSum, bufFile+" or "+managedDir+"/<file>", func(rel string) bool {
		return rel == bufFile || fs.ValidPath(rel) && strings.HasPrefix(rel, managedDir+"/") && rel != kitSum
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w; delete it, then run go tool hgmx gen", kitSum, err)
	}
	plan, err := planCopy(root, kitSum, rec, k.shipped)
	if err != nil {
		return nil, err
	}
	in := &inputs{plan: plan}
	for _, f := range plan.files {
		_, ok := k.shipped[f.rel]
		switch {
		case !f.drift.blocks():
		case !ok:
			in.refused = append(in.refused, f.rel+": edited, and hgmx no longer ships it; move site changes under proto/<site>/, then delete it")
		default:
			in.refused = append(in.refused, f.rel+": differs from hgmx's copy; managed contracts are not edited in sites, so restore it with git checkout or delete it")
		}
	}
	if in.temps, err = fs.Glob(root.FS(), tempGlob(bufFile)); err != nil {
		return nil, err
	}
	err = fs.WalkDir(root.FS(), managedDir, func(rel string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil || d.IsDir() {
			return err
		}
		_, ok := k.shipped[rel]
		_, had := rec[rel]
		switch {
		case ok, had, rel == kitSum, goOutput(d.Name()):
		case temp(d.Name()):
			in.temps = append(in.temps, rel)
		default:
			in.refused = append(in.refused, rel+": not shipped by hgmx; proto/hgmx holds only hgmx's contracts, so move it under proto/<site>/")
		}
		return nil
	})
	return in, err
}

func (in *inputs) stale() []string {
	var out []string
	for _, f := range in.plan.files {
		switch f.drift {
		case missing:
			out = append(out, f.rel+": missing")
		case stale:
			out = append(out, f.rel+": out of date")
		case dropped:
			out = append(out, f.rel+": no longer shipped by hgmx")
		}
	}
	if len(out) == 0 && !in.plan.recorded() {
		out = append(out, kitSum+": out of date")
	}
	for _, rel := range in.temps {
		out = append(out, rel+": interrupted write")
	}
	return out
}

func (in *inputs) sync(root *os.Root) (wrote, removed int, err error) {
	if len(in.refused) > 0 {
		return 0, 0, errors.New(strings.Join(in.refused, "\n"))
	}
	for _, rel := range in.temps {
		if err := root.Remove(rel); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return 0, removed, err
		}
		removed++
	}
	wrote, n, err := in.plan.apply()
	return wrote, removed + n, err
}

func goOutput(name string) bool {
	return strings.HasSuffix(name, ".pb.go") || strings.HasSuffix(name, ".connect.go")
}

func (k *kit) generate(ctx context.Context) ([]change, error) {
	tmp, err := os.MkdirTemp("", "hgmx-gen-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	t := tools{log: k.log}
	var changes []change
	var errs []error
	for _, p := range k.platforms {
		out := filepath.Join(tmp, p.name)
		if err := k.run(ctx, t, p, out); err != nil {
			errs = append(errs, fmt.Errorf("%s generation: %w", p.name, err))
			continue
		}
		c, err := k.diff(p, out)
		if err != nil {
			return nil, err
		}
		changes = append(changes, c...)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return changes, nil
}

type plugin struct {
	Local         []string `json:"local,omitempty"`
	ProtocBuiltin string   `json:"protoc_builtin,omitempty"`
	ProtocPath    string   `json:"protoc_path,omitempty"`
	Out           string   `json:"out"`
	Opt           []string `json:"opt,omitempty"`
}

func (k *kit) run(ctx context.Context, t tools, p platform, out string) error {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	var plugins []plugin
	args := []string{"generate", "-o", out}
	managed := true
	switch p.name {
	case "go":
		managed = false
		if !k.self {
			ok, err := k.hasSources()
			if err != nil || !ok {
				return err
			}
			args = append(args, "--exclude-path", managedDir)
		}
		plugins = []plugin{
			{Local: []string{"go", "run", protocGenGo}, Out: ".", Opt: []string{"paths=source_relative"}},
			{Local: []string{"go", "run", connectGo}, Out: ".", Opt: []string{"paths=source_relative"}},
		}
	case "android":
		java, err := exec.LookPath("java")
		if err != nil {
			return errors.New("java is not on PATH; Kotlin Connect generation needs JDK 21")
		}
		protoc, err := t.protoc(ctx)
		if err != nil {
			return err
		}
		jar, err := t.connectKotlin(ctx)
		if err != nil {
			return err
		}
		plugins = []plugin{
			{ProtocBuiltin: "java", ProtocPath: protoc, Out: ".", Opt: []string{"lite"}},
			{Local: []string{java, "-jar", jar}, Out: "."},
		}
	case "ios":
		m := macHost()
		if err := t.swift(ctx, m); err != nil {
			return err
		}
		opt := []string{"Visibility=Public", "FileNaming=PathToUnderscores"}
		plugins = []plugin{
			{Local: m.plugin(swiftTools[0]), Out: ".", Opt: opt},
			{Local: m.plugin(swiftTools[1]), Out: ".", Opt: append(opt, "GenerateAsyncMethods=true", "GenerateCallbackMethods=false")},
		}
	}
	template := map[string]any{"version": "v2", "plugins": plugins}
	if managed {
		template["managed"] = map[string]any{
			"enabled": true,
			"disable": []map[string]string{{"file_option": "java_package"}, {"file_option": "go_package"}},
		}
	}
	data, err := json.Marshal(template)
	if err != nil {
		return err
	}
	cmd := bufCmd(ctx, k.dir, append(args, "--template", string(data))...)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("buf generate: %w\n%s", err, tail(string(msg), 30))
	}
	return nil
}

func (k *kit) hasSources() (bool, error) {
	managed := k.path(managedDir)
	found := false
	err := filepath.WalkDir(k.path("proto"), func(name string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case d.IsDir() && name == managed:
			return filepath.SkipDir
		case !d.IsDir() && strings.HasSuffix(name, ".proto"):
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return found, err
}

func (k *kit) diff(p platform, out string) ([]change, error) {
	want := map[string][]byte{}
	err := filepath.WalkDir(out, func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(out, name)
		if err != nil {
			return err
		}
		want[path.Join(p.out, filepath.ToSlash(rel))], err = os.ReadFile(name)
		return err
	})
	if err != nil {
		return nil, err
	}
	var changes []change
	err = filepath.WalkDir(k.path(p.out), func(name string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil || d.IsDir() {
			return err
		}
		if p.name == "go" && !goOutput(d.Name()) && !temp(d.Name()) {
			return nil
		}
		rel, err := filepath.Rel(k.dir, name)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		data, ok := want[rel]
		if !ok {
			why := "no source"
			if temp(d.Name()) {
				why = "interrupted write"
			}
			changes = append(changes, change{rel, nil, why})
			return nil
		}
		delete(want, rel)
		got, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, data) {
			changes = append(changes, change{rel, data, "out of date"})
		}
		return nil
	})
	for _, rel := range slices.Sorted(maps.Keys(want)) {
		changes = append(changes, change{rel, want[rel], "missing"})
	}
	return changes, err
}

func (k *kit) apply(root *os.Root, changes []change) error {
	for _, c := range changes {
		if c.data != nil {
			if err := writeRoot(root, c.path, c.data); err != nil {
				return err
			}
			k.log.Debug("Generated", slog.String("file", c.path), slog.String("reason", c.why))
			continue
		}
		if err := root.Remove(c.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		k.log.Debug("Removed", slog.String("file", c.path), slog.String("reason", c.why))
		for dir := path.Dir(c.path); dir != "." && dir != "proto"; dir = path.Dir(dir) {
			if root.Remove(dir) != nil {
				break
			}
		}
	}
	return nil
}

func (k *kit) gen(ctx context.Context) error {
	root, err := os.OpenRoot(k.dir)
	if err != nil {
		return err
	}
	defer root.Close()
	in, err := k.inputs(root)
	if err != nil {
		return err
	}
	wrote, removed, err := in.sync(root)
	if err != nil {
		return err
	}
	generated, err := k.generate(ctx)
	if err != nil {
		return fmt.Errorf("%w\nno generated files were replaced", err)
	}
	if err := k.apply(root, generated); err != nil {
		return err
	}
	for _, c := range generated {
		if c.data == nil {
			removed++
		} else {
			wrote++
		}
	}
	if wrote+removed == 0 {
		k.log.Info("Contracts up to date")
		return nil
	}
	k.log.Info("Contracts generated", slog.Int("written", wrote), slog.Int("removed", removed))
	return nil
}

func (k *kit) stale(ctx context.Context) ([]string, error) {
	root, err := os.OpenRoot(k.dir)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	in, err := k.inputs(root)
	if err != nil {
		return nil, err
	}
	generated, err := k.generate(ctx)
	if err != nil {
		return nil, err
	}
	stale := in.stale()
	for _, c := range generated {
		stale = append(stale, c.path+": "+c.why)
	}
	return stale, nil
}

func (k *kit) contracts(ctx context.Context) error {
	root, err := os.OpenRoot(k.dir)
	if err != nil {
		return err
	}
	defer root.Close()
	in, err := k.inputs(root)
	if err != nil {
		return err
	}
	for _, line := range in.refused {
		fmt.Fprintln(os.Stdout, line)
	}
	if len(in.refused) > 0 {
		return fmt.Errorf("edited or unowned managed files: %d", len(in.refused))
	}
	for _, args := range [][]string{{"format", "--diff", "--exit-code"}, {"lint"}} {
		cmd := bufCmd(ctx, k.dir, args...)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("buf %s: %w", args[0], err)
		}
	}
	return k.breaking(ctx)
}

type compatibility struct {
	Base     string    `json:"base"`
	Releases []release `json:"releases"`
}

type release struct {
	Platform string `json:"platform"`
	Build    int    `json:"build"`
	Commit   string `json:"commit"`
	Majors   []int  `json:"majors"`
}

func (k *kit) compatibility() (compatibility, error) {
	var c compatibility
	data, err := os.ReadFile(k.path(compatFile))
	if errors.Is(err, fs.ErrNotExist) {
		return c, fmt.Errorf("%s is missing; native apps need it, with \"releases\": [] until the first store release", compatFile)
	}
	if err != nil {
		return c, err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&c); err != nil {
		return c, fmt.Errorf("%s: %w", compatFile, err)
	}
	if c.Base == "" || c.Releases == nil {
		return c, fmt.Errorf("%s: want a base ref and a releases list", compatFile)
	}
	for i, r := range c.Releases {
		if !slices.ContainsFunc(platforms[1:], func(p platform) bool { return p.name == r.Platform }) ||
			r.Build < 1 || len(r.Commit) != 40 || len(r.Majors) == 0 {
			return c, fmt.Errorf("%s: release %d wants platform android or ios, a positive build, a full commit hash and supported majors", compatFile, i)
		}
	}
	return c, nil
}

func (k *kit) breaking(ctx context.Context) error {
	if !k.apps {
		return nil
	}
	c, err := k.compatibility()
	if err != nil {
		return err
	}
	if len(c.Releases) == 0 {
		k.log.Info("No store releases yet; no release baselines to check", slog.String("file", compatFile))
		return nil
	}
	retained := fmt.Errorf("%s records released builds, but hgmx has no retained-client checks yet; they must exist before a release is recorded", compatFile)
	type baseline struct{ name, commit string }
	var baselines []baseline
	for _, r := range c.Releases {
		if exec.CommandContext(ctx, "git", "-C", k.dir, "cat-file", "-e", r.Commit+"^{commit}").Run() != nil {
			return errors.Join(fmt.Errorf("release baseline %s build %d (%s) is not in this clone; fetch full history", r.Platform, r.Build, r.Commit), retained)
		}
		baselines = append(baselines, baseline{fmt.Sprintf("%s build %d", r.Platform, r.Build), r.Commit})
	}
	out, err := exec.CommandContext(ctx, "git", "-C", k.dir, "merge-base", "HEAD", c.Base).Output()
	if err != nil {
		return errors.Join(fmt.Errorf("merge-base with %s: %w; fetch %s", c.Base, err, c.Base), retained)
	}
	baselines = append(baselines, baseline{"merge-base with " + c.Base, strings.TrimSpace(string(out))})
	var failed []string
	for _, b := range baselines {
		paths, err := exec.CommandContext(ctx, "git", "-C", k.dir, "ls-tree", "--name-only", b.commit, "--", bufFile).Output()
		if err != nil {
			return fmt.Errorf("git ls-tree %s: %w", b.commit, err)
		}
		if len(bytes.TrimSpace(paths)) == 0 {
			k.log.Info("No contract at baseline; nothing to break", slog.String("baseline", b.name))
			continue
		}
		cmd := bufCmd(ctx, k.dir, "breaking", "--against", ".git#ref="+b.commit)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if cmd.Run() != nil {
			failed = append(failed, b.name)
		}
	}
	if len(failed) > 0 {
		return errors.Join(fmt.Errorf("breaking changes against %s", strings.Join(failed, ", ")), retained)
	}
	return retained
}

func (k *kit) pins() []string {
	var bad []string
	for _, pin := range []struct{ path, min string }{{protobufGoPath, protocGenGoRT}, {connectGoPath, connectGoRT}} {
		if v := k.version(pin.path); v != "" && semver.Compare(v, pin.min) < 0 {
			bad = append(bad, fmt.Sprintf("go.mod: %s %s is older than the generator's %s", pin.path, v, pin.min))
		}
	}
	for _, p := range k.platforms {
		switch p.name {
		case "android":
			bad = append(bad, k.androidPins()...)
		case "ios":
			bad = append(bad, k.iosPins()...)
		}
	}
	return bad
}

func (k *kit) androidPins() []string {
	data, err := os.ReadFile(k.path(androidLock))
	if err != nil {
		return []string{androidLock + ": missing; run ./gradlew :app:dependencies --write-locks in apps/android"}
	}
	var bad []string
	seen := map[string]bool{}
	for line := range strings.Lines(string(data)) {
		coord, _, _ := strings.Cut(line, "=")
		parts := strings.Split(coord, ":")
		if len(parts) != 3 {
			continue
		}
		module, v := parts[0]+":"+parts[1], parts[2]
		for prefix, want := range mobileRuntimes {
			if module == prefix || strings.HasPrefix(module, prefix+"-") {
				seen[prefix] = true
				if v != want {
					bad = append(bad, fmt.Sprintf("%s: %s %s; generated code needs %s", androidLock, module, v, want))
				}
			}
		}
	}
	for _, prefix := range slices.Sorted(maps.Keys(mobileRuntimes)) {
		if !seen[prefix] {
			bad = append(bad, fmt.Sprintf("%s: no %s; generated code needs %s", androidLock, prefix, mobileRuntimes[prefix]))
		}
	}
	return bad
}

func (k *kit) iosPins() []string {
	data, err := os.ReadFile(k.path(iosLock))
	if err != nil {
		return []string{iosLock + ": missing; resolve the Swift packages in apps/ios with xcodebuild -resolvePackageDependencies"}
	}
	var resolved struct {
		Pins []struct {
			Identity string
			State    struct{ Revision, Version string }
		}
	}
	if err := json.Unmarshal(data, &resolved); err != nil {
		return []string{fmt.Sprintf("%s: %v", iosLock, err)}
	}
	got := map[string]struct{ Revision, Version string }{}
	for _, p := range resolved.Pins {
		got[p.Identity] = p.State
	}
	var bad []string
	for _, t := range swiftTools {
		pin, ok := got[t.name]
		switch {
		case !ok:
			bad = append(bad, fmt.Sprintf("%s: no %s; generated code needs %s", iosLock, t.name, t.version))
		case pin.Version != t.version || pin.Revision != t.commit:
			bad = append(bad, fmt.Sprintf("%s: %s %s; generated code needs %s at %s", iosLock, t.name, pin.Version, t.version, t.commit))
		}
	}
	return bad
}
