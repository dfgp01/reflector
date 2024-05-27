package inner

import (
	"fmt"
	"reflector/dao/mysql"
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
	// err = dao1.AutoCreateTable(&TestUser{}, &ProductOrder{})
	// if err != nil {
	// 	panic(err)
	// }
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

func CreateOrder() {
	//多条记录
	orders := []*ProductOrder{
		{ProductID: 1000, UserId: 4, Price: 10.25, Amount: 10},
		{ProductID: 1000, UserId: 4, Price: 10.25, Amount: 20},
		{ProductID: 1001, UserId: 5, Price: 5.36, Amount: 99},
	}
	for i := 0; i < 10; i++ {
		orders = append(orders, orders[:]...)
	}

	//calc total
	for _, val := range orders {
		val.Total = val.Price * float32(val.Amount)
	}

	err := dao.Create(orders)
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range orders {
		fmt.Println(user.ID)
	}

	fmt.Println(len(orders))
}

func Query() {
	param := &ManagerOperLog{Method: "/group", OperId: 1000}
	var result []*ManagerOperLog
	err := dao.OrmQuery(param, &result)
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range result {
		fmt.Println(val)
	}
}

func Check() {
	var tables []string
	var total int
	var result []string
	db := dao.Driver()
	err := db.Raw("SHOW TABLES").Scan(&tables).Error
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range tables {
		err = db.Select("COUNT(1)").Table(val).Scan(&total).Error
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, fmt.Sprintf("%s %d\n", val, total))
	}
	fmt.Println(result)
}
