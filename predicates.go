package tokencraft

import (
	"strings"
)

var opens = map[string]string{
	"{": "}",
	"(": ")",
	"[": "]",
}

var closes = map[string]string{
	"}": "{",
	")": "(",
	"]": "[",
}

func FindAllSpelling(tokens []Token, spelling string) []Token {
	result := []Token{}
	for _, token := range tokens {
		if token.Spelling == spelling {
			result = append(result, token)
		}
	}
	return result
}

func FindNext(tokens []Token, spelling string) *Token {
	for _, token := range tokens {
		if token.Spelling == spelling {
			return &token
		}
	}
	return nil
}

func FindAllSpellings(tokens []Token, spellings []string) [][]Token {
	result := [][]Token{}
	for i, token := range tokens {
		if token.Spelling == spellings[0] {
			match := MatchIgnoreWhitespaces(tokens[i:], spellings)
			if match != nil {
				result = append(result, match)
			}
		}
	}
	return result
}

func FindAllSpellingsIndex(tokens []Token, spellings []string) [][]int {
	result := [][]int{}
	for i, token := range tokens {
		if token.Spelling == spellings[0] {
			match := MatchIgnoreWhitespacesIndex(tokens[i:], spellings)
			for j := range match {
				match[j] += i
			}
			if match != nil {
				result = append(result, match)
			}
		}
	}
	return result
}

func MatchIgnoreWhitespacesIndex(tokens []Token, spellings []string) []int {
	found := 0
	location := 0
	result := []int{}
	for found < len(spellings) && location < len(tokens) {
		if tokens[location].Type == WHITESPACE {
			location++
		} else if tokens[location].Spelling == spellings[found] {
			result = append(result, location)
			found++
			location++
		} else {
			return nil
		}
	}
	if found < len(spellings) {
		return nil
	}
	return result
}

func MatchIgnoreWhitespaces(tokens []Token, spellings []string) []Token {
	match := MatchIgnoreWhitespacesIndex(tokens, spellings)
	if match == nil {
		return nil
	}
	result := make([]Token, len(match))
	for i, index := range match {
		result[i] = tokens[index]
	}
	return result
}

func FindMatchingParen(tokens []Token) (*Token, int) {
	if len(tokens) == 0 {
		// No tokens
		return nil, -1
	}
	open := tokens[0]
	if _, ok := opens[open.Spelling]; !ok {
		// Not an open token
		return nil, -1
	}
	stack := []Token{open}
	for i := 1; i < len(tokens); i++ {
		candidate := tokens[i]
		if _, ok := opens[candidate.Spelling]; ok {
			stack = append(stack, candidate)
		} else if closeFor, ok := closes[candidate.Spelling]; ok {
			if len(stack) == 0 {
				// No open token
				return nil, -1
			}
			popped := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if popped.Spelling != closeFor {
				// Incoherent parenthesis
				return nil, -1
			}
			if len(stack) == 0 {
				return &candidate, i
			}
		}
	}
	// No close token
	return nil, -1
}

func DropWhitespaces(tokens []Token) []Token {
	result := []Token{}
	for _, token := range tokens {
		if token.Type != WHITESPACE {
			result = append(result, token)
		}
	}
	return result
}

func DropComments(tokens []Token) []Token {
	result := []Token{}
	for _, token := range tokens {
		if token.Type != COMMENT {
			result = append(result, token)
		}
	}
	return result
}

func JoinSpellings(tokens []Token, seperator string) string {
	result := []string{}
	for _, token := range tokens {
		result = append(result, token.Spelling)
	}
	return strings.Join(result, seperator)
}

func Split(tokens []Token, seperator string) [][]Token {
	result := [][]Token{}
	current := []Token{}
	for _, token := range tokens {
		if token.Spelling == seperator {
			result = append(result, current)
			current = []Token{}
		} else {
			current = append(current, token)
		}
	}
	result = append(result, current)
	return result
}

func SplitParenAware(tokens []Token, seperator string) [][]Token {
	result := [][]Token{}
	current := []Token{}
	skip := -1
	for i, token := range tokens {
		if i < skip {
			continue
		}
		if token.Spelling == seperator {
			result = append(result, current)
			current = []Token{}
		} else if _, ok := opens[token.Spelling]; ok {
			found, skipUntil := FindMatchingParen(tokens[i:])
			if found == nil {
				// unable to parse - matching paren missing
				return nil
			}
			skip = skipUntil + i
			current = append(current, tokens[i:skip]...)
		} else {
			current = append(current, token)
		}
	}
	result = append(result, current)
	return result
}

func Strip(tokens []Token) []Token {
	start := 0
	for start < len(tokens) && tokens[start].Type == WHITESPACE {
		start++
	}
	end := len(tokens) - 1
	for end >= 0 && tokens[end].Type == WHITESPACE {
		end--
	}
	return tokens[start : end+1]
}

func NextNoWhitespaceIndex(tokens []Token, index int) int {
	for index < len(tokens) && tokens[index].Type == WHITESPACE {
		index++
	}
	return index
}

func NextNoWhitespace(tokens []Token, index int) *Token {
	index = NextNoWhitespaceIndex(tokens, index)
	if index >= len(tokens) {
		return nil
	}
	return &tokens[index]
}
