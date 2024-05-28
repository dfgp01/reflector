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

	//param
	ManagerOperLog struct {
		Method        string `query:"show;_%;group=3"`
		RequestMethod string `query:"distinct"`
		Status        int    `query:"distinct"`
		OperId        int    `query:"count;having=gt"`
		Group         string `query:"count"`
		Url           string `query:"distinct=oper_url"`
	}
)
