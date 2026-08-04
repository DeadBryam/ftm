//go:build !windows

package i18n

import "os"

var localeEnvVars = []string{"LC_ALL", "LC_MESSAGES", "LANG"}

func detectSystemLang() string {
	for _, key := range localeEnvVars {
		if code := normalizeLocale(os.Getenv(key)); code != "" {
			return parseLangCode(code)
		}
	}

	if code := detectPlatformLang(); code != "" {
		return parseLangCode(code)
	}

	return DefaultLang
}
