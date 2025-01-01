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

type BusStrategyConfigDict struct {
	service.Service
}

// GetPage 获取BusStrategyConfigDict列表
func (e *BusStrategyConfigDict) GetPage(c *dto.BusStrategyConfigDictGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyConfigDict, count *int64) error {
	var err error
	var data models.BusStrategyConfigDict

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyConfigDictService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyConfigDict对象
func (e *BusStrategyConfigDict) Get(d *dto.BusStrategyConfigDictGetReq, p *actions.DataPermission, model *models.BusStrategyConfigDict) error {
	var data models.BusStrategyConfigDict

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyConfigDict error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyConfigDict对象
func (e *BusStrategyConfigDict) Insert(c *dto.BusStrategyConfigDictInsertReq) error {
	var err error
	var data models.BusStrategyConfigDict
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyConfigDictService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyConfigDict对象
func (e *BusStrategyConfigDict) Update(c *dto.BusStrategyConfigDictUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyConfigDict{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyConfigDictService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyConfigDict
func (e *BusStrategyConfigDict) Remove(d *dto.BusStrategyConfigDictDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyConfigDict

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyConfigDict error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
