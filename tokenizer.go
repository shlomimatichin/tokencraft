package tokencraft

type HashMode int

const (
	HASH_IS_DIRECTIVE HashMode = iota
	HASH_IS_COMMENT
	HASH_IS_SPECIAL
)

type TokenType int

const (
	SPECIAL TokenType = iota
	IDENTIFIER
	QUOTES
	COMMENT
	DIRECTIVE
	WHITESPACE
)

type Token struct {
	Type         TokenType
	BeginsOffset int
	BeginsLine   int
	BeginColumn  int
	Spelling     string
	TokenIndex   int
}

func Tokenize(data string, hashMode HashMode) []Token {
	e := newEater(data)
	tokens := []Token{}

	for !e.done() {
		c := e.current()
		next := e.next()
		offset := e.offset
		line := e.line
		column := e.column
		if wordCharacter(c) {
			tokens = append(tokens, Token{IDENTIFIER, offset, line, column, eatIdentifier(e), len(tokens)})
		} else if c == '\'' || c == '"' || c == '`' {
			tokens = append(tokens, Token{QUOTES, offset, line, column, eatQuotes(e), len(tokens)})
		} else if c == '/' && next == '/' {
			tokens = append(tokens, Token{COMMENT, offset, line, column, eatUntil(e, '\n', true), len(tokens)})
		} else if c == '/' && next == '*' {
			tokens = append(tokens, Token{COMMENT, offset, line, column, eatUntilAfter(e, "*/"), len(tokens)})
		} else if c == '#' {
			switch hashMode {
			case HASH_IS_DIRECTIVE:
				tokens = append(tokens, Token{DIRECTIVE, offset, line, column, eatUntil(e, '\n', true), len(tokens)})
			case HASH_IS_COMMENT:
				tokens = append(tokens, Token{COMMENT, offset, line, column, eatUntil(e, '\n', true), len(tokens)})
			case HASH_IS_SPECIAL:
				tokens = append(tokens, Token{SPECIAL, offset, line, column, "#", len(tokens)})
				e.advance()
			}
		} else if c == '<' && next == '<' {
			tokens = append(tokens, Token{SPECIAL, offset, line, column, "<<", len(tokens)})
			e.advance()
			e.advance()
		} else if c == ':' && next == ':' {
			tokens = append(tokens, Token{SPECIAL, offset, line, column, "::", len(tokens)})
			e.advance()
			e.advance()
		} else if whitespace(c) {
			tokens = append(tokens, Token{WHITESPACE, offset, line, column, eatWhitespace(e), len(tokens)})
		} else {
			tokens = append(tokens, Token{SPECIAL, offset, line, column, data[offset : offset+1], len(tokens)})
			e.advance()
		}
	}
	return tokens
}
