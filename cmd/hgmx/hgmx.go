package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/nosvagor/hgmx"
)

func main() {
	code := run(os.Stdin, os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usageText = `usage: hgmx <command> [<args>...]

hgmx - A component management tool for Go Templ and HTMX projects.

commands:
  info      	Displays information about the hgmx environment
  init      	Initializes a new hgmx project
  palette   	Generates a color palette based on the input hex color
  version 	-v  Prints the versio
`

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(stderr, usageText)
		return 64
	}

	switch args[1] {
	case "info":
		return infoCmd(stdout, stderr, args[2:])
	case "init":
		return initCmd(stdout, stderr, args[2:])
	case "palette":
		return paletteCmd(stdout, stderr, args[2:])
	// TODO: Add 'add' command
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, hgmx.Version())
		return 0
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, usageText)
		return 0
	}

	fmt.Fprintf(stderr, "Unknown command: %q\n", args[1])
	fmt.Fprint(stderr, usageText)
	return 64
}
