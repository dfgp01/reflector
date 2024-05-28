package v2

import (
	"reflect"
	"reflector/internal/convert"
)

var (

	// 缓存type映射
	typeMapper = make(map[string]IObject)

	// 为防止循环依赖，暂存struct-name，建议预热，避免并发读写，
	clzCache = make(map[string]*StrutObject)

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
	CreateObjectInfo(_int)
	CreateObjectInfo(_int8)
	CreateObjectInfo(_int16)
	CreateObjectInfo(_int32)
	CreateObjectInfo(_int64)

	CreateObjectInfo(_uint)
	CreateObjectInfo(_uint8)
	CreateObjectInfo(_uint16)
	CreateObjectInfo(_uint32)
	CreateObjectInfo(_uint64)

	CreateObjectInfo(_float32)
	CreateObjectInfo(_float64)

	CreateObjectInfo(_ints)
	CreateObjectInfo(_int8s)
	CreateObjectInfo(_int16s)
	CreateObjectInfo(_int32s)
	CreateObjectInfo(_int64s)

	CreateObjectInfo(_uints)
	CreateObjectInfo(_uint8s)
	CreateObjectInfo(_uint16s)
	CreateObjectInfo(_uint32s)
	CreateObjectInfo(_uint64s)

	CreateObjectInfo(_float32s)
	CreateObjectInfo(_float64s)

	CreateObjectInfo(_bool)
	CreateObjectInfo(_string)

}

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
	} else if IsNumber(k) {
		return Number
	} else if isString(k) {
		return String
	} else if k == reflect.Ptr {
		//指针的下级只能是struct
		if t.Elem().Kind() == reflect.Struct {
			return Ptr
		}
		return Invalid
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

// 这里是入口
func CreateObjectInfo(v interface{}) (IObject, error) {
	tp := reflect.TypeOf(v)
	//is pointer head
	if tp.Kind() == reflect.Ptr {
		return createPtrObject(&BaseObject{tp: Ptr, refType: tp})
	}
	return createObjectInfo(tp)
}

func createObjectInfo(t reflect.Type) (IObject, error) {
	//是否已有
	if o, ok := typeMapper[t.Name()]; ok {
		return o, nil
	}

	tp := refType(t)
	if tp == Invalid {
		return nil, model.ErrInvalidObjectType
	}
	obj := &BaseObject{tp: tp, refType: t}
	switch tp {
	case Bool:
		return &BoolObject{BaseObject: obj}, nil
	case Number:
		return &NumberObject{BaseObject: obj}, nil
	case String:
		return &StringObject{BaseObject: obj}, nil
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
			//continue
			return nil, err
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
		if err == model.ErrInvalidObjectType {
			return nil, model.ErrInvalidSliceType
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
		return nil, model.ErrInvalidMapKeyType
	}
	valObj, err := createObjectInfo(obj.refType.Elem())
	if err != nil {
		return nil, err
	}
	if !availableMapValType(valObj) {
		return nil, model.ErrInvalidMapValType
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
		return nil, model.ErrInvalidPtrType
	}
	return &PtrObject{
		HasSubObject: &HasSubObject{
			BaseObject: obj, sub: sub,
		},
	}, nil
}

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
			val.SetInt(convert.StringToInt64(str))
		} else if isUint(val.Kind()) {
			val.SetUint(convert.StringToUint64(str))
		} else {
			val.SetFloat(convert.StringToFloat64(str))
		}
	default:
		return
	}
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
