package i18n

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
