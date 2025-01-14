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

type BusStrategyConfigSchema struct {
	service.Service
}

// GetPage 获取BusStrategyConfigSchema列表
func (e *BusStrategyConfigSchema) GetPage(c *dto.BusStrategyConfigSchemaGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyConfigSchema, count *int64) error {
	var err error
	var data models.BusStrategyConfigSchema

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyConfigSchemaService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyConfigSchema对象
func (e *BusStrategyConfigSchema) Get(d *dto.BusStrategyConfigSchemaGetReq, p *actions.DataPermission, model *models.BusStrategyConfigSchema) error {
	var data models.BusStrategyConfigSchema

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyConfigSchema error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyConfigSchema对象
func (e *BusStrategyConfigSchema) Insert(c *dto.BusStrategyConfigSchemaInsertReq) error {
	var err error
	var data models.BusStrategyConfigSchema
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyConfigSchemaService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyConfigSchema对象
func (e *BusStrategyConfigSchema) Update(c *dto.BusStrategyConfigSchemaUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyConfigSchema{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyConfigSchemaService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyConfigSchema
func (e *BusStrategyConfigSchema) Remove(d *dto.BusStrategyConfigSchemaDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyConfigSchema

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyConfigSchema error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
