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
		return nil, err
	}

	segments := NewComplexField(segment, SegmentValidator)

	return &Hl7Parser{
		enc: enc,
		mssg: &Message{
			raw:          source,
			ComplexField: segments,
		},
	}, nil

}

func (p *Hl7Parser) Parse() (*Message, error) {

	mssg := p.mssg
	//manually adding header
	err := mssg.Push(NewSimpleField(mssg.raw[:3]))
	if err != nil {
		return nil, err
	}

	//manually adding encoding
	err = mssg.Push(NewSimpleField(mssg.raw[4:8]))
	if err != nil {
		return nil, err
	}

	last := segment
	i, l := 9, len(mssg.raw)

	for ; i < l;  {

		nextF, consumed, err := next(mssg.raw[i:],p.enc)
		if err != nil {

			if err != errNoMoreData {
				return nil, err	
			}
			//treat end of data as simple field
			nextF = segment
			
		}

		value := mssg.raw[i:i+consumed]
		i = i + consumed

		switch {
		case last == segment && nextF == segment:
			err = mssg.Push(NewSimpleField(value))
		}


	}



	return nil, nil
}

func next(source []byte, enc *Encoding) (FieldType, int, error) {

	l,k := len(source),0


	for ; k < l; k++ {

		v := source[k]
		switch v {
		case enc.Field:
			return Simple, k, nil
		case enc.Component:
			return Component, k, nil
		case enc.Repeated:
			return Component, k, nil
		case enc.Component:
			return SubComponent, k, nil
		case CR:
			return segment, k, nil
		case enc.Escaping:
			//continue
			k++
		}
	}

	return 0, k + 1, errNoMoreData
}
