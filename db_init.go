package gormmapper

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var connection *DB

/**
 * 初始化数据库
 * @param
 * @return 
 */
func Initialize(configs map[string]*DBConfig) {
	// 实例化
	NewDB(configs)

	// ping
	go Ping()
}


/**
 * 新建db
 * @param
 * @return 
 */
func NewDB(configs map[string]*DBConfig) *DB {
	connection = &DB{
		DefaultName: "default",
	}

	pool := make(map[string]*gorm.DB, 0)
	for key, config := range configs {
		url 	:= connection.getUrl(config)
		maxIdleConns := config.MaxIdleConns
		maxOpenConns := config.MaxOpenConns
		enableLog 	 := config.EnableLog

		cfg := &gorm.Config{};
		db, err := gorm.Open(mysql.Open(url), cfg);
		if err != nil {
			//app.Logger().Errorf("连接数据库失败 error: %s", err.Error())
			panic("连接数据库失败 error: " + err.Error())
			return nil
		}

		if enableLog {
			db.Debug()
		}

		d, _ := db.DB()
		d.SetMaxIdleConns(maxIdleConns)
		d.SetMaxOpenConns(maxOpenConns)

		//if err = db.AutoMigrate(models...).Error; nil != err {
		//	log.Errorf("auto migrate tables failed: %s", err.Error())
		//}
		//return

		mapID := connection.DBID(key)
		pool[mapID] = db
	}

	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return "mnt_" + defaultTableName
	//}

	connection.pool = pool
	return connection
}


/**
 * 对外暴露
 * @param
 * @return
 */
func Connection() *DB {
	return connection
}