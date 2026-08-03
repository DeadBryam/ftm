//go:build !windows

package i18n

import "os"

// localeEnvVars is the POSIX precedence order: LC_ALL overrides everything,
// LC_MESSAGES governs program messages specifically, and LANG is the fallback
// default. Reading only LANG gets the wrong answer for anyone who overrides
// their message language without changing the rest of their locale.
var localeEnvVars = []string{"LC_ALL", "LC_MESSAGES", "LANG"}

func detectSystemLang() string {
	for _, key := range localeEnvVars {
		if code := normalizeLocale(os.Getenv(key)); code != "" {
			return parseLangCode(code)
		}
	}

	// GUI apps get an empty environment on macOS, so there is a platform
	// fallback. Elsewhere this reports nothing.
	if code := detectPlatformLang(); code != "" {
		return parseLangCode(code)
	}

	return DefaultLang
}
