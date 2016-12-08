package experimental

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func _TestSimple(t *testing.T) {
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

	errPrefixMessage := "deseliarize error : "

	var int16 int16
	if err := f(&int16, "Int16.pack"); err != nil {
		t.Error(err)
	}
	if int16 != -16 {
		t.Error(errPrefixMessage, int16)
	}

	var Int int
	if err := f(&Int, "Int32.pack"); err != nil {
		t.Error(err)
	}
	if Int != -32 {
		t.Error(errPrefixMessage, Int)
	}

	var Int32 int32
	if err := f(&Int32, "Int32.pack"); err != nil {
		t.Error(err)
	}
	if Int32 != -32 {
		t.Error(errPrefixMessage, Int32)
	}

	var Int64 int64
	if err := f(&Int64, "Int64.pack"); err != nil {
		t.Error(err)
	}
	if Int64 != -64 {
		t.Error(errPrefixMessage, Int64)
	}

	var Uint16 uint16
	if err := f(&Uint16, "UInt16.pack"); err != nil {
		t.Error(err)
	}
	if Uint16 != 16 {
		t.Error(errPrefixMessage, Uint16)
	}

	var Uint uint
	if err := f(&Uint, "UInt32.pack"); err != nil {
		t.Error(err)
	}
	if Uint != 32 {
		t.Error(errPrefixMessage, Uint)
	}

	var Uint32 uint32
	if err := f(&Uint32, "UInt32.pack"); err != nil {
		t.Error(err)
	}
	if Uint32 != 32 {
		t.Error(errPrefixMessage, Uint32)
	}

	var Uint64 uint64
	if err := f(&Uint64, "UInt64.pack"); err != nil {
		t.Error(err)
	}
	if Uint64 != 64 {
		t.Error(errPrefixMessage, Uint64)
	}

	var Float32 float32
	if err := f(&Float32, "Single.pack"); err != nil {
		t.Error(err)
	}
	if Float32 != 1.23456 {
		t.Error(errPrefixMessage, Float32)
	}

	var Float64 float64
	if err := f(&Float64, "Double.pack"); err != nil {
		t.Error(err)
	}
	if Float64 != 2.3456789 {
		t.Error(errPrefixMessage, Float64)
	}

	var Bool bool
	if err := f(&Bool, "Boolean.pack"); err != nil {
		t.Error(err)
	}
	if Bool != false {
		t.Error(errPrefixMessage, Bool)
	}

	// byte
	var Byte uint8
	if err := f(&Byte, "Byte.pack"); err != nil {
		t.Error(err)
	}
	if Byte != 255 {
		t.Error(errPrefixMessage, Byte)
	}

	// sbyte
	var Sbyte int8
	if err := f(&Sbyte, "SByte.pack"); err != nil {
		t.Error(err)
	}
	if Sbyte != -127 {
		t.Error(errPrefixMessage, Sbyte)
	}

	var String string
	if err := f(&String, "String.pack"); err != nil {
		t.Error(err)
	}
	if String != "This is simple pack." {
		t.Error(errPrefixMessage, String)
	}

	Time := time.Time{}
	if err := f(&Time, "DateTime.pack"); err != nil {
		t.Error(err)
	}
	if Time != time.Unix(1480846414, 631973000) {
		t.Error(errPrefixMessage, Time)
	}

}

