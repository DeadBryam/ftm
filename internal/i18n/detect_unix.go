//go:build !windows

package i18n

import (
	"os"

	"golang.org/x/text/language"
)

func detectSystemLang() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		return DefaultLang
	}

	tag, err := language.Parse(lang)
	if err != nil {
		return DefaultLang
	}

	base, _ := tag.Base()
	return parseLangCode(base.String())
}
