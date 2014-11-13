package gohl7_test

import(
	"github.com/yehezkel/gohl7"
	"testing"
	"os"
)

var(
	file_path = "test.hl7"
)


func BenchmarkLongMessage(b *testing.B){

	file, err := os.Open(file_path) // For read access.	
	if err != nil {
		b.Fatal(err)
	}

	//reset timmer
	b.ResetTimer()

	defer file.Close()

	for i := 0; i < b.N; i++{

		//skip parser initialization from the benchmark
		//as the buffio.Scanner is creating its buffers
		b.StopTimer()
		parser, err := gohl7.NewParser(file)
		if err != nil{
			b.Fatal(err)
		}
		b.StartTimer()
		
		//this is the function where all the logic happends
		_, err = parser.Parse()
		if err != nil{
			b.Fatal(err)
		}

		//stop timer to seek the file
		b.StopTimer()
		_, err = file.Seek(0,0)
		if err != nil{
			b.Fatal(err)
		}
		b.StartTimer()
	}
}

