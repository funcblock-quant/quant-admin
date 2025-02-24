package service

import (
	"errors"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	"strconv"
)

type BusDexCexTriangularObserver struct {
	service.Service
}

// GetPage 获取BusDexCexTriangularObserver列表
func (e *BusDexCexTriangularObserver) GetPage(c *dto.BusDexCexTriangularObserverGetPageReq, p *actions.DataPermission, list *[]dto.BusDexCexTriangularObserverGetPageResp, count *int64) error {
	var err error
	var data models.BusDexCexTriangularObserver
	e.Log.Infof("e[GetPage], data: %+v", data)
	tx := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	statuses := []int{1, 2, 3} // 需要查询的多个状态
	tx.Where("status IN (?)", statuses)
	err = tx.Debug().Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetPage error:%s \r\n", err)
		return err
	}

	if len(*list) == 0 {
		return nil
	}

	if config.ApplicationConfig.Mode == "dev" {
		// dev环境不调用grpc
		return nil
	}

	for i := range *list {
		observerId := (*list)[i].InstanceId // 使用 (*list)[i] 访问原始元素
		id := (*list)[i].Id
		state, err := client.GetObserverState(observerId)
		if err != nil {
			e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
			continue
		}
		e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", observerId, state)
		buyOnDex := state.GetBuyOnDex()
		cexSellPrice, dexBuyPrice := e.calculate_dex_cex_price(buyOnDex, true)
		e.Log.Infof("[buy on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexSellPrice, dexBuyPrice)

		sellOnDex := state.GetSellOnDex()
		cexBuyPrice, dexSellPrice := e.calculate_dex_cex_price(sellOnDex, false)
		e.Log.Infof("[sell on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexBuyPrice, dexSellPrice)

		if cexSellPrice-dexBuyPrice > 0 {
			//获取最新的价差记录统计信息，设置价差持续时间
			dexBuyData := models.BusDexCexPriceSpreadStatistics{}
			err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					(*list)[i].DexBuyDiffDuration = "0"
				} else {
					e.Log.Errorf("db error:%s", err)
					return err
				}
			}
			(*list)[i].DexBuyDiffDuration = dexBuyData.Duration
		}
		if dexSellPrice-cexBuyPrice > 0 {
			dexSellData := models.BusDexCexPriceSpreadStatistics{}
			err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					(*list)[i].DexBuyDiffDuration = "0"
				} else {
					e.Log.Errorf("db error:%s", err)
					return err
				}
			}
			(*list)[i].DexSellDiffDuration = dexSellData.Duration
		}

		buyOnDexProfit := *buyOnDex.CexSellQuoteAmount - *buyOnDex.CexBuyQuoteAmount
		sellOnDexProfit := *sellOnDex.CexSellQuoteAmount - *sellOnDex.CexBuyQuoteAmount

		//(*list)[i].ProfitOfBuyOnDex = strconv.FormatFloat(buyOnDexProfit, 'f', 6, 64)
		//(*list)[i].ProfitOfSellOnDex = strconv.FormatFloat(sellOnDexProfit, 'f', 6, 64)
		//(*list)[i].CexSellPrice = strconv.FormatFloat(cexSellPrice, 'f', 6, 64)
		//(*list)[i].DexBuyPrice = strconv.FormatFloat(dexBuyPrice, 'f', 6, 64)
		//(*list)[i].DexBuyDiffPrice = strconv.FormatFloat(cexSellPrice-dexBuyPrice, 'f', 6, 64)
		//(*list)[i].CexBuyPrice = strconv.FormatFloat(cexBuyPrice, 'f', 6, 64)
		//(*list)[i].DexSellPrice = strconv.FormatFloat(dexSellPrice, 'f', 6, 64)
		//(*list)[i].DexSellDiffPrice = strconv.FormatFloat(dexSellPrice-cexBuyPrice, 'f', 6, 64)
		(*list)[i].ProfitOfBuyOnDex = buyOnDexProfit
		(*list)[i].ProfitOfSellOnDex = sellOnDexProfit
		(*list)[i].CexSellPrice = cexSellPrice
		(*list)[i].DexBuyPrice = dexBuyPrice
		(*list)[i].DexBuyDiffPrice = cexSellPrice - dexBuyPrice
		(*list)[i].CexBuyPrice = cexBuyPrice
		(*list)[i].DexSellPrice = dexSellPrice
		(*list)[i].DexSellDiffPrice = dexSellPrice - cexBuyPrice
	}

	return nil
}

