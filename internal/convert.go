package internal

import (
	"reflect"
	"strconv"
)

// 将字符串转为 int, int16...等类型
func StringToNumber(numberStr string, tp reflect.Kind) interface{} {
	switch tp {
	case reflect.Int:
		return StringToInt(numberStr)
	case reflect.Int8:
		return StringToInt8(numberStr)
	case reflect.Int16:
		return StringToInt16(numberStr)
	case reflect.Int32:
		return StringToInt32(numberStr)
	case reflect.Int64:
		return StringToInt64(numberStr)
	case reflect.Uint:
		return StringToUint(numberStr)
	case reflect.Uint8:
		return StringToUint8(numberStr)
	case reflect.Uint16:
		return StringToUint16(numberStr)
	case reflect.Uint32:
		return StringToUint32(numberStr)
	case reflect.Uint64:
		return StringToUint64(numberStr)
	case reflect.Float32:
		return StringToFloat32(numberStr)
	case reflect.Float64:
		return StringToFloat64(numberStr)
	default:
		return 0
	}
}

func StringToInt(numberStr string) int {
	i64, _ := strconv.ParseInt(numberStr, 10, 64)
	return int(i64)
}

func StringToInt8(numberStr string) int8 {
	i64, _ := strconv.ParseInt(numberStr, 10, 8)
	return int8(i64)
}

func StringToInt16(numberStr string) int16 {
	i64, _ := strconv.ParseInt(numberStr, 10, 16)
	return int16(i64)
}

func StringToInt32(numberStr string) int32 {
	i64, _ := strconv.ParseInt(numberStr, 10, 32)
	return int32(i64)
}

func StringToInt64(numberStr string) int64 {
	i64, _ := strconv.ParseInt(numberStr, 10, 64)
	return i64
}

func StringToUint(numberStr string) uint {
	ui64, _ := strconv.ParseUint(numberStr, 10, 64)
	return uint(ui64)
}

func StringToUint8(numberStr string) uint8 {
	ui64, _ := strconv.ParseUint(numberStr, 10, 8)
	return uint8(ui64)
}

func StringToUint16(numberStr string) uint16 {
	ui64, _ := strconv.ParseUint(numberStr, 10, 16)
	return uint16(ui64)
}

func StringToUint32(numberStr string) uint32 {
	ui64, _ := strconv.ParseUint(numberStr, 10, 32)
	return uint32(ui64)
}

func StringToUint64(numberStr string) uint64 {
	ui64, _ := strconv.ParseUint(numberStr, 10, 64)
	return ui64
}

func StringToFloat32(numberStr string) float32 {
	f64, _ := strconv.ParseFloat(numberStr, 32)
	return float32(f64)
}

func StringToFloat64(numberStr string) float64 {
	f64, _ := strconv.ParseFloat(numberStr, 64)
	return f64
}
