# hgmx

hgmx is my shared Go toolkit for server-rendered sites using templ, htmx and Tailwind.
Each site imports only what it uses.
[CONVENTIONS.md](CONVENTIONS.md) holds the site-building rules; godoc holds the API details.

## Start a Site

Use Go 1.26 or later, in a new project directory:

```bash
go install github.com/cogikyo/hgmx/cmd/hgmx@latest
go mod init example.com/site
hgmx init
go get github.com/cogikyo/hgmx
go mod tidy
go mod vendor
go build -o site .
./site
```

`init` writes `./main.go` and `./views/`, including generated templ code and static assets.
It refuses to overwrite an existing `main.go` or `views/`.
The scaffold serves Home, About and a full-page 404, with system-ui and OFL-licensed Iosevka fonts.
Pass `-static views/static` to read assets from disk while editing.

The generated `main.go` starts with `telemetry.FromEnv`, which returns a usable pipeline even when export is disabled or configuration fails.
It logs configuration errors locally, then composes the site in its `run` function:

```go
func run(p *telemetry.Pipeline, addr, dir string) error {
	site, err := views.New(views.Config{Version: version, Static: dir, Head: p.Beacon().Head()})
	if err != nil {
		return err
	}
	h := server.Site{
		Version:   version,
		CSP:       csp,
		Fallback:  http.HandlerFunc(views.InternalError),
		Telemetry: p,
	}.Handler(site)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)
	log := p.Logger()
	log.Info("listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, h, log), p.Shutdown)
}
```

The imports, `version` and `csp` declarations are in that generated file.
`server.Site` owns the middleware, health check and beacon; `views.New` returns the application's mux.

## Packages

- **`server`** composes the root handler, version reloads, recovery and graceful shutdown.
- **`render`** buffers templ HTML with ETags and conditional requests.
- **`telemetry`** exports OTLP/HTTP traces, logs and metrics from `OTEL_*` settings, with local logging and a browser beacon.
- **`static`** serves file sets under hashed URLs, with a disk-backed dev mode.
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
