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

type BusStrategyBaseInfo struct {
	service.Service
}

// GetPage 获取BusStrategyBaseInfo列表
func (e *BusStrategyBaseInfo) GetPage(c *dto.BusStrategyBaseInfoGetPageReq, p *actions.DataPermission, list *[]models.BusStrategyBaseInfo, count *int64) error {
	var err error
	var data models.BusStrategyBaseInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyBaseInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Get(d *dto.BusStrategyBaseInfoGetReq, p *actions.DataPermission, model *models.BusStrategyBaseInfo) error {
	var data models.BusStrategyBaseInfo

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyBaseInfo error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Insert(c *dto.BusStrategyBaseInfoInsertReq) error {
	var err error
	var data models.BusStrategyBaseInfo
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyBaseInfoService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyBaseInfo对象
func (e *BusStrategyBaseInfo) Update(c *dto.BusStrategyBaseInfoUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyBaseInfo{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyBaseInfoService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyBaseInfo
func (e *BusStrategyBaseInfo) Remove(d *dto.BusStrategyBaseInfoDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyBaseInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyBaseInfo error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
