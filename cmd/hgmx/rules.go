package main

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/version"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/a-h/templ"
	parser "github.com/a-h/templ/parser/v2"
	"github.com/cogikyo/hgmx"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

type rule struct {
	name string
	why  string
	fix  string
}

var (
	ruleHgmxLatest = &rule{
		"hgmx-latest",
		"every site stays on the latest hgmx master while hgmx is being built out",
		"run go tool hgmx update",
	}
	ruleHgmxTool = &rule{
		"hgmx-tool",
		"the site's go.mod pins the hgmx library and CLI together so they cannot drift",
		"run go tool hgmx update; for a first install, run go get -tool github.com/cogikyo/hgmx/cmd/hgmx@master",
	}
	ruleHgmxSkills = &rule{
		"hgmx-skills",
		"every site carries the agent skills of its pinned hgmx version unedited, and keeps site-specific guidance in SITE.md beside each skill",
		"run go tool hgmx skills",
	}
	ruleJustfile = &rule{
		"layout/no-justfile",
		"hgmx commands replace per-site justfiles, so every site runs the same gen, dev, check and migrate",
		"delete it; run go tool hgmx gen, go tool hgmx dev, go tool hgmx check or go tool hgmx migrate instead",
	}
	ruleInternal = &rule{
		"layout/no-internal",
		"a site is a main module nobody imports, so internal/ only adds a path segment",
		"move its packages to routes/, middleware/, store/, views/ or a domain package at the root",
	}
	ruleVendor = &rule{
		"layout/no-vendor",
		"builds and deploys fetch modules from go.mod and go.sum, so a vendor/ copy only adds files that drift from them",
		"run git rm -r vendor/, then go mod tidy",
	}
	ruleRootGo = &rule{
		"layout/no-root-go",
		"the module root holds repo files such as go.mod, credentials.json, buf.yaml, proto/ and apps/; it is not a Go package",
		"move package main files to cmd/server/, and other code into the package that owns it",
	}
	ruleCmd = &rule{
		"layout/cmd",
		"dev, migrate, credentials and deploys build the site from cmd/server, and cmd/ names every command",
		"move the command into cmd/<name>/, or the library code out of cmd/",
	}
	ruleGoVersion = &rule{
		"layout/go-version",
		"sites stay on at least the Go release hgmx requires",
		"raise the go directive",
	}
	ruleTemplVersion = &rule{
		"templ/version",
		"hgmx runs the templ generator it was built with, so the site's templ runtime must match the code it generates",
		"match the site's templ version to hgmx's",
	}
	ruleDirection = &rule{
		"imports/direction",
		"dependencies point one way: routes wires middleware, store, domain packages and views; views only render",
		"move the code to the package that owns it, or pass the value in as a parameter",
	}
	ruleStoreOnly = &rule{
		"gorm/store-only",
		"store/ owns the database; other packages call Store methods and test errors with pg.ErrNotFound and the other sentinels",
		"add a Store method for the query and call it",
	}
	ruleLiveModels = &rule{
		"migrate/no-live-models",
		"a migration replays against the schema of its own time; live models change later and break fresh databases",
		"declare the columns the step needs inside the step, or write SQL",
	}
	ruleAutoMigrate = &rule{
		"gorm/no-automigrate",
		"schema changes run only as migrate steps during deploys, never at startup",
		"write a migrate.Step in store/migrations",
	}
	ruleGlobal = &rule{
		"gorm/no-global",
		"a package-level *gorm.DB hides who owns the connection and lets tests share state",
		"keep the *gorm.DB inside store.Store and pass the Store in",
	}
	ruleSlogContext = &rule{
		"slog/context",
		"the *Context variants carry the request's trace and span IDs into the log record",
		"call the *Context variant with the context in scope",
	}
	ruleExit = &rule{
		"exit/main-only",
		"os.Exit and log.Fatal skip deferred calls and graceful shutdown",
		"return the error and let main exit",
	}
	ruleInlineJS = &rule{
		"templ/no-inline-js",
		"hgmx's base CSP blocks inline handlers and script bodies, and hx-on evaluates code",
		"move the behaviour to hyperscript (_=) or a .js module served from static",
	}
	ruleCSP = &rule{
		"csp-literal",
		"hgmx sends a strict base Content-Security-Policy on every response, and a hand-written header drifts from it",
		"set overrides in server.Site.CSP, or for one response use server.CSP{…}.Header()",
	}
	ruleBoost = &rule{
		"htmx/boost",
		"hgmx's navigation script expects one hx-boost value, and copied literals drift from it",
		"write hx-boost:inherited={ scripts.Boost }",
	}
	ruleIcons = &rule{
		"hgmx/icons",
		"sites serve only the icons they copy into views/icons, and catalog names follow Lucide's renames",
		"use the replacement the finding names",
	}
	ruleIconsSource = &rule{
		"hgmx/icons-source",
		"one glyph handler serves one icon set, so a site serves either hgmx's whole catalog or its copied-in " + iconsDir + ", never both",
		"keep one: a preview site deletes " + iconsDir + " and uses catalog names; any other site drops the " + catalogPath + " import",
	}
	ruleIconsDrift = &rule{
		"hgmx/icons-drift",
		"the site's copied-in icons match the catalog of its pinned hgmx version",
		"run go tool hgmx icons update",
	}
	ruleIgnore = &rule{
		"hgmx/ignore",
		"every ignore names the rule it silences and why, and goes away with its finding",
		"write //hgmx:ignore <rule> <reason> on or above the finding's line, or delete the directive",
	}
)

