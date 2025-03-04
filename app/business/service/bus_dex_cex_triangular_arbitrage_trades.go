package service

import (
	"errors"
	"fmt"
	"math"
	lark "quanta-admin/common/notification"
	ext "quanta-admin/config"
	"strconv"
	"strings"
	"time"

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
		)

	// 1. Profit 计算字段过滤
	if c.MinProfit != "" && c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) BETWEEN ? AND ?", c.MinProfit, c.MaxProfit)
	} else if c.MinProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) >= ?", c.MinProfit)
	} else if c.MaxProfit != "" {
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) <= ?", c.MaxProfit)
	}

	if c.MinProfitPercent != "" && c.MaxProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) BETWEEN ? AND ?", c.MinProfitPercent, c.MaxProfitPercent)
	} else if c.MinProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) >= ?", c.MinProfitPercent)
	} else if c.MaxProfitPercent != "" {
		query = query.Where("((strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount)/strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount) <= ?", c.MaxProfitPercent)
	}

	if c.IsSuccess {
		query = query.Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = 1")
	} else {
		query = query.Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = 0 or strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success = 0 or strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = 0")
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
func (e *StrategyDexCexTriangularArbitrageTrades) Get(d *dto.StrategyDexCexTriangularArbitrageTradesGetReq, p *actions.DataPermission, model *dto.StrategyDexCexTriangularArbitrageTradesGetDetailResp) error {
	var data models.StrategyDexCexTriangularArbitrageTrades

	e.Log.Infof("get detail for trade %s \r\n", d.Id)

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(&model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStrategyDexCexTriangularArbitrageTrades error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	e.Log.Infof("get trade %v \r\n", model)
	var oppo models.StrategyDexCexTriangularArbitrageOpportunities
	err = e.Orm.Model(&oppo).
		Scopes(
			actions.Permission(oppo.TableName(), p),
		).
		Where("opportunity_id = ?", model.OpportunityId).
		First(&oppo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStrategyDexCexTriangularArbitrageTrades error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	e.Log.Infof("get oppo %v \r\n", oppo)

	model.DexPoolType = oppo.DexPoolType
	model.DexPoolId = oppo.DexPoolId
	model.CexTargetAsset = oppo.CexTargetAsset
	model.CexQuoteAsset = oppo.CexQuoteAsset
	model.OppoDexSolAmount = oppo.DexSolAmount
	model.OppoDexTargetAmount = oppo.DexTargetAmount
	model.OppoCexSellQuantity = oppo.CexSellQuantity
	model.OppoCexBuyQuantity = oppo.CexBuyQuantity
	model.OppoCexSellQuoteAmount = oppo.CexSellQuoteAmount
	model.OppoCexBuyQuoteAmount = oppo.CexBuyQuoteAmount

	return nil
}

// GetDexCexTriangularTraderStatistics 获取StrategyDexCexTriangularArbitrageTrades对象统计信息
func (e *StrategyDexCexTriangularArbitrageTrades) GetDexCexTriangularTraderStatistics(d *dto.StrategyDexCexTriangularArbitrageTradesGetStatisticsReq, p *actions.DataPermission, model *dto.StrategyDexCexTriangularArbitrageTradesGetStatisticsResp) error {
	e.Log.Info("开始获取套利交易统计信息")

	// 计算 24 小时前的时间
	last24Hours := time.Now().Add(-24 * time.Hour)
	db := e.Orm
	// **单个查询获取多个统计数据**
	row := db.Raw(`
			SELECT 
				COUNT(*) AS totalTrade,
				COALESCE(SUM(CASE WHEN dex_success = 1 AND cex_sell_success = 1 AND cex_buy_success = 1 THEN 1 ELSE 0 END),0) AS totalSuccessTrade,
				COALESCE(SUM(CASE WHEN created_at >= ? THEN 1 ELSE 0 END),0) AS dailyTotalTrade,
				COALESCE(SUM(CASE WHEN dex_success = 1 AND cex_sell_success = 1 AND cex_buy_success = 1 AND created_at >= ? THEN 1 ELSE 0 END),0) AS dailyTotalSuccessTrade,
				COALESCE(SUM(cex_sell_quote_amount - cex_buy_quote_amount), 0) AS totalProfit,
				COALESCE(SUM(CASE WHEN created_at >= ? THEN (cex_sell_quote_amount - cex_buy_quote_amount) ELSE 0 END), 0) AS dailyTotalProfit,
				COALESCE(SUM(cex_sell_quote_amount + cex_buy_quote_amount), 0) AS totalTradeVolume,
				COALESCE(SUM(CASE WHEN created_at >= ? THEN (cex_sell_quote_amount + cex_buy_quote_amount) ELSE 0 END), 0) AS dailyTotalTradeVolume
			FROM strategy_dex_cex_triangular_arbitrage_trades
		`, last24Hours, last24Hours, last24Hours, last24Hours).Row()

	err := row.Scan(
		&model.TotalTrade,
		&model.TotalSuccessTrade,
		&model.DailyTotalTrade,
		&model.DailyTotalSuccessTrade,
		&model.TotalProfit,
		&model.DailyTotalProfit,
		&model.TotalTradeVolume,
		&model.DailyTotalTradeVolume,
	)

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	// 失败套利次数计算（总套利 - 成功套利）
	model.TotalFailedTrade = model.TotalTrade - model.TotalSuccessTrade
	model.DailyTotalFailedTrade = model.DailyTotalTrade - model.DailyTotalSuccessTrade

	// dailyProfitChangePercent
	model.DailyProfitChangePercent = calculateDailyProfitChangePercent(model.DailyTotalProfit, model.TotalProfit)
	model.DailyVolumeChangePercent = calculateDailyProfitChangePercent(model.DailyTotalTradeVolume, model.TotalTradeVolume)

	return nil
}

func calculateDailyProfitChangePercent(dailyProfitStr, totalProfitStr string) string {
	// 将 string 转换为 float64
	dailyProfit, err1 := strconv.ParseFloat(dailyProfitStr, 64)
	totalProfit, err2 := strconv.ParseFloat(totalProfitStr, 64)

	// 检查转换错误
	if err1 != nil || err2 != nil {
		return "error"
	}

	// 计算前一天的总利润
	profitBeforeToday := totalProfit - dailyProfit

	// 处理特殊情况：
	// 1. 总利润和当日利润都为 0（无盈利）
	// 2. 仅有当日利润，之前利润为 0（新开始）
	if totalProfit == 0 {
		return "0"
	}

	// 避免分母为 0，改用绝对值来计算增长率
	if profitBeforeToday == 0 {
		if dailyProfit > 0 {
			return "100"
		} else {
			return "-100"
		}
	}

	changePercent := (dailyProfit / math.Abs(profitBeforeToday)) * 100

	return fmt.Sprintf("%.2f", changePercent)
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

func (e *StrategyDexCexTriangularArbitrageTrades) ScanTrades() error {
	e.Log.Infof("开始扫描订单")
	var data models.StrategyDexCexTriangularArbitrageTrades
	var trades []dto.StrategyDexCexTriangularArbitrageTradesGetPageResp
	// 获取当前时间
	now := time.Now()

	// 计算 1 分钟前的时间
	oneMinuteAgo := now.Add(-time.Minute)

	err := e.Orm.Model(&data).
		Select("strategy_dex_cex_triangular_arbitrage_trades.*, opportunities.cex_target_asset as symbol").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.created_at >= ?", oneMinuteAgo).
		Find(&trades).Error

	if err != nil {
		e.Log.Errorf("获取订单失败，数据库异常:%+v \n", err)
		return err
	}

	if len(trades) == 0 {
		e.Log.Infof("未成交新订单")
		return nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("新成交 %d 笔套利，交易信息如下：\n", len(trades)))

	for i, trade := range trades {
		var buySide string
		if trade.BuyOnDex == "0" {
			buySide = "CEX"
		} else {
			buySide = "DEX"
		}

		cexSellAmount, err1 := strconv.ParseFloat(trade.CexSellQuoteAmount, 64)
		cexBuyAmount, err2 := strconv.ParseFloat(trade.CexBuyQuoteAmount, 64)
		if err1 != nil || err2 != nil {
			return err
		}

		builder.WriteString(fmt.Sprintf(
			"%d. 币对: %s，买方: %s，利润: %.8f USDT，交易时间(UTC): %s\n",
			i+1, trade.Symbol, buySide, cexSellAmount-cexBuyAmount, trade.CreatedAt.Format("2006-01-02 15:04:05"),
		))
	}

	notificationMsg := builder.String()
	config := ext.ExtConfig
	larkClient := lark.NewLarkRobotAlert(config)
	e.Log.Infof("lark notificationMsg:%s \n", notificationMsg)
	err = larkClient.SendLarkAlert(notificationMsg)
	if err != nil {
		e.Log.Infof("lark 推送消息失败")
		return err
	}
	return nil
}
