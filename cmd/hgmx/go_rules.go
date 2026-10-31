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

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	gormPath = "gorm.io/gorm"
	domain   = "domain packages"
)

var layers = []struct {
	dir string
	may []string
}{
	{"routes", []string{"routes", "middleware", "store", "views", "views/pages", domain}},
	{"middleware", []string{"middleware", "store", domain}},
	{"store/migrations", []string{"store/migrations"}},
	{"store", []string{"store"}},
	{"views/pages", []string{"views", "views/pages"}},
	{"views", []string{"views"}},
}

type sitePkg struct {
	pkg    *packages.Package
	rel    string
	layer  string
	may    []string
	domain bool
}

func (a *audit) classify(p *packages.Package) sitePkg {
	mod := a.site.mod.Module.Mod.Path
	rel := strings.TrimPrefix(strings.TrimPrefix(p.PkgPath, mod), "/")
	sp := sitePkg{pkg: p, rel: rel}
	if p.Name == "main" {
		sp.layer = "main"
		return sp
	}
	for _, l := range layers {
		if rel == l.dir || strings.HasPrefix(rel, l.dir+"/") {
			sp.layer, sp.may = l.dir, l.may
			return sp
		}
	}
	sp.layer, _, _ = strings.Cut(rel, "/")
	sp.may = []string{sp.layer, "store"}
	sp.domain = true
	return sp
}

func (sp sitePkg) allows(to sitePkg) bool {
	return sp.layer == "main" || slices.Contains(sp.may, to.layer) || to.domain && slices.Contains(sp.may, domain)
}

func (sp sitePkg) policy() string {
	may := strings.Join(sp.may, ", ")
	if i := strings.LastIndex(may, ", "); i >= 0 {
		may = may[:i] + " and " + may[i+2:]
	}
	return fmt.Sprintf("%s may import site packages only from %s", sp.layer, may)
}

func (sp sitePkg) String() string {
	if sp.rel == sp.layer || sp.domain {
		return sp.rel
	}
	return fmt.Sprintf("%s (%s)", sp.rel, sp.layer)
}

func (a *audit) goRules() error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:   a.site.dir,
		Env:   a.site.env(),
		Tests: true,
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
	for _, p := range pkgs {
		if strings.HasSuffix(p.ID, ".test") {
			continue
		}
		if err := a.typeRules(p); err != nil {
			return err
		}
		if p.ID == p.PkgPath && (p.PkgPath == mod || strings.HasPrefix(p.PkgPath, mod+"/")) {
			site[p.PkgPath] = a.classify(p)
		}
	}
	for _, sp := range site {
		if err := a.imports(sp, site); err != nil {
			return err
		}
	}
	return nil
}

func (a *audit) goReport(pos token.Position, r *rule, msg, fix string) error {
	rel, err := filepath.Rel(a.site.dir, pos.Filename)
	if err != nil || strings.HasPrefix(rel, "..") {
		return nil
	}
	f := finding{file: filepath.ToSlash(rel), line: pos.Line, col: pos.Column, rule: r, msg: msg, fix: fix}
	if strings.HasSuffix(f.file, "_templ.go") {
		m, err := a.templMap(f.file)
		if err != nil {
			return err
		}
		if line, col, ok := m.source(f.line, f.col); ok {
			f.file, f.line, f.col = m.src, line, col
		}
	}
	a.report(f)
	return nil
}

func (a *audit) imports(sp sitePkg, site map[string]sitePkg) error {
	p := sp.pkg
	for _, f := range p.Syntax {
		for _, spec := range f.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			report := func(r *rule, msg string) error {
				return a.goReport(p.Fset.Position(spec.Pos()), r, msg, "")
			}
			to, internal := site[path]
			var chain []string
			if internal && sp.layer == "store/migrations" {
				chain = storeChain(to, site, map[string]bool{})
			}
			switch {
			case strings.HasPrefix(path, "gorm.io/"):
				if sp.layer != "store" && sp.layer != "store/migrations" {
					err = report(ruleStoreOnly, fmt.Sprintf("%s imports %s; gorm stays inside store/", sp.rel, path))
				}
			case path == "net/http":
				if sp.layer == "views" || sp.layer == "views/pages" {
					err = report(ruleDirection, fmt.Sprintf("%s imports net/http; views render, and handlers live in routes", sp.rel))
				}
			case chain != nil:
				err = report(ruleLiveModels, fmt.Sprintf("%s imports %s", sp.rel, strings.Join(chain, ", which imports ")))
			case internal && !sp.allows(to):
				err = report(ruleDirection, fmt.Sprintf("%s imports %s; %s", sp.rel, to, sp.policy()))
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func storeChain(p sitePkg, site map[string]sitePkg, seen map[string]bool) []string {
	if p.layer == "store" {
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

func (a *audit) typeRules(p *packages.Package) error {
	var err error
	report := func(pos token.Pos, r *rule, msg, fix string) {
		if err == nil {
			err = a.goReport(p.Fset.Position(pos), r, msg, fix)
		}
	}
	for _, f := range p.Syntax {
		test := strings.HasSuffix(p.Fset.File(f.Pos()).Name(), "_test.go")
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if test || d.Tok != token.VAR {
					continue
				}
				for _, spec := range d.Specs {
					for _, name := range spec.(*ast.ValueSpec).Names {
						obj := p.TypesInfo.Defs[name]
						if obj != nil && isNamed(deref(obj.Type()), gormPath, "DB") {
							report(name.Pos(), ruleGlobal, fmt.Sprintf("package-level %s holds a %s", name.Name, types.TypeString(obj.Type(), (*types.Package).Name)), "")
						}
					}
				}
			case *ast.FuncDecl:
				if d.Body == nil {
					continue
				}
				exitOK := d.Recv == nil && (d.Name.Name == "main" && p.Name == "main" || d.Name.Name == "TestMain" && test)
				ast.Inspect(d.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if r, msg, fix := callRule(p, call, exitOK); r != nil {
							report(call.Pos(), r, msg, fix)
						}
					}
					return true
				})
			}
		}
	}
	return err
}

func callRule(p *packages.Package, call *ast.CallExpr, exitOK bool) (r *rule, msg, fix string) {
	fn, ok := typeutil.Callee(p.TypesInfo, call).(*types.Func)
	if !ok || fn.Pkg() == nil {
		return nil, "", ""
	}
	name, callee := fn.Name(), types.ExprString(call.Fun)
	switch fn.Pkg().Path() {
	case gormPath:
		if name == "AutoMigrate" {
			return ruleAutoMigrate, "AutoMigrate changes the schema outside a migration", ""
		}
	case "os":
		if name == "Exit" && !exitOK {
			return ruleExit, "os.Exit outside main", ""
		}
	case "log":
		if strings.HasPrefix(name, "Fatal") && !exitOK {
			return ruleExit, callee + " outside main", ""
		}
	case "log/slog":
		switch name {
		case "Debug", "Info", "Warn", "Error":
			if ctx := contextInScope(p.Types, call.Pos()); ctx != "" {
				return ruleSlogContext,
					fmt.Sprintf("%s drops the context; %s is in scope", callee, strings.TrimSuffix(ctx, ".Context()")),
					fmt.Sprintf("call %sContext(%s, ...)", callee, ctx)
			}
		}
	}
	return nil, "", ""
}

func contextInScope(pkg *types.Package, pos token.Pos) string {
	top := pkg.Scope()
	for s := top.Innermost(pos); s != nil && s != top; s = s.Parent() {
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
