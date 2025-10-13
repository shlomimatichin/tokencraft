package tokencraft

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTokenizePanicRepro(t *testing.T) {
	// Minimal reproduction of panic from /tmp/t.ts
	// Bug: eatUntilAfter panics with index out of range [-1] when
	// processing unclosed multi-line comment starting with "/*"
	testCases := []struct {
		name     string
		code     string
		expected int
	}{
		{"unclosed comment 2 chars", "/*", 1},
		{"unclosed comment 3 chars", "/* ", 1},
		{"unclosed comment with text", "/* hello", 1},
		{"closed comment", "/* test */", 1},
		{"comment then code", "/* test */ x", 3}, // comment, whitespace, identifier
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := Tokenize(tc.code, HASH_IS_DIRECTIVE)
			// Should not panic
			if len(tokens) < 1 {
				t.Errorf("Expected at least 1 token, got %d", len(tokens))
			}
			if tokens[0].Type != COMMENT {
				t.Errorf("Expected first token to be COMMENT, got %v", tokens[0].Type)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	code := "class Name {\n" +
		" public:\n" +
		" void it(int b) {\n" +
		"  TRACE_INFO(\"Hello \"<<7);\n" +
		"  //comment 1\n" +
		" /* comment 2 */\n" +
		"#define Bu To\\\n" +
		"Ba\n" +
		"}};"

	tokens := Tokenize(code, HASH_IS_DIRECTIVE)
	toString := map[TokenType]string{
		SPECIAL:    "SPECIAL",
		IDENTIFIER: "IDENTIFIER",
		QUOTES:     "QUOTES",
		COMMENT:    "COMMENT",
		DIRECTIVE:  "DIRECTIVE",
		WHITESPACE: "WHITESPACE",
	}
	for _, token := range tokens {
		fmt.Printf("{%s, %d, %d, %d, \"%s\", %d},\n", toString[token.Type], token.BeginsOffset, token.BeginsLine, token.BeginColumn, token.Spelling, token.TokenIndex)
	}
	expected := []Token{
		{IDENTIFIER, 0, 1, 1, "class", 0},
		{WHITESPACE, 5, 1, 6, " ", 1},
		{IDENTIFIER, 6, 1, 7, "Name", 2},
		{WHITESPACE, 10, 1, 11, " ", 3},
		{SPECIAL, 11, 1, 12, "{", 4},
		{WHITESPACE, 12, 1, 13, "\n ", 5},
		{IDENTIFIER, 14, 2, 2, "public", 6},
		{SPECIAL, 20, 2, 8, ":", 7},
		{WHITESPACE, 21, 2, 9, "\n ", 8},
		{IDENTIFIER, 23, 3, 2, "void", 9},
		{WHITESPACE, 27, 3, 6, " ", 10},
		{IDENTIFIER, 28, 3, 7, "it", 11},
		{SPECIAL, 30, 3, 9, "(", 12},
		{IDENTIFIER, 31, 3, 10, "int", 13},
		{WHITESPACE, 34, 3, 13, " ", 14},
		{IDENTIFIER, 35, 3, 14, "b", 15},
		{SPECIAL, 36, 3, 15, ")", 16},
		{WHITESPACE, 37, 3, 16, " ", 17},
		{SPECIAL, 38, 3, 17, "{", 18},
		{WHITESPACE, 39, 3, 18, "\n  ", 19},
		{IDENTIFIER, 42, 4, 3, "TRACE_INFO", 20},
		{SPECIAL, 52, 4, 13, "(", 21},
		{QUOTES, 53, 4, 14, "\"Hello \"", 22},
		{SPECIAL, 61, 4, 22, "<<", 23},
		{IDENTIFIER, 63, 4, 24, "7", 24},
		{SPECIAL, 64, 4, 25, ")", 25},
		{SPECIAL, 65, 4, 26, ";", 26},
		{WHITESPACE, 66, 4, 27, "\n  ", 27},
		{COMMENT, 69, 5, 3, "//comment 1", 28},
		{WHITESPACE, 80, 5, 14, "\n ", 29},
		{COMMENT, 82, 6, 2, "/* comment 2 */", 30},
		{WHITESPACE, 97, 6, 17, "\n", 31},
		{DIRECTIVE, 98, 7, 1, "#define Bu To\\\nBa", 32},
		{WHITESPACE, 115, 8, 3, "\n", 33},
		{SPECIAL, 116, 9, 1, "}", 34},
		{SPECIAL, 117, 9, 2, "}", 35},
		{SPECIAL, 118, 9, 3, ";", 36},
	}
	for i, token := range tokens {
		if !reflect.DeepEqual(token, expected[i]) {
			t.Errorf("%d: got %q, want %q", i, token, expected[i])
		}
	}
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("got %q, want %q", tokens, expected)
	}
}
