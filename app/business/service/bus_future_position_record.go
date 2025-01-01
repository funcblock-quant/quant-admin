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

type BusFuturePositionRecord struct {
	service.Service
}

// GetPage 获取BusFuturePositionRecord列表
func (e *BusFuturePositionRecord) GetPage(c *dto.BusFuturePositionRecordGetPageReq, p *actions.DataPermission, list *[]models.BusFuturePositionRecord, count *int64) error {
	var err error
	var data models.BusFuturePositionRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusFuturePositionRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusFuturePositionRecord对象
func (e *BusFuturePositionRecord) Get(d *dto.BusFuturePositionRecordGetReq, p *actions.DataPermission, model *models.BusFuturePositionRecord) error {
	var data models.BusFuturePositionRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusFuturePositionRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusFuturePositionRecord对象
func (e *BusFuturePositionRecord) Insert(c *dto.BusFuturePositionRecordInsertReq) error {
	var err error
	var data models.BusFuturePositionRecord
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusFuturePositionRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusFuturePositionRecord对象
func (e *BusFuturePositionRecord) Update(c *dto.BusFuturePositionRecordUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusFuturePositionRecord{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusFuturePositionRecordService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusFuturePositionRecord
func (e *BusFuturePositionRecord) Remove(d *dto.BusFuturePositionRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusFuturePositionRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusFuturePositionRecord error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
