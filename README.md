# hgmx

A customizable component-based framework for building HDA applications in Go, inspired by [shadcn-ui](https://ui.shadcn.com/). 

Utilizes:

- [HTMX](https://htmx.org/) &rarr; core hypermedia-driven application (HDA) driver
- [_hyperscript](https://hyperscript.org/) &rarr; simple client side scripting 
- [Go Templ](https://templ.guide/) &rarr; Go HTML template builder 
- [Tailwind](https://tailwindcss.com/) &rarr; styling
- [Motion](https://motion.dev/) &rarr; animation library


## Installation

```bash
go install github.com/nosvagor/hgmx/cmd/hgmx@latest
```

## CLI Usage

Initialize a new project:

```bash
hgmx init
```

Symlink components to another project (useful for forking own version)

```bash
hgmx link .
```