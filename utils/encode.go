package utils

import (
	"reflector/encode"
	"reflector/internal"
)

// 提供一个默认的序列化
func Encoder(v interface{}, serializer ...encode.ISerializer) ([]byte, error) {
	if v == nil {
		return nil, encode.ErrEncoder
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
}
