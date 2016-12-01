package experimental

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {

}

func TestCheck(t *testing.T) {

	packData, err := fileToBytes("Primitive.pack")
	if err != nil {
		t.Error(err)
	}
	t.Log("data size : ", len(packData))

	size := getDataSizeFromZFData(packData)
	t.Log(size)
	indexNum := getDataIndexFromZFData(packData)
	t.Log(indexNum)

	/*
		off := binary.LittleEndian.Uint32(packData[8:12])
		val := binary.LittleEndian.Uint16(packData[off : off+2])
		v := int16(val)
		t.Log("val ", v)
	*/

	type stTest struct {
		Int16  int16
		Rune   rune
		Byte   byte
		Int8   int8
		String string
	}
	st := stTest{}
	cconv(st, t)

	type st2 struct {
		Int16          int16
		Int            int
		Int64          int64
		Uint16         uint16
		Uint           uint
		Uint64         uint64
		Float          float32
		Double         float64
		Bool           bool
		Uint8          byte
		Int8           int8 // Sbyte
		Char           rune
		TimeSpan       []int
		DateTime       []int
		DateTimeOffset []int
		String         string
	}
	stt := &st2{}
	rv := reflect.ValueOf(stt)
	rvst := rv.Elem()

	t.Log(rv.Type())
	t.Log(rvst.Type())

	// have to be pointer
	for i := 0; i < rvst.NumField(); i++ {
		start := 8 + i*4
		off := binary.LittleEndian.Uint32(packData[start : start+4])
		filed := rvst.Field(i)
		t.Log(filed.Type())
		ds(filed, packData, off, t)
	}

	for i := 0; i < rvst.NumField(); i++ {
		filed := rvst.Field(i)
		t.Log(filed.Interface())
	}

	if !Hole() {
		t.Log("this is log")
		t.Errorf("%s", "can you show?")
	}
}

func fileToBytes(fileName string) ([]byte, error) {

	file, err := os.Open("./pack/" + fileName)
	defer file.Close()
	if err != nil {
		return []byte(""), err
	}

	fi, err := file.Stat()
	if err != nil {
		return []byte(""), err
	}

	b := make([]byte, fi.Size())
	n, err := file.Read(b)
	if err != nil {
		return []byte(""), err
	}

	// size check
	if n != len(b) {
		return []byte(""), errors.New("size wrong!!")
	}

	return b, nil
}

func getDataSizeFromZFData(data []byte) uint32 {
	offset := 0
	size := binary.LittleEndian.Uint32(data[offset : offset+4])
	return size
}

func getDataIndexFromZFData(data []byte) uint32 {
	offset := 4
	num := binary.LittleEndian.Uint32(data[offset : offset+4])
	return num
}

