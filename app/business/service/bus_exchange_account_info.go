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

type BusExchangeAccountInfo struct {
	service.Service
}

// GetPage 获取BusExchangeAccountInfo列表
func (e *BusExchangeAccountInfo) GetPage(c *dto.BusExchangeAccountInfoGetPageReq, p *actions.DataPermission, list *[]models.BusExchangeAccountInfo, count *int64) error {
	var err error
	var data models.BusExchangeAccountInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusExchangeAccountInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Get(d *dto.BusExchangeAccountInfoGetReq, p *actions.DataPermission, model *models.BusExchangeAccountInfo) error {
	var data models.BusExchangeAccountInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusExchangeAccountInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Insert(c *dto.BusExchangeAccountInfoInsertReq) error {
	var err error
	var data models.BusExchangeAccountInfo
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusExchangeAccountInfoService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusExchangeAccountInfo对象
func (e *BusExchangeAccountInfo) Update(c *dto.BusExchangeAccountInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusExchangeAccountInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusExchangeAccountInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusExchangeAccountInfo
func (e *BusExchangeAccountInfo) Remove(d *dto.BusExchangeAccountInfoDeleteReq, p *actions.DataPermission) error {
	var data models.BusExchangeAccountInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusExchangeAccountInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
