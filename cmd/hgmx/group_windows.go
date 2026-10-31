package main

import (
	"os"
	"os/exec"
)

func ownGroup(cmd *exec.Cmd) {}

func killGroup(p *os.Process) {
	p.Kill()
}
