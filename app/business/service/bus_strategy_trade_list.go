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

type BusStrategyTradeList struct {
	service.Service
}

// GetPage 获取BusStrategyTradeList列表
func (e *BusStrategyTradeList) GetPage(c *dto.BusStrategyTradeListGetPageReq, p *actions.DataPermission, list *[]dto.BusStrategyTradeListGetPageResp, count *int64) error {
	var err error
	var data models.BusStrategyTradeList

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Joins("LEFT JOIN bus_strategy_symbol_group ON bus_strategy_symbol_group.id = bus_strategy_trade_list.symbol_group_id").
		Select("bus_strategy_trade_list.*, bus_strategy_symbol_group.group_name").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusStrategyTradeListService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusStrategyTradeList对象
func (e *BusStrategyTradeList) Get(d *dto.BusStrategyTradeListGetReq, p *actions.DataPermission, model *models.BusStrategyTradeList) error {
	var data models.BusStrategyTradeList

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusStrategyTradeList error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusStrategyTradeList对象
func (e *BusStrategyTradeList) Insert(c *dto.BusStrategyTradeListInsertReq) error {
	var err error

	var existingData models.BusStrategyTradeList
	err = e.Orm.Where("symbol_group_id = ? AND symbol = ?", c.SymbolGroupId, c.Symbol).First(&existingData).Error
	if err == nil {
		// 已存在同名记录
		e.Log.Infof("Record already exists with SymbolGroupId: %d and Symbol: %s", c.SymbolGroupId, c.Symbol)
		return nil // 返回 nil 表示操作成功，不再新增
	}

	var data models.BusStrategyTradeList
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusStrategyTradeListService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusStrategyTradeList对象
func (e *BusStrategyTradeList) Update(c *dto.BusStrategyTradeListUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusStrategyTradeList{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusStrategyTradeListService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusStrategyTradeList
func (e *BusStrategyTradeList) Remove(d *dto.BusStrategyTradeListDeleteReq, p *actions.DataPermission) error {
	var data models.BusStrategyTradeList

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusStrategyTradeList error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
