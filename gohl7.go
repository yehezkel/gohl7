//A single pass parser for HL7 version 2.X.X.
//No usage of regular expresions to extract the values.
//The package supports all the separators including repeated fields and different encoding
//characters.
//
//Currently non-ascii HL7 messages are not fully supported due to the lack
//of documentation and message samples, if you can provide some information
//it will be more than welcome.
//
//To properly understand the package i suggest to refer first to the main interfaces:
//Hl7DataType and Hl7ComposeType, then take a look at the types that implement them:
//SimpleField, SubComponent, Component, Repeated, Segment. Then you can see the Parser.
//
//For examples of usage checkout the code and take a look at the gohl7_test file, particularly the
//first test function
//
//Design Decisions
//
//After inspecting the code you may find the same implementation of the Field function  for the SubComponent,
//Component, Repeated, and Segment, i was able to implement thouse types as follow to avoid the code repetition:
//		type ComposedType []Hl7DataType
//		type SubComponent struct{ComposedType}
//		type Component struct{ComposedType}
//
//		func(c ComposedType) Field(index int)(Hl7DataType, bool){
//			//only implement the function once
//		}
//The above implemention will force me to do the New.. or Make.. function for each of the types or
//force the package clients to start doing:
//		r := &SubComponent{...to verbose...} or
//		r := new(Component) not bad but then you lose some extra features of having and slice as underlying type as len or indexing
//Instead i decided to duplicate the function implementation but be able to do:
//		var s Segment
//		s.AppendValue(..some value..)
//		s.AppendValue(..other value..)
//		l := len(s) //or s[0]
//Also having an slice as underlying type allowes me the implement methods receiving pointers and not pointers and still implement an interface and be able to
//work not always with pointer variables, see how Component Field method has no pointer as a receiver but AppendField does, and still the type Component implement Hl7ComposedType
//interface
package gohl7

//The Hl7DataType is implemented by all the types
//including a simple value between field separators
//it returns the value contained at the index position by the HL7 field
//or nil and false in case the position does not exists
type Hl7DataType interface{
	Field(index int) (Hl7DataType, bool)
}


//The Hl7ComposedType is only implemented by the HL7 types that can actually
//contain other data types inside them: Subcomponent, Component, Repeated and Segment.
//It provides the functionality of appending a new value into them
//the method may return a validation error, as an example when
//a repeated type is been inserted inside a subcomponent data type
type Hl7ComposedType interface{
	Hl7DataType
	AppendValue(v Hl7DataType) (error)
}

