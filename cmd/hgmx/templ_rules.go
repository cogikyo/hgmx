package main

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/a-h/templ/generator"
	parser "github.com/a-h/templ/parser/v2"
	"github.com/a-h/templ/parser/v2/visitor"
)

var rawScript = regexp.MustCompile(`(?is)templ\.Raw\(.*<script`)

var dataScripts = map[string]bool{
	"application/json":    true,
	"application/ld+json": true,
	"text/hyperscript":    true,
}

func (a *audit) templRules() error {
	for _, rel := range a.templs {
		t, err := parser.Parse(a.site.path(rel))
		if err != nil {
			return fmt.Errorf("parse %s: %w", rel, err)
		}
		report := func(r parser.Range, msg string) {
			a.report(finding{file: rel, line: int(r.From.Line) + 1, col: int(r.From.Col) + 1, rule: ruleInlineJS, msg: msg})
		}
		attr := func(key parser.AttributeKey) error {
			k, ok := key.(parser.ConstantAttributeKey)
			if !ok {
				return nil
			}
			if msg := handlerAttr(k.Name); msg != "" {
				report(k.NameRange, msg)
			}
			return nil
		}
		expr := func(e parser.Expression) error {
			if rawScript.MatchString(e.Value) {
				report(e.Range, "templ.Raw renders a <script> element")
			}
			return nil
		}

		v := visitor.New()
		v.ConstantAttribute = func(n *parser.ConstantAttribute) error { return attr(n.Key) }
		v.BoolConstantAttribute = func(n *parser.BoolConstantAttribute) error { return attr(n.Key) }
		v.ExpressionAttribute = func(n *parser.ExpressionAttribute) error { return attr(n.Key) }
		v.BoolExpressionAttribute = func(n *parser.BoolExpressionAttribute) error { return attr(n.Key) }
		v.ScriptTemplate = func(n *parser.ScriptTemplate) error {
			report(n.Range, fmt.Sprintf("script template %s", n.Name.Value))
			return nil
		}
		v.GoCode = func(n *parser.GoCode) error { return expr(n.Expression) }
		v.StringExpression = func(n *parser.StringExpression) error { return expr(n.Expression) }
		v.CallTemplateExpression = func(n *parser.CallTemplateExpression) error { return expr(n.Expression) }
		element := v.TemplElementExpression
		v.TemplElementExpression = func(n *parser.TemplElementExpression) error {
			if err := expr(n.Expression); err != nil {
				return err
			}
			return element(n)
		}
		script := v.ScriptElement
		v.ScriptElement = func(n *parser.ScriptElement) error {
			if inlineScript(n) {
				report(n.Range, "inline <script> body")
			}
			return script(n)
		}
		if err := t.Visit(v); err != nil {
			return err
		}
	}
	return nil
}

func handlerAttr(name string) string {
	lower := strings.ToLower(name)
	switch {
	case lower == "hx-on" || lower == "data-hx-on" ||
		strings.HasPrefix(lower, "hx-on:") || strings.HasPrefix(lower, "hx-on-") ||
		strings.HasPrefix(lower, "data-hx-on:") || strings.HasPrefix(lower, "data-hx-on-"):
		return fmt.Sprintf("%s attribute evaluates code", name)
	case len(lower) > 2 && strings.HasPrefix(lower, "on") && isLetters(lower[2:]):
		return fmt.Sprintf("%s handler attribute", name)
	}
	return ""
}

func isLetters(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func inlineScript(n *parser.ScriptElement) bool {
	body := false
	for _, part := range n.Contents {
		if part.GoCode != nil || part.Value != nil && strings.TrimSpace(*part.Value) != "" {
			body = true
		}
	}
	if !body {
		return false
	}
	for _, a := range n.Attributes {
		ca, ok := a.(*parser.ConstantAttribute)
		if !ok {
			continue
		}
		if k, ok := ca.Key.(parser.ConstantAttributeKey); ok && strings.EqualFold(k.Name, "type") && dataScripts[strings.ToLower(strings.TrimSpace(ca.Value))] {
			return false
		}
	}
	return true
}

type templMap struct {
	src       string
	sm        *parser.SourceMap
	raw, disk *token.File
	rawToks   []goTok
	diskToks  []goTok
}

type goTok struct {
	off int
	tok token.Token
	lit string
}

func (a *audit) templMap(rel string) (*templMap, error) {
	if m, ok := a.maps[rel]; ok {
		return m, nil
	}
	if a.maps == nil {
		a.maps = map[string]*templMap{}
	}
	m := &templMap{src: strings.TrimSuffix(rel, "_templ.go") + ".templ"}
	a.maps[rel] = m
	t, err := parser.Parse(a.site.path(m.src))
	if errors.Is(err, fs.ErrNotExist) {
		return m, nil
	}
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", m.src, err)
	}
	var raw bytes.Buffer
	out, err := generator.Generate(t, &raw, generator.WithFileName(m.src))
	if err != nil {
		return nil, fmt.Errorf("generate %s: %w", m.src, err)
	}
	disk, err := os.ReadFile(a.site.path(rel))
	if err != nil {
		return nil, err
	}
	m.sm = out.SourceMap
	m.raw, m.rawToks = goTokens(raw.Bytes())
	m.disk, m.diskToks = goTokens(disk)
	return m, nil
}

func goTokens(src []byte) (*token.File, []goTok) {
	f := token.NewFileSet().AddFile("", -1, len(src))
	var s scanner.Scanner
	s.Init(f, src, nil, 0)
	var toks []goTok
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			return f, toks
		}
		if tok != token.SEMICOLON {
			toks = append(toks, goTok{f.Offset(pos), tok, lit})
		}
	}
}

func unique(toks []goTok, t goTok) int {
	at := -1
	for i, u := range toks {
		if u.tok == t.tok && u.lit == t.lit {
			if at >= 0 {
				return -1
			}
			at = i
		}
	}
	return at
}

func (m *templMap) source(line, col int) (int, int, bool) {
	if m.sm == nil || line > m.disk.LineCount() {
		return 0, 0, false
	}
	off := m.disk.Offset(m.disk.LineStart(line)) + col - 1
	i, ok := slices.BinarySearchFunc(m.diskToks, off, func(t goTok, off int) int { return cmp.Compare(t.off, off) })
	if !ok {
		return 0, 0, false
	}
	t, j := m.diskToks[i], -1
	switch {
	case t.tok == token.STRING:
		if unique(m.diskToks, t) >= 0 {
			j = unique(m.rawToks, t)
		}
	case len(m.rawToks) == len(m.diskToks) && m.rawToks[i].tok == t.tok && m.rawToks[i].lit == t.lit:
		j = i
	}
	if j < 0 {
		return 0, 0, false
	}
	p := m.raw.Position(m.raw.Pos(m.rawToks[j].off))
	src, ok := m.sm.SourcePositionFromTarget(uint32(p.Line-1), uint32(p.Column-1))
	return int(src.Line) + 1, int(src.Col) + 1, ok
}
