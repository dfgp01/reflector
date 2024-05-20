package encode

import (
	"errors"
	"reflect"
	"reflector/internal"

	"google.golang.org/protobuf/proto"
)

const (
	MaxPacketSize = 1 << 19   //单个包最大512k
	MaxSegSize    = 1<<16 - 1 //单个数据段最大64k-1
)

/*
*
仅用于pb数组的二进制协议文
将多个[]byte合并为一个[]byte，类似包的概念，能用于pb数组的序列化
*/
type Packet struct {
	l      int    //内有多少个数据段
	b      []byte //头两项（16位）代表包长，例如[0,3,2,2,4] -> len=3; data=[2,2,4]
	cursor int
}

func (p *Packet) append(data []byte) error {
	size := len(data)
	if size == 0 {
		return nil
	}
	if size > MaxSegSize {
		return ErrMaxSeg
	}
	if len(p.b)+size > MaxPacketSize {
		return ErrMaxPacket
	}
	//16位整数拆成两个byte
	p.b = append(p.b, byte(size>>8), byte(size))
	p.b = append(p.b, data...)
	p.l++
	return nil
}

func (p *Packet) split() [][]byte {
	if len(p.b) == 0 {
		return nil
	}
	p.cursor = 0
	var bs [][]byte
	for {
		b := p.next()
		if len(b) == 0 {
			break
		}
		bs = append(bs, b)
		p.l++
	}
	return bs
}

func (p *Packet) next() []byte {
	i := p.cursor
	if i >= len(p.b) {
		return nil
	}
	size := int(p.b[i])<<8 + int(p.b[i+1])
	start := p.cursor + 2
	p.cursor = start + size
	return p.b[start:p.cursor]
}

func (p *Packet) data() []byte {
	return p.b[:]
}

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
