package tokencraft

import (
	"os"
	"path/filepath"
)

type TokenizedFile struct {
	Tokens []Token
	Filename string
	Basename string
}

func TokenizeFile(filename string, hashMode HashMode) (*TokenizedFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &TokenizedFile{
		Tokens: Tokenize(string(data), hashMode),
		Filename: filename,
		Basename: filepath.Base(filename),
	}, nil
}
