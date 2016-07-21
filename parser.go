package gohl7

import (
	"bufio"
	"errors"
	"io"
)

const (
	HEADER_LABEL = "MSH"
	CR           = '\r'
	NL           = '\n'
)

var (
	errUnexpectedToken = errors.New("After Segment Header only field separator")
	errMssgHeader      = errors.New("Invalid Message Header")
	errHeaderLength    = errors.New("Invalid Segment Header Length")
	errEscape          = errors.New("Invalid Escape Character")
	errBadEncoding     = errors.New("Invalid Encoding")
)

//the actual parser, and interesting point to notice is that even that the HL7 specification states that only '\r' is allow as segment separator
//in the real world HL7 messages are modified on windows or linux so the parser suports \n, \r and \r\n as segment separator.
type Parser struct {
	scanner *bufio.Scanner

	field        byte
	component    byte
	repeated     byte
	escape       byte
	subcomponent byte

	last    byte
	current byte

	segments []Segment
	sgmt     Segment
	rep      Repeated
	cmp      Component
	scmp     SubComponent
}

func NewParser(r io.Reader) (*Parser, error) {

	l := 5
	buffer := make([]byte, l)

	_, err := r.Read(buffer[:4])
	if err != nil {
		return nil, err
	}

	if string(buffer[:3]) != HEADER_LABEL {
		return nil, errMssgHeader
	}

	field_separator, i := buffer[3], 0

	for ; i < l; i++ {
		_, err = r.Read(buffer[i : i+1])
		if err != nil {
			return nil, err
		}

		if buffer[i] == field_separator {
			break
		}
	}

	if i == l && buffer[4] != field_separator {
		return nil, errBadEncoding
	}

	//reset to 0 unused bytes
	for ; i < l; i++ {
		buffer[i] = 0
	}

	return &Parser{
		scanner:      bufio.NewScanner(r),
		field:        field_separator,
		component:    buffer[0],
		repeated:     buffer[1],
		escape:       buffer[2],
		subcomponent: buffer[3],
	}, nil
}

func (p *Parser) Parse() ([]Segment, error) {

	encoding, err := encodingToField(
		p.field, p.component,
		p.repeated, p.escape,
		p.subcomponent,
	)

	if err != nil {
		return nil, err
	}

	header := Segment([]Hl7DataType{
		SimpleField(HEADER_LABEL),
		encoding,
	})

	p.sgmt = header
	p.last = p.field
	p.current = p.field

	p.scanner.Split(p.split)

	for p.scanner.Scan() {
		switch p.current {
		case NL:
			if p.last != CR {
				err = p.appendSegment()
			}
		case CR:
			err = p.appendSegment()
		case p.field:
			err = p.appendField()
		case p.component:
			err = p.appendComponent()
		case p.subcomponent:
			err = p.appendSubComponent()
		case p.repeated:
			err = p.appendRepeated()
		}

		if err != nil {
			return nil, err
		}
	}

	return p.segments, p.scanner.Err()
}

func (p *Parser) appendSegment() (err error) {

	if err = p.appendField(); err != nil {
		return err
	}

	p.segments = append(p.segments, p.sgmt)
	p.sgmt = nil
	return err
}

func (p *Parser) appendField() (err error) {

	var value Hl7DataType

	if p.rep != nil {
		if err = p.appendRepeated(); err != nil {
			return err
		}

		value = p.rep
		p.rep = nil
	} else {

		if p.last == p.component ||
			p.last == p.subcomponent {

			if err = p.appendComponent(); err != nil {
				return err
			}

			value = p.cmp
			p.cmp = nil
		} else {
			value = p.simpleField()
		}
	}

	return p.sgmt.AppendValue(value)
}

func (p *Parser) appendRepeated() (err error) {
	var value Hl7DataType

	if p.last == p.field || p.last == p.repeated {
		value = p.simpleField()
	} else {
		err = p.appendComponent()
		if err != nil {
			return err
		}

		value = p.cmp
		p.cmp = nil
	}

	return p.rep.AppendValue(value)
}

func (p *Parser) appendComponent() (err error) {
	var value Hl7DataType

	if p.scmp != nil {
		err = p.appendSubComponent()
		if err != nil {
			return err
		}

		value = p.scmp
		p.scmp = nil
	} else {
		value = p.simpleField()
	}

	return p.cmp.AppendValue(value)
}

func (p *Parser) appendSubComponent() (err error) {

	return p.scmp.AppendValue(
		p.simpleField(),
	)
}

func (p *Parser) simpleField() SimpleField {

	raw := p.scanner.Bytes()
	s := make(SimpleField, len(raw))
	copy(s, raw)

	return s
}

func (p *Parser) split(data []byte, atEOF bool) (advance int, token []byte, err error) {

	i, escape := 0, false

OuterLoop:
	for ; i < len(data); i++ {
		switch data[i] {
		case p.escape:
			if (i == len(data)-1 && atEOF) ||
				(i < (len(data)-1) &&
					data[i+1] != p.escape &&
					data[i+1] != p.field &&
					data[i+1] != p.component &&
					data[i+1] != p.subcomponent &&
					data[i+1] != p.repeated) {
				return advance, token, errEscape
			}

			i++
			escape = true
		case NL:
			if p.current == NL ||
				(p.current == CR && i != 0) {
				return advance, token, errUnexpectedToken
			}
			break OuterLoop
		case p.field:
			if (p.current == NL || p.current == CR) &&
				i != 3 {
				return advance, token, errHeaderLength
			}
			break OuterLoop
		case p.component, p.repeated, p.subcomponent, CR:
			if p.current == NL || p.current == CR {
				return advance, token, errUnexpectedToken
			}
			break OuterLoop
		}
	}

	found := (i < len(data))

	if atEOF {
		token = data
		advance = len(data)
	}

	if found {
		token = data[:i]
		advance = i + 1
	}

	if atEOF || found {
		p.last = p.current
		p.current = CR

		if found {
			p.current = data[i]
		}

		if !escape {
			return advance, token, err
		}

		j, s := 0, 0
		for ; j < i; j, s = j+1, s+1 {
			if data[j] == p.escape {
				j++
			}
			data[s] = data[j]
		}

		return advance, data[:s], err
	}

	return advance, token, err
}

func encodingToField(field, component, repeated, escape, subcomponent byte) (SimpleField, error) {

	tmp := []byte{field, component, repeated, escape, subcomponent}

	//remove not specified encoding bytes
	for i := 0; i < 5; i++ {
		if tmp[i] == 0 {
			tmp = tmp[:i]
			break
		}
	}

	if len(tmp) == 0 {
		//this error is so naive no need to add to global scope
		return nil, errors.New("Invalid encoding  bytes")
	}

	//check that all the encoding bytes are unique
	for i := 0; i < len(tmp); i++ {
		for j := i + 1; j < len(tmp); j++ {
			if tmp[i] == tmp[j] {
				return nil, errBadEncoding
			}
		}
	}

	for _, e := range tmp {
		if e == NL || e == CR {
			return nil, errBadEncoding
		}
	}

	return SimpleField(tmp[1:]), nil
}
