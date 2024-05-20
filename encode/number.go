package encode

import (
	"fmt"
	"reflect"
	"reflector/internal"
	"reflector/model"
	"strings"
)

var (
	NumberSerializer = &numberSerializer{}
)

type (
	numberSerializer struct{}
)

// []int, []int32... []uint64...
func (s *numberSerializer) Marshal(v interface{}) ([]byte, error) {

	t, val := internal.TV(v)

	//is number?
	if internal.IsNumber(t.Kind()) {
		return []byte(fmt.Sprintf("%v,", val.Interface())), nil
	}

	//is slice number?
	if internal.IsNumberSlice(t) {
		l := val.Len()
		if l == 0 {
			return nil, nil
		}
		builder := strings.Builder{}
		for i := 0; i < l; i++ {
			numberStr := fmt.Sprintf("%v,", val.Index(i).Interface())
			builder.WriteString(numberStr)
		}
		result := builder.String()
		result = result[:len(result)-1] //去掉最后一个逗号
		return []byte(result), nil
	}

	return nil, model.ErrNotNumberSlice
}

// &[]int, &[]float32...
func (s *numberSerializer) UnMarshal(data []byte, dest interface{}) error {
	if len(data) == 0 {
		return nil
	}
	t, v := reflect.TypeOf(dest), reflect.ValueOf(dest)
	if t.Kind() != reflect.Ptr {
		return ErrNotPtrType
	}
	tp := t.Elem()

	//is number?
	if internal.IsNumber(tp.Kind()) {
		number := internal.StringToNumber(string(data), tp.Kind())
		v.Elem().Set(reflect.ValueOf(number))
		return nil
	}

	//is slice number?
	if internal.IsNumberSlice(tp) {
		var (
			//ptr->slice->number
			kind               = tp.Elem().Kind()
			numberStrs         = strings.Split(string(data), ",")
			numbersUnknownType []interface{}
		)

		for _, numberStr := range numberStrs {
			numbersUnknownType = append(numbersUnknownType, internal.StringToNumber(numberStr, kind))
		}
		internal.MakeSliceAndAppend(t, val, numbersUnknownType...)
		return nil
	}

	return ErrNotNumberSlice

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
