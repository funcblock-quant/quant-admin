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

type BusDexCexDepositWithdrawRecord struct {
	service.Service
}

// GetPage 获取BusDexCexDepositWithdrawRecord列表
func (e *BusDexCexDepositWithdrawRecord) GetPage(c *dto.BusDexCexDepositWithdrawRecordGetPageReq, p *actions.DataPermission, list *[]models.BusDexCexDepositWithdrawRecord, count *int64) error {
	var err error
	var data models.BusDexCexDepositWithdrawRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexDepositWithdrawRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexCexDepositWithdrawRecord对象
func (e *BusDexCexDepositWithdrawRecord) Get(d *dto.BusDexCexDepositWithdrawRecordGetReq, p *actions.DataPermission, model *models.BusDexCexDepositWithdrawRecord) error {
	var data models.BusDexCexDepositWithdrawRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexDepositWithdrawRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexDepositWithdrawRecord对象
func (e *BusDexCexDepositWithdrawRecord) Insert(c *dto.BusDexCexDepositWithdrawRecordInsertReq) error {
    var err error
    var data models.BusDexCexDepositWithdrawRecord
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexDepositWithdrawRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexCexDepositWithdrawRecord对象
func (e *BusDexCexDepositWithdrawRecord) Update(c *dto.BusDexCexDepositWithdrawRecordUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.BusDexCexDepositWithdrawRecord{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("BusDexCexDepositWithdrawRecordService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除BusDexCexDepositWithdrawRecord
func (e *BusDexCexDepositWithdrawRecord) Remove(d *dto.BusDexCexDepositWithdrawRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexDepositWithdrawRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveBusDexCexDepositWithdrawRecord error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
