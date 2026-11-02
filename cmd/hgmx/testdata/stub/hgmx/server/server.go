package server

type CSP map[string][]string

func (c CSP) Header() string { return "" }
