//go:build !windows && !darwin

package i18n

// detectPlatformLang has no equivalent outside macOS: on Linux and the BSDs the
// LC_* and LANG variables are the system language.
func detectPlatformLang() string { return "" }
