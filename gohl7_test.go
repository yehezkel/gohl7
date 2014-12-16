package gohl7_test

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
		_, err := gohl7.NewParser([]byte(v))
		if err == nil{
			t.Fatalf("Expecting error with header %s\n",v)
		}
	}
}

func TestSample(t *testing.T){
	data := []byte("MSH|^~\\&||bbbb\\||c^s&s~a1a1a1\rPID|435|431|433\nEVN|A28")
	parser, err := gohl7.NewParser(data)

	if err != nil{
		t.Fatalf("Unexpected error %s\n",err)
	}

	_, err = parser.Parse()

	if err != nil{
		t.Fatalf("Unexpected error %s\n",err)
	}
}
