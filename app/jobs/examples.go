package jobs

import (
	"fmt"
	"quanta-admin/app/business/daos"
	"quanta-admin/app/business/models"
	businessService "quanta-admin/app/business/service"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/instance_service"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/config"
	"quanta-admin/notification/lark"
	"strconv"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
)

// InitJob
// 需要将定义的struct 添加到字典中；
// 字典 key 可以配置到 自动任务 调用目标 中；
func InitJob() {
	jobList = map[string]JobExec{
		"ExamplesOne":                  ExamplesOne{},
		"InstanceInspection":           InstanceInspection{},
		"PriceTriggerInspection":       PriceTriggerInspection{},
		"PriceTriggerExpireInspection": PriceTriggerExpireInspection{},
		"DexCexObserverInspection":     DexCexObserverInspection{},
		// ...
	}
}

// ExamplesOne
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExamplesOne struct {
}

func (t ExamplesOne) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesOne exec success"
	// TODO: 这里需要注意 Examples 传入参数是 string 所以 arg.(string)；请根据对应的类型进行转化；
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			//testLarkNotification()
			fmt.Println("string", arg.(string))
			fmt.Println(str, arg.(string))
		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}

func testLarkNotification() {
	fmt.Println("测试lark通知")
	notification := lark.NewLarkTextNotification(new(string), "测试")
	notification.SendNotification()
}

// InstanceInspection 实例巡检，防止策略端服务重启后实例下线
type InstanceInspection struct{}

func (t InstanceInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore InstanceInspection exec success\r\n"
	// 获取所有instance对应的grpc server
	service := daos.BusStrategyInstanceDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}
	var registeredStrategyList []models.BusStrategyBaseInfo
	err := service.GetRegisteredInstanceList(&registeredStrategyList)
	if err != nil {
		log.Errorf("GetRegisteredInstanceList error: %v", err)
		return err
	}
	for _, strategy := range registeredStrategyList {
		grpcServiceName := config.ExtConfig.GetGrpcWithURL(strategy.GrpcEndpoint)
		instancesResp, err := client.ListInstance(grpcServiceName)
		if err != nil {
			log.Errorf("strategy: %s,  grpc service: %s get running instances error: %v\r\n", strategy.StrategyName, grpcServiceName, err.Error())
			continue
		}
		existIds := instancesResp.GetInstanceIds()
		log.Infof("strategy: %s, get running instances: %v\r\n", strategy.StrategyName, existIds)
		var instances []models.BusStrategyInstance
		err = service.GetRunningInstanceByStrategyId(strategy.Id, &instances)
		if err != nil {
			log.Errorf("GetRunningInstanceByStrategyId error: %v\r\n", err)
			continue
		}

		for _, instance := range instances {
			if contains(existIds, strconv.Itoa(instance.Id)) {
				continue
			}
			log.Infof("instance id : %d not exists in grpc service, restart\r\n", instance.Id)
			var config models.BusStrategyInstanceConfig
			err := service.GetInstanceConfigByInstanceId(instance.Id, &config)
			if err != nil {
				log.Errorf("GetInstanceConfigByInstanceId error: %v", err)
				continue
			}
			var instanceType instance_service.InstanceType
			//启动实例
			if instance.Type == "0" { // 观察者
				instanceType = instance_service.InstanceType_OBSERVER_INSTANCE
			} else if instance.Type == "1" { // 交易者
				instanceType = instance_service.InstanceType_TRADER_INSTANCE
			} else {
				log.Errorf("unsupport instance type :%s \r\n", instance.Type)
				continue
			}

			configStruct := config.SchemaText
			newInstance, err := client.StartNewInstance(grpcServiceName, strconv.Itoa(instance.Id), instanceType, &configStruct)
			if err != nil {
				log.Errorf("StartNewInstance error: %v", err)
				continue
			}
			log.Infof("start new instance : %s successful\r\n", newInstance)
		}

	}
	log.Infof(str)
	return nil
}

// PriceTriggerInspection 实例巡检，防止策略端服务重启后实例下线
type PriceTriggerInspection struct{}