func (e *BusDexCexTriangularObserver) calculate_dex_cex_price(priceState *pb.ObserverState, isDexBuy bool) (float64, float64) {
	var cexPrice float64      // TRUMP/USDT
	var dexPrice float64      //TRUMP/USDT
	var cexQuotePrice float64 // 例如：TRUMP/USDT
	var cexSolPrice float64   //SOL/USDT
	if isDexBuy {
		// dex买入
		if priceState.CexSellQuantity != nil && priceState.CexSellQuoteAmount != nil && *priceState.CexSellQuantity != 0 {
			cexQuotePrice = *priceState.CexSellQuoteAmount / *priceState.CexSellQuantity
		} else {
			// 处理 nil 或除数为 0 的情况，避免 panic
			cexQuotePrice = 0
		}

		if priceState.CexBuyQuantity != nil && priceState.CexBuyQuoteAmount != nil && *priceState.CexBuyQuantity != 0 {
			cexSolPrice = *priceState.CexBuyQuoteAmount / *priceState.CexBuyQuantity
		} else {
			// 处理 nil 或除数为 0 的情况，避免 panic
			cexSolPrice = 0
		}
	} else {
		// dex卖出
		if priceState.CexBuyQuantity != nil && priceState.CexBuyQuoteAmount != nil && *priceState.CexBuyQuantity != 0 {
			cexQuotePrice = *priceState.CexBuyQuoteAmount / *priceState.CexBuyQuantity
		} else {
			// 处理 nil 或除数为 0 的情况，避免 panic
			cexQuotePrice = 0
		}

		if priceState.CexSellQuantity != nil && priceState.CexSellQuoteAmount != nil && *priceState.CexSellQuantity != 0 {
			cexSolPrice = *priceState.CexSellQuoteAmount / *priceState.CexSellQuantity
		} else {
			// 处理 nil 或除数为 0 的情况，避免 panic
			cexSolPrice = 0
		}
	}

	cexPrice = cexQuotePrice

	var dexSolPrice float64 //TRUMP/WSOL
	if priceState.DexSolAmount != nil && priceState.DexTargetAmount != nil && *priceState.DexTargetAmount != 0 {
		dexSolPrice = *priceState.DexSolAmount / *priceState.DexTargetAmount
	} else {
		// 处理 nil 或除数为 0 的情况，避免 panic
		dexSolPrice = 0
	}

	if cexSolPrice != 0 && dexSolPrice != 0 {
		dexPrice = dexSolPrice * cexSolPrice
	} else {
		dexPrice = 0
	}

	return cexPrice, dexPrice
}

