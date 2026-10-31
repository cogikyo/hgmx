package main

import (
	"bufio"
	"bytes"
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/version"
	"io"
	"io/fs"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"golang.org/x/mod/modfile"
)

type rule struct {
	name string
	why  string
	fix  string
}

var (
	ruleJustfile = &rule{
		"layout/no-justfile",
		"hgmx commands replace per-site justfiles, so every site runs the same gen, dev, check and migrate",
		"delete it; run hgmx gen, hgmx dev, hgmx check or hgmx migrate instead",
	}
	ruleInternal = &rule{
		"layout/no-internal",
		"a site is a main module nobody imports, so internal/ only adds a path segment",
		"move its packages to routes/, middleware/, store/, views/ or a domain package at the root",
	}
	ruleGoVersion = &rule{
		"layout/go-version",
		"sites stay on the Go release hgmx is built with",
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
		"the CSP allows script-src 'self' only, so inline handlers and script bodies are blocked, and hx-on evaluates code",
		"move the behaviour to hyperscript (_=) or a .js module served from static",
	}
	ruleIgnore = &rule{
		"hgmx/ignore",
		"every ignore names the rule it silences and why, and goes away with its finding",
		"write //hgmx:ignore <rule> <reason> on or above the finding's line, or delete the directive",
	}
)

var rules = []*rule{
	ruleJustfile, ruleInternal, ruleGoVersion, ruleTemplVersion,
	ruleDirection, ruleStoreOnly, ruleLiveModels,
	ruleAutoMigrate, ruleGlobal, ruleSlogContext, ruleExit,
	ruleInlineJS, ruleIgnore,
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
}

func (a *audit) report(f finding) {
	a.findings = append(a.findings, f)
}

func runRules(s site) ([]finding, error) {
	a := &audit{site: s}
	if err := a.tree(); err != nil {
		return nil, err
	}
	a.gomod()
	if err := a.templRules(); err != nil {
		return nil, err
	}
	if err := a.goRules(); err != nil {
		return nil, err
	}
	return a.resolve(), nil
}

func (a *audit) tree() error {
	if err := a.directives("go.mod"); err != nil {
		return err
	}
	return a.site.walk(func(rel string, d fs.DirEntry) error {
		name := d.Name()
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
}

func (a *audit) gomod() {
	f := a.site.mod
	if lang := version.Lang(runtime.Version()); lang != "" && f.Go != nil && version.Compare("go"+f.Go.Version, lang) < 0 {
		want := strings.TrimPrefix(lang, "go")
		a.report(finding{
			file: "go.mod", line: f.Go.Syntax.Start.Line, col: f.Go.Syntax.Start.LineRune,
			rule: ruleGoVersion,
			msg:  fmt.Sprintf("go %s is older than go %s, which hgmx is built with", f.Go.Version, want),
			fix:  fmt.Sprintf("run go mod edit -go=%s, then go mod tidy", want),
		})
	}
	r := a.site.templ()
	if r == nil || r.Mod.Version == templ.Version() {
		return
	}
	fix := fmt.Sprintf("run go get %s@%s", templPath, templ.Version())
	if version.Compare(goVersion(r.Mod.Version), goVersion(templ.Version())) > 0 {
		fix = fmt.Sprintf("install an hgmx built with templ %s, or run go get %s@%s", r.Mod.Version, templPath, templ.Version())
	}
	a.report(finding{
		file: "go.mod", line: r.Syntax.Start.Line, col: modCol(r.Syntax),
		rule: ruleTemplVersion,
		msg:  fmt.Sprintf("site requires templ %s; hgmx generates with %s", r.Mod.Version, templ.Version()),
		fix:  fix,
	})
}

func goVersion(semver string) string {
	return "go" + strings.TrimPrefix(semver, "v")
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
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(nil, len(data)+1)
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
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
	return sc.Err()
}

func (a *audit) resolve() []finding {
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
		case !ig.used:
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
