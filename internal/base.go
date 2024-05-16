package internal

import (
	"reflect"
)

// 定义反射对象信息
type IObject interface {

	//是否有效类型
	ValidType() bool

	//声明的类型
	Type() ObjType

	//反射的类型对象
	RefType() reflect.Type

	//子类型，ptr, slice, map用，其中map返回map-val类型
	SubType() ObjType

	//子结构，ptr, slice, map用，其中map返回map-valType
	Sub() IObject

	//key类型，map用
	KeyType() ObjType

	//value类型，map用
	ValType() ObjType
}

// 基础类型对象
type BaseObject struct {
	tp      ObjType
	refType reflect.Type
}

func (o *BaseObject) ValidType() bool       { return o.tp != Invalid }
func (o *BaseObject) Type() ObjType         { return o.tp }
func (o *BaseObject) RefType() reflect.Type { return o.refType }
func (o *BaseObject) Make(v interface{})    {}
func (o *BaseObject) SubType() ObjType      { return Invalid }
func (o *BaseObject) Sub() IObject          { return nil }
func (o *BaseObject) KeyType() ObjType      { return Invalid }
func (o *BaseObject) ValType() ObjType      { return Invalid }

// 字段信息
type StructFieldObject struct {
	IObject                     //字段的对象信息
	field   reflect.StructField //field反射信息
}

func (s *StructFieldObject) Name() string {
	return s.field.Name
}

// 结构体类型
type StrutObject struct {
	*BaseObject
	fields []*StructFieldObject
}

func (o *StrutObject) Name() string                 { return o.refType.Name() }
func (o *StrutObject) Fields() []*StructFieldObject { return o.fields }

// 暂时这样写
type InterfaceObject struct {
	*BaseObject
}

type HasSubObject struct {
	*BaseObject
	sub IObject
}

func (o *HasSubObject) SubType() ObjType { return o.sub.Type() }
func (o *HasSubObject) Sub() IObject     { return o.sub }

// 数组、切片类型
type SliceObject struct {
	*HasSubObject
}

func (o *SliceObject) ValidType() bool { return o.sub.ValidType() }

// map类型
type MapObject struct {
	*HasSubObject //sub 就是value-object
	keyObj        IObject
}

func (o *MapObject) ValidType() bool  { return o.keyObj.ValidType() && o.sub.ValidType() }
func (o *MapObject) KeyType() ObjType { return o.keyObj.Type() }
func (o *MapObject) ValType() ObjType { return o.sub.Type() }

type PtrObject struct {
	*HasSubObject //sub 就是指向的类型
}
