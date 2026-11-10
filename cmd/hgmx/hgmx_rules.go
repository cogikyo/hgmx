package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func (a *audit) hgmx() error {
	if a.site.mod.Module.Mod.Path == hgmxPath {
		return nil
	}
	if !slices.ContainsFunc(a.site.mod.Tool, func(t *modfile.Tool) bool { return t.Path == hgmxTool }) {
		a.report(finding{file: "go.mod", rule: ruleHgmxTool, msg: "missing tool " + hgmxTool})
	}
	replaces := a.site.hgmxReplaces()
	for _, r := range replaces {
		a.report(finding{
			file: "go.mod", line: r.Syntax.Start.Line, col: r.Syntax.Start.LineRune, rule: ruleHgmxLatest,
			msg: "replace " + r.Old.Path + " overrides the pinned hgmx version",
			fix: "remove this replace; use go.work to develop against a local hgmx checkout",
		})
	}
	if len(replaces) > 0 || a.site.version(hgmxPath) == "" {
		return nil
	}
	latest, err := latestHgmx(context.Background())
	if err != nil {
		// Keep offline site work possible when master cannot be resolved.
		newLogger(logLevel, os.Stderr).Warn("Could not resolve hgmx master; skipping hgmx-latest", slog.String("error", err.Error()))
		return nil
	}
	a.hgmxLatest(latest)
	return nil
}

func (r hgmxRefs) current(version string) bool {
	if rev, err := module.PseudoVersionRev(version); err == nil {
		return rev == r.master[:12]
	}
	return r.tags[version] == r.master
}

func (a *audit) hgmxLatest(latest hgmxRefs) {
	v := a.site.version(hgmxPath)
	if !latest.current(v) {
		a.report(finding{file: "go.mod", rule: ruleHgmxLatest,
			msg: fmt.Sprintf("site requires hgmx %s; master resolves to %s", v, latest.master[:12])})
	}
}
