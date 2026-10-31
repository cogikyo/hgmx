//go:build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func ownGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killGroup(p *os.Process) {
	syscall.Kill(-p.Pid, syscall.SIGTERM)
}
