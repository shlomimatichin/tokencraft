package tokencraft

func RebuildWithWhitespace(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}
	offset := tokens[0].BeginsOffset
	length := tokens[len(tokens)-1].BeginsOffset + len(tokens[len(tokens)-1].Spelling) - offset
	result := make([]rune, length)
	for i := range result {
		result[i] = ' '
	}
	for _, token := range tokens {
		for i, c := range token.Spelling {
			result[token.BeginsOffset-offset+i] = c
		}
	}
	return string(result)
}