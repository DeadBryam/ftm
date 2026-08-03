package i18n

import "testing"

// seed replaces the package-level store for a single test and restores it
// afterwards, since translations are global process state.
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
	// Each argument replaces one occurrence, so a repeated placeholder keeps
	// its second instance literal.
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
		// The first supported tag wins, even if an unsupported one precedes it.
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

// The embedded locales must actually load, and both languages must define the
// same keys, or the UI silently falls back to English in places.
func TestEmbeddedLocalesLoadWithMatchingKeys(t *testing.T) {
	seed(t, map[string]map[string]string{}, DefaultLang)

	if err := Load(); err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	en := GetTranslations(LangEN)
	es := GetTranslations(LangES)

	if len(en) == 0 {
		t.Fatal("English locale loaded no keys")
	}
	if len(es) == 0 {
		t.Fatal("Spanish locale loaded no keys")
	}

	for key := range en {
		if _, ok := es[key]; !ok {
			t.Errorf("key %q is missing from es.yaml", key)
		}
	}
	for key := range es {
		if _, ok := en[key]; !ok {
			t.Errorf("key %q is in es.yaml but not en.yaml", key)
		}
	}
}
