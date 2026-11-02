package main

import (
	"bytes"
	"fmt"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"regexp"
	"strings"

	"github.com/a-h/templ"
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
		report := func(at parser.Range, r *rule, msg string) {
			a.report(finding{file: rel, line: int(at.From.Line) + 1, col: int(at.From.Col) + 1, rule: r, msg: msg})
		}
		attr := func(key parser.AttributeKey) error {
			k, ok := key.(parser.ConstantAttributeKey)
			if !ok {
				return nil
			}
			if msg := handlerAttr(k.Name); msg != "" {
				report(k.NameRange, ruleInlineJS, msg)
			}
			return nil
		}
		expr := func(e parser.Expression) error {
			if rawScript.MatchString(e.Value) {
				report(e.Range, ruleInlineJS, "templ.Raw renders a <script> element")
			}
			return nil
		}

		boost := func(key parser.AttributeKey, value string) {
			if k, ok := key.(parser.ConstantAttributeKey); ok && isBoost(k.Name) && value != "false" {
				report(k.NameRange, ruleBoost, fmt.Sprintf("literal %s value", k.Name))
			}
		}

		v := visitor.New()
		v.ConstantAttribute = func(n *parser.ConstantAttribute) error {
			boost(n.Key, n.Value)
			return attr(n.Key)
		}
		v.BoolConstantAttribute = func(n *parser.BoolConstantAttribute) error { return attr(n.Key) }
		v.ExpressionAttribute = func(n *parser.ExpressionAttribute) error {
			if tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, n.Expression.Value); err == nil && tv.Value != nil && tv.Value.Kind() == constant.String {
				boost(n.Key, constant.StringVal(tv.Value))
			}
			return attr(n.Key)
		}
		v.BoolExpressionAttribute = func(n *parser.BoolExpressionAttribute) error { return attr(n.Key) }
		v.ScriptTemplate = func(n *parser.ScriptTemplate) error {
			report(n.Range, ruleInlineJS, fmt.Sprintf("script template %s", n.Name.Value))
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
				report(n.Range, ruleInlineJS, "inline <script> body")
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

func isBoost(name string) bool {
	lower := strings.TrimPrefix(strings.ToLower(name), "data-")
	return lower == "hx-boost" || strings.HasPrefix(lower, "hx-boost:")
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
	disk, err := os.ReadFile(a.site.path(gen))
	if err != nil {
		return nil
	}
	var raw bytes.Buffer
	out, err := generator.Generate(t, &raw, generator.WithVersion(templ.Version()), generator.WithFileName(rel))
	if err != nil {
		return fmt.Errorf("generate %s: %w", rel, err)
	}
	formatted, err := format.Source(raw.Bytes())
	if err != nil || !bytes.Equal(formatted, disk) {
		return nil
	}
	a.overlay[a.site.path(gen)] = raw.Bytes()
	a.maps[gen] = out.SourceMap
	return nil
}
