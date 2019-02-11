package gohl7

import (
	"errors"
)

var (
	ErrEmptyChildren = errors.New("No child to pop off")
)

type FieldType byte

const (
	FieldType Simple = iota
	Repeated
	Component
	SubComponent
)

type Field interface {
	Type() FieldType
}

type ContainerField interface {
	Field
	Pop() (Field, error)
	Push(Field) error
}

type validator interface {
	ValidatePush(parent, child ContainerField) error
}

type SimpleField struct {
	v []byte
}

func (v *SimpleField) Type() FieldType {
	return Simple
}

type ComplexField struct {
	fieldType FieldType
	children  []Field
	validator func(Field, Field) error
}

func (f *ComplexField) Type() FieldType {

	return f.fieldType
}

func (f *ComplexField) Pop() (Field, error) {

	l := len(f.children) - 1
	if l < 0 {
		return nil, ErrEmptyChildren
	}

	last = f.children[l]
	f.children = f.children[:l]

	return last, nil
}

func (f *ComplexField) Push(child Field) (err error) {

	if f.validator != nil && err = f.validator(f, child); err {
		return err
	}

	f.children = append(f.children, child)

	return
}
