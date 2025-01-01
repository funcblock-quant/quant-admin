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

type BusStrategyWatchList struct {
	service.Service
}

// GetPage 获取BusStrategyWatchList列表
func (e *BusStrategyWatchList) GetPage(c *dto.BusStrategyWatchListGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyWatchList, count *int64) error {
	var err error
	var data models.BusStrategyWatchList

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyWatchListService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyWatchList对象
func (e *BusStrategyWatchList) Get(d *dto.BusStrategyWatchListGetReq, p *actions.DataPermission, model *models.BusStrategyWatchList) error {
	var data models.BusStrategyWatchList

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyWatchList error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyWatchList对象
func (e *BusStrategyWatchList) Insert(c *dto.BusStrategyWatchListInsertReq) error {
	var err error
	var data models.BusStrategyWatchList
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyWatchListService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyWatchList对象
func (e *BusStrategyWatchList) Update(c *dto.BusStrategyWatchListUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyWatchList{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyWatchListService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyWatchList
func (e *BusStrategyWatchList) Remove(d *dto.BusStrategyWatchListDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyWatchList

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyWatchList error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
