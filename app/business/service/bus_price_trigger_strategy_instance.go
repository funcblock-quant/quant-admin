package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	"strconv"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusPriceTriggerStrategyInstance struct {
	service.Service
}

// GetPage 获取BusPriceTriggerStrategyInstance列表
func (e *BusPriceTriggerStrategyInstance) GetPage(c *dto.BusPriceTriggerStrategyInstanceGetPageReq, p *actions.DataPermission, list *[]dto.BusPriceTriggerStrategyResp, count *int64) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance
	var detail models.BusPriceMonitorForOptionHedging

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error

	for index, strategy := range *list {
		details := make([]models.BusPriceMonitorForOptionHedging, 0)
		err := e.Orm.Model(&detail).Where("strategy_instance_id = ?", strategy.Id).Order("id desc").Find(&details).Error
		if err != nil {
			e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get details error:%s \r\n", err)
			return err
		}
		(*list)[index].Details = details
	}
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Get(d *dto.BusPriceTriggerStrategyInstanceGetReq, p *actions.DataPermission, model *models.BusPriceTriggerStrategyInstance) error {
	var data models.BusPriceTriggerStrategyInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusPriceTriggerStrategyInstance error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Insert(c *dto.BusPriceTriggerStrategyInstanceInsertReq) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance
	c.Generate(&data)
	data.Status = "created"
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService error:%s \r\n", tx.Error)
		return tx.Error
	}

	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService Insert error:%s \r\n", err)
		return err
	}
	e.Log.Infof("instance id : %d\r\n", data.Id)

	//创建成功后， 自动通过grpc启动
	request := &trigger_service.StartTriggerRequest{
		InstanceId: strconv.Itoa(data.Id),
		OpenPrice:  c.OpenPrice,
		ClosePrice: c.ClosePrice,
		Side:       c.Side,
		Amount:     c.Amount,
		Symbol:     c.Symbol,
		StopTime:   strconv.FormatInt(c.CloseTime.UnixMilli(), 10),
	}
	_, err = client.StartInstance(request)
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("Service grpc start error:%s \r\n", err)
		return err
	}

	tx.Commit()

	err = e.Orm.Model(&data).Update("status", "started").Error
	if err != nil {
		e.Log.Errorf("start trigger update status failed:%s \r\n", err)
		return err
	}

	return nil
}

// Update 修改BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Update(c *dto.BusPriceTriggerStrategyInstanceUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusPriceTriggerStrategyInstance{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusPriceTriggerStrategyInstance
func (e *BusPriceTriggerStrategyInstance) Remove(d *dto.BusPriceTriggerStrategyInstanceDeleteReq, p *actions.DataPermission) error {
	var data models.BusPriceTriggerStrategyInstance

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusPriceTriggerStrategyInstance error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
