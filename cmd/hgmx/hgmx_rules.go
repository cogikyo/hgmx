package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
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
			fix: "remove this replace; use go.work to develop against ../hgmx",
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
	return a.hgmxLatest(latest)
}

func (r hgmxRefs) current(version string) bool {
	if rev, err := module.PseudoVersionRev(version); err == nil {
		return rev == r.master[:12]
	}
	return r.tags[version] == r.master
}

func (a *audit) hgmxLatest(latest hgmxRefs) error {
	v := a.site.version(hgmxPath)
	if !latest.current(v) {
		a.report(finding{file: "go.mod", rule: ruleHgmxLatest,
			msg: fmt.Sprintf("site requires hgmx %s; master resolves to %s", v, latest.master[:12])})
	}
	data, err := os.ReadFile(a.site.path("vendor/modules.txt"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for line := range strings.Lines(string(data)) {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != "#" || fields[1] != hgmxPath || !semver.IsValid(fields[2]) {
			continue
		}
		v := fields[2]
		if !latest.current(v) {
			a.report(finding{file: "vendor/modules.txt", rule: ruleHgmxLatest,
				msg: fmt.Sprintf("vendor contains hgmx %s; master resolves to %s", v, latest.master[:12])})
		}
		break
	}
	return nil
}
