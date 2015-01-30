package gohl7

import (
	"errors"
)

const (
	HEADER_LABEL      = "MSH"
	CR           byte = '\r'
	NL           byte = '\n'
)

var (
	errUnexpectedToken = errors.New("After Segment Header only field separator")
	errMssgHeader      = errors.New("Invalid Message Header")
	errHeaderLength    = errors.New("Invalid Segment Header Length")
	errEscape          = errors.New("Invalid Escape Character")
	eod                = errors.New("end of data")
)

type Parser struct {
	buffer []byte
	r      int
	c      int

	encoding *Encoding
	err      error

	last    byte
	current byte

	segments []*Segment
	sgmt     *Segment
	rep      *Repeated
	cmp      *Component
	scmp     *SubComponent
}

func NewParser(buffer []byte) (*Parser, error) {

	if len(buffer) < 3 || string(buffer[:3]) != HEADER_LABEL {
		return nil, errMssgHeader
	}

	encoding, advance, err := newEncoding(buffer[3:])

	if err != nil {
		return nil, err
	}

	return &Parser{
		buffer:   buffer,
		r:        3 + advance,
		c:        3 + advance,
		encoding: encoding,
	}, nil
}

func (p *Parser) Parse() ([]*Segment, error) {

	if p.err == eod {
		return p.segments, nil
	}

	if p.err != nil {
		return nil, p.err
	}

	enc, _ := p.encoding.ToSimpleField()

	if p.err != nil {
		return nil, p.err
	}

	p.sgmt = &Segment{
		fields: []Hl7DataType{&SimpleField{value: []byte(HEADER_LABEL)}, enc},
	}

	p.last = p.encoding.Field
	p.current = p.encoding.Field

	var token []byte
	var err error

	for {
		token, err = p.scan()
		if err != nil {
			break
		}

		simple := NewSimpleField(token)

		switch p.current {
		case NL:
			if p.last != CR {
				err = p.appendSegment(simple)
			}
		case CR:
			err = p.appendSegment(simple)
		case p.encoding.Field:
			err = p.appendField(simple)
		case p.encoding.Repeated:
			err = p.appendRepeated(simple)
		case p.encoding.Component:
			err = p.appendComponent(simple)
		case p.encoding.Subcomponent:
			err = p.appendSubComponent(simple)
		}

		if err != nil {
			break
		}

	}

	p.err = err
	if err == eod {
		err = nil
	}
	return p.segments, err
}

func (p *Parser) appendSegment(last Hl7DataType) (err error) {
	if err = p.appendField(last); err != nil {
		return err
	}

	p.segments = append(p.segments, p.sgmt)
	p.sgmt = nil

	return err
}

func (p *Parser) appendField(value Hl7DataType) (err error) {

	if p.rep != nil {

		err = p.appendRepeated(value)
		if err != nil {
			return err
		}
		value = p.rep
		p.rep = nil

	} else {

		if p.last == p.encoding.Component ||
			p.last == p.encoding.Subcomponent {

			err = p.appendComponent(value)
			if err != nil {
				return err
			}

			value = p.cmp
			p.cmp = nil
		}
	}

	if p.sgmt == nil {
		p.sgmt = new(Segment)
	}
	return p.sgmt.AppendValue(value)
}

func (p *Parser) appendRepeated(value Hl7DataType) (err error) {

	if p.last == p.encoding.Component ||
		p.last == p.encoding.Subcomponent {

		err = p.appendComponent(value)
		if err != nil {
			return err
		}

		value = p.cmp
		p.cmp = nil
	}

	if p.rep == nil {
		p.rep = new(Repeated)
	}

	return p.rep.AppendValue(value)
}

func (p *Parser) appendComponent(value Hl7DataType) (err error) {

	if p.scmp != nil {
		err = p.appendSubComponent(value)
		if err != nil {
			return err
		}

		value = p.scmp
		p.scmp = nil
	}

	if p.cmp == nil {
		p.cmp = new(Component)
	}

	return p.cmp.AppendValue(value)
}

func (p *Parser) appendSubComponent(value Hl7DataType) (err error) {

	if p.scmp == nil {
		p.scmp = new(SubComponent)
	}

	return p.scmp.AppendValue(value)
}

func (p *Parser) Encoding() (*Encoding, error) {
	if p.err != nil {
		return nil, p.err
	}

	return p.encoding, nil
}

func (p *Parser) scan() ([]byte, error) {

	r, c, l := p.r, p.c, len(p.buffer)
	buffer := p.buffer[:]

	if r > l {
		return nil, eod
	}

OuterLoop:
	for ; r < l; r, c = r+1, c+1 {

		buffer[c] = buffer[r]

		if !p.encoding.IsToken(buffer[r]) {
			continue
		}

		switch buffer[r] {
		case p.encoding.Field:
			if (p.current == NL || p.current == CR) &&
				r-p.r != 3 {
				return nil, errHeaderLength
			}
			break OuterLoop
		case p.encoding.Component, p.encoding.Repeated, p.encoding.Subcomponent, CR:
			if p.current == NL || p.current == CR {
				return nil, errUnexpectedToken
			}
			break OuterLoop
		case NL:
			if p.current == NL ||
				(p.current == CR && (r-p.r) != 0) {
				return nil, errUnexpectedToken
			}
			break OuterLoop
		case p.encoding.Escaping:
			next := r + 1
			if next == l ||
				!p.encoding.IsToken(buffer[next]) {
				return nil, errEscape
			}
			buffer[c] = buffer[next]
			r++
		}
	}

	found := CR
	if r < l {
		found = buffer[r]
	}

	p.r = r + 1
	token := buffer[p.c:c]
	p.c = c + 1

	p.last = p.current
	p.current = found
	return token, nil
}
