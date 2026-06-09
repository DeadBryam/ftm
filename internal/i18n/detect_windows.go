//go:build windows

package i18n

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName = kernel32.NewProc("GetUserDefaultLocaleName")
)

const localeNameMaxLength = 85

func detectSystemLang() string {
	var buf [localeNameMaxLength]uint16
	ret, _, _ := procGetUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(localeNameMaxLength),
	)
	if ret == 0 {
		return DefaultLang
	}

	lang := syscall.UTF16ToString(buf[:])
	if len(lang) < 2 {
		return DefaultLang
	}

	return parseLangCode(lang[:2])
}
