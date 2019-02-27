package gohl7

import (
	//"bitbucket.org/yehezkel/gohl7"
	//"testing"
	//"io/ioutil"
	//"os"
	"log"
	"testing"
)



func TestSimpleRepeated(t *testing.T) {

	raw := []byte("MSH|^~\\&|aaa|bbb1~bbb2~bbb3\r")

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

	raw := []byte("MSH|^~\\&|aaa|bbb1^bbb2")

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

func TestRepatedComponentMessage(t *testing.T) {

	raw := []byte("MSH|^~\\&|aaa|rrr1~bbb1^bbb2^bbb3~ccc1\r")

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
func TestRepatedComponentSubcComponentMessage(t *testing.T) {

	raw := []byte("MSH|^~\\&|aaa|rrr1~bbb1^ssss&")

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