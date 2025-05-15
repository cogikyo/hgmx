package hgmx

import (
	"embed"
)

//go:embed .version
var version string

func Version() string {
	return version
}

//go:embed library/**
var LibraryFS embed.FS

func Library() embed.FS {
	return LibraryFS
}
