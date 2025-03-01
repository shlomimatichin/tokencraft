package tokencraft

import "testing"

func TestRebuildWithWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected string
	}{
		{
			name:     "empty tokens",
			tokens:   []Token{},
			expected: "",
		},
		{
			name: "single token",
			tokens: []Token{
				{BeginsOffset: 0, Spelling: "hello"},
			},
			expected: "hello",
		},
		{
			name: "tokens with space",
			tokens: []Token{
				{BeginsOffset: 0, Spelling: "hello"},
				{BeginsOffset: 6, Spelling: "world"},
			},
			expected: "hello world",
		},
		{
			name: "tokens with multiple spaces",
			tokens: []Token{
				{BeginsOffset: 0, Spelling: "int"},
				{BeginsOffset: 4, Spelling: "main"},
				{BeginsOffset: 9, Spelling: "()"},
			},
			expected: "int main ()",
		},
		{
			name: "tokens with newlines",
			tokens: []Token{
				{BeginsOffset: 0, Spelling: "if"},
				{BeginsOffset: 3, Spelling: "("},
				{BeginsOffset: 4, Spelling: "true"},
				{BeginsOffset: 8, Spelling: ")"},
				{BeginsOffset: 10, Spelling: "{"},
				{BeginsOffset: 14, Spelling: "}"},
			},
			expected: "if (true) {   }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RebuildWithWhitespace(tt.tokens)
			if got != tt.expected {
				t.Errorf("RebuildWithWhitespace() = %q, want %q", got, tt.expected)
			}
		})
	}
}
