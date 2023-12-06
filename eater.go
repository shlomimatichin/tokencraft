package tokencraft

type Eater struct {
	data   string
	offset int
	line   int
	column int
}

func newEater(data string) *Eater {
	return &Eater{data, 0, 1, 1}
}

func (e *Eater) done() bool {
	return e.offset >= len(e.data)
}

func (e *Eater) current() byte {
	return e.data[e.offset]
}

func (e *Eater) next() byte {
	var c byte = 0
	if e.offset < len(e.data)-1 {
		c = e.data[e.offset+1]
	}
	return c
}

func (e *Eater) advance() {
	c := e.data[e.offset]
	if c == '\n' {
		e.line++
		e.column = 1
	} else {
		e.column++
	}
	e.offset++
}
