package inner

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type (
	TestUser struct {
		gorm.Model
		Name        string       `dao:"match=fuzzy_all,sort=asc"`
		Age         uint8        `dao:"match=gt,sort=asc"` // 一个未签名的8位整数
		Birthday    *time.Time   // A pointer to time.Time, can be null
		ActivatedAt sql.NullTime // Uses sql.NullTime for nullable time fields
		//DefaultAt   time.Time    // 非指针time类型，测试空值
	}

	ProductOrder struct {
		gorm.Model
		ProductID uint    `dao:"show;group=1"`
		UserId    uint    `dao:"show;group=2;asc=2"`
		Price     float32 `dao:"gt"`
		Amount    uint
		Total     float32 `dao:"sum;lt;desc;having=gt"`
	}
)