func isPrimitive(value reflect.Value, t *testing.T) bool {
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

// deserialize
func ds(st reflect.Value, data []byte, offset uint32, t *testing.T) {

	t.Log("-----------------ds-------------------")
	t.Log(st.Type())

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
		_v := data[offset]
		v := int8(_v)
		st.Set(reflect.ValueOf(v))
		// if int8
	// todo : implement

	case reflect.Int16:
		// Int16 [short(2)]
		_v := binary.LittleEndian.Uint16(data[offset : offset+2])
		v := int16(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Int32:
		// TODO : if rune

		if isRune {
			// rune [ushort(2)]
			b := []byte{data[offset], data[offset+1], 0, 0}
			_v := binary.LittleEndian.Uint32(b)
			v := rune(_v)
			st.Set(reflect.ValueOf(v))
		} else {
			// Int32 [int(4)]
			_v := binary.LittleEndian.Uint32(data[offset : offset+4])
			v := int32(_v)
			st.Set(reflect.ValueOf(v))
		}

	case reflect.Int:
		// Int32 [int(4)]
		_v := binary.LittleEndian.Uint32(data[offset : offset+4])
		v := int(_v)
		st.Set(reflect.ValueOf(v))

	case reflect.Int64:
		// Int64 [long(8)]
		_v := binary.LittleEndian.Uint64(data[offset : offset+8])
		v := int64(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint8: //
		_v := data[offset]
		v := uint8(_v)
		st.Set(reflect.ValueOf(v))
	// if byte uint8

	case reflect.Uint16: // Uint16 / Char
		v := binary.LittleEndian.Uint16(data[offset : offset+2])
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint32:
		v := binary.LittleEndian.Uint32(data[offset : offset+4])
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint:
		_v := binary.LittleEndian.Uint32(data[offset : offset+4])
		v := uint(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Uint64:
		v := binary.LittleEndian.Uint64(data[offset : offset+8])
		st.SetUint(v)
	//rv.Set(v)

	case reflect.Float32: // Single
		_v := binary.LittleEndian.Uint32(data[offset : offset+4])
		v := math.Float32frombits(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Float64: // Double
		_v := binary.LittleEndian.Uint64(data[offset : offset+8])
		v := math.Float64frombits(_v)
		st.Set(reflect.ValueOf(v))
	//rv.Set(v)

	case reflect.Bool:
		b := data[offset : offset+1]
		if b[0] == 0x01 {
			st.SetBool(true)
		} else if b[0] == 0x00 {
			st.SetBool(false)
		}

	case reflect.String:
		l := binary.LittleEndian.Uint32(data[offset : offset+4])
		end := uint32(offset+4) + l
		v := string(data[offset+4 : end])
		st.SetString(v)

	case reflect.Struct:
		t.Log("this is struct")
		t.Log(st.NumField())

		for i := 0; i < st.NumField(); i++ {
			v := st.Field(i)
			cconv(v.Interface(), t)
		}

	case reflect.Slice, reflect.Array:
		t.Log("this is slice array")
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
		t.Log("this is map")

	default:
		t.Log("unknown....")
	}

}

// disp debug
func cconv(st interface{}, t *testing.T) {
	rv := reflect.ValueOf(st)

	t.Log("------------------------------------")
	t.Log(rv.Type())

	switch rv.Kind() {
	case reflect.Int8:
		// if int8
		t.Log("this is int 8")

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		// if rune int32
		t.Log("this is int")
		//v := rv.Int()
		//mv = int(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// if byte uint8
		t.Log("this is uint")
		//v := rv.Uint()
		//mv = int(v)

	case reflect.Float32, reflect.Float64:
		t.Log("this is float")
		//v := rv.Float()
		//mv = float32(v)

	case reflect.String:
		t.Log("this is string")
		//mv = rv.String()

	case reflect.Struct:
		t.Log("this is struct")
		t.Log(rv.NumField())

		for i := 0; i < rv.NumField(); i++ {
			v := rv.Field(i)
			cconv(v.Interface(), t)
		}

	case reflect.Bool:
		t.Log("this is bool")
		//mv = rv.Bool()

	case reflect.Slice, reflect.Array:
		t.Log("this is slice array")
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
		t.Log("this is map")
		/*
			mm := value.(map[interface{}]interface{})

			var itemsKey interface{} = "_items"
			var sizeKey interface{} = "_size"

			// 中身が配列で構成されている場合、配列にして返す
			iFace, isArray := mm[itemsKey]
			if isArray {
				array := iFace.([]interface{})
				var v []interface{}
				size := mm[sizeKey].(int64)
				for i := int64(0); i < size; i++ {
					log.Debug(mapping(array[i]))
					v = append(v, mapping(array[i]))
				}
				mv = v
				break
			}

			// mapを新規作成する
			var newMap = map[string]interface{}{}
			for k, v := range mm {
				s := k.(string)
				newMap[s] = mapping(v)
			}
			mv = newMap
		*/
	default:
		t.Log("unknown....")
	}

}

/*
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
*/
