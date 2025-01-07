package service

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusStrategyBaseInfo struct {
	service.Service
}

// GetPage 获取BusStrategyBaseInfo列表
func (e *BusStrategyBaseInfo) GetPage(c *dto.BusStrategyBaseInfoGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyBaseInfo, count *int64) error {
	var err error
	var data models.BusStrategyBaseInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyBaseInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Get(d *dto.BusStrategyBaseInfoGetReq, p *actions.DataPermission, model *models.BusStrategyBaseInfo) error {
	var data models.BusStrategyBaseInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyBaseInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Insert(c *dto.BusStrategyBaseInfoInsertReq) error {
	var err error
	var data models.BusStrategyBaseInfo
	var configs = make([]models.BusStrategyConfigDict, 0, len(c.Configurations))
	// 启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. 生成主表数据，并插入主表
	c.Status = strconv.Itoa(1) // 新注册策略默认为“已注册”
	c.Generate(&data)
	if err = tx.Create(&data).Error; err != nil {
		tx.Rollback() // 插入主表失败，回滚事务
		e.Log.Errorf("Error while inserting BusStrategyBaseInfo: %v", err)
		return err
	}

	// 2. 保存配置表数据，每个配置项需要关联主表的 ID
	for _, configReq := range c.Configurations {
		var config models.BusStrategyConfigDict
		configReq.Generate(&config)

		// 为每个配置设置关联主表的ID
		config.StrategyId = strconv.Itoa(data.Id) // 关联主表ID

		configs = append(configs, config)
	}
	// 将配置数据插入配置表
	if err := tx.CreateInBatches(&configs, len(configs)).Error; err != nil {
		tx.Rollback() // 插入配置失败，回滚事务
		e.Log.Errorf("Error while inserting config: %v", err)
		return err
	}

	// 3. 提交事务
	err = tx.Commit().Error
	return err
}

// Update 修改BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Update(c *dto.BusStrategyBaseInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyBaseInfo{}
	var configs = make([]models.BusStrategyConfigDict, 0, len(c.Configurations))

	// 启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 查询当前记录
	err = tx.First(&data, c.Id).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("Error finding BusStrategyBaseInfo: %v", err)
		return err
	}

	// 更新主表字段
	c.Generate(&data)
	if err := tx.Save(&data).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Error updating BusStrategyBaseInfo: %v", err)
		return err
	}

	// 2. 删除旧的配置数据
	log.Infof("strategy_id: %d", data.Id)
	if err := tx.Where("strategy_id = ?", data.Id).Delete(&models.BusStrategyConfigDict{}).Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("Error deleting old configurations: %v", err)
		return err
	}

	// 3. 插入新的配置数据
	for _, configReq := range c.Configurations {
		var config models.BusStrategyConfigDict
		configReq.Generate(&config)

		// 关联主表 ID
		config.StrategyId = strconv.Itoa(data.Id)

		// 插入配置数据
		if err := tx.Create(&config).Error; err != nil {
			tx.Rollback()
			e.Log.Errorf("Error inserting new config: %v", err)
			return err
		}
		configs = append(configs, config)
	}

	// 4. 提交事务
	err = tx.Commit().Error
	if err != nil {
		e.Log.Errorf("Transaction commit failed: %v", err)
		return err
	}

	return nil
}

// Remove 删除BusStrategyBaseInfo
func (e *BusStrategyBaseInfo) Remove(d *dto.BusStrategyBaseInfoDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyBaseInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyBaseInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
