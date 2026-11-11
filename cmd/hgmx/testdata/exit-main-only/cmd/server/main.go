package main

import (
	"log"
	"os"
)

func main() {
	defer func() { os.Exit(run()) }()
}

func run() int {
	log.Fatal("boom")
	return 0
}
