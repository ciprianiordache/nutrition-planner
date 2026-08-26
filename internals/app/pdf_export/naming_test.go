package pdf_export

import "testing"

func TestSanitizePathSegment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fallback string
		want     string
	}{
		{"plain name", "Andreea Popescu", "client", "Andreea Popescu"},
		{"slash", "Client/Weird", "client", "Client-Weird"},
		{"backslash and colon", `C:\evil`, "client", "C-evil"},
		{"wildcards and quotes", `a*b?c"d<e>f|g`, "client", "a-b-c-d-e-f-g"},
		{"leading/trailing dots and spaces", "  .hidden.  ", "client", "hidden"},
		{"empty falls back", "", "client", "client"},
		{"only illegal chars collapse to a single dash", "///", "client", "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizePathSegment(tt.input, tt.fallback)
			if got != tt.want {
				t.Errorf("sanitizePathSegment(%q, %q) = %q, want %q", tt.input, tt.fallback, got, tt.want)
			}
		})
	}
}
