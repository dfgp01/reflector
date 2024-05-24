package internal

import (
	"errors"
	"reflect"
)

/**
*	原包来自 serversdk/generic/*	还有很多没拿过来优化
 */

var (
	ErrInvalidObjectType = errors.New("invalid object type")
	ErrNotClassType      = errors.New("invalid struct type")
	ErrInvalidSliceType  = errors.New("invalid slice type")
	ErrInvalidMapKeyType = errors.New("invalid map key type")
	ErrInvalidMapValType = errors.New("invalid map value type")
	ErrInvalidPtrType    = errors.New("invalid ptr type")
	ErrCheckType         = errors.New("check type error")
)

// 获取反射信息
func TV(v interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(v), reflect.ValueOf(v)
}

// 获取反射对象信息，参数范围：
//  1. bool, int, string 等基础类型，以及数组切片等
//  2. *struct{}, []*struct{}
//  3. map[string]string, map[interface{}]interface{}
func ReadIn(v interface{}, mustPtr bool) (reflect.Type, reflect.Value, error) {
	tp, val := TV(v)
	if mustPtr && tp.Kind() != reflect.Ptr {
		return nil, val, ErrInvalidPtrType
	}
	head := tp
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	err := checkType(tp)
	if err != nil {
		return nil, val, err
	}
	return head, val, nil
}

func MakeSliceAndAppend(ptrSlice reflect.Type, ptrSliceV reflect.Value, data ...interface{}) {

	//ptr->slice
	var (
		sliceT = ptrSlice.Elem()
		slice  = ptrSliceV.Elem()
	)

	//这种做法可行
	//slice.Set(reflect.MakeSlice(sliceT, 0, 1))

	//这种也可以
	newPtr := reflect.New(sliceT)
	slice.Set(newPtr.Elem())

	if len(data) == 0 {
		return
	}

	for _, d := range data {
		slice.Set(reflect.Append(slice, reflect.ValueOf(d)))
	}

	//var d []int64
	//t, ptr, slice := reflect.TypeOf(&d), reflect.ValueOf(&d), reflect.ValueOf(d)
	//fmt.Println("berfore:", t, ptr, slice, ptr.CanSet(), slice.CanSet(), ptr.Elem().CanSet())
	//-----------------------上面代码中
	//ptr.type=*[]int64
	//ptr.CanSet()=false，slice.CanSet()=false，但是ptr->slice.CanSet()=true
	//也许这就是指针的玄机吧
}

func MakeMapAndSet(ptrMap reflect.Type, ptrMapV reflect.Value, args ...interface{}) {

	//ptr->map
	var (
		mpT = ptrMap.Elem()
		mp  = ptrMapV.Elem()
	)

	//这种做法可行
	mp.Set(reflect.MakeMap(mpT))

	//这种也可以
	//newPtr := reflect.New(mpT)
	//mp.Set(newPtr.Elem())

	//参数格式错误
	if len(args) == 0 && len(args)%2 > 0 {
		return
	}

	//todo 要做类型检查
	for i := 0; i < len(args); i += 2 {
		mp.SetMapIndex(reflect.ValueOf(args[i]), reflect.ValueOf(args[i+1]))
	}

}

// 接受参数：&Struct{}或Struct{}
func StructIter(dest interface{}, fn func(field reflect.StructField, value reflect.Value)) error {
	t, v, err := ReadIn(dest, false)
	if err != nil {
		return err
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ErrNotClassType
	}
	for i := 0; i < t.NumField(); i++ {
		ft, fv := t.Field(i), v.Field(i)
		//ft.Anonymous的情况？
		fn(ft, fv)
	}
	return nil
}
