package model

import (
	"errors"
	"fmt"
	"reflect"
	"reflector/internal"
	"reflector/internal/convert"
	"sort"
	"strings"
)

const (
	sortA string = "asc"  //升序
	sortD string = "desc" //降序

	show    string = "show"   //select查询字段
	column  string = "column" //指定列名
	groupBy string = "group"  //group by 列
	having  string = "having" //having条件

	sum      string = "sum"      //求和函数
	avg      string = "avg"      //平均值函数
	max      string = "max"      //最大值函数
	min      string = "min"      //最小值函数
	distinct string = "distinct" //去重计数
	count    string = "count"    //计数

	equal    string = "eq" //等于
	notEqual string = "ne" //不等于，!=同<>
	in       string = "in" //in查询
	notIn    string = "ni" //not in查询

	lessThan     string = "lt" //小于
	lessEqual    string = "le" //小于等于
	greaterThan  string = "gt" //大于
	greaterEqual string = "ge" //大于等于

	fuzzyAll    string = "%_%" //全模糊查询
	fuzzyPrefix string = "_%"  //前模糊查询
	fuzzySuffix string = "%_"  //后模糊查询

	splitor string = ";"
	kvPair  string = "="

	tagName string = "dao"
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

func (q *Query) Params() ([]*Column, *Filter, []*Column, *Filter, []*Column) {
	return q.cols, q.rootCond, q.groups, q.rootHaving, q.sorts
}

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
		options := strings.Split(field.Tag.Get(tagName), splitor)
		if len(options) == 0 {
			return
		}

		columnField := &Column{
			Name: convert.CamelToSnakeCase(field.Name),
			Val:  value.Interface(),
		}

		//dao-tags
		for _, option := range options {
			var (
				op, exp string
			)
			pairs := strings.Split(option, kvPair)
			op = pairs[0]
			if len(pairs) > 1 {
				exp = pairs[1]
			}

			//有值才处理
			if !value.IsZero() {
				switch op {
				case equal, notEqual, in, notIn:
					//字符串、数值皆可
					q.and(columnField, op)
				case greaterThan, greaterEqual, lessThan, lessEqual:
					//仅限数值
					q.and(columnField, op)
				case fuzzyAll, fuzzyPrefix, fuzzySuffix:
					//仅限字符串
					q.and(columnField, op)
				case having:
					q.havingAnd(columnField, exp)
				}
			}

			switch op {
			case show, sum, count, min, max, avg, distinct:
				q.addColumn(columnField, op, exp)
			case groupBy:
				var groupNo int
				if exp != "" {
					groupNo = convert.StringToInt(exp)
				}
				q.addGroup(columnField, groupNo)
			case sortA, sortD:
				var sortNo int
				if exp != "" {
					sortNo = convert.StringToInt(exp)
				}
				q.addSort(columnField, op, sortNo)
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

func (q *Query) and(columnField *Column, sym string) {
	f := newFilter(columnField, sym)
	if f == nil {
		q.err = errors.New("invalid filter symbol")
		return
	}
	q.conds = append(q.conds, f)
}

func (q *Query) havingAnd(columnField *Column, sym string) {
	f := newFilter(columnField, sym)
	if f == nil {
		q.err = errors.New("invalid having symbol")
		return
	}
	q.havings = append(q.havings, f)
}

func newFilter(columnField *Column, sym string) *Filter {
	f := &Filter{Col: columnField}
	switch sym {
	case equal:
		f.Mc = MatchEq
	case notEqual:
		f.Mc = MatchNe
	case in:
		f.Mc = MatchIn
	case notIn:
		f.Mc = MatchNotIn
	case lessThan:
		f.Mc = MatchLt
	case lessEqual:
		f.Mc = MatchLe
	case greaterThan:
		f.Mc = MatchGt
	case greaterEqual:
		f.Mc = MatchGe
	case fuzzyAll:
		f.Mc = MatchFuzzyAll
	case fuzzyPrefix:
		f.Mc = MatchFuzzyPrefix
	case fuzzySuffix:
		f.Mc = MatchFuzzySuffix
	default:
		//error match
		return nil
	}
	return f
}

func (q *Query) Debug() {

	if q.err != nil {
		fmt.Printf("error: %v\n", q.err)
	}
	fmt.Print("cols: ")
	for _, val := range q.cols {
		fmt.Printf("%s %v, ", val.Name, val.Aggr)
	}
	fmt.Print("\n")

	fmt.Print("group: ")
	for _, val := range q.groups {
		fmt.Printf("%s, ", val.Name)
	}
	fmt.Print("\n")

	fmt.Print("sorts: ")
	for _, val := range q.sorts {
		fmt.Printf("%s %v, ", val.Name, val.So)
	}
	fmt.Print("\n")

	fmt.Printf("having: %v\n", q.rootHaving)
	for _, val := range q.rootHaving.Sub {
		fmt.Printf(" %v %v %v, \n", val.Col.Name, val.Mc, val.Col.Val)
	}
	fmt.Print("\n")

	fmt.Printf("cond: %v\n", q.rootCond)
	for _, val := range q.rootCond.Sub {
		fmt.Printf(" %v %v %v, \n", val.Col.Name, val.Mc, val.Col.Val)
	}
	fmt.Print("\n")
}
