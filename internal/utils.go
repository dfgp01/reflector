package internal

import (
	"reflect"
	"reflector/model"
	"strconv"
)

/**
*	原包来自 serversdk/generic/*	还有很多没拿过来优化

	重新整理一下思绪，我们的预期实现结果为：

			"string" <-> map[string]proto 需要自定义string的分隔符等，否则可能有逗号冲突
			map[string]string <-> map[string]proto 和redis可以提供原始map，问题不大
			其实可能分两步："string" <-> map[string]string <-> map[string]proto 需要指定ProtoSerializer，否则默认用json
			暂定的分隔符格式为："key1<->value1<->key2<->value2"
			转换链："string" -> map[string]string -> map[int]interface

		提供给ORM的处理接口：
			pager参数：不需要反射，但需要封装
				Pager{pageNo, size, total, totalPage, param<interface>, resp<interface>}
				param应为传入参，resp应为指针参，如 var a []*User -> &a
			这样我们要对param进行解析，对resp进行析构
			一些query包装：
				query(param<interface>, resp<interface>)，其中param可拆为：
					Cond{sort, group..., PagerParam{pageNo, size, total, totalPage}, param<interface>}}

		外部接口：In和Out，设计handler机制

		hub(in, out, args...)	args不确定要填什么
		hub("1,2,3", &[]uint16)		需要通过in和out识别它们的类型：string, ptr-slice-number
*/

/**
*	额外的快捷接口，逻辑和这包里的关系不大
 */

// 获取反射信息
func TV(v interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(v), reflect.ValueOf(v)
}

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

// 获取反射对象信息，参数范围：
//  1. bool, int, string 等基础类型，以及数组切片等
//  2. *struct{}, []*struct{}
//  3. map[string]string, map[interface{}]interface{}
func ReadIn(v interface{}) (*Object, error) {
	tp, val := TV(v)
	head := tp
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	err := checkType(tp)
	if err != nil {
		return nil, err
	}
	return &Object{
		RefTp:   head,
		RefVal:  val,
		HeadPtr: head.Kind() == reflect.Ptr,
	}, nil
}

type Object struct {
	RefTp   reflect.Type
	RefVal  reflect.Value
	HeadPtr bool
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
