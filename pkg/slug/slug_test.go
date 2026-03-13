package slug

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Fix login bug", "fix-login-bug"},
		{"Something does not work!!!", "something-does-not-work"},
		{"Add OAuth (GitHub)", "add-oauth-github"},
		{"Fix   multiple     spaces", "fix-multiple-spaces"},
		{"Café login failure", "cafe-login-failure"},
		{"Añadir autenticación", "anadir-autenticacion"},
		{"  Leading and trailing  ", "leading-and-trailing"},
		{"snake_case_title", "snake-case-title"},
		{"multiple---dashes", "multiple-dashes"},
	}

	for _, tt := range tests {
		got := Generate(tt.input)

		if got != tt.want {
			t.Errorf("Generate(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
