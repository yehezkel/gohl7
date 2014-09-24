package gohl7

import(
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

type Hl7DataType interface{
	Field(index int) (Hl7DataType, bool)
}

type Hl7ComposedType interface{
	Hl7DataType
	AppendValue(v Hl7DataType) (error) 
}

type SimpleField 	  []byte
type SubComponent	  []SimpleField
type Component  	  []Hl7DataType
type Repeated  	  	  []Hl7DataType
type Segment 		  []Hl7DataType


func (simple SimpleField) Field(index int) (Hl7DataType, bool){
	if index != 0{
		return nil, false
	}

	return simple, true
}

func (s SubComponent) Field(index int) (Hl7DataType, bool){
	l := len(s)
	
	if index < 0 || index >= l{
		return nil, false
	}

	return s[index], true
}

func (s *SubComponent) AppendValue(v Hl7DataType) (error){

	value, ok := v.(SimpleField)

	if !ok{
		return errSubComponentType
	}

	*s = append(*s,value)

	return nil
}

func (c Component) Field(index int) (Hl7DataType, bool){
	l := len(c)
	
	if index < 0 || index >= l{
		return nil, false
	}

	return c[index], true
}

func (c *Component) AppendValue(v Hl7DataType) (err error){

	switch v.(type){
		case SimpleField, SubComponent: err = nil
		default: err = errComponentType 
	}

	if err != nil{
		return
	}

	*c = append(*c,v)

	return nil
}

func (r Repeated) Field(index int) (Hl7DataType, bool){
	l := len(r)
	
	if index < 0 || index >= l{
		return nil, false
	}

	return r[index], true
}

func (r *Repeated) AppendValue(v Hl7DataType) (err error){

	switch v.(type){
		case SimpleField, SubComponent, Component: 
			err = nil
		default: err = errRepeatType 
	}

	if err != nil{
		return
	}

	*r = append(*r,v)

	return nil
}

func (s Segment) Field(index int) (Hl7DataType, bool){
	l := len(s)
	
	if index < 0 || index >= l{
		return nil, false
	}

	return s[index], true
}

func (s *Segment) AppendValue(v Hl7DataType) (err error){

	switch v.(type){
		case SimpleField, SubComponent, Component, Repeated: 
			err = nil
		default: err = errRepeatType 
	}

	if err != nil{
		return
	}

	*s = append(*s,v)

	return nil
}