var rules = []*rule{
	ruleHgmxLatest, ruleHgmxTool, ruleHgmxSkills,
	ruleJustfile, ruleInternal, ruleVendor, ruleRootGo, ruleCmd, ruleGoVersion, ruleTemplVersion,
	ruleDirection, ruleStoreOnly, ruleLiveModels,
	ruleAutoMigrate, ruleGlobal, ruleSlogContext, ruleExit,
	ruleCSP, ruleInlineJS, ruleBoost, ruleIcons, ruleIconsSource, ruleIconsDrift, ruleIgnore,
}

type finding struct {
	file      string
	line, col int
	rule      *rule
	msg       string
	fix       string
}

func (f finding) String() string {
	pos := f.file
	if f.line > 0 {
		pos = fmt.Sprintf("%s:%d:%d", f.file, f.line, f.col)
	}
	return fmt.Sprintf("%s [%s] %s", pos, f.rule.name, f.msg)
}

func (f finding) print(w io.Writer) {
	fmt.Fprintf(w, "%s\n  why: %s\n  fix: %s\n  see: CONVENTIONS.md#%s\n",
		f, f.rule.why, cmp.Or(f.fix, f.rule.fix), strings.ReplaceAll(f.rule.name, "/", "-"))
}

type ignore struct {
	file   string
	line   int
	col    int
	rule   string
	reason string
	bad    string
	used   bool
}

type audit struct {
	site     site
	findings []finding
	ignores  []*ignore
	templs   []string
	overlay  map[string][]byte
	maps     map[string]*parser.SourceMap
}

func (a *audit) report(f finding) {
	a.findings = append(a.findings, f)
}

func runRules(s site) ([]finding, error) {
	a := &audit{site: s, overlay: map[string][]byte{}, maps: map[string]*parser.SourceMap{}}
	if err := a.tree(); err != nil {
		return nil, err
	}
	if err := a.gomod(); err != nil {
		return nil, err
	}
	if err := a.hgmx(); err != nil {
		return nil, err
	}
	if err := a.skills(); err != nil {
		return nil, err
	}
	if err := a.templRules(); err != nil {
		return nil, err
	}
	if err := a.iconRules(); err != nil {
		return nil, err
	}
	if err := a.goRules(); err != nil {
		return a.resolve(false), err
	}
	return a.resolve(true), nil
}

func (a *audit) tree() error {
	if err := a.directives("go.mod"); err != nil {
		return err
	}
	if _, err := os.Stat(a.site.path("vendor")); err == nil {
		a.report(finding{file: "vendor/", rule: ruleVendor, msg: "vendor/ directory"})
	}
	dirs := map[string][]clause{}
	err := a.site.walk(func(rel string, d fs.DirEntry) error {
		name := d.Name()
		if !d.IsDir() {
			switch dir, ext := path.Dir(rel), path.Ext(name); {
			case dir == "." && ext == ".go":
				a.report(finding{file: rel, rule: ruleRootGo, msg: "Go file in the module root"})
			case dir == "." && ext == ".templ":
				a.report(finding{file: rel, rule: ruleRootGo, msg: "templ file in the module root"})
			case ext == ".go":
				if c, ok := packageClause(a.site.path(rel)); ok {
					dirs[dir] = append(dirs[dir], c)
				}
			}
		}
		switch {
		case d.IsDir():
			if name == "internal" {
				a.report(finding{file: rel + "/", rule: ruleInternal, msg: "internal/ directory"})
			}
		case strings.EqualFold(name, "justfile") || strings.EqualFold(name, ".justfile"):
			a.report(finding{file: rel, rule: ruleJustfile, msg: "justfile"})
		case strings.HasSuffix(name, ".templ"):
			a.templs = append(a.templs, rel)
			return a.directives(rel)
		case strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_templ.go"):
			return a.directives(rel)
		}
		return nil
	})
	if err != nil {
		return err
	}
	a.commands(dirs)
	return nil
}

