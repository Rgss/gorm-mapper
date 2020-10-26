package gormmapper

import "time"

// database config
type DBConfig struct {
	User         string        // 数据库用户名
	Pass         string        // 数据库密码
	Host         string        // 主机
	Port         int           // 端口
	DbName       string        // 数据库名
	Charset      string        // 字符编码
	MaxIdleConns int           // 最大空闲连接数
	MaxOpenConns int           // 最大打开连接数
	MaxLifetime  time.Duration // 连接最大可复用的时
	MaxIdleTime  time.Duration // 链接最大空闲时间
	EnableLog    bool          // 是否开启debug
}
