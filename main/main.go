package main

import (
	"github.com/Rgss/gorm-mapper"
	"gorm.io/gorm"
	"log"
)

func main() {

	initDB()

	user := new(User)
	tem := &TEM{}
	builder := gormmapper.SearcherBuilder(&user)
	log.Printf("builder: %v", builder)
	log.Printf("tem: %v", tem)

	//user.Username = "imp"
	//user.Password = "123456"
	//num := tem.InsertSelective(user)
	//log.Printf("num: %v", num)
	//log.Printf("user: %v", user)

	//where := map[string]interface{}{"id_gt": 32}
	//builder = builder.Debug().Where(where).Build()
	//log.Printf("builder: %v", builder.GetParsedValue())
	//err := tem.SelectOneBySearcher(builder, user)
	//log.Printf("err: %v", err)
	//log.Printf("user: %v", user)

	//users := make([]User, 0)
	//where := map[string]interface{}{"id_gt": 32}
	//builder = builder.Debug().Where(where).Build()
	//err, pager := tem.SelectPageBySearcher(builder, &users)
	//log.Printf("err: %v", err)
	//log.Printf("pager: %v", pager)
	//log.Printf("users: %v", users)

	//fields := []string{"id", "username", "password"}
	//users := make([]User, 0)
	//where := map[string]interface{}{"id_gt": 31}
	//builder = builder.Debug().Where(where).Build()
	//err, pager := tem.PreSelectFields(fields).SelectPageBySearcher(builder, &users)
	//log.Printf("err: %v", err)
	//log.Printf("pager: %v", pager)
	//log.Printf("users: %v", users)

	//fields := []string{"Username", "Password"}
	//user.Password = "abc"
	//user.Status = 0
	//num := tem.PreUpdateFields(fields).UpdateByPrimaryKey(31, user)
	//log.Printf("num: %v", num)

	// where := gormmapper.WhereBuilder().Put("id_gt", 30).Put("status", 1)
	//entity := new(User)
	//where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
	//log.Printf("builder: %v", builder)
	//builder.Where(where).Debug().Build()
	//tem.SelectBySearcher(builder, &entity)
	//log.Printf("entity: %v", entity)

	//account := &entity.Account{}
	//mapper.AccountMapper.SelectByPrimaryKey(1, account)
	//log.Printf("account: %v", account)

	//
	m := gormmapper.MapperBuilder()
	gen := gormmapper.MapperGeneratorBuilder(*m)
	//gen.MapperPackage("mapper")
	//gen.MapperPath("E:\\imp\\go\\src\\github.com\\Rgss\\gorm-mapper\\tests\\mapper")
	//gen.MapperPathAutoSignleton(true)
	//gen.Start()

	gen.EntityPackage("entity")
	gen.EntityPath("/Users/zhanglong/data/go/src/github.com/Rgss/gorm-mapper/main/entity")
	gen.MapperPackage("mapper")
	gen.MapperPath("/Users/zhanglong/data/go/src/github.com/Rgss/gorm-mapper/main/mapper")
	gen.MapperPathAutoSignleton(true)
	gen.Start()

	log.Printf("tem: %v", tem)
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
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		EnableLog:    false,
	}

	gormmapper.CreateConnection("default", config, &gorm.Config{})
}
