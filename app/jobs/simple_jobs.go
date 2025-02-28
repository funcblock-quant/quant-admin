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
		fmt.Println("GetLatestSpreadData Job running")
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

	c.AddFunc("@every 5s", func() {
		// 每5s一次，获查看链上链下套利实例水位调节情况
		fmt.Println("Monitor WaterLevel Job running")
		s := service.BusDexCexTriangularObserver{}
		orm := sdk.Runtime.GetDbByKey("*")
		s.Orm = orm
		log := logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{})
		s.Log = log
		err := s.MonitorWaterLevelToStartTrader()
		if err != nil {
			fmt.Errorf("Monitor WaterLevel Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 5s", func() {
		// 每5s一次，获查看交易中的实例的水位健康状态
		fmt.Println("MonitorWaterLevelToStopTrader Job running")
		s := service.BusDexCexTriangularObserver{}
		orm := sdk.Runtime.GetDbByKey("*")
		s.Orm = orm
		log := logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{})
		s.Log = log
		err := s.MonitorWaterLevelToStopTrader()
		if err != nil {
			fmt.Errorf("MonitorWaterLevelToStopTrader Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 1m", func() {
		// 每5s一次，启动全局水位调整功能
		fmt.Println("StartGlobalWaterLevelConfig Job running")
		s := service.BusDexCexTriangularObserver{}
		orm := sdk.Runtime.GetDbByKey("*")
		s.Orm = orm
		log := logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{})
		s.Log = log
		err := s.StartGlobalWaterLevel()
		if err != nil {
			fmt.Errorf("StartGlobalWaterLevelConfig Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 1m", func() {
		// 每1m一次，对于交易进行简单封控和告警
		fmt.Println("Start Scan dex cex trade records Job running")
		s := service.StrategyDexCexTriangularArbitrageTrades{}
		orm := sdk.Runtime.GetDbByKey("*")
		s.Orm = orm
		log := logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{})
		s.Log = log
		err := s.ScanTrades()
		if err != nil {
			fmt.Errorf("StartGlobalWaterLevelConfig Job run failed, err:%v\n", err)
		}
	})

	c.Start()
	select {}
}
