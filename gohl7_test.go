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
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
							newSimpleStr("bbb1"), newSimpleStr("bbb2"), newSimpleStr(""), newSimpleStr("bbb3"),
						),
					),
				),

				newComplexFieldWithChildren(segment, SegmentValidator,
					newSimpleStr("TMP"),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
							newSimpleStr("ddd1"), newSimpleStr(""),
						),
					),
					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
							newSimpleStr(""), newSimpleStr("ddd2"),
						),
					),

					newComplexFieldWithChildren(Component, ComponentValidator,
						newComplexFieldWithChildren(SubComponent, ComponentValidator,
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

/*func TestSample(t *testing.T) {
	data := []byte("MSH|^~\\&||bbbb\\||c^s&s~a1a1a1\rPID|435|431|433\nEVN|A28")
	parser, err := gohl7.NewParser(data)

	if err != nil {
		t.Fatalf("Unexpected error %s\n", err)
	}

	_, err = parser.Parse()

	if err != nil {
		t.Fatalf("Unexpected error %s\n", err)
	}
}

func TestMultipleSegments(t *testing.T) {
	tests := []struct {
		mssg  []byte
		count int
	}{
		{[]byte("MSH|^~\\&|\rEVN|A21"), 2},
		{[]byte("MSH|^~\\&|\nEVN|A22"), 2},
		{[]byte("MSH|^~\\&|\r\nEVN|A23|"), 2},
	}

	for _, v := range tests {
		parser, err := gohl7.NewParser(v.mssg)

		if err != nil {
			t.Fatal(err)
		}

		segments, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if len(segments) != v.count {
			t.Fatalf("expecting %d amount of segments, got %d on %s\n", v.count, len(segments), v.mssg)
		}
	}

}

func TestMssgSubComponent(t *testing.T) {
	tests := []struct {
		mssg   []byte
		cindex int      //component index
		sindex int      //subcomponent index
		values []string //subcomponent values
	}{
		{[]byte("MSH|^~\\&|c1^s1&s2^c2|last"), 2, 1, []string{"s1", "s2"}},
		{[]byte("MSH|^~\\&|s1&s2^c2|last"), 2, 0, []string{"s1", "s2"}},
		{[]byte("MSH|^~\\&|&s2|last"), 2, 0, []string{"", "s2"}},
	}

	for _, v := range tests {
		parser, err := gohl7.NewParser(v.mssg)
		if err != nil {
			t.Fatal(err)
		}

		segments, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}

		if len(segments) != 1 {
			t.Fatalf("unexpected segments length\n")
		}

		field, ok := segments[0].Field(v.cindex)
		if !ok {
			t.Fatalf("Expecting field at index %d\n", v.cindex)
		}

		cmp, ok := field.(*gohl7.Component)
		if !ok {
			t.Fatalf("Expecting *gohl7.Component at position: %d\n", v.cindex)
		}

		field, ok = cmp.Field(v.sindex)
		if !ok {
			t.Fatalf("Expecting field at index: %d\n", v.sindex)
		}

		scmp, ok := field.(*gohl7.SubComponent)
		if !ok {
			t.Fatalf("Expecting subcomponent at index %d\n", v.sindex)
		}

		for index, expected := range v.values {
			field, ok := scmp.Field(index)
			if !ok {
				t.Fatalf("expecting simple field: %s at position %d", expected, index)
			}

			simple, ok := field.(*gohl7.SimpleField)

			if simple.String() != expected {
				t.Fatalf("expecting: %s got %s", expected, simple)
			}
		}

	}
}

func TestMssgComponent(t *testing.T) {
	tests := []struct {
		mssg   []byte
		index  int      //component index
		values []string //expected values on each component field
	}{
		{[]byte("MSH|^~\\&|c1^c2^c3|test"), 2, []string{"c1", "c2", "c3"}},
		{[]byte("MSH|^~\\&|c1^c2^"), 2, []string{"c1", "c2", ""}},
	}

	for _, v := range tests {
		parser, err := gohl7.NewParser(v.mssg)
		if err != nil {
			t.Fatal(err)
		}

		segments, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}

		if len(segments) < 1 {
			t.Fatalf("unexpected empty segments\n")
		}
		field, ok := segments[0].Field(v.index)
		if !ok {
			t.Fatalf("component does not exist at index %d of %s\n", v.index, v.mssg)
		}

		cmp, ok := field.(*gohl7.Component)
		if !ok {
			t.Fatalf("expecting *gohl7.Component got: %T", field)
		}

		for index, expected := range v.values {
			field, ok := cmp.Field(index)
			if !ok {
				t.Fatalf("expecting simple field: %s at position %d", expected, index)
			}

			simple, ok := field.(*gohl7.SimpleField)

			if simple.String() != expected {
				t.Fatalf("expecting: %s got %s", expected, simple)
			}
		}
	}
}

func TestMssgRepeated(t *testing.T) {
	tests := []struct {
		mssg   []byte
		index  int
		values []string
	}{
		{[]byte("MSH|^~\\&|r1~r2|end"), 2, []string{"r1", "r2"}},
		{[]byte("MSH|^~\\&|r1~r2"), 2, []string{"r1", "r2"}},
		{[]byte("MSH|^~\\&|r1~"), 2, []string{"r1", ""}},
	}

	for _, v := range tests {
		parser, err := gohl7.NewParser(v.mssg)
		if err != nil {
			t.Fatal(err)
		}
		segments, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}

		if len(segments) == 0 {
			t.Fatalf("unexpected empty segments\n")
		}
		field, ok := segments[0].Field(v.index)
		if !ok {
			t.Fatalf("repeated field does not exist at index %d of %s\n", v.index, v.mssg)
		}
		rep, ok := field.(*gohl7.Repeated)
		if !ok {
			t.Fatalf("expecting *gohl7.Repeated got: %T\n", field)
		}
		for index, expected := range v.values {
			field, ok := rep.Field(index)
			if !ok {
				t.Fatalf("expecting simple field: %s at position %d\n", expected, index)
			}

			simple, ok := field.(*gohl7.SimpleField)
			if simple.String() != expected {
				t.Fatalf("expecting: %s got %s", expected, simple)
			}
		}
	}
}*/
