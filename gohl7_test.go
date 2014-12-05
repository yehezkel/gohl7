package gohl72_test

import(
	"testing"
	"bitbucket.org/yehezkel/gohl7"
)

func TestBadHeader(t *testing.T){
	tests := []string{
		"M||",
		"||",
		"WRONG||",
		"",
	}

	for _, v := range tests{
		_, err := gohl72.NewParser([]byte(v))
		if err == nil{
			t.Fatalf("Expecting error with header %s\n",v)
		}
	}
}
