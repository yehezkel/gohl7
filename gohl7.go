package gohl7

// Function CleanRawValue return the escaped field free value of a simple field
// this function does create a copy of the value, if an in-place escape char replacement is
// need it (for saving memory) use RawValue+Encoding.Clean.
// this function could have receive a *SimpleField instead of a Field, but to free the caller
// of the type assertion we are doing it here.
func CleanRawValue(simple Field, enc *Encoding) ([]byte, error) {

	if simple == nil {
		return []byte{}, nil
	}

	s, ok := simple.(*SimpleField)
	if !ok {
		//error defined on utils.go
		return nil, ErrNonInternalImplementation
	}

	l, j := len(s.v), 0
	escapedValue := make([]byte, l)

	for i := 0; i < l; j++ {

		if s.v[i] == enc.Escaping {
			i++
			//not checking for out of bounds because this will imply a wrongly formatted
			//field, which should have be cought be the parser
		}

		escapedValue[j] = s.v[i]
		i++
	}

	return escapedValue[:j], nil
}

// Function RawValue gives you raw byte of a simple field.
// this function could have been a *SimpleField method but to discorage
// it overuse, it has been done similar to CleanRawValue.
func RawValue(simple *SimpleField) []byte {

	if simple == nil {
		return []byte{}
	}

	return simple.v
}

// Function ContainerFieldChildren returns the list of children from a ContainerField
// Given the fact this function is tied to the ComplexField implementation, it could have
// have been implemented as a method but to keep things consistent with SimpleField is done
// as a separte method
func ContainerFieldChildren(container Field) ([]Field, error) {

	if container == nil {
		return nil, nil
	}

	complexF, ok := container.(*ComplexField)
	if !ok {
		//error defined on utils.go
		return nil, ErrNonInternalImplementation
	}

	return complexF.children, nil

}
