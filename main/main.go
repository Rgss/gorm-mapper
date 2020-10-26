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
	log.Printf("tem: %v", tem)

	//user.Username = "imp"
	//user.Password = "123456"
	//num := tem.InsertSelective(user)
	//log.Printf("num: %v", num)
	//log.Printf("user: %v", user)

	//where := map[string]interface{}{"id_gt": 32}
	//builder = builder.Debug().Where(where).Build()
	//log.Printf("builder: %v", builder.GetParsedValue())
	//err := tem.SelectOneBySearchBuilder(builder, user)
	//log.Printf("err: %v", err)
	//log.Printf("user: %v", user)

	//users := make([]User, 0)
	//where := map[string]interface{}{"id_gt": 32}
	//builder = builder.Debug().Where(where).Build()
	//err, pager := tem.SelectPageBySearchBuilder(builder, &users)
	//log.Printf("err: %v", err)
	//log.Printf("pager: %v", pager)
	//log.Printf("users: %v", users)

	//fields := []string{"id", "username", "password"}
	//users := make([]User, 0)
	//where := map[string]interface{}{"id_gt": 31}
	//builder = builder.Debug().Where(where).Build()
	//err, pager := tem.PreSelectFields(fields).SelectPageBySearchBuilder(builder, &users)
	//log.Printf("err: %v", err)
	//log.Printf("pager: %v", pager)
	//log.Printf("users: %v", users)

	//fields := []string{"id", "username", "password"}
	//user.Password = "abc"
	//user.Status = 0
	//num := tem.PreUpdateFields(fields).UpdateByPrimaryKey(32, user)
	//log.Printf("num: %v", num)

	// where := gormmapper.WhereBuilder().Put("id_gt", 30).Put("status", 1)
	//entity := new(User)
	//where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
	//log.Printf("builder: %v", builder)
	//builder.Where(where).Debug().Build()
	//tem.SelectBySearchBuilder(builder, &entity)
	//log.Printf("entity: %v", entity)

	//
	m := gormmapper.MapperBuilder()
	gen := gormmapper.MapperGeneratorBuilder(*m)
	gen.EntityPackage("entity")
	gen.EntityPath("E:\\imp\\go\\src\\github.com\\Rgss\\gorm-mapper\\tests\\entity")
	//gen.MapperPackage("repository")
	//gen.MapperPath("E:\\imp\\go\\src\\github.com\\Rgss\\gorm-mapper\\tests\\repository")
	gen.MapperPackage("mapper")
	gen.MapperPath("E:\\imp\\go\\src\\github.com\\Rgss\\gorm-mapper\\tests\\mapper")
	gen.MapperPathAutoSignleton(true)
	gen.Start()

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
		Pass:         "",
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
