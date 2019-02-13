package gohl7

/*import (
	"errors"
)

var (
	errSubComponentType = errors.New(
		"Hl7 SubComponent may only contain SimpleFields",
	)
	errComponentType = errors.New(
		"Hl7 Component may only contain SimpleFields and SubComponents",
	)
	errRepeatType = errors.New(
		"Hl7 Repeated may only contain SimpleFields, SubComponents, Component",
	)
)

type Hl7DataType interface {
	Field(index int) (Hl7DataType, bool)
}

type Hl7ComposedType interface {
	Hl7DataType
	AppendValue(v Hl7DataType) error
}

type SimpleField struct {
	value []byte
}

type SubComponent struct {
	fields []Hl7DataType
}

type Component struct {
	fields []Hl7DataType
}

type Repeated struct {
	fields []Hl7DataType
}

type Segment struct {
	fields []Hl7DataType
}

func NewSimpleField(value []byte) *SimpleField {
	return &SimpleField{
		value: value,
	}
}

func (simple *SimpleField) Field(index int) (Hl7DataType, bool) {

	if index != 0 {
		return nil, false
	}

	return simple, true
}

func (simple *SimpleField) String() string {
	return string(simple.value)
}

func (s *SubComponent) Field(index int) (Hl7DataType, bool) {
	l := len(s.fields)

	if index < 0 || index >= l {
		return nil, false
	}

	return s.fields[index], true
}

func (s *SubComponent) AppendValue(v Hl7DataType) (err error) {

	_, ok := v.(*SimpleField)

	if !ok {
		return errSubComponentType
	}

	s.fields = append(s.fields, v)

	return
}

func (c *Component) Field(index int) (Hl7DataType, bool) {
	l := len(c.fields)

	if index < 0 || index >= l {
		return nil, false
	}

	return c.fields[index], true
}

func (c *Component) AppendValue(v Hl7DataType) (err error) {

	switch v.(type) {
	case *SimpleField, *SubComponent:
		err = nil
	default:
		return errComponentType
	}

	c.fields = append(c.fields, v)

	return
}

//Hl7DataType Field method implementation
func (r *Repeated) Field(index int) (Hl7DataType, bool) {
	l := len(r.fields)

	if index < 0 || index >= l {
		return nil, false
	}

	return r.fields[index], true
}

//Hl7ComposedType AppendValue implementation
func (r *Repeated) AppendValue(v Hl7DataType) (err error) {

	switch v.(type) {
	case *SimpleField, *SubComponent, *Component:
		err = nil
	default:
		return errRepeatType
	}

	r.fields = append(r.fields, v)

	return
}

//Hl7DataType Field method implementation
func (s *Segment) Field(index int) (Hl7DataType, bool) {
	l := len(s.fields)

	if index < 0 || index >= l {
		return nil, false
	}

	return s.fields[index], true
}

//Hl7ComposedType AppendValue implementation
func (s *Segment) AppendValue(v Hl7DataType) (err error) {

	switch v.(type) {
	case *SimpleField, *SubComponent, *Component, *Repeated:
		err = nil
	default:
		return errRepeatType
	}

	s.fields = append(s.fields, v)

	return
}*/
