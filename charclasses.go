package tokencraft

func digitCharacter(c byte) bool {
	return c >= '0' && c <= '9'
}

func alphabeticCharacter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func wordCharacter(c byte) bool {
	return alphabeticCharacter(c) || digitCharacter(c) || c == '_'
}

func whitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t'
}
