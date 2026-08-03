//go:build windows

package i18n

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName = kernel32.NewProc("GetUserDefaultLocaleName")
)

const localeNameMaxLength = 85

func detectSystemLang() string {
	// An explicit locale override wins. Git Bash, MSYS and WSL-style shells all
	// set these on Windows, and a user who sets them means them.
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if code := normalizeLocale(os.Getenv(key)); code != "" {
			return parseLangCode(code)
		}
	}

	if code := normalizeLocale(userDefaultLocaleName()); code != "" {
		return parseLangCode(code)
	}

	return DefaultLang
}

// userDefaultLocaleName returns the user's locale as a BCP 47 name such as
// "es-ES", or "" if the call fails.
func userDefaultLocaleName() string {
	var buf [localeNameMaxLength]uint16

	ret, _, _ := procGetUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(localeNameMaxLength),
	)
	if ret == 0 {
		return ""
	}

	return syscall.UTF16ToString(buf[:])
}
