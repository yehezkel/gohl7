package gohl7

import (
	"errors"
)

var (
	//ErrEmptyChildren = errors.New("No child to pop off")
	ErrSegmentChild      = errors.New("Invalid Segment Child")
	ErrInvalidSegmenId   = errors.New("Invalid Segment Id")
	ErrRepeatedChild     = errors.New("Invalid Repeated Child")
	ErrComponentChild    = errors.New("Invalid Component Child")
	ErrSubComponentChild = errors.New("Invalid SubComponent Child")
)

func SegmentValidator(parent, child Field) error {

	t := child.Type()

	if t != Simple && t != Repeated && t != Component {
		return ErrSegmentChild
	}

	//extra validation if Field interface implementation
	//are the one provided with this package
	parentField, okP := parent.(*ComplexField)
	childField, okC := parent.(*SimpleField)

	if !okP || !okC {
		return nil
	}

	if len(parentField.children) == 0 && len(childField.v) != 3 {
		return ErrInvalidSegmenId
	}

	return nil
}

func ComponentValidator(parent, child Field) error {

	t := child.Type()

	if t != Simple && t != SubComponent {
		return ErrComponentChild
	}

	return nil
}

func RepeatedValidator(parent, child Field) error {

	if child.Type() != Component {
		return ErrRepeatedChild
	}

	return nil
}

func SubComponentValidator(parent, child Field) error {

	if child.Type() != Simple {
		return ErrSubComponentChild
	}

	return nil
}
