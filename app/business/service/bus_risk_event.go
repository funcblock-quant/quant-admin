package service

import (
	"errors"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusRiskEvent struct {
	service.Service
}

// GetPage 获取BusRiskEvent列表
func (e *BusRiskEvent) GetPage(c *dto.BusRiskEventGetPageReq, p *actions.DataPermission, list *[]models.BusRiskEvent, count *int64) error {
	var err error
	var data models.BusRiskEvent

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusRiskEventService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusRiskEvent对象
func (e *BusRiskEvent) Get(d *dto.BusRiskEventGetReq, p *actions.DataPermission, model *models.BusRiskEvent) error {
	var data models.BusRiskEvent

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusRiskEvent error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusRiskEvent对象
func (e *BusRiskEvent) Insert(c *dto.BusRiskEventInsertReq) error {
	var err error
	var data models.BusRiskEvent
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusRiskEventService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusRiskEvent对象
func (e *BusRiskEvent) Update(c *dto.BusRiskEventUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusRiskEvent{}

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		)

	updateData := map[string]interface{}{
		"is_recovered": true,
		"recovered_by": c.UpdateBy,
		"recovered_at": time.Now(),
	}

	err = db.
		Where("id = ?", c.Id).
		Updates(&updateData).Error
	if err != nil {
		e.Log.Errorf("BusRiskEventService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusRiskEvent
func (e *BusRiskEvent) Remove(d *dto.BusRiskEventDeleteReq, p *actions.DataPermission) error {
	var data models.BusRiskEvent

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusRiskEvent error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
