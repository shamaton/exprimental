package experimental

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type deserializer struct {
	data   []byte
	offset uint32
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
	if t.Kind() == reflect.Struct {
		return ds.deserializeStruct(t)
	}
	// primitive?

	return nil
}

func (d *deserializer) deserializeStruct(t reflect.Value) error {
	dataLen := len(d.data)
	if dataLen < minStructDataSize {
		return fmt.Errorf("data size is not enough: %s", dataLen)
	}
	// size
	size := binary.LittleEndian.Uint32(d.read_s4())
	if size != uint32(dataLen) {
		return fmt.Errorf("data size is wrong [%s != %s]", size, dataLen)
	}

	// index
	// todo : implement

	for i := 0; i < t.NumField(); i++ {
		indexOffset := 8 + i*4
		d.offset = binary.LittleEndian.Uint32(d.data[indexOffset : indexOffset+4])
		filed := t.Field(i)
		d.deserialize(filed)

		fmt.Println(filed.Interface())
	}
	return nil
}

func Serialize(holder interface{}) ([]byte, error) {
	return []byte(""), nil
}

func (d *deserializer) read_s1() byte {
	defer d.addOffset(1)
	return d.data[d.offset]
}

func (d *deserializer) read_s2() []byte {
	rb := uint32(2)
	defer d.addOffset(rb)
	return d.data[d.offset : d.offset+rb]
}

func (d *deserializer) read_s4() []byte {
	rb := uint32(4)
	defer d.addOffset(rb)
	return d.data[d.offset : d.offset+rb]
}

func (d *deserializer) read_s8() []byte {
	rb := uint32(8)
	defer d.addOffset(rb)
	return d.data[d.offset : d.offset+rb]
}

func (d *deserializer) addOffset(add uint32) {
	d.offset += add
}

func (d *deserializer) deserialize(st reflect.Value) {

	fmt.Println("--------->", st.Type())

	/*
		switch i.(type) {
		case int16:
			t.Log("aaaa")
			_v := binary.LittleEndian.Uint16(data[offset : offset+2])
			v := int16(_v)
			i = v

		case int:
			t.Log("bbbb")
			_v := binary.LittleEndian.Uint32(data[offset : offset+4])
			v := int(_v)
			i = v

		}
		return
	*/
	isRune := false
	i := st.Interface()
	switch i.(type) {
	case rune:
		isRune = true
	}

	switch st.Kind() {
	case reflect.Int8:
		_v := d.read_s1()
		v := int8(_v)
		st.Set(reflect.ValueOf(v))
	// if int8
	// todo : implement

	case reflect.Int16:
		// Int16 [short(2)]
		_v := binary.LittleEndian.Uint16(d.read_s2())
		v := int16(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Int32:
		// TODO : if rune

		if isRune {
			// rune [ushort(2)]
			b := []byte{d.data[d.offset], d.data[d.offset+1], 0, 0}
			_v := binary.LittleEndian.Uint32(b)
			v := rune(_v)
			st.Set(reflect.ValueOf(v))
		} else {
			// Int32 [int(4)]
			_v := binary.LittleEndian.Uint32(d.read_s4())
			v := int32(_v)
			st.Set(reflect.ValueOf(v))
		}

	case reflect.Int:
		// Int32 [int(4)]
		_v := binary.LittleEndian.Uint32(d.read_s4())
		// NOTE : double cast
		v := int(int32(_v))
		st.Set(reflect.ValueOf(v))

	case reflect.Int64:
		// Int64 [long(8)]
		_v := binary.LittleEndian.Uint64(d.read_s8())
		v := int64(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint8: //
		_v := d.read_s1()
		v := uint8(_v)
		st.Set(reflect.ValueOf(v))
	// if byte uint8

	case reflect.Uint16: // Uint16 / Char
		v := binary.LittleEndian.Uint16(d.read_s2())
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint32:
		v := binary.LittleEndian.Uint32(d.read_s4())
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint:
		_v := binary.LittleEndian.Uint32(d.read_s4())
		v := uint(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint64:
		v := binary.LittleEndian.Uint64(d.read_s8())
		st.SetUint(v)
	//rv.Set(v)

	case reflect.Float32: // Single
		_v := binary.LittleEndian.Uint32(d.read_s4())
		v := math.Float32frombits(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Float64: // Double
		_v := binary.LittleEndian.Uint64(d.read_s8())
		v := math.Float64frombits(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Bool:
		b := d.read_s1()
		if b == 0x01 {
			st.SetBool(true)
		} else if b == 0x00 {
			st.SetBool(false)
		}

	case reflect.String:
		l := binary.LittleEndian.Uint32(d.read_s4())
		end := uint32(d.offset) + l
		v := string(d.data[d.offset:end])
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
}
