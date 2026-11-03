// Package skills embeds hgmx agent skills for installation by the hgmx CLI.
package skills

import "embed"

//go:embed all:hgmx*
var FS embed.FS
