//go:build windows

package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func applyUpdate(execPath, tmpPath string) error {
	exeName := filepath.Base(execPath)
	exeDir := filepath.Dir(execPath)
	newPath := filepath.Join(exeDir, exeName+".new")

	if err := os.Rename(tmpPath, newPath); err != nil {
		return fmt.Errorf("stage new binary: %w", err)
	}

	batPath := filepath.Join(os.TempDir(), "ftm-update.bat")
	bat := fmt.Sprintf(
		`@echo off
:loop
timeout /t 1 /nobreak >nul
move /y "%s" "%s" >nul 2>&1
if errorlevel 1 goto loop
del "%%~f0"
`, newPath, execPath)

	if err := os.WriteFile(batPath, []byte(bat), 0644); err != nil {
		return fmt.Errorf("write updater script: %w", err)
	}

	cmd := exec.Command("cmd", "/c", "start", "", "/b", batPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch updater script: %w", err)
	}
	return nil
}
