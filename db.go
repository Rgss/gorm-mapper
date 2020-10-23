package gormmapper

import (
	"fmt"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type DB struct {
	DefaultName string
	pool        map[string]*gorm.DB
}

/**
 * connection url
 * @param
 * @return 
 */
func (mc *DB) getUrl(config *DBConfig) string {
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s", config.User, config.Pass, config.Host, config.Port, config.DbName, config.Charset)
	return url
}

/**
 * 获取数据库连接
 * @param
 * @return 
 */
func (mc *DB) DB(args... string) *gorm.DB {
	dbName := mc.DefaultName
	for _, val := range args {
		dbName = val
	}

	mapID := mc.DBID(dbName)
	if _, ok := mc.pool[mapID]; !ok {
		log.Printf("[error] db %v is not exists", dbName)
	}

	return mc.pool[mapID]
}


/**
 * id
 * @param
 * @return 
 */
func (mc *DB) DBID(dbName string) string {
	dbName = "db_" + md5(dbName)
	return dbName
}
