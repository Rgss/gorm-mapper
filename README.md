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
 
 type TEM struct {
	gormmapper.Mapper
 }

 user := new(User)
 tem := &TEM{}
 builder := gormmapper.Builder(&user)

 # 新增数据
 //user.Username = "imp"
 //user.Password = "123456"
 //num := tem.Insert(user)	// 返回受影响记录数

 # 读取数据
 // where := gormmapper.WhereBuilder().Put("id_gt", 30).Put("status", 1)
 //entity := new(User)
 //where := gormmapper.WhereBuilder().AddOperator(gormmapper.OperatorGT("id", "1")).AddOperator(gormmapper.OperatorLIKE("username", "%imp%"))
 //builder.Where(where).Debug().Build()
 //tem.SelectBySearchBuilder(builder, &entity)
 
 ```