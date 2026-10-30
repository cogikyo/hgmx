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
- **No JSON API** feeds client rendering, and no route negotiates between JSON and HTML.

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

```html
<body
  hx-boost:inherited="target:#main select:#main swap:'innerHTML strip show:top showTarget:body' transition:true"
  hx-select-oob:inherited="#nav-links:outerMorph"
  hx-config:inherited='{"headers":{"HGMX-Version":"<rev>"}}'
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
@media (prefers-reduced-motion: reduce) {
  ::view-transition-group(*),
  ::view-transition-old(*),
  ::view-transition-new(*) {
    animation: none;
  }
}
```

**Use it this way:**

- **Links are plain `<a href>`**, so they work without JS and get boosted.
  - Use explicit `hx-get` and `hx-target` only where a link behaves differently from navigation.
  - Don't put `hx-get` with `hx-push-url` on nav links.
- **Handlers always return the full page:** no `HX-Request` branch, no `Vary`, no partial template.
- **The body takes no `hx-sync`** when loading `scripts.Navigation`.
- **Version reloads pair `server.VersionConfig` with `server.PageVersion`.**
  - Emit `server.VersionConfig(version)` on `<body>` and wrap page handlers with `PageVersion`.
  - A GET or HEAD carrying a different build version gets a full reload; requests without the header pass through.
  - `scripts.Navigation` re-sends the config headers on history restores.
- **Per-page regions outside `#main`**, such as a breadcrumb, are listed in `hx-select-oob` and morphed.
  - Morph only regions whose attributes only the server writes.
  - Morph strips classes and inline styles that JS set.
- **Persistent elements outside `#main`** keep their DOM node and JS state across clicks, Back and Forward.
  - Examples are a nav controller, Motion state, or a background canvas.
- **Validation and 404 pages use the full layout.** htmx 4 selects their `#main` as on success.
  - Recovery serves the fixed `InternalError` page for a 500; it has no nav.
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
  - `#main` carries `view-transition-name`, so persistent chrome doesn't flash.
- **Page transitions take 120 ms**, and reduced motion turns them off.
- **Function never depends on animation.**
- **Motion handles choreographed component animation**, such as a logo, menu or icons.
  - View transitions handle page-to-page changes; the two don't overlap.

## Scripting

- **hyperscript** covers small, local behaviour next to its markup.
  - Examples are toggles, dismissal, event relays (`send hovered to …`), timed classes and event waits.
  - Repeated logic becomes a `behavior` in a `._hs` file, attached with `install`.
- **JS modules** cover widgets with real logic, such as a scroll controller or Motion choreography.
  - They live in `.js` files served through `static`, whose set-hash URLs keep relative imports working.
- **No inline scripts:** no inline modules, classic scripts, importmaps or `templ.Raw("<script>…")`.
  - Inline scripts rerun on every swap; external modules evaluate once and cache.
  - CSP then needs only `script-src 'self'`, without nonces.
- **Module JS is idempotent across swaps.**
  - Delegate events on `document`, or initialise from `htmx:after:process`.
  - Don't guard with `data-*` "ready" markers; morph removes attributes the server didn't send.
- **One runtime per responsibility:** no `hx-live` or Alpine next to hyperscript.
- **No `hx-on`:** it evaluates code, which strict CSP blocks.
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
- **The templ version in `go.mod` is the one authority.** Sites require templ's runtime as usual.
  - Run the generator at that version: `go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate`.
- **No templ `tool` line:** Go would vendor the generator's whole dependency tree.
  - It grew a site's `vendor/` from about 0.8 MB to 16 MB.
  - hgmx keeps its own `tool` line because it has no `vendor/`.
- **Generated templ and CSS are committed**, and regenerated and checked before each commit.
  - templ checks its output with `generate -check` and `fmt -fail .`, run at the same version.

## Tailwind v4

- **Standalone CLI 4.3.3** with one input, `main.css`.
  - It starts with `@import "tailwindcss" source(none);` plus explicit `@source` for `**/*.templ` and `**/*.go`.
- **Palette sites** map `palette` output to `:root` variables, then to `@theme inline` semantic tokens.
  - Generate from a base and primary colour, with an optional accent, and inspect `Palette.Clashes()`.
  - Tier labels express contrast targets, not fixed lightness; unreachable tiers stay in the output.
  - Shared components use semantic utilities only.
