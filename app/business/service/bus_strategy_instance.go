package service

import (
	"errors"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"strconv"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/instance_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	ext "quanta-admin/config"
)

type BusStrategyInstance struct {
	service.Service
}

// GetPage 获取BusStrategyInstance列表
func (e *BusStrategyInstance) GetPage(c *dto.BusStrategyInstanceGetPageReq, p *actions.DataPermission, list *[]dto.BusStrategyInstanceGetPageResp, count *int64) error {
	var err error
	var data models.BusStrategyInstance

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Joins("LEFT JOIN bus_strategy_base_info ON bus_strategy_base_info.id = bus_strategy_instance.strategy_id"). // 关联策略表
		Select("bus_strategy_instance.*, bus_strategy_base_info.strategy_name").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyInstance对象
func (e *BusStrategyInstance) Get(d *dto.BusStrategyInstanceGetReq, p *actions.DataPermission, model *dto.BusStrategyInstanceGetResp) error {
	var data models.BusStrategyInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyInstance error:%s \r\n", err)
		return err
	}

	var configs []models.BusStrategyInstanceConfig
	err = e.Orm.Model(&models.BusStrategyInstanceConfig{}).
		Where("strategy_instance_id=?", d.Id).
		Find(&configs).Error

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if len(configs) > 0 {
		model.Schema = configs[0]
	}

	return nil
}

// Insert 创建BusStrategyInstance对象
func (e *BusStrategyInstance) Insert(c *dto.BusStrategyInstanceInsertReq) error {
	var err error
	var data models.BusStrategyInstance
	var instanceConfig = models.BusStrategyInstanceConfig{}

	// 启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	c.Generate(&data)
	if err = tx.Create(&data).Error; err != nil {
		tx.Rollback() // 插入主表失败，回滚事务
		e.Log.Errorf("Error while inserting BusStrategyInstance: %v", err)
		return err
	}

	// 2. 保存配置表数据，每个配置项需要关联主表的 ID
	c.Schema.StrategyInstanceId = strconv.Itoa(data.Id)
	c.Schema.Generate(&instanceConfig)
	e.Log.Infof("instanceConfig: %+v", instanceConfig)
	if err := tx.Create(&instanceConfig).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Error while inserting instance config: %v", err)
		return err
	}

	// 3. 提交事务
	err = tx.Commit().Error
	return err
}

// Update 修改BusStrategyInstance对象
func (e *BusStrategyInstance) Update(c *dto.BusStrategyInstanceUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyInstance{}
	e.Log.Infof("c:%+v ", *c)
	// 启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	tx.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)
	e.Log.Infof("data:%+v ", data)

	db := tx.Save(&data)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("BusStrategyInstanceService update error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 2. 删除旧的配置数据
	log.Infof("strategy_instance_id: %d", data.Id)
	if err := tx.Where("strategy_instance_id = ?", data.Id).Delete(&models.BusStrategyInstanceConfig{}).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Error deleting old configurations: %v", err)
		return err
	}

	var instanceConfig = models.BusStrategyInstanceConfig{}
	// 3. 保存配置表数据，每个配置项需要关联主表的 ID

	c.Schema.StrategyInstanceId = strconv.Itoa(data.Id)
	c.Schema.Generate(&instanceConfig)
	// 将配置数据插入配置表
	e.Log.Infof("instanceConfig: %+v", instanceConfig)
	if err := tx.Create(&instanceConfig).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Error while update instance config: %v", err)
		return err
	}

	// 4. 提交事务
	err = tx.Commit().Error
	if err != nil {
		e.Log.Errorf("Transaction commit failed: %v", err)
		return err
	}

	return nil
}

// StartInstance 启动BusStrategyInstance
func (e *BusStrategyInstance) StartInstance(d *dto.BusStrategyInstanceStartReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstance
	// 查询策略实例
	err := e.Orm.Model(&data).First(&data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("can not find bus strategy_instance_id:%d \r\n", d.Id)
		return err
	}

	strategy := models.BusStrategyBaseInfo{}
	// 获取策略基本信息
	err = e.Orm.Model(&strategy).First(&strategy, "id=?", data.StrategyId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("can not find bus strategy_id:%s \r\n", data.StrategyId)
		return err
	}

	//获取策略实例配置
	instanceConfig := models.BusStrategyInstanceConfig{}
	err = e.Orm.Model(&instanceConfig).First(&instanceConfig, "strategy_instance_id=?", data.Id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("can not find instance config :%s \r\n", data.Id)
		return err
	}

	endpoint := strategy.GrpcEndpoint
	config := ext.Extend{}
	serviceName := config.GetGrpcWithURL(endpoint)
	if serviceName == "" {
		e.Log.Errorf("can not find bus grpc service:%s \r\n", endpoint)
		return errors.New("can not find grpc service info")
	}

	var instanceType instance_service.InstanceType
	//启动实例
	if data.Type == "0" { // 观察者
		instanceType = instance_service.InstanceType_OBSERVER_INSTANCE
	} else if data.Type == "1" { // 交易者
		instanceType = instance_service.InstanceType_TRADER_INSTANCE
	} else {
		e.Log.Errorf("unsupport instance type :%s \r\n", data.Type)
		return errors.New("unsupported instance type")
	}

	configStruct := instanceConfig.SchemaText

	_, err = client.StartNewInstance(serviceName, strconv.Itoa(data.Id), instanceType, &configStruct)
	if err != nil {
		e.Log.Errorf("start instance occur grpc error :%v \r\n", err)
		return errors.New("start instance failed")
	}

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where("id = ?", d.GetId()).
		Updates(map[string]interface{}{
			"start_run_time": time.Now(),
			"stop_run_time":  nil,
			"status":         1,
		})
	if err := db.Error; err != nil {
		//如果保存数据库失败了，则需要尝试停止实例
		client.StopInstance(serviceName, strconv.Itoa(data.Id))
		e.Log.Errorf("Service BusStrategyInstance start error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权修改该数据")
	}
	return nil
}

