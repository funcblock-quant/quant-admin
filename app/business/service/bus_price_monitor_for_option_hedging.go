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

type BusPriceMonitorForOptionHedging struct {
	service.Service
}

// GetPage 获取BusPriceMonitorForOptionHedging列表
func (e *BusPriceMonitorForOptionHedging) GetPage(c *dto.BusPriceMonitorForOptionHedgingGetPageReq, p *actions.DataPermission, list *[]models.BusPriceMonitorForOptionHedging, count *int64) error {
	var err error
	var data models.BusPriceMonitorForOptionHedging

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusPriceMonitorForOptionHedgingService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusPriceMonitorForOptionHedging对象
func (e *BusPriceMonitorForOptionHedging) Get(d *dto.BusPriceMonitorForOptionHedgingGetReq, p *actions.DataPermission, model *models.BusPriceMonitorForOptionHedging) error {
	var data models.BusPriceMonitorForOptionHedging

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusPriceMonitorForOptionHedging error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusPriceMonitorForOptionHedging对象
func (e *BusPriceMonitorForOptionHedging) Insert(c *dto.BusPriceMonitorForOptionHedgingInsertReq) error {
	var err error
	var data models.BusPriceMonitorForOptionHedging
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceMonitorForOptionHedgingService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusPriceMonitorForOptionHedging对象
func (e *BusPriceMonitorForOptionHedging) Update(c *dto.BusPriceMonitorForOptionHedgingUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusPriceMonitorForOptionHedging{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusPriceMonitorForOptionHedgingService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusPriceMonitorForOptionHedging
func (e *BusPriceMonitorForOptionHedging) Remove(d *dto.BusPriceMonitorForOptionHedgingDeleteReq, p *actions.DataPermission) error {
	var data models.BusPriceMonitorForOptionHedging

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusPriceMonitorForOptionHedging error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
