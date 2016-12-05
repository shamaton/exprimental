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

	// todo : pointer to pointer

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

	_, err := ds.deserialize(t, 0)
	return err
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

func (d *deserializer) deserialize(st reflect.Value, offset uint32) (uint32, error) {
	var err error

	fmt.Println("--------->", st.Type())

	switch st.Kind() {
	case reflect.Int8:
		b, o := d.read_s1(offset)
		v := int8(b)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Int16:
		// Int16 [short(2)]
		b, o := d.read_s2(offset)
		_v := binary.LittleEndian.Uint16(b)
		v := int16(_v)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

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
			b, o := d.read_s4(offset)
			_v := binary.LittleEndian.Uint32(b)
			v := int32(int32(_v))
			st.Set(reflect.ValueOf(v))
			// update
			offset = o
		}

	case reflect.Int:
		// Int32 [int(4)]
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		// NOTE : double cast
		v := int(int32(_v))
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Int64:
		// Int64 [long(8)]
		b, o := d.read_s8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := int64(_v)
		st.SetInt(v)
		// update
		offset = o

	case reflect.Uint8:
		// byte in cSharp
		_v, o := d.read_s1(offset)
		v := uint8(_v)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint16:
		// Uint16 / Char
		b, o := d.read_s2(offset)
		v := binary.LittleEndian.Uint16(b)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint32:
		b, o := d.read_s4(offset)
		v := binary.LittleEndian.Uint32(b)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint:
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := uint(_v)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint64:
		b, o := d.read_s8(offset)
		v := binary.LittleEndian.Uint64(b)
		st.SetUint(v)
		// update
		offset = o

	case reflect.Float32:
		// Single
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := math.Float32frombits(_v)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Float64:
		// Double
		b, o := d.read_s8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := math.Float64frombits(_v)
		st.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Bool:
		b, o := d.read_s1(offset)
		if b == 0x01 {
			st.SetBool(true)
		} else if b == 0x00 {
			st.SetBool(false)
		}
		// update
		offset = o

	case reflect.String:
		b, o := d.read_s4(offset)
		l := binary.LittleEndian.Uint32(b)
		v := string(d.data[o : o+l])
		st.SetString(v)
		// update
		offset = o + l

	case reflect.Struct:
		if isDateTime(st) {
			b, o1 := d.read_s8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.read_s4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			v := time.Unix(int64(seconds), int64(nanos))
			//fmt.Println(int64(seconds), int64(nanos))
			st.Set(reflect.ValueOf(v))
			// update
			offset = o2
		} else {

			/*
				t.Log("this is struct")
				t.Log(st.NumField())

				for i := 0; i < st.NumField(); i++ {
					v := st.Field(i)
					cconv(v.Interface(), t)
				}
			*/
		}

	case reflect.Slice:
		// element type
		e := st.Type().Elem()

		// length
		b, offset := d.read_s4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return offset, nil
		}

		o := offset
		tmpSlice := reflect.MakeSlice(st.Type(), l, l)

		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			tmpSlice.Index(i).Set(v)
		}
		st.Set(tmpSlice)

		// update
		offset = o

	case reflect.Array:
		// element type
		e := st.Type().Elem()

		// length
		b, offset := d.read_s4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return offset, nil
		}
		if l != st.Len() {
			return 0, fmt.Errorf("Array Length is different : data[%d] array[%d]", l, st.Len())
		}

		o := offset
		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}
			st.Index(i).Set(v)
		}

		// update
		offset = o

	case reflect.Map:
	//t.Log("this is map")

	default:
		//t.Log("unknown....")
	}

	return offset, err
}
