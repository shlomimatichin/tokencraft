package tokencraft

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestAllMethods(t *testing.T) {
	code := `package tokencraft

type Eater struct {
	data   string
	offset int
	line   int
	column int
}

func newEater(data string) *Eater {
	return &Eater{data, 0, 1, 1}
}

func (e Eater) nonPointer() {
}

func (e *Eater) done() bool {
	return e.offset >= len(e.data)
}

func (e *Eater) current() byte {
	return e.data[e.offset]
}

func (e *Eater) next() (string, error) {
}
`
	methods := GolangAllMethods(Tokenize(code, HASH_IS_SPECIAL))
	assert.Assert(t, is.Len(methods, 4))
	assert.Assert(t, is.Equal(false, methods[0].TypePointer))
	assert.Assert(t, is.Equal(true, methods[1].TypePointer))
	assert.Assert(t, is.Equal(true, methods[2].TypePointer))
	assert.Assert(t, is.Equal(true, methods[3].TypePointer))
	assert.Assert(t, is.Equal(methods[0].Name, "nonPointer"))
	assert.Assert(t, is.Equal(methods[1].Name, "done"))
	assert.Assert(t, is.Equal(methods[2].Name, "current"))
	assert.Assert(t, is.Equal(methods[3].Name, "next"))
	assert.Assert(t, is.Contains(JoinSpellings(methods[2].BodySpan, ""), "e.data[e.offset]"))
}

func TestChain(t *testing.T) {
	code := `a.b().c(d(0), f(1, 2)).g()`
	chain := GolangCallChain(Tokenize(code, HASH_IS_SPECIAL), 0)
	assert.Assert(t, is.Len(chain, 3))
	assert.Assert(t, is.Equal(chain[0].Name.Spelling, "b"))
	assert.Assert(t, is.Equal(chain[1].Name.Spelling, "c"))
	assert.Assert(t, is.Equal(chain[2].Name.Spelling, "g"))
}
