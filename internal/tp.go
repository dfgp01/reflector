package internal

import (
	"reflect"
)

// 简单代表一下数据类型
type ObjType int

const (
	Invalid ObjType = iota //不支持的类型，如chan、func等
	Bool
	Number
	String
	Struct
	Slice //slice或array
	Map
	Interface //定义时为interace{}类型，运行时不确定
	Ptr       //引用带头
)

// 类型归纳
func refType(t reflect.Type) ObjType {
	k := t.Kind()
	if isBool(k) {
		return Bool
	} else if isNumber(k) {
		return Number
	} else if isString(k) {
		return String
	} else if k == reflect.Ptr {
		return Ptr
	} else if k == reflect.Struct {
		return Struct
	} else if k == reflect.Slice || k == reflect.Array {
		return Slice
	} else if k == reflect.Map {
		return Map
	} else if k == reflect.Interface {
		return Interface
	} else {
		return Invalid
	}
}

func hasElem(t reflect.Type) bool {
	return t.Kind() == reflect.Array || t.Kind() == reflect.Map || t.Kind() == reflect.Slice || t.Kind() == reflect.Pointer || t.Kind() == reflect.Ptr
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
func isNumber(k reflect.Kind) bool {
	return isInteger(k) || isFloat(k)
}

// 字符串
func isString(k reflect.Kind) bool {
	return k == reflect.String
}

// 所有基础类型，也是目前支持的map-key类型
func isBaseKind(k reflect.Kind) bool {
	return isBool(k) || isString(k) || isNumber(k)
}

// 目前map-key可支持的类型
func availableMapKeyType(o IObject) bool {
	return o.Type() == Interface || isBaseType(o)
}

// 目前map-val可支持的类型
func availableMapValType(o IObject) bool {
	return o.Type() == Interface || isBaseType(o)
}

func isBaseType(o IObject) bool {
	return o.Type() == Bool || o.Type() == Number || o.Type() == String
}

// 暂时先不用reflect.Value了
func createObjectInfo(t reflect.Type) (IObject, error) {
	tp := refType(t)
	if tp == Invalid {
		return nil, ErrInvalidObjectType
	}
	obj := &BaseObject{tp: tp, refType: t}
	switch tp {
	case Struct:
		return createStructObject(obj)
	case Slice:
		return createSliceObject(obj)
	case Map:
		return createMapObject(obj)
	case Interface:
		return &InterfaceObject{BaseObject: obj}, nil
	case Ptr:
		return createPtrObject(obj)
	default:
		return obj, nil
	}
}

func createStructObject(obj *BaseObject) (*StrutObject, error) {

	var (
		structT = obj.refType
		fields  []*StructFieldObject
	)
	if c, ok := clzCache[structT.Name()]; ok {
		return c, nil
	}

	for i := 0; i < structT.NumField(); i++ {
		sf := structT.Field(i)
		fieldObj, err := createObjectInfo(sf.Type)
		if err != nil {
			//忽略这个字段
			continue
		}
		fields = append(fields, &StructFieldObject{
			IObject: fieldObj,
			field:   sf,
		})
	}

	s := &StrutObject{
		BaseObject: obj, fields: fields,
	}

	clzCache[structT.Name()] = s
	return s, nil
}

func createSliceObject(obj *BaseObject) (*SliceObject, error) {
	//sub type
	sub, err := createObjectInfo(obj.refType.Elem())
	if err != nil {
		if err == ErrInvalidObjectType {
			return nil, ErrInvalidSliceType
		}
		return nil, err
	}
	return &SliceObject{
		HasSubObject: &HasSubObject{
			BaseObject: obj, sub: sub,
		},
	}, nil
}

func createMapObject(obj *BaseObject) (*MapObject, error) {
	//key类型限制，目前只做基础类型
	keyObj, err := createObjectInfo(obj.refType.Key())
	if err != nil {
		return nil, err
	}
	if !availableMapKeyType(keyObj) {
		return nil, ErrInvalidMapKeyType
	}
	valObj, err := createObjectInfo(obj.refType.Elem())
	if err != nil {
		return nil, err
	}
	if !availableMapValType(valObj) {
		return nil, ErrInvalidMapValType
	}

	return &MapObject{
		HasSubObject: &HasSubObject{
			BaseObject: obj, sub: valObj,
		},
		keyObj: keyObj,
	}, nil
}

func createPtrObject(obj *BaseObject) (*PtrObject, error) {

	//开始递归构建
	sub, err := createObjectInfo(obj.refType.Elem())
	if err != nil {
		return nil, err
	}
	if sub.Type() == Ptr {
		//只允许一层指针引用
		return nil, ErrInvalidPtrType
	}
	return &PtrObject{
		HasSubObject: &HasSubObject{
			BaseObject: obj, sub: sub,
		},
	}, nil
}
