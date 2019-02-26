package gohl7

import (
    "errors"
)

var (
    errNoMatchingTypes = errors.New("No matching type")
    errNonInternalImplementation = errors.New("Provided fields do not match internal implementation")
    errNonMatchingChildren = errors.New("Different amount of childrens")
)

//function newSimpleStr is an internal utility function
//mostly for testing
func newSimpleStr(s string) *SimpleField {

    return NewSimpleField([]byte(s))
}


// function newComplexFieldWithChildren is an internal utility function
// mostly for testing
//error code is ignore on purpose as its assume the call is internal
func newComplexFieldWithChildren(t FieldType, v func(Field, Field) error, children ...Field) *ComplexField {

    complexF := NewComplexField(t,v)

    for _, child := range children {
        _ = complexF.Push(child)
    }

    return complexF
}

//function deepEqual compare two fields recursivly for equality
//a return value of nil implies equal a non nil implies difference
func deepEqual(l,r Field) (error) {

    if l.Type() != r.Type() {
        return errNoMatchingTypes
    }

    if l.Type() == Simple {
        return nil
    }

    complexL,okL := l.(*ComplexField)
    complexR,okR := r.(*ComplexField)

    if !okL || !okR {
        return errNonInternalImplementation
    }

    childrenL, childrenR := complexL.children, complexR.children
    if len(childrenL) != len(childrenR) {
        return errNonMatchingChildren
    }

    for i:=0; i < len(childrenL); i++ {
        err := deepEqual(childrenL[i],childrenR[i])
        if err != nil {
            return err
        }
    }

    return nil
}