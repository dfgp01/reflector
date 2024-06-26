package utils

import (
	"reflector/encode"
	"reflector/internal"
)

// 提供一个默认的反序列化
func Decoder(data []byte, dest interface{}, serializer ...encode.ISerializer) error {
	if dest == nil {
		return encode.ErrDecoder
	}
	if len(data) == 0 {
		return nil
	}

	//默认为json处理
	var ser encode.ISerializer
	if len(serializer) > 0 {
		ser = serializer[0]
	} else {
		ser = encode.JsonSerializer
	}

	refTp, _, err := internal.ReadIn(dest, true)
	if err != nil {
		return err
	}
	tp := refTp.Elem()

	//映射目标是字符串或字符串数组
	if internal.IsString(tp.Kind()) || internal.IsStringSlice(tp) {
		ser = encode.StringSerializer
	}

	//映射目标为数字或数字数组
	if internal.IsNumber(tp.Kind()) || internal.IsNumberSlice(tp) {
		ser = encode.NumberSerializer
	}

	//映射目标为结构体数组，且指定了序列化器为Protobuf
	if internal.IsPtrStructSlice(tp) {
		if ser == encode.ProtoSerializer {
			ser = encode.ProtoSliceSerializer
		}
	}

	return ser.UnMarshal(data, dest)
}
