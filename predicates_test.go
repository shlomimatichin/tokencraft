package tokencraft

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
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
	found := FindAllSpellings(tokens, []string{"void", "it", "("})
	expected := [][]Token{
		{
			{IDENTIFIER, 23, 3, 2, "void", 9},
			{IDENTIFIER, 28, 3, 7, "it", 11},
			{SPECIAL, 30, 3, 9, "(", 12},
		},
	}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("got %q, want %q", tokens, expected)
	}
}

func TestParensMatch(t *testing.T) {
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
	open := tokens[4]
	if open.Spelling != "{" {
		t.Errorf("got %q, want %q", open.Spelling, "{")
	}
	close, i := FindMatchingParen(tokens[4:])
	if close == nil {
		t.Errorf("got %q, want %q", close, "}")
	}
	if close.Spelling != "}" {
		t.Errorf("got %q, want %q", close.Spelling, "}")
	}
	if close.TokenIndex != 35 {
		t.Errorf("got %q, want %q", close.TokenIndex, 35)
	}
	if i != 35-4 {
		t.Errorf("got %q, want %q", i, 35-4)
	}
	if close.BeginsOffset != 117 {
		t.Errorf("got %q, want %q", close.BeginsOffset, 117)
	}
}

func TestSplitParenAware(t *testing.T) {
	code := "a,b"
	tokens := Tokenize(code, HASH_IS_DIRECTIVE)
	parts := SplitParenAware(tokens, ",")
	if len(parts) != 2 {
		t.Errorf("got %q, want %q", len(parts), 2)
	}
	if JoinSpellings(parts[0], "") != "a" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[0], ""), "a")
	}
	if JoinSpellings(Strip(parts[1]), "") != "b" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[1], ""), "b")
	}

	code = "a(1,2),b"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	parts = SplitParenAware(tokens, ",")
	if len(parts) != 2 {
		t.Errorf("got %q, want %q", len(parts), 2)
	}
	if JoinSpellings(parts[0], "") != "a(1,2)" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[0], ""), "a(1,2)")
	}
	if JoinSpellings(Strip(parts[1]), "") != "b" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[1], ""), "b")
	}

	code = "a(b, c, d(e, f), g)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	open := tokens[1]
	if open.Spelling != "(" {
		t.Errorf("got %q, want %q", open.Spelling, "(")
	}
	close, _ := FindMatchingParen(tokens[1:])
	if close == nil {
		t.Errorf("got %q, want %q", close, ")")
	}
	if close.Spelling != ")" {
		t.Errorf("got %q, want %q", close.Spelling, "}")
	}
	if close.BeginsOffset != 18 {
		t.Errorf("got %q, want %q", close.TokenIndex, 18)
	}
	inParens := tokens[open.TokenIndex+1 : close.TokenIndex]
	parts = SplitParenAware(inParens, ",")
	if len(parts) != 4 {
		t.Errorf("got %d, want %d", len(parts), 4)
	}
	if JoinSpellings(parts[0], "") != "b" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[0], ""), "b")
	}
	if JoinSpellings(Strip(parts[1]), "") != "c" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[1], ""), "c")
	}
	if JoinSpellings(Strip(parts[2]), "") != "d(e, f)" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[2], ""), "d(e, f)")
	}
	if JoinSpellings(Strip(parts[3]), "") != "g" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[3], ""), "g")
	}
}

func TestSplitBug(t *testing.T) {
	code := "tenantactivity_mv | take 10; tenantactivity_mv | project EventID"
	tokens := Tokenize(code, HASH_IS_DIRECTIVE)
	tokens = DropWhitespaces(tokens)
	parts := Split(tokens, ";")
	if len(parts) != 2 {
		t.Errorf("got %d, want %d", len(parts), 2)
	}
	if JoinSpellings(parts[0], "") != "tenantactivity_mv|take10" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[0], ""), "tenantactivity_mv|take10")
	}
	if JoinSpellings(parts[1], "") != "tenantactivity_mv|projectEventID" {
		t.Errorf("got %q, want %q", JoinSpellings(parts[1], ""), "tenantactivity_mv|projectEventID")
	}
}

