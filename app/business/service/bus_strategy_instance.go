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

type BusStrategyInstance struct {
	service.Service
}

// GetPage 获取BusStrategyInstance列表
func (e *BusStrategyInstance) GetPage(c *dto.BusStrategyInstanceGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyInstance, count *int64) error {
	var err error
	var data models.BusStrategyInstance

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyInstance对象
func (e *BusStrategyInstance) Get(d *dto.BusStrategyInstanceGetReq, p *actions.DataPermission, model *models.BusStrategyInstance) error {
	var data models.BusStrategyInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyInstance error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyInstance对象
func (e *BusStrategyInstance) Insert(c *dto.BusStrategyInstanceInsertReq) error {
	var err error
	var data models.BusStrategyInstance
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyInstanceService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyInstance对象
func (e *BusStrategyInstance) Update(c *dto.BusStrategyInstanceUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyInstance{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyInstanceService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyInstance
func (e *BusStrategyInstance) Remove(d *dto.BusStrategyInstanceDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyInstance

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyInstance error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
