package gohl72

import(
	"errors"
)

const (
	HEADER_LABEL = "MSH"
	CR 			 = '\r'
	NL 			 = '\n'
)

var (
	errUnexpectedToken = errors.New("After Segment Header only field separator")
	errMssgHeader      = errors.New("Invalid Message Header")
	errHeaderLength    = errors.New("Invalid Segment Header Length")
	errEscape          = errors.New("Invalid Escape Character")
)

type Parser struct{
	buffer []byte
	r int

	encoding *Encoding
	err error

	last 			byte
	current			byte

	segments       []*Segment
	sgmt 		   *Segment
	rep            Repeated
	cmp			   Component
	scmp	       SubComponent
}

func NewParser(buffer []byte) (*Parser, error){


	if len(buffer) < 3 || string(buffer[:3]) != HEADER_LABEL{
		return nil, errMssgHeader
	}

	encoding, advance, err := newEncoding(buffer[3:])

	if err != nil{
		return nil,err
	}

	return &Parser{
		buffer: buffer,
		r: 3 + advance,
		encoding: encoding,	
	}, nil	
}

func(p *Parser) Parse()([]*Segment, error){
	if p.err != nil{
		return nil, p.err
	}

	enc, _ := p.encoding.ToSimpleField()

	if p.err != nil{
		return nil, p.err
	}

	p.sgmt = &Segment{
		fields: []Hl7DataType{&SimpleField{value:[]byte(HEADER_LABEL)}, enc},
	}

	p.last = p.encoding.Field
	p.current = p.encoding.Field

	p.segments = append(p.segments, p.sgmt)

	return p.segments, nil
}

func(p *Parser) Encoding() (*Encoding, error){
	if p.err != nil{
		return nil, p.err
	}

	return p.encoding, nil	
}