// Get 获取BusDexCexTriangularObserver对象
func (e *BusDexCexTriangularObserver) Get(d *dto.BusDexCexTriangularObserverGetReq, p *actions.DataPermission, model *dto.BusDexCexTriangularObserverDetailResp) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexTriangularObserver error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if config.ApplicationConfig.Mode == "dev" {
		// dev环境不调用grpc
		return nil
	}

	// 获取最新价差数据
	observerId := model.InstanceId
	id := model.Id
	state, err := client.GetObserverState(observerId)
	if err != nil {
		e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
		return nil
	}
	e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", observerId, state)
	buyOnDex := state.GetBuyOnDex()
	cexSellPrice, dexBuyPrice := e.calculate_dex_cex_price(buyOnDex, true)
	e.Log.Infof("[buy on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexSellPrice, dexBuyPrice)

	sellOnDex := state.GetSellOnDex()
	cexBuyPrice, dexSellPrice := e.calculate_dex_cex_price(sellOnDex, false)
	e.Log.Infof("[sell on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexBuyPrice, dexSellPrice)

	buyOnDexProfit := *buyOnDex.CexSellQuoteAmount - *buyOnDex.CexBuyQuoteAmount
	sellOnDexProfit := *sellOnDex.CexSellQuoteAmount - *sellOnDex.CexBuyQuoteAmount

	model.ProfitOfBuyOnDex = buyOnDexProfit
	model.ProfitOfSellOnDex = sellOnDexProfit
	model.CexSellPrice = cexSellPrice
	model.DexBuyPrice = dexBuyPrice
	model.DexBuyDiffPrice = cexSellPrice - dexBuyPrice
	if cexSellPrice-dexBuyPrice > 0 {
		//获取最新的价差记录统计信息，设置价差持续时间
		dexBuyData := models.BusDexCexPriceSpreadStatistics{}
		err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			model.DexBuyDiffDuration = "0"
		}
		if err != nil {
			e.Log.Errorf("db error:%s", err)
			return err
		}
		model.DexBuyDiffDuration = dexBuyData.Duration
	}
	if dexSellPrice-cexBuyPrice > 0 {
		dexSellData := models.BusDexCexPriceSpreadStatistics{}
		err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			model.DexBuyDiffDuration = "0"
		}
		if err != nil {
			e.Log.Errorf("db error:%s", err)
			return err
		}
		model.DexSellDiffDuration = dexSellData.Duration
	}
	model.CexBuyPrice = cexBuyPrice
	model.DexSellPrice = dexSellPrice
	model.DexSellDiffPrice = dexSellPrice - cexBuyPrice

	return nil
}

// GetSymbolList 获取BusDexCexTriangularObserver所有币种列表
func (e *BusDexCexTriangularObserver) GetSymbolList(p *actions.DataPermission, list *[]dto.DexCexTriangularObserverSymbolListResp) error {
	var err error
	var data models.BusDexCexTriangularObserver

	err = e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Select("symbol").
		Group("symbol").
		Where("status = ?", 1).
		Debug().Find(list).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexTriangularObserver对象
func (e *BusDexCexTriangularObserver) Insert(c *dto.BusDexCexTriangularObserverInsertReq) error {
	var err error
	var data models.BusDexCexTriangularObserver
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// BatchInsert 创建BusDexCexTriangularObserver对象
func (e *BusDexCexTriangularObserver) BatchInsert(c *dto.BusDexCexTriangularObserverBatchInsertReq) error {
	var data models.BusDexCexTriangularObserver
	var err error
	baseTokens := c.TargetToken
	if len(baseTokens) == 0 {
		return errors.New("empty baseTokens")
	}

	var successStartedCount int
	for _, baseToken := range baseTokens {
		//循环创建监听

		c.Generate(&data, baseToken)
		tx := e.Orm.Begin()
		err = tx.Create(&data).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService Insert error:%s \r\n", err)
			continue
		}

		var ammConfig = pb.DexConfig{}
		var amberConfig = pb.AmberObserverConfig{}
		var arbitrageConfig = pb.ObserverParams{}
		c.GenerateAmmConfig(&ammConfig)
		c.GenerateAmberConfig(&amberConfig)
		c.GenerateObserverParams(&arbitrageConfig)
		if config.ApplicationConfig.Mode == "dev" {
			// dev环境不调用grpc
		} else {
			instanceId := strconv.Itoa(data.Id)
			err = client.StartNewArbitragerClient(&instanceId, &amberConfig, &ammConfig, &arbitrageConfig)
			if err != nil {
				e.Log.Errorf("Service BatchInsert error:%s \r\n", err)
				tx.Rollback()
				continue
			}
		}

		tx.Where("id = ?", data.Id).Update("status", 1)
		err = tx.Commit().Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService Insert error:%s \r\n", err)
			continue
		}
		successStartedCount += 1
	}

	if successStartedCount == 0 {
		return errors.New("创建失败")
	}
	return nil
}

