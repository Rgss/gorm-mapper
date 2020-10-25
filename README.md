# gorm-mapper
**简介**  

gorm-mapper 是一个基于gorm的便捷映射器，更加方便的进行数据库操作

<br>

**数据库基本操作**

```
type User struct {  
 	Id         int    `gorm:"not null;primary_key:id;AUTO_INCREMENT" json:"id" form:"id"`  
 	Username   string `gorm:"column:username;not null;default:''" json:"username" form:"username"`  
 	Password   string `gorm:"column:password;not null;default:''" json:"-" form:"password"`  
 	Status     int    `gorm:"column:status;not null;default:1" json:"status" form:"status"`   
 	UpdateTime int    `gorm:"column:update_time; not null; default:0" json:"-" form:"updateTime"`  
 	CreateTime int    `gorm:"column:create_time;not null;default:0" json:"createTime" form:"createTime"`  
 }

 func (User) TableName() string {  
 	&nbsp;&nbsp;&nbsp;&nbsp;return "user"  
 }  
 
 type TestMapper struct {
	gormmapper.Mapper
 }

 user := new(User)
 testMapper := &TestMapper{}
 builder := gormmapper.Builder(&User{})

 # 新增数据
 //user.Username = "imp"
 //user.Password = "123456"
 //num := testMapper.Insert(user)	// 返回受影响记录数

 # 读取单条数据
 //user := new(User)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.SelectOneBySearchBuilder(builder, &user)
 //testMapper.SelectByPrimaryKey(1)
 
 # 分页读取数据
 //users := make([]User{}, 0)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.SelectBySearchBuilder(builder, &users) // 多条读取
 //_, pager := testMapper.SelectPageBySearchBuilder(builder, &users)
 
 # 更新数据
 //user := new(User)
 //user.Password = "123456"
 //testMapper.UpdateByPrimaryKey(1, user)
 //testMapper.UpdateSelectiveByPrimaryKey(1, user) // 选择性字段更新, 为空 0 等不更新
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //testMapper.UpdateBySearchBuilder(builder, user)  // 根据SearchBuilder修改

 // 删除数据
 //testMapper.DeleteByPrimaryKey(1)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1"))
 //builder.Where(where).Limit(1).Debug().Build()
 //testMappser.DeleteBySearchBuilder(builder)

 ```
 
 **WhereBuilder**

 **SearchBuilder**