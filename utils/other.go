package utils

import (
	"reflect"
	"reflector/model"
)

/**
*	暂定
 */

// 通过反射获得struct的名，参数类型范围：struct{}, *struct{}, []struct{}, []*struct{}
func GetClassName(v interface{}) (string, error) {
	t := reflect.TypeOf(v)
	return digClassName(t)
}

func digClassName(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Struct:
		return t.Name(), nil
	case reflect.Ptr, reflect.Array, reflect.Slice:
		return digClassName(t.Elem())
	default:
		return "", model.ErrNotClassType
	}
}
