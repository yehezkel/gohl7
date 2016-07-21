package gohl7_test

import (
	"bytes"
	"github.com/yehezkel/gohl7"
	"strings"
	"testing"
)

// This test is just an example of the package
//usage
func TestFirst(t *testing.T) {
	r := strings.NewReader("MSH|^~\\&||bbbb|c^s&s~a1a1a1\rPID|435|431|433\nEVN|A28")
	parser, err := gohl7.NewParser(r)

	if err != nil {
		t.Fatal(err)
	}

	sgments, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(sgments); i++ {
		t.Logf("%s\n", sgments[i])
	}
}

func TestSimpleField(t *testing.T) {
	tests := []struct {
		value  string
		index  int
		result string
		ok     bool
	}{
		{"test", 0, "test", true},
		{"test", 1, "", false},
		{"", 0, "", true},
		{"", 1, "", false},
	}

	for _, v := range tests {
		var hl7Simple gohl7.Hl7DataType = gohl7.SimpleField(v.value)

		result, ok := hl7Simple.Field(v.index)

		sp, _ := result.(gohl7.SimpleField)

		if string(sp) != v.result ||
			ok != v.ok {

			t.Fatalf("expecting %d, %t at index %d\n", v.result, v.ok, v.index)
		}
	}
}

func TestSubComponent(t *testing.T) {

	var r gohl7.Repeated
	var c gohl7.Component
	var s gohl7.Segment
	var sc gohl7.SubComponent

	tests := []struct {
		value gohl7.Hl7DataType
		err   bool
	}{
		{gohl7.SimpleField("test"), false},
		{gohl7.SimpleField("test2"), false},
		{r, true},
		{c, true},
		{s, true},
	}

	var err error
	for _, v := range tests {
		err = sc.AppendValue(v.value)
		if err != nil && !v.err ||
			err == nil && v.err {
			t.Fatal("Unexpecting error appending to Subcomponent")
		}
	}

	count := 0

	for i, v := range tests {
		if v.err {
			count++
			continue
		}

		result, ok := sc.Field(i)
		if !ok {
			t.Fatal("Fail Field method on SubComponent")
		}

		simpleField1 := result.(gohl7.SimpleField)
		simpleField2 := v.value.(gohl7.SimpleField)
		if string(simpleField1) != string(simpleField2) {
			t.Fatal("Expecting %s go %s\n", simpleField2, simpleField1)
		}
	}

	_, ok := sc.Field(len(tests) - count)
	if ok {
		t.Fatal("Expecting no value on position %d\n", len(tests)-count)
	}
}

func TestComponent(t *testing.T) {
	var c gohl7.Component
	var c1 gohl7.Component
	var sub gohl7.SubComponent
	var r gohl7.Repeated
	var s gohl7.Segment

	err := c.AppendValue(r)
	if err == nil {
		t.Fatal("Component can not contain a Repeated Field")
	}

	err = c.AppendValue(s)
	if err == nil {
		t.Fatal("Component can not contain a Segment")
	}

	err = c.AppendValue(c1)
	if err == nil {
		t.Fatal("Component can not contain a Component")
	}

	err = c.AppendValue(sub)
	if err != nil {
		t.Fatal("Component may contain a SubComponent")
	}

	simpleField := gohl7.SimpleField("test")
	err = c.AppendValue(simpleField)
	if err != nil {
		t.Fatal("Component may contain a SimpleField")
	}

	_, ok := c.Field(2)
	if ok {
		t.Fatal("Component Field out of bound should return false")
	}

	result, ok := c.Field(1)
	if !ok {
		t.Fatal("Bad indexing on Component")
	}

	simple, ok := result.(gohl7.SimpleField)

	if !ok {
		t.Fatal("Bad return value on Field method on Component")
	}

	if !bytes.Equal(simple, simpleField) {
		t.Fatal("Bad return field on Component")
	}
	// if simple != simpleField{
	// 	t.Fatal("Bad return field on Component")
	// }
}

func TestRepeated(t *testing.T) {
	var c gohl7.Component
	var sub gohl7.SubComponent
	var r gohl7.Repeated
	var s gohl7.Segment

	err := r.AppendValue(r)
	if err == nil {
		t.Fatal("Repeated can not contain a Repeated Field")
	}

	err = r.AppendValue(s)
	if err == nil {
		t.Fatal("Repeated can not contain a Segment")
	}

	err = r.AppendValue(c)
	if err != nil {
		t.Fatal("Repeated may contain a Component Field")
	}

	err = r.AppendValue(sub)
	if err != nil {
		t.Fatal("Repeated may contain a SubComponent Field")
	}

	simpleField := gohl7.SimpleField("test")
	err = r.AppendValue(simpleField)
	if err != nil {
		t.Fatal("Repeated may contain a SimpleField Field")
	}

	_, ok := r.Field(3)
	if ok {
		t.Fatal("Repeated Field out of bound should return false")
	}

	result, ok := r.Field(2)
	if !ok {
		t.Fatal("Bad indexing on Repeated")
	}

	simple, ok := result.(gohl7.SimpleField)

	if !ok {
		t.Fatal("Bad return value on Field method on Repeated")
	}

	if !bytes.Equal(simple, simpleField) {
		t.Fatal("Bad return field on Repeated")
	}
}

