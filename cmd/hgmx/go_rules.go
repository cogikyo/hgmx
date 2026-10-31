package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

const gormPath = "gorm.io/gorm"

type layer string

const (
	layerMain       layer = "main"
	layerRoutes     layer = "routes"
	layerMiddleware layer = "middleware"
	layerStore      layer = "store"
	layerMigrations layer = "store/migrations"
	layerViews      layer = "views"
	layerPages      layer = "views/pages"
)

var allowed = map[layer]struct {
	to   []layer
	rule string
}{
	layerRoutes:     {[]layer{layerRoutes, layerMiddleware, layerStore, layerViews, layerPages}, "routes may import middleware, store, views and domain packages"},
	layerMiddleware: {[]layer{layerMiddleware, layerStore}, "middleware may import store and domain packages"},
	layerStore:      {[]layer{layerStore}, "store may import only store/..."},
	layerMigrations: {[]layer{layerMigrations}, "store/migrations may import only store/migrations/..."},
	layerViews:      {[]layer{layerViews}, "views and its components may import only views/... (not views/pages) and hgmx"},
	layerPages:      {[]layer{layerViews, layerPages}, "views/pages may import only views/... and hgmx"},
}

type sitePkg struct {
	pkg   *packages.Package
	rel   string
	layer layer
}

func (a *audit) layerOf(p *packages.Package) (string, layer) {
	mod := a.site.mod.Module.Mod.Path
	if p.Name == "main" {
		return strings.TrimPrefix(strings.TrimPrefix(p.PkgPath, mod), "/"), layerMain
	}
	rel := strings.TrimPrefix(p.PkgPath, mod+"/")
	switch {
	case under(rel, "routes"):
		return rel, layerRoutes
	case under(rel, "middleware"):
		return rel, layerMiddleware
	case under(rel, "store/migrations"):
		return rel, layerMigrations
	case under(rel, "store"):
		return rel, layerStore
	case under(rel, "views/pages"):
		return rel, layerPages
	case under(rel, "views"):
		return rel, layerViews
	}
	return rel, layer("domain " + strings.SplitN(rel, "/", 2)[0])
}

func under(rel, dir string) bool {
	return rel == dir || strings.HasPrefix(rel, dir+"/")
}

func (l layer) domain() bool {
	return strings.HasPrefix(string(l), "domain ")
}

func (a *audit) goRules() error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes |
			packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:        a.site.dir,
		Env:        a.site.env(),
		BuildFlags: a.site.flags(),
		Tests:      true,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return err
	}
	var errs []error
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		for _, e := range p.Errors {
			errs = append(errs, errors.New(e.Error()))
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("load packages: %w", errors.Join(errs...))
	}

	mod := a.site.mod.Module.Mod.Path
	site := map[string]sitePkg{}
	var roots []*packages.Package
	for _, p := range pkgs {
		if strings.HasSuffix(p.ID, ".test") {
			continue
		}
		roots = append(roots, p)
		if p.ID == p.PkgPath && (p.PkgPath == mod || strings.HasPrefix(p.PkgPath, mod+"/")) {
			rel, l := a.layerOf(p)
			site[p.PkgPath] = sitePkg{p, rel, l}
		}
	}
	for _, sp := range site {
		a.imports(sp, site)
	}

	graph, err := checker.Analyze([]*analysis.Analyzer{siteAnalyzer}, roots, nil)
	if err != nil {
		return err
	}
	for act := range graph.All() {
		if !act.IsRoot {
			continue
		}
		if act.Err != nil {
			return act.Err
		}
		for _, d := range act.Diagnostics {
			r := ruleByName(d.Category)
			pos := act.Package.Fset.Position(d.Pos)
			fix := ""
			if len(d.SuggestedFixes) > 0 {
				fix = d.SuggestedFixes[0].Message
			}
			a.goReport(pos, r, d.Message, fix)
		}
	}
	return nil
}

func ruleByName(name string) *rule {
	i := slices.IndexFunc(rules, func(r *rule) bool { return r.name == name })
	return rules[i]
}

func (a *audit) goReport(pos token.Position, r *rule, msg, fix string) {
	rel, err := filepath.Rel(a.site.dir, pos.Filename)
	if err != nil || strings.HasPrefix(rel, "..") {
		return
	}
	a.report(finding{file: filepath.ToSlash(rel), line: pos.Line, col: pos.Column, rule: r, msg: msg, fix: fix})
}

func (a *audit) imports(sp sitePkg, site map[string]sitePkg) {
	p := sp.pkg
	for _, f := range p.Syntax {
		for _, spec := range f.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			pos := p.Fset.Position(spec.Pos())
			to, internal := site[path]
			switch {
			case path == gormPath || strings.HasPrefix(path, "gorm.io/"):
				if sp.layer != layerStore && sp.layer != layerMigrations {
					a.goReport(pos, ruleStoreOnly, fmt.Sprintf("%s imports %s; gorm stays inside store/", sp.rel, path), "")
				}
			case path == "net/http":
				if sp.layer == layerViews || sp.layer == layerPages {
					a.goReport(pos, ruleDirection, fmt.Sprintf("%s imports net/http; views render, and handlers live in routes", sp.rel), "")
				}
			case internal && sp.layer == layerMigrations:
				if chain := storeChain(to, site, map[string]bool{}); chain != nil {
					msg := fmt.Sprintf("%s imports %s", sp.rel, strings.Join(chain, ", which imports "))
					a.goReport(pos, ruleLiveModels, msg, "")
				} else if to.layer != layerMigrations {
					a.goReport(pos, ruleDirection, fmt.Sprintf("%s imports %s; %s", sp.rel, to.rel, allowed[sp.layer].rule), "")
				}
			case internal && !edgeAllowed(sp.layer, to.layer):
				desc := allowed[sp.layer].rule
				if sp.layer.domain() {
					desc = "a domain package may import only itself and store"
				}
				a.goReport(pos, ruleDirection, fmt.Sprintf("%s imports %s; %s", sp.rel, to.rel, desc), "")
			}
		}
	}
}

