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
	case "pt":
		return LangPT
	default:
		return DefaultLang
	}
}

func normalizeLocale(value string) string {
	value = strings.TrimSpace(value)

	if idx := strings.IndexByte(value, '@'); idx != -1 {
		value = value[:idx]
	}
	if idx := strings.IndexByte(value, '.'); idx != -1 {
		value = value[:idx]
	}

	switch strings.ToUpper(value) {
	case "", "C", "POSIX":
		return ""
	}

	value = strings.ReplaceAll(value, "_", "-")

	tag, err := language.Parse(value)
	if err != nil {
		return ""
	}

	base, _ := tag.Base()
	return base.String()
}
