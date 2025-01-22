package jobs

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"quanta-admin/app/business/daos"
	"quanta-admin/app/business/models"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/trigger_service"
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
		"MonitorArbitrageOpportunity":  MonitorArbitrageOpportunity{}, //监控套利机会定时任务
		"InstanceInspection":           InstanceInspection{},
		"PriceTriggerInspection":       PriceTriggerInspection{},
		"PriceTriggerExpireInspection": PriceTriggerExpireInspection{},
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

// MonitorArbitrageOpportunity
type MonitorArbitrageOpportunity struct {
}

func (t MonitorArbitrageOpportunity) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore MonitorArbitrageOpportunity exec success"

	fmt.Println(str)
	return nil
}

// InstanceInspection 实例巡检，防止策略端服务重启后实例下线
type InstanceInspection struct{}

func (t InstanceInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore InstanceInspection exec success"
	instance, err := client.ListInstance("market-making")
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	fmt.Printf("instance:%+v\n", instance.InstanceIds)
	fmt.Printf(str)
	return nil
}

// PriceTriggerInspection 实例巡检，防止策略端服务重启后实例下线
type PriceTriggerInspection struct{}

func (t PriceTriggerInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore PriceTriggerInspection exec success"
	instanceIds, err := client.ListInstances()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	fmt.Printf("instanceIds:%+v\n", instanceIds)
	service := daos.BusPriceTriggerInstanceDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}
	apiConfigService := daos.BusPriceTriggerApiConfigDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}

	instances := make([]models.BusPriceTriggerStrategyInstance, 0)
	err = service.GetInstancesList(&instances)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	for _, instance := range instances {
		if instance.Status == "started" && !contains(instanceIds, strconv.Itoa(instance.Id)) && instance.CloseTime.After(time.Now()) {
			//中台状态为started，但是策略端没有，则需要重启
			apiConfig := models.BusPriceTriggerStrategyApikeyConfig{}
			err := apiConfigService.GetApiConfigById(instance.ApiConfig, &apiConfig)
			if err != nil {
				fmt.Printf("重启 instance id: %d 失败, 异常信息：%v", instance.Id, err.Error())
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
			}

			_, err = client.StartInstance(request)
			if err != nil {
				fmt.Errorf("Service grpc start error:%s \r\n", err)
				continue
			}
		}
	}
	fmt.Printf(str)
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

type PriceTriggerExpireInspection struct{}

func (t PriceTriggerExpireInspection) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore PriceTriggerExpireInspection exec success"
	fmt.Println("开始执行price-trigger 过期扫描任务")
	service := daos.BusPriceTriggerInstanceDAO{
		Db: sdk.Runtime.GetDbByKey("*"),
	}

	instances := make([]models.BusPriceTriggerStrategyInstance, 0)
	err := service.GetInstancesList(&instances)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	expiredIds := make([]string, 0)
	for _, instance := range instances {
		if instance.Status == "started" && instance.CloseTime.Before(time.Now()) {
			//超过close time，自动关停
			expiredIds = append(expiredIds, strconv.Itoa(instance.Id))
		}
	}
	fmt.Println("过期任务id：", expiredIds)
	err = service.ExpireInstanceWithIds(expiredIds)
	if err != nil {
		fmt.Printf("关停过期下单实例失败, 异常信息：%v\n", err.Error())
	}

	fmt.Printf(str)
	return nil
}
