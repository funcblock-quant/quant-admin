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

type BusStrategyInstanceConfig struct {
	service.Service
}

// GetPage 获取BusStrategyInstanceConfig列表
func (e *BusStrategyInstanceConfig) GetPage(c *dto.BusStrategyInstanceConfigGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyInstanceConfig, count *int64) error {
	var err error
	var data models.BusStrategyInstanceConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyInstanceConfigService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyInstanceConfig对象
func (e *BusStrategyInstanceConfig) Get(d *dto.BusStrategyInstanceConfigGetReq, p *actions.DataPermission, model *models.BusStrategyInstanceConfig) error {
	var data models.BusStrategyInstanceConfig

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyInstanceConfig error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyInstanceConfig对象
func (e *BusStrategyInstanceConfig) Insert(c *dto.BusStrategyInstanceConfigInsertReq) error {
	var err error
	var data models.BusStrategyInstanceConfig
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyInstanceConfigService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyInstanceConfig对象
func (e *BusStrategyInstanceConfig) Update(c *dto.BusStrategyInstanceConfigUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyInstanceConfig{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyInstanceConfigService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyInstanceConfig
func (e *BusStrategyInstanceConfig) Remove(d *dto.BusStrategyInstanceConfigDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstanceConfig

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyInstanceConfig error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
