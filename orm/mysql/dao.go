package mysql

import (
	"errors"
	"fmt"
	"reflector/orm/model"

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

//根据条件查询

func (d *Dao) OrmQuery(q *model.Query) error {
	q.Hook(func(cols []*model.Column, conds *model.Filter, group []*model.GroupField, having *model.Filter) {
		if len(cols) > 0 {
			var selectCols []string
			for _, col := range cols {
				selectCols = append(selectCols, aggrSelect(col.Name, col.Aggr))
			}
			q.engine.Select(aggrs)
		}
		if len(conds) > 0 {
			q.engine.Where(conds[0].field.Name, conds[0].field.Val)
		}
	})
	return q.Err()
}

func aggrSelect(name string, aggr model.Aggregation) string {
	switch aggr {
	case model.AggrCount:
		return fmt.Sprintf("count(%s) as %s", name, name)
	case model.AggrCountDistinct:
		return fmt.Sprintf("count(distinct %s) as %s", name, name)
	case model.AggrSum:
		return fmt.Sprintf("sum(%s) as %s", name, name)
	case model.AggrAvg:
		return fmt.Sprintf("avg(%s) as %s", name, name)
	case model.AggrMax:
		return fmt.Sprintf("max(%s) as %s", name, name)
	case model.AggrMin:
		return fmt.Sprintf("min(%s) as %s", name, name)
	}
	return name
}
