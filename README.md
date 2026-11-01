# hgmx

hgmx is my shared Go toolkit for server-rendered sites using templ, htmx and Tailwind.
Each site imports only what it uses.
[CONVENTIONS.md](CONVENTIONS.md) holds the site-building rules; godoc holds the API details.

## Start a Site

Use Go 1.27 or later and put the standalone Tailwind CLI 4.3.3 on PATH.
In a new project directory:

```bash
go install github.com/cogikyo/hgmx/cmd/hgmx@master
go mod init example.com/site
hgmx init
go get github.com/cogikyo/hgmx@master
go mod tidy
go mod vendor
hgmx check
go build -o site .
./site
```

`init` writes `main.go`, `routes/`, `store/` and `views/`, including migrations, generated templ code and static assets.
It requires `go.mod`, refuses to overwrite any scaffold target and pins templ to hgmx's compiled-in version.
The scaffold serves Home, About and a full-page 404, with system-ui and OFL-licensed Iosevka fonts.
`main` builds views with `views.New`, then passes the mux from `routes.New(store, views)` to `server.Site` and `server.Run`.
`server.Site.InternalError` is the panic (500) page, not the application's 404 page.
The database pool is lazy: the scaffold serves without Postgres.
Set `PGHOST`, `PGDATABASE` and `PGUSER` for the site's database and `OTEL_*` for telemetry export.
Serving never migrates; the site binary accepts only `migrate up`, `down` and `status`.
See [the deploy contract](CONVENTIONS.md#deploy-contract) before deploying schema changes.

## Site Commands

Run these from the site's module root; there is no justfile.

| Command                     | Task                                                                |
| --------------------------- | ------------------------------------------------------------------- |
| `hgmx gen`                  | Generate templ Go and minified CSS                                  |
| `hgmx dev [-- site flags…]` | Watch templ and CSS; restart the site with disk-backed static files |
| `hgmx fmt`                  | Format Go and templ                                                 |
| `hgmx migrate …`            | Run the site's migration command through `go run .`                 |
| `hgmx check`                | Check generation, rules and formatting; then vet, build and test    |

Generation and formatting use templ in-process; the site's templ version must match hgmx's.
`check` compares generated output to the tree without writing.
Database tests use `pgtest.Open(t)`; put `initdb` and `pg_ctl` on PATH.
Missing executables or a socket bind failure skip the test unless `HGMX_PGTEST=require`; other startup failures fail.

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
hgmx palette 'night=#222536,#0184fe' -f tailwind -o views/static/css/colors.css
```

The input is `<name>=<base>,<primary>[,<accent>]`; repeat it for multiple themes.
Formats are `tailwind` (default), `css` and `lua`; Lua takes one theme.
Without `-o`, output goes to stdout.

## Build the CLI

From this repository: `go build -o hgmx ./cmd/hgmx`.

`hgmx.Version()` reports the main module's Go build-info version only when that module is hgmx itself, as in the CLI.
It returns empty in a site; set the site's own version separately.