func (t PriceTriggerInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore PriceTriggerInspection exec success\r\n"
	instanceIds, err := client.ListInstances()
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	log.Infof("instanceIds:%+v\n", instanceIds)
	service := daos.BusPriceTriggerInstanceDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}
	apiConfigService := daos.BusPriceTriggerApiConfigDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}

	instances := make([]models.BusPriceTriggerStrategyInstance, 0)
	err = service.GetInstancesList(&instances)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	for _, instance := range instances {
		if instance.Status == "started" && !contains(instanceIds, strconv.Itoa(instance.Id)) && instance.CloseTime.After(time.Now()) {
			// 中台状态为started，但是策略端没有，则需要重启
			// 但是目前需要做一个校验，重启前需要判断，是不是有止盈单，如果有，就可能是止盈暂停了。这个时候需要暂停中台的任务，而不是去重启

			apiConfig := models.BusPriceTriggerStrategyApikeyConfig{}
			err := apiConfigService.GetApiConfigById(instance.ApiConfig, &apiConfig)
			if err != nil {
				log.Errorf("重启 instance id: %d 失败, 异常信息：%v \r\n", instance.Id, err.Error())
				continue
			}

			var results []models.BusPriceMonitorForOptionHedging
			oneSecondAgo := time.Now().Add(-2 * time.Second)

			timestamp := oneSecondAgo.Unix()
			log.Infof("oneSecondAgo (Unix Timestamp - Seconds): %d", timestamp)

			err = service.GetStopProfitTradingRecord(instance.Id, oneSecondAgo, &results)
			if err != nil {
				log.Errorf("获取止盈记录失败 error:%s \r\n", err)
				continue
			}

			log.Infof("example jobs get trading for instance : %d, trading: %+v \n", instance.Id, results)

			if len(results) > 0 {
				// 如果有止盈单，修改实例状态
				err = service.UpdateInstancePaused(instance.Id)
				if err != nil {
					log.Errorf("update price trigger instance paused error:%s \r\n", err)
					continue
				}
				log.Infof("update price trigger instance,  instanceId: %d, status: %s \n", instance.Id, "paused")
				// 如果走到了这一步，则就直接跳过本次循环，不去重启
				continue
			}

			apiConfigReq := trigger_service.APIConfig{
				ApiKey:    apiConfig.ApiKey,
				SecretKey: apiConfig.SecretKey,
				Exchange:  apiConfig.Exchange,
			}

			profitTargetConfig := trigger_service.ProfitTargetConfig{
				InstanceId: strconv.Itoa(instance.Id),
			}
			profitTargetValue, err := strconv.ParseFloat(instance.ProfitTargetPrice, 64)
			if err != nil {
				fmt.Println("转换失败:", err)
				return err
			}

			if instance.ProfitTargetType == "LIMIT" {
				//限价止盈
				profitTargetConfig.ProfitTargetType = trigger_service.ProfitTargetType_LIMIT
				profitTargetConfig.Config = &trigger_service.ProfitTargetConfig_LimitConfig{
					LimitConfig: &trigger_service.LimitTypeConfig{
						ProfitTargetPrice: profitTargetValue,
						//LossTargetPrice:   req.LossTargetPrice,
					},
				}
			} else if instance.ProfitTargetType == "FLOATING" {
				//浮动止盈
				profitTargetConfig.ProfitTargetType = trigger_service.ProfitTargetType_FLOATING
				profitTargetConfig.Config = &trigger_service.ProfitTargetConfig_FloatingConfig{
					FloatingConfig: &trigger_service.FloatingTypeConfig{
						CallbackRatio: *instance.CallbackRatio,
						CutoffRatio:   *instance.CutoffRatio,
						MinProfit:     *instance.MinProfit,
					},
				}
			}
			execConfig := trigger_service.ExecuteConfig{
				InstanceId: strconv.Itoa(instance.Id),
				ExecuteNum: uint32(instance.ExecuteNum),
				DelayTime:  uint32(instance.DelayTime),
			}

			request := &trigger_service.StartTriggerRequest{
				InstanceId:         strconv.Itoa(instance.Id),
				OpenPrice:          instance.OpenPrice,
				ClosePrice:         instance.ClosePrice,
				Side:               instance.Side,
				Amount:             instance.Amount,
				Symbol:             instance.Symbol,
				StopTime:           strconv.FormatInt(instance.CloseTime.UnixMilli(), 10),
				ApiConfig:          &apiConfigReq,
				UserId:             instance.ExchangeUserId,
				ProfitTargetConfig: &profitTargetConfig,
				ExecuteConfig:      &execConfig,
				CloseOrderType:     instance.CloseOrderType,
			}
			log.Infof("start price trigger instance with profit target type :%+v\n", request.ProfitTargetConfig.ProfitTargetType)
			log.Infof("start price trigger instance request:%+v\n", request)
			_, err = client.StartTriggerInstance(request)
			if err != nil {
				log.Errorf("Service grpc start error:%s \r\n", err)
				continue
			}
		}
	}
	log.Infof(str)
	return nil
}

func contains(instanceIds []string, target string) bool {
	for _, id := range instanceIds {
		if id == target {
			return true
		}
	}
	return false
}

func containsObserver(observerInfos []*pb.BasicInfo, target string) (bool, bool) {
	for _, info := range observerInfos {
		if *info.InstanceId == target {
			return true, *info.TraderEnabled
		}
	}
	return false, false
}

func containsWaterLevelInstance(waterLevelInstances *waterLevelPb.InstanceListResponse, target string) bool {
	if waterLevelInstances == nil {
		return false
	}
	for _, id := range waterLevelInstances.InstanceIds {
		if id == target {
			return true
		}
	}
	return false
}

type PriceTriggerExpireInspection struct{}

