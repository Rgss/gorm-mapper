package entity

// gorm-mapper auto generate entity
// GatewayApi
type GatewayApi struct {
	Id          string `gorm:"column:id; primaryKey; not null;" json:"id" form:"id"`
	Path        string `gorm:"column:path; not null;" json:"path" form:"path"`
	ServiceId   string `gorm:"column:service_id; default:0;" json:"serviceId" form:"serviceId"`
	Url         string `gorm:"column:url; default:0;" json:"url" form:"url"`
	Retryable   int8   `gorm:"column:retryable; default:0;" json:"retryable" form:"retryable"`
	Enabled     int8   `gorm:"column:enabled; not null;" json:"enabled" form:"enabled"`
	StripPrefix int    `gorm:"column:strip_prefix; default:0;" json:"stripPrefix" form:"stripPrefix"`
	ApiName     string `gorm:"column:api_name; default:0;" json:"apiName" form:"apiName"`
}

// tablename
func (e GatewayApi) TableName() string {
	return "gateway_api"
}
