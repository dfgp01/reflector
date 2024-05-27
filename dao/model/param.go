package model

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
		//nothing
	}
	q.sorts = append(q.sorts, columnField)
}
