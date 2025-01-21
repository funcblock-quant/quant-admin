package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusPriceTriggerStrategyApikeyConfig struct {
	service.Service
}

// GetPage 获取BusPriceTriggerStrategyApikeyConfig列表
func (e *BusPriceTriggerStrategyApikeyConfig) GetPage(c *dto.BusPriceTriggerStrategyApikeyConfigGetPageReq, p *actions.DataPermission, list *[]models.BusPriceTriggerStrategyApikeyConfig, count *int64, userId *int) error {
	var err error
	var data models.BusPriceTriggerStrategyApikeyConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Where("user_id = ?", userId).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyApikeyConfigService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusPriceTriggerStrategyApikeyConfig对象
func (e *BusPriceTriggerStrategyApikeyConfig) Get(d *dto.BusPriceTriggerStrategyApikeyConfigGetReq, p *actions.DataPermission, model *models.BusPriceTriggerStrategyApikeyConfig) error {
	var data models.BusPriceTriggerStrategyApikeyConfig

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusPriceTriggerStrategyApikeyConfig error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusPriceTriggerStrategyApikeyConfig对象
func (e *BusPriceTriggerStrategyApikeyConfig) Insert(c *dto.BusPriceTriggerStrategyApikeyConfigInsertReq) error {
	var err error
	var data models.BusPriceTriggerStrategyApikeyConfig
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyApikeyConfigService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusPriceTriggerStrategyApikeyConfig对象
func (e *BusPriceTriggerStrategyApikeyConfig) Update(c *dto.BusPriceTriggerStrategyApikeyConfigUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusPriceTriggerStrategyApikeyConfig{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyApikeyConfigService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusPriceTriggerStrategyApikeyConfig
func (e *BusPriceTriggerStrategyApikeyConfig) Remove(d *dto.BusPriceTriggerStrategyApikeyConfigDeleteReq, p *actions.DataPermission) error {
	var data models.BusPriceTriggerStrategyApikeyConfig

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusPriceTriggerStrategyApikeyConfig error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