// Update 修改BusDexCexTriangularObserver对象
func (e *BusDexCexTriangularObserver) Update(c *dto.BusDexCexTriangularObserverUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusDexCexTriangularObserver{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusDexCexTriangularObserver
func (e *BusDexCexTriangularObserver) Remove(d *dto.BusDexCexTriangularObserverDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexTriangularObserver

	instanceId := strconv.Itoa(d.Ids)
	err := client.StopArbitragerClient(instanceId)
	if err != nil {
		e.Log.Errorf("暂停监视器失败 error:%s \r\n", err)
		return err
	}

	e.Log.Infof("grpc请求暂停监视器成功")

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())

	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusDexCexTriangularObserver error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// StartTrader 启动交易功能
func (e *BusDexCexTriangularObserver) StartTrader(c *dto.BusDexCexTriangularObserverStartTraderReq) error {
	var data models.BusDexCexTriangularObserver
	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}
	if data.IsTrading {
		e.Log.Infof("实例：%s 交易功能已开启，跳过grpc调用", data.InstanceId)
		return nil
	}

	// step1 先启动水位调节实例
	tokenConfig := &waterLevelPb.TokenConfig{
		Currency:               data.TargetToken,
		PubKey:                 *data.TokenMint,
		OwnerProgram:           *data.OwnerProgram,
		Decimals:               uint32(data.Decimals),
		AlertThreshold:         fmt.Sprintf("%v", c.AlertThreshold),
		BuyTriggerThreshold:    fmt.Sprintf("%v", c.BuyTriggerThreshold),
		TargetBalanceThreshold: fmt.Sprintf("%v", c.TargetBalanceThreshold),
		SellTriggerThreshold:   fmt.Sprintf("%v", c.SellTriggerThreshold),
	}

	clientRequest := &waterLevelPb.StartInstanceRequest{
		InstanceId:   strconv.Itoa(data.Id),
		ExchangeType: data.ExchangeType,
		TokenConfig:  tokenConfig,
	}

	_, err = client.StartWaterLevelInstance(clientRequest)
	if err != nil {
		e.Log.Errorf("启动水位调节失败:%s \r\n", err)
		return err
	}
	e.Log.Infof("水位调节启动成功")

	//exchangeType, err := data.GetExchangeTypeForTrader()
	//if err != nil {
	//	e.Log.Errorf("获取ExchangeType参数异常，:%s \r\n", err)
	//	return err
	//}

	slippageBpsUint, err := strconv.ParseUint(*c.SlippageBps, 10, 32)
	if err != nil {
		e.Log.Errorf("slippageBps: %v\n", slippageBpsUint)
		return errors.New("error slippageBps")
	}
	log.Infof("slippageBps: %v\n", slippageBpsUint)

	//priorityFee := uint64(*c.PriorityFee * 1_000_000_000)
	//jitoFee := uint64(*c.JitoFee * 1_000_000_000)

	//err = requestStartTrader(c, priorityFee, jitoFee, slippageBpsUint, data, exchangeType, e)
	//if err != nil {
	//	return err
	//}

	// 启动水位调节后，更新数据库中的相关参数
	updateData := map[string]interface{}{
		"is_trading":               false,
		"alert_threshold":          c.AlertThreshold,
		"buy_trigger_threshold":    c.BuyTriggerThreshold,
		"target_balance_threshold": c.TargetBalanceThreshold,
		"sell_trigger_threshold":   c.SellTriggerThreshold,
		"slippage_bps":             slippageBpsUint,
		"status":                   2, // 水位调节中
	}

	if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", data.InstanceId).
		Updates(updateData).Error; err != nil {
		e.Log.Errorf("更新实例参数失败：%s", err)
		return err
	}

	e.Log.Infof("实例：%s 参数已成功更新", data.InstanceId)
	return nil
}

// MonitorWaterLevelToStartTrader 监控水位以启动交易
func (e *BusDexCexTriangularObserver) MonitorWaterLevelToStartTrader() error {
	var data models.BusDexCexTriangularObserver
	var list []models.BusDexCexTriangularObserver
	err := e.Orm.Model(&data).
		//status = 2 表示水位调节中，并且交易功能暂停
		Where("status = ? and is_trading = ?", 2, false).
		Find(&list).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	for _, instance := range list {
		instanceId := instance.Id
		queryReq := &waterLevelPb.InstantId{
			InstanceId: strconv.Itoa(instanceId),
		}
		waterLevelState, err := client.GetWaterLevelInstanceState(queryReq)
		if err != nil {
			e.Log.Errorf("get WaterLevelInstanceState error: instanceId:%d, error msg:%s \r\n", instance.Id, err)
			continue
		}

		traderSwitch := waterLevelState.TraderSwitch
		if traderSwitch {
			e.Log.Infof("waterlevel state for instancId: %d is: success", instanceId)
			e.Log.Infof("currency: %s, cex balance:%s,  dex balance: %s", waterLevelState.Currency, waterLevelState.CexAccountBalance, waterLevelState.ChainWalletBalance)
			//开启交易功能

			err = requestStartTrader(&instance, e)
			if err != nil {
				e.Log.Errorf("start trader error:%s \r\n", err)
				return err
			}

			// 启动成功后，更新状态
			updateData := map[string]interface{}{
				"is_trading": true,
				"status":     3, // 交易已启动
			}
			if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
				Where("id = ?", instance.Id).
				Updates(updateData).Error; err != nil {
				e.Log.Errorf("更新实例参数失败：%s", err)
				continue
			}
			e.Log.Infof("实例：%s 参数已成功更新", data.InstanceId)
		}
	}
	return nil
}

func requestStartTrader(instance *models.BusDexCexTriangularObserver, e *BusDexCexTriangularObserver) error {

	exchangeType, err := instance.GetExchangeTypeForTrader()
	if err != nil {
		e.Log.Errorf("获取ExchangeType参数异常，:%s \r\n", err)
		return err
	}

	slippageBpsUint, err := strconv.ParseUint(instance.SlippageBps, 10, 64)
	if err != nil {
		e.Log.Errorf("转换失败: %s", err)
	}
	e.Log.Infof("slippageBps: %v\n", slippageBpsUint)

	amberTraderConfig := &pb.AmberTraderConfig{
		ExchangeType: &exchangeType,
	}
	traderParams := &pb.TraderParams{
		SlippageBps: &slippageBpsUint,
	}
	if config.ApplicationConfig.Mode != "dev" {
		instanceId := strconv.Itoa(instance.Id)
		err := client.EnableTrader(instanceId, amberTraderConfig, traderParams)
		if err != nil {
			e.Log.Errorf("GRPC 启动Trader for instanceId:%s 失败，异常:%s \r\n", instance.InstanceId, err)
			return err
		}
	}
	return nil
}

// MonitorWaterLevelToStopTrader 监控水位以暂停交易
func (e *BusDexCexTriangularObserver) MonitorWaterLevelToStopTrader() error {
	var data models.BusDexCexTriangularObserver
	var list []models.BusDexCexTriangularObserver
	err := e.Orm.Model(&data).
		//status = 3 表示已启动交易，并且交易功能开启中
		Where("status = ? and is_trading = ?", 3, true).
		Find(&list).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	for _, instance := range list {
		instanceId := strconv.Itoa(instance.Id)
		queryReq := &waterLevelPb.InstantId{
			InstanceId: instanceId,
		}
		waterLevelState, err := client.GetWaterLevelInstanceState(queryReq)
		if err != nil {
			e.Log.Errorf("get WaterLevelInstanceState error: instanceId:%d, error msg:%s \r\n", instance.Id, err)
			continue
		}

		traderSwitch := waterLevelState.TraderSwitch
		e.Log.Infof("currency: %s, cex balance:%s,  dex balance: %s", waterLevelState.Currency, waterLevelState.CexAccountBalance, waterLevelState.ChainWalletBalance)
		if !traderSwitch {
			e.Log.Infof("waterlevel state for instancId: %d is: failed", instanceId)
			//关闭交易功能

			err = client.DisableTrader(instanceId)
			if err != nil {
				e.Log.Errorf("grpc暂停实例：:%s 交易功能失败，异常：%s \r\n", instanceId, err)
				return err
			}

			// 暂停交易成功后，更新状态
			updateData := map[string]interface{}{
				"is_trading": false,
				"status":     2, // 水位调节中
			}
			if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
				Where("id = ?", instance.Id).
				Updates(updateData).Error; err != nil {
				e.Log.Errorf("更新实例参数失败：%s", err)
				continue
			}
			e.Log.Infof("实例：%s 参数已成功更新", data.InstanceId)
		}
	}
	return nil
}

// StopTrader 停止交易功能
func (e *BusDexCexTriangularObserver) StopTrader(c *dto.BusDexCexTriangularObserverStopTraderReq) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}
	if !data.IsTrading {
		e.Log.Infof("实例：%s 交易功能未开启，跳过grpc调用", data.InstanceId)
		return nil
	}
	if config.ApplicationConfig.Mode != "dev" {
		err = client.DisableTrader(data.InstanceId)
		if err != nil {
			e.Log.Errorf("grpc暂停实例：:%s 交易功能失败，异常：%s \r\n", data.InstanceId, err)
			return err
		}

	}

	// 更新observer的isTrading = false
	err = e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		Update("is_trading", false).Error
	if err != nil {
		e.Log.Errorf("更新数据库实例:%s 交易状态失败，异常信息：%s \r\n", data.InstanceId, err)
		return err
	}

	return nil
}

