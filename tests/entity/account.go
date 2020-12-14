package entity

// gorm-mapper auto generate entity
// Account
type Account struct {
	Id         int    `gorm:"column:id; primaryKey; not null; autoIncrement;" json:"id" form:"id"`
	Username   string `gorm:"column:username; not null;" json:"username" form:"username"`
	Password   string `gorm:"column:password; not null;" json:"password" form:"password"`
	Nickname   string `gorm:"column:nickname; not null;" json:"nickname" form:"nickname"`
	Phone      int    `gorm:"column:phone; not null;" json:"phone" form:"phone"`
	Avatar     string `gorm:"column:avatar; not null;" json:"avatar" form:"avatar"`
	Desc       string `gorm:"column:desc; not null;" json:"desc" form:"desc"`
	Group      int8   `gorm:"column:group; not null;" json:"group" form:"group"`
	Status     int8   `gorm:"column:status; not null;" json:"status" form:"status"`
	UpdateTime int    `gorm:"column:update_time; not null;" json:"updateTime" form:"updateTime"`
	CreateTime int    `gorm:"column:create_time; not null;" json:"createTime" form:"createTime"`
}

// tablename
func (e Account) TableName() string {
	return "account"
}