func edgeAllowed(from, to layer) bool {
	switch {
	case from == layerMain:
		return true
	case from.domain():
		return to == from || to == layerStore
	case to.domain():
		return from == layerRoutes || from == layerMiddleware
	}
	return slices.Contains(allowed[from].to, to)
}

func storeChain(p sitePkg, site map[string]sitePkg, seen map[string]bool) []string {
	if p.layer == layerStore {
		return []string{p.rel}
	}
	seen[p.pkg.PkgPath] = true
	for path := range p.pkg.Imports {
		next, ok := site[path]
		if !ok || seen[path] {
			continue
		}
		if chain := storeChain(next, site, seen); chain != nil {
			return append([]string{p.rel}, chain...)
		}
	}
	return nil
}

var siteAnalyzer = &analysis.Analyzer{
	Name: "hgmx",
	Doc:  "reports hgmx site rules that need type information",
	Run:  analyze,
}

func analyze(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		test := strings.HasSuffix(pass.Fset.File(f.Pos()).Name(), "_test.go")
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if !test && d.Tok == token.VAR {
					globals(pass, d)
				}
			case *ast.FuncDecl:
				if d.Body == nil {
					continue
				}
				exitOK := d.Recv == nil && (d.Name.Name == "main" && pass.Pkg.Name() == "main" || d.Name.Name == "TestMain" && test)
				ast.Inspect(d.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						calls(pass, call, exitOK)
					}
					return true
				})
			}
		}
	}
	return nil, nil
}

func globals(pass *analysis.Pass, d *ast.GenDecl) {
	for _, spec := range d.Specs {
		for _, name := range spec.(*ast.ValueSpec).Names {
			obj := pass.TypesInfo.Defs[name]
			if obj != nil && isNamed(deref(obj.Type()), gormPath, "DB") {
				pass.Report(analysis.Diagnostic{Pos: name.Pos(), Category: ruleGlobal.name,
					Message: fmt.Sprintf("package-level %s holds a %s", name.Name, types.TypeString(obj.Type(), (*types.Package).Name))})
			}
		}
	}
}

func calls(pass *analysis.Pass, call *ast.CallExpr, exitOK bool) {
	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return
	}
	name := fn.Name()
	switch fn.Pkg().Path() {
	case gormPath:
		if name == "AutoMigrate" {
			pass.Report(analysis.Diagnostic{Pos: call.Pos(), Category: ruleAutoMigrate.name, Message: "AutoMigrate changes the schema outside a migration"})
		}
	case "os":
		if name == "Exit" && !exitOK {
			pass.Report(analysis.Diagnostic{Pos: call.Pos(), Category: ruleExit.name, Message: "os.Exit outside main"})
		}
	case "log":
		if strings.HasPrefix(name, "Fatal") && !exitOK {
			pass.Report(analysis.Diagnostic{Pos: call.Pos(), Category: ruleExit.name, Message: fmt.Sprintf("%s outside main", callName(call))})
		}
	case "log/slog":
		switch name {
		case "Debug", "Info", "Warn", "Error":
		default:
			return
		}
		ctx := contextInScope(pass, call.Pos())
		if ctx == "" {
			return
		}
		pass.Report(analysis.Diagnostic{
			Pos:      call.Pos(),
			Category: ruleSlogContext.name,
			Message:  fmt.Sprintf("%s drops the context; %s is in scope", callName(call), strings.TrimSuffix(ctx, ".Context()")),
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: fmt.Sprintf("call %sContext(%s, ...)", callName(call), ctx),
			}},
		})
	}
}

func callName(call *ast.CallExpr) string {
	return types.ExprString(call.Fun)
}

func contextInScope(pass *analysis.Pass, pos token.Pos) string {
	pkg := pass.Pkg.Scope()
	for s := pkg.Innermost(pos); s != nil && s != pkg; s = s.Parent() {
		for _, name := range s.Names() {
			v, ok := s.Lookup(name).(*types.Var)
			if !ok || v.Pos() >= pos {
				continue
			}
			switch t := v.Type(); {
			case isNamed(t, "context", "Context"):
				return name
			case isNamed(deref(t), "net/http", "Request") && t != deref(t):
				return name + ".Context()"
			}
		}
	}
	return ""
}

func deref(t types.Type) types.Type {
	if p, ok := types.Unalias(t).(*types.Pointer); ok {
		return p.Elem()
	}
	return t
}

func isNamed(t types.Type, pkg, name string) bool {
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == pkg && obj.Name() == name
}
