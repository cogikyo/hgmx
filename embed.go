package hgmx

import (
	"embed"
	"runtime/debug"
)

// Version reads hgmx's module version from Go build info, or returns empty.
// Local builds can report "(devel)" instead of a release tag.
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Path != "github.com/cogikyo/hgmx" {
		return ""
	}
	return info.Main.Version
}

// GoMod is hgmx's module file, used to check the site's minimum Go version.
//
//go:embed go.mod
var GoMod []byte

//go:embed library/**
var LibraryFS embed.FS
