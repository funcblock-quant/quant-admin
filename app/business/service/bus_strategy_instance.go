package service

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
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
		).Joins("LEFT JOIN bus_strategy_base_info ON bus_strategy_base_info.id = bus_strategy_instance.strategy_id").            // 关联策略表
		Joins("LEFT JOIN bus_exchange_account_group ON bus_exchange_account_group.id = bus_strategy_instance.account_group_id"). // 关联group表
		Select("bus_strategy_instance.*, bus_exchange_account_group.group_name, bus_strategy_base_info.strategy_name").
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
	e.Log.Infof("configs %+v, 配置数量 %d \n", configs, len(configs))
	model.Configs = configs

	return nil
}

// Insert 创建BusStrategyInstance对象
func (e *BusStrategyInstance) Insert(c *dto.BusStrategyInstanceInsertReq) error {
	var err error
	var data models.BusStrategyInstance
	var instanceConfigs = make([]models.BusStrategyInstanceConfig, 0, len(c.Configurations))
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
	for _, configReq := range c.Configurations {
		var config models.BusStrategyInstanceConfig
		configReq.Generate(&config)

		// 为每个配置设置关联主表的ID
		config.StrategyInstanceId = strconv.Itoa(data.Id) // 关联主表ID
		instanceConfigs = append(instanceConfigs, config)
	}

	// 将配置数据插入配置表
	if err = tx.CreateInBatches(&instanceConfigs, len(instanceConfigs)).Error; err != nil {
		tx.Rollback() // 插入配置失败，回滚事务
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
		e.Log.Errorf("BusStrategyInstanceService Save error:%s \r\n", err)
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

	var instanceConfigs = make([]models.BusStrategyInstanceConfig, 0, len(c.Configurations))
	// 2. 保存配置表数据，每个配置项需要关联主表的 ID
	for _, configReq := range c.Configurations {
		var config models.BusStrategyInstanceConfig
		configReq.Generate(&config)

		// 为每个配置设置关联主表的ID
		config.StrategyInstanceId = strconv.Itoa(data.Id) // 关联主表ID
		instanceConfigs = append(instanceConfigs, config)
	}

	// 将配置数据插入配置表
	if err = tx.CreateInBatches(&instanceConfigs, len(instanceConfigs)).Error; err != nil {
		tx.Rollback() // 插入配置失败，回滚事务
		e.Log.Errorf("Error while inserting instance config: %v", err)
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

	//TODO 后期这里需要真实去启动实例
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
		e.Log.Errorf("Service BusStrategyInstance start error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权修改该数据")
	}
	return nil
}

// StopInstance 暂停BusStrategyInstance
func (e *BusStrategyInstance) StopInstance(d *dto.BusStrategyInstanceStopReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstance

	//TODO 后期这里需要真实去停用实例
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
	if d.StrategyInstanceIds != nil && len(d.StrategyInstanceIds) > 0 {
		//1. 获取所有账户组
		var strategyInstances []models.BusStrategyInstance
		instanceQuery := e.Orm.Model(&data).Scopes(actions.Permission(data.TableName(), p))
		if err := instanceQuery.Where("id in ?", d.StrategyInstanceIds).
			Find(&strategyInstances).Error; err != nil {
			e.Log.Errorf("Service QueryInstanceDashboard error:%s", err)
		}

		//totalBeginBalance := 0
		//for _, strategyInstance := range strategyInstances {
		//
		//}

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
