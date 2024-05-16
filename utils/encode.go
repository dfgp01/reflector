package utils

import "reflector/encode"

// 提供一个默认的反序列化
func Decoder(dest interface{}, data []byte, s ...encode.ISerializer) error {
	if dest == nil {
		return encode.ErrDecoder
	}
	if len(data) == 0 {
		return nil
	}
	var ser encode.ISerializer
	if len(s) > 0 {
		ser = s[0]
	} else {
		ser = encode.JsonSerializer
	}
	return ser.UnMarshal(data, dest)
}

// 提供一个默认的序列化
func Encoder(v interface{}, s ...encode.ISerializer) ([]byte, error) {
	if v == nil {
		return nil, encode.ErrEncoder
	}
	var ser encode.ISerializer
	if len(s) > 0 {
		ser = s[0]
	} else {
		ser = encode.JsonSerializer
	}
	return ser.Marshal(v)
}
