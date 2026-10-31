package views

import (
	"embed"
	"html/template"
	"io/fs"
	"os"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/scripts"
	"github.com/cogikyo/hgmx/static"
)

//go:embed static
var files embed.FS

const Assets = "/static/"

type Config struct {
	Version string
	Static  string        // Development directory; empty uses embedded assets.
	Head    template.HTML // Trusted HTML appended to the document head.
}

type Site struct {
	Static  *static.Server
	Scripts *scripts.Handler
	head    template.HTML
	version string
}

func New(c Config) (*Site, error) {
	s := &Site{Scripts: scripts.New(scripts.Navigation), head: c.Head, version: c.Version}
	if c.Static != "" {
		s.Static = static.Dev(Assets, os.DirFS(c.Static))
		return s, nil
	}
	sub, err := fs.Sub(files, "static")
	if err != nil {
		return nil, err
	}
	if s.Static, err = static.New(Assets, sub); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Site) Page(path, title string, content templ.Component) templ.Component {
	return layout(s, path, title, content)
}

const InternalError = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Internal Server Error</title>
</head>
<body>
<main id="main">
<h1>Internal Server Error</h1>
<p>Something went wrong on our side. <a href="/">Go home</a>.</p>
</main>
</body>
</html>
`