func TestSDS(t *testing.T) {
	f := func(in interface{}, out interface{}, isDispByte bool) error {
		d, err := Serialize(in)
		if err != nil {
			return err
		}
		if isDispByte {
			t.Log(in, " -- to byte --> ", d)
		}
		if err := Deserialize(out, d); err != nil {
			return err
		}
		return nil
	}
	_p := func(in interface{}, out interface{}) string {
		return fmt.Sprint("value different [in]:", in, " [out]:", out)
	}

	var rInt8 int8
	vInt8 := int8(-8)
	if err := f(vInt8, &rInt8, false); err != nil {
		t.Error(err)
	}
	if vInt8 != rInt8 {
		t.Error(_p(vInt8, rInt8))
	}
	t.Log(rInt8)

	var rInt16 int16
	vInt16 := int16(-16)
	if err := f(vInt16, &rInt16, false); err != nil {
		t.Error(err)
	}
	if vInt16 != rInt16 {
		t.Error(_p(vInt16, rInt16))
	}

	var rInt int
	vInt := -65535
	if err := f(vInt, &rInt, false); err != nil {
		t.Error(err)
	}
	if vInt != rInt {
		t.Error(_p(vInt, rInt))
	}

	var rInt32 int32
	vInt32 := int32(-32)
	if err := f(vInt32, &rInt32, false); err != nil {
		t.Error(err)
	}
	if vInt32 != rInt32 {
		t.Error(_p(vInt32, rInt32))
	}

	var rInt64 int64
	vInt64 := int64(-64)
	if err := f(vInt64, &rInt64, false); err != nil {
		t.Error(err)
	}
	if vInt64 != rInt64 {
		t.Error(_p(vInt64, rInt64))
	}
	t.Log(rInt64)

	var rUint8 uint8
	vUint8 := uint8(math.MaxUint8)
	if err := f(vUint8, &rUint8, false); err != nil {
		t.Error(err)
	}
	if vUint8 != rUint8 {
		t.Error(_p(vUint8, rUint8))
	}
	t.Log(rUint8)

	var rUint16 uint16
	vUint16 := uint16(math.MaxUint16)
	if err := f(vUint16, &rUint16, false); err != nil {
		t.Error(err)
	}
	if vUint16 != rUint16 {
		t.Error(_p(vUint16, rUint16))
	}
	t.Log(rUint16)

	var rUint uint
	vUint := uint(math.MaxUint32 / 2)
	if err := f(vUint, &rUint, false); err != nil {
		t.Error(err)
	}
	if vUint != rUint {
		t.Error(_p(vUint, rUint))
	}
	t.Log(rUint)

	var rUint32 uint32
	vUint32 := uint32(math.MaxUint32)
	if err := f(vUint32, &rUint32, false); err != nil {
		t.Error(err)
	}
	if vUint32 != rUint32 {
		t.Error(_p(vUint32, rUint32))
	}
	t.Log(rUint32)

	var rUint64 uint64
	vUint64 := uint64(math.MaxUint64)
	if err := f(vUint64, &rUint64, false); err != nil {
		t.Error(err)
	}
	if vUint64 != rUint64 {
		t.Error(_p(vUint64, rUint64))
	}
	t.Log(rUint64)

	var rFloat32 float32
	vFloat32 := float32(math.MaxFloat32)
	if err := f(vFloat32, &rFloat32, false); err != nil {
		t.Error(err)
	}
	if vFloat32 != rFloat32 {
		t.Error(_p(vFloat32, rFloat32))
	}
	t.Log(rFloat32)

	var rFloat64 float64
	vFloat64 := float64(math.MaxFloat64)
	if err := f(vFloat64, &rFloat64, false); err != nil {
		t.Error(err)
	}
	if vFloat64 != rFloat64 {
		t.Error(_p(vFloat64, rFloat64))
	}
	t.Log(rFloat64)

	var rBool bool
	vBool := true
	if err := f(vBool, &rBool, false); err != nil {
		t.Error(err)
	}
	if vBool != rBool {
		t.Error(_p(vBool, rBool))
	}
	t.Log(rBool)

	var rString string
	vString := "this string serialize and deserialize."
	if err := f(vString, &rString, false); err != nil {
		t.Error(err)
	}
	if vString != rString {
		t.Error(_p(vString, rString))
	}
	t.Log(rString)

	// todo : more array/slice test cases
	var rIntArr []int
	vIntArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MinInt32}
	if err := f(vIntArr, &rIntArr, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vIntArr, rIntArr) {
		t.Error(_p(vIntArr, rIntArr))
	}
	t.Log(rIntArr)

	var rStrArr []string
	vStrArr := []string{"this", "is", "string", "array", ".", "can", "you", "see", "?"}
	if err := f(vStrArr, &rStrArr, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vStrArr, rStrArr) {
		t.Error(_p(vStrArr, rStrArr))
	}
	t.Log(rStrArr)

	var rArrEmpty []string
	vArrEmpty := []string{}
	if err := f(vArrEmpty, &rArrEmpty, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vArrEmpty, rArrEmpty) {
		t.Error(_p(vArrEmpty, rArrEmpty))
	}
	t.Log(rArrEmpty)

	/*
		var _rUint8 int8
		_vUint8 := int8(-8)
		if err := f(_vUint8, &_rUint8, false); err != nil {
			t.Error(err)
		}
		if _vUint8 != _rUint8 {
			t.Error(_p(_vUint8, _rUint8))
		}
		t.Log(_rUint8)
	*/
	type childchild struct {
		String string
		Floats []float32
	}
	type child struct {
		Int   int
		Time  time.Time
		Child childchild
	}
	type st struct {
		Int16  int16
		Int    int
		Int64  int64
		Uint16 uint16
		Uint   uint
		Uint64 uint64
		Float  float32
		Double float64
		Bool   bool
		Uint8  byte
		Int8   int8
		String string
		Time   time.Time
		Child  child
	}
	vSt := &st{
		Int:    -32,
		Int8:   -8,
		Int16:  -16,
		Int64:  -64,
		Uint:   32,
		Uint8:  8,
		Uint16: 16,
		Uint64: 64,
		Float:  1.23,
		Double: 2.3456,
		Bool:   true,
		String: "hello",
		Time:   time.Now(),
		Child: child{
			Int:   1234567,
			Time:  time.Now(),
			Child: childchild{String: "this is child in child", Floats: []float32{1.2, 3.4, 5.6}},
		},
	}
	rSt := st{}
	if err := f(vSt, &rSt, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*vSt, rSt) {
		t.Error(_p(*vSt, rSt))
	}

	t.Log(rSt)
	t.Log("stst ", unsafe.Sizeof(*vSt), " : ", unsafe.Sizeof(rSt))

	// pointer test mmmm...
	hoge := new(int)
	*hoge = 123
	fuga := new(int)
	rrrr := reflect.ValueOf(&fuga)
	t.Log(rrrr.Type().Elem())
	if err := f(&hoge, &fuga, false); err != nil {
		t.Error(err)
	}
	t.Log(hoge, *hoge, fuga, *fuga)

}