func TestSegment(t *testing.T) {
	var c gohl7.Component
	var sub gohl7.SubComponent
	var r gohl7.Repeated
	var s gohl7.Segment

	simpleField := gohl7.SimpleField("MSH")
	err := s.AppendValue(simpleField)
	if err != nil {
		t.Fatal("Segment may contain a SimpleField Field")
	}

	err = s.AppendValue(r)
	if err != nil {
		t.Fatal("Segment may contain a Repeated Field")
	}

	err = s.AppendValue(s)
	if err == nil {
		t.Fatal("Segment can not contain a Segment")
	}

	err = s.AppendValue(c)
	if err != nil {
		t.Fatal("Segment may contain a Component Field")
	}

	err = s.AppendValue(sub)
	if err != nil {
		t.Fatal("Segment may contain a SubComponent Field")
	}

	_, ok := s.Field(4)
	if ok {
		t.Fatal("Segment Field out of bound should return false")
	}

	result, ok := s.Field(0)
	if !ok {
		t.Fatal("Bad indexing on Segment")
	}

	simple, ok := result.(gohl7.SimpleField)

	if !ok {
		t.Fatal("Bad return value on Field method on Segment")
	}

	if !bytes.Equal(simple, simpleField) {
		t.Fatal("Bad return field on Repeated")
	}
}

func TestBadMssgHeader(t *testing.T) {
	tests := []string{
		"M||",
		"||",
		"WRONG||",
		"",
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		_, err := gohl7.NewParser(r)
		if err == nil {
			t.Fatalf("Expecting error with header %s\n", v)
		}
	}
}

func TestMssHeader(t *testing.T) {
	r := strings.NewReader("MSH|^~\\&|")
	_, err := gohl7.NewParser(r)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBadEncoding(t *testing.T) {
	tests := []string{
		"MSH|^~\\&sdfs",
		"MSH|^^\\&|",
		"MSH|\r~\\&|",
		"MSH|\n~\\&|",
		"MSH|\n",
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		parser, err1 := gohl7.NewParser(r)
		var err2 error

		if err1 == nil {
			_, err2 = parser.Parse()
		}

		if err1 == nil && err2 == nil {
			t.Fatalf("Expecting error with encoding %s\n", v)
		}
	}
}

func TestMultipleSegments(t *testing.T) {
	tests := []string{
		"MSH|^~\\&|\rEVN|A28",
		"MSH|^~\\&|\nEVN|A28",
		"MSH|^~\\&|\r\nEVN|A28",
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		parser, err := gohl7.NewParser(r)

		if err != nil {
			t.Fatal(err)
		}

		segments, err := parser.Parse()

		if len(segments) != 2 {
			t.Fatalf("Expecting 2 segments with %s\n", v)
		}
	}
}

func TestBadSegmentHeader(t *testing.T) {
	tests := []string{
		"MSH|^~\\&|\rEVN1|A28",
		"MSH|^~\\&|\rE|A28",
		"MSH|^~\\&|\r|A28",
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		parser, err := gohl7.NewParser(r)

		if err != nil {
			t.Fatal(err)
		}

		_, err = parser.Parse()

		if err == nil {
			t.Fatalf("Expecting error on segment header in %s\n", v)
		}
	}
}

func TestBadSeparatorAfterHeader(t *testing.T) {
	tests := []string{
		"MSH|^~\\&|\rEVN^A28|",
		"MSH|^~\\&|\rEVN\r",
		"MSH|^~\\&|\rEVN\n",
		"MSH|^~\\&|\rEVN\\a",
		"MSH|^~\\&|\rEVN~a",
		"MSH|^~\\&|\rEVN&a",
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		parser, err := gohl7.NewParser(r)

		if err != nil {
			t.Fatal(err)
		}

		_, err = parser.Parse()

		if err == nil {
			t.Fatalf("Expecting error on segment header in %s\n", v)
		}
	}
}

func TestBadEscaping(t *testing.T) {
	tests := []string{
		`MSH|^~\&|Test\b|`,
		`MSH|^~\|Test\f`,
	}

	for _, v := range tests {
		r := strings.NewReader(v)
		parser, err := gohl7.NewParser(r)

		if err != nil {
			t.Fatal(err)
		}

		_, err = parser.Parse()

		if err == nil {
			t.Fatalf("Expecting error on segment header in %s\n", v)
		}
	}
}

func TestEscaping(t *testing.T) {

	tests := []struct {
		mssg  string
		count int
		index int
		value string
	}{
		{`MSH|^~\&|Te\\st`, 3, 2, `Te\st`},
		{`MSH|^~\&|Te\^st`, 3, 2, `Te^st`},
		{`MSH|^~\&|Te\~st`, 3, 2, `Te~st`},
		{`MSH|^~\&|Te\|st`, 3, 2, `Te|st`},
		{`MSH|^~\&|Te\&st`, 3, 2, `Te&st`},
	}

	for _, v := range tests {
		r := strings.NewReader(v.mssg)
		parser, err := gohl7.NewParser(r)

		if err != nil {
			t.Fatal(err)
		}

		segments, err := parser.Parse()

		if err != nil {
			t.Fatal(err)
		}

		if len(segments) != 1 {
			t.Fatalf("expecting only one segment on %s\n", v.mssg)
		}

		s := segments[0]

		if len(s) != v.count {
			t.Fatalf("expecting %d fields on segment: %s\n", v.count, v.mssg)
		}

		field, ok := s.Field(v.index)

		if !ok {
			t.Fatalf("expecting field at index %d in %s\n", v.index, v.mssg)
		}

		str, ok := (field).(gohl7.SimpleField)

		if !ok {
			t.Fatalf("expecting SimpleField on field %d in %s\n", v.index, v.mssg)
		}

		if string(str) != v.value {
			t.Fatalf("expecting %s after escaping %s\n", v.value, s[v.index])
		}
	}
}
