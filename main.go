package main

import (
	"fmt"
	"reflect"
	"reflector/utils"
)

type (
	User struct {
		Age      int                         `json:"age,omitempty"`
		Name     string                      `json:"name,omitempty"`
		Children []*User                     `json:"children,omitempty"`
		Parent   *User                       `json:"parent,omitempty"`
		Mate     map[string]*User            `json:"mates,omitempty"`
		Mates    map[string][]*User          `json:"mates,omitempty"`
		Adv      map[interface{}]interface{} `json:"adv,omitempty"`
	}
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

func ref5() {
	//基础类型测试已通过，下面进行json和protobuf测试
	var a float32 = 111.333
	var b []float64 = []float64{1.32, 2.7778}
	var c = "1,2,3,4,5"
	var d string = "aaaaab"
	var e []string = []string{"1", "2", "3,", "4,7", "5"}

	err := utils.Decoder([]byte("1.222,4,5.655"), &b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a, b, c, d, e)
}
