package gormmapper

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func Ping() {
	l, _ := time.LoadLocation("PRC")
	o := cron.WithLocation(l)
	c := cron.New(o, cron.WithSeconds())

	// 定时检查
	c.AddFunc("@every 30m", func() {
		ping()
	})

	c.Start()
}

func ping() {
	log.Printf("[info] database ping ...")
	for key, value := range connections {
		sd, err := value.db.DB()
		if err != nil {
			log.Printf("[error] db#%v ping error1: %v", key, err.Error())
		}

		err = sd.Ping()
		if err != nil {
			log.Printf("[error] db#%v ping error2: %v", key, err.Error())
		}
	}
}
