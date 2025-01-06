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

type BusStrategySymbolGroup struct {
	service.Service
}

// GetPage 获取BusStrategySymbolGroup列表
func (e *BusStrategySymbolGroup) GetPage(c *dto.BusStrategySymbolGroupGetPageReq, p *actions.DataPermission, list *[]models.BusStrategySymbolGroup, count *int64) error {
	var err error
	var data models.BusStrategySymbolGroup

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategySymbolGroupService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategySymbolGroup对象
func (e *BusStrategySymbolGroup) Get(d *dto.BusStrategySymbolGroupGetReq, p *actions.DataPermission, model *models.BusStrategySymbolGroup) error {
	var data models.BusStrategySymbolGroup

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategySymbolGroup error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategySymbolGroup对象
func (e *BusStrategySymbolGroup) Insert(c *dto.BusStrategySymbolGroupInsertReq) error {
	var err error
	var data models.BusStrategySymbolGroup
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategySymbolGroupService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategySymbolGroup对象
func (e *BusStrategySymbolGroup) Update(c *dto.BusStrategySymbolGroupUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategySymbolGroup{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategySymbolGroupService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategySymbolGroup
func (e *BusStrategySymbolGroup) Remove(d *dto.BusStrategySymbolGroupDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategySymbolGroup

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategySymbolGroup error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
