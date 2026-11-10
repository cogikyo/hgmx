package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
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
	{"api", []string{"api", "proto", "store", domain}},
	{"proto", []string{"proto"}},
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
		Dir:        a.site.dir,
		Env:        a.site.env(),
		BuildFlags: []string{readonly},
		Overlay:    a.overlay,
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
	for _, p := range pkgs {
		if strings.HasSuffix(p.ID, ".test") {
			continue
		}
		a.typeRules(p)
		if p.ID == p.PkgPath && (p.PkgPath == mod || strings.HasPrefix(p.PkgPath, mod+"/")) {
			site[p.PkgPath] = a.classify(p)
		}
	}
	for _, sp := range site {
		a.imports(sp, site)
	}
	return nil
}

func (a *audit) goReport(pos token.Position, r *rule, msg, fix string) {
	rel, err := filepath.Rel(a.site.dir, pos.Filename)
	if err != nil || strings.HasPrefix(rel, "..") {
		return
	}
	f := finding{file: filepath.ToSlash(rel), line: pos.Line, col: pos.Column, rule: r, msg: msg, fix: fix}
	if sm := a.maps[f.file]; sm != nil {
		if src, ok := sm.SourcePositionFromTarget(uint32(f.line-1), uint32(f.col-1)); ok {
			f.file = strings.TrimSuffix(f.file, "_templ.go") + ".templ"
			f.line, f.col = int(src.Line)+1, int(src.Col)+1
		}
	}
	a.report(f)
}

func (a *audit) imports(sp sitePkg, site map[string]sitePkg) {
	p := sp.pkg
	for _, f := range p.Syntax {
		if generated(p, f) {
			continue
		}
		for _, spec := range f.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			report := func(r *rule, msg string) {
				a.goReport(p.Fset.Position(spec.Path.Pos()), r, msg, "")
			}
			to, internal := site[path]
			var chain []string
			if internal && sp.layer == "store/migrations" {
				chain = storeChain(to, site, map[string]bool{})
			}
			switch {
			case strings.HasPrefix(path, "gorm.io/"):
				if sp.layer != "store" && sp.layer != "store/migrations" {
					report(ruleStoreOnly, fmt.Sprintf("%s imports %s; gorm stays inside store/", sp.rel, path))
				}
			case path == "net/http":
				if sp.layer == "views" || sp.layer == "views/pages" {
					report(ruleDirection, fmt.Sprintf("%s imports net/http; views render, and handlers live in routes", sp.rel))
				}
			case chain != nil:
				report(ruleLiveModels, fmt.Sprintf("%s imports %s", sp.rel, strings.Join(chain, ", which imports ")))
			case internal && !sp.allows(to):
				report(ruleDirection, fmt.Sprintf("%s imports %s; %s", sp.rel, to, sp.policy()))
			}
		}
	}
}

func generated(p *packages.Package, f *ast.File) bool {
	return ast.IsGenerated(f) && !strings.HasSuffix(p.Fset.File(f.Pos()).Name(), "_templ.go")
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

func (a *audit) typeRules(p *packages.Package) {
	report := func(pos token.Pos, r *rule, msg, fix string) {
		a.goReport(p.Fset.Position(pos), r, msg, fix)
	}
	for _, f := range p.Syntax {
		if generated(p, f) {
			continue
		}
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
					switch n := n.(type) {
					case *ast.CallExpr:
						if r, msg, fix := callRule(p, n, exitOK); r != nil {
							report(n.Pos(), r, msg, fix)
						}
					case *ast.AssignStmt:
						for i, lhs := range n.Lhs {
							ix, ok := lhs.(*ast.IndexExpr)
							if !ok || !isNamed(p.TypesInfo.TypeOf(ix.X), "net/http", "Header") {
								continue
							}
							var values []ast.Expr
							if len(n.Rhs) == len(n.Lhs) {
								values = elems(n.Rhs[i])
							}
							if cspLiteral(p, ix.Index, values...) {
								report(lhs.Pos(), ruleCSP, types.ExprString(lhs)+" sets Content-Security-Policy", "")
							}
						}
					case *ast.CompositeLit:
						if !isNamed(p.TypesInfo.TypeOf(n), "net/http", "Header") {
							break
						}
						for _, e := range n.Elts {
							if kv, ok := e.(*ast.KeyValueExpr); ok && cspLiteral(p, kv.Key, elems(kv.Value)...) {
								report(kv.Pos(), ruleCSP, "http.Header literal sets Content-Security-Policy", "")
							}
						}
					}
					return true
				})
			}
		}
	}
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
	case "net/http":
		if recv := fn.Signature().Recv(); (name == "Set" || name == "Add") && recv != nil && isNamed(recv.Type(), "net/http", "Header") && cspLiteral(p, call.Args[0], call.Args[1]) {
			return ruleCSP, callee + " sets Content-Security-Policy", ""
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

func cspLiteral(p *packages.Package, key ast.Expr, values ...ast.Expr) bool {
	v := p.TypesInfo.Types[key].Value
	if v == nil || v.Kind() != constant.String || !strings.EqualFold(constant.StringVal(v), "Content-Security-Policy") {
		return false
	}
	return len(values) == 0 || slices.ContainsFunc(values, func(e ast.Expr) bool { return !cspHeader(p, e) })
}

func elems(e ast.Expr) []ast.Expr {
	if c, ok := ast.Unparen(e).(*ast.CompositeLit); ok {
		return c.Elts
	}
	return nil
}

func cspHeader(p *packages.Package, e ast.Expr) bool {
	call, ok := ast.Unparen(e).(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := typeutil.Callee(p.TypesInfo, call).(*types.Func)
	if !ok || fn.Name() != "Header" {
		return false
	}
	recv := fn.Signature().Recv()
	return recv != nil && isNamed(recv.Type(), hgmxPath+"/server", "CSP")
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
