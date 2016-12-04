package experimental

import (
	"encoding/binary"
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	f := func(i interface{}, fileName string) error {
		d, err := fileToBytes(fileName)
		if err != nil {
			return err
		}
		if err := Deserialize(i, d); err != nil {
			return err
		}
		return nil
	}

	var int int
	if err := f(&int, "Int32.pack"); err != nil {
		t.Error(err)
	}
	if int != -32 {
		t.Error("deseliarize error : ", int)
	}

}

func _TestCheck(t *testing.T) {

	packData, err := fileToBytes("Primitive.pack")
	if err != nil {
		t.Error(err)
	}
	t.Log("data size : ", len(packData))

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

	st3 := &st2{}
	Deserialize(st3, packData)

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
