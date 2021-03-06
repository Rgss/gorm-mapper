# gorm-mapper
**简介**  

gorm-mapper 是一个基于gorm的便捷映射器，更加方便的进行数据库操作

<br>

**依赖**  
1、gorm v2


<br>

**安装**  
```
go get -u github.com/Rgss/gorm-mapper  
```

<br>

## 数据库基本操作    

```
# 初始化
// 创建test数据库连接
config := &gormmapper.DBConfig{
    User:         "root",
    Pass:         "",
    Host:         "127.0.0.1",
    Port:         3306,
    DbName:       "test",
    Charset:      "utf8",
    MaxIdleConns: 10,
    MaxOpenConns: 10,
    EnableLog:    false,
}
gormmapper.CreateConnection("default", config, &gorm.Config{})


type User struct {  
 	Id         int    `gorm:"not null;primary_key:id;AUTO_INCREMENT" json:"id" form:"id"`  
 	Username   string `gorm:"column:username;not null;default:''" json:"username" form:"username"`  
 	Password   string `gorm:"column:password;not null;default:''" json:"-" form:"password"`  
 	Status     int    `gorm:"column:status;not null;default:1" json:"status" form:"status"`   
 	UpdateTime int    `gorm:"column:update_time; not null; default:0" json:"-" form:"updateTime"`  
 	CreateTime int    `gorm:"column:create_time;not null;default:0" json:"createTime" form:"createTime"`  
 }

 func (User) TableName() string {  
    return "user"  
 }  
 
 type TestMapper struct {
	gormmapper.Mapper
 }

 user := new(User)
 testMapper := &TestMapper{}
 builder := gormmapper.SearcherBuilder(&User{})

 # 新增数据
 //user.Username = "imp"
 //user.Password = "123456"
 //num := testMapper.Insert(user)	// 返回受影响记录数

 # 读取单条数据
 //user := new(User)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.SelectOneBySearcher(builder, &user)
 //testMapper.SelectByPrimaryKey(1)
 
 # 读取多条数据 
 //users := make([]User{}, 0)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.SelectBySearcher(builder, &users) // 多条读取
 //_, pager := testMapper.SelectPageBySearcher(builder, &users) //分页读取
 
 # 更新数据
 //user := new(User)
 //user.Password = "123456"
 //testMapper.UpdateByPrimaryKey(1, user)
 //testMapper.UpdateSelectiveByPrimaryKey(1, user) // 选择性字段更新, 为空 0 等不更新
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.UpdateBySearcher(builder, user)  // 根据Searcher修改

 # 删除数据
 //testMapper.DeleteByPrimaryKey(1)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1"))
 //builder.Where(where).Limit(1).Debug().Build()
 //testMappser.DeleteBySearcher(builder)


# 多数据源切换
// 创建account库连接
accountConfig := &gormmapper.DBConfig{
    User:         "root",
    Pass:         "",
    Host:         "127.0.0.1",
    Port:         3306,
    DbName:       "account",
    Charset:      "utf8",
    MaxIdleConns: 10,
    MaxOpenConns: 10,
    EnableLog:    false,
}
gormmapper.CreateConnection("account", accountConfig, &gorm.Config{})

gormmapper.DatabaseSource("account").SelectByPrimaryKey(1)  // account 库
gormmapper.SelectByPrimaryKey(2)  // test 库

 ```
 <br>
 
## Where    
 结构体Where是一个基于map的实现，主要用于搜索条件的构建。
 ```
 # 初始化 
 where := gormmapper.WhereBuilder()

 # 基于map的条件添加
 where.Put("status", 1).Put("create_time_gt", 100)

 # 基于操作符的条件添加，具体操作符可参考operator.go文件
 where.AddOperator(gormmapper.OperatorGT("create_time", 100)).AddOperator(gormmapper.OperatorEQ("status", 1))
 
 ```
 <br>  
   
## Searcher  
 结构体Searcher主要用于mapper中sql各属性映射的构建起。
 ```
 # 初始化
 builder := gormmapper.Builder(&User{});
 builder.Where(where).Debug().Page(1).Sort("id", "desc").Sort("create_time", "desc")

 # 格式化条件生成
 builder.Build()
```
 <br>  
   
 ## MapperGenrator  
 用于根据数据表结构，生成对象实体的生成器。
 ```
 # 初始化
 m := gormmapper.MapperBuilder()
 gen := gormmapper.MapperGeneratorBuilder(*m)
 gen.EntityPackage("entity")  // 设置实体报名
 gen.EntityPath("/data/go/src/github.com/Rgss/gorm-mapper/main/entity") // 设置实体路径
 gen.MapperPackage("mapper")    // 设置映射器包名
 gen.MapperPath("/data/go/src/github.com/Rgss/gorm-mapper/main/mapper") // 设置映射器包路径
 gen.MapperPathAutoSignleton(true)
 gen.Start()
```
 
 <br>
   
 **其它**  
 欢迎issues