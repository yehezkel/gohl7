package gohl72_test

import(
	"gohl72"
	"testing"
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