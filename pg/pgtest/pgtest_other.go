//go:build !unix

package pgtest

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func lockFile(*os.File) error {
	return fmt.Errorf("%w: file locking is not supported on %s", unavailable, runtime.GOOS)
}

func tryLockFile(*os.File) (bool, error) {
	return false, lockFile(nil)
}

func detach(*exec.Cmd) {}
