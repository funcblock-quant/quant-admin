package jobs

import (
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"google.golang.org/protobuf/proto"
	"quanta-admin/app/business/daos"
	"quanta-admin/app/business/models"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/instance_service"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/config"
	"quanta-admin/notification/lark"
	"strconv"
	"time"
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
				log.Infof("instance id : %d exists in grpc service, skip restart\r\n", instance.Id)
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
	fmt.Printf(str)
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
			//中台状态为started，但是策略端没有，则需要重启
			apiConfig := models.BusPriceTriggerStrategyApikeyConfig{}
			err := apiConfigService.GetApiConfigById(instance.ApiConfig, &apiConfig)
			if err != nil {
				log.Errorf("重启 instance id: %d 失败, 异常信息：%v \r\n", instance.Id, err.Error())
				continue
			}

			apiConfigReq := trigger_service.APIConfig{
				ApiKey:    apiConfig.ApiKey,
				SecretKey: apiConfig.SecretKey,
				Exchange:  apiConfig.Exchange,
			}
			request := &trigger_service.StartTriggerRequest{
				InstanceId: strconv.Itoa(instance.Id),
				OpenPrice:  instance.OpenPrice,
				ClosePrice: instance.ClosePrice,
				Side:       instance.Side,
				Amount:     instance.Amount,
				Symbol:     instance.Symbol,
				StopTime:   strconv.FormatInt(instance.CloseTime.UnixMilli(), 10),
				ApiConfig:  &apiConfigReq,
				UserId:     instance.ExchangeUserId,
			}

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
	service := daos.BusDexCexTriangularObserverDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}

	observerInfos, err := client.ListArbitragerClient()
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
				err := restartObserver(&observer)
				if err != nil {
					continue
				}
			}

		} else if observer.Status == "2" {
			// 水位调节中的，需要校验1. observer是否开启，水位调节是否开启
			log.Infof("observer: %d\n status is 2, check observer and water level \r\n", observer.Id)
			if exist, _ := containsObserver(observerInfos, strconv.Itoa(observer.Id)); !exist {
				// 服务端不存在的，重启
				err := restartObserver(&observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}
			if !containsWaterLevelInstance(waterLevelInstances, strconv.Itoa(observer.Id)) {
				// 服务端不存在的，重启
				err := restartWaterLevel(&observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}
		} else if observer.Status == "3" {
			// 中台显示状态为3，启动交易中，但需要根据isTrading字段一起判断
			log.Infof("observer: %d\n status is: 3, is trading is: %t, check observer, water level, trader \r\n", observer.Id, observer.IsTrading)
			exist, isTrading := containsObserver(observerInfos, strconv.Itoa(observer.Id))
			if !exist {
				// 服务端不存在的，重启
				err := restartObserver(&observer)
				if err != nil {
					//如果重启失败，则不进行下一步水位调节开启
					continue
				}
			}

			if !containsWaterLevelInstance(waterLevelInstances, strconv.Itoa(observer.Id)) {
				// 服务端不存在的，重启
				err := restartWaterLevel(&observer)
				if err != nil {
					//如果重启失败，则下次再重启
					continue
				}
			}

			if observer.IsTrading && !isTrading {
				// 如果实例开启了交易，但是服务端没有启动交易功能，则还需要重启交易功能
				err := restartTrader(&observer)
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

func GenerateAmberConfig(observer *models.BusDexCexTriangularObserver, amberConfig *pb.AmberObserverConfig) error {
	amberConfig.ExchangeType = &observer.ExchangeType
	amberConfig.TakerFee = proto.Float64(*observer.TakerFee)

	amberConfig.TargetToken = &observer.TargetToken
	amberConfig.QuoteToken = &observer.QuoteToken

	if observer.Depth != "" {
		depthInt, err := strconv.Atoi(observer.Depth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		amberConfig.BidDepth = proto.Int32(int32(depthInt))
		amberConfig.AskDepth = proto.Int32(int32(depthInt))
	}
	return nil
}

func restartObserver(observer *models.BusDexCexTriangularObserver) error {
	slippageBpsUint, err := strconv.ParseUint(observer.SlippageBps, 10, 32)
	if err != nil {
		log.Infof("slippageBps: %v\n", slippageBpsUint)
		return err
	}

	maxArraySize := new(uint32)
	*maxArraySize = uint32(observer.MaxArraySize) //默认5， clmm使用参数

	dexConfig := &pb.DexConfig{}
	if observer.DexType == "RAY_AMM" {
		dexConfig.Config = &pb.DexConfig_RayAmm{
			RayAmm: &pb.RayAmmConfig{
				Pool:      observer.AmmPoolId,
				TokenMint: observer.TokenMint,
			},
		}
	} else if observer.DexType == "RAY_CLMM" {
		dexConfig.Config = &pb.DexConfig_RayClmm{
			RayClmm: &pb.RayClmmConfig{
				Pool:         observer.AmmPoolId,
				TokenMint:    observer.TokenMint,
				MaxArraySize: maxArraySize,
			},
		}
	}

	transactionFee := &pb.TransactionFee{
		PriorityFee: &observer.PriorityFee,
		JitoFee:     &observer.JitoFee,
	}

	arbitrageConfig := &pb.ObserverParams{
		MinSolAmount:             observer.MinSolAmount,
		MaxSolAmount:             observer.MaxSolAmount,
		TriggerProfitQuoteAmount: observer.MinProfit,
		TxFee:                    transactionFee,
	}

	amberConfig := &pb.AmberObserverConfig{}
	GenerateAmberConfig(observer, amberConfig)

	instanceId := strconv.Itoa(observer.Id)
	log.Infof("restart observer success with params: dexConfig: %+v\n, arbitrageConfig: %+v\n", dexConfig, arbitrageConfig)
	err = client.StartNewArbitragerClient(&instanceId, amberConfig, dexConfig, arbitrageConfig)
	if err != nil {
		log.Infof("restart observer throw grpc error: %v\n", err)
		return err
	}
	return nil
}

func restartTrader(observer *models.BusDexCexTriangularObserver) error {
	instanceId := strconv.Itoa(observer.Id)
	slippageBpsUint, err := strconv.ParseUint(observer.SlippageBps, 10, 32)
	if err != nil {
		log.Infof("slippageBps: %v\n", slippageBpsUint)
		return err
	}

	trader, err := observer.GetExchangeTypeForTrader()
	if err != nil {
		log.Errorf("get exchange type for trader failed: %v\n", err)
		return err
	}
	amberTraderConfig := &pb.AmberTraderConfig{
		ExchangeType: &trader,
	}
	traderParams := &pb.TraderParams{
		SlippageBps: &slippageBpsUint,
	}
	err = client.EnableTrader(instanceId, amberTraderConfig, traderParams)
	if err != nil {
		log.Errorf("restart instance: %d trader error: %v\n", instanceId, err)
		return err
	} else {
		log.Infof("restart instance: %d trader success\n", instanceId)
		return nil
	}
}

func restartWaterLevel(observer *models.BusDexCexTriangularObserver) error {
	tokenConfig := &waterLevelPb.TokenConfig{
		Currency:               observer.TargetToken,
		PubKey:                 *observer.TokenMint,
		OwnerProgram:           *observer.OwnerProgram,
		Decimals:               uint32(observer.Decimals),
		AlertThreshold:         strconv.FormatFloat(*observer.AlertThreshold, 'f', -1, 64),
		BuyTriggerThreshold:    strconv.FormatFloat(*observer.BuyTriggerThreshold, 'f', -1, 64),
		TargetBalanceThreshold: strconv.FormatFloat(*observer.TargetBalanceThreshold, 'f', -1, 64),
		SellTriggerThreshold:   strconv.FormatFloat(*observer.SellTriggerThreshold, 'f', -1, 64),
	}

	clientRequest := &waterLevelPb.StartInstanceRequest{
		InstanceId:   strconv.Itoa(observer.Id),
		ExchangeType: observer.ExchangeType,
		TokenConfig:  tokenConfig,
	}

	log.Infof("restart water level with req: %v \r\n", clientRequest)
	_, err := client.StartWaterLevelInstance(clientRequest)
	if err != nil {
		log.Errorf("启动水位调节失败:%s \r\n", err)
		return err
	}
	log.Infof("水位调节启动成功")
	return nil
}
