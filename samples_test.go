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

	raw := []byte("MSH|^~\\&|aaa|bbbb\rTMP|123|456")

	parser, err := NewHl7Parser(raw)

	if err != nil {
		log.Panic(err)
	}

	msg, err := parser.Parse()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("%s\n", msg.ComplexField)
}

func TestSimpleComponentMessage(t *testing.T) {

	raw := []byte("MSH|^~\\&|aaa|bbb1^bbb2\rTMP|123a^123b^123c|456")

	parser, err := NewHl7Parser(raw)

	if err != nil {
		log.Panic(err)
	}

	msg, err := parser.Parse()

	if err != nil {
		log.Panic(err)
	}

	log.Printf("%s\n", msg.ComplexField)
}
