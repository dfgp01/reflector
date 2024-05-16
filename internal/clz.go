package internal

import (
	"reflect"
)

// 输入字符串转换具体值，未完
func (o *BaseObject) StringIn(str string, val reflect.Value) {
	switch o.Type() {
	case Bool:
		if str == "true" || str == "1" {
			val.SetBool(true)
		} else {
			val.SetBool(false)
		}
	case String:
		val.SetString(str)
	case Number:
		if isInt(val.Kind()) {
			val.SetInt(StringToInt64(str))
		} else if isUint(val.Kind()) {
			val.SetUint(StringToUint64(str))
		} else {
			val.SetFloat(StringToFloat64(str))
		}
	default:
		return
	}
}

// 初始化新对象，clz、slice、map用，v是带&引用的对象，可能要在别的地方写
func (o *PtrObject) Make(v interface{}) {

	//不严谨，需细化，struct
	val := reflect.ValueOf(v)
	ins := reflect.New(o.refType.Elem())
	val.Set(ins)

	//不严谨，需细化，map
	val2 := reflect.ValueOf(v)
	ins2 := reflect.MakeMap(o.refType)
	val2.Set(ins2)

	//slice用
	ins3 := reflect.New(o.refType)
	val3 := reflect.ValueOf(v)
	val3.Set(ins3.Elem())
}

func IsEmptyVal(o IObject) bool {
	// if isNumber(reflect.Kind(o.Type())) {
	// 	return utils.ParseUint64(o.GetValue()) == 0
	// } else if o.Type() == String {
	// 	return o.GetValue() == ""
	// } else {
	// 	return o.GetValue() == nil
	// }
	return false
}
