package gohl7

import (
	"errors"
)

var (
	ErrMssgHeader   = errors.New("Invalid Message Header")
	ErrMssgEncoding = errors.New("Invalid message encoding field")
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
