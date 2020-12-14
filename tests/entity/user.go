package entity

// gorm-mapper auto generate entity
// User
type User struct {
	Id       int    `gorm:"column:id; primaryKey; not null; autoIncrement;" json:"id" form:"id"`
	Username string `gorm:"column:username; default:;" json:"username" form:"username"`
	Password string `gorm:"column:password; default:;" json:"password" form:"password"`
	Name     string `gorm:"column:name; default:;" json:"name" form:"name"`
	Age      int    `gorm:"column:age; default:0;" json:"age" form:"age"`
	Nick     string `gorm:"column:nick; default:0;" json:"nick" form:"nick"`
	Created  int    `gorm:"column:created; default:0;" json:"created" form:"created"`
}

// tablename
func (e User) TableName() string {
	return "user"
}
