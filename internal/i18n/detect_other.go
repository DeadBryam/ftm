//go:build !windows && !darwin

package i18n

func detectPlatformLang() string { return "" }
