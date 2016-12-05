package experimental

import (
	"fmt"
	"reflect"
)

type serializer struct {
	data []byte
}

func createSerializer() *serializer {
	return &serializer{}
}

func Serialize(holder interface{}) ([]byte, error) {
	d := createSerializer()

	// todo : pointer to pointer

	t := reflect.ValueOf(holder)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fmt.Println(t.Type())

	d.serialize(t)

	return []byte(""), nil
}

func (d *serializer) serialize(rv reflect.Value) []byte {
	switch rv.Kind() {
	case reflect.Int:
		b := make([]byte, 4)

		v := rv.Interface().(int)
		b[0] = byte(v >> 24)
		b[1] = byte(v >> 16)
		b[2] = byte(v >> 8)
		b[3] = byte(v)
		fmt.Println(b)
		return b
	}
	return []byte("")
}
