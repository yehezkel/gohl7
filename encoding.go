package gohl7

import (
	"errors"
)

var (
	errBadEncoding = errors.New("Invalid Encoding")
)

type Encoding struct {
	Field        byte
	Component    byte
	Repeated     byte
	Escaping     byte
	Subcomponent byte
}

func newEncoding(buffer []byte) (*Encoding, int, error) {

	l, i := len(buffer), 0

	if l == 0 {
		return nil, 0, errBadEncoding
	}

	encoding := new(Encoding)

	encoding.Field = buffer[i]

	i++
	if i == l {
		return nil, i, errBadEncoding
	}

	if buffer[i] == encoding.Field {
		return encoding, i + 1, nil
	}

	encoding.Component = buffer[i]

	i++
	if i == l {
		return nil, i, errBadEncoding
	}

	if buffer[i] == encoding.Field {
		return encoding, i + 1, nil
	}

	encoding.Repeated = buffer[i]

	i++
	if i == l {
		return nil, i, errBadEncoding
	}

	if buffer[i] == encoding.Field {
		return encoding, i + 1, nil
	}

	encoding.Escaping = buffer[i]

	i++
	if i == l {
		return nil, i, errBadEncoding
	}

	if buffer[i] == encoding.Field {
		return encoding, i + 1, nil
	}

	encoding.Subcomponent = buffer[i]

	i++
	if i == l || buffer[i] != encoding.Field {
		return nil, i, errBadEncoding
	}

	return encoding, i + 1, nil
}

func (enc *Encoding) ToSimpleField() (*SimpleField, error) {
	tmp := []byte{
		enc.Field,
		enc.Component,
		enc.Repeated,
		enc.Escaping,
		enc.Subcomponent,
	}

	for i := 0; i < 5; i++ {
		if tmp[i] == 0 {
			tmp = tmp[:i]
			break
		}
	}

	//check that at leat the field encoding exist
	if len(tmp) == 0 {
		return nil, errBadEncoding
	}

	//check that all the encoding bytes are unique
	//check that no encoding character is NL or CR
	for i := 0; i < len(tmp); i++ {

		if tmp[i] == NL || tmp[i] == CR {
			return nil, errBadEncoding
		}

		for j := i + 1; j < len(tmp); j++ {
			if tmp[i] == tmp[j] {
				return nil, errBadEncoding
			}
		}
	}

	return &SimpleField{
		value: tmp[1:],
	}, nil
}

func (enc *Encoding) IsToken(b byte) bool {
	return (b == enc.Field ||
		b == enc.Escaping ||
		b == enc.Component ||
		b == enc.Repeated ||
		b == enc.Subcomponent ||
		b == CR ||
		b == NL)
}
