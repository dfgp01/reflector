package model

import "errors"

type Pager struct {
	No    int32 `json:"no,omitempty"`    //当前页码
	Rows  int32 `json:"rows,omitempty"`  //每页数据量
	Count int32 `json:"count,omitempty"` //总共多少页
	Total int32 `json:"total,omitempty"` //总共多少条数据
}

func (p *Pager) init() {
	if p.No <= 0 {
		p.No = 1
	}
	if p.Rows <= 0 {
		p.Rows = 100
	}
}

func (p *Pager) Fill(total int32) {
	p.init()
	p.Total = total
	p.Count = total / p.Rows
	mod := total % p.Rows
	if mod > 0 {
		p.Count++
	}
	if p.No > p.Count {
		p.No = p.Count
	}
}

type Aggregation int

const (
	AggrNone Aggregation = iota
	AggrCount
	AggrCountDistinct
	AggrSum
	AggrAvg
	AggrMax
	AggrMin
)

type Column struct {
	Aggr    Aggregation //0=groupby
	Name    string
	Val     interface{}
	So      SortOrder //升序降序
	ColName string    //指定列
	groupNo int       //group by 位置顺序
	sortNo  int       //sort 位置顺序
}

type SortOrder int

const (
	SortNone SortOrder = iota
	SortAsc
	SortDesc
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

	fuzzyAll    string = "%"  //全模糊查询
	fuzzyPrefix string = "_%" //前模糊查询
	fuzzySuffix string = "%_" //后模糊查询

	splitor string = ";"
	kvPair  string = "="
	tagName string = "query"
)

func (q *Query) addColumn(columnField *Column, aggr, colName string) {
	if colName != "" {
		columnField.ColName = colName
	} else {
		columnField.ColName = columnField.Name
	}
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

func (q *Query) addGroup(columnField *Column, no int) {
	columnField.groupNo = no
	q.groups = append(q.groups, columnField)
}
func (q *Query) addSort(columnField *Column, sort string, no int) {
	columnField.sortNo = no
	switch sort {
	case sortA:
		columnField.So = SortAsc
	case sortD:
		columnField.So = SortDesc
	default:
		columnField.So = SortAsc
	}
	q.sorts = append(q.sorts, columnField)
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
