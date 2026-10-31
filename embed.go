package hgmx

import (
	"embed"
	"runtime/debug"
)

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

//go:embed go.mod
var GoMod []byte

//go:embed library/**
var LibraryFS embed.FS
