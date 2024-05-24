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

	cols  []*Column
	group []*GroupField
	sorts []*SortField

	havings    []*Filter
	rootHaving *Filter

	conds    []*Filter
	rootCond *Filter
}

func (q *Query) Err() error {
	return q.err
}

func (q *Query) SetPage(p *Pager) *Query {
	q.p = p
	return q
}

// 反射解析模型，转换为sql查询参数
func (q *Query) Model(dest interface{}) *Query {
	//传进来一个&User{}，解析，然后换成filters?????
	//没有tag的话没法做

	q.err = internal.StructIter(dest, func(field reflect.StructField, value reflect.Value) {
		ts := strings.Split(field.Tag.Get(tagName), splitor)
		if len(ts) == 0 {
			return
		}

		columnField := &Column{
			Name: convert.CamelToSnakeCase(field.Name),
			Val:  value.Interface(),
		}

		//dao-tags
		for _, option := range ts {
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
				}

				//having条件
				if strings.Contains(option, having) {
					havingPair := strings.Split(option, kvPair)
					if len(havingPair) == 2 {
						q.havingAnd(columnField, havingPair[1])
					}
				}
			}

			//指定column名：
			if strings.Contains(option, column) {
				columnPair := strings.Split(option, kvPair)
				if len(columnPair) == 2 {
					field.Name = columnPair[1]
				}
			}

			switch option {
			case show, sum, count, min, max, avg, distinct:
				q.addColumn(columnField, option)
			}

			//group by 字段
			var (
				groupNo int
			)
			if strings.Contains(option, groupBy) {
				groupPair := strings.Split(option, kvPair)
				if len(groupPair) == 2 {
					groupNo = convert.StringToInt(groupPair[1])
				}
				q.addGroup(columnField, groupNo)
			}

			//排序字段
			var (
				sortStr string
				no      int
			)
			if strings.Contains(option, sortA) || strings.Contains(option, sortD) {

				sortPair := strings.Split(option, kvPair)
				if len(sortPair) == 2 {
					no = convert.StringToInt(sortPair[1])
				}
				sortStr = sortPair[0]
				q.addSort(columnField, sortStr, no)
			}
		}

	})
	if q.err != nil {
		return q
	}
	sort.Slice(q.group, func(i, j int) bool {
		return q.group[i].no < q.group[j].no
	})
	sort.Slice(q.sorts, func(i, j int) bool {
		return q.sorts[i].no < q.sorts[j].no
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

func (q *Query) addGroup(columnField *Column, no int) {
	q.group = append(q.group, &GroupField{columnField, no})
}

func (q *Query) addColumn(columnField *Column, aggr string) {
	switch aggr {
	case sum:
		columnField.Aggr = AggrSum
	case count:
		columnField.Aggr = AggrCount
	case min:
		columnField.Aggr = AggrMin
	case max:
		columnField.Aggr = AggrMax
	case avg:
		columnField.Aggr = AggrAvg
	case distinct:
		columnField.Aggr = AggrCountDistinct
	default:
		columnField.Aggr = AggrNone
	}
	q.cols = append(q.cols, columnField)
}

func (q *Query) addSort(columnField *Column, sort string, no int) {
	s := &SortField{columnField, SortNone, no}
	switch sort {
	case sortA:
		s.or = SortAsc
	case sortD:
		s.or = SortDesc
	default:
		//nothing
	}
	q.sorts = append(q.sorts, s)
}

func newFilter(columnField *Column, sym string) *Filter {
	f := &Filter{col: columnField}
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

func (q *Query) Hook(fn func(cols []*Column, conds []*Filter, having *Filter)) {

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
	for _, val := range q.group {
		fmt.Printf("%s, ", val.Name)
	}
	fmt.Print("\n")

	fmt.Print("sorts: ")
	for _, val := range q.sorts {
		fmt.Printf("%s %v, ", val.Name, val.or)
	}
	fmt.Print("\n")

	fmt.Print("having: %v\n", q.rootHaving)
	for _, val := range q.rootHaving.Sub {
		fmt.Printf(" %v %v %v, \n", val.col.Name, val.Mc, val.col.Val)
	}
	fmt.Print("\n")

	fmt.Printf("cond: %v\n", q.rootCond)
	for _, val := range q.rootCond.Sub {
		fmt.Printf(" %v %v %v, \n", val.col.Name, val.Mc, val.col.Val)
	}
	fmt.Print("\n")
}
