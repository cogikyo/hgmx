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
go build -o site .
./site
```

`init` writes `main.go`, `routes/`, `store/` and `views/`, including migrations, generated templ code and static assets.
It requires `go.mod`, refuses to overwrite any scaffold target and pins templ to hgmx's compiled-in version.
The scaffold serves Home, About and a full-page 404, with system-ui and OFL-licensed Iosevka fonts.
`main` builds views with `views.New`, then passes the mux from `routes.New(store, views)` to `server.Site` and `server.Run`.
`server.Site` supplies the base CSP and a fixed panic (500) page; set `InternalError` to custom HTML or leave it empty for hgmx's default.
The database pool is lazy: the scaffold serves without Postgres.
Set `PGHOST`, `PGDATABASE` and `PGUSER` for the site's database and `OTEL_*` for telemetry export.
Serving never migrates; use `migrate up`, `down [track]` or `status`.
See [the deploy contract](CONVENTIONS.md#deploy-contract) before deploying schema changes.

## Site Commands

Run these from the site's module root; there is no justfile.

| Command                             | Task                                                                |
| ----------------------------------- | ------------------------------------------------------------------- |
| `go tool hgmx gen`                  | Generate templ Go and minified CSS                                  |
| `go tool hgmx dev [-- site flags…]` | Watch templ and CSS; restart the site with disk-backed static files |
| `go tool hgmx fmt`                  | Format Go and templ                                                 |
| `go tool hgmx migrate …`            | Run the site's migration command through `go run .`                 |
| `go tool hgmx update`               | Update hgmx to master and tidy; re-vendor if present                |
| `go tool hgmx credentials`          | Build the site and write `credentials.json`                         |
| `go tool hgmx credentials check …`  | Test SMTP and Google client authentication                          |
| `go tool hgmx check`                | Check generated files, credentials, rules, format, vet, build, test |

Generation and formatting use templ in-process; the site's templ version must match hgmx's.
`check` compares generated output to the tree without writing.
Its `credentials` stage fails if the site-root `credentials.json` is missing or differs from the site's declarations and registered flags, including formatting differences.
Declare secrets in a site-owned `credential.Registry`, read through its handles, then run `go tool hgmx credentials` and commit the file with the code, even for sites with no secrets.
See [credential wiring](CONVENTIONS.md#hgmx-credentials) for the early command dispatch.

`go tool hgmx credentials check [name...] [-- site flags...]` requires a current manifest and checks the named SMTP/Google credentials, or all declared SMTP/Google targets when no names are given.
It never prints secrets; SMTP sends no mail, and the Google probe does not test browser sign-in.
Ordinary `check` does not run these live provider checks.

Database tests use `pgtest.Open(t)`; put `initdb` and `pg_ctl` on PATH.
Missing executables or a socket bind failure skip the test unless `HGMX_PGTEST=require`; other startup failures fail.

## Agent Skills

hgmx ships harness-neutral agent skills for building sites and integrating Stripe payments.
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

## Packages

- **`server`** composes the root handler, version reloads, recovery and graceful shutdown.
- **`render`** buffers templ HTML with ETags and conditional requests.
- **`telemetry`** exports OTLP/HTTP traces, logs and metrics from `OTEL_*` settings, with local logging and a browser beacon.
- **`pg`** opens Postgres from `PG*` settings and supplies IDs, JSON, enums and query scopes; `pg/pgtest` supplies test databases.
- **`migrate`** runs Go-defined schema steps with transactions and a session lock.
- **`static`** serves file sets at hashed URLs or fixed root URLs, with a disk-backed dev mode.
- **`scripts`** serves selected htmx, hx-preload, navigation support, hyperscript and Motion scripts in dependency order.
- **`palette`** generates contrast-tier scales and semantic tokens from base, primary and optional accent colours.
- **`theme`** serves palette stylesheets, the theme switcher and explorer.
- **`icons`** provides 441 Lucide templ components with Motion 14.7.1 hover animations; see the [motion grammar](icons/GRAMMAR.md) when adding one.
- **`webapp`** builds an install manifest and icons from a PNG, without offline support.

## Generate a Palette

```bash
go tool hgmx palette 'night=#222536,#0184fe' -f tailwind -o views/static/css/colors.css
```

The input is `<name>=<base>,<primary>[,<accent>]`; repeat it for multiple themes.
Formats are `tailwind` (default), `css` and `lua`; Lua takes one theme.
Without `-o`, output goes to stdout.

## Build the CLI

From this repository: `go build -o hgmx ./cmd/hgmx`.

`hgmx.Version()` reports the main module's Go build-info version only when that module is hgmx itself, as in the CLI.
It returns empty in a site; set the site's own version separately.

## Regenerate RPC Code

Regenerate committed RPC Go sources from `proto/`:

```bash
GOWORK=off CGO_ENABLED=0 go run github.com/bufbuild/buf/cmd/buf@v1.75.0 generate
```
