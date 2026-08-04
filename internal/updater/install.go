package updater

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Method string

const (
	MethodSelf     Method = "self"
	MethodStore    Method = "store"
	MethodHomebrew Method = "homebrew"
	MethodDownload Method = "download"
)

var ErrNotSelfUpdatable = errors.New("this installation cannot replace its own binary")

var desktopBuild bool

func MarkDesktop() {
	desktopBuild = true
}

func DetectMethod() Method {
	execPath, err := os.Executable()
	if err != nil {
		return MethodDownload
	}
	if resolved, err := filepath.EvalSymlinks(execPath); err == nil {
		execPath = resolved
	}

	return methodFor(execPath, runtime.GOOS, desktopBuild)
}

func methodFor(execPath, goos string, desktop bool) Method {
	path := strings.ToLower(strings.ReplaceAll(execPath, `\`, "/"))

	if goos == "windows" && strings.Contains(path, "/windowsapps/") {
		return MethodStore
	}

	if strings.Contains(path, "/cellar/") || strings.Contains(path, "/homebrew/") {
		return MethodHomebrew
	}

	if desktop || strings.Contains(path, ".app/contents/macos/") {
		return MethodDownload
	}

	return MethodSelf
}
