package i18n

import (
	"strings"

	"golang.org/x/text/language"
)

func parseLangCode(code string) string {
	switch code {
	case "en":
		return LangEN
	case "es":
		return LangES
	default:
		return DefaultLang
	}
}

// normalizeLocale extracts the language subtag from a POSIX locale value such
// as "es_ES.UTF-8@euro", returning "" when the value names no language.
//
// The stripping has to happen before parsing: language.Parse rejects the
// codeset and modifier suffixes outright, so the plain $LANG of nearly every
// Unix shell ("es_ES.UTF-8") fails to parse and falls back to English.
func normalizeLocale(value string) string {
	value = strings.TrimSpace(value)

	// Drop the @modifier and .codeset suffixes.
	if idx := strings.IndexByte(value, '@'); idx != -1 {
		value = value[:idx]
	}
	if idx := strings.IndexByte(value, '.'); idx != -1 {
		value = value[:idx]
	}

	// "C" and "POSIX" mean "no localisation", not a language.
	switch strings.ToUpper(value) {
	case "", "C", "POSIX":
		return ""
	}

	// POSIX separates language and territory with "_", BCP 47 with "-".
	value = strings.ReplaceAll(value, "_", "-")

	tag, err := language.Parse(value)
	if err != nil {
		return ""
	}

	base, _ := tag.Base()
	return base.String()
}
