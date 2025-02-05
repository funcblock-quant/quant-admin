package jobs

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/robfig/cron/v3"
	"quanta-admin/app/business/service"
)

func InitSimpleJob() {
	c := cron.New()
	fmt.Printf("Init simple job \n")
	c.AddFunc("@every 5s", func() {
		// 每5s一次，获取最新的实时价差
		s := service.BusDexCexPriceSpreadData{}
		orm := sdk.Runtime.GetDbByKey("*")
		s.Orm = orm
		log := logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{})
		s.Log = log
		err := s.GetLatestSpreadData()
		if err != nil {
			fmt.Errorf("GetLatestSpreadData failed, err:%v\n", err)
		}
	})

	c.Start()
	select {}
}
