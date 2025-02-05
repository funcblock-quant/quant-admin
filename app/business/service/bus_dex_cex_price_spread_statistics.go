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

type BusDexCexPriceSpreadStatistics struct {
	service.Service
}

// GetPage 获取BusDexCexPriceSpreadStatistics列表
func (e *BusDexCexPriceSpreadStatistics) GetPage(c *dto.BusDexCexPriceSpreadStatisticsGetPageReq, p *actions.DataPermission, list *[]models.BusDexCexPriceSpreadStatistics, count *int64) error {
	var err error
	var data models.BusDexCexPriceSpreadStatistics

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadStatisticsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexCexPriceSpreadStatistics对象
func (e *BusDexCexPriceSpreadStatistics) Get(d *dto.BusDexCexPriceSpreadStatisticsGetReq, p *actions.DataPermission, model *models.BusDexCexPriceSpreadStatistics) error {
	var data models.BusDexCexPriceSpreadStatistics

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexPriceSpreadStatistics error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexPriceSpreadStatistics对象
func (e *BusDexCexPriceSpreadStatistics) Insert(c *dto.BusDexCexPriceSpreadStatisticsInsertReq) error {
	var err error
	var data models.BusDexCexPriceSpreadStatistics
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadStatisticsService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexCexPriceSpreadStatistics对象
func (e *BusDexCexPriceSpreadStatistics) Update(c *dto.BusDexCexPriceSpreadStatisticsUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusDexCexPriceSpreadStatistics{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadStatisticsService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusDexCexPriceSpreadStatistics
func (e *BusDexCexPriceSpreadStatistics) Remove(d *dto.BusDexCexPriceSpreadStatisticsDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexPriceSpreadStatistics

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusDexCexPriceSpreadStatistics error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