// UpdateObserver 更新observer 参数
func (e *BusDexCexTriangularObserver) UpdateObserver(c *dto.BusDexCexTriangularUpdateObserverParamsReq) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	priorityFee := uint64(*c.PriorityFee * 1_000_000_000)
	jitoFee := uint64(*c.JitoFee * 1_000_000_000)
	transactionFee := &pb.TransactionFee{
		PriorityFee: &priorityFee,
		JitoFee:     &jitoFee,
	}

	observerParams := &pb.ObserverParams{
		MinSolAmount:             c.MinSolAmount,
		MaxSolAmount:             c.MaxSolAmount,
		TriggerProfitQuoteAmount: c.TriggerProfitQuoteAmount,
		TxFee:                    transactionFee,
	}
	if config.ApplicationConfig.Mode != "dev" {
		// dev环境不调用grpc
		instanceId := strconv.Itoa(data.Id)
		err = client.UpdateObserverParams(instanceId, observerParams)
		if err != nil {
			e.Log.Errorf("grpc更新实例：:%s Observer参数失败，异常：%s \r\n", data.InstanceId, err)
			return err
		}
	}

	// 更新observer的参数
	updateData := map[string]interface{}{
		"min_sol_amount": c.MinSolAmount,
		"max_sol_amount": c.MaxSolAmount,
		"min_profit":     c.TriggerProfitQuoteAmount,
		"priority_fee":   priorityFee,
		"jito_fee":       jitoFee,
	}
	err = e.Orm.Model(&data).
		Where("id = ?", data.Id).
		Updates(updateData).Error
	if err != nil {
		e.Log.Errorf("更新数据库实例:%d 交易状态失败，异常信息：%s \r\n", data.Id, err)
		return err
	}

	return nil
}

