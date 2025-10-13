package tokencraft

func eatIdentifier(e *Eater) string {
	beginsOffset := e.offset
	for !e.done() && wordCharacter(e.current()) {
		e.advance()
	}
	return e.data[beginsOffset:e.offset]
}

func eatWhitespace(e *Eater) string {
	beginsOffset := e.offset
	for !e.done() && whitespace(e.current()) {
		e.advance()
	}
	return e.data[beginsOffset:e.offset]
}

func eatUntil(e *Eater, lookFor byte, backslashEscape bool) string {
	beginsOffset := e.offset
	for !e.done() {
		if e.current() == '\\' && backslashEscape {
			e.advance()
			if !e.done() {
				e.advance()
			}
			continue
		}
		if e.current() == lookFor {
			break
		}
		e.advance()
	}
	return e.data[beginsOffset:e.offset]
}

func eatUntilAfter(e *Eater, lookFor string) string {
	beginsOffset := e.offset
	for i := 0; i < len(lookFor) - 1 && !e.done(); i++ {
		e.advance()
	}
	for !e.done() {
		// Ensure we have enough data to look back
		if e.offset < len(lookFor) {
			e.advance()
			continue
		}
		equal := true
		for i := 0; i < len(lookFor); i++ {
			if e.data[e.offset - len(lookFor) + i] != lookFor[i] {
				equal = false
				break
			}
		}
		if equal {
			break
		}
		e.advance()
	}
	return e.data[beginsOffset:e.offset]
}

func eatQuotes(e *Eater) string {
	beginsOffset := e.offset
	open := e.current()
	if open != '\'' && open != '"' && open != '`' {
		panic("eatQuotes called with invalid open quote")
	}
	e.advance()
	for !e.done() {
		if e.current() == '\\' && !e.done() {
			e.advance()
			if !e.done() {
				e.advance()
			}
			continue
		}
		if e.current() == open {
			e.advance()
			break
		}
		e.advance()
	}
	return e.data[beginsOffset:e.offset]
}
