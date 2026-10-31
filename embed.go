package hgmx

import (
	"embed"
	"runtime/debug"
)

// Version reads hgmx's module version from Go build info, or returns empty.
// Local builds can report "(devel)" instead of a release tag.
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if info.Main.Path == "github.com/cogikyo/hgmx" {
		return info.Main.Version
	}
	for _, d := range info.Deps {
		if d.Path == "github.com/cogikyo/hgmx" {
			return d.Version
		}
	}
	return ""
}

// GoMod is hgmx's module file, used to check the site's minimum Go version.
//
//go:embed go.mod
var GoMod []byte

//go:embed library/**
var LibraryFS embed.FS