func TestFindMatchingClosingParen(t *testing.T) {
	// Test basic matching
	code := "a(b)"
	tokens := Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the closing paren at index 3
	closeIndex := 3
	open, i := FindMatchingClosingParen(tokens, closeIndex)
	if open == nil {
		t.Errorf("got nil, want opening paren")
	}
	if open.Spelling != "(" {
		t.Errorf("got %q, want %q", open.Spelling, "(")
	}
	if i != 1 {
		t.Errorf("got index %d, want %d", i, 1)
	}

	// Test nested parentheses
	code = "a(b(c)d)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the outer closing paren
	closeIndex = -1
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Spelling == ")" {
			closeIndex = i
			break
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open == nil {
		t.Errorf("got nil, want opening paren")
	}
	if open.Spelling != "(" {
		t.Errorf("got %q, want %q", open.Spelling, "(")
	}
	// Should match the first opening paren
	if tokens[i].Spelling != "(" || i != 1 {
		t.Errorf("got index %d with spelling %q, want index 1 with '('", i, tokens[i].Spelling)
	}

	// Test inner nested paren
	code = "a(b(c)d)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the inner closing paren
	closeCount := 0
	closeIndex = -1
	for i, tok := range tokens {
		if tok.Spelling == ")" {
			closeCount++
			if closeCount == 1 {
				closeIndex = i
				break
			}
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open == nil {
		t.Errorf("got nil, want opening paren")
	}
	if open.Spelling != "(" {
		t.Errorf("got %q, want %q", open.Spelling, "(")
	}
	// Should match the second opening paren
	if i != 3 {
		t.Errorf("got index %d, want %d", i, 3)
	}

	// Test with braces
	code = "class Name { void it() { } }"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the last closing brace
	closeIndex = -1
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Spelling == "}" {
			closeIndex = i
			break
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open == nil {
		t.Errorf("got nil, want opening brace")
	}
	if open.Spelling != "{" {
		t.Errorf("got %q, want %q", open.Spelling, "{")
	}

	// Test with brackets
	code = "array[index[0]]"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the outer closing bracket
	closeIndex = -1
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Spelling == "]" {
			closeIndex = i
			break
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open == nil {
		t.Errorf("got nil, want opening bracket")
	}
	if open.Spelling != "[" {
		t.Errorf("got %q, want %q", open.Spelling, "[")
	}

	// Test invalid index (negative)
	code = "a(b)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	open, i = FindMatchingClosingParen(tokens, -1)
	if open != nil {
		t.Errorf("got %v, want nil for invalid negative index", open)
	}
	if i != -1 {
		t.Errorf("got %d, want -1 for invalid negative index", i)
	}

	// Test invalid index (out of bounds)
	open, i = FindMatchingClosingParen(tokens, 100)
	if open != nil {
		t.Errorf("got %v, want nil for out of bounds index", open)
	}
	if i != -1 {
		t.Errorf("got %d, want -1 for out of bounds index", i)
	}

	// Test non-closing paren token
	code = "a(b)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	open, i = FindMatchingClosingParen(tokens, 0) // 'a' is not a closing paren
	if open != nil {
		t.Errorf("got %v, want nil for non-closing paren", open)
	}
	if i != -1 {
		t.Errorf("got %d, want -1 for non-closing paren", i)
	}

	// Test mismatched parentheses
	code = "a(b]"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the closing bracket
	closeIndex = -1
	for i, tok := range tokens {
		if tok.Spelling == "]" {
			closeIndex = i
			break
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open != nil {
		t.Errorf("got %v, want nil for mismatched parens", open)
	}
	if i != -1 {
		t.Errorf("got %d, want -1 for mismatched parens", i)
	}

	// Test no matching opening paren
	code = "b)"
	tokens = Tokenize(code, HASH_IS_DIRECTIVE)
	// Find the closing paren
	closeIndex = -1
	for i, tok := range tokens {
		if tok.Spelling == ")" {
			closeIndex = i
			break
		}
	}
	open, i = FindMatchingClosingParen(tokens, closeIndex)
	if open != nil {
		t.Errorf("got %v, want nil when no opening paren exists", open)
	}
	if i != -1 {
		t.Errorf("got %d, want -1 when no opening paren exists", i)
	}
}
