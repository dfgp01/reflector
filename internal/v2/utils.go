package v2

import "reflect"

// 获取反射对象信息，参数范围：
//  1. bool, int, string 等基础类型，以及数组切片等
//  2. *struct{}, []*struct{}
//  3. map[string]string, map[interface{}]interface{}
func ReadIn(v interface{}) (*ObjectWrapper, error) {
	obj, err := CreateObjectInfo(v)
	if err != nil {
		return nil, err
	}
	return &ObjectWrapper{
		Root:    obj,
		HeadPtr: obj.Type() == Ptr,
		Val:     reflect.ValueOf(v),
	}, nil
}

type ObjectWrapper struct {
	Root    IObject
	HeadPtr bool
	Val     reflect.Value
	Name    string
}

func (o *ObjectWrapper) Invalid() bool {
	return o.Root == nil
}

func (o *ObjectWrapper) CheckType(typeChains ...ObjType) error {
	head := o.Root
	for _, tp := range typeChains {
		if head == nil {
			return ErrCheckType
		}
		if tp != head.Type() {
			return ErrCheckType
		}
		head = head.Sub()
	}
	return nil
}

func MakeSliceAndAppend(ptrSlice IObject, ptrSliceV reflect.Value, data ...interface{}) {

	//ptr->slice->number
	var (
		sliceT = ptrSlice.RefType().Elem()
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

func MakeMapAndSet(ptrMap IObject, ptrMapV reflect.Value, args ...interface{}) {

	//ptr->slice->number
	var (
		mpT = ptrMap.RefType().Elem()
		mp  = ptrMapV.Elem()
	)

	//这种做法可行
	mp.Set(reflect.MakeMap(mpT))

	//这种也可以
	//newPtr := reflect.New(mpT)
	//mp.Set(newPtr.Elem())

	if len(args) == 0 && len(args)%2 > 0 {
		return
	}

	//todo 要做类型检查
	for i := 0; i < len(args); i += 2 {
		mp.SetMapIndex(reflect.ValueOf(args[i]), reflect.ValueOf(args[i+1]))
	}

}
