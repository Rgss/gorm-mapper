package entity

// gorm-mapper auto generate entity
// Test
type Test struct {
	Id   int    `gorm:"column:id; default:0;" json:"id" form:"id"`
	Name string `gorm:"column:name; not null;" json:"name" form:"name"`
	Time int    `gorm:"column:time; not null;" json:"time" form:"time"`
}

// tablename
func (e Test) TableName() string {
	return "test"
}