func (t PriceTriggerExpireInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore PriceTriggerExpireInspection exec success\r\n"
	log.Infof("开始执行price-trigger 过期扫描任务")
	service := daos.BusPriceTriggerInstanceDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}

	instances := make([]models.BusPriceTriggerStrategyInstance, 0)
	err := service.GetInstancesList(&instances)
	if err != nil {
		log.Infof(err.Error())
		return err
	}
	expiredIds := make([]string, 0)
	for _, instance := range instances {
		if instance.Status == "started" && instance.CloseTime.Before(time.Now()) {
			//超过close time，自动关停
			expiredIds = append(expiredIds, strconv.Itoa(instance.Id))
		}
	}
	log.Infof("过期任务id：", expiredIds)
	if len(expiredIds) > 0 {
		err = service.ExpireInstanceWithIds(expiredIds)
		if err != nil {
			fmt.Printf("关停过期下单实例失败, 异常信息：%v\n", err.Error())
		}
	}
	log.Infof(str)
	return nil
}

type DexCexObserverInspection struct{}

func (t DexCexObserverInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore DexCexObserverInspection exec success \r\n"
	log.Infof("开始执行dex-cex 检查任务")
	db := sdk.Runtime.GetDbByKey("*")
	service := daos.BusDexCexTriangularObserverDAO{
		Db: db,
	}

	observerInfos, err := client.ListArbitragerClient()
	if err != nil {
		log.Errorf("grpc获取监控实例失败, %+v \n", err)
	}

	log.Infof("observerInfos:%+v\n", observerInfos)

	waterLevelInstances, err := client.ListWaterLevelInstance()
	log.Infof("waterLevelInstances:%+v\n", waterLevelInstances)

	observers := make([]models.BusDexCexTriangularObserver, 0)
	err = service.GetObserverList(&observers)
	if err != nil {
		log.Infof(err.Error())
		return err
	}

	for _, observer := range observers {
		log.Infof("observer:%+v\n", observer)
		if observer.Status == "4" {
			//已停止的直接跳过
			log.Infof("observer: %d\n status is stopped, skip \r\n", observer.Id)
			continue
		}

		if observer.Status == "1" {
			// 已开启observer的，只需要校验是否启动observer
			log.Infof("observer: %d\n status is 1, check observer \r\n", observer.Id)
			if exists, _ := containsObserver(observerInfos, strconv.Itoa(observer.Id)); !exists {
				// 服务端不存在，重启
				err = businessService.DoStartObserver(&observer)
				if err != nil {
					continue
				}
			}

		} else if observer.Status == "2" {
			// 水位调节中的，需要校验1. observer是否开启，水位调节是否开启
			log.Infof("observer: %d\n status is 2, check observer and water level \r\n", observer.Id)
			exist, isTrading := containsObserver(observerInfos, strconv.Itoa(observer.Id))
			if !exist {
				// 服务端不存在的，重启
				err = businessService.DoStartObserver(&observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}

			// 校验交易功能是否被风控
			result, _ := businessService.CheckIsTradeBlockedByRiskControl(observer.Id)
			if result {
				// 交易被风控，不进行下一步水位调节开启
				log.Infof("observer: %d\n is blocked by risk control, skip water level \r\n", observer.Id)
				continue
			}

			if !containsWaterLevelInstance(waterLevelInstances, strconv.Itoa(observer.Id)) {
				// 服务端不存在的，重启
				err = businessService.DoStartTokenWaterLevel(db, &observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}

			if !isTrading {
				// 如果实例水位调节中，按照讨论，不暂停交易的话，恢复现场就需要启动交易，但是服务端没有启动交易功能，则还需要重启交易功能
				err = businessService.DoStartTrader(db, &observer)
				if err != nil {
					//如果重启失败，则下次再重启
					continue
				}
			}

		} else if observer.Status == "3" {
			// 中台显示状态为3，启动交易中，但需要根据isTrading字段一起判断
			log.Infof("observer: %d\n status is: 3, is trading is: %t, check observer, water level, trader \r\n", observer.Id, observer.IsTrading)
			exist, isTrading := containsObserver(observerInfos, strconv.Itoa(observer.Id))
			if !exist {
				// 服务端不存在的，重启
				err = businessService.DoStartObserver(&observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}

			// 校验交易功能是否被风控
			result, _ := businessService.CheckIsTradeBlockedByRiskControl(observer.Id)
			if result {
				// 交易被风控，不进行下一步水位调节开启和交易开启
				log.Infof("observer: %d\n is blocked by risk control, skip water level and trading \r\n", observer.Id)
				continue
			}

			if !containsWaterLevelInstance(waterLevelInstances, strconv.Itoa(observer.Id)) {
				// 服务端不存在的，重启
				err = businessService.DoStartTokenWaterLevel(db, &observer)
				if err != nil {
					//如果重启失败，则下次再重启
					continue
				}
			}

			if observer.IsTrading && !isTrading {
				// 如果实例开启了交易，但是服务端没有启动交易功能，则还需要重启交易功能
				err = businessService.DoStartTrader(db, &observer)
				if err != nil {
					//如果重启失败，则下次再重启
					continue
				}
			}
		}

	}

	log.Infof(str)
	return nil
}
