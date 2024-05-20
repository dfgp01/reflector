package internal

import (
	"reflect"
	"reflector/model"
)

var (

	// 缓存type映射
	typeMapper = make(map[string]reflect.Type)

	// 为防止循环依赖，暂存struct-name，建议预热，避免并发读写
	clzCache = make(map[string]bool)

	_int      int
	_int8     int8
	_int16    int16
	_int32    int32
	_int64    int64
	_uint     uint
	_uint8    uint8
	_uint16   uint16
	_uint32   uint32
	_uint64   uint64
	_float32  float32
	_float64  float64
	_ints     []int
	_int8s    []int8
	_int16s   []int16
	_int32s   []int32
	_int64s   []int64
	_uints    []uint
	_uint8s   []uint8
	_uint16s  []uint16
	_uint32s  []uint32
	_uint64s  []uint64
	_float32s []float32
	_float64s []float64
	_bool     bool
	_string   string
)

func init() {
	//缓存基础类型反射
	hot(_int)
	hot(_int8)
	hot(_int16)
	hot(_int32)
	hot(_int64)

	hot(_uint)
	hot(_uint8)
	hot(_uint16)
	hot(_uint32)
	hot(_uint64)

	hot(_float32)
	hot(_float64)

	hot(_ints)
	hot(_int8s)
	hot(_int16s)
	hot(_int32s)
	hot(_int64s)

	hot(_uints)
	hot(_uint8s)
	hot(_uint16s)
	hot(_uint32s)
	hot(_uint64s)

	hot(_float32s)
	hot(_float64s)

	hot(_bool)
	hot(_string)
}

func hot(v interface{}) {
	checkType(reflect.TypeOf(v))
}

// 不支持的类型
func InValid(k reflect.Kind) bool {
	return k == reflect.Invalid || k == reflect.Func || k == reflect.Chan
}

// 有下级元素
func HasElem(k reflect.Kind) bool {
	return k == reflect.Map ||
		k == reflect.Array || k == reflect.Slice ||
		k == reflect.Pointer || k == reflect.Ptr
}

func isBool(k reflect.Kind) bool {
	return k == reflect.Bool
}

// 有符号
func isInt(k reflect.Kind) bool {
	return k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64
}

// 无符号
func isUint(k reflect.Kind) bool {
	return k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64
}

// 浮点
func isFloat(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64 || k == reflect.Complex64 || k == reflect.Complex128
}

// 所有整数
func isInteger(k reflect.Kind) bool {
	return isInt(k) || isUint(k)
}

// 所有数值
func IsNumber(k reflect.Kind) bool {
	return isInteger(k) || isFloat(k)
}

// 所有数值数组
func IsNumberSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsNumber(t.Elem().Kind())
}

// 字符串
func IsString(k reflect.Kind) bool {
	return k == reflect.String
}

// 字符串数组
func IsStringSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsString(t.Elem().Kind())
}

// 所有基础类型，也是目前支持的map-key类型
func IsBaseType(k reflect.Kind) bool {
	return isBool(k) || IsString(k) || IsNumber(k)
}

// 所有基础类型数组
func IsBaseSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsBaseType(t.Elem().Kind())
}

func IsPtrStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func IsPtrStructSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && IsPtrStruct(t.Elem())
}

func checkType(t reflect.Type) error {
	//是否已有
	if _, ok := typeMapper[t.Name()]; ok {
		return nil
	}

	if InValid(t.Kind()) {
		return model.ErrInvalidObjectType
	}

	if HasElem(t.Kind()) {
		switch t.Kind() {
		case reflect.Ptr:
			//指针合法性
			if t.Elem().Kind() != reflect.Struct {
				return model.ErrInvalidPtrType
			}
		case reflect.Map:
			//map-key合法性
			if IsBaseType(t.Key().Kind()) {
				return model.ErrInvalidMapKeyType
			}
		}
		err := checkType(t.Elem())
		if err != nil {
			return err
		}
	}

	//struct检查
	if t.Kind() == reflect.Struct {
		//防循环依赖
		if _, ok := clzCache[t.Name()]; ok {
			return nil
		} else {
			clzCache[t.Name()] = true
		}
		for i := 0; i < t.NumField(); i++ {
			if err := checkType(t.Field(i).Type); err != nil {
				return err
			}
		}
	}

	//interface检查

	typeMapper[t.Name()] = t
	return nil
}
