# reflector
golang 反射小助手

#   预期序列化目标

转换目标|转换目标|支持情况|说明
:---:|:---:|:---:|:---:
string|number|支持|使用strconv
string|[]number|支持|用逗号隔开
|||
json|struct、proto|encoding/json
json|[]struct、[]proto|encoding/json
json|map[string]interface{}|encoding/json
json|[]map[string]interface{}|encoding/json
json|map[string][]interface{}|不支持|interface{}已含[]interface{}
json|[]map[string][]interface{}|不支持|同上
json|map[number]interface{}|不支持|json没有非字符串的key
json|map[number][]interface{}|不支持|同上
|||
bytes|proto|google.golang.org/protobuf/proto|protobuf字节码
bytes|[]proto|支持|protobuf字节码+包长
bytes|map[string]proto|缺
bytes|map[number]proto|缺
bytes|map[string][]proto|缺
bytes|map[number][]proto|缺
|||
struct|map[string]interface{}|不支持|可使用github.com/mitchellh/mapstructure
struct|map[number]interface|不支持|struct没有非字符串的key


***

#   使用示例，参考refector/demo包下的代码


*   数值的序列化，float和其他int, uint均支持
  
```
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
```

*   数组的序列化和反序列化，[]float和其他[]int, []uint均支持
  ```
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
```

*   一般proto序列化
```
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
```

*   proto数组的序列化
```
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
```