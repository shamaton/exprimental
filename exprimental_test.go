package experimental

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/shamaton/zeroformatter"
	"github.com/shamaton/zeroformatter/char"
	"github.com/shamaton/zeroformatter/datetimeoffset"
)

var emptyByte = []byte("")

func TestSimple(t *testing.T) {

	var Int16 int16
	Int16Ans := int16(-16)
	if err := checkRoutine(&Int16, &Int16Ans, "Int16.pack"); err != nil {
		t.Error(err)
	}

	var Int int
	IntAns := -32
	if err := checkRoutine(&Int, &IntAns, "Int32.pack"); err != nil {
		t.Error(err)
	}

	var Int32 int32
	Int32Ans := int32(-32)
	if err := checkRoutine(&Int32, &Int32Ans, "Int32.pack"); err != nil {
		t.Error(err)
	}

	var Int64 int64
	Int64Ans := int64(-64)
	if err := checkRoutine(&Int64, &Int64Ans, "Int64.pack"); err != nil {
		t.Error(err)
	}

	var Uint16 uint16
	Uint16Ans := uint16(16)
	if err := checkRoutine(&Uint16, &Uint16Ans, "UInt16.pack"); err != nil {
		t.Error(err)
	}

	var Uint uint
	UintAns := uint(32)
	if err := checkRoutine(&Uint, &UintAns, "UInt32.pack"); err != nil {
		t.Error(err)
	}

	var Uint32 uint32
	Uint32Ans := uint32(32)
	if err := checkRoutine(&Uint32, &Uint32Ans, "UInt32.pack"); err != nil {
		t.Error(err)
	}

	var Uint64 uint64
	Uint64Ans := uint64(64)
	if err := checkRoutine(&Uint64, &Uint64Ans, "UInt64.pack"); err != nil {
		t.Error(err)
	}

	var Float32 float32
	Float32Ans := float32(1.23456)
	if err := checkRoutine(&Float32, &Float32Ans, "Single.pack"); err != nil {
		t.Error(err)
	}

	var Float64 float64
	Float64Ans := float64(2.3456789)
	if err := checkRoutine(&Float64, &Float64Ans, "Double.pack"); err != nil {
		t.Error(err)
	}

	var Bool bool
	BoolAns := false
	if err := checkRoutine(&Bool, &BoolAns, "Boolean.pack"); err != nil {
		t.Error(err)
	}

	var Byte uint8
	ByteAns := uint8(math.MaxUint8)
	if err := checkRoutine(&Byte, &ByteAns, "Byte.pack"); err != nil {
		t.Error(err)
	}

	var Sbyte int8
	SbyteAns := int8(-127)
	if err := checkRoutine(&Sbyte, &SbyteAns, "SByte.pack"); err != nil {
		t.Error(err)
	}

	String := ""
	StringAns := "This is simple pack."
	if err := checkRoutine(&String, &StringAns, "String.pack"); err != nil {
		t.Error(err)
	}

	var Char char.Char
	CharAns := char.Char('a')
	if err := checkRoutine(&Char, &CharAns, "Char.pack"); err != nil {
		t.Error(err)
	}

	Time := time.Time{}
	TimeAns := time.Unix(1480846414, 631973000)
	if err := checkRoutine(&Time, &TimeAns, "DateTime.pack"); err != nil {
		t.Error(err)
	}

	TimeOffset := datetimeoffset.DateTimeOffset{}
	TimeOffsetAns := datetimeoffset.Unix(1480846414, 681594000)
	if err := checkRoutine(&TimeOffset, &TimeOffsetAns, "DateTimeOffset.pack"); err != nil {
		t.Error(err)
	}

	Duration := time.Duration(0)
	DurationAns := time.Duration(10 * time.Millisecond)
	if err := checkRoutine(&Duration, &DurationAns, "TimeSpan.pack"); err != nil {
		t.Error(err)
	}

}

