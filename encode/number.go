package encode

import (
	"fmt"
	"reflect"
	"reflector/internal"
	"reflector/tools"
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

	t, val, err := internal.ReadIn(v, false)
	if err != nil {
		return nil, err
	}

	//is number?
	if internal.IsNumber(t.Kind()) {
		return []byte(fmt.Sprintf("%v", val.Interface())), nil
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

	return nil, ErrNotNumberSlice
}

// &[]int, &[]float32...
func (s *numberSerializer) UnMarshal(data []byte, dest interface{}) error {
	if len(data) == 0 {
		return nil
	}
	t, v, err := internal.ReadIn(dest, true)
	if err != nil {
		return err
	}
	head, t := t, t.Elem()

	//is number?
	if internal.IsNumber(t.Kind()) {
		number := tools.StringToNumber(string(data), t.Kind())
		v.Elem().Set(reflect.ValueOf(number))
		return nil
	}

	//is slice number?
	if internal.IsNumberSlice(t) {
		var (
			//ptr->slice->number
			kind               = t.Elem().Kind()
			numberStrs         = strings.Split(string(data), ",")
			numbersUnknownType []interface{}
		)

		for _, numberStr := range numberStrs {
			numbersUnknownType = append(numbersUnknownType, tools.StringToNumber(numberStr, kind))
		}
		internal.MakeSliceAndAppend(head, v, numbersUnknownType...)
		return nil
	}

	return ErrNotNumberSlice
}
