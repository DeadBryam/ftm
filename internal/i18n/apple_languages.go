package i18n

import "strings"

//
//	(
//	    "en-SV",
//	    "es-SV"
//	)
//

func parseAppleLanguages(out string) string {
	for _, line := range strings.Split(out, "\n") {
		entry := strings.TrimSpace(line)
		entry = strings.Trim(entry, "(),")
		entry = strings.TrimSpace(entry)
		entry = strings.Trim(entry, `"`)

		if entry == "" {
			continue
		}
		if code := normalizeLocale(entry); code != "" {
			return code
		}
	}

	return ""
}
