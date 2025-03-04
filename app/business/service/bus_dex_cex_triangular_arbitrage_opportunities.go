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

type StrategyDexCexTriangularArbitrageOpportunities struct {
	service.Service
}

// GetPage 获取StrategyDexCexTriangularArbitrageOpportunities列表
func (e *StrategyDexCexTriangularArbitrageOpportunities) GetPage(c *dto.StrategyDexCexTriangularArbitrageOpportunitiesGetPageReq, p *actions.DataPermission, list *[]models.StrategyDexCexTriangularArbitrageOpportunities, count *int64) error {
	var err error
	var data models.StrategyDexCexTriangularArbitrageOpportunities

	query := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)

	if c.MinProfit != "" && c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) BETWEEN ? AND ?", c.MinProfit, c.MaxProfit)
	} else if c.MinProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) >= ?", c.MinProfit)
	} else if c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) <= ?", c.MaxProfit)
	}

	if c.MinProfitPercent != "" && c.MaxProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) BETWEEN ? AND ?", c.MinProfitPercent, c.MaxProfitPercent)
	} else if c.MinProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) >= ?", c.MinProfitPercent)
	} else if c.MaxProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_opportunities.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_opportunities.cex_buy_quote_amount) <= ?", c.MaxProfitPercent)
	}

	if c.Symbol != "" {
		query = query.Where("strategy_dex_cex_triangular_arbitrage_opportunities.cex_target_asset = ?", c.Symbol)
	}

	err = query.Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageOpportunitiesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取StrategyDexCexTriangularArbitrageOpportunities对象
func (e *StrategyDexCexTriangularArbitrageOpportunities) Get(d *dto.StrategyDexCexTriangularArbitrageOpportunitiesGetReq, p *actions.DataPermission, model *models.StrategyDexCexTriangularArbitrageOpportunities) error {
	var data models.StrategyDexCexTriangularArbitrageOpportunities

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStrategyDexCexTriangularArbitrageOpportunities error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建StrategyDexCexTriangularArbitrageOpportunities对象
func (e *StrategyDexCexTriangularArbitrageOpportunities) Insert(c *dto.StrategyDexCexTriangularArbitrageOpportunitiesInsertReq) error {
	var err error
	var data models.StrategyDexCexTriangularArbitrageOpportunities
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageOpportunitiesService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改StrategyDexCexTriangularArbitrageOpportunities对象
func (e *StrategyDexCexTriangularArbitrageOpportunities) Update(c *dto.StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.StrategyDexCexTriangularArbitrageOpportunities{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("StrategyDexCexTriangularArbitrageOpportunitiesService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除StrategyDexCexTriangularArbitrageOpportunities
func (e *StrategyDexCexTriangularArbitrageOpportunities) Remove(d *dto.StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq, p *actions.DataPermission) error {
	var data models.StrategyDexCexTriangularArbitrageOpportunities

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveStrategyDexCexTriangularArbitrageOpportunities error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
