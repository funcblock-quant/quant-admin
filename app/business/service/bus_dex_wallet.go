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

type BusDexWallet struct {
	service.Service
}

// GetPage 获取BusDexWallet列表
func (e *BusDexWallet) GetPage(c *dto.BusDexWalletGetPageReq, p *actions.DataPermission, list *[]models.BusDexWallet, count *int64) error {
	var err error
	var data models.BusDexWallet

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexWalletService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexWallet对象
func (e *BusDexWallet) Get(d *dto.BusDexWalletGetReq, p *actions.DataPermission, model *models.BusDexWallet) error {
	var data models.BusDexWallet

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexWallet error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexWallet对象
func (e *BusDexWallet) Insert(c *dto.BusDexWalletInsertReq) error {
    var err error
    var data models.BusDexWallet
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexWalletService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexWallet对象
func (e *BusDexWallet) Update(c *dto.BusDexWalletUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.BusDexWallet{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("BusDexWalletService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除BusDexWallet
func (e *BusDexWallet) Remove(d *dto.BusDexWalletDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexWallet

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveBusDexWallet error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
