package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	"testing"
	"io/ioutil"
	"os"
)

var(
	sampleFile = "test.hl7"
)


func BenchmarkLongMessage(b *testing.B){

	file, err := os.Open(sampleFile) // For read access.
	if err != nil {
		b.Fatal(err)
	}

	raw,err := ioutil.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}

    err = file.Close()

	//reset timmer
	b.ResetTimer()
    if err != nil {
        b.Fatal(err)
    }

	for i := 0; i < b.N; i++{

		parser, err := NewHl7Parser(raw)
		if err != nil{
			b.Fatal(err)
		}

		//this is the function where all the logic happends
	    _, err = parser.Parse()
		if err != nil{
			b.Fatal(err)
		}
	}
}
