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

type BusSpotOrderRecord struct {
	service.Service
}

// GetPage 获取BusSpotOrderRecord列表
func (e *BusSpotOrderRecord) GetPage(c *dto.BusSpotOrderRecordGetPageReq, p *actions.DataPermission, list *[]models.BusSpotOrderRecord, count *int64) error {
	var err error
	var data models.BusSpotOrderRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusSpotOrderRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusSpotOrderRecord对象
func (e *BusSpotOrderRecord) Get(d *dto.BusSpotOrderRecordGetReq, p *actions.DataPermission, model *models.BusSpotOrderRecord) error {
	var data models.BusSpotOrderRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusSpotOrderRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusSpotOrderRecord对象
func (e *BusSpotOrderRecord) Insert(c *dto.BusSpotOrderRecordInsertReq) error {
	var err error
	var data models.BusSpotOrderRecord
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusSpotOrderRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusSpotOrderRecord对象
func (e *BusSpotOrderRecord) Update(c *dto.BusSpotOrderRecordUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusSpotOrderRecord{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusSpotOrderRecordService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusSpotOrderRecord
func (e *BusSpotOrderRecord) Remove(d *dto.BusSpotOrderRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusSpotOrderRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusSpotOrderRecord error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
