# hgmx

hgmx is my shared Go toolkit for server-rendered sites using templ, htmx and Tailwind.
Each site imports only what it uses.
[CONVENTIONS.md](CONVENTIONS.md) holds the site-building rules; godoc holds the API details.

## Start a Site

Use Go 1.27 or later and put the standalone Tailwind CLI 4.3.3 on PATH.
In a new project directory:

```bash
go mod init example.com/site
go get -tool github.com/cogikyo/hgmx/cmd/hgmx@master
go tool hgmx init
go mod tidy
go tool hgmx check
go build -o site ./cmd/server
./site
```

`init` writes `.gitignore` (ignoring `/server`), `cmd/server/main.go`, `credentials.json`, `routes/`, `store/` and `views/`.
It requires `go.mod`, pins templ to hgmx's version and refuses existing scaffold targets except `.gitignore`.
It keeps existing `.gitignore` content and appends `/server` if missing.
The [site layout](CONVENTIONS.md#layout) keeps the root free of Go packages and process wiring beside `cmd/server/main.go`.
The scaffold serves Home, About and a full-page 404 without Postgres, using system-ui and OFL-licensed Iosevka fonts.
Set `PGHOST`, `PGDATABASE` and `PGUSER` for the database and `OTEL_*` for telemetry.
Serving never migrates; read [the deploy contract](CONVENTIONS.md#deploy-contract) before schema changes.

## Site Commands

Run these from the site's module root; there is no justfile.

| Command                             | Task                                                                |
| ----------------------------------- | ------------------------------------------------------------------- |
| `go tool hgmx gen`                  | Generate templ Go, minified CSS and app contract code               |
| `go tool hgmx dev [-- site flags…]` | Watch templ and CSS; restart the site with disk-backed static files |
| `go tool hgmx fmt`                  | Format Go and templ                                                 |
| `go tool hgmx migrate …`            | Run `go run ./cmd/server migrate …`                                 |
| `go tool hgmx update`               | Update hgmx to master and tidy                                      |
| `go tool hgmx credentials`          | Build the site and write `credentials.json`                         |
| `go tool hgmx credentials check …`  | Test SMTP and Google client authentication                          |
| `go tool hgmx check`                | Check generated files, rules, credentials, format, vet, build, test |

`check` compares generated output to the tree without writing.
Commit `credentials.json` even for zero-secret sites; the credentials stage compares it byte-for-byte.
See [credential wiring and live checks](CONVENTIONS.md#hgmx-credentials); ordinary `check` skips live provider checks.

Database tests use `pgtest.Open(t)`; put `initdb` and `pg_ctl` on PATH.
Missing executables or a socket bind failure skip the test unless `HGMX_PGTEST=require`; other startup failures fail.

## Agent Skills

hgmx ships harness-neutral agent skills for sites, Connect APIs, SwiftUI apps, Compose apps and Stripe payments.
Run `go tool hgmx skills` to install them in `.agents/skills/`; `go tool hgmx update` re-syncs them when updating hgmx.
When `.claude/` exists, hgmx also links the skills into `.claude/skills/`.
Put site-specific rules in `SITE.md` beside a skill's `SKILL.md`; they override the skill, and hgmx never changes them.
Do not edit managed copies: the `hgmx-skills` check reports drift, and sync stops on edits with instructions to move them into `SITE.md`.

## Accounts

`auth` provides per-site passwordless accounts with passkeys, email codes and revocable sessions.
See the [package docs](https://pkg.go.dev/github.com/cogikyo/hgmx/auth) for migration setup, `Config`, modes and styling hooks.
With the database, mail queue and config ready:

```go
a, err := auth.New(db, queue, cfg)
if err != nil {
	return err
}
a.Mount(mux)
mux.Handle("GET /members", a.Require(auth.Anyone, members))
```

`auth` also serves a versioned Connect-RPC API for SwiftUI and Compose apps; see [Serving an App API](CONVENTIONS.md#serving-an-app-api) for wiring.

## Start a Phone App

The starters target iOS 17+ with SwiftUI and Android 9+ with Compose.
One Go backend serves the website and the native API; the website still uses HTML, not the app API.
First [register each app and mount the API](CONVENTIONS.md#serving-an-app-api) in the site.
Use that registration's public ID and HTTPS API origin on a domain name when adopting its starter:

```bash
go tool hgmx init apps android --name "Example" --origin https://example.com --app-id example-android --package com.example.app
go tool hgmx init apps ios --name "Example" --origin https://example.com --app-id example-ios --bundle-id com.example.app
go tool hgmx gen
go tool hgmx check
```

Each command copies its platform project into `apps/android/` or `apps/ios/`, substitutes the identity and records the source hgmx version.
It also creates `apps/compatibility.json` if absent.
The copy is site-owned: rerunning `init apps` refuses an existing platform directory and leaves its edits alone.
Port later starter fixes deliberately; `hgmx update` does not sync handwritten app code.
Android accepts optional `--google-client-id` for the Google web/server client; iOS rejects it because its starter has no Google sign-in yet.

For generation prerequisites, output ownership and partial local checks, see [App Generation and Checks](CONVENTIONS.md#app-generation-and-checks).
After generation, open `apps/ios/Starter.xcodeproj` with its `Starter` scheme, or build Android from `apps/android/` with `./gradlew :app:lintDebug :app:assembleDebug`.
The shipped [iOS skill](skills/hgmx-ios/SKILL.md) and [Android skill](skills/hgmx-android/SKILL.md) identify the files to change for a screen, identity and token storage.
The [Connect skill](skills/hgmx-connect/SKILL.md) covers the site's service boundary.
Before planning a store release, read [App Launch Requirements](CONVENTIONS.md#app-launch-requirements); neither starter is submission-ready.

## Packages

- **`server`** composes the root handler, version reloads, recovery and graceful shutdown.
- **`render`** buffers templ HTML with ETags and conditional requests.
- **`telemetry`** exports OTLP/HTTP traces, logs and metrics from `OTEL_*` settings, with local logging and a browser beacon.
- **`pg`** opens Postgres from `PG*` settings and supplies IDs, JSON, enums and query scopes; `pg/pgtest` supplies test databases.
- **`migrate`** runs Go-defined schema steps with transactions and a session lock.
- **`static`** serves file sets at hashed URLs or fixed root URLs, with a disk-backed dev mode.
- **`scripts`** serves selected htmx, hx-preload, navigation support, hyperscript and Motion scripts in dependency order.
- **`palette`** generates OKLCH shades and semantic colour tokens from one seed.
- **`icons`** holds the full catalog for previews; sites [copy selected icons](CONVENTIONS.md#hgmx-icons) with `go tool hgmx icons add <names>` and import `icons/glyph` for rendering and hover animations.
- **`webapp`** builds an install manifest and icons from a PNG, without offline support.

## Generate a Palette

From this repository's root, overwrite `library/views/static/css/colors.css`:

```bash
go run ./cmd/hgmx palette '#222536'
```

Add `--lua` to write a Lua palette to stdout instead.
Live theme switching is not available.

## Build the CLI

From this repository: `go build -o hgmx ./cmd/hgmx`.

`hgmx.Version()` reports the main module's Go build-info version only when that module is hgmx itself, as in the CLI.
It returns empty in a site; set the site's own version separately.

## Regenerate RPC Code

From this repository's root, regenerate committed RPC Go sources from `proto/`:

```bash
go run ./cmd/hgmx gen
```

The hidden `hgmx buf` command is a CI-only maintainer entry point to hgmx's pinned Buf.
The compatibility job uses it against the latest release tag; site authors use `gen` and `check`.

## Test App Adoption

From this repository, use `starters/fixture/` before changing generators, contracts, dependency pins or starters:

```bash
starters/fixture/adopt.sh /tmp/opencode/hgmx-fixture android ios
starters/fixture/bare.sh /tmp/opencode/hgmx-bare-ios ios
```

Both destinations must be empty.
`adopt.sh <dest> [android] [ios]` defaults to both platforms and creates a temporary site using this hgmx checkout, with auth, a Notes RPC and native Notes screens.
It exercises identity substitution and generation for a site-owned contract; `bare.sh <dest> <android|ios>` isolates the starter with only hgmx's contracts.
Build the adopted apps with the platform SDKs and run `go tool hgmx check` in the fixture; adoption alone does not prove compilation or a live call.

`starters/fixture/seeded.sh <fixture> [seed...]` copies an adopted fixture to scratch and expects specific failures for generated drift, pins, dependency locks and contract edits.
Android seeds also cover managed-input edits, breaking changes and release-baseline failures; those release seeds need a committed fixture with local history.
Run the seeds after a clean check and native build, with JDK/Android SDK or Xcode available for the selected platforms.
These checks do not prove physical passkeys, production associations or store readiness.
