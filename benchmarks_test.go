package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	"testing"
	//"io/ioutil"
	//"os"
)

/*var(
	file_path = "test.hl7"
)


func BenchmarkLongMessage(b *testing.B){

	file, err := os.Open(file_path) // For read access.
	if err != nil {
		b.Fatal(err)
	}

	data,err := ioutil.ReadAll(file)

	//b.Fatalf("length %s\n",data[:10])

	if err != nil {
		b.Fatal(err)
	}
	buffer := make([]byte, len(data))

	//reset timmer
	b.ResetTimer()

	defer file.Close()

	for i := 0; i < b.N; i++{

		//skip copy
		b.StopTimer()
		_ = copy(buffer, data)
		b.StartTimer()

		parser, err := gohl7.NewParser(buffer)
		if err != nil{
			b.Fatal(err)
		}


		//this is the function where all the logic happends
		_, err = parser.Parse()
		if err != nil{
			b.Fatal(err)
		}
	}
}*/
