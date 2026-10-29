# hgmx

hgmx is my shared Go toolkit.
An abstraction lands here once every app needs it.
Each site imports only what it uses.
[CONVENTIONS.md](CONVENTIONS.md) sets how a site uses these packages, htmx, templ and Tailwind.

## Build and Run

With Go 1.26 or later, from this repo:

```bash
go build -o hgmx ./cmd/hgmx
./hgmx --help
```

To install instead, use `go install github.com/cogikyo/hgmx/cmd/hgmx@latest`.

## CLI Usage

- `hgmx init` copies a working two-page scaffold from the embedded `library/` into `./views/`.
  - It copies `views.go`, `views.templ`, `icon.png` and `static/`.
  - Before the first build, run `go get github.com/cogikyo/hgmx` to require the hgmx packages.
  - It copies no generated code, so pin templ once with `go get -tool github.com/a-h/templ/cmd/templ@v0.3.1100`.
  - Then run `go tool templ generate`, and build.
  - `views.New(dev)` builds the `Site`; `Site.Scripts` loads htmx, hx-preload and _hyperscript.
  - `Site.Mount` registers the static, scripts and webapp handlers and the pages `GET /{$}` and `GET /about`.
  - It also registers `GET /` as a full-page 404 through the layout.
- `hgmx palette '#ecc45d'` overwrites `library/static/css/colors.css` in this repo.
- `hgmx palette --lua '#ecc45d'` writes a Lua palette to stdout.

The palette generator is also available as `github.com/cogikyo/hgmx/palette`.

## Static Files

`github.com/cogikyo/hgmx/static` serves an `fs.FS` at `<prefix><sethash>/<path>`, where one hash covers every file.
Relative CSS and ES module references that stay within the file set inherit the hash.

```go
sub, err := fs.Sub(embedded, "static")
if err != nil {
	return err
}
assets, err := static.New("/static/", sub)
if err != nil {
	return err
}
mux.Handle("/static/", assets)

css, err := assets.URL("css/main.min.css")
```

- **Caching:** a stale hash serves the current file with `no-cache`, not a 404.
  - Current snapshot hashes get immutable caching; existing files carry ETags.
- **Dev:** `static.Dev("/static/", os.DirFS("views/static"))` reads files on each request with `no-cache`.
- **Templates:** `URL` returns `(string, error)`, so templ takes `href={ assets.URL("css/main.min.css") }` and html/template takes it as a function.

## Scripts

`github.com/cogikyo/hgmx/scripts` embeds htmx, the htmx `hx-preload` extension, _hyperscript, and Motion, pinned with sha256 in `scripts/lock.json`.

```go
js, err := scripts.New(scripts.Preload, scripts.Hyperscript)
if err != nil {
	return err
}
mux.Handle(scripts.Path, js)
```

`Preload` includes `HTMX` automatically.
`Head()` emits deferred scripts in `Library` declaration order, regardless of argument order.

`go generate ./scripts` downloads files that are missing or that don't match the lock.
`go run ./scripts/fetch -check scripts` only verifies the vendored files; the release workflow runs it before publishing.

## Installable Web App

`github.com/cogikyo/hgmx/webapp` builds a manifest and icons for browser installation, without adding offline support.
It generates the icons from one square PNG of at least 512 px per side and serves the set under `/app/`, apart from the site's own static files.

```go
w, err := webapp.New(webapp.App{
	Name:       "hgmx",
	ShortName:  "hgmx",
	Theme:      "#222536",
	Background: "#222536",
	Icon:       iconPNG,
})
if err != nil {
	return err
}
mux.Handle(webapp.Path, w)
```

- **Colours:** `Theme` and `Background` must be `#RRGGBB`.
  - `Background` fills the maskable icons and the 180 px apple-touch-icon.
- **Identity:** the manifest URL changes when the manifest or icon set changes.
  - The manifest's fixed `id: "/"` preserves the installed app's identity.
- **Location:** `start_url` and `scope` are fixed at `/`.