func TestArray(t *testing.T) {

	IntSlice := []int{}
	if err := checkRoutine(&IntSlice, &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MaxInt32}, "ListInt.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(IntSlice)

	IntArr := [10]int32{}
	if err := checkRoutine(&IntArr, &[10]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MaxInt32}, "ListInt.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(IntArr)

	FloatSlice := []float32{}
	if err := checkRoutine(&FloatSlice, &[]float32{1.2, 3.4, 5.6, 7.8}, "ListFloat.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(FloatSlice)

	FloatArray := [4]float32{}
	if err := checkRoutine(&FloatArray, &[4]float32{1.2, 3.4, 5.6, 7.8}, "ListFloat.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(FloatArray)

	StringSlice := []string{}
	if err := checkRoutine(&StringSlice, &[]string{"Can", "you", "see", "this", "array", "message", "?"}, "ListString.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(StringSlice)

	StringArray := [7]string{}
	if err := checkRoutine(&StringArray, &[7]string{"Can", "you", "see", "this", "array", "message", "?"}, "ListString.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(StringArray)

	EmptySlice := []uint64{}
	if err := checkRoutine(&EmptySlice, &[]uint64{}, "ListEmpty.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(EmptySlice)

	EmptyArray := [0]uint64{}
	if err := checkRoutine(&EmptyArray, &[0]uint64{}, "ListEmpty.pack"); err != nil {
		t.Error(err)
	}
	//t.Log(EmptyArray)

}

func TestMap(t *testing.T) {
	mpInt := map[int]int{}
	mpIntAns := map[int]int{1: 2, 3: 4}
	if err := checkRoutine(&mpInt, &mpIntAns, "MapInt.pack"); err != nil {
		t.Error(err)
	}

	mpStr := map[string]string{}
	mpStrAns := map[string]string{"one": "two", "three": "four", "five": "six"}
	if err := checkRoutine(&mpStr, &mpStrAns, "MapString.pack"); err != nil {
		t.Error(err)
	}
}

func TestStruct(t *testing.T) {

	type Struct struct {
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
		Int8           int8
		Char           char.Char
		TimeSpan       time.Duration
		DateTime       time.Time
		DateTimeOffset datetimeoffset.DateTimeOffset
		String         string
	}

	st := &Struct{}
	stAns := &Struct{
		Int16:          math.MinInt16,
		Int:            math.MinInt32,
		Int64:          math.MinInt64,
		Uint16:         math.MaxUint16,
		Uint:           math.MaxUint32,
		Uint64:         math.MaxUint64,
		Float:          -math.MaxFloat32,
		Double:         math.MaxFloat64,
		Bool:           true,
		Uint8:          math.MaxUint8,
		Int8:           math.MaxInt8,
		Char:           char.Char('a'),
		TimeSpan:       time.Duration(1 * time.Millisecond),
		DateTime:       time.Unix(1480846414, 795326000),
		DateTimeOffset: datetimeoffset.Unix(1480846414, 795326000),
		String:         "Hello!! Can you see this text?",
	}
	if err := checkRoutine(st, stAns, "Primitive.pack"); err != nil {
		t.Error(err)
	}

	type StructPointer struct {
		Int16          *int16
		Int            *int
		Int64          *int64
		Uint16         *uint16
		Uint           *uint
		Uint64         *uint64
		Float          *float32
		Double         *float64
		Bool           *bool
		Uint8          *byte
		Int8           *int8
		Char           *char.Char
		TimeSpan       *time.Duration
		DateTime       *time.Time
		DateTimeOffset *datetimeoffset.DateTimeOffset
		String         *string
	}

	stp := &StructPointer{}
	Int16 := int16(math.MinInt16)
	Int := int(math.MinInt32)
	Int64 := int64(math.MinInt64)
	Uint16 := uint16(math.MaxUint16)
	Uint := uint(math.MaxUint32)
	Uint64 := uint64(math.MaxUint64)
	Float := float32(-math.MaxFloat32)
	Double := float64(math.MaxFloat64)
	Bool := true
	Uint8 := uint8(math.MaxUint8)
	Int8 := int8(math.MaxInt8)
	Char := char.Char('a')
	TimeSpan := time.Duration(1 * time.Millisecond)
	DateTime := time.Unix(1480846414, 795326000)
	DateTimeOffset := datetimeoffset.Unix(1480846414, 795326000)
	String := "Hello!! Can you see this text?"

	stpAns := &StructPointer{
		Int16:          &Int16,
		Int:            &Int,
		Int64:          &Int64,
		Uint16:         &Uint16,
		Uint:           &Uint,
		Uint64:         &Uint64,
		Float:          &Float,
		Double:         &Double,
		Bool:           &Bool,
		Uint8:          &Uint8,
		Int8:           &Int8,
		Char:           &Char,
		TimeSpan:       &TimeSpan,
		DateTime:       &DateTime,
		DateTimeOffset: &DateTimeOffset,
		String:         &String,
	}
	if err := checkRoutine(stp, stpAns, "Primitive.pack"); err != nil {
		t.Error(err)
	}
}

func checkRoutine(holder interface{}, answer interface{}, fileName string) error {
	d, err := fileToBytes("zeroformatter/" + fileName)
	if err != nil {
		return err
	}
	if err := zeroformatter.Deserialize(holder, d); err != nil {
		return err
	}
	if !reflect.DeepEqual(holder, answer) {
		return fmt.Errorf("data is not correct!! please check. \n%v", holder)
	}
	return nil
}

func fileToBytes(fileName string) ([]byte, error) {

	file, err := os.Open("pack/" + fileName)
	defer file.Close()
	if err != nil {
		return emptyByte, err
	}

	fi, err := file.Stat()
	if err != nil {
		return emptyByte, err
	}

	b := make([]byte, fi.Size())
	n, err := file.Read(b)
	if err != nil {
		return emptyByte, err
	}

	// size check
	if n != len(b) {
		return emptyByte, errors.New("size wrong!!")
	}

	return b, nil
}
