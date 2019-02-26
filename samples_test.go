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

    build := newComplexFieldWithChildren(
        message, MessageValidator,

        newComplexFieldWithChildren(segment,SegmentValidator,
            newSimpleStr("MSH"),newSimpleStr("^~\\&"), newSimpleStr("aaa"), newSimpleStr("bbbb"),
        ),

        newComplexFieldWithChildren(segment,SegmentValidator,
            newSimpleStr("TMP"),newSimpleStr("123"), newSimpleStr("456"),
        ),
    )

	parser, err := NewHl7Parser(raw)

	if err != nil {
		log.Panic(err)
	}

	msg, err := parser.Parse()

	if err != nil {
		log.Panic(err)
	}

    err = deepEqual(msg.ComplexField,build)
    if err != nil {
        log.Panic(err)
    }
}

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
