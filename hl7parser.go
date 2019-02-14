package gohl7

import (
	"errors"
	//"log"
)

var (
	ErrMssgHeader     = errors.New("Invalid Message Header")
	ErrMssgEncoding   = errors.New("Invalid message encoding field")
	ErrUnexpectedCase = errors.New("Unexpected case, implementation error")
	errNoMoreData     = errors.New("no more data")
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

	segments := NewComplexField(message, MessageValidator)

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
	currentSegment := NewComplexField(segment, SegmentValidator)

	//manually adding header
	err := currentSegment.Push(NewSimpleField(mssg.raw[:3]))
	if err != nil {
		return nil, err
	}

	//manually adding encoding
	err = currentSegment.Push(NewSimpleField(mssg.raw[4:8]))
	if err != nil {
		return nil, err
	}

	last := Simple
	i, l := 9, len(mssg.raw)

	for i < l {

		nextF, consumed, err := next(mssg.raw[i:], p.enc)

		if err != nil {

			if err != errNoMoreData {
				return nil, err
			}
			//treat end of data as segment
			nextF = segment

		}

		value := mssg.raw[i : i+consumed]

		i = i + consumed + 1

		switch {
		case last == segment && nextF == Simple:
			//segment header case
			fallthrough

		case last == Simple && nextF == Simple:
			err = currentSegment.Push(NewSimpleField(value))

		case last == Simple && nextF == segment:
			//append last part of segment
			err = currentSegment.Push(NewSimpleField(value))
			//add current segment to the message
			mssg.Push(currentSegment)
			//create new segment
			currentSegment = NewComplexField(segment, SegmentValidator)

		case last == Simple && nextF == Component:
			//create complex field component
			complexF := NewComplexField(Component, ComponentValidator)
			//push component into segment
			currentSegment.Push(complexF)
			//append simple field
			complexF.Push(NewSimpleField(value))

		case last == Simple && nextF == SubComponent:

			//create complex field component
			componentF := NewComplexField(Component, ComponentValidator)
			//push component into segment
			currentSegment.Push(componentF)
			//create comples field subcomponent
			subcomponentF := NewComplexField(SubComponent, SubComponentValidator)
			//push subcomponent into component
			componentF.Push(subcomponentF)
			//push simple field into subcomponent
			subcomponentF.Push(NewSimpleField(value))

		case last == Component && nextF == Simple:
			fallthrough
		case last == Component && nextF == Component:
			err = pushChildToLastChild(currentSegment, NewSimpleField(value))
		case last == Component && nextF == segment:
			err = pushChildToLastChild(currentSegment, NewSimpleField(value))
			//add current segment to the message
			mssg.Push(currentSegment)
			//create new segment
			currentSegment = NewComplexField(segment, SegmentValidator)

		case last == segment && nextF == segment:
			//special case for CR+EndOfData
			//noop
		default:
			err = ErrUnexpectedCase
		}

		if err != nil {
			return nil, err
		}

		last = nextF

	}

	return mssg, nil
}

func next(source []byte, enc *Encoding) (FieldType, int, error) {

	l, k := len(source), 0

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

//Given a complex field append a new child to its last child
//in other word: add a grandchild
//in other words: push child to last child
func pushChildToLastChild(parent *ComplexField, newChild Field) error {

	lastField, err := parent.Pop()
	if err != nil {
		return err
	}

	//converting last field back to *ComplexField
	complexF, ok := lastField.(*ComplexField)
	if !ok {
		return ErrUnexpectedCase
	}

	//push new child
	err = complexF.Push(newChild)
	if err != nil {
		return err
	}

	//push back the last complex field
	return parent.Push(complexF)
}
