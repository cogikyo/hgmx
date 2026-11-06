# Site Conventions

These are the default ways my Go sites do each common job.
Read this page instead of copying an older site.
A rule holds unless its own text names an exception.

## Hypermedia Model

- **The server owns all state that matters.** HTML carries that state and the available actions.
  - Client JS holds only client-only state, such as open or closed and animation phase.
- **Every address-bar URL returns a complete page** when requested directly.
  - htmx 4 refetches pages on Back and Forward instead of keeping snapshots, so this is required.
- **A route is a resource, and its state lives in the URL.**
  - Search uses the page route itself, as shown below, not a separate fragment route.
- **Invalid input** re-renders the full page with the error and status 422, never bare text.
- **The server renders UI state it knows**, such as the active nav link, selection and counts.
  - JS never reconciles it; the server writes `aria-current="page"` and the like.
- **No JSON API** feeds browser rendering, and no route negotiates between JSON and HTML.

**Search-as-you-type** filters the page in place.
`GET /items?q=` renders the filtered full page, so reload and share keep the query.
This spelling was measured with htmx 4.0.0 in Chromium and Firefox, under the boost from [Navigation](#navigation).
It relies on `scripts.Navigation`, described there.

```html
<form action="/items" role="search">
  <input
    type="search"
    name="q"
    value="…"
    hx-get="/items"
    hx-trigger="input changed delay:200ms"
    hx-select="#items"
    hx-target="#results"
    hx-swap="innerHTML"
    hx-sync="queue last"
    hx-select-oob=""
    hx-replace-url="true"
  />
</form>
<section id="results"><div id="items">…</div></section>
```

- **Target a stable wrapper** (`#results`, `innerHTML`) and select the inner list (`#items`).
  - htmx 4.0.0 resolves a queued request's target when it fires.
  - A swapped-out target leaves the next response in a detached node, a stale list under the new URL.
- **`hx-sync="queue last"`** keeps the newest keystroke while a search is in flight.
- **`hx-select-oob=""`:** without it, the inherited body value re-morphs the nav on every keystroke.
  - Any value on the element overrides the inherited one.
- **`hx-replace-url`:** without it, the URL doesn't update.
  - `hx-push-url` would add a history entry per keystroke.
- **The form is the fallback:** Enter submits it through the boost.
  - With JS off, the plain form GET filters the page.

## Navigation

Every page route returns the full page, and htmx picks what to swap.
This spelling was verified with htmx 4.0.0 in Chromium and Firefox; WebKit is untested.

```templ
<body
  hx-boost:inherited={ scripts.Boost }
  hx-select-oob:inherited="#nav-links:outerMorph"
  hx-config:inherited={ server.VersionConfig(version) }
>
  <nav id="nav">
    …
    <ul id="nav-links">
      <!-- server-rendered aria-current -->
    </ul>
  </nav>
  <main id="main" hx-history-elt>…</main>
</body>
```

```css
#main {
  view-transition-name: main;
}
:root {
  view-transition-name: none;
}
::view-transition {
  pointer-events: none;
}
::view-transition-group(main),
::view-transition-old(main),
::view-transition-new(main) {
  animation-duration: 120ms;
}
```

**Use it this way:**

- **Links are plain `<a href>`**, so they work without JS and get boosted.
  - Use explicit `hx-get` and `hx-target` only where a link behaves differently from navigation.
  - Don't put `hx-get` with `hx-push-url` on nav links.
- **Handlers always return the full page:** no `HX-Request` branch, no `Vary`, no partial template.
- **The body takes no `hx-sync`** when loading `scripts.Navigation`.
- **Leave `view-transition-name` off `#main`** if fixed or positioned layers must interleave with its content.
  The name creates a stacking context even outside transitions, so overlays can block links.
  With the root name also off, boost crossfades nothing visible.
- **Version reloads pair `server.VersionConfig` with `server.PageVersion`.**
  - Emit `server.VersionConfig(version)` on `<body>` and wrap page handlers with `PageVersion`.
  - A GET or HEAD carrying a different build version gets a full reload; requests without the header pass through.
  - `scripts.Navigation` re-sends the config headers on history restores.
  - Auth mounts `/auth/*` without `PageVersion`, so per-route wrappers leave those pages without version reloads.
    On browser-only sites, wrap the root handler returned by `Site.Handler` with `PageVersion` to cover them too.
- **Per-page regions outside `#main`**, such as a breadcrumb, are listed in `hx-select-oob` and morphed.
  - Morph only regions whose attributes only the server writes.
  - Morph strips classes and inline styles that JS set.
- **Persistent elements outside `#main`** keep their DOM node and JS state across clicks, Back and Forward.
  - Examples are a nav controller, Motion state, or a background canvas.
- **Validation and 404 pages use the full layout.** htmx 4 selects their `#main` as on success.
  - Recovery serves fixed `InternalError` HTML for a 500; hgmx's default has no nav.
  - A boosted 500 leaves the previous page marked in the nav; this is accepted.
- **Links to non-page responses** get `hx-boost="false"` or `download`; otherwise `#main` ends up empty.
  - These include feeds, PDFs, `/healthz`, downloads and pages without `#main`.
- **Inheritance is explicit in htmx 4:** add `:inherited` on purpose and expect nothing implicit.

**Load `scripts.Navigation` for the boost.**
Its implementation is [scripts/hgmx.js](scripts/hgmx.js).
In htmx 4.0.0, it drops a response once its body arrives if its source left the page or its request was aborted.
It also drops an older navigation response when a newer navigation has started, ordered by trigger time.
This happens before redirects, history, title or swap; a newer navigation alone does not cancel widget responses.
After a navigation swap, focus moves to the main heading or `main`, unless already inside `main`.

**`hx-sync` cannot replace this:**

- **`body:replace` loses track from the third request:** an aborted request's cleanup frees the queue slot.
- **Widgets without their own `hx-sync`** share the body queue and abort each other.
- **Each abort logs `htmx:error`.**
- **The preload promise ignores the abort signal.**

**Each part of the spelling prevents a measured failure:**

- **`innerHTML strip`:** plain `innerHTML` nests `<main>` in `<main>`.
  - `innerMorph` throws `HierarchyRequestError`.
  - Morphing `#main` strips JS-set attributes and carries listeners across pages.
- **`show:top showTarget:body`:** the boost default scrolls `#main` under a sticky nav.
  - `show:window:top` isn't parsed; `scroll-margin-top` on `#main` also works.
- **`hx-select-oob` on `<body>`:** Back and Forward re-request from `[hx-history-elt]`, which inherits from `<body>`.
  - Placed on the nav or inside the boost value, it leaves the nav stale after Back.
- **Root crossfade off:** otherwise the nav ghosts mid-transition.
- **No `hx-push-url:inherited` on `<body>`:** restores then push new entries and destroy Forward.

**Costs:** each navigation fetches the whole page.
Inline scripts in `#main` run on every swap, restores included, which is why [Scripting](#scripting) bans them.
Text typed into an input in `#main` during a navigation is lost.

The htmx 4.0.0 source differs from its docs.
Boosted swaps default to `config.defaultSwap` (`innerHTML`), not `outerSync`.
Restores always use `outerSync` and send no `HX-Request` or `HX-Target`.
Preload requests never carry `HX-Request-Type`.

**Fragments are a fallback** for a site that truly needs them: use explicit `hx-get` with `hx-preload`, branch on `HX-Target: main#main` with `Vary: HX-Target`.
Never branch on `HX-Request-Type`; preload requests lack it, so the full page would land inside `#main`.

### htmx: Boost

**Use `hx-boost:inherited={ scripts.Boost }` with `scripts.Navigation`.**
The shared value keeps navigation settings in step with hgmx.
The `htmx/boost` check rejects literal boost values, including string expressions, except `"false"` to opt out.

## Prefetch and Transitions

- **htmx sites load `scripts.Navigation`**, which includes htmx and hx-preload.
  - It preloads boosted anchors on mousedown and touchstart and reuses the response for 5 s.
  - Plain links need no preload attributes.
- **Don't add Speculation Rules** to links hx-preload covers: they fetch twice, with no documented shared cache.
- **Sites without htmx add no prefetch** while they have only a page or two.
  - When such a site grows pages, it sends a `Speculation-Rules` response header.
  - Inline rules would need `'inline-speculation-rules'` in the CSP.
  - It also adds `@view-transition { navigation: auto; }`.
- **Only `#main` animates**, and the root animation is off.
  - `#main` carries `view-transition-name` unless the stacking exception in [Navigation](#navigation) applies.
- **Page transitions take 120 ms**, and reduced motion turns them off through `scripts.Navigation`, which checks each swap.
- **Function never depends on animation.**
- **Motion handles choreographed component animation**, such as a logo, menu or icons.
  - View transitions handle page-to-page changes; the two don't overlap.

## Scripting

- **hyperscript** covers small, local behaviour next to its markup.
  - Examples are toggles, dismissal, event relays (`send hovered to …`), timed classes and event waits.
  - Repeated logic becomes a `behavior` in a `._hs` file, attached with `install`.
- **JS modules** cover widgets with real logic, such as a scroll controller or Motion choreography.
  - They live in `.js` files served through `static`, whose set-hash URLs keep relative imports working.

### templ: No Inline JS

- **No inline scripts:** no inline modules, classic scripts, importmaps or `templ.Raw("<script>…")`.
  - Inline scripts rerun on every swap; external modules evaluate once and cache.
  - hgmx's base CSP blocks inline JS.
  - Data scripts (`application/json`, `application/ld+json`) and `text/hyperscript` are allowed.
- **Module JS is idempotent across swaps.**
  - Delegate events on `document`, or initialise from `htmx:after:process`.
  - Don't guard with `data-*` "ready" markers; morph removes attributes the server didn't send.
- **One runtime per responsibility:** no `hx-live` or Alpine next to hyperscript.
- **No handler attributes:** no `on*=`, `hx-on`, `hx-on:` or `hx-on-`, including `data-hx-on` forms.
  - They evaluate code, which strict CSP blocks; hyperscript `_=` is allowed.
- **No templ `script` templates**; they are legacy.

## templ

- **Components take ordinary parameters**; a page takes one view-model struct built by its handler.
  - No mandatory Props struct and no style or customisation interfaces.
- **Wrappers use `{ children... }`**; named slots are `templ.Component` parameters.
- **Pass-through attributes** (`aria-*`, `data-*`, `hx-*`) arrive as a `templ.Attributes` spread.
- **Dynamic URLs in `hx-*` attributes** go through `templ.URL`.
- **Data for JS** goes in `data-*` attributes or `templ.JSONScript`.
- **Context carries only request-scoped values**, such as the current path or user, read through typed getters.
  - Process constants such as the build version are package variables set once at startup.
- **No templ `css` components**; Tailwind owns styling.
- **Render pages with `render.HTML`**, buffered with an ETag and conditional handling for 200 GET and HEAD responses.
  - On a render error, return if the request was cancelled; otherwise panic into recovery.
  - Stream only for a measured need.
- **Generated templ and CSS are committed**, and regenerated and checked before each commit.
  - Run `go tool hgmx gen`; `go tool hgmx check` checks generation and formatting without writing.

### templ: Version

- **The site's required templ version equals hgmx's compiled-in version.** `go tool hgmx init` sets it in `go.mod`.
  - `go tool hgmx check` reports a mismatch and the command that fixes it.
- **hgmx runs templ generation and formatting in-process.** Use `go tool hgmx gen`, `go tool hgmx dev` and `go tool hgmx fmt`.
  - Sites add no templ `tool` line and do not run `go run templ@version`.

## Tailwind v4

- **Standalone CLI 4.3.3** with one input, `main.css`.
  - It starts with `@import "tailwindcss" source(none);` plus explicit `@source` for `**/*.templ` and `**/*.go`.
- **Palette sites** import `palette` CSS output, which defines `:root` shades and `@theme` tokens.
  - Generate from one hex seed.
  - Shared components use semantic utilities only.
- **Single-brand sites** without generated palettes use brand-named `@theme` tokens.
  - Semantic names would only add indirection there.
- **Each colour has one definition.** JS reads tokens with `getComputedStyle` instead of repeating hex values.
- **Classes built in Go are complete literals** chosen by map or branch, so Tailwind's scan finds them.
  - Class lists use templ's list form, never `+` concatenation.

```templ
var sizes = map[Size]string{
	Small: "px-2 py-1 text-sm",
	Large: "px-4 py-2 text-base",
}

templ Button(size Size, active bool, attrs templ.Attributes) {
	<button class={ "rounded-md", sizes[size], templ.KV("ring-2", active) } { attrs... }>
		{ children... }
	</button>
}
```

- **No safelist:** data-driven colours render from the value itself, as an SVG `fill` or a `style`.
  - A `style` attribute needs a `Site.CSP` override: `"style-src-attr": {"'unsafe-inline'"}`.
- **No tailwind-merge.** Components expose typed variants instead of accepting override classes.
  - A caller's extra classes must not conflict with the component's own.
- **`@layer components`** is only for single-element patterns or HTML you don't author.
  - Rendered markdown is one example of HTML you don't author.
- **Use `@apply` rarely.**
- **Load the typography plugin** only where `prose` is used.

## Routing and Server

- **`net/http.ServeMux` with method patterns** routes a 5–20 page site; no router dependency.
- **Mount hgmx handlers at their full prefix**, without `http.StripPrefix`.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /{$}", home)
mux.HandleFunc("GET /docs", docs)
mux.HandleFunc("GET /docs/{path...}", doc)
mux.HandleFunc("GET /items", items) // ?q= search uses the page route
mux.HandleFunc("GET /items/{name}", item)
mux.Handle("/static/", assets)
mux.Handle("GET /favicon.ico", assets.File("favicon/favicon.ico"))
mux.Handle(scripts.Path, js)
mux.Handle(webapp.Path, app)
```

- **`server.Site{Version, CSP, Allow, InternalError, Telemetry, RPC}.Handler(mux)` owns composition.**
  - Application routes run telemetry, recovery and `http.CrossOriginProtection` before the mux.
  - `/healthz` and the `/_v/` beacon bypass that application chain; the beacon mounts only when telemetry is enabled.
  - Supply non-nil `Telemetry`; use `telemetry.Disabled` to opt out of telemetry.
- **`InternalError` is fixed panic-page HTML, not the 404 page.**
  Empty uses hgmx's default; recovery never renders a template.
  - It sends 500 with `no-store` and an accurate `Content-Length`, with no body for HEAD.
  - The application's mux owns its 404 page.
- **Safe methods never mutate.** Cross-origin protection guards the application routes.
- **Security headers cover every route**, including `/healthz` and the mounted `/_v/` beacon.
  - `server.Site` sets hgmx's base CSP with site overrides, `nosniff`, referrer and permissions policies, and opener isolation.
- **`Site.Allow []string`** changes named default Permissions-Policy entries from `()` to `(self)`.
  - Unknown names panic when `Handler` is built.
  - Supported names are `accelerometer`, `browsing-topics`, `camera`, `geolocation`, `gyroscope`, `magnetometer`, `microphone`, `payment` and `usb`.
- **Use `server.New` and `server.Run(ctx, srv, hooks...)`** for the server lifecycle.
  - `New` sets header and idle timeouts; read and write timeouts stay unset for uploads and streams.
  - `Run` drains for up to 5 seconds, then closes on drain failure and runs shutdown hooks.
  - All hooks share a context ending 8 seconds after shutdown starts; after bind failure they get the full 8 seconds.
  - After a drain timeout, handlers may outlive `Run`; closing their resources can race them.
- **`telemetry.FromEnv(local)` never returns nil**, even with export disabled or invalid configuration.
  - Log its configuration error through `Logger()`; local logs and access lines still work.
  - Export is OTLP HTTP/protobuf, configured by supported `OTEL_*` environment variables.
  - Pass `p.Shutdown`, then `func(ctx context.Context) error { return pg.Close(ctx, db) }` to `server.Run`.
  - `Shutdown` logs dropped telemetry and returns nil if the final export fails or times out.
- **`server.Site` serves `/healthz` as `ok <version>`** with `no-store`; set the build version via `-X main.version`.

### Serving an App API

Register native clients in `Host.Apps` under their HTTPS API origin; never move an app ID to another host or platform.
Build `api := accounts.API()`, register site services on `api.Server()`, and classify each procedure with `api.Allow` using generated procedure constants.
Use `api.Mount(mux)`, which rejects missing or unmatched access policies; direct Connect mounting bypasses auth's checks.
Mount browser auth with `accounts.Mount(mux)` and set `server.Site.RPC` to `api.Server()` before building `Site.Handler(mux)` for RPC-aware recovery and unknown-service faults.

Authenticated site handlers read the user with `auth.Caller(ctx)`; `Public` and `Anonymous` calls do not resolve a user.
RPC accepts native bearers, never browser cookies, and requires `HGMX-App-ID` and `HGMX-App-Build` even for bootstrap.
Keep `PageVersion` and other HTML-only middleware on browser handlers.

Set `Config.Erase` from `store/` when account deletion needs site-data cleanup, using the supplied transaction while auth holds its locks.
An error rolls back cleanup and account deletion.
This is the transaction-passing exception: both changes must commit together.

### CSP Literal

**Set site overrides in `server.Site.CSP`, a `map[string][]string`.**
Each entry replaces one base directive's sources or adds a directive outside the base.
For a form POST that redirects off-origin, include the destination in the `form-action` override; Chromium checks redirects too.
A nil map keeps the base; CSP cannot be disabled.
Unknown directives and unsafe source tokens panic when `Handler` is built.
Each source token must be non-empty visible ASCII without spaces, semicolons or commas.

The base policy is:

```http
Content-Security-Policy: default-src 'none'; script-src 'self' https://static.cloudflareinsights.com; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; manifest-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'
```

This includes the sources needed by `webapp` and hgmx's telemetry beacon.
Cloudflare's automatically injected beacon reports to the site's own `/cdn-cgi/rum`, so `connect-src 'self'` covers it.
This follows [Cloudflare's FAQ](https://developers.cloudflare.com/web-analytics/faq/#what-do-i-need-to-add-to-my-content-security-policy-csp); it has not been browser-verified here.

**For one response, set `Content-Security-Policy` only with `server.CSP{…}.Header()`.**
`Header()` serializes the base with those overrides and performs the same validation.
The `csp-literal` check rejects direct header values that are not calls to `server.CSP.Header()`.

## Layout

- **`main.go`** holds flags, wiring and the `migrate` dispatch; serving never migrates.
  - It builds the store and `views.New` result, then calls `routes.New(store, views)` for the mux.
- **`routes/routes.go`** holds the URL table in `New(*store.Store, *views.Site) *http.ServeMux`.
  - Handlers in `routes/` parse input, call the store and render a full page.
- **`store/`** holds `Store`, models and queries; `store/migrations/` holds `All []migrate.Step`.
- **`views/`** owns embedded static files and the page shell.
  - `views/pages/<name>/` holds templ and its view-model types, without HTTP handlers.
- **`views/components/<name>/`** holds components copied from hgmx and the site's own.
- **Add `middleware/`** when application routes need shared request handling.
  - Add a root domain package when business rules need a home outside handlers and queries.
  - Don't scaffold empty packages or add a DI framework.
- **Don't commit `vendor/`**, except on existing sites until the deploy system confirms removal is safe.
  - Nix deploys will fetch modules using a pinned per-site `vendorHash` maintained by kiln.

### Imports: Direction

Site-internal imports point in these directions; `main` wires them together.
Each package may also import its own layer, subject to the page restriction below.

| Package                | May import                                                  |
| ---------------------- | ----------------------------------------------------------- |
| `routes`               | `middleware`, `store`, domain packages, `views`, page views |
| `api`                  | `store`, domain packages, generated `proto/...`             |
| `middleware`           | `store`, domain packages                                    |
| A domain package       | Its own domain and `store`                                  |
| `store`                | `store/...`, excluding `store/migrations/...`               |
| `store/migrations`     | `store/migrations/...` only                                 |
| `views` and components | `views/...`, excluding `views/pages/...`                    |
| `views/pages/...`      | `views/...`, including other page views                     |

**Views only render:** they do not import `net/http`, routes, middleware, store or domain packages.
The table governs site-internal edges; standard-library and external imports remain available, except for the `net/http` and gorm restrictions.

### Layout: No Internal

**No `internal/` directories.** A site is a main module nobody imports; `internal/` adds no useful boundary.
Put packages at the root under the names above.

### Layout: Go Version

**Use at least hgmx's Go release**, currently `go 1.27.0` in `go.mod`.
Raise the site's `go` directive when upgrading hgmx requires it.

## Data

### Gorm: Store Only

- **Models and gorm imports stay in `store/`**, including its migrations.
  - `main` passes `pg.Open`'s result to `store.New` and `pg.Close` without naming a gorm type.
- **`store.Store` methods take `ctx` first** and query through `s.db.WithContext(ctx)`.
  - Routes hold `*store.Store` and call its methods.
- **Test outcomes with `errors.Is`**, using `pg.ErrNotFound`, `pg.ErrDuplicate` or `pg.ErrForeignKey`.
- **A Store method owns its transactions**, using `s.db.WithContext(ctx).Transaction`.
  - No generic repository interface, gorm-forwarding methods or transaction-passing machinery.
- **Use gorm's `clause.OnConflict` for upserts.** Keep the conflict target and updates in the Store method.
- **Raw SQL carries a one-line reason** for what the query builder cannot express clearly.
- **Timestamps use `timestamptz` only.** hgmx supplies UTC timestamps.

### Gorm: No Global

**Keep `*gorm.DB` in `store.Store`, not a package variable.** Pass the Store to its callers so ownership and test isolation stay visible.
The check exempts test files.

### Column Types and Queries

- **`pg.ID` is a UUIDv7** when made with `pg.NewID`; parse external IDs with `pg.ParseID`.
  - A zero ID saves as SQL NULL; scanning NULL gives the zero ID.
- **`pg.JSON[T]` stores `V` as `jsonb`** and marshals as `V`, without a wrapper object.
  - SQL NULL scans to the zero value of `T`; saving writes non-NULL JSON, even for a nil value.
  - Use it on `NOT NULL DEFAULT` columns; it does not preserve SQL NULL.
- **`pg.Enum[T]`** validates string values and creates or extends a Postgres enum in a migration.
  - Keep each migration's values local to that migration, and use a new value in a later step.
- **`pg.Page.Scope`** clamps a negative offset to zero, defaults a non-positive limit to 20 and caps it at 100.
  - `pg.List[T]` counts the supplied query, then fetches the requested page.
- **`pg.Like`** trims the query and matches a case-insensitive substring, escaping `%`, `_` and `\`.
  - Empty or whitespace-only input leaves the query unchanged.
- **Use `db.Where("x IN ?", vs)` for membership.** An empty or nil slice produces `IN (NULL)` and matches no rows.
- **`pg.Sort`** uses an allowlist of keys to columns; an unknown key adds no order.

## Migrations

### Gorm: No AutoMigrate

**Site schema changes are `migrate.Step` entries in `store/migrations.All`.** Never call `AutoMigrate` or migrate at server startup.
Versions must be exactly **1, 2, 3, … in list order within each track**; append steps without deleting or reordering applied ones.

Pass `migrate.Track{Steps: migrations.All}` as the unnamed site track in the `[]migrate.Track` passed to `migrate.Command`, after any named library tracks.
The site track keeps `goose_db_version`; named tracks use `hgmx_<name>_migrations`.

- **`Up` is required.** A nil `Down` makes rollback fail and leaves the history row in place.
- **Each step runs in a transaction** with a 5-second `lock_timeout`, unless `NoTx` is set.
- **`NoTx` steps must be rerunnable:** the operation and its history row are not atomic.
  - Keep one concurrent statement per step, plus the cleanup needed to retry it.
  - For an index, run `DROP INDEX CONCURRENTLY IF EXISTS`, then `CREATE INDEX CONCURRENTLY`.
  - Use separate `Exec` calls; a multi-statement call cannot run these outside a transaction.
- **`migrate.Command`** accepts `up`, `down [track]` and `status`.
  - `up` runs tracks in slice order; `status` lists every track.
  - `down` rolls back one site step; `down auth` rolls back one step in the named `auth` track.
  - With steps registered, all commands require `PGDATABASE` and ping the database.
  - With steps registered, one session lock covers the whole command, including `status`.
  - `status` can wait behind a deploy; it is not a lock-free probe.
  - The pool needs at least two open connections, or three when running `NoTx` steps.

### Migrate: No Live Models

**Migrations never import live models**, directly or through another package.
Declare the columns a step needs inside that step, or write SQL.
_A fresh database must replay each step against the schema of its own time._

### Deploy Contract

**Kiln runs migrations only during deploys**, using the new binary's `<site> migrate up` in a transient unit with the site's user and `PG*` environment.
A non-zero exit, including a lock timeout, keeps the old binary version running.
Production never runs `down`; rollback reverts binaries only.

**While auth schema v5 is live, rollback to v0.3.0 or older weakens credential isolation.**
Those binaries accept a copied native app token as a browser cookie on the same origin because they do not check `client_id`.
Recover by deploying forward.

**Expand, then contract:** add schema compatible with the old binary, deploy code that uses it, and remove obsolete schema in a later deploy after old binaries no longer need it.
An applied migration can remain after a later step fails, so every intermediate schema must support the old binary.

**After Google identities exist, use the current CLI for account recovery.**
An older binary's `auth reset` leaves linked Google accounts in place; deploy forward first if the current CLI is not available.

**Auth v4 cutover requires a fresh sign-in.**
Sessions without an origin are refused, so everyone signs in once more.
Pending codes, links, Google flows and passkey ceremonies from the previous binary must restart.

**Applying auth v4 under a running previous binary** changes cached `SELECT *` and `RETURNING *` result types.
Each affected cached statement fails once per pooled connection with `cached plan must not change result type`; retrying that statement on the same connection succeeds.
A retry can hit another connection's stale plan first.

**Rollback past the multi-host release weakens isolation.**
The previous binary accepts new sessions without host binding or `Host.Allow`, and its sign-out-others reaches sibling hosts.
Recover by deploying forward, not by staying rolled back; sessions created during rollback have no origin and are refused again on roll-forward.

## Commands

### hgmx: Latest

**Every site stays on the latest hgmx master while hgmx is being built out.**
`go tool hgmx check` compares the commits in `go.mod` and existing `vendor/modules.txt` with master from one `git ls-remote` lookup of master and tags, using peeled commits for annotated tags.
If resolution fails, including offline, missing git or after a 10-second timeout, it warns with the cause and skips this rule.

Run `go tool hgmx update` to resolve master with git and run `go get -tool github.com/cogikyo/hgmx/cmd/hgmx@<12-char commit>`, then `go mod tidy`.
It runs `go mod vendor` only when `vendor/` already exists.
Update fails if master cannot be resolved or its module version is older than the site's current requirement.
It prints the old → new hgmx version and suggests `go tool hgmx gen` before `go tool hgmx check` if templ changed.

### hgmx: Tool

**Declare `tool github.com/cogikyo/hgmx/cmd/hgmx` in the site's `go.mod` and use `go tool hgmx`.**
The library and CLI then use the same module version and cannot drift.
For a first install, run `go get -tool github.com/cogikyo/hgmx/cmd/hgmx@master`.

### hgmx: Skills

**Each site commits the agent skills of its pinned hgmx version, unedited.**
`go tool hgmx skills` writes them to `.agents/skills/<name>/` and records a SHA-256 of each written file in `.agents/skills/hgmx.sum`, in `sha256sum` format.
When `.claude/` exists, it also links `.claude/skills/<name>` to the skill folder, relative to the physical `.claude/skills/` folder.
`go tool hgmx update` runs it after the module bump.

Put site-specific guidance in `SITE.md` beside a skill; hgmx never writes or removes it.
Edited managed files or conflicting unowned files stop the sync before skill files are changed.
The message tells you to move changes into `SITE.md`, delete the file, then rerun `go tool hgmx skills`.
Sync re-hashes skill files before replacing or removing them; a change found then stops the sync, but earlier writes remain.
Writes use sibling `.<name>.hgmx-<random>` temp files before renaming.
A skill hgmx no longer ships loses its unedited files and its managed link; its folder stays while `SITE.md` remains.
The `hgmx-skills` rule in `go tool hgmx check` fails on missing, stale, edited or dropped managed files, incorrect links, or stale fingerprints.

### Layout: No Justfile

**No justfile anywhere.** Run the shared hgmx commands from the site's module root.

| Command                             | Task                                                                        |
| ----------------------------------- | --------------------------------------------------------------------------- |
| `go tool hgmx init`                 | Write `main.go`, `routes/`, `store/` and `views/` in a module with `go.mod` |
| `go tool hgmx gen`                  | Generate templ and minified CSS                                             |
| `go tool hgmx dev [-- site flags…]` | Watch templ and CSS, restarting `go run . -static views/static`             |
| `go tool hgmx fmt`                  | Format Go and templ files                                                   |
| `go tool hgmx migrate …`            | Run `go run . migrate …` for development                                    |
| `go tool hgmx credentials`          | Build the site and write `credentials.json`                                 |
| `go tool hgmx credentials check …`  | Test SMTP and Google client authentication                                  |
| `go tool hgmx check`                | Check generated files, credentials, rules, format, vet, build, test         |

`init` refuses if any scaffold target already exists and pins the site's templ requirement to hgmx's.
It rewrites site imports and regenerates templ Go files for the site's module path.
Tailwind generation uses `tailwindcss` from PATH when `views/static/css/main.css` exists.
`check` compares generated output to the tree in its gen stage without writing.
It stops at the first failed stage and never writes formatted files.
For package loading and Go commands, `check` sets `GOWORK=off` and passes `-mod=vendor` when `vendor/` exists or `-mod=readonly` otherwise, overriding `GOFLAGS`.
`dev` and `migrate` keep the caller's workspace settings.

### hgmx: Credentials

**Declare credentials unconditionally in a site-owned `credential.Registry`.**
Use `Feature(name, condition)` for requirements that follow runtime flags or env, or `credential.Always` for unconditional requirements.
Keep those conditions consistent with the site's runtime branches; Go serializes them for Nix without evaluating them.
`Declare(definition, requirement, shared)` copies the metadata and returns a handle; read values through `handle.Read()`.
Use `"site"` for site-local secrets or a shared label such as `"smtp"`.

Register server flags before dispatching the site's `credentials` subcommand to `credential.Command`, before opening databases, queues or telemetry.
Pass that same FlagSet and a `credential.Checks` map keyed by each SMTP/Google handle's `Name()`, with callbacks using the site's provider settings.
The scaffold includes the dispatch and an empty registry for sites with no secrets.

Run **`go tool hgmx credentials`** after changing declarations or registered flags, then commit the site-root `credentials.json` with the code, even for zero-secret sites.
It builds the site and exports metadata without secret values; cogictl uses the file to select deployment secrets.
The `credentials` stage of `go tool hgmx check` compares that output byte-for-byte without writing; a missing file or any difference, including formatting, fails.
Invalid declarations, unknown referenced flags and site build failures also fail the stage.

Use `go tool hgmx credentials check [name...] [-- site flags...]` for live SMTP/Google client checks with a current manifest.
With no names, it checks all declared SMTP/Google targets; missing values or settings fail instead of skipping them.
It never prints secrets; SMTP sends no mail, and Google success verifies token-endpoint client authentication, not browser sign-in.
Ordinary `go tool hgmx check` does not run live provider checks.

## Env

- **`PGHOST`, `PGDATABASE`, `PGUSER`** configure Postgres through libpq-compatible `PG*` settings.
  - hgmx defines no database URL variable.
- **`pg.Open` opens a lazy pool**, so serving can start without a reachable database.
  - The first query establishes connectivity.
  - Close it on shutdown with `pg.Close`.
  - SQL logs retain placeholders and never include bound parameter values.
- **`OTEL_*`** configures telemetry; export uses HTTP/protobuf.
  - Set `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` and `OTEL_RESOURCE_ATTRIBUTES` as needed.
  - An unset endpoint or `OTEL_SDK_DISABLED=true` disables export; local logging still works.
- **`HGMX_TRUSTED_PROXY`** is the proxy allowed to supply `CF-Connecting-IP`: an IP or CIDR; invalid values stop startup.
  - `clientip.From` reads it once at startup; auth's per-IP limits use that address.
  - When unset, the connection's peer address is used instead.
- **`HGMX_AUTH_KEY`** is the fallback for a handle declared from `auth.KeyCredential`.
  - `$CREDENTIALS_DIRECTORY/auth-key` takes priority; use a per-site key of at least 32 bytes.
- **`smtp-password`** is read through a handle declared from `mail.PasswordCredential`.
  - `$CREDENTIALS_DIRECTORY/smtp-password` takes priority; if missing, `HGMX_SMTP_PASSWORD` is read instead.
- **`google-client-secret`** is read through a handle declared from `auth.GoogleCredential` and passed to `auth.Google.ClientSecret`.
  - `auth.New` requires both `ClientID` and `ClientSecret` when `Config.Google` is set; it does not read credentials.
  - `$CREDENTIALS_DIRECTORY/google-client-secret` takes priority; if missing, `HGMX_GOOGLE_CLIENT_SECRET` is read instead.
- **`HGMX_SMTP_HOST`, `HGMX_SMTP_USER`, `HGMX_SMTP_FROM`** configure the site's `mail.SMTP` sender.
  - Use port 465, the default; `mail.SMTP` uses implicit TLS, not STARTTLS.
  - Cloudflare Email Sending uses host `smtp.mx.cloudflare.net` and user `api_token`.
  - One account-scoped token sends for every onboarded domain.
  - Set `HGMX_SMTP_FROM` to an address on the site's `mail.<domain>`.

## Tests

**The scaffold's `routes.TestNew` builds the mux.**
A missing file passed to `static.Server.File` or a pattern conflict fails `go tool hgmx check`.

**Use `pgtest.Open(t)` for database tests.** It supplies a fresh database and non-superuser role per call, then cleans both up with the test.
Production and test connections share hgmx's gorm configuration.

- **Put `initdb` and `pg_ctl` on PATH.** Missing executables or a socket bind failure skip the test.
  - Set `HGMX_PGTEST=require` to fail instead; CI requires it.
  - All other startup failures fail the test.
- **The cluster is shared per user**, under `$TMPDIR/hgmx-pgtest-<uid>` (`/tmp` if `TMPDIR` is unset).
  - Startup is file-locked; tests reuse it across binaries and leave it running.
  - It listens only on a Unix socket, whose path must fit within 107 bytes.
  - Set `TMPDIR` to a shorter path if needed.
- **Each test binary holds a `flock` lease for its lifetime**, named by a random owner ID.
  - Its databases and roles carry that ID; ownership works across process namespaces.
  - Another binary reclaims them and the lease file only after taking the lease without blocking.
  - Resources without a matching owner lease are left untouched.

## Logging and Exit

### slog: Context

**Use slog's `DebugContext`, `InfoContext`, `WarnContext` and `ErrorContext` when a context or request is in scope.** Pass `ctx` or `r.Context()` so the record carries the request's trace and span IDs.

### Exit: Main Only

**Only `main` may call `os.Exit` or `log.Fatal*`; `TestMain` is also allowed.** Return errors elsewhere so deferred cleanup and graceful shutdown can run.

### hgmx: Ignore

**An exception names its rule and reason:** put `//hgmx:ignore <rule> <reason>` on the finding's line or the line above.
Go-rule findings from generated templ code map to the `.templ` source line; put the ignore there, not in `_templ.go`.
A missing reason, unknown rule or unused ignore is a finding; remove the directive when its finding goes away.

## Components

Copy styled components; import infrastructure.

- **Copied and owned by the site:** templ markup, variant-to-class maps and the component's behaviour JS.
  - Copies use semantic tokens, so restyling rarely needs edits.
- **Imported as Go modules:** infrastructure, and any small helper a copied component calls.
  - Copied helpers drift out of step with the module, so helpers stay imported.
- **A component moves into hgmx** only when a second site needs it.

The infrastructure packages are:

- **`server`** owns HTTP composition, recovery, version reloads and shutdown.
- **`render`** buffers templ HTML with cache validation.
- **`telemetry`** owns local logging, OTLP export and browser beacons.
- **`pg`** owns Postgres connections, column types and query scopes; `pg/pgtest` supplies test databases.
- **`migrate`** runs Go-defined schema steps under a session lock.
- **`static`** serves the site's files at set-hash URLs or fixed root URLs.
- **`scripts`** serves selected scripts and navigation support with ordered head tags.
- **`webapp`** serves the install manifest and icons.
- **`palette`** generates OKLCH shades and semantic colour tokens from one seed.
- **`icons`** provides SVG templ components and hover-animation assets.

Put the selected handlers' `Head()` output in the document head.
It returns `template.HTML`, which html/template renders as-is and templ takes as `@templ.Raw(...)`.
The icons handler's `Head()` belongs after the `scripts.Motion` tag; the icon hover helper also needs hyperscript.
