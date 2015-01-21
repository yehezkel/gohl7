package gohl7_test

import (
	"bitbucket.org/yehezkel/gohl7"
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
		{[]byte("MSH|^~\\&|\rEVN|A28"), 2},
		{[]byte("MSH|^~\\&|\nEVN|A28"), 2},
		{[]byte("MSH|^~\\&|\r\nEVN|A28"), 2},
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
