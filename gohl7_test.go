package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	//"fmt"
	"testing"
)

func TestBadHeader(t *testing.T) {
	tests := []string{
		"M||",
		"||",
		"WRONG||",
		"",
	}

	for _, v := range tests {
		_, err := NewHl7Parser([]byte(v))
		if err == nil {
			t.Fatalf("Expecting error with header %s\n", v)
		}
	}
}

func TestSimpleFieldMessage(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&||\rTMP||"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"), newSimpleStr(""), newSimpleStr(""),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"), newSimpleStr(""), newSimpleStr(""),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa|bbbb\rTMP|123|456"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"), newSimpleStr("aaa"), newSimpleStr("bbbb"),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"), newSimpleStr("123"), newSimpleStr("456"),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestRepeatedFieldSimple(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|~"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr(""), newSimpleStr(""),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa|bbb1~bbb2~bbb3\rTMP|ddd1~|~ddd2|~"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"), newSimpleStr("aaa"),
					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr("bbb1"), newSimpleStr("bbb2"), newSimpleStr("bbb3"),
					),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"),
					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr("ddd1"), newSimpleStr(""),
					),
					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr(""), newSimpleStr("ddd2"),
					),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr(""), newSimpleStr(""),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestComponentFieldSimple(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|^\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr(""), newSimpleStr(""),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa^aaa1|bbb1^bbb2^^bbb3\rTMP|ddd1^|^ddd2|^"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr("aaa"), newSimpleStr("aaa1"),
					),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr("bbb1"), newSimpleStr("bbb2"), newSimpleStr(""), newSimpleStr("bbb3"),
					),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr("ddd1"), newSimpleStr(""),
					),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr(""), newSimpleStr("ddd2"),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr(""), newSimpleStr(""),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestSubComponentFieldSimple(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|&\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa&aaa1|bbb1&bbb2&&bbb3\rTMP|ddd1&|&ddd2|&"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
							newSimpleStr("aaa"), newSimpleStr("aaa1"),
						),
					),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("bbb1"), newSimpleStr("bbb2"), newSimpleStr(""), newSimpleStr("bbb3"),
						),
					),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("ddd1"), newSimpleStr(""),
						),
					),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr(""), newSimpleStr("ddd2"),
						),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestRepeatedComponentField(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|^~^"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa1^aaa2~bbb1^bbb2~|~ddd1^ddd2|eee1~^eee2|fff1^~\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr("aaa1"), newSimpleStr("aaa2"),
						),
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr("bbb1"), newSimpleStr("bbb2"),
						),
						newSimpleStr(""),
					),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr(""),
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr("ddd1"), newSimpleStr("ddd2"),
						),
					),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr("eee1"),
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr(""), newSimpleStr("eee2"),
						),
					),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr("fff1"), newSimpleStr(""),
						),
						newSimpleStr(""),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestComponentSubComponentField(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|&^&"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr(""), newSimpleStr(""),
						),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|aaa1&aaa2^bbb1&bbb2^|^ddd1&ddd2|eee1^&eee2|fff1&^\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("aaa1"), newSimpleStr("aaa2"),
						),
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("bbb1"), newSimpleStr("bbb2"),
						),
						newSimpleStr(""),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr(""),
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("ddd1"), newSimpleStr("ddd2"),
						),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newSimpleStr("eee1"),
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr(""), newSimpleStr("eee2"),
						),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, SubComponentValidator,
							newSimpleStr("fff1"), newSimpleStr(""),
						),
						newSimpleStr(""),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestRepeatedSubComponentField(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|a1&a2~\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newComplexFieldWithChildren(SubComponent, SubComponentValidator,
								newSimpleStr("a1"), newSimpleStr("a2"),
							),
						),
						newSimpleStr(""),
					),
				),
			),
		},
		{
			[]byte("MSH|^~\\&|~a1&a2\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newSimpleStr(""),
						newComplexFieldWithChildren(Component, ComponentValidator,
							newComplexFieldWithChildren(SubComponent, SubComponentValidator,
								newSimpleStr("a1"), newSimpleStr("a2"),
							),
						),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestRepeatedComponentSubComponentField(t *testing.T) {

	table := []struct {
		input  []byte
		result Field
	}{
		{
			[]byte("MSH|^~\\&|^&~|aa1^bb1&bb2~cc1\r"),
			newComplexFieldWithChildren(
				message, MessageValidator,

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("MSH"), newSimpleStr("^~\\&"),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr(""),
							newComplexFieldWithChildren(SubComponent, SubComponentValidator,
								newSimpleStr(""), newSimpleStr(""),
							),
						),
						newSimpleStr(""),
					),

					newComplexFieldWithChildren(Repeated, RepeatedValidator,
						newComplexFieldWithChildren(Component, ComponentValidator,
							newSimpleStr("aa1"),
							newComplexFieldWithChildren(SubComponent, SubComponentValidator,
								newSimpleStr("bb1"), newSimpleStr("bb2"),
							),
						),
						newSimpleStr("cc1"),
					),
				),
			),
		},
	}

	for _, test := range table {

		raw := test.input
		parser, err := NewHl7Parser(raw)

		if err != nil {
			t.Fatal(err)
		}

		msg, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		err = deepEqual(msg.ComplexField, test.result)
		if err != nil {
			t.Fatal(err)
		}
	}

}
