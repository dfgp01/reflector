package encode

import "errors"

const (
	MaxPacketSize = 1 << 19   //单个包最大512k
	MaxSegSize    = 1<<16 - 1 //单个数据段最大64k-1
)

var (
	ErrMaxPacket = errors.New("max packet size")
	ErrMaxSeg    = errors.New("max seg size")
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
