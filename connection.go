package gormmapper

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Connection
type connection struct {
	name    string
	db      *gorm.DB
	config  *gorm.Config
	conPool *gorm.ConnPool
}

// instance
var connections map[string]*connection

/**
 * 返回db
 * @param
 * @return
 */
func (c *connection) DB() *gorm.DB {
	return c.db
}

/**
 * 初始化
 * @param
 * @return
 */
func init() {

	// 实例列表
	connections = make(map[string]*connection)

	// ping
	go Ping()
}

/**
 * 创建DB
 * @param
 * @return
 */
func CreateDB(name string, d *gorm.DB) {
	connection := &connection{
		name: name,
		db:   d,
	}
	connections[name] = connection
}

/**
 * 创建连接
 * @param
 * @return
 */
func CreateConnection(name string, dbConfig *DBConfig, config *gorm.Config) *connection {
	connection := createConnection(name, dbConfig, config)
	if connection == nil {
		panic("create connection error.")
		return nil
	}

	connections[name] = connection
	return connection
}

/**
 * createConnection
 * @param
 * @return
 */
func createConnection(name string, dbConfig *DBConfig, config *gorm.Config) *connection {
	dsn := BuildDSN(dbConfig)
	d, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Printf("gorm.Open err: %v", err.Error())
		return nil
	}

	if dbConfig.EnableLog {
		d.Debug()
	}

	sd, err := d.DB()
	if err != nil {
		log.Printf("gorm.Open.DB() err: %v", err.Error())
		return nil
	}

	if dbConfig.MaxOpenConns > 0 {
		sd.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}

	if dbConfig.MaxIdleConns > 0 {
		sd.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}

	if dbConfig.MaxLifetime > 0 {
		sd.SetConnMaxLifetime(dbConfig.MaxLifetime)
	}

	if dbConfig.MaxIdleTime > 0 {
		sd.SetConnMaxIdleTime(dbConfig.MaxIdleTime)
	}

	connection := &connection{
		name: name,
		db:   d,
	}
	return connection
}

/**
 * 返回连接
 * @param
 * @return
 */
func Connection(args ...string) *connection {
	name := "default"
	for _, val := range args {
		name = val
	}

	if _, ok := connections[name]; !ok {
		s := fmt.Sprintf("the connection %v does not exists.", name)
		panic(s)
		return nil
	}
	return connections[name]
}

/**
 * connection url
 * @param
 * @return
 */
func BuildDSN(config *DBConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s", config.User, config.Pass, config.Host, config.Port, config.DbName, config.Charset)
}