- **Single-brand sites** without palette swapping use brand-named `@theme` tokens.
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
  - A `style` attribute also needs a CSP that allows inline styles.
- **No tailwind-merge.** Components expose typed variants instead of accepting override classes.
  - A caller's extra classes must not conflict with the component's own.
- **`@layer components`** is only for single-element patterns or HTML you don't author.
  - Rendered markdown is one example of HTML you don't author.
- **Use `@apply` rarely.**
- **Load the typography plugin** only where `prose` is used.

## Routing and Server

Sites use Go 1.27 or later.

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
mux.Handle(scripts.Path, js)
mux.Handle(webapp.Path, app)
```

- **`server.Site{Version, CSP, Fallback, Telemetry}.Handler(mux)` owns composition.**
  - Health and the `/_v/` beacon sit outside the middleware chain.
  - The chain is telemetry, security headers, recovery, `http.CrossOriginProtection`, then an unwrapped `*http.ServeMux`.
  - Supply a non-nil telemetry pipeline and a fixed fallback that renders without calling failing templates.
- **Safe methods never mutate.** Cross-origin protection guards the application routes.
- **`server.Site` sends security headers:** the site's CSP, `nosniff`, referrer and permissions policies, and opener isolation.
- **CSP starts from `default-src 'none'`** and allows `script-src 'self'` for scripts.
  - Sites using `webapp` add `manifest-src 'self'` and `img-src 'self'`; these were checked in Chromium.
  - The telemetry beacon needs `script-src 'self'` and `connect-src 'self'`.
- **Use `server.New` and `server.Run(ctx, srv, hooks...)`** for the server lifecycle.
  - `New` sets header and idle timeouts; read and write timeouts stay unset for uploads and streams.
  - `Run` drains for up to 10 seconds, then closes on drain failure and runs shutdown hooks.
  - All hooks share a fresh 5-second context; they run even after bind failure.
  - After a drain timeout, handlers may outlive `Run`; closing their resources can race them.
- **`telemetry.FromEnv(local)` never returns nil**, even with export disabled or invalid configuration.
  - Log its configuration error through `Logger()`; local logs and access lines still work.
  - Export is OTLP HTTP/protobuf, configured by supported `OTEL_*` environment variables.
  - Pass `p.Shutdown` to `server.Run` so export flushes after the HTTP drain.
- **`server.Site` serves `/healthz` as `ok <version>`** with `no-store`; set the build version via `-X main.version`.

## Layout

- **`./main.go`** holds flags and wiring; `hgmx init` also writes `./views/`.
- **Feature folders:** a page's handler, view model and templ live together in `views/pages/<name>/`.
  - Don't split a page between a handlers tree and a views tree.
- **`views/components/<name>/`** holds components copied from hgmx and the site's own.
- **No service or store layers** without a database, and no DI framework.
- **Commit `vendor/`:** deployed sites build offline from it.

## Components

Copy styled components; import infrastructure.

- **Copied and owned by the site:** templ markup, variant-to-class maps and the component's behaviour JS.
  - Copies use semantic tokens, so restyling rarely needs edits.
- **Imported as Go modules:** infrastructure, and any small helper a copied component calls.
  - Copied helpers drift out of step with the module, so helpers stay imported.
  - hgmx only adds to helpers, so a `go.mod` bump never breaks an untouched copy.
- **A component moves into hgmx** only when a second site needs it.

The infrastructure packages are:

- **`server`** owns HTTP composition, recovery, version reloads and shutdown.
- **`render`** buffers templ HTML with cache validation.
- **`telemetry`** owns local logging, OTLP export and browser beacons.
- **`static`** serves the site's files under a set-hash URL.
- **`scripts`** serves selected scripts and navigation support with ordered head tags.
- **`webapp`** serves the install manifest and icons.
- **`palette`** generates contrast-tier scales and semantic colour tokens.
- **`theme`** serves palette stylesheets, switching and the explorer.
- **`icons`** provides SVG templ components and hover-animation assets.

Put the selected handlers' `Head()` output in the document head.
It returns `template.HTML`, which html/template renders as-is and templ takes as `@templ.Raw(...)`.
The icons handler's `Head()` belongs after the `scripts.Motion` tag; the icon hover helper also needs hyperscript.
