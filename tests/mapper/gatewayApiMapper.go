package mapper

import gormmapper "github.com/Rgss/gorm-mapper"

// global instance
var GatewayApiMapper = &gatewayApiMapper{}

// the mapper is generated by gorm-mapper automatically
// gatewayApiMapper
type gatewayApiMapper struct {
	gormmapper.Mapper
}
