package repository

import gormmapper "github.com/Rgss/gorm-mapper"

// global instance
var UserRepository = &userRepository{}

// the mapper is generated by gorm-mapper automatically
// mapper
type userRepository struct {
	gormmapper.Mapper
}
