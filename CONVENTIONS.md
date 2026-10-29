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
It relies on the stale-response extension described there.

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
::view-transition-old(root),
::view-transition-new(root) {
  animation: none;
}
::view-transition-group(main),
::view-transition-old(main),
::view-transition-new(main) {
  animation-duration: 200ms;
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
- **Per-page regions outside `#main`**, such as a breadcrumb, are listed in `hx-select-oob` and morphed.
  - Morph only regions whose attributes only the server writes.
  - Morph strips classes and inline styles that JS set.
- **Persistent elements outside `#main`** keep their DOM node and JS state across clicks, Back and Forward.
  - Examples are a nav controller, Motion state, or a background canvas.
- **Error and 404 pages are full pages.** htmx 4 swaps 4xx and 5xx bodies, so it selects their `#main`.
- **Links to non-page responses** get `hx-boost="false"` or `download`; otherwise `#main` ends up empty.
  - These include feeds, PDFs, `/healthz`, downloads and pages without `#main`.
- **Inheritance is explicit in htmx 4:** add `:inherited` on purpose and expect nothing implicit.

**Drop responses from elements that left the page.**
In htmx 4.0.0, navigation and an element's own requests run in separate queues, and history updates before the swap.
A late response from an element that navigation removed would still overwrite the URL and title.
Register this extension from a static module (see [Scripting](#scripting)):

```js
htmx.registerExtension("stale", {
  htmx_after_request: (elt) => elt.isConnected,
});
```

It drops the response once its body has arrived, before redirects, history, title or swap.
Abort-based `hx-sync` fixes log an htmx error per abort or slow navigation down.

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
Text typed into an input in `#main` during a navigation is lost, though focus is kept.

The htmx 4.0.0 source differs from its docs.
Boosted swaps default to `config.defaultSwap` (`innerHTML`), not `outerSync`.
Restores always use `outerSync` and send no `HX-Request` or `HX-Target`.
Preload requests never carry `HX-Request-Type`.

**Fragments are a fallback** for a site that truly needs them: use explicit `hx-get` with `hx-preload`, branch on `HX-Target: main#main` with `Vary: HX-Target`.
Never branch on `HX-Request-Type`; preload requests lack it, so the full page would land inside `#main`.

## Prefetch and Transitions

- **htmx sites use hx-preload** (`scripts.Preload`).
  - It preloads boosted anchors on mousedown and touchstart and reuses the response for 5 s.
  - Plain links need no preload attributes.
- **Don't add Speculation Rules** to links hx-preload covers: they fetch twice, with no documented shared cache.
- **Sites without htmx add no prefetch** while they have only a page or two.
  - When such a site grows pages, it sends a `Speculation-Rules` response header.
  - Inline rules would need `'inline-speculation-rules'` in the CSP.
  - It also adds `@view-transition { navigation: auto; }`.
- **Only `#main` animates**, and the root animation is off.
  - `#main` carries `view-transition-name`, so persistent chrome doesn't flash.
- **Transitions take 200 ms or less**, and reduced motion turns them off.
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
- **Rendering is buffered** (the default `templ.Handler`) unless a measured need calls for streaming.
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

Sites use Go 1.26 or later.

- **`net/http.ServeMux` with method patterns** routes a 5–20 page site; no router dependency.
- **Mount hgmx handlers at their full prefix**, without `http.StripPrefix`.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /{$}", home)
mux.HandleFunc("GET /docs", docs)
mux.HandleFunc("GET /docs/{path...}", doc)
mux.HandleFunc("GET /items", items) // ?q= search uses the page route
mux.HandleFunc("GET /items/{name}", item)
mux.HandleFunc("GET /healthz", healthz)
mux.Handle("/static/", assets)
mux.Handle(scripts.Path, js)
mux.Handle(webapp.Path, app)
```

- **Middleware is `func(http.Handler) http.Handler`**, applied in one visible chain.
  - The order from outside in: recover, log, security headers, `CrossOriginProtection`, mux.
- **`http.NewCrossOriginProtection()`** guards any site with forms or mutations; safe methods never mutate.
- **Every site sends security headers:** its own CSP, `nosniff`, `Referrer-Policy` and `Permissions-Policy`.
- **CSP starts from `default-src 'none'`** and allows `script-src 'self'` for scripts.
  - Sites using `webapp` add `manifest-src 'self'` and `img-src 'self'`; these were checked in Chromium.
- **Set `ReadHeaderTimeout` explicitly.**
- **On a signal, call `Shutdown` with a deadline**, and exit only after `Shutdown` returns.
  - `Serve` returns `ErrServerClosed` at once, while active requests are still draining.
- **`/healthz` returns `ok <shortrev>`** with `no-store`; the version comes from `-X main.version`.

## Layout

- **`cmd/main.go`** holds flags and wiring, and `internal/server` holds the server setup.
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

- **`static`** serves the site's files under a set-hash URL.
- **`scripts`** serves the vendored scripts and their head tags.
- **`webapp`** serves the install manifest and icons.
- **`palette`** generates the colour palette.

Put `scripts.Handler.Head()` and `webapp.Handler.Head()` in the document head.
Both return `template.HTML`, which html/template renders as-is and templ takes as `@templ.Raw(...)`.