func (a *audit) gomod() error {
	own, err := modfile.Parse("go.mod", hgmx.GoMod, nil)
	if err != nil {
		return err
	}
	f := a.site.mod
	if f.Go != nil && version.Compare("go"+f.Go.Version, "go"+own.Go.Version) < 0 {
		a.report(finding{
			file: "go.mod", line: f.Go.Syntax.Start.Line, col: f.Go.Syntax.Start.LineRune,
			rule: ruleGoVersion,
			msg:  fmt.Sprintf("go %s is older than hgmx's go %s", f.Go.Version, own.Go.Version),
			fix:  fmt.Sprintf("run go mod edit -go=%s, then go mod tidy", own.Go.Version),
		})
	}
	r := a.site.templ()
	if r == nil || r.Mod.Version == templ.Version() {
		return nil
	}
	fix := fmt.Sprintf("run go get %s@%s", templPath, templ.Version())
	if semver.Compare(r.Mod.Version, templ.Version()) > 0 {
		fix = fmt.Sprintf("run go tool hgmx update, or run go get %s@%s", templPath, templ.Version())
	}
	a.report(finding{
		file: "go.mod", line: r.Syntax.Start.Line, col: modCol(r.Syntax),
		rule: ruleTemplVersion,
		msg:  fmt.Sprintf("site requires templ %s; hgmx generates with %s", r.Mod.Version, templ.Version()),
		fix:  fix,
	})
	return nil
}

func modCol(l *modfile.Line) int {
	if l.InBlock {
		return l.Start.LineRune
	}
	return l.Start.LineRune + len("require ")
}

func (a *audit) directives(rel string) error {
	data, err := os.ReadFile(a.site.path(rel))
	if err != nil {
		return err
	}
	n := 0
	for line := range strings.Lines(string(data)) {
		n++
		line = strings.TrimRight(line, "\r\n")
		i := strings.Index(line, "//hgmx:")
		if i < 0 || strings.TrimLeft(line[:i], " \t") != "" && !strings.ContainsAny(line[i-1:i], " \t") {
			continue
		}
		d, ok := ast.ParseDirective(token.NoPos, line[i:])
		if !ok {
			continue
		}
		ig := &ignore{file: rel, line: n, col: i + 1}
		name, reason, _ := strings.Cut(d.Args, " ")
		ig.rule, ig.reason = name, strings.TrimSpace(reason)
		switch {
		case d.Name != "ignore":
			ig.bad = fmt.Sprintf("unknown directive //hgmx:%s", d.Name)
		case name == "":
			ig.bad = "ignore names no rule"
		case !slices.ContainsFunc(rules, func(r *rule) bool { return r.name == name }):
			ig.bad = fmt.Sprintf("ignore names unknown rule %s", name)
		case ig.reason == "":
			ig.bad = fmt.Sprintf("ignore of %s gives no reason", name)
		}
		a.ignores = append(a.ignores, ig)
	}
	return nil
}

func (a *audit) resolve(complete bool) []finding {
	seen := map[string]bool{}
	var out []finding
	for _, f := range a.findings {
		key := f.String()
		if seen[key] {
			continue
		}
		seen[key] = true
		if a.suppressed(f) {
			continue
		}
		out = append(out, f)
	}
	for _, ig := range a.ignores {
		switch {
		case ig.bad != "":
			out = append(out, finding{file: ig.file, line: ig.line, col: ig.col, rule: ruleIgnore, msg: ig.bad})
		// A failed package load can hide the finding an ignore suppresses.
		case !ig.used && complete:
			out = append(out, finding{file: ig.file, line: ig.line, col: ig.col, rule: ruleIgnore,
				msg: fmt.Sprintf("unused ignore: no %s finding on this or the next line", ig.rule)})
		}
	}
	slices.SortFunc(out, func(a, b finding) int {
		return cmp.Or(cmp.Compare(a.file, b.file), cmp.Compare(a.line, b.line), cmp.Compare(a.col, b.col), cmp.Compare(a.rule.name, b.rule.name))
	})
	return out
}

func (a *audit) suppressed(f finding) bool {
	hit := false
	for _, ig := range a.ignores {
		if ig.bad == "" && ig.file == f.file && ig.rule == f.rule.name && (ig.line == f.line || ig.line == f.line-1) {
			ig.used = true
			hit = true
		}
	}
	return hit
}
