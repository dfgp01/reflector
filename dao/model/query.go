package model

import (
	"reflect"
	"reflector/internal"
	"reflector/tools"
	"sort"
	"strings"
)

type Query struct {
	err error
	p   *Pager

	cols   []*Column
	groups []*Column
	sorts  []*Column

	havings    []*Filter
	rootHaving *Filter

	conds    []*Filter
	rootCond *Filter
}

func (q *Query) Err() error {
	return q.err
}

func (q *Query) Pager() *Pager {
	return q.p
}

func (q *Query) Columns() []*Column { return q.cols }
func (q *Query) Cond() *Filter      { return q.rootCond }
func (q *Query) Groups() []*Column  { return q.groups }
func (q *Query) Having() *Filter    { return q.rootHaving }
func (q *Query) Sorts() []*Column   { return q.sorts }

func NewQuery(pager ...*Pager) *Query {
	if len(pager) > 0 {
		return &Query{p: pager[0]}
	}
	return &Query{}
}

func (q *Query) SetPage(p *Pager) *Query {
	p.init()
	q.p = p
	return q
}

// 反射解析模型，转换为sql查询参数，需要带有dao-tag的struct
func (q *Query) Model(param interface{}) *Query {
	q.err = internal.StructIter(param, func(field reflect.StructField, value reflect.Value) {
		columnField := &Column{
			Name: tools.CamelToSnakeCase(field.Name),
			Val:  value.Interface(),
		}

		options := tagStrToMap(field.Tag.Get(tagName))

		//dao-tags
		for option, exp := range options {

			//有值才处理
			if !value.IsZero() {
				switch option {
				case equal, notEqual, in, notIn:
					//字符串、数值皆可
					q.and(columnField, option)
				case greaterThan, greaterEqual, lessThan, lessEqual:
					//仅限数值
					q.and(columnField, option)
				case fuzzyAll, fuzzyPrefix, fuzzySuffix:
					//仅限字符串
					q.and(columnField, option)
				case having:
					q.havingAnd(columnField, exp)
				}
			}

			switch option {
			case show, sum, count, min, max, avg, distinct:
				q.addColumn(columnField, option, exp)
			case groupBy:
				var groupNo int
				if exp != "" {
					groupNo = tools.StringToInt(exp)
				}
				q.addGroup(columnField, groupNo)
			case sortA, sortD:
				var sortNo int
				if exp != "" {
					sortNo = tools.StringToInt(exp)
				}
				q.addSort(columnField, option, sortNo)
			}
		}

	})
	if q.err != nil {
		return q
	}
	sort.Slice(q.groups, func(i, j int) bool {
		return q.groups[i].groupNo < q.groups[j].groupNo
	})
	sort.Slice(q.sorts, func(i, j int) bool {
		return q.sorts[i].sortNo < q.sorts[j].sortNo
	})
	if len(q.conds) > 0 {
		q.rootCond = And(q.conds...)
	}
	if len(q.havings) > 0 {
		q.rootHaving = And(q.havings...)
	}
	return q
}

func tagStrToMap(tag string) map[string]string {
	kvs := strings.Split(tag, splitor)
	if len(kvs) == 0 {
		return nil
	}
	m := make(map[string]string)
	for _, val := range kvs {
		kv := strings.Split(val, kvPair)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		} else {
			m[kv[0]] = ""
		}
	}
	return m
}
