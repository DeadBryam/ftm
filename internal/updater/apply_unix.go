//go:build !windows

package updater

import (
	"os"
	"os/exec"
	"runtime"
)

func applyUpdate(execPath, tmpPath string) error {
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return err
	}

	oldPath := execPath + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(execPath, oldPath); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, execPath); err != nil {
		_ = os.Rename(oldPath, execPath)
		os.Remove(tmpPath)
		return err
	}

	_ = os.Remove(oldPath)

	if runtime.GOOS == "darwin" {
		_ = exec.Command("xattr", "-d", "com.apple.quarantine", execPath).Run()
	}
	return nil
}
