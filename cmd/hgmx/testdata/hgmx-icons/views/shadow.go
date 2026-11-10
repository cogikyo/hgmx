package views

import "github.com/cogikyo/hgmx/icons"

var _ = icons.Plus

func shadow(icons struct {
	Path       string
	BadgeAlert bool
}) {
	println(icons.Path, icons.BadgeAlert)
}
