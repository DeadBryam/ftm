package i18n

import "testing"

func seed(t *testing.T, translations map[string]map[string]string, lang string) {
	t.Helper()

	store.mu.Lock()
	previous := store.translations
	store.translations = translations
	store.mu.Unlock()

	langMu.Lock()
	previousLang := currentLang
	currentLang = lang
	langMu.Unlock()

	t.Cleanup(func() {
		store.mu.Lock()
		store.translations = previous
		store.mu.Unlock()

		langMu.Lock()
		currentLang = previousLang
		langMu.Unlock()
	})
}

func TestTFallsBackToEnglishThenKey(t *testing.T) {
	seed(t, map[string]map[string]string{
		LangEN: {"start": "Start", "only_in_en": "English only"},
		LangES: {"start": "Iniciar"},
	}, LangES)

	if got := T("start"); got != "Iniciar" {
		t.Errorf("T(start) = %q, want the Spanish string", got)
	}
	if got := T("only_in_en"); got != "English only" {
		t.Errorf("T(only_in_en) = %q, want the English fallback", got)
	}
	if got := T("missing_everywhere"); got != "missing_everywhere" {
		t.Errorf("T(missing_everywhere) = %q, want the key itself", got)
	}
}

func TestTFSubstitutesPlaceholders(t *testing.T) {
	seed(t, map[string]map[string]string{
		LangEN: {
			"greet":  "Hello {0}, you have {1} tunnels",
			"repeat": "{0} and {0}",
			"none":   "no placeholders",
		},
	}, LangEN)

	if got := TF("greet", "Bryan", 3); got != "Hello Bryan, you have 3 tunnels" {
		t.Errorf("TF(greet) = %q", got)
	}
	if got := TF("none", "unused"); got != "no placeholders" {
		t.Errorf("TF(none) = %q, want the string unchanged", got)
	}

	if got := TF("repeat", "x"); got != "x and {0}" {
		t.Errorf("TF(repeat) = %q", got)
	}
}

func TestSetLanguageWithFallbackRejectsUnknown(t *testing.T) {
	seed(t, map[string]map[string]string{
		LangEN: {"start": "Start"},
		LangES: {"start": "Iniciar"},
	}, LangEN)

	SetLanguageWithFallback(LangES)
	if got := CurrentLanguage(); got != LangES {
		t.Errorf("CurrentLanguage() = %q, want %q", got, LangES)
	}

	SetLanguageWithFallback("klingon")
	if got := CurrentLanguage(); got != DefaultLang {
		t.Errorf("CurrentLanguage() = %q after an unknown language, want %q", got, DefaultLang)
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		header string
		want   string
	}{
		{"", DefaultLang},
		{"es", LangES},
		{"es-ES", LangES},
		{"es-MX,es;q=0.9,en;q=0.8", LangES},
		{"en-US,en;q=0.9", LangEN},
		{"fr-FR,fr;q=0.9", DefaultLang},
		{"not a language tag", DefaultLang},

		{"de-DE,es;q=0.7", LangES},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			if got := ParseAcceptLanguage(tt.header); got != tt.want {
				t.Fatalf("ParseAcceptLanguage(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

func TestParseLangCode(t *testing.T) {
	tests := map[string]string{
		"en": LangEN,
		"es": LangES,
		"fr": DefaultLang,
		"":   DefaultLang,
	}

	for code, want := range tests {
		if got := parseLangCode(code); got != want {
			t.Errorf("parseLangCode(%q) = %q, want %q", code, got, want)
		}
	}
}

func TestEmbeddedLocalesLoadWithMatchingKeys(t *testing.T) {
	seed(t, map[string]map[string]string{}, DefaultLang)

	if err := Load(); err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	en := GetTranslations(DefaultLang)
	if len(en) == 0 {
		t.Fatal("default locale loaded no keys")
	}

	for _, lang := range SupportedLanguages() {
		if lang == DefaultLang {
			continue
		}

		other := GetTranslations(lang)
		if len(other) == 0 {
			t.Errorf("%s.yaml loaded no keys", lang)
			continue
		}

		for key := range en {
			if _, ok := other[key]; !ok {
				t.Errorf("key %q is missing from %s.yaml", key, lang)
			}
		}
		for key := range other {
			if _, ok := en[key]; !ok {
				t.Errorf("key %q is in %s.yaml but not %s.yaml", key, lang, DefaultLang)
			}
		}
	}
}