func (e *BusStrategyInstance) BatchStartInstance(d *dto.BusStrategyInstanceBatchStartReq, p *actions.DataPermission, failedCount *int) error {
	ids := d.GetId()
	*failedCount = 0
	if idsList, ok := ids.([]int); ok {
		for _, id := range idsList {
			req := &dto.BusStrategyInstanceStartReq{
				Id: id,
			}
			err := e.StartInstance(req, p)
			if err != nil {
				e.Log.Warnf("batch start instance: %d error:%s \r\n", id, err)
				*failedCount++
				continue
			}
		}
	} else {
		e.Log.Errorf("batch start instance error with type assert \r\n")
		return errors.New("batch start instance error")
	}

	return nil
}

// StopInstance 暂停BusStrategyInstance
func (e *BusStrategyInstance) StopInstance(d *dto.BusStrategyInstanceStopReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstance
	// 查询策略实例
	err := e.Orm.Model(&data).First(&data, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("can not find bus strategy_instance_id:%d \r\n", d.Id)
		return err
	}

	strategy := models.BusStrategyBaseInfo{}
	// 获取策略基本信息
	err = e.Orm.Model(&strategy).First(&strategy, "id=?", data.StrategyId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		e.Log.Errorf("can not find bus strategy_id:%s \r\n", data.StrategyId)
		return err
	}
	endpoint := strategy.GrpcEndpoint
	config := ext.Extend{}
	serviceName := config.GetGrpcWithURL(endpoint)
	if serviceName == "" {
		e.Log.Errorf("can not find bus strategy_id:%s \r\n", data.StrategyId)
		return errors.New("can not find grpc service info")
	}
	err = client.StopInstance(serviceName, strconv.Itoa(data.Id))
	if err != nil {
		e.Log.Errorf("stop instance failed , grpc error:%v \r\n", err)
		return errors.New("stop instance failed")
	}

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where("id = ?", d.GetId()).
		Updates(map[string]interface{}{
			"stop_run_time": time.Now(),
			"status":        0,
		})
	if err := db.Error; err != nil {
		e.Log.Errorf("Service BusStrategyInstance stop error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权修改该数据")
	}
	return nil
}

// Remove 删除BusStrategyInstance
func (e *BusStrategyInstance) Remove(d *dto.BusStrategyInstanceDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstance

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyInstance error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// QueryInstanceDashboard 查询策略实例dashboard数据
func (e *BusStrategyInstance) QueryInstanceDashboard(d *dto.BusStrategyInstanceDashboardGetReq, p *actions.DataPermission, resp *dto.BusStrategyInstanceDashboardGetResp) error {
	var data models.BusStrategyInstance
	var arbitrageData models.BusArbitrageRecord
	// 查询所有策略
	query := e.Orm.Model(&arbitrageData).
		Scopes(actions.Permission(arbitrageData.TableName(), p))
	//TODO dashboard面板
	if d.StrategyInstanceIds != nil && len(d.StrategyInstanceIds) > 0 {
		//1. 获取所有账户组
		var strategyInstances []models.BusStrategyInstance
		instanceQuery := e.Orm.Model(&data).Scopes(actions.Permission(data.TableName(), p))
		if err := instanceQuery.Where("id in ?", d.StrategyInstanceIds).
			Find(&strategyInstances).Error; err != nil {
			e.Log.Errorf("Service QueryInstanceDashboard error:%s", err)
		}

		query = query.Where("strategy_instance_id in ?", d.StrategyInstanceIds)

		// 查询策略数据并映射到结构
		if err := query.Select("id, name, status"). // 只选择必要的字段
								Find(&strategyInstances).Error; err != nil {
			e.Log.Errorf("QueryInstanceDashboard Find error:%s", err)
			return err
		}

		// 统计总数
		var totalCount int64
		if err := query.Count(&totalCount).Error; err != nil {
			e.Log.Errorf("QueryInstanceDashboard Count error:%s", err)
			return err
		}

	} else {
		// 查询指定的策略数据

	}
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service QueryInstanceDashboard error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// ParseInstanceConfig 将实例配置转成struct
func (e *BusStrategyInstance) ParseInstanceConfig(config models.BusStrategyInstanceConfig) (*structpb.Struct, error) {
	yamlString := config.SchemaText
	e.Log.Infof("yamlString : %v", yamlString)
	var result map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlString), &result)
	e.Log.Infof("result map : %v", result)
	if err != nil {
		e.Log.Errorf("Failed to parse YAML: %v", err)
		return nil, err
	}
	protoStruct, err := structpb.NewStruct(result)
	if err != nil {
		e.Log.Errorf("Failed to convert to Struct: %v", err)
		return nil, err
	}
	return protoStruct, nil
}
