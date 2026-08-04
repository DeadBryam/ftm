//go:build !windows

package i18n

import "testing"

func TestDetectSystemLangEnvPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		lcAll   string
		lcMsg   string
		lang    string
		want    string
		skipMac bool
	}{
		{
			name:  "LC_ALL wins over everything",
			lcAll: "es_ES.UTF-8",
			lcMsg: "en_US.UTF-8",
			lang:  "en_US.UTF-8",
			want:  LangES,
		},
		{
			name:  "LC_MESSAGES beats LANG",
			lcMsg: "es_ES.UTF-8",
			lang:  "en_US.UTF-8",
			want:  LangES,
		},
		{
			name: "LANG is the fallback",
			lang: "es_ES.UTF-8",
			want: LangES,
		},
		{
			name: "english stays english",
			lang: "en_GB.UTF-8",
			want: LangEN,
		},
		{
			name: "brazilian portuguese resolves to pt",
			lang: "pt_BR.UTF-8",
			want: LangPT,
		},
		{
			name: "untranslated language falls back to the default",
			lang: "fr_FR.UTF-8",
			want: DefaultLang,
		},
		{
			name:  "C locale is skipped in favour of the next variable",
			lcAll: "C",
		},
		{
			name: "empty environment falls back to the default",
			// On macOS this consults AppleLanguages instead, which depends on
			// the machine running the test.
			want:    DefaultLang,
			skipMac: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipMac && detectPlatformLang() != "" {
				t.Skip("macOS reports a platform language; environment fallback is not reachable")
			}

			t.Setenv("LC_ALL", tt.lcAll)
			t.Setenv("LC_MESSAGES", tt.lcMsg)
			t.Setenv("LANG", tt.lang)

			if got := detectSystemLang(); got != tt.want {
				t.Fatalf("detectSystemLang() = %q, want %q", got, tt.want)
			}
		})
	}
}
