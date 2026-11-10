package appid

import "regexp"

var pattern = regexp.MustCompile(`^[a-z0-9]+([.-][a-z0-9]+)*$`)

func Valid(id string) bool {
	return len(id) <= 64 && pattern.MatchString(id)
}
