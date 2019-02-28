package gohl7

import (
	"errors"
)

var (
	errBadEncoding      = errors.New("Invalid Encoding")
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
