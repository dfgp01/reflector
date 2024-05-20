package model

import "errors"

var (
	ErrNotClassType      = errors.New("not class type")
	ErrInvalidObjectType = errors.New("invalid object type")
	ErrInvalidSliceType  = errors.New("invalid slice type")
	ErrInvalidMapKeyType = errors.New("invalid map key type")
	ErrInvalidMapValType = errors.New("invalid map value type")
	ErrInvalidPtrType    = errors.New("invalid ptr type")
	ErrCheckType         = errors.New("check type error")

	ErrMaxPacket = errors.New("max packet size")
	ErrMaxSeg    = errors.New("max seg size")

	ErrProtobuf = errors.New("can not convert protobuf message")
	ErrDecoder  = errors.New("dest must be pointer")
	ErrEncoder  = errors.New("obj is nil")

	ErrNotNumberSlice = errors.New("not number or number slice")
)
