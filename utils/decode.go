package utils

import (
	"reflector/encode"
	"reflector/internal"
	"reflector/model"
	"strings"
)

// 提供一个默认的反序列化
func Decoder(data []byte, dest interface{}, serializer ...encode.ISerializer) error {
	if dest == nil {
		return model.ErrDecoder
	}
	if len(data) == 0 {
		return nil
	}

	var ser encode.ISerializer
	if len(serializer) > 0 {
		ser = serializer[0]
	} else {
		ser = encode.JsonSerializer
	}

	wrap, err := internal.ReadIn(dest)
	if err != nil {
		return err
	}
	head, tp := wrap.RefTp, wrap.RefTp
	if wrap.HeadPtr {
		tp = head.Elem()
	}

	if internal.IsString(tp.Kind()) {
		return []byte(v.(string)), nil
	} else if internal.IsStringSlice(tp) {
		return []byte(strings.Join(v.([]string), ",")), nil
	}

	if internal.IsNumber(tp.Kind()) || internal.IsNumberSlice(tp) {
		ser = encode.NumberSerializer
	}

	if internal.IsPtrStructSlice(tp) {
		if ser == encode.ProtoSerializer {
			ser = encode.ProtoSliceSerializer
		}
	}

	return ser.UnMarshal(data, dest)
}
