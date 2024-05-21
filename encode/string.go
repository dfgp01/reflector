package encode

import (
	"reflect"
	"reflector/internal"
	"reflector/model"
	"strings"
)

var (
	StringSerializer = &stringSerializer{}
)

type (
	stringSerializer struct{}
)

// []string
func (s *stringSerializer) Marshal(v interface{}) ([]byte, error) {

	t, val := internal.TV(v)

	//is string?
	if internal.IsString(t.Kind()) {
		return []byte(val.Interface().(string)), nil
	}

	//is slice string?
	if internal.IsStringSlice(t) {
		return []byte(strings.Join(val.Interface().([]string), ",")), nil
	}

	return nil, model.ErrNotStringSlice
}

// &[]string
func (s *stringSerializer) UnMarshal(data []byte, dest interface{}) error {
	if len(data) == 0 {
		return nil
	}
	t, v := internal.TV(dest)

	//is not ptr
	if t.Kind() != reflect.Ptr {
		return model.ErrInvalidPtrType
	}
	head := t
	t = t.Elem()

	// is string?
	if internal.IsString(t.Kind()) {
		v.Elem().SetString(string(data))
		return nil
	}

	// is slice string?
	if internal.IsStringSlice(t) {
		var (
			ss    = strings.Split(string(data), ",")
			sList []interface{}
		)
		for _, s := range ss {
			sList = append(sList, s)
		}
		internal.MakeSliceAndAppend(head, v, sList...)
		return nil
	}

	return model.ErrNotStringSlice
}
