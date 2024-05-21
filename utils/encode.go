package utils

import (
	"reflector/encode"
	"reflector/internal"
	"reflector/model"
)

// 提供一个默认的序列化
func Encoder(v interface{}, serializer ...encode.ISerializer) ([]byte, error) {
	if v == nil {
		return nil, model.ErrEncoder
	}

	var ser encode.ISerializer
	if len(serializer) > 0 {
		ser = serializer[0]
	} else {
		ser = encode.JsonSerializer
	}

	refTp, _, err := internal.ReadIn(v, false)
	if err != nil {
		return nil, err
	}
	head, tp := refTp, refTp
	if internal.IsPointer(refTp) {
		tp = head.Elem()
	}

	if internal.IsString(tp.Kind()) || internal.IsStringSlice(tp) {
		ser = encode.StringSerializer
	}

	if internal.IsNumber(tp.Kind()) || internal.IsNumberSlice(tp) {
		ser = encode.NumberSerializer
	}

	if internal.IsPtrStructSlice(tp) {
		if ser == encode.ProtoSerializer {
			ser = encode.ProtoSliceSerializer
		}
	}

	return ser.Marshal(v)

	//is number slice?
	//is number?

	//proto
	//proto slice

	//map[string]interface{}

	//(int, int8, 16, 32, 64, f32, f64) * 2 * 2 = 28种
	//目前只需要加NumberSerializer和ProtoExtSerializer即可

	//从最简单的开始，number or []number
	//NumberSerializer
	//然后判断是否指定了ProtoSerializer，然后还要根据情况选择ProtoExtSerializer
	//最后默认就是用JsonSerializer
}
