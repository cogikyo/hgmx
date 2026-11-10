package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	own := a.site.mod.Module.Mod.Path + "/" + iconsDir
	local, managed, err := iconSet(a.site.path(iconsDir))
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	var files []*ast.File
	err = a.site.walk(func(rel string, d fs.DirEntry) error {
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
		uses := false
		for _, spec := range f.Imports {
			p, _ := strconv.Unquote(spec.Path.Value)
			uses = uses || p == catalogPath || p == own
			if p == catalogPath && local {
				a.goReport(fset.Position(spec.Path.Pos()), ruleIconsSource,
					fmt.Sprintf("imports hgmx's full catalog while %s holds a local icon set; one glyph handler serves only one", iconsDir), "")
			}
		}
		if uses {
			files = append(files, f)
		}
		return nil
	})
	if err != nil || len(files) == 0 && !managed {
		return err
	}
	cat, err := a.site.catalog()
	if err != nil {
		return err
	}
	svgs, err := fs.Glob(cat.fsys, "svg/*.svg")
	if err != nil {
		return err
	}
	current := map[string]bool{}
	for _, rel := range svgs {
		name, _ := relIcon(rel)
		current[glyph.Ident(name)] = true
	}
	renamed := map[string]string{}
	for alias, name := range cat.renamed {
		if id := glyph.Ident(alias); !current[id] {
			renamed[id] = glyph.Ident(name)
		}
	}
	for _, f := range files {
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
	}
	if managed {
		return a.iconDrift(cat)
	}
	return nil
}

func iconSet(dir string) (local, managed bool, err error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return false, false, nil
	}
	for _, e := range entries {
		name := e.Name()
		switch {
		case e.IsDir():
		case name == iconsList:
			managed = true
		case name == iconsGo:
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return false, false, err
			}
			managed = managed || strings.HasPrefix(string(data), iconsHeader)
			local = true
		case strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".templ"):
			local = true
		}
	}
	return local || managed, managed, err
}

func (a *audit) iconDrift(cat iconCatalog) error {
	root, err := os.OpenRoot(a.site.path(iconsDir))
	if err != nil {
		return err
	}
	defer root.Close()
	if _, err := root.Stat(iconsList); errors.Is(err, fs.ErrNotExist) {
		a.report(finding{
			file: iconsDir + "/" + iconsList,
			rule: ruleIconsDrift,
			msg:  iconsList + " is missing from a generated icon set",
			fix:  "restore it from version control, or recreate the list with go tool hgmx icons add <names>",
		})
		return nil
	}
	res, err := syncIcons(root, cat, iconsCheck, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", iconsDir, err)
	}
	for _, it := range res.items {
		if !it.verdict.drift() {
			continue
		}
		msg, fix := it.text(true)
		a.report(finding{file: iconsDir + "/" + iconRel(it.name, "svg"), rule: ruleIconsDrift, msg: it.name + " " + msg, fix: fix})
	}
	for _, rel := range res.stale {
		a.report(finding{file: iconsDir + "/" + rel, rule: ruleIconsDrift, msg: rel + " out of date"})
	}
	return nil
}

type unloaded struct{}

func (unloaded) Import(string) (*types.Package, error) {
	return nil, errors.New("not loaded")
}
