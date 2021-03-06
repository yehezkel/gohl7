package gohl7

import (
	"errors"
	//temporal debug location
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

	//the i <= l instead of i < l is to handle wrong end formated files that
	//do not end on \r Ex: MSH|^~\&|~
	for i <= l {
		//temporal debug location
		//log.Printf("len: %d current: %d left %s\n", l, i, mssg.raw[i:])
		nextF, consumed, err := next(mssg.raw[i:], p.enc)
		//temporal debug location
		//log.Printf("len: %d current: %d consumed: %d err: %s\n", l, i, consumed, err)
		if err != nil {

			if err != errNoMoreData {
				return nil, err
			}
			//treat end of data as segment
			nextF = segment
			err = nil

		}

		value := mssg.raw[i : i+consumed]

		i = i + consumed + 1

		switch {
		case last == segment && nextF == Simple:
			//segment header case
			fallthrough

		case last == Simple && nextF == Simple:
			fallthrough

		case last == Simple && nextF == segment:
			//append last part of segment
			err = currentSegment.Push(NewSimpleField(value))

		case last == Simple && nextF == Component:
			//create complex field component
			complexF := NewComplexField(Component, ComponentValidator)
			//push component into segment
			err = currentSegment.Push(complexF)
			//append simple field
			if err == nil {
				complexF.Push(NewSimpleField(value))
			}

		case last == Simple && nextF == Repeated:

			//create complex field component
			complexF := NewComplexField(Repeated, RepeatedValidator)
			//push repeated into segment
			err = currentSegment.Push(complexF)
			//append simple field
			if err == nil {
				complexF.Push(NewSimpleField(value))
			}

		case last == Simple && nextF == SubComponent:

			//create complex field component
			componentF := NewComplexField(Component, ComponentValidator)
			//push component into segment
			err = currentSegment.Push(componentF)
			//create complex field subcomponent
			subcomponentF := NewComplexField(SubComponent, SubComponentValidator)
			//push subcomponent into component
			if err == nil {
				err = componentF.Push(subcomponentF)
			}
			//push simple field into subcomponent
			if err == nil {
				err = subcomponentF.Push(NewSimpleField(value))
			}

		case last == Repeated && nextF == segment:
			fallthrough

		case last == Repeated && nextF == Simple:
			fallthrough

		case last == Repeated && nextF == Repeated:
			err = pushChildToLastChild(currentSegment, NewSimpleField(value))

		case last == Repeated && nextF == Component:

			complexF := NewComplexField(Component, ComponentValidator)
			err = complexF.Push(NewSimpleField(value))
			if err == nil {
				err = pushChildToLastChild(currentSegment, complexF)
			}

		case last == Repeated && nextF == SubComponent:

			//create complex field component
			componentF := NewComplexField(Component, ComponentValidator)
			//create complex field subcomponent
			subcomponentF := NewComplexField(SubComponent, SubComponentValidator)
			//push subcomponent into component
			err = componentF.Push(subcomponentF)
			//push simple field into subcomponent
			if err == nil {
				err = subcomponentF.Push(NewSimpleField(value))
			}
			//push new component into Repeated field
			if err == nil {
				err = pushChildToLastChild(currentSegment, componentF)
			}

		case last == Component && nextF == segment:
			fallthrough

		case last == Component && nextF == Simple:
			fallthrough

		case last == Component && nextF == Component:

			var complexF *ComplexField
			complexF, err = getLastComplexChild(currentSegment)
			if err != nil {
				break
			}

			repeatedParent := (complexF.Type() == Repeated)

			if repeatedParent {
				err = pushChildToLastChild(complexF, NewSimpleField(value))
			} else {
				err = pushChildToLastChild(currentSegment, NewSimpleField(value))
			}

		case last == Component && nextF == SubComponent:

			var complexF *ComplexField

			complexF, err = getLastComplexChild(currentSegment)
			if err != nil {
				break
			}
			repeatedParent := (complexF.Type() == Repeated)

			if repeatedParent {
				//then the last child has to be the component
				complexF, err = getLastComplexChild(complexF)
			}

			//complex should reference the current component

			//create complex field subcomponent
			subcomponentF := NewComplexField(SubComponent, SubComponentValidator)
			//push subcomponent into component
			if err == nil {
				err = complexF.Push(subcomponentF)
			}

			//push simple field into subcomponent
			if err == nil {
				subcomponentF.Push(NewSimpleField(value))
			}

		case last == Component && nextF == Repeated:

			var complexF *ComplexField
			complexF, err = popLastComplexChild(currentSegment)
			if err != nil {
				break
			}

			//already a repeated field
			if complexF.Type() == Repeated {
				//put it back, no need to check for error
				currentSegment.Push(complexF)
				//then its last child has to be the component
				complexF, err = getLastComplexChild(complexF)

			} else {
				//build a new repeated field
				repeatedF := NewComplexField(Repeated, RepeatedValidator)
				//push existing component to it
				err = repeatedF.Push(complexF)

				//push repeated field to segment
				if err == nil {
					err = currentSegment.Push(repeatedF)
				}

			}
			//push to component the simple field
			if err == nil {
				complexF.Push(NewSimpleField(value))
			}

		case last == SubComponent && nextF == Simple:
			fallthrough

		case last == SubComponent && nextF == Component:
			fallthrough

		case last == SubComponent && nextF == SubComponent:
			fallthrough

		case last == SubComponent && nextF == segment:

			var complexF *ComplexField
			complexF, err = getLastComplexChild(currentSegment)
			if err != nil {
				break
			}

			repeatedParent := (complexF.Type() == Repeated)

			if repeatedParent {
				//then the last child has to be the component
				complexF, err = getLastComplexChild(complexF)
			}

			//complexF should reference the current component
			if err == nil {
				err = pushChildToLastChild(complexF, NewSimpleField(value))
			}

		case last == SubComponent && nextF == Repeated:

			var complexF *ComplexField
			complexF, err = popLastComplexChild(currentSegment)
			if err != nil {
				break
			}

			repeatedParent := (complexF.Type() == Repeated)

			//already on repeated field
			if repeatedParent {
				//put it back, no need to check for error
				currentSegment.Push(complexF)
				//then the last child has to be the component
				complexF, err = getLastComplexChild(complexF)

			} else {
				//build a new repeated field
				repeatedF := NewComplexField(Repeated, RepeatedValidator)
				//push existing component to it
				err = repeatedF.Push(complexF)
				//push repeated field to segment
				if err == nil {
					err = currentSegment.Push(repeatedF)
				}
			}

			//complexF should reference the current component
			if err == nil {
				err = pushChildToLastChild(complexF, NewSimpleField(value))
			}

		case last == segment && nextF == segment:
			//special case for CR+EndOfData
			//noop
		default:
			err = ErrUnexpectedCase
		}

		//case to handle new segments, excluding special case above.
		if err == nil && nextF == segment && last != segment {

			//add current segment to the message
			err = mssg.Push(currentSegment)
			//create new segment
			currentSegment = NewComplexField(segment, SegmentValidator)
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
			return Repeated, k, nil
		case enc.Subcomponent:
			return SubComponent, k, nil
		case CR:
			return segment, k, nil
		case enc.Escaping:
			//continue
			k++
		}
	}

	return 0, k, errNoMoreData
}

//Given a complex field append a new child to its last child
//in other word: add a grandchild
//in other words: push child to last child
func pushChildToLastChild(parent *ComplexField, newChild Field) error {

	complexF, err := getLastComplexChild(parent)
	if err != nil {
		return err
	}

	//push new child
	return complexF.Push(newChild)
}

func popLastComplexChild(parent *ComplexField) (*ComplexField, error) {

	lastField, err := parent.Pop()
	if err != nil {
		return nil, err
	}

	//converting last field back to *ComplexField
	complexF, ok := lastField.(*ComplexField)
	if !ok {
		return nil, ErrUnexpectedCase
	}

	return complexF, nil

}

//getLastComplexChild get the reference of the last complex child. helper function
func getLastComplexChild(parent *ComplexField) (*ComplexField, error) {

	//pop the child
	child, err := popLastComplexChild(parent)
	if err != nil {
		return child, err
	}

	//push the child back
	return child, parent.Push(child)
}
