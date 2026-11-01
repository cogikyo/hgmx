//go:build unix

package pgtest

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

func lockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

func tryLockFile(f *os.File) (bool, error) {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		return false, nil
	}
	return err == nil, err
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
