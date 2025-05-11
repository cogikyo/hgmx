package hgmx

import (
	"embed"
)

//go:embed .version
var version string

//go:embed library/**
var LibraryFS embed.FS

func Version() string {
	return version
}

func Library() embed.FS {
	return LibraryFS
}
