package experimental

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"
)

type deserializer struct {
	data []byte
}

func createDeserializer(data []byte) *deserializer {
	return &deserializer{
		data: data,
	}
}

const minStructDataSize = 9

func Deserialize(holder interface{}, data []byte) error {
	ds := createDeserializer(data)

	t := reflect.ValueOf(holder)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", holder)
	}

	t = t.Elem()

	// Struct
	if t.Kind() == reflect.Struct && !isDateTime(t) {
		return ds.deserializeStruct(t)
	}
	// Primitive?
	//if isPrimitive(t) {
	return ds.deserialize(t, 0)
	//}

	//return nil
}

func (d *deserializer) deserializeStruct(t reflect.Value) error {
	dataLen := len(d.data)
	if dataLen < minStructDataSize {
		return fmt.Errorf("data size is not enough: %s", dataLen)
	}
	// size
	b, _ := d.read_s4(0)
	size := binary.LittleEndian.Uint32(b)
	if size != uint32(dataLen) {
		return fmt.Errorf("data size is wrong [%s != %s]", size, dataLen)
	}

	// index
	// todo : implement

	for i := 0; i < t.NumField(); i++ {
		indexOffset := 8 + i*4
		dataOffset := binary.LittleEndian.Uint32(d.data[indexOffset : indexOffset+4])
		filed := t.Field(i)
		d.deserialize(filed, dataOffset)

		fmt.Println(filed.Interface())
	}
	return nil
}

func isDateTime(value reflect.Value) bool {
	i := value.Interface()
	switch i.(type) {
	case time.Time:
		return true
	}
	return false
}

func isPrimitive(value reflect.Value) bool {
	switch value.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool,
		reflect.String:
		// TODO : datetime etc...
		return true
	}
	return false
}

func Serialize(holder interface{}) ([]byte, error) {
	return []byte(""), nil
}

func (d *deserializer) read_s1(index uint32) (byte, uint32) {
	rb := uint32(1)
	return d.data[index], index + rb
}

func (d *deserializer) read_s2(index uint32) ([]byte, uint32) {
	rb := uint32(2)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s4(index uint32) ([]byte, uint32) {
	rb := uint32(4)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s8(index uint32) ([]byte, uint32) {
	rb := uint32(8)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) deserialize(st reflect.Value, offset uint32) error {

	fmt.Println("--------->", st.Type())

	if isDateTime(st) {
		b, offset := d.read_s8(offset)
		seconds := binary.LittleEndian.Uint64(b)
		b, _ = d.read_s4(offset)
		nanos := binary.LittleEndian.Uint32(b)
		v := time.Unix(int64(seconds), int64(nanos))
		//fmt.Println(int64(seconds), int64(nanos))
		st.Set(reflect.ValueOf(v))

		return nil
	}

	switch st.Kind() {
	case reflect.Int8:
		b, _ := d.read_s1(offset)
		v := int8(b)
		st.Set(reflect.ValueOf(v))

	case reflect.Int16:
		// Int16 [short(2)]
		b, _ := d.read_s2(offset)
		_v := binary.LittleEndian.Uint16(b)
		v := int16(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Int32:
		// TODO : if rune
		/*
			if isRune {
				// rune [ushort(2)]
				b := []byte{d.data[d.offset], d.data[d.offset+1], 0, 0}
				_v := binary.LittleEndian.Uint32(b)
				v := rune(_v)
				st.Set(reflect.ValueOf(v))
			} else*/
		{
			// Int32 [int(4)]
			b, _ := d.read_s4(offset)
			_v := binary.LittleEndian.Uint32(b)
			v := int32(int32(_v))
			st.Set(reflect.ValueOf(v))
		}

	case reflect.Int:
		// Int32 [int(4)]
		b, _ := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		// NOTE : double cast
		v := int(int32(_v))
		st.Set(reflect.ValueOf(v))

	case reflect.Int64:
		// Int64 [long(8)]
		b, _ := d.read_s8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := int64(_v)
		st.SetInt(v)

	case reflect.Uint8:
		// byte in cSharp
		_v, _ := d.read_s1(offset)
		v := uint8(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Uint16:
		// Uint16 / Char
		b, _ := d.read_s2(offset)
		v := binary.LittleEndian.Uint16(b)
		st.Set(reflect.ValueOf(v))

	case reflect.Uint32:
		b, _ := d.read_s4(offset)
		v := binary.LittleEndian.Uint32(b)
		st.Set(reflect.ValueOf(v))

	case reflect.Uint:
		b, _ := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := uint(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Uint64:
		b, _ := d.read_s8(offset)
		v := binary.LittleEndian.Uint64(b)
		st.SetUint(v)

	case reflect.Float32:
		// Single
		b, _ := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := math.Float32frombits(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Float64:
		// Double
		b, _ := d.read_s8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := math.Float64frombits(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Bool:
		b, _ := d.read_s1(offset)
		if b == 0x01 {
			st.SetBool(true)
		} else if b == 0x00 {
			st.SetBool(false)
		}

	case reflect.String:
		b, offset := d.read_s4(offset)
		l := binary.LittleEndian.Uint32(b)
		v := string(d.data[offset : offset+l])
		st.SetString(v)

	case reflect.Struct:
	/*
		t.Log("this is struct")
		t.Log(st.NumField())

		for i := 0; i < st.NumField(); i++ {
			v := st.Field(i)
			cconv(v.Interface(), t)
		}
	*/

	case reflect.Slice, reflect.Array:
	//t.Log("this is slice array")
	/*
		var v []interface{}
		for i := 0; i < rv.Len(); i++ {
			iFace := rv.Index(i).Interface()
			if iFace != nil {
				v = append(v, mapping(iFace))
			}
		}
		mv = v
		return mv
	*/

	case reflect.Map:
	//t.Log("this is map")

	default:
		//t.Log("unknown....")
	}

	return nil
}
