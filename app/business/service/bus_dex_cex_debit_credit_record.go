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

type BusDexCexDebitCreditRecord struct {
	service.Service
}

// GetPage 获取BusDexCexDebitCreditRecord列表
func (e *BusDexCexDebitCreditRecord) GetPage(c *dto.BusDexCexDebitCreditRecordGetPageReq, p *actions.DataPermission, list *[]models.BusDexCexDebitCreditRecord, count *int64) error {
	var err error
	var data models.BusDexCexDebitCreditRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexDebitCreditRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexCexDebitCreditRecord对象
func (e *BusDexCexDebitCreditRecord) Get(d *dto.BusDexCexDebitCreditRecordGetReq, p *actions.DataPermission, model *models.BusDexCexDebitCreditRecord) error {
	var data models.BusDexCexDebitCreditRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexDebitCreditRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexDebitCreditRecord对象
func (e *BusDexCexDebitCreditRecord) Insert(c *dto.BusDexCexDebitCreditRecordInsertReq) error {
    var err error
    var data models.BusDexCexDebitCreditRecord
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexDebitCreditRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexCexDebitCreditRecord对象
func (e *BusDexCexDebitCreditRecord) Update(c *dto.BusDexCexDebitCreditRecordUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.BusDexCexDebitCreditRecord{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("BusDexCexDebitCreditRecordService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除BusDexCexDebitCreditRecord
func (e *BusDexCexDebitCreditRecord) Remove(d *dto.BusDexCexDebitCreditRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexDebitCreditRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveBusDexCexDebitCreditRecord error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
