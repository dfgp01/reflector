package encode

import (
	"errors"
	"reflect"
	"reflector/internal"

	"google.golang.org/protobuf/proto"
)

var (
	ProtoSliceSerializer = &protoSliceSerializer{}
	ErrProtobufSlice     = errors.New("can not convert protobuf slice")
)

type protoSliceSerializer struct{}

// v is []proto.Message
func (s *protoSliceSerializer) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.([]proto.Message)
	if !ok {
		return nil, ErrProtobufSlice
	}
	pk := &Packet{}
	for _, m := range msg {
		seg, err := proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		pk.append(seg)
	}
	if pk.l != len(msg) {
		//错误处理
	}
	return pk.data(), nil
}

// ptr is &[]proto.Message
func (s *protoSliceSerializer) UnMarshal(data []byte, ptr interface{}) error {
	//类型判断：&[]proto.Message
	// msg, ok := ptr.(*[]proto.Message)
	// if !ok {
	// 	return ErrProtobufSlice
	// }

	//如果上面不行，用这种......是否需要接入reflector.internal里的sdk?
	internal.In(ptr)
	tp := reflect.TypeOf(ptr)
	if tp.Kind() != reflect.Ptr {
		return ErrProtobufSlice
	}
	tp = tp.Elem()
	if tp.Kind() != reflect.Slice {
		return ErrProtobufSlice
	}
	//head-ptr here
	tp = tp.Elem()
	if tp.Kind() != reflect.Ptr {
		return ErrProtobufSlice
	}
	if tp.Elem().Kind() != reflect.Struct {
		return ErrProtobufSlice
	}

	pk := &Packet{b: data}
	segs := pk.split()
	if len(segs) == 0 {
		return nil
	}

	// for _, val := range segs {
	// 	//通过反射初始化proto
	// 	err := ProtoSerializer.UnMarshal(val, m)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	//追加到slice
	// 	*msg = append(*msg, m)
	// }
	return nil
}
