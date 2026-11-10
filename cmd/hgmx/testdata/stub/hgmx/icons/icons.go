package icons

type Name string

const (
	UserCircle    Name = "user-circle"
	TriangleAlert Name = "triangle-alert"
	Plus          Name = "plus"
)

const Path = "/icons/"

var Hover = map[string]any{}

func New() any { return nil }
