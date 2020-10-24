package main

import (
	"github.com/Rgss/gorm-mapper"
	"log"
)

func main() {

	initDB()

	user := new(User)
	tem := &TEM{}
	builder := gormmapper.Builder(&user)
	log.Printf("builder: %v", builder)

	//user.Username = "imp"
	//user.Password = "123456"
	//num := tem.InsertSelective(user)
	//log.Printf("num: %v", num)
	//log.Printf("user: %v", user)

	where := map[string]interface{}{"id": 32}
	builder = builder.Debug().Where(where).Build()
	err := tem.SelectOneBySearchBuilder(builder, user)
	log.Printf("err: %v", err)
	log.Printf("user: %v", user)

}

type TEM struct {
	gormmapper.Mapper
}

type User struct {
	Id         int    `gorm:"not null;primary_key:id;AUTO_INCREMENT" json:"id" form:"id"`
	Username   string `gorm:"column:username;not null;default:''" json:"username" form:"username" binding:"required"`
	Password   string `gorm:"column:password;not null;default:''" json:"-" form:"password"`
	Status     int    `gorm:"column:status;not null;default:1" json:"status" form:"status"`
	UpdateTime int    `gorm:"column:update_time; not null; default:0" json:"-" form:"updateTime"`
	CreateTime int    `gorm:"column:create_time;not null;default:0" json:"createTime" form:"createTime"`
}

func (User) TableName() string {
	return "user"
}

func initDB() {
	config := &gormmapper.DBConfig{
		User:         "root",
		Pass:         "123456",
		Host:         "127.0.0.1",
		Port:         3306,
		DbName:       "test",
		Charset:      "utf8",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		EnableLog:    false,
	}

	configs := make(map[string]*gormmapper.DBConfig)
	configs["default"] = config
	gormmapper.Initialize(configs)
}
