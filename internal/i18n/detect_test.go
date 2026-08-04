package i18n

import "testing"

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{

		{"es_ES.UTF-8", "es"},
		{"en_US.UTF-8", "en"},
		{"es_ES.utf8", "es"},
		{"es_ES@euro", "es"},
		{"es_ES.UTF-8@euro", "es"},

		{"es_ES", "es"},
		{"es-ES", "es"},
		{"es", "es"},
		{"en", "en"},

		{"fr_FR.UTF-8", "fr"},
		{"pt_BR.UTF-8", "pt"},

		{"C", ""},
		{"POSIX", ""},
		{"C.UTF-8", ""},
		{"c", ""},

		{"", ""},
		{"   ", ""},
		{"nonsense value", ""},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := normalizeLocale(tt.value); got != tt.want {
				t.Fatalf("normalizeLocale(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestResolveLanguage(t *testing.T) {

	if got := ResolveLanguage("es"); got != "es" {
		t.Errorf("ResolveLanguage(es) = %q, want es", got)
	}
	if got := ResolveLanguage("en"); got != "en" {
		t.Errorf("ResolveLanguage(en) = %q, want en", got)
	}

	for _, input := range []string{"", "auto"} {
		got := ResolveLanguage(input)
		if got != LangEN && got != LangES {
			t.Errorf("ResolveLanguage(%q) = %q, want a supported language", input, got)
		}
	}
}
