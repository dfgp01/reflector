package inner

import (
	"fmt"
	"reflector/orm/model"
	"reflector/orm/mysql"
	"time"
)

var c = &mysql.Config{
	Host: "8.134.168.80:3304",
	Db:   "component_power",
	User: "root",
	Pass: "dfaklrejqlg413u43dfhfhs",
}

var dao *mysql.Dao

func init() {
	Open()
}

func Open() {
	dao1, err := mysql.Open(c)
	if err != nil {
		panic(err)
	}
	err = dao1.AutoCreateTable(&TestUser{})
	if err != nil {
		panic(err)
	}
	dao = dao1
}

func Create() {
	//多条记录
	n := time.Now()
	users := []*TestUser{
		{Name: "Jinzhu", Age: 18, Birthday: &n},
		{Name: "Jackson", Age: 19},
		{Name: "Joe", Age: 20},
	}
	users[2].ID = 1

	//采用save()，影响条数是4条，3条insert和1条update，不能凭rowsAffected判断是否成功
	err := dao.Create(users)
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		fmt.Println(user.ID)
	}
}

// 注意 First、Take、Last的区别
func GetOne() {

	var user TestUser
	// 获取第一条记录（主键升序）
	dao.Driver().First(&user)
	// SELECT * FROM users ORDER BY id LIMIT 1;

	// 获取一条记录，没有指定排序字段
	dao.Driver().Take(&user)
	// SELECT * FROM users LIMIT 1;

	// 获取最后一条记录（主键降序）
	dao.Driver().Last(&user)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;

	//如果你想避免ErrRecordNotFound错误，你可以使用Find，比如db.Limit(1).Find(&user)，Find方法可以接受struct和slice的数据。
}

func Query() {
	order := &ProductOrder{
		ProductID: 123, UserId: 456, Price: 10086.25, Amount: 99, Total: 65535.99,
	}
	q := &model.Query{}
	q.Model(order).AddColumn("user_count", model.AggrCount).Debug()
}
