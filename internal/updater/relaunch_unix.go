//go:build !windows

package updater

import (
	"os"
	"syscall"
)

func relaunch(execPath string) error {
	return syscall.Exec(execPath, append([]string{execPath}, os.Args[1:]...), os.Environ())
}
