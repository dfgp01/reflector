package mysql

import (
	"errors"
	"fmt"
	"reflector/dao/model"

	"gorm.io/gorm"
)

var (
	ErrParamType    = errors.New("invalid param type")
	ErrConfig       = errors.New("invalid config")
	ErrRowsAffected = errors.New("rows affected error")
)

type Dao struct {
	engine *gorm.DB
}

// 自定义创建
func NewDao(db *gorm.DB) *Dao {
	return &Dao{engine: db}
}

// 默认预设创建
func Open(c *Config) (*Dao, error) {
	if c == nil || c.User == "" || c.Pass == "" || c.Host == "" || c.Db == "" {
		return nil, ErrConfig
	}
	db, err := defaultConn(c)
	if err != nil {
		return nil, err
	}
	return NewDao(db), nil
}

// 自定义操作
func (d *Dao) Driver() *gorm.DB {
	return d.engine
}

// 关闭连接
func (d *Dao) Close() error {
	return close(d.engine)
}

// 自动建表
func (d *Dao) AutoCreateTable(model ...interface{}) error {
	return d.engine.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model...)
}

// 创建记录，若已存在，则更新，可以插入多条
func (d *Dao) Create(record interface{}) error {
	if record == nil {
		return ErrParamType
	}
	result := d.engine.Save(record)
	if result.Error != nil {
		return result.Error
	}

	//受影响的行数不对，可能发生了一些问题
	if result.RowsAffected <= 0 {
		return ErrRowsAffected
	}
	return nil
}

// 根据条件查询，param是带有query标签的struct，dest是gorm的&struct或&[]*struct模型
func (d *Dao) OrmQuery(param, dest interface{}, pager ...*model.Pager) error {
	q := model.NewQuery(pager...).Model(param)
	if q.Err() != nil {
		return q.Err()
	}

	db := d.engine

	//where条件
	cond := q.Cond()
	if cond != nil {
		if len(cond.Sub) > 0 {
			for _, filter := range cond.Sub {
				str, v := condPair(filter)
				db = db.Where(str, v)
			}
		} else {
			str, v := condPair(cond)
			db = db.Where(str, v)
		}
	}

	//having条件
	having := q.Having()
	if having != nil {
		if len(having.Sub) > 0 {
			for _, filter := range having.Sub {
				str, v := condPair(filter)
				db = db.Having(str, v)
			}
		} else {
			str, v := condPair(cond)
			db = db.Having(str, v)
		}
	}

	//group by 字段
	groups := q.Groups()
	if len(groups) > 0 {
		for _, col := range groups {
			db = db.Group(col.Name)
		}
	}

	//分页查询，先查总数
	page := q.Pager()
	if page != nil {
		var total int32
		err := db.Select("COUNT(1) as total").First(&total).Error
		if err != nil {
			return err
		}
		page.Fill(total)
	}

	var ss []string
	for _, c := range q.Columns() {
		ss = append(ss, aggrSelect(c))
	}
	//selectCols := strings.Join(ss, ",")
	db = db.Select(ss)

	sorts := q.Sorts()
	if len(sorts) > 0 {
		for _, col := range sorts {
			db = db.Order(sortBy(col.Name, col.So))
		}
	}

	if page != nil {
		db = db.Limit(int(page.Rows * (page.No - 1))).Offset(int(page.Rows))
	}
	if err := db.Find(dest).Error; err != nil {
		return err
	}

	//pager
	return nil
}

func aggrSelect(c *model.Column) string {
	switch c.Aggr {
	case model.AggrCount:
		return fmt.Sprintf("count(`%s`) as `%s`", c.ColName, c.Name)
	case model.AggrCountDistinct:
		return fmt.Sprintf("count(distinct `%s`) as `%s`", c.ColName, c.Name)
	case model.AggrSum:
		return fmt.Sprintf("sum(`%s`) as `%s`", c.ColName, c.Name)
	case model.AggrAvg:
		return fmt.Sprintf("avg(`%s`) as `%s`", c.ColName, c.Name)
	case model.AggrMax:
		return fmt.Sprintf("max(`%s`) as `%s`", c.ColName, c.Name)
	case model.AggrMin:
		return fmt.Sprintf("min(`%s`) as `%s`", c.ColName, c.Name)
	}
	return c.Name
}

func condPair(filter *model.Filter) (string, interface{}) {
	var (
		str  string
		name = filter.Col.Name
		val  = filter.Col.Val
	)
	switch filter.Mc {
	case model.MatchNe:
		str = fmt.Sprintf("`%s` != ?", name)
	case model.MatchLt:
		str = fmt.Sprintf("`%s` < ?", name)
	case model.MatchLe:
		str = fmt.Sprintf("`%s` <= ?", name)
	case model.MatchGt:
		str = fmt.Sprintf("`%s` > ?", name)
	case model.MatchGe:
		str = fmt.Sprintf("`%s` >= ?", name)
	case model.MatchIn:
		str = fmt.Sprintf("`%s` in (?)", name)
	case model.MatchNotIn:
		str = fmt.Sprintf("`%s` not in (?)", name)
	case model.MatchFuzzyAll:
		str = fmt.Sprintf("`%s` like ?", name)
		val = fmt.Sprintf("%%%v%%", val)
	case model.MatchFuzzyPrefix:
		str = fmt.Sprintf("`%s` like ?", name)
		val = fmt.Sprintf("%v%%", val)
	case model.MatchFuzzySuffix:
		str = fmt.Sprintf("`%s` like ?", name)
		val = fmt.Sprintf("%%%v", val)
	case model.MatchEq:
		fallthrough
	default:
		str = fmt.Sprintf("%s = ?", name)
	}
	return str, val
}

func sortBy(name string, sort model.SortOrder) string {
	switch sort {
	case model.SortDesc:
		return fmt.Sprintf("`%s` %s", name, "desc")
	default:
		return fmt.Sprintf("`%s`", name)
	}
}
