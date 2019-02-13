package gohl7

import (
	"errors"
)

var (
	ErrMssgHeader   = errors.New("Invalid Message Header")
	ErrMssgEncoding = errors.New("Invalid message encoding field")
	errNoMoreData   = errors.New("no more data")
)

const (
	HEADER_LABEL      = "MSH"
	CR           byte = '\r'
	NL           byte = '\n'
)

type Hl7Parser struct {
	enc  *Encoding
	mssg *Message
}

func NewHl7Parser(source []byte) (*Hl7Parser, error) {

	l := len(source)

	if l < 3 || string(source[:3]) != HEADER_LABEL {
		return nil, ErrMssgHeader
	}

	if l < 8 {
		return nil, ErrMssgEncoding
	}

	enc, err := ParseEncoding(source[3:8])
	if err != nil {
		return nil.err
	}

	return &Hl7Parser{
		enc: enc,
		mssg: &Message{
			raw: source,
		},
	}

}

func (p *Hl7Parser) Parse() (Message, error) {

	return nil, nil
}

func next(source []byte, enc *Encoding) (FieldType, int, err) {

	l := len(source)

	for k = 0; k < l; k++ {

		v = source[k]
		switch v {
		case enc.Field:
			return Simple, k, nil
		case enc.Component:
			return Component, k, nil
		case enc.Repeated:
			return Component, k, nil
		case enc.Component:
			return Subcomponent, k, nil
		case CR:
			return Segment, k, nil
		case enc.Escaping:
			//continue
			k++
		}
	}

	return 0, k + 1, errNoMoreData
}
