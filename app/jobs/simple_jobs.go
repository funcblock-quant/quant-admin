package jobs

import (
	"fmt"
	"quanta-admin/app/business/service"
	"sync"
	"time"

	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/robfig/cron/v3"
)

var (
	jobLock sync.Mutex // 防止任务重入
)

func InitSimpleJob() {
	c := cron.New()
	fmt.Printf("Init simple job \n")
	orm := sdk.Runtime.GetDbByKey("*")               // 复用数据库连接
	log := logger.NewHelper(sdk.Runtime.GetLogger()) // 复用日志对象
	c.AddFunc("@every 5s", func() {

		// 每5s一次，获取最新的实时价差
		fmt.Println("GetLatestSpreadData Job running")
		s := service.BusDexCexPriceSpreadData{}
		s.Orm = orm
		s.Log = log
		err := s.GetLatestSpreadData()
		if err != nil {
			log.Errorf("GetLatestSpreadData failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 5s", func() {
		// 每5s一次，获查看链上链下套利实例水位调节情况
		fmt.Println("Monitor WaterLevel Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.MonitorWaterLevelToStartTrader()
		if err != nil {
			log.Errorf("Monitor WaterLevel Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 5s", func() {
		// 每5s一次，获查看交易中的实例的水位健康状态
		fmt.Println("MonitorWaterLevelToStopTrader Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.MonitorWaterLevelToStopTrader()
		if err != nil {
			log.Errorf("MonitorWaterLevelToStopTrader Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 1m", func() {
		// 每5s一次，启动全局水位调整功能
		fmt.Println("StartGlobalWaterLevelConfig Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.StartGlobalWaterLevel()
		if err != nil {
			log.Errorf("StartGlobalWaterLevelConfig Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 1m", func() {
		// 每1m一次，对于交易进行简单封控和告警
		fmt.Println("Start Scan dex cex trade records Job running")
		s := service.StrategyDexCexTriangularArbitrageTrades{}
		s.Orm = orm
		s.Log = log
		err := s.ScanTrades()
		if err != nil {
			log.Errorf("StartGlobalWaterLevelConfig Job run failed, err:%v\n", err)
		}
	})

	c.AddFunc("@every 5s", func() {
		// 每5s一次，查看期权下单策略的执行次数是否到期
		fmt.Println("Monitor Price Trigger ExecuteNum Job running")
		s := service.BusPriceTriggerStrategyInstance{}
		s.Orm = orm
		s.Log = log
		err := s.MonitorExecuteNum()
		if err != nil {
			log.Errorf("Monitor Price Trigger ExecuteNum Job run failed, err:%v\n", err)
		}

	})

	c.AddFunc("@every 5s", func() {

		// 每5s一次风控检查
		fmt.Println("CheckRisk Control Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.CheckRiskControl()
		if err != nil {
			log.Errorf("CheckRiskControl Job run failed, err:%v\n", err)
		}

	})

	c.AddFunc("@every 5s", func() {

		// 每5s一次检查存在未处理的暂停交易的风控事件，并检查对应的交易实例是否暂停，作为补偿
		fmt.Println("Check Exist Risk Event Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.CheckExistRiskEvent()
		if err != nil {
			log.Errorf("Check Exist Risk Event Job run failed, err:%v\n", err)
		}

	})

	c.AddFunc("@every 5s", func() {

		// 每5s一次检查被blocking的交易实例，如果所有风控都被恢复，也需要重新启动交易
		fmt.Println("Check Blocking Instance Job running")
		s := service.BusDexCexTriangularObserver{}
		s.Orm = orm
		s.Log = log
		err := s.CheckBlockingInstance()
		if err != nil {
			log.Errorf("Check Exist Risk Event Job run failed, err:%v\n", err)
		}

	})

	c.Start()
	select {}
}

// 包装任务，防止崩溃
func safeWrapper(name string, job func()) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Job %s crashed: %v", name, r)
				jobLock.Unlock()
			}
		}()
		start := time.Now()
		job()
		logger.Infof("Job %s finished in %v", name, time.Since(start))
	}
}
