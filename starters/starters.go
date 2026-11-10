package starters

import "embed"

//go:embed apps apps/android/.gitignore apps/ios/.gitignore
var FS embed.FS
