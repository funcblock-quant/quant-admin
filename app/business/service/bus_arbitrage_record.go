package service

import (
	"errors"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusArbitrageRecord struct {
	service.Service
}

// GetPage 获取BusArbitrageRecord列表
func (e *BusArbitrageRecord) GetPage(c *dto.BusArbitrageRecordGetPageReq, p *actions.DataPermission, list *[]models.BusArbitrageRecord, count *int64) error {
	var err error
	var data models.BusArbitrageRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusArbitrageRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusArbitrageRecord对象
func (e *BusArbitrageRecord) Get(d *dto.BusArbitrageRecordGetReq, p *actions.DataPermission, model *models.BusArbitrageRecord) error {
	var data models.BusArbitrageRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusArbitrageRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusArbitrageRecord对象
func (e *BusArbitrageRecord) Insert(c *dto.BusArbitrageRecordInsertReq) error {
	var err error
	var data models.BusArbitrageRecord
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusArbitrageRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusArbitrageRecord对象
func (e *BusArbitrageRecord) Update(c *dto.BusArbitrageRecordUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusArbitrageRecord{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusArbitrageRecordService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusArbitrageRecord
func (e *BusArbitrageRecord) Remove(d *dto.BusArbitrageRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusArbitrageRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusArbitrageRecord error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
