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
