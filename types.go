package gohl7

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyChildren = errors.New("No child to pop off")
)

type FieldType byte

const (
	Simple FieldType = iota
	Repeated
	Component
	SubComponent

	//segment and message are defined as Field type
	//to help generalizing the implementation of the parser
	//give that this not an actual field type its unexported
	segment
	message
)

type Field interface {
	Type() FieldType
}

type ContainerField interface {
	Field
	Pop() (Field, error)
	Push(Field) error
}

type SimpleField struct {
	v []byte
}

func NewSimpleField(v []byte) *SimpleField {

	return &SimpleField{
		v: v,
	}
}

func (s *SimpleField) Type() FieldType {
	return Simple
}

func (s *SimpleField) String() string {
	return string(s.v)
}

type ComplexField struct {
	fieldType FieldType
	children  []Field
	validator func(Field, Field) error
}

func NewComplexField(t FieldType, v func(Field, Field) error) *ComplexField {

	return &ComplexField{
		fieldType: t,
		validator: v,
	}
}

func (f *ComplexField) Type() FieldType {

	return f.fieldType
}

func (f *ComplexField) Pop() (Field, error) {

	l := len(f.children) - 1
	if l < 0 {
		return nil, ErrEmptyChildren
	}

	last := f.children[l]
	f.children = f.children[:l]

	return last, nil
}

func (f *ComplexField) Push(child Field) (err error) {

	if f.validator != nil {
		err = f.validator(f, child)
		if err != nil {
			return err
		}
	}

	f.children = append(f.children, child)

	return
}

func (f *ComplexField) String() string {

	return fmt.Sprintf("Field Type :%#v Childrens: %s\n", f.fieldType, f.children)
}

func IsSimpleField(f Field) bool {

	return f.Type() == Simple
}

func IsComplexField(f Field) bool {

	t := f.Type()

	//TODO: not sure if Segment should be added here
	return t == Component ||
		t == Repeated ||
		t == SubComponent
}

type Message struct {
	raw []byte
	//reusing the push pop logic for segments
	*ComplexField
}
