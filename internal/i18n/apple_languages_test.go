package i18n

import "testing"

func TestParseAppleLanguages(t *testing.T) {
	tests := []struct {
		name string
		out  string
		want string
	}{
		{
			name: "quoted array",
			out:  "(\n    \"en-SV\",\n    \"es-SV\"\n)\n",
			want: "en",
		},
		{
			name: "spanish first",
			out:  "(\n    \"es-ES\",\n    \"en-US\"\n)\n",
			want: "es",
		},
		{
			name: "unquoted entries",
			out:  "(\n    es,\n    en\n)\n",
			want: "es",
		},
		{
			name: "single entry",
			out:  "(\n    \"es-MX\"\n)\n",
			want: "es",
		},
		{
			name: "empty array",
			out:  "(\n)\n",
			want: "",
		},
		{
			name: "garbage",
			out:  "could not read domain",
			want: "",
		},
		{
			name: "empty",
			out:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAppleLanguages(tt.out); got != tt.want {
				t.Fatalf("parseAppleLanguages(%q) = %q, want %q", tt.out, got, tt.want)
			}
		})
	}
}
