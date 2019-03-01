package gohl7

import (
	"errors"
)

var (
	ErrBadEncoding      = errors.New("Invalid Encoding")
	ErrRepeatedEncoding = errors.New("Invalid Encoding, repeated chars")
)

type Encoding struct {
	Field        byte
	Component    byte
	Repeated     byte
	Escaping     byte
	Subcomponent byte
}

func ParseEncoding(buffer []byte) (*Encoding, error) {

	l := len(buffer)
	if l != 5 {
		return nil, ErrBadEncoding
	}

	e := &Encoding{
		Field:        buffer[0],
		Component:    buffer[1],
		Repeated:     buffer[2],
		Escaping:     buffer[3],
		Subcomponent: buffer[4],
	}

	//checking for duplicates
	bag := make(map[byte]bool)
	for _, v := range buffer {
		_, ok := bag[v]
		if ok {
			return nil, ErrRepeatedEncoding
		}

		bag[v] = true
	}

	return e, nil

}

// Function Clean removes all unescape/clean the given input
// the process is done in-place no buffer is made
func (enc *Encoding) Clean(input []byte) []byte {

	l, j := len(input), 0

	for i := 0; i < l; j++ {

		if input[i] == enc.Escaping {
			i++
			//not checking for out of bounds because this will imply a wrongly formatted
			//field, which should have be cought be the parser
		}

		input[j] = input[i]
		i++
	}

	return input[:j]
}
