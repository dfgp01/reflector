package encode

import (
	"errors"
	"fmt"
	"reflector/internal"
	"strings"
)

var (
	NumberArraySerializer = &numberArraySerializer{}
	ErrNotNumberSlice     = errors.New("not number slice")
)

type (
	numberArraySerializer struct{}
)

// []int, []int32... []uint64...
func (s *numberArraySerializer) Marshal(v interface{}) ([]byte, error) {
	wr, err := internal.In(v)
	if err != nil {
		return nil, err
	}
	if err = wr.CheckType(internal.Slice, internal.Number); err != nil {
		return nil, err
	}
	if wr.Val.IsZero() {
		return nil, nil
	}
	l := wr.Val.Len()

	builder := strings.Builder{}
	for i := 0; i < l; i++ {
		numberStr := fmt.Sprintf("%v,", wr.Val.Index(i).Interface())
		builder.WriteString(numberStr)
	}
	result := builder.String()
	result = result[:len(result)-1] //去掉最后一个逗号
	return []byte(result), nil
}

// &[]int, &[]float32...
func (s *numberArraySerializer) UnMarshal(data []byte, ptr interface{}) error {
	if len(data) == 0 {
		return nil
	}

	wr, err := internal.In(ptr)
	if err != nil {
		return err
	}
	if err = wr.CheckType(internal.Ptr, internal.Slice, internal.Number); err != nil {
		return err
	}

	var (
		//ptr->slice->number
		kind               = wr.Root.RefType().Elem().Elem().Kind()
		numberStrs         = strings.Split(string(data), ",")
		numbersUnknownType []interface{}
	)

	for _, numberStr := range numberStrs {
		numbersUnknownType = append(numbersUnknownType, internal.StringToNumber(numberStr, kind))
	}
	internal.MakeSliceAndAppend(wr.Root, wr.Val, numbersUnknownType...)

	return nil
}
