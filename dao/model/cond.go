package model

type LogicCondition int

const (
	CondNone LogicCondition = iota //默认，只有filter<=1时
	CondAnd
	CondOr
)

type MatchCondition int

const (
	MatchErr         MatchCondition = iota
	MatchEq                         //=
	MatchNe                         //!=
	MatchLt                         //<
	MatchLe                         //<=
	MatchGt                         //>
	MatchGe                         //>=
	MatchRange                      // between(n, m)	暂时不做
	MatchIn                         // in [...]
	MatchNotIn                      // !in [...]
	MatchFuzzyAll                   // like %str%
	MatchFuzzyPrefix                //like str%
	MatchFuzzySuffix                // like %str
)

type Filter struct {
	Lc  LogicCondition
	Mc  MatchCondition
	Sub []*Filter
	Col *Column
}

// and 至少提供两组filter
func And(filters ...*Filter) *Filter {
	return &Filter{
		Lc:  CondAnd,
		Sub: filters,
	}
}

// or 至少提供两组filter
func Or(filters ...*Filter) *Filter {
	return &Filter{
		Lc:  CondOr,
		Sub: filters,
	}
}

// 暂时这样处理
func Equal(k string, v interface{}) *Filter { return filterMatch(k, v, MatchEq) }

func NotEqual(k string, v interface{}) *Filter { return filterMatch(k, v, MatchNe) }

func LessThan(k string, v interface{}) *Filter { return filterMatch(k, v, MatchLt) }

func LessEqual(k string, v interface{}) *Filter { return filterMatch(k, v, MatchLe) }

func GreaterThan(k string, v interface{}) *Filter { return filterMatch(k, v, MatchGt) }

func GreaterEqual(k string, v interface{}) *Filter { return filterMatch(k, v, MatchGe) }

func Range(k string, v interface{}) *Filter { return filterMatch(k, v, MatchRange) }

func In(k string, v interface{}) *Filter { return filterMatch(k, v, MatchIn) }

func NotIn(k string, v interface{}) *Filter { return filterMatch(k, v, MatchNotIn) }

func Fuzzy(k string, v interface{}) *Filter { return filterMatch(k, v, MatchFuzzyAll) }

func FuzzyPrefix(k string, v interface{}) *Filter { return filterMatch(k, v, MatchFuzzyPrefix) }

func FuzzySuffix(k string, v interface{}) *Filter { return filterMatch(k, v, MatchFuzzySuffix) }

func filterMatch(k string, v interface{}, match MatchCondition) *Filter {
	return &Filter{
		Col: &Column{
			Name: k,
			Val:  v,
		},
		Mc: match,
	}
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
