# hgmx

A CLI for copying hgmx's customizable templates into Go HDA applications, inspired by [shadcn-ui](https://ui.shadcn.com/).

Utilizes:

- [HTMX](https://htmx.org/) &rarr; core hypermedia-driven application (HDA) driver
- [_hyperscript](https://hyperscript.org/) &rarr; simple client side scripting
- [Go Templ](https://templ.guide/) &rarr; Go HTML template builder
- [Tailwind](https://tailwindcss.com/) &rarr; styling
- [Motion](https://motion.dev/) &rarr; animation library

## Build and Run

With Go 1.27 or later, from this repo:

```bash
go build -o hgmx ./cmd/hgmx
./hgmx --help
```

To install instead, use `go install github.com/cogikyo/hgmx/cmd/hgmx@latest`.

## CLI Usage

- `hgmx init` copies templates from the embedded `library/` into `./views/`.
- `hgmx palette '#ecc45d'` overwrites `library/static/css/colors.css` in this repo.
- `hgmx palette --lua '#ecc45d'` writes a Lua palette to stdout.
- `hgmx link --input '../hgmx/library/*' --output 'views/*'` replaces destination files with symlinks.

The palette generator is also available as `github.com/cogikyo/hgmx/palette`.
