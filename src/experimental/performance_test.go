package experimental

import (
	"testing"

	"github.com/shamaton/zeroformatter"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type BenchMarkStruct struct {
	Int    int
	Uint   uint
	Float  float32
	Double float64
	Bool   bool
	String string
}

var s = BenchMarkStruct{
	Int:    123,
	Uint:   456,
	Float:  1.234,
	Double: 6.789,
	Bool:   true,
	String: "this is text.",
}

func BenchmarkPackZeroformatter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := s
		if _, err := zeroformatter.Serialize(t); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackMsgpack(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := s
		if _, err := msgpack.Marshal(t); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnpackZeroformatter(b *testing.B) {
	d, _ := zeroformatter.Serialize(s)

	for n := 0; n < b.N; n++ {
		t := BenchMarkStruct{}
		if err := zeroformatter.Deserialize(&t, d); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnpackMsgpack(b *testing.B) {

	d, _ := msgpack.Marshal(s)

	for n := 0; n < b.N; n++ {
		t := BenchMarkStruct{}
		if err := msgpack.Unmarshal(d, &t); err != nil {
			b.Fatal(err)
		}
	}
}

func TestCheck(t *testing.T) {

	/*
		d, err := fileToBytes("msgpack/" + "Comparison.pack")
		if err != nil {
			t.Error(err)
		}
		t.Log(d)
	*/

	msgData, err := msgpack.Marshal(s)
	if err != nil {
		t.Error(err)
	}

	msgSt := BenchMarkStruct{}
	err = msgpack.Unmarshal(msgData, &msgSt)
	if err != nil {
		t.Error(err)
	}
	t.Log(msgSt)

	zeroData, err := zeroformatter.Serialize(s)
	if err != nil {
		t.Error(err)
	}
	zeroSt := BenchMarkStruct{}
	err = zeroformatter.Deserialize(&zeroSt, zeroData)
	if err != nil {
		t.Error(err)
	}
	t.Log(zeroSt)
}
