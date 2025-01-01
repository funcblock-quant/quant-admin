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

type BusTriangleOrdersRecord struct {
	service.Service
}

// GetPage 获取BusTriangleOrdersRecord列表
func (e *BusTriangleOrdersRecord) GetPage(c *dto.BusTriangleOrdersRecordGetPageReq, p *actions.DataPermission, list *[]models.BusTriangleOrdersRecord, count *int64) error {
	var err error
	var data models.BusTriangleOrdersRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusTriangleOrdersRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusTriangleOrdersRecord对象
func (e *BusTriangleOrdersRecord) Get(d *dto.BusTriangleOrdersRecordGetReq, p *actions.DataPermission, model *models.BusTriangleOrdersRecord) error {
	var data models.BusTriangleOrdersRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusTriangleOrdersRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusTriangleOrdersRecord对象
func (e *BusTriangleOrdersRecord) Insert(c *dto.BusTriangleOrdersRecordInsertReq) error {
	var err error
	var data models.BusTriangleOrdersRecord
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusTriangleOrdersRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusTriangleOrdersRecord对象
func (e *BusTriangleOrdersRecord) Update(c *dto.BusTriangleOrdersRecordUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusTriangleOrdersRecord{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusTriangleOrdersRecordService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusTriangleOrdersRecord
func (e *BusTriangleOrdersRecord) Remove(d *dto.BusTriangleOrdersRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusTriangleOrdersRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusTriangleOrdersRecord error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
