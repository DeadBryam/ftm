package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sthbryan/ftm/internal/i18n"
)

func Uninstall() error {
	binaryPath, err := exec.LookPath("ftm")
	if err != nil {
		return fmt.Errorf("%s", i18n.T("uninstall_not_found"))
	}

	absPath, err := filepath.EvalSymlinks(binaryPath)
	if err != nil {
		absPath = binaryPath
	}

	fmt.Println(i18n.TF("uninstall_removing", absPath))
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("%s", i18n.TF("uninstall_error", err.Error()))
	}

	fmt.Println(i18n.T("uninstall_success"))
	return nil
}
