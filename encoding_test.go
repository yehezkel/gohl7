package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	//"fmt"
	"testing"
)

func TestBadEncodingDef(t *testing.T) {

	table := []struct {
		input string
		err   error
	}{
		{
			"",
			ErrBadEncoding,
		},
		{
			"abcd",
			ErrBadEncoding,
		},
		{
			"||~\\&",
			ErrRepeatedEncoding,
		},
	}

	for _, test := range table {

		_, err := ParseEncoding([]byte(test.input))

		if err != test.err {
			t.Errorf("Expecting error %s got %s", test.err, err)
		}
	}

}

func TestCleanEncoding(t *testing.T) {

	table := []struct {
		encoding string
		input    string
		output   string
	}{
		{
			`|^~\&`,
			`aaa\^`,
			`aaa^`,
		},
		{
			`|^~\&`,
			`aa\\`,
			`aa\`,
		},
	}

	for _, test := range table {

		enc, err := ParseEncoding([]byte(test.encoding))
		if err != nil {
			t.Fatal(err)
		}

		out := enc.Clean([]byte(test.input))

		if string(out) != test.output {
			t.Errorf("Expecting clean value %s got %s", test.output, out)
		}
	}

}
