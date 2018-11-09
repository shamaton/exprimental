package experimental

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/zeroformatter"
	"github.com/shamaton/zeroformatter/datetimeoffset"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type BenchChild struct {
	Int    int
	String string
}

type BenchMarkStruct struct {
	Int    int
	Uint   uint
	Float  float32
	Double float64
	Bool   bool
	String string
	Array  []int
	Map    map[string]string
	Map2   map[string]string
	Child  BenchChild
}

var s = BenchMarkStruct{
	Int:    -123,
	Uint:   456,
	Float:  1.234,
	Double: 6.789,
	Bool:   true,
	String: "this is text.",
	Array:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:    map[string]string{"this": "1", "is": "2", "map": "3"},
	Map2:   map[string]string{"this": "4", "is": "5", "map2": "6"},
	Child:  BenchChild{Int: 123456, String: "this is struct of child"},
}

var zeroData, _ = zeroformatter.Serialize(s)
var msgData, _ = msgpack.Marshal(s)

func BenchmarkMMM(b *testing.B) {

	aaa := map[int]int{}
	for i := 0; i < 10000; i++ {
		aaa[i] = i + 1
	}

	for n := 0; n < b.N; n++ {
		if _, err := zeroformatter.Serialize(aaa); err != nil {
			b.Fatal(err)
		}
	}
}

func _BenchmarkAAA(b *testing.B) {
	a := []int{}
	for n := 0; n < b.N; n++ {
		a = append(a, n)
	}
}

func _BenchmarkBBB(b *testing.B) {
	a := make([]int, 0, b.N)
	for n := 0; n < b.N; n++ {
		a = append(a, n)
	}
}
func _BenchmarkCCC(b *testing.B) {
	a := make([]int, b.N, b.N)
	for n := 0; n < b.N; n++ {
		a[n] = n
	}
}

func BenchmarkNNN(b *testing.B) {

	aaa := map[int]int{}
	for i := 0; i < 10000; i++ {
		aaa[i] = i + 1
	}

	for n := 0; n < b.N; n++ {
		if _, err := msgpack.Marshal(aaa); err != nil {
			b.Fatal(err)
		}
	}
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
	for n := 0; n < b.N; n++ {
		t := BenchMarkStruct{}
		if err := zeroformatter.Deserialize(&t, zeroData); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnpackZeroformatterDelay(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := BenchMarkStruct{}
		if _, err := zeroformatter.DelayDeserialize(&t, zeroData); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnpackMsgpack(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := BenchMarkStruct{}
		if err := msgpack.Unmarshal(msgData, &t); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDuration(b *testing.B) {
	a := time.Duration(1)

	for n := 0; n < b.N; n++ {
		if _, err := zeroformatter.Serialize(a); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDateTimeOffset(b *testing.B) {
	a := datetimeoffset.Now()

	for n := 0; n < b.N; n++ {
		if _, err := zeroformatter.Serialize(a); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDateTime(b *testing.B) {
	a := time.Now()

	for n := 0; n < b.N; n++ {
		if _, err := zeroformatter.Serialize(a); err != nil {
			b.Fatal(err)
		}
	}
}

func TestCheck(t *testing.T) {

	type hoge struct {
		A int
		B int
	}
	ab := hoge{}
	abr := reflect.ValueOf(&ab)
	abr = abr.Elem()

	rr := reflect.ValueOf(&ab.A)
	rr = rr.Elem()

	for i := 0; i < abr.NumField(); i++ {
		t := abr.Field(i)
		if reflect.DeepEqual(t.Addr(), rr.Addr()) {
			fmt.Println("correct!!")
		}

		ta := t.Addr()
		ra := rr.Addr()
		fmt.Println(t.Addr().Pointer(), ra.Pointer())
		if ta.Pointer() == ra.Pointer() {
			fmt.Println("correct2!!")
		}
	}
	fmt.Println(rr.Addr())

	sha := hoge{
		A: 1234,
		B: 5678,
	}
	shab, err := zeroformatter.Serialize(sha)
	if err != nil {
		t.Error(err)
	}

	shasha := hoge{}
	dds, err := zeroformatter.DelayDeserialize(&shasha, shab)
	if err != nil {
		t.Error(err)
	}
	/*
		err = dds.DeserializeByElement(&shasha.B, &shasha.A)
		if err != nil {
			t.Error(err)
		}
	*/
	err = dds.DeserializeByIndex(1, 0)
	if err != nil {
		t.Error(err)
	}

	// test case
	/*
		var DummyA int
		err = dds.DeserializeByElement(&DummyA)
		if err != nil {
			t.Error(err)
		}
	*/

	t.Log(shasha)

	if ans, err := dds.IsDeserialized(&shasha.B); ans && err == nil {
		t.Log("is deserialzed")
	}
	/*
		d, err := fileToBytes("zeroformatter/" + "MapInt.pack")
		if err != nil {
			t.Error(err)
		}
		t.Log(d)

		dd, err := fileToBytes("zeroformatter/" + "MapString.pack")
		if err != nil {
			t.Error(err)
		}
		t.Log(dd)
	*/

	type Struct struct {
		String string
	}
	h := Struct{String: "zeroformatter"}

	d, err := zeroformatter.Serialize(h)
	if err != nil {
		// log.Fatal(err)
		t.Error(err)
	}
	r := Struct{}
	err = zeroformatter.Deserialize(&r, d)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)

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
