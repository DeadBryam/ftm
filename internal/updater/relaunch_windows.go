//go:build windows

package updater

import (
	"os"
	"os/exec"
)

func relaunch(execPath string) error {
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	os.Exit(0)

	return nil
}
