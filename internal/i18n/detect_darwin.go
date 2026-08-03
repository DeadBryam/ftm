//go:build darwin

package i18n

import (
	"context"
	"os/exec"
	"time"
)

// detectPlatformLang reads the user's preferred language from macOS.
//
// An app launched from Finder or bundled as a .app inherits none of the shell
// environment, so LANG and the LC_* variables are all empty there and the
// POSIX path reports nothing. Reading AppleLanguages is what makes the desktop
// app follow the system language instead of always falling back to English.
func detectPlatformLang() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return ""
	}

	return parseAppleLanguages(string(out))
}
