package service

import (
	"errors"
	"quanta-admin/app/business/constant"
	"quanta-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusDexCexTriangularArbitrageRecord struct {
	service.Service
}

// GetPage 获取BusDexCexTriangularArbitrageRecord列表
func (e *BusDexCexTriangularArbitrageRecord) GetPage(c *dto.BusDexCexTriangularArbitrageRecordGetPageReq, p *actions.DataPermission, list *[]models.BusDexCexTriangularArbitrageRecord, count *int64) error {
	var err error
	var data models.BusDexCexTriangularArbitrageRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexTriangularArbitrageRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexCexTriangularArbitrageRecord对象
func (e *BusDexCexTriangularArbitrageRecord) Get(d *dto.BusDexCexTriangularArbitrageRecordGetReq, p *actions.DataPermission, model *models.BusDexCexTriangularArbitrageRecord) error {
	var data models.BusDexCexTriangularArbitrageRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexTriangularArbitrageRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexTriangularArbitrageRecord对象
func (e *BusDexCexTriangularArbitrageRecord) Insert(c *dto.BusDexCexTriangularArbitrageRecordInsertReq) error {
	var err error
	var data models.BusDexCexTriangularArbitrageRecord
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexTriangularArbitrageRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexCexTriangularArbitrageRecord对象
func (e *BusDexCexTriangularArbitrageRecord) Update(c *dto.BusDexCexTriangularArbitrageRecordUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusDexCexTriangularArbitrageRecord{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusDexCexTriangularArbitrageRecordService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusDexCexTriangularArbitrageRecord
func (e *BusDexCexTriangularArbitrageRecord) Remove(d *dto.BusDexCexTriangularArbitrageRecordDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexTriangularArbitrageRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusDexCexTriangularArbitrageRecord error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *BusDexCexTriangularArbitrageRecord) QueryArbitrageOpportunityList(d *dto.BusArbitrageOpportunityGetReq, p *actions.DataPermission, list *[]dto.BusArbitrageOpportunityGetResp) error {
	var err error
	// 获取选中的交易所信息
	isntanceId := d.StrategyInstanceId
	var instance = models.BusStrategyInstance{}
	err = e.Orm.Model(&models.BusStrategyInstance{}).First(&instance, isntanceId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("no record found")
			return nil
		}
		e.Log.Errorf("QueryArbitrageOpportunityList Query error:%s \r\n", err)
		return err
	}
	e.Log.Infof("strategy instance:%+v", instance)

	model := models.BusDexCexTriangularArbitrageRecord{}
	tx := e.Orm.Scopes(
		actions.Permission(model.TableName(), p)).
		Model(&model).
		Select("dex_symbol as symbol, dex_platform, cex_platform, count(*) as count").
		Group("dex_symbol, dex_platform, cex_platform")
	tx.Where("type = ?", constant.ARBITRAGE_TYPE_SIMULATE)
	if len(d.Symbols) > 0 {
		tx.Where("dex_symbol in ?", d.Symbols)
	}
	if d.BeginTime != "" && d.EndTime != "" {
		tx.Where("created_at > ? AND created_at < ?", d.BeginTime, d.EndTime)
	}
	if d.MinProfit != "" {
		decimalMinProfit := utils.ParseString2DBDecimal(d.MinProfit)
		tx.Where("(quote_token_profit + base_token_profit) > ?", decimalMinProfit)
	}
	if d.MaxProfit != "" {
		decimalMaxProfit := utils.ParseString2DBDecimal(d.MaxProfit)
		tx.Where("(quote_token_profit + base_token_profit) <?", decimalMaxProfit)
	}

	err = tx.Debug().Find(list).Error
	if err != nil {
		e.Log.Errorf("QueryArbitrageOpportunityList Query error:%s \r\n", err)
		return err
	}

	e.Log.Infof("QueryArbitrageOpportunityList: %v\r\n", list)
	return err
}
