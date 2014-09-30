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

//SimpleField represent the string values between the field separator
//as follows: PID|This this a simple field|
//The underlaying type of SimpleField is an slice of bytes and not strings
//to be able to modify its encoding after the parsing stage
type SimpleField 	  []byte

//SubComponent represent the values typically between & encoding field
//as an example: PID|subcomponent1&subcomponent2|
type SubComponent	  []SimpleField

//Component represent the values typically between ^ encoding field
//as an example: PID|Component1^Component2|
type Component  	  []Hl7DataType

//Repeated represent the values typically between ~ encoding field
//as an example: PID|Component1^Component2~Component11^Component22|
type Repeated  	  	  []Hl7DataType

//Segment is implemented as no more than an slice containing the types defined
//above
type Segment 		  []Hl7DataType



//It may not be clear why SimpleField must implement this method as it only contain
//a single value at index 0, the reason for this is because the HL7 specification
//states that empty fields may be excluded when generating the HL7 so
// a value like this: PID|123&^~ may be and should be written as PID|123
//so order to maintain an standard interface between the hl7 consummers that
//suport more values on the segments the SimpleField return itself while we
//are trying to get the value at index 0
func (simple SimpleField) Field(index int) (Hl7DataType, bool){
	if index != 0{
		return nil, false
	}

	return simple, true
}

//Hl7DataType Field method implementation
func (s SubComponent) Field(index int) (Hl7DataType, bool){
	l := len(s)

	if index < 0 || index >= l{
		return nil, false
	}

	return s[index], true
}

//Hl7ComposedType AppendValue implementation
func (s *SubComponent) AppendValue(v Hl7DataType) (error){

	value, ok := v.(SimpleField)

	if !ok{
		return errSubComponentType
	}

	*s = append(*s,value)

	return nil
}

//Hl7DataType Field method implementation
func (c Component) Field(index int) (Hl7DataType, bool){
	l := len(c)

	if index < 0 || index >= l{
		return nil, false
	}

	return c[index], true
}

//Hl7ComposedType AppendValue implementation
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

//Hl7DataType Field method implementation
func (r Repeated) Field(index int) (Hl7DataType, bool){
	l := len(r)

	if index < 0 || index >= l{
		return nil, false
	}

	return r[index], true
}

//Hl7ComposedType AppendValue implementation
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

//Hl7DataType Field method implementation
func (s Segment) Field(index int) (Hl7DataType, bool){
	l := len(s)

	if index < 0 || index >= l{
		return nil, false
	}

	return s[index], true
}

//Hl7ComposedType AppendValue implementation
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