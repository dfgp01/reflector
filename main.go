package main

import (
	"fmt"
	"reflect"
	"reflector/encode"
	"reflector/internal"
)

type (
	User struct {
		Age      int
		Name     string
		Children []*User
		Parent   *User
		Mate     map[string]*User
		Mates    map[string][]*User
		Adv      map[interface{}]interface{}
	}
	IPer interface {
		Peak([]byte) string
	}

	Ali IPer
)

func main() {

	ref5()
}

func ref() {
	var u1 map[string]*User
	var u2 map[string][]*User

	t1 := reflect.TypeOf(u1)
	v1 := reflect.ValueOf(u1)

	t2 := reflect.TypeOf(u2)
	v2 := reflect.ValueOf(u2)

	fmt.Printf(" 1: %v,			2: %v\n", u1, u2)
	fmt.Printf("type1: %v,		type2: %v\n", t1.Kind(), t2.Kind())
	fmt.Printf("value1: %v, 	value2: %v\n", v1.Kind(), v2.Kind())

	e1, e2 := v1.Elem(), v2.Elem()
	fmt.Printf("e1: %v,		e2: %v\n", e1.Kind(), e2.Kind())
	fmt.Println(e1.IsValid(), e2.IsValid()) //true, false
	//fmt.Println(e1.IsZero(), e2.IsZero())   //e2.IsZero()和IsNil()都会崩溃
}

func ref2() {
	var (
		a []int
		b = []int64{1, 2, 3}
	)

	t1 := reflect.TypeOf(a)
	v1 := reflect.ValueOf(a)

	t2 := reflect.TypeOf(b)
	v2 := reflect.ValueOf(b)

	//slice, slice, true, true, true
	fmt.Printf("a: tk: %v, vk: %v, valid: %v, zero: %v, nil: %v\n", t1.Kind(), v1.Kind(), v1.IsValid(), v1.IsZero(), v1.IsNil())
	//slice, slice, true, false, false
	fmt.Printf("b: tk: %v, vk: %v, valid: %v, zero: %v, nil: %v\n", t2.Kind(), v2.Kind(), v2.IsValid(), v2.IsZero(), v2.IsNil())

	fmt.Println()
	fmt.Println("--------------------------elem---------------------------")
	fmt.Println()

	t1 = t1.Elem()
	//v1.Elem() will panic, isValid:true, isZero:true, isNil:true
	t2 = t2.Elem()
	//v2.Elem() will panic,
	l := v2.Len()
	v2 = v2.Index(2)

	//int
	fmt.Printf("a: tk: %v, vk: zero and nil\n", t1.Kind())
	//int64, 3, true, false, isNil() will panic
	fmt.Printf("b: tk: %v, len: %v, index(2){ valid: %v, zero: %v }\n", t2.Kind(), l, v2.IsValid(), v2.IsZero())

	// var a = [3]*User{Age int}	扩展 a = []interface {User, Map, Slice}
	// tp = slice -> clz 没了
	// v = slice -> clz but [3]
}

func ref3() {
	var a User // User zero=true,  User{} zero=true, User{Age:1} zero=false，IsZero()的逻辑...
	var b *User

	t1 := reflect.TypeOf(a)
	v1 := reflect.ValueOf(a)

	t2 := reflect.TypeOf(b)
	v2 := reflect.ValueOf(b)

	//struct, struct, true, true, isNil() will panic
	fmt.Printf("a: tk: %v, vk: %v, valid: %v, zero: %v\n", t1.Kind(), v1.Kind(), v1.IsValid(), v1.IsZero())
	//ptr, ptr, true, true, true
	fmt.Printf("b: tk: %v, vk: %v, valid: %v, zero: %v, nil: %v\n", t2.Kind(), v2.Kind(), v2.IsValid(), v2.IsZero(), v2.IsNil())

	fmt.Println()
	fmt.Println("--------------------------elem---------------------------")
	fmt.Println()

	t2 = t2.Elem()
	v2 = v2.Elem()
	//struct, invalid, false, isZero() and isNil will panic
	fmt.Printf("b: tk: %v, vk: %v, valid: %v\n", t2.Kind(), v2.Kind(), v2.IsValid())
}

func ref4() {
	//var a *User
	//var b = User{Age: 1}
	//var c = &User{Age: 123, Name: "123"}

	// d := []float32{1.1, 2.2, 3.3, 4.3444}
	// bs, err := encode.StringSerializer.Marshal(d)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(bs, string(bs))

	var d []uint16
	bs := []byte("18,58,19,5")
	err := encode.NumberArraySerializer.UnMarshal(bs, &d)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)
}

func ref5() {
	var a []map[internal.IObject][]*User
	t(a)  //[]uint16
	t(&a) //*[]uint16
}

func t(v interface{}) {
	aa := reflect.TypeOf(v).String()
	fmt.Println(aa)
}
