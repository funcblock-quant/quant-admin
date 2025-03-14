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
		query = query.Where("(strategy_dex_cex_triangular_arbitrage_trades.dex_success = 1 or strategy_dex_cex_triangular_arbitrage_trades.dex_tx_fee is not null) and (strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success = 0 or strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = 0)")
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
		e.Log.Error("lark 推送消息失败")
		return err
	}
	return nil
}

// DailyTradeSnapshot 每日交易快照
func (e *StrategyDexCexTriangularArbitrageTrades) DailyTradeSnapshot() error {
	e.Log.Infof("开始生成每日套利快照")
	yesterday := time.Now().AddDate(0, 0, -1)
	// 获取当天时间范围
	snapshotDate := yesterday.Format("2006-01-02")
	// 获取前一天 00:00:00
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.Local)
	// 获取前一天 23:59:59
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.Local)

	fmt.Println("前一天开始时间:", startOfDay)
	fmt.Println("前一天结束时间:", endOfDay)

	var snapshots []models.BusDexCexDailyTradeStatisticSnapshot
	instanceIdSet := make(map[string]bool)

	// 查询当日成功交易的 instanceId
	var tradedInstanceIds []string
	e.Orm.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("DISTINCT instance_id").
		Where("dex_success = ? AND cex_buy_success = ? AND cex_sell_success =?", 1, 1, 1).
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Pluck("instance_id", &tradedInstanceIds)

	for _, id := range tradedInstanceIds {
		instanceIdSet[id] = true
	}

	var activeInstanceIds []string
	e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Select("id").
		Where("status in ?", []int{INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING}).
		Pluck("id", &activeInstanceIds)

	for _, id := range activeInstanceIds {
		instanceIdSet[id] = true
	}
	// 获取交易开启状态的币种信息
	var instances []struct {
		Id          string
		Symbol      string
		TargetToken string
		QuoteToken  string
	}

	e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Select("id, symbol, target_token, quote_token").
		Where("id IN ?", getKeys(instanceIdSet)).
		Scan(&instances)

	var markdownContent string
	markdownContent += fmt.Sprintf("📊 每日交易快照\n📅 日期：%s\n\n", snapshotDate)
	markdownContent += "        | 币对 | 成交笔数 | 总成交量 | 当天利润 | 利润增长率 |\n"
	markdownContent += "        |------|--------|---------|---------|---------|\n"

	var allTrades int
	var allVolume, allProfit, allPreviousProfit, allProfitGrowthRate float64
	// 统计每个币种的总成交量和总收益
	for _, inst := range instances {
		var totalTrade int
		var totalVolume, totalProfit, previousTotalProfit, profitGrowthRate float64

		e.Orm.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
			Select("COUNT(*) AS totalTrade, COALESCE(SUM(cex_sell_quote_amount), 0) AS total_volume, COALESCE(SUM(cex_sell_quote_amount - cex_buy_quote_amount), 0) AS total_profit").
			Where("instance_id = ?", inst.Id).
			Where("dex_success = ? AND cex_buy_success = ? AND cex_sell_success =?", 1, 1, 1).
			Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
			Row().Scan(&totalTrade, &totalVolume, &totalProfit)

		// 查询今天之前的所有利润
		e.Orm.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
			Select("COALESCE(SUM(cex_sell_quote_amount - cex_buy_quote_amount), 0) AS previous_total_profit").
			Where("instance_id = ?", inst.Id).
			Where("dex_success = ? AND cex_buy_success = ? AND cex_sell_success =?", 1, 1, 1).
			Where("created_at < ?", startOfDay).
			Row().Scan(&previousTotalProfit)

		e.Log.Infof("previous_total_profit : %d", previousTotalProfit)
		// 计算当天利润增长百分比
		if previousTotalProfit > 0 {
			profitGrowthRate = (totalProfit - previousTotalProfit) / previousTotalProfit
		} else {
			profitGrowthRate = 0
		}

		// 组装快照数据（无成交的数据 TotalVolume 和 TotalProfit 仍然为 0）
		snapshots = append(snapshots, models.BusDexCexDailyTradeStatisticSnapshot{
			InstanceID:   inst.Id,
			SnapshotDate: snapshotDate,
			Symbol:       inst.Symbol,
			TargetToken:  inst.TargetToken,
			QuoteToken:   inst.QuoteToken,
			TotalVolume:  totalVolume,
			TotalProfit:  totalProfit,
		})

		allTrades += totalTrade
		allVolume += totalVolume
		allProfit += totalProfit
		allPreviousProfit += previousTotalProfit

		// 拼接lark通知消息
		markdownContent += fmt.Sprintf("        | %s/%s | %d | %.2f | $%.2f | %.2f%% |\n",
			inst.TargetToken, inst.QuoteToken, totalTrade, totalVolume,
			totalProfit, profitGrowthRate*100)
	}

	if allPreviousProfit > 0 {
		allProfitGrowthRate = (allProfit - allPreviousProfit) / allPreviousProfit
	} else {
		allProfitGrowthRate = 0
	}

	markdownContent += "        |------|--------|---------|---------|---------|\n"
	markdownContent += fmt.Sprintf("汇总 |   x   | %d | %.2f | $%.2f | %.2f%% |\n",
		allTrades, allVolume,
		allProfit, allProfitGrowthRate*100)

	config := ext.ExtConfig
	larkClient := lark.NewLarkRobotAlert(config)
	e.Log.Infof("lark notificationMsg:%s \n", markdownContent)
	err := larkClient.SendLarkAlert(markdownContent)
	if err != nil {
		e.Log.Infof("lark 推送消息失败")
	}

	return e.Orm.Create(&snapshots).Error
}

// getKeys 获取 map 的 key 列表
func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