// UpdateTrader 更新Trader 参数
func (e *BusDexCexTriangularObserver) UpdateTrader(c *dto.BusDexCexTriangularUpdateTraderParamsReq) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}
	//priorityFee := uint64(*c.PriorityFee * 1_000_000_000)
	//jitoFee := uint64(*c.JitoFee * 1_000_000_000)
	slippageBpsUint, err := strconv.ParseUint(*c.SlippageBps, 10, 32)
	if err != nil {
		e.Log.Errorf("slippageBps: %v\n", slippageBpsUint)
		return errors.New("error slippageBps")
	}
	log.Infof("slippageBps: %v\n", slippageBpsUint)
	//txBuildParams := &pb.TransactionFee{
	//	PriorityFee: &priorityFee,
	//	JitoFee:     &jitoFee,
	//}
	traderParams := &pb.TraderParams{
		SlippageBps: &slippageBpsUint,
	}
	if config.ApplicationConfig.Mode != "dev" {
		err = client.UpdateTraderParams(data.InstanceId, traderParams)
		if err != nil {
			e.Log.Errorf("grpc更新实例：:%s trader参数失败，异常：%s \r\n", data.InstanceId, err)
			return err
		}
	}

	updateData := map[string]interface{}{
		"slippage_bps": slippageBpsUint,
	}
	// 更新observer的trader相关参数
	if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", data.Id).
		Updates(updateData).Error; err != nil {
		e.Log.Errorf("更新实例参数失败：%s", err)
		return err
	}

	e.Log.Infof("实例：%s 参数已成功更新", data.InstanceId)

	return nil
}
