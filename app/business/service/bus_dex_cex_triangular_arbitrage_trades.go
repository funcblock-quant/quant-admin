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

type StrategyDexCexTriangularArbitrageTrades struct {
	service.Service
}

// GetPage 获取StrategyDexCexTriangularArbitrageTrades列表
func (e *StrategyDexCexTriangularArbitrageTrades) GetPage(c *dto.StrategyDexCexTriangularArbitrageTradesGetPageReq, p *actions.DataPermission, list *[]dto.StrategyDexCexTriangularArbitrageTradesGetPageResp, count *int64) error {
	var err error
	var data models.StrategyDexCexTriangularArbitrageTrades

	query := e.Orm.Model(&data).
		Select("strategy_dex_cex_triangular_arbitrage_trades.*, opportunities.cex_target_asset as symbol").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Where("strategy_dex_cex_triangular_arbitrage_trades.error IS NULL OR strategy_dex_cex_triangular_arbitrage_trades.error = ''")

	// 1. Profit 计算字段过滤
	if c.MinProfit != "" && c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) BETWEEN ? AND ?", c.MinProfit, c.MaxProfit)
	} else if c.MinProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) >= ?", c.MinProfit)
	} else if c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) <= ?", c.MaxProfit)
	}

	if c.Symbol != "" {
		query = query.Where("opportunities.cex_target_asset = ?", c.Symbol)
	}

	err = query.Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageTradesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取StrategyDexCexTriangularArbitrageTrades对象
func (e *StrategyDexCexTriangularArbitrageTrades) Get(d *dto.StrategyDexCexTriangularArbitrageTradesGetReq, p *actions.DataPermission, model *models.StrategyDexCexTriangularArbitrageTrades) error {
	var data models.StrategyDexCexTriangularArbitrageTrades

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStrategyDexCexTriangularArbitrageTrades error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建StrategyDexCexTriangularArbitrageTrades对象
func (e *StrategyDexCexTriangularArbitrageTrades) Insert(c *dto.StrategyDexCexTriangularArbitrageTradesInsertReq) error {
	var err error
	var data models.StrategyDexCexTriangularArbitrageTrades
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageTradesService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改StrategyDexCexTriangularArbitrageTrades对象
func (e *StrategyDexCexTriangularArbitrageTrades) Update(c *dto.StrategyDexCexTriangularArbitrageTradesUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.StrategyDexCexTriangularArbitrageTrades{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageTradesService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除StrategyDexCexTriangularArbitrageTrades
func (e *StrategyDexCexTriangularArbitrageTrades) Remove(d *dto.StrategyDexCexTriangularArbitrageTradesDeleteReq, p *actions.DataPermission) error {
	var data models.StrategyDexCexTriangularArbitrageTrades

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStrategyDexCexTriangularArbitrageTrades error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
