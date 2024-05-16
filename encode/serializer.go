package encode

import (
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	JsonSerializer  = &jsonSerializer{}
	ProtoSerializer = &protoSerializer{}
	ErrProtobuf     = errors.New("can not convert protobuf message")
	ErrDecoder      = errors.New("dest must be pointer")
	ErrEncoder      = errors.New("obj is nil")
)

type (
	ISerializer interface {
		Marshal(v interface{}) ([]byte, error)
		UnMarshal(data []byte, ptr interface{}) error
	}
	jsonSerializer  struct{}
	protoSerializer struct{}
)

func (s *jsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *jsonSerializer) UnMarshal(data []byte, ptr interface{}) error {
	return json.Unmarshal(data, ptr)
}

func (s *protoSerializer) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, ErrProtobuf
	}
	return proto.Marshal(msg)
}

func (s *protoSerializer) UnMarshal(data []byte, ptr interface{}) error {
	msg, ok := ptr.(proto.Message)
	if !ok {
		return ErrProtobuf
	}
	return proto.Unmarshal(data, msg)
}
