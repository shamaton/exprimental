package experimental

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

const (
	byte1 uint32 = 1 << iota
	byte2
	byte4
	byte8
)

const (
	intByte1 = 1 << iota
	intByte2
	intByte4
	intByte8
)

type serializer struct {
	data []byte
}

func createSerializer() *serializer {
	return &serializer{}
}

func Serialize(holder interface{}) ([]byte, error) {
	d := createSerializer()

	t := reflect.ValueOf(holder)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	var b []byte
	if t.Kind() == reflect.Struct && !isDateTime(t) {
		startOffset := (2 + t.NumField()) * intByte4
		b = make([]byte, startOffset, t.Type().Size())
		fmt.Println(len(b))
		fmt.Println(t.Type().Size())
		d.serializeStruct(t, &b, startOffset)
		fmt.Println(len(b))
	} else {
		b = d.serialize(t)
	}
	//fmt.Println(t.Type())

	return b, nil
}

func (d *serializer) serializeStruct(rv reflect.Value, b *[]byte, offset int) {
	nf := rv.NumField()
	index := 2 * intByte4
	for i := 0; i < nf; i++ {
		// todo : サイズ受け取りたい
		ab := d.serialize(rv.Field(i))
		*b = append(*b, ab...)

		(*b)[index], (*b)[index+1], (*b)[index+2], (*b)[index+3] = byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24)
		index += intByte4
		offset += len(ab)
	}
	// size
	si := len(*b)
	(*b)[0], (*b)[1], (*b)[2], (*b)[3] = byte(si), byte(si>>8), byte(si>>16), byte(si>>24)
	(*b)[4], (*b)[5], (*b)[6], (*b)[7] = byte(nf), byte(nf>>8), byte(nf>>16), byte(nf>>24)
}

func (d *serializer) serialize(rv reflect.Value) []byte {
	var ret []byte

	switch rv.Kind() {
	case reflect.Int8:
		b := make([]byte, byte1)
		v := rv.Int()
		b[0] = byte(v)
		ret = b

	case reflect.Int16:
		b := make([]byte, byte2)
		v := rv.Int()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		ret = b

	case reflect.Int32, reflect.Int:
		b := make([]byte, byte4)
		v := rv.Int()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Int64:
		b := make([]byte, byte8)
		v := rv.Int()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		ret = b

	case reflect.Uint8:
		b := make([]byte, byte1)
		v := rv.Uint()
		b[0] = byte(v)
		ret = b

	case reflect.Uint16:
		b := make([]byte, byte2)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		ret = b

	case reflect.Uint32, reflect.Uint:
		b := make([]byte, byte4)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Uint64:
		b := make([]byte, byte8)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		ret = b

	case reflect.Float32:
		b := make([]byte, byte4)

		v := math.Float32bits(float32(rv.Float()))
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Float64:
		b := make([]byte, byte8)

		v := math.Float64bits(rv.Float())
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		ret = b

	case reflect.Bool:
		b := make([]byte, byte1)

		if rv.Bool() {
			b[0] = 0x01
		} else {
			b[0] = 0x00
		}
		ret = b

	case reflect.String:
		str := rv.String()
		l := uint32(len(str))
		b := make([]byte, 0, l+byte4)
		b = append(b, byte(l), byte(l>>8), byte(l>>16), byte(l>>24))

		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&str))
		b = append(b, strBytes...)
		ret = b

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			// first : know element size
			fb := d.serialize(rv.Index(0))

			// second : make byte array
			size := uint32(l*len(fb)) + byte4
			b := make([]byte, 0, size)

			// third : append data
			b = append(b, byte(l), byte(l>>8), byte(l>>16), byte(l>>24))
			b = append(b, fb...)

			for i := 1; i < l; i++ {
				ab := d.serialize(rv.Index(i))
				b = append(b, ab...)
			}
			ret = b
		} else {
			// only make length info
			b := make([]byte, byte4)
			ret = b
		}

	case reflect.Ptr:
	case reflect.Uintptr:
		// todo : error
	}

	return ret
}

/*
	s1 := time.Now()
	e1 := time.Now()
	fmt.Println("1:", e1.Sub(s1).Nanoseconds())
*/
