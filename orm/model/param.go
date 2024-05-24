package model

type Pager struct {
	No    int32 `json:"no,omitempty"`    //当前页码
	Rows  int32 `json:"rows,omitempty"`  //每页数据量
	Count int32 `json:"count,omitempty"` //总共多少页
	Total int32 `json:"total,omitempty"` //总共多少条数据
}

func (p *Pager) Init() {
	if p.No <= 0 {
		p.No = 1
	}
	if p.Rows <= 0 {
		p.Rows = 100
	}
}

func (p *Pager) Fill(total int32) {
	p.Init()
	p.Total = total
	p.Count = total / p.Rows
	mod := total % p.Rows
	if mod > 0 {
		p.Count++
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
	Aggr Aggregation //0=groupby
	Name string
	Val  interface{}
}

type ISortable interface {
	SortNo() int
}

type Sortables []ISortable

func (sb Sortables) Len() int {
	return len(sb)
}

func (sb Sortables) Less(i, j int) bool {
	return sb[i].SortNo() < sb[j].SortNo()
}

func (sb Sortables) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

type SortOrder int

const (
	SortNone SortOrder = iota
	SortAsc
	SortDesc
)

type SortField struct {
	*Column
	or SortOrder
	no int
}

func (s *SortField) SortNo() int {
	return s.no
}

type GroupField struct {
	*Column
	no int
}

func (g *GroupField) SortNo() int {
	return g.no
}

func (q *Query) AddColumn(name string, aggr Aggregation) *Query {
	q.cols = append(q.cols, &Column{Aggr: aggr, Name: name})
	return q
}
