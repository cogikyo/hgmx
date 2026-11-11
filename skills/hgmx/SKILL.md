---
name: hgmx
description: Build, change, or review a Go site that uses github.com/cogikyo/hgmx, including site conventions, shared packages, native app routing, generation, checks, updates, and managed agent skills.
---

# Build a site on hgmx

hgmx is a shared Go toolkit for server-rendered sites using templ, htmx, and Tailwind.
Each site imports the infrastructure it needs and owns its pages and copied components.
Use **hgmx-connect** for native API contracts and service wiring, **hgmx-ios** for SwiftUI, and **hgmx-android** for Compose.
Browser pages remain server-rendered HTML and do not consume the native API.

## Read the source of truth

Read the conventions from the site's selected hgmx module before changing the site:

```bash
cat "$(go list -m -f '{{.Dir}}' github.com/cogikyo/hgmx)/CONVENTIONS.md"
```

Use `go doc` for each hgmx package and symbol you touch, with its full import path, such as `github.com/cogikyo/hgmx/server`.
Inspect the site's code and configuration for its current behavior; do not copy conventions from an older site.
Re-read the module's conventions after an update.

## Daily loop

Run commands from the site's module root through `go tool hgmx`.
Declare `tool github.com/cogikyo/hgmx/cmd/hgmx` in `go.mod` so the CLI and library share one version.
Keep the module root free of `.go` and `.templ` files, process wiring beside `cmd/server/main.go`, other commands in `cmd/<name>/` and library code outside `cmd/`.

1. Run **`go tool hgmx check`** to find the current failures.
   It checks contracts, generation, runtime pins, site rules, credentials, format, vet, build and tests without writing fixes; contracts and pins run only when present.
2. Keep the site on **latest hgmx master** while hgmx is being built out.
   Run `go tool hgmx update` when behind; it updates the library and tool, tidies modules, and re-syncs skills.
   A failed remote lookup makes the latest-version check warn and skip, so an offline pass does not prove currency.
3. Use **`go tool hgmx dev`** while editing.
   Run `go tool hgmx gen` and `go tool hgmx fmt` to write generated and formatted files.
4. Run **`go tool hgmx check`** again before handing off.
   Report any failed or skipped check, including database tests skipped for missing local prerequisites.
   Native `--skip` makes check exit nonzero; it is not a substitute for platform generation or builds.

## Managed skills

`go tool hgmx skills` installs the module's skills in `.agents/skills/` and links them into `.claude/skills/` when `.claude/` exists.
`check` reports stale or edited managed copies; sync refuses to overwrite a file the site edited.
Do not edit managed files to tailor a site: put local rules in `SITE.md` beside the relevant `SKILL.md`.
hgmx leaves `SITE.md` untouched.

If `SITE.md` sits beside this file, read it; its rules override this skill.
