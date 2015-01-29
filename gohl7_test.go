package gohl7_test

import (
	"bitbucket.org/yehezkel/gohl7"
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
		_, err := gohl7.NewParser([]byte(v))
		if err == nil {
			t.Fatalf("Expecting error with header %s\n", v)
		}
	}
}

func TestSample(t *testing.T) {
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
		{[]byte("MSH|^~\\&|\r\nEVN|A23"), 2},
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
			t.Fatalf("expecting *gohtl7.Component")
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
