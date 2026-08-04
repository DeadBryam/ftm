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
