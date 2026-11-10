package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	"github.com/cogikyo/hgmx/icons"
	"github.com/cogikyo/hgmx/icons/glyph"
)

const catalogPath = hgmxPath + "/icons"

var glyphMoves = map[string]string{
	"New":     "glyph.New(icons.Set)",
	"Path":    "glyph.Path",
	"Hover":   "glyph.Hover",
	"Handler": "glyph.Handler",
}

func (a *audit) iconRules() error {
	cat, err := readCatalog(icons.FS)
	if err != nil {
		return err
	}
	current := map[string]bool{}
	for _, n := range icons.All() {
		current[glyph.Ident(string(n))] = true
	}
	renamed := map[string]string{}
	for alias, name := range cat.renamed {
		if id := glyph.Ident(alias); !current[id] {
			renamed[id] = glyph.Ident(name)
		}
	}
	own := a.site.mod.Module.Mod.Path + "/" + iconsDir
	fset := token.NewFileSet()
	return a.site.walk(func(rel string, d fs.DirEntry) error {
		if d.IsDir() || !strings.HasSuffix(rel, ".go") {
			return nil
		}
		var src any
		if data, ok := a.overlay[a.site.path(rel)]; ok {
			src = data
		}
		f, err := parser.ParseFile(fset, a.site.path(rel), src, parser.SkipObjectResolution)
		if err != nil {
			return nil
		}
		if ast.IsGenerated(f) && !strings.HasSuffix(rel, "_templ.go") {
			return nil
		}
		if !slices.ContainsFunc(f.Imports, func(spec *ast.ImportSpec) bool {
			p, _ := strconv.Unquote(spec.Path.Value)
			return p == catalogPath || p == own
		}) {
			return nil
		}
		info := &types.Info{Uses: map[*ast.Ident]types.Object{}}
		conf := types.Config{Importer: unloaded{}, Error: func(error) {}}
		conf.Check(f.Name.Name, fset, []*ast.File{f}, info)
		ast.Inspect(f, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			x, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			pkg, ok := info.Uses[x].(*types.PkgName)
			if !ok {
				return true
			}
			p := pkg.Imported().Path()
			if p != catalogPath && p != own {
				return true
			}
			old := x.Name + "." + sel.Sel.Name
			report := func(msg, fix string) { a.goReport(fset.Position(sel.Pos()), ruleIcons, msg, fix) }
			if to, ok := renamed[sel.Sel.Name]; ok {
				fix := "write " + x.Name + "." + to
				if p == own {
					fix = "run go tool hgmx icons update --rename, then " + fix
				}
				report(fmt.Sprintf("%s was renamed to %s.%s, following Lucide", old, x.Name, to), fix)
			}
			if to, ok := glyphMoves[sel.Sel.Name]; ok && p == catalogPath {
				report(fmt.Sprintf("%s moved to %s, which serves only the site's icons", old, to),
					fmt.Sprintf("import %q and write %s; copy the site's icons into %s with go tool hgmx icons add <names>", glyphPath, to, iconsDir))
			}
			return true
		})
		return nil
	})
}

type unloaded struct{}

func (unloaded) Import(string) (*types.Package, error) {
	return nil, errors.New("not loaded")
}
