package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
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
		if err := a.overlayTempl(t, rel); err != nil {
			return err
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

func (a *audit) overlayTempl(t *parser.TemplateFile, rel string) error {
	gen := strings.TrimSuffix(rel, ".templ") + "_templ.go"
	if _, err := os.Stat(a.site.path(gen)); err != nil {
		return nil
	}
	var raw bytes.Buffer
	out, err := generator.Generate(t, &raw, generator.WithFileName(rel))
	if err != nil {
		return fmt.Errorf("generate %s: %w", rel, err)
	}
	a.overlay[a.site.path(gen)] = raw.Bytes()
	a.maps[gen] = out.SourceMap
	return nil
}
