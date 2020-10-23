package gormmapper

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func Ping() {
	loc, _  := time.LoadLocation("PRC")
	option := cron.WithLocation(loc)
	dbcron  := cron.New(option, cron.WithSeconds())

	// 检查心跳
	dbcron.AddFunc("@every 60m", func() {
		ping()
	})

	dbcron.Start()
}

func ping() {
	log.Printf("[info] database ping ...")
	for key, value := range connection.pool  {
		log.Printf("key: %v, value: %v", key, value)
		//err := value.DB().Ping()
		//if err != nil {
		//	log.Printf("[error] db#%v ping error: %v", key, err.Error())
		//}
	}
}