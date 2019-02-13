package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	//"testing"
	//"io/ioutil"
	//"os"
	"log"
	"testing"
)


func TestSimpleMessage(t *testing.T) {

	raw := []byte("MSH|^~\\&|aaa|bbbb")

	parser, err := NewHl7Parser(raw)
	if err != nil {
		log.Panic(err)
	}

	msg, err := parser.Parse()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("%#v\n",msg)
}


