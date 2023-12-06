package tokencraft

type GolangMethod struct {
	Type           string
	TypePointer    bool
	Name           string
	StartLine      int
	ParametersSpan []Token
	ReturnTypeSpan []Token
	BodySpan       []Token
}

func GolangAllMethods(tokens []Token) []GolangMethod {
	methods := []GolangMethod{}
	for _, match := range FindAllSpellings(tokens, []string{"func", "("}) {
		if match[0].BeginColumn != 1 {
			continue
		}
		bindSpanOpenIndex := match[len(match)-1]
		bindSpanEnd, _ := FindMatchingParen(tokens[bindSpanOpenIndex.TokenIndex:])
		if bindSpanEnd == nil {
			continue
		}
		bindSpan := DropWhitespaces(tokens[bindSpanOpenIndex.TokenIndex+1 : bindSpanEnd.TokenIndex])
		if len(bindSpan) > 3 || len(bindSpan) < 2 || bindSpan[0].Type != IDENTIFIER || bindSpan[len(bindSpan)-1].Type != IDENTIFIER {
			continue
		}
		if len(bindSpan) == 3 && bindSpan[1].Spelling != "*" {
			continue
		}
		method := GolangMethod{
			StartLine:   match[0].BeginsLine,
			Type:        bindSpan[len(bindSpan)-1].Spelling,
			TypePointer: len(bindSpan) == 3,
		}
		nameTokenIndex := NextNoWhitespaceIndex(tokens, bindSpanEnd.TokenIndex+1)
		if nameTokenIndex >= len(tokens) {
			continue
		}

		nameToken := tokens[nameTokenIndex]
		if nameToken.Type != IDENTIFIER {
			continue
		}
		method.Name = nameToken.Spelling
		parensSpanOpenIndex := NextNoWhitespaceIndex(tokens, nameTokenIndex+1)
		if parensSpanOpenIndex >= len(tokens) {
			continue
		}

		parensClose, _ := FindMatchingParen(tokens[parensSpanOpenIndex:])
		if parensClose == nil {
			continue
		}
		method.ParametersSpan = tokens[parensSpanOpenIndex+1 : parensClose.TokenIndex]
		nextIndex := NextNoWhitespaceIndex(tokens, parensClose.TokenIndex+1)
		if nextIndex >= len(tokens) {
			continue
		}

		nextToken := tokens[nextIndex]
		if nextToken.Spelling == "{" {
			// no return type
		} else if nextToken.Spelling == "(" {
			// return type
			returnSpanClose, _ := FindMatchingParen(tokens[nextIndex:])
			if returnSpanClose == nil {
				continue
			}
			method.ReturnTypeSpan = tokens[nextIndex+1 : returnSpanClose.TokenIndex]
			nextIndex = NextNoWhitespaceIndex(tokens, returnSpanClose.TokenIndex+1)
			if nextIndex >= len(tokens) {
				continue
			}
		} else {
			bodyStart := FindNext(tokens[nextIndex:], "{")
			if bodyStart == nil {
				continue
			}
			method.ReturnTypeSpan = Strip(tokens[nextIndex:bodyStart.TokenIndex])
			nextIndex = bodyStart.TokenIndex
		}

		nextToken = tokens[nextIndex]
		if nextToken.Spelling != "{" {
			continue
		}

		bodyClose, _ := FindMatchingParen(tokens[nextIndex:])
		method.BodySpan = tokens[nextIndex+1 : bodyClose.TokenIndex]
		methods = append(methods, method)
	}
	return methods
}

type GolangChainedCall struct {
	StartLine      int
	Name           Token
	ParametersSpan []Token
	Parameters     [][]Token
}

func GolangCallChain(tokens []Token, startPos int) []GolangChainedCall {
	startPos = NextNoWhitespaceIndex(tokens, startPos+1)
	result := []GolangChainedCall{}
	for startPos < len(tokens) && tokens[startPos].Spelling == "." {
		nameTokenIndex := NextNoWhitespaceIndex(tokens, startPos+1)
		if nameTokenIndex >= len(tokens) {
			return result
		}
		nameToken := tokens[nameTokenIndex]
		opensIndex := NextNoWhitespaceIndex(tokens, nameTokenIndex+1)
		if opensIndex >= len(tokens) {
			return result
		}
		opens := tokens[opensIndex]
		if opens.Spelling != "(" {
			return result
		}
		close, closeIndex := FindMatchingParen(tokens[opensIndex:])
		if close == nil {
			return result
		}
		startPos = NextNoWhitespaceIndex(tokens, opensIndex+closeIndex+1)
		parametersSpan := tokens[opensIndex+1 : opensIndex+closeIndex]

		parameters := SplitParenAware(parametersSpan, ",")
		for i := range parameters {
			parameters[i] = Strip(parameters[i])
		}
		result = append(result, GolangChainedCall{
			StartLine:      nameToken.BeginsLine,
			Name:           nameToken,
			ParametersSpan: parametersSpan,
			Parameters:     parameters,
		})
	}
	return result
}
