package inner

import (
	"fmt"
	"reflector/encode"
	"reflector/utils"
)

type (
	User struct {
		Age      int     `json:"age,omitempty"`
		Name     string  `json:"name,omitempty"`
		Children []*User `json:"children,omitempty"`
		Parent   *User   `json:"parent,omitempty"`
	}
)

// 整型的序列化和反序列化
func IntSerialization() {

	var number int16 = 355
	bs, err := utils.Encoder(number)
	if err != nil {
		fmt.Println(err)
	}

	//result: [51 53 53] 355
	fmt.Println(bs, string(bs))

	str := "65535"
	err = utils.Decoder([]byte(str), &number)
	if err != nil {
		fmt.Println(err)
	}

	//result: 32767
	fmt.Println(number)
}

// 整型数组的序列化和反序列化
func IntSliceSerialization() {

	var number []int32 = []int32{1, 2, 3, 4, 5}
	bs, err := utils.Encoder(number)
	if err != nil {
		fmt.Println(err)
	}

	//result: [49 44 50 44 51 44 52 44 53] 1,2,3,4,5
	fmt.Println(bs, string(bs))

	str := "1234,5678,65535666"
	err = utils.Decoder([]byte(str), &number)
	if err != nil {
		fmt.Println(err)
	}

	//result: [1234 5678 65535666]
	fmt.Println(number)
}

// 浮点型数组的序列化和反序列化
func FloatSliceSerialization() {

	var number []float32 = []float32{1.34, 2.555555, 6.78988}
	bs, err := utils.Encoder(number)
	if err != nil {
		fmt.Println(err)
	}

	//result: [49 46 51 52 44 50 46 53 53 53 53 53 53 44 54 46 55 56 57 56 56] 1.34,2.555555,6.78988
	fmt.Println(bs, string(bs))

	str := "1234.1234,5678.56789,65535666"
	err = utils.Decoder([]byte(str), &number)
	if err != nil {
		fmt.Println(err)
	}

	//result: [1234.1234 5678.568 6.5535664e+07]
	fmt.Println(number)
}

// json序列化，采用默认实现，无需关注
func JsonSerialization() {
	user := &User{
		Name:     "Tom",
		Age:      18,
		Children: []*User{{Name: "Jerry"}},
		Parent:   &User{Name: "Jack"},
	}
	bs, err := utils.Encoder(user)
	if err != nil {
		fmt.Println(err)
	}

	//result: {"age":18,"name":"Tom","children":[{"name":"Jerry"}],"parent":{"name":"Jack"}}
	fmt.Println(string(bs))

	//反序列化
	var user2 User
	err = utils.Decoder(bs, &user2)
	if err != nil {
		fmt.Println(err)
	}

	//result: {Age:18 Name:Tom Children:[0xc000076340] Parent:0xc000076380}
	fmt.Printf("%+v\n", user2)
}

// 一般proto序列化
func ProtoSerialization() {
	stu := &Student{
		Name: "Tom",
		Age:  18,
		Hobbies: []string{
			"football",
			"basketball",
		},
	}
	bs, err := utils.Encoder(stu)
	if err != nil {
		fmt.Println(err)
	}

	//由于没有指定序列器，默认采用json
	//result: {"name":"Tom","age":18,"hobbies":["football","basketball"]}
	fmt.Println(bs, string(bs))

	//采用proto处理
	bs, err = utils.Encoder(stu, encode.ProtoSerializer)
	if err != nil {
		fmt.Println(err)
	}

	//反序列化，记得采用ProtoSerializer
	var stu2 Student
	err = utils.Decoder(bs, &stu2, encode.ProtoSerializer)
	if err != nil {
		fmt.Println(err)
	}

	//result: {state:{NoUnkeyedLiterals:{} DoNotCompare:[] DoNotCopy:[] atomicMessageInfo:0xc00000c408} sizeCache:0 unknownFields:[] Name:Tom Age:18 Hobbies:[football basketball]}
	fmt.Printf("%+v\n", stu2)
}

// proto数组的序列化
func ProtoSliceSerialization() {
	stus := []*Student{
		{
			Name:    "Tom",
			Age:     18,
			Hobbies: []string{"football", "basketball"},
		},
		{
			Name:    "Jerry",
			Age:     20,
			Hobbies: []string{"golang"},
		},
	}

	//采用proto处理，也可以用protoSliceSerializer
	bs, err := utils.Encoder(stus, encode.ProtoSerializer)
	if err != nil {
		fmt.Println(err)
	}

	//反序列化，记得用ProtoSerializer或protoSliceSerializer
	var stus2 []*Student
	err = utils.Decoder(bs, &stus2, encode.ProtoSerializer)
	if err != nil {
		fmt.Println(err)
	}

	//result: [name:"Tom"  age:18  hobbies:"football"  hobbies:"basketball" name:"Jerry"  age:20  hobbies:"golang"]
	fmt.Printf("%+v\n", stus2)
}
