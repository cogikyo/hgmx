package appid

import "regexp"

var (
	pattern = regexp.MustCompile(`^[a-z0-9]+([.-][a-z0-9]+)*$`)
	bundle  = regexp.MustCompile(`^[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)+$`)
	pkg     = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*(\.[A-Za-z][A-Za-z0-9_]*)+$`)
)

func Valid(id string) bool {
	return len(id) <= 64 && pattern.MatchString(id)
}

func ValidBundle(id string) bool {
	return bundle.MatchString(id)
}

func ValidPackage(name string) bool {
	return pkg.MatchString(name)
}
