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

type GatemtMakerOrders struct {
	service.Service
}

// GetPage 获取GatemtMakerOrders列表
func (e *GatemtMakerOrders) GetPage(c *dto.GatemtMakerOrdersGetPageReq, p *actions.DataPermission, list *[]models.GatemtMakerOrders, count *int64) error {
	var err error
	var data models.GatemtMakerOrders

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GatemtMakerOrdersService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GatemtMakerOrders对象
func (e *GatemtMakerOrders) Get(d *dto.GatemtMakerOrdersGetReq, p *actions.DataPermission, model *models.GatemtMakerOrders) error {
	var data models.GatemtMakerOrders

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGatemtMakerOrders error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GatemtMakerOrders对象
func (e *GatemtMakerOrders) Insert(c *dto.GatemtMakerOrdersInsertReq) error {
	var err error
	var data models.GatemtMakerOrders
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GatemtMakerOrdersService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GatemtMakerOrders对象
func (e *GatemtMakerOrders) Update(c *dto.GatemtMakerOrdersUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.GatemtMakerOrders{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("GatemtMakerOrdersService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除GatemtMakerOrders
func (e *GatemtMakerOrders) Remove(d *dto.GatemtMakerOrdersDeleteReq, p *actions.DataPermission) error {
	var data models.GatemtMakerOrders

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveGatemtMakerOrders error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
