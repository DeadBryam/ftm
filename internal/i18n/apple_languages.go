package i18n

import "strings"

// parseAppleLanguages pulls the first usable language out of the output of
// `defaults read -g AppleLanguages`, an old-style plist array:
//
//	(
//	    "en-SV",
//	    "es-SV"
//	)
//
// Entries are not always quoted, and the first entry is the user's preferred
// language.
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
