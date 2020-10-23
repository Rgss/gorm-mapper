package main

import "github.com/Rgss/gorm-mapper"

func main() {

	initDB()

	te := &TE{}
	sb := gormmapper.NSearchBuilder()

	te.SelectOneBySearchBuilder()

}

type TE struct {
	gormmapper.Mapper
}

func initDB() {
	config := &gormmapper.DBConfig{
		User:         "root",
		Pass:         "123456",
		Host:         "127.0.0.1",
		Port:         3306,
		DbName:       "test",
		Charset:      "utf8",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		EnableLog:    false,
	}

	configs := make(map[string]*gormmapper.DBConfig)
	configs["default"] = config
	gormmapper.Initialize(configs)
}