func _TestArray(t *testing.T) {

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
	errSliceMessage := "slice deseliarize error : "
	errArrayMessage := "array deseliarize error : "

	IntSlice := []int{}
	if err := f(&IntSlice, "ListInt.pack"); err != nil {
		t.Error(err)
	}
	t.Log(IntSlice)
	if !reflect.DeepEqual(IntSlice, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MaxInt32}) {
		t.Error(errSliceMessage, IntSlice)
	}

	IntArr := [10]int32{}
	if err := f(&IntArr, "ListInt.pack"); err != nil {
		t.Error(err)
	}
	t.Log(IntArr)
	if !reflect.DeepEqual(IntArr, [10]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MaxInt32}) {
		t.Error(errArrayMessage, IntArr)
	}

	FloatSlice := []float32{}
	if err := f(&FloatSlice, "ListFloat.pack"); err != nil {
		t.Error(err)
	}
	t.Log(FloatSlice)
	if !reflect.DeepEqual(FloatSlice, []float32{1.2, 3.4, 5.6, 7.8}) {
		t.Error(errSliceMessage, FloatSlice)
	}

	FloatArray := [4]float32{}
	if err := f(&FloatArray, "ListFloat.pack"); err != nil {
		t.Error(err)
	}
	t.Log(FloatArray)
	if !reflect.DeepEqual(FloatArray, [4]float32{1.2, 3.4, 5.6, 7.8}) {
		t.Error(errArrayMessage, FloatArray)
	}

	StringSlice := []string{}
	if err := f(&StringSlice, "ListString.pack"); err != nil {
		t.Error(err)
	}
	t.Log(StringSlice)
	if !reflect.DeepEqual(StringSlice, []string{"Can", "you", "see", "this", "array", "message", "?"}) {
		t.Error(errSliceMessage, StringSlice)
	}

	StringArray := [7]string{}
	if err := f(&StringArray, "ListString.pack"); err != nil {
		t.Error(err)
	}
	t.Log(StringArray)
	if !reflect.DeepEqual(StringArray, [7]string{"Can", "you", "see", "this", "array", "message", "?"}) {
		t.Error(errArrayMessage, StringArray)
	}

	EmpltySlice := []uint64{}
	if err := f(&EmpltySlice, "ListEmpty.pack"); err != nil {
		t.Error(err)
	}
	t.Log(EmpltySlice)
	if !reflect.DeepEqual(EmpltySlice, []uint64{}) {
		t.Error(errSliceMessage, EmpltySlice)
	}

	EmpltyArray := [0]uint64{}
	if err := f(&EmpltyArray, "ListEmpty.pack"); err != nil {
		t.Error(err)
	}
	t.Log(EmpltyArray)
	if !reflect.DeepEqual(EmpltyArray, [0]uint64{}) {
		t.Error(errArrayMessage, EmpltyArray)
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
		DateTime       time.Time
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
