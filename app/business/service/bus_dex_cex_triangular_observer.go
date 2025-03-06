package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	common "quanta-admin/common/models"
	lark "quanta-admin/common/notification"
	ext "quanta-admin/config"
	"sort"
	"strconv"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
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
		id := (*list)[i].Id
		state, err := client.GetObserverState(strconv.Itoa(id))
		if err != nil {
			e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
			continue
		}
		e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", strconv.Itoa(id), state)
		buyOnDex := state.GetBuyOnDex()
		var cexSellPrice, dexBuyPrice, buyOnDexProfit float64
		if buyOnDex != nil {
			cexSellPrice, dexBuyPrice = e.calculate_dex_cex_price(buyOnDex, true)
			buyOnDexProfit = *buyOnDex.CexSellQuoteAmount - *buyOnDex.CexBuyQuoteAmount
		} else {
			// 处理 buyOnDex 为空的情况，例如设置默认值或跳过计算
			cexSellPrice = 0
			dexBuyPrice = 0
			buyOnDexProfit = 0
		}

		sellOnDex := state.GetSellOnDex()
		var cexBuyPrice, dexSellPrice, sellOnDexProfit float64
		if sellOnDex != nil {
			cexBuyPrice, dexSellPrice = e.calculate_dex_cex_price(sellOnDex, false)
			sellOnDexProfit = *sellOnDex.CexSellQuoteAmount - *sellOnDex.CexBuyQuoteAmount
		} else {
			// 处理 sellOnDex 为空的情况，例如设置默认值或跳过计算
			cexBuyPrice = 0
			dexSellPrice = 0
			sellOnDexProfit = 0
		}

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
	id := strconv.Itoa(model.Id)
	state, err := client.GetObserverState(id)
	if err != nil {
		e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
		return nil
	}
	e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", id, state)
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
		Unscoped().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Select("symbol").
		Group("symbol").
		Where("status IN ?", []int{1, 2, 3}).
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

		tx.Model(&data).Where("id = ?", data.Id).Update("status", 1)
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

	// 判断是否被风控
	if data.IsTradingBlocked {
		e.Log.Infof("实例：%s 已被风控，跳过grpc调用", data.InstanceId)
		return errors.New("实例已被风控,请前往风控中心手动处理或等风控自动解除")
	}

	data.AlertThreshold = c.AlertThreshold
	data.BuyTriggerThreshold = c.BuyTriggerThreshold
	data.SellTriggerThreshold = c.SellTriggerThreshold
	data.MinDepositAmountThreshold = c.MinDepositAmountThreshold
	data.MinWithdrawAmountThreshold = c.MinWithdrawAmountThreshold
	// step1 先启动水位调节实例
	err = StartTokenWaterLevelWithCheckExists(&data)
	if err != nil {
		return err
	}

	//exchangeType, err := data.GetExchangeTypeForTrader()
	//if err != nil {
	//	e.Log.Errorf("获取ExchangeType参数异常，:%s \r\n", err)
	//	return err
	//}

	//slippageBpsUint, err := strconv.ParseUint(*c.SlippageBps, 10, 32)
	//if err != nil {
	//	e.Log.Errorf("slippageBps: %v\n", slippageBpsUint)
	//	return errors.New("error slippageBps")
	//}
	//log.Infof("slippageBps: %v\n", slippageBpsUint)

	priorityFeeRate := *c.PriorityFeeRate
	//err = StartTrader(c, priorityFee, jitoFee, slippageBpsUint, data, exchangeType, e)
	//if err != nil {
	//	return err
	//}

	// 启动水位调节后，更新数据库中的相关参数
	updateData := map[string]interface{}{
		"is_trading":             false,
		"alert_threshold":        c.AlertThreshold,
		"buy_trigger_threshold":  c.BuyTriggerThreshold,
		"sell_trigger_threshold": c.SellTriggerThreshold,
		//"slippage_bps":           slippageBpsUint,
		"priority_fee_rate": priorityFeeRate,
		"jito_fee_rate":     c.JitoFeeRate,
		"status":            2, // 水位调节中
	}

	if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", data.Id).
		Updates(updateData).Error; err != nil {
		e.Log.Errorf("更新实例参数失败：%s", err)
		return err
	}

	e.Log.Infof("实例：%s 参数已成功更新", data.Id)
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
			e.Log.Infof("currency: %s, cex spot balance:%s, cex margin balance:%s, dex balance: %s", waterLevelState.Currency, waterLevelState.SpotAccountBalance, waterLevelState.MarginAccountBalance, waterLevelState.ChainWalletBalance)
			//开启交易功能

			err = StartTrader(&instance)
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
		e.Log.Infof("currency: %s, cex spot balance:%s, cex margin balance:%s, dex balance: %s", waterLevelState.Currency, waterLevelState.SpotAccountBalance, waterLevelState.MarginAccountBalance, waterLevelState.ChainWalletBalance)
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
func (e *BusDexCexTriangularObserver) StopTrader(c *dto.BusDexCexTriangularObserverStopTraderReq, isTradingBlocked bool) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}
	if !data.IsTrading {
		e.Log.Infof("实例：%s 交易功能未开启，跳过grpc调用", c.InstanceId)
		return nil
	}
	instanceId := strconv.Itoa(c.InstanceId)
	if config.ApplicationConfig.Mode != "dev" {
		err = client.DisableTrader(instanceId)
		if err != nil {
			e.Log.Errorf("grpc暂停实例：:%d 交易功能失败，异常：%s \r\n", c.InstanceId, err)
			return err
		}

		stopReq := &waterLevelPb.InstantId{
			InstanceId: instanceId,
		}
		err = client.StopWaterLevelInstance(stopReq)
		if err != nil {
			e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", c.InstanceId, err)
			return err
		}
	}

	// 更新observer的isTrading = false
	updateData := map[string]interface{}{
		"is_trading":         false,
		"status":             1,
		"is_trading_blocked": isTradingBlocked,
	}

	err = e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		Updates(updateData).Error
	if err != nil {
		e.Log.Errorf("更新数据库实例:%d 交易状态失败，异常信息：%s \r\n", c.InstanceId, err)
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

	//triggerHoldingMsUint := uint64(c.TriggerHoldingMs)
	observerParams := &pb.ObserverParams{
		MinQuoteAmount: c.MinQuoteAmount,
		MaxQuoteAmount: c.MaxQuoteAmount,
		//SlippageRate:      c.SlippageBpsRate,
		ProfitTriggerRate: c.ProfitTriggerRate,
		//TriggerHoldingMs:  &triggerHoldingMsUint,
	}
	if config.ApplicationConfig.Mode != "dev" {
		// dev环境不调用grpc
		instanceId := strconv.Itoa(data.Id)
		err = client.UpdateObserverParams(instanceId, observerParams)
		if err != nil {
			e.Log.Errorf("grpc更新实例：:%d Observer参数失败，异常：%s \r\n", data.Id, err)
			return err
		}
	}

	// 更新observer的参数
	updateData := map[string]interface{}{
		"min_quote_amount":    c.MinQuoteAmount,
		"max_quote_amount":    c.MaxQuoteAmount,
		"profit_trigger_rate": c.ProfitTriggerRate,
		"trigger_holding_ms":  c.TriggerHoldingMs,
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
	//slippageBpsUint, err := strconv.ParseUint(*c.SlippageBps, 10, 64)
	//if err != nil {
	//	e.Log.Errorf("转换失败: %s", err)
	//}
	//e.Log.Infof("slippageBps: %v\n", slippageBpsUint)
	//slippageBpsFloat := float64(slippageBpsUint) / 10000.0

	traderParams := &pb.TraderParams{
		//Slippage:    &slippageBpsFloat,
		SlippageRate:    c.SlippageBpsRate,
		PriorityFeeRate: c.PriorityFeeRate,
		JitoFeeRate:     c.JitoFeeRate,
	}
	if config.ApplicationConfig.Mode != "dev" {
		instanceId := strconv.Itoa(data.Id)
		err = client.UpdateTraderParams(instanceId, traderParams)
		if err != nil {
			e.Log.Errorf("grpc更新实例：:%s trader参数失败，异常：%s \r\n", instanceId, err)
			return err
		}
	}

	updateData := map[string]interface{}{
		"slippage_bps_rate": c.SlippageBpsRate,
		"priority_fee_rate": c.PriorityFeeRate,
		"jito_fee_rate":     *c.JitoFeeRate,
	}
	// 更新observer的trader相关参数
	if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", data.Id).
		Updates(updateData).Error; err != nil {
		e.Log.Errorf("更新实例参数失败：%s", err)
		return err
	}

	e.Log.Infof("实例：%d 参数已成功更新", data.Id)

	return nil
}

// UpdateWaterLevel 更新WaterLevel 参数
func (e *BusDexCexTriangularObserver) UpdateWaterLevel(c *dto.BusDexCexTriangularUpdateWaterLevelParamsReq) error {
	var data models.BusDexCexTriangularObserver

	err := e.Orm.Model(&data).
		Where("id = ?", c.InstanceId).
		First(&data).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	data.BuyTriggerThreshold = c.BuyTriggerThreshold
	data.SellTriggerThreshold = c.SellTriggerThreshold
	data.AlertThreshold = c.AlertThreshold
	data.MinDepositAmountThreshold = c.MinDepositAmountThreshold
	data.MinWithdrawAmountThreshold = c.MinWithdrawAmountThreshold
	err = UpdateTokenWaterLevel(&data)
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"alert_threshold":               c.AlertThreshold,
		"buy_trigger_threshold":         c.BuyTriggerThreshold,
		"sell_trigger_threshold":        c.SellTriggerThreshold,
		"min_deposit_amount_threshold":  c.MinDepositAmountThreshold,
		"min_withdraw_amount_threshold": c.MinWithdrawAmountThreshold,
	}
	// 更新observer的trader相关参数
	if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", data.Id).
		Updates(updateData).Error; err != nil {
		e.Log.Errorf("更新实例参数失败：%s", err)
		return err
	}

	e.Log.Infof("实例：%d 参数已成功更新", data.Id)

	return nil
}

// GetGlobalWaterLevelState 获取全局WaterLevel 状态及当前参数
func (e *BusDexCexTriangularObserver) GetGlobalWaterLevelState() (*dto.BusDexCexTriangularGlobalWaterLevelStateResp, error) {
	var data models.BusCommonConfig
	resp := &dto.BusDexCexTriangularGlobalWaterLevelStateResp{}
	//step 1 : 查询sol以及稳定币的水位调节是否启动
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		//return nil, errors.New("水位调节服务不可用，请稍后刷新重试")
	} else {
		log.Infof("waterLevelInstances:%+v\n", waterLevelInstances)
		for _, instanceId := range waterLevelInstances.InstanceIds {
			if instanceId == "SOLANA" {
				// solana 已经启动水位调节
				resp.SolWaterLevelState = true
			} else if instanceId == "USDT" {
				// 稳定币 已经启动水位调节
				resp.StableCoinWaterLevelState = true
			}
		}
	}

	// 如果水位调节服务挂了，返回对应的错误，给到前端水位调节不可用的提示之类的。
	//step 2 : 封装全局的水位调节启动结果以及配置到响应体
	err = e.Orm.Model(&models.BusCommonConfig{}).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&data).Error

	resp.SolWaterLevelConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{}
	if err != nil {
		e.Log.Errorf("获取solana水位调节参数失败")
	} else {
		configJsonStr := data.ConfigJson
		e.Log.Infof("获取到solana水位调节参数：%s\r\n", configJsonStr)
		var configMap map[string]interface{}

		// 解析 JSON
		err := json.Unmarshal([]byte(configJsonStr), &configMap)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
		} else {
			var solMinDepositAmountThreshold, solMinWithdrawAmountThreshold float64
			alertThreshold := configMap["alertThreshold"].(float64)
			buyTriggerThreshold := configMap["buyTriggerThreshold"].(float64)
			sellTriggerThreshold := configMap["sellTriggerThreshold"].(float64)
			if v, ok := configMap["minDepositAmountThreshold"].(float64); ok {
				solMinDepositAmountThreshold = v
			} else {
				solMinDepositAmountThreshold = 0
			}

			if v, ok := configMap["minWithdrawAmountThreshold"].(float64); ok {
				solMinWithdrawAmountThreshold = v
			} else {
				solMinWithdrawAmountThreshold = 0
			}
			resp.SolWaterLevelConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{
				AlertThreshold:             &alertThreshold,
				BuyTriggerThreshold:        &buyTriggerThreshold,
				SellTriggerThreshold:       &sellTriggerThreshold,
				MinDepositAmountThreshold:  &solMinDepositAmountThreshold,
				MinWithdrawAmountThreshold: &solMinWithdrawAmountThreshold,
			}
		}
	}

	var stableData models.BusCommonConfig
	err = e.Orm.Model(&stableData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&stableData).Error

	resp.StableCoinWaterLevelConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{}
	if err != nil {
		e.Log.Errorf("获取稳定币水位调节参数失败")
	} else {
		configJsonStr := stableData.ConfigJson
		e.Log.Infof("获取到稳定币水位调节参数：%s\r\n", configJsonStr)
		var configMap map[string]interface{}

		// 解析 JSON
		err := json.Unmarshal([]byte(configJsonStr), &configMap)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
		} else {
			alertThreshold := configMap["alertThreshold"].(float64)
			resp.StableCoinWaterLevelConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{
				AlertThreshold: &alertThreshold,
			}
		}
	}

	return resp, nil
}

// UpdateGlobalWaterLevelConfig 更新全局WaterLevel 参数
func (e *BusDexCexTriangularObserver) UpdateGlobalWaterLevelConfig(req *dto.BusDexCexTriangularUpdateGlobalWaterLevelConfigReq) error {
	var data models.BusCommonConfig
	var instance models.BusDexCexTriangularObserver

	solWaterLevelConfigJsonStr, err := json.Marshal(req.SolWaterLevelConfig)
	if err != nil {
		e.Log.Errorf("JSON序列化失败,%s", err)
		return err
	}

	stableCoinWaterLevelConfigJsonStr, err := json.Marshal(req.StableCoinWaterLevelConfig)
	if err != nil {
		e.Log.Errorf("JSON序列化失败,%s", err)
		return err
	}

	//获取当前启动的所有水位调节实例
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		return err
	}

	var quoteTokens []string
	err = e.Orm.Model(&instance).
		Where("status = ? AND is_trading = ?", 3, true).
		Distinct().
		Pluck("quote_token", &quoteTokens).
		Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	isSolanaStarted := false
	for _, instanceId := range waterLevelInstances.InstanceIds {
		if instanceId == "SOLANA" {
			// solana 已经启动水位调节
			isSolanaStarted = true
			break
		}
	}

	if !isSolanaStarted {
		tokenConfig := &waterLevelPb.TokenThresholdConfig{
			AlertThreshold:             strconv.FormatFloat(*req.SolWaterLevelConfig.AlertThreshold, 'f', -1, 64),
			BuyTriggerThreshold:        strconv.FormatFloat(*req.SolWaterLevelConfig.BuyTriggerThreshold, 'f', -1, 64),
			SellTriggerThreshold:       strconv.FormatFloat(*req.SolWaterLevelConfig.SellTriggerThreshold, 'f', -1, 64),
			MinDepositAmountThreshold:  strconv.FormatFloat(*req.SolWaterLevelConfig.MinDepositAmountThreshold, 'f', -1, 64),
			MinWithdrawAmountThreshold: strconv.FormatFloat(*req.SolWaterLevelConfig.MinWithdrawAmountThreshold, 'f', -1, 64),
		}

		clientRequest := &waterLevelPb.StartInstanceRequest{
			InstanceId:           "SOLANA",
			ExchangeType:         "Binance",
			Currency:             "SOL",
			CurrencyType:         0, // token
			TokenThresholdConfig: tokenConfig,
		}

		e.Log.Infof("启动solana全局水位调节 req: %v \r\n", clientRequest)
		_, err = client.StartWaterLevelInstance(clientRequest)
		if err != nil {
			e.Log.Errorf("启动solana全局水位调节失败:%s \r\n", err)
			return errors.New("更新Solana全局水位调节失败")
		}
		e.Log.Infof("启动solana全局水位调节启动成功")
	} else {
		//如果已经启动了，要尝试更新
		tokenConfig := &waterLevelPb.TokenThresholdConfig{
			AlertThreshold:             strconv.FormatFloat(*req.SolWaterLevelConfig.AlertThreshold, 'f', -1, 64),
			BuyTriggerThreshold:        strconv.FormatFloat(*req.SolWaterLevelConfig.BuyTriggerThreshold, 'f', -1, 64),
			SellTriggerThreshold:       strconv.FormatFloat(*req.SolWaterLevelConfig.SellTriggerThreshold, 'f', -1, 64),
			MinDepositAmountThreshold:  strconv.FormatFloat(*req.SolWaterLevelConfig.MinDepositAmountThreshold, 'f', -1, 64),
			MinWithdrawAmountThreshold: strconv.FormatFloat(*req.SolWaterLevelConfig.MinWithdrawAmountThreshold, 'f', -1, 64),
		}
		updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
			InstanceId:           "SOLANA",
			CurrencyType:         0, // token
			TokenThresholdConfig: tokenConfig,
		}
		err = client.UpdateWaterLevelInstance(updateReq)
		if err != nil {
			e.Log.Errorf("启动solana全局水位调节失败:%s \r\n", err)
			return errors.New("更新Solana全局水位调节失败")
		}
		e.Log.Infof("更新solana全局水位调节启动成功")
	}

	// 保存配置到数据库
	err = e.Orm.Model(&data).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局solana水位调节参数")
			// 如果不存在，则新增
			data = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  common.GLOBAL_SOLANA_WATER_LEVEL_KEY,
				ConfigJson: string(solWaterLevelConfigJsonStr),
			}
			err = e.Orm.Create(&data).Error
			if err != nil {
				e.Log.Error("保存solana水位调节参数失败")
				return err
			}
		} else {
			e.Log.Errorf("db error: %s", err)
			return err
		}
	}

	updateData := map[string]interface{}{
		"config_json": string(solWaterLevelConfigJsonStr),
	}

	e.Log.Infof("更新参数：%s\n", string(solWaterLevelConfigJsonStr))
	err = e.Orm.Model(&data).
		Where("id = ?", data.Id).
		Updates(updateData).Error

	/***********稳定币部分**************/
	for _, quoteToken := range quoteTokens {
		isStarted := false
		for _, instanceId := range waterLevelInstances.InstanceIds {
			if instanceId == quoteToken {
				// 该稳定币已经启动水位调节
				isStarted = true
				break
			}
		}
		if !isStarted {
			tokenConfig := &waterLevelPb.StableCoinThresholdConfig{
				AlertThreshold: strconv.FormatFloat(*req.StableCoinWaterLevelConfig.AlertThreshold, 'f', -1, 64),
			}

			clientRequest := &waterLevelPb.StartInstanceRequest{
				InstanceId:                quoteToken,
				ExchangeType:              "Binance",
				Currency:                  quoteToken,
				CurrencyType:              1, // 稳定币
				StableCoinThresholdConfig: tokenConfig,
			}

			e.Log.Infof("启动稳定币 %s 全局水位调节 req: %v \r\n", quoteToken, clientRequest)
			_, err = client.StartWaterLevelInstance(clientRequest)
			if err != nil {
				e.Log.Errorf("启动稳定币 %s 全局水位调节 失败:%s \r\n", quoteToken, err)
				return err
			}
			e.Log.Infof("启动稳定币 %s 全局水位调节 成功", quoteToken)
		} else {
			//如果已经启动了，要尝试更新
			tokenConfig := &waterLevelPb.StableCoinThresholdConfig{
				AlertThreshold: strconv.FormatFloat(*req.StableCoinWaterLevelConfig.AlertThreshold, 'f', -1, 64),
			}
			updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
				InstanceId:                quoteToken,
				CurrencyType:              1, // 稳定币
				StableCoinThresholdConfig: tokenConfig,
			}
			err = client.UpdateWaterLevelInstance(updateReq)
			if err != nil {
				e.Log.Errorf("更新稳定币 %s 全局水位调节 失败:%s \r\n", quoteToken, err)
				return errors.New("更新稳定币全局水位调节失败")
			}
			e.Log.Infof("更新稳定币全局水位调节启动成功")
		}
	}

	var stableData models.BusCommonConfig
	// 稳定币水位调节参数处理
	err = e.Orm.Model(&stableData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&stableData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局稳定币水位调节参数")
			// 如果不存在，则新增
			stableData = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY,
				ConfigJson: string(stableCoinWaterLevelConfigJsonStr),
			}
			err = e.Orm.Create(&stableData).Error
			if err != nil {
				e.Log.Error("保存稳定币水位调节参数失败")
				return err
			}
		} else {
			e.Log.Errorf("db error: %s", err)
			return err
		}
	}

	updateData = map[string]interface{}{
		"config_json": string(stableCoinWaterLevelConfigJsonStr),
	}

	e.Log.Infof("更新参数：%s\n", string(stableCoinWaterLevelConfigJsonStr))
	err = e.Orm.Model(&stableData).
		Where("id = ?", stableData.Id).
		Updates(updateData).Error

	return nil
}

// GetGlobalRiskConfigState 获取全局风控 状态及当前参数
func (e *BusDexCexTriangularObserver) GetGlobalRiskConfigState() (*dto.BusDexCexTriangularUpdateGlobalRiskConfig, error) {
	var data models.BusCommonConfig
	resp := &dto.BusDexCexTriangularUpdateGlobalRiskConfig{}
	//step 2 : 封装全局的风控参数配置到响应体
	err := e.Orm.Model(&models.BusCommonConfig{}).
		Where("category = ? and config_key = ?", common.DEX_CEX_RISK_COTROL, common.RISK_CONTROL_CONFIG_KEY).
		Order("created_at desc").
		First(&data).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Errorf("获取风控参数失败:%s \r\n", err)
			return nil, err
		}
		e.Log.Errorf("不存在全局风控参数")
	} else {
		configJsonStr := data.ConfigJson
		e.Log.Infof("获取到全局风控参数：%s\r\n", configJsonStr)

		// 解析 JSON
		err := json.Unmarshal([]byte(configJsonStr), resp)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
			return nil, err
		}
	}

	return resp, nil
}

// UpdateGlobalRiskConfig 更新全局风控 参数
func (e *BusDexCexTriangularObserver) UpdateGlobalRiskConfig(req *dto.BusDexCexTriangularUpdateGlobalRiskConfig) error {
	var data models.BusCommonConfig

	riskConfigJsonStr, err := json.Marshal(req)
	if err != nil {
		e.Log.Errorf("JSON序列化失败,%s", err)
		return err
	}

	// 保存配置到数据库
	err = e.Orm.Model(&data).
		Where("category = ? and config_key = ?", common.DEX_CEX_RISK_COTROL, common.RISK_CONTROL_CONFIG_KEY).
		Order("created_at desc").
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局风控参数")
			// 如果不存在，则新增
			data = models.BusCommonConfig{
				Category:   common.DEX_CEX_RISK_COTROL,
				ConfigKey:  common.RISK_CONTROL_CONFIG_KEY,
				ConfigJson: string(riskConfigJsonStr),
			}
			err = e.Orm.Create(&data).Error
			if err != nil {
				e.Log.Error("更新全局风控失败")
				return err
			}
		} else {
			e.Log.Errorf("db error: %s", err)
			return err
		}
	}

	updateData := map[string]interface{}{
		"config_json": string(riskConfigJsonStr),
	}

	e.Log.Infof("更新参数：%s\n", string(riskConfigJsonStr))
	err = e.Orm.Model(&data).
		Where("id = ?", data.Id).
		Updates(updateData).Error
	if err != nil {
		e.Log.Errorf("更新全局风控参数失败：%s", err)
		return err
	}
	return nil
}

// StartGlobalWaterLevel 启动全局水位调整功能
func (e *BusDexCexTriangularObserver) StartGlobalWaterLevel() error {
	var data models.BusDexCexTriangularObserver
	var quoteTokens []string
	err := e.Orm.Model(&data).
		Where("status = ? AND is_trading = ?", 3, true).
		Distinct().
		Pluck("quote_token", &quoteTokens).
		Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	//if len(quoteTokens) == 0 {
	//	e.Log.Infof("没有交易中的币对，不需要启动全局水位调整功能，跳过")
	//	return nil
	//}

	var solData models.BusCommonConfig
	err = e.Orm.Model(&solData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&solData).Error

	var solanaAlertThreshold, solBuyTriggerThreshold, solSellTriggerThreshold, stableCoinAlertThreshold, solMinDepositAmountThreshold, solMinWithdrawAmountThreshold float64

	if err != nil {
		e.Log.Errorf("获取solana水位调节参数失败, 跳过本次启动全局水位调节任务")
		return err
	} else {
		configJsonStr := solData.ConfigJson
		e.Log.Infof("获取到solana水位调节参数：%s\r\n", configJsonStr)
		var configMap map[string]interface{}

		// 解析 JSON
		err := json.Unmarshal([]byte(configJsonStr), &configMap)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
			return err
		} else {
			solanaAlertThreshold = configMap["alertThreshold"].(float64)
			solBuyTriggerThreshold = configMap["buyTriggerThreshold"].(float64)
			solSellTriggerThreshold = configMap["sellTriggerThreshold"].(float64)
			if v, ok := configMap["minDepositAmountThreshold"].(float64); ok {
				solMinDepositAmountThreshold = v
			} else {
				solMinDepositAmountThreshold = 0
			}

			if v, ok := configMap["minWithdrawAmountThreshold"].(float64); ok {
				solMinWithdrawAmountThreshold = v
			} else {
				solMinWithdrawAmountThreshold = 0
			}
		}
	}

	var stableData models.BusCommonConfig
	err = e.Orm.Model(&stableData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&stableData).Error

	if err != nil {
		e.Log.Errorf("获取稳定币水位调节参数失败, 跳过本次启动全局水位调节任务")
		return err
	} else {
		configJsonStr := stableData.ConfigJson
		e.Log.Infof("获取到稳定币水位调节参数：%s\r\n", configJsonStr)
		var configMap map[string]interface{}

		// 解析 JSON
		err := json.Unmarshal([]byte(configJsonStr), &configMap)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
			return err
		} else {
			stableCoinAlertThreshold = configMap["alertThreshold"].(float64)
		}
	}

	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		return err
	}

	isSolanaStarted := false
	for _, instanceId := range waterLevelInstances.InstanceIds {
		if instanceId == "SOLANA" {
			// solana 已经启动水位调节
			isSolanaStarted = true
			break
		}
	}

	if !isSolanaStarted {
		tokenConfig := &waterLevelPb.TokenThresholdConfig{
			AlertThreshold:             strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64),
			BuyTriggerThreshold:        strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64),
			SellTriggerThreshold:       strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64),
			MinDepositAmountThreshold:  strconv.FormatFloat(solMinDepositAmountThreshold, 'f', -1, 64),
			MinWithdrawAmountThreshold: strconv.FormatFloat(solMinWithdrawAmountThreshold, 'f', -1, 64),
		}

		clientRequest := &waterLevelPb.StartInstanceRequest{
			InstanceId:           "SOLANA",
			ExchangeType:         "Binance",
			Currency:             "SOL",
			CurrencyType:         0, // token
			TokenThresholdConfig: tokenConfig,
		}

		e.Log.Infof("启动solana全局水位调节 req: %v \r\n", clientRequest)
		_, err = client.StartWaterLevelInstance(clientRequest)
		if err != nil {
			e.Log.Errorf("启动solana全局水位调节失败:%s \r\n", err)
			return err
		}
		e.Log.Infof("启动solana全局水位调节启动成功")
	} else {
		//如果已经启动了，要尝试更新
		// 先获取当前实例的参数
		solanaStateReq := &waterLevelPb.InstantId{
			InstanceId: "SOLANA",
		}
		solanaState, err := client.GetWaterLevelInstanceState(solanaStateReq)
		if err != nil {
			e.Log.Errorf("获取state 失败，直接更新")
			tokenConfig := &waterLevelPb.TokenThresholdConfig{
				AlertThreshold:             strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64),
				BuyTriggerThreshold:        strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64),
				SellTriggerThreshold:       strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64),
				MinDepositAmountThreshold:  strconv.FormatFloat(solMinDepositAmountThreshold, 'f', -1, 64),
				MinWithdrawAmountThreshold: strconv.FormatFloat(solMinWithdrawAmountThreshold, 'f', -1, 64),
			}
			updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
				InstanceId:           "SOLANA",
				CurrencyType:         0, // token
				TokenThresholdConfig: tokenConfig,
			}
			client.UpdateWaterLevelInstance(updateReq)
		} else {
			e.Log.Infof("从服务端获取到solana水位调节实例参数：%v \n", solanaState)
			oldParams := solanaState.InstanceParams.TokenThresholdConfig
			if oldParams.AlertThreshold == strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64) &&
				oldParams.BuyTriggerThreshold == strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64) &&
				oldParams.BuyTriggerThreshold == strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64) &&
				oldParams.MinDepositAmountThreshold == strconv.FormatFloat(solMinDepositAmountThreshold, 'f', -1, 64) &&
				oldParams.MinWithdrawAmountThreshold == strconv.FormatFloat(solMinWithdrawAmountThreshold, 'f', -1, 64) {
				// 参数一致，不需要更新
				e.Log.Infof("solana 全局水位调节参数一致，不需要更新，跳过")
			} else {
				// 参数不一致，更新
				tokenConfig := &waterLevelPb.TokenThresholdConfig{
					AlertThreshold:             strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64),
					BuyTriggerThreshold:        strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64),
					SellTriggerThreshold:       strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64),
					MinDepositAmountThreshold:  strconv.FormatFloat(solMinDepositAmountThreshold, 'f', -1, 64),
					MinWithdrawAmountThreshold: strconv.FormatFloat(solMinWithdrawAmountThreshold, 'f', -1, 64),
				}
				updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
					InstanceId:           "SOLANA",
					CurrencyType:         0, // token
					TokenThresholdConfig: tokenConfig,
				}
				client.UpdateWaterLevelInstance(updateReq)
			}
		}

	}

	log.Infof("waterLevelInstances:%+v\n", waterLevelInstances)

	for _, quoteToken := range quoteTokens {
		isStarted := false
		for _, instanceId := range waterLevelInstances.InstanceIds {
			if instanceId == quoteToken {
				// 该稳定币已经启动水位调节
				isStarted = true
				break
			}
		}
		if !isStarted {
			tokenConfig := &waterLevelPb.StableCoinThresholdConfig{
				AlertThreshold: strconv.FormatFloat(stableCoinAlertThreshold, 'f', -1, 64),
			}

			clientRequest := &waterLevelPb.StartInstanceRequest{
				InstanceId:                quoteToken,
				ExchangeType:              "Binance",
				Currency:                  quoteToken,
				CurrencyType:              1, // 稳定币
				StableCoinThresholdConfig: tokenConfig,
			}

			e.Log.Infof("启动稳定币 %s 全局水位调节 req: %v \r\n", quoteToken, clientRequest)
			_, err = client.StartWaterLevelInstance(clientRequest)
			if err != nil {
				e.Log.Errorf("启动稳定币 %s 全局水位调节 失败:%s \r\n", quoteToken, err)
				return err
			}
			e.Log.Infof("启动稳定币 %s 全局水位调节 成功", quoteToken)
		} else {
			//如果已经在服务端启动了，则需要比对参数，判断要不要更新
			stableStateReq := &waterLevelPb.InstantId{
				InstanceId: quoteToken,
			}
			stableState, err := client.GetWaterLevelInstanceState(stableStateReq)
			if err != nil {
				e.Log.Errorf("获取state 失败，直接更新")
				tokenConfig := &waterLevelPb.StableCoinThresholdConfig{
					AlertThreshold: strconv.FormatFloat(stableCoinAlertThreshold, 'f', -1, 64),
				}
				updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
					InstanceId:                quoteToken,
					CurrencyType:              1, // token
					StableCoinThresholdConfig: tokenConfig,
				}
				client.UpdateWaterLevelInstance(updateReq)
			} else {
				e.Log.Infof("从服务端获取到稳定币： %s 水位调节实例参数：%v \n", quoteToken, stableState)
				oldParams := stableState.InstanceParams.StableCoinThresholdConfig
				if oldParams.AlertThreshold == strconv.FormatFloat(stableCoinAlertThreshold, 'f', -1, 64) {
					// 参数一致，不需要更新
					e.Log.Infof("稳定币 %s 全局水位调节参数一致，不需要更新，跳过", quoteToken)
				} else {
					tokenConfig := &waterLevelPb.StableCoinThresholdConfig{
						AlertThreshold: strconv.FormatFloat(stableCoinAlertThreshold, 'f', -1, 64),
					}
					updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
						InstanceId:                quoteToken,
						CurrencyType:              1, // token
						StableCoinThresholdConfig: tokenConfig,
					}
					client.UpdateWaterLevelInstance(updateReq)
				}
			}
		}

	}
	return nil
}

func StartObserver(observer *models.BusDexCexTriangularObserver) error {

	maxArraySize := new(uint32)
	*maxArraySize = uint32(observer.MaxArraySize) //默认5， clmm使用参数

	dexConfig := &pb.DexConfig{}
	if observer.DexType == "RAY_AMM" {
		dexConfig.Config = &pb.DexConfig_RayAmm{
			RayAmm: &pb.RayAmmConfig{
				Pool:      observer.AmmPoolId,
				TokenMint: observer.TokenMint,
			},
		}
	} else if observer.DexType == "RAY_CLMM" {
		dexConfig.Config = &pb.DexConfig_RayClmm{
			RayClmm: &pb.RayClmmConfig{
				Pool:         observer.AmmPoolId,
				TokenMint:    observer.TokenMint,
				MaxArraySize: maxArraySize,
			},
		}
	}

	//triggerHoldingMsUint := uint64(observer.TriggerHoldingMs)
	arbitrageConfig := &pb.ObserverParams{
		MinQuoteAmount: observer.MinQuoteAmount,
		MaxQuoteAmount: observer.MaxQuoteAmount,
		//SlippageRate:      observer.SlippageBpsRate,
		ProfitTriggerRate: observer.ProfitTriggerRate,
		//TriggerHoldingMs:  &triggerHoldingMsUint,
	}

	amberConfig := &pb.AmberObserverConfig{}
	GenerateAmberConfig(observer, amberConfig)

	instanceId := strconv.Itoa(observer.Id)
	log.Infof("restart observer success with params: dexConfig: %+v\n, arbitrageConfig: %+v\n", dexConfig, arbitrageConfig)
	err := client.StartNewArbitragerClient(&instanceId, amberConfig, dexConfig, arbitrageConfig)
	if err != nil {
		log.Errorf("restart observer throw grpc error: %v\n", err)
		return err
	}
	return nil
}

func GenerateAmberConfig(observer *models.BusDexCexTriangularObserver, amberConfig *pb.AmberObserverConfig) error {
	amberConfig.ExchangeType = &observer.ExchangeType
	amberConfig.TakerFee = proto.Float64(*observer.TakerFee)

	amberConfig.TargetToken = &observer.TargetToken
	amberConfig.QuoteToken = &observer.QuoteToken

	if observer.Depth != "" {
		depthInt, err := strconv.Atoi(observer.Depth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		amberConfig.BidDepth = proto.Int32(int32(depthInt))
		amberConfig.AskDepth = proto.Int32(int32(depthInt))
	}
	return nil
}

func StartTrader(instance *models.BusDexCexTriangularObserver) error {

	exchangeType, err := instance.GetExchangeTypeForTrader()
	if err != nil {
		log.Errorf("获取ExchangeType参数异常，:%s \r\n", err)
		return err
	}

	//slippageBpsUint, err := strconv.ParseUint(instance.SlippageBps, 10, 64)
	//if err != nil {
	//	log.Errorf("转换失败: %s", err)
	//}
	//log.Infof("slippageBps: %v\n", slippageBpsUint)
	//slippageBpsFloat := float64(slippageBpsUint) / 10000.0

	amberTraderConfig := &pb.AmberTraderConfig{
		ExchangeType: &exchangeType,
	}

	jitoFee := instance.JitoFeeRate
	traderParams := &pb.TraderParams{
		//Slippage:    &slippageBpsFloat,
		SlippageRate:    instance.SlippageBpsRate,
		PriorityFeeRate: instance.PriorityFeeRate,
		JitoFeeRate:     jitoFee,
	}
	if config.ApplicationConfig.Mode != "dev" {
		instanceId := strconv.Itoa(instance.Id)
		err := client.EnableTrader(instanceId, amberTraderConfig, traderParams)
		if err != nil {
			log.Errorf("GRPC 启动Trader for instanceId:%d 失败，异常:%s \r\n", instance.Id, err)
			return err
		}
	}
	return nil
}

// StartTokenWaterLevelWithCheckExists 启动水位调节，校验是否存在
func StartTokenWaterLevelWithCheckExists(observer *models.BusDexCexTriangularObserver) error {
	instances, err := client.ListWaterLevelInstance()
	if err != nil {
		log.Errorf("water level 服务不可用，:%s \r\n", err)
		return err
	}
	ids := instances.InstanceIds
	isExist := false
	for _, id := range ids {
		if id == strconv.Itoa(observer.Id) {
			isExist = true
			break
		}
	}
	if isExist {
		//说明已经启动实例了，此时更新实例参数
		err = UpdateTokenWaterLevel(observer)
		if err != nil {
			log.Errorf("grpc启动实例失败，异常信息:%s \r\n", err)
			return err
		}
	} else {
		//说明未启动实例，此时启动新的实例
		err = StartTokenWaterLevel(observer)
		if err != nil {
			log.Errorf("grpc启动实例失败，异常信息:%s \r\n", err)
			return err
		}
	}
	return nil
}

// StartTokenWaterLevel 启动水位调节，不校验是否存在
func StartTokenWaterLevel(observer *models.BusDexCexTriangularObserver) error {
	tokenConfig := &waterLevelPb.TokenThresholdConfig{
		AlertThreshold:             strconv.FormatFloat(*observer.AlertThreshold, 'f', -1, 64),
		BuyTriggerThreshold:        strconv.FormatFloat(*observer.BuyTriggerThreshold, 'f', -1, 64),
		SellTriggerThreshold:       strconv.FormatFloat(*observer.SellTriggerThreshold, 'f', -1, 64),
		MinDepositAmountThreshold:  strconv.FormatFloat(*observer.MinDepositAmountThreshold, 'f', -1, 64),
		MinWithdrawAmountThreshold: strconv.FormatFloat(*observer.MinWithdrawAmountThreshold, 'f', -1, 64),
	}

	clientRequest := &waterLevelPb.StartInstanceRequest{
		InstanceId:           strconv.Itoa(observer.Id),
		ExchangeType:         observer.ExchangeType,
		Currency:             observer.TargetToken,
		CurrencyType:         0, // token
		PubKey:               observer.TokenMint,
		TokenThresholdConfig: tokenConfig,
	}

	log.Infof("restart water level with req: %v \r\n", clientRequest)
	_, err := client.StartWaterLevelInstance(clientRequest)
	if err != nil {
		log.Errorf("启动水位调节失败:%s \r\n", err)
		return err
	}
	log.Infof("水位调节启动成功")
	return nil
}

// UpdateTokenWaterLevel 更新水位调节
func UpdateTokenWaterLevel(observer *models.BusDexCexTriangularObserver) error {
	instanceId := strconv.Itoa(observer.Id)

	tokenConfig := &waterLevelPb.TokenThresholdConfig{
		AlertThreshold:             strconv.FormatFloat(*observer.AlertThreshold, 'f', -1, 64),
		BuyTriggerThreshold:        strconv.FormatFloat(*observer.BuyTriggerThreshold, 'f', -1, 64),
		SellTriggerThreshold:       strconv.FormatFloat(*observer.SellTriggerThreshold, 'f', -1, 64),
		MinDepositAmountThreshold:  strconv.FormatFloat(*observer.MinDepositAmountThreshold, 'f', -1, 64),
		MinWithdrawAmountThreshold: strconv.FormatFloat(*observer.MinWithdrawAmountThreshold, 'f', -1, 64),
	}

	waterLevelParams := &waterLevelPb.UpdateInstanceParamsRequest{
		InstanceId:           instanceId,
		CurrencyType:         0, //Token
		TokenThresholdConfig: tokenConfig,
	}

	err := client.UpdateWaterLevelInstance(waterLevelParams)
	if err != nil {
		log.Errorf("grpc更新实例：:%s water level参数失败，异常：%s \r\n", instanceId, err)
		return err
	}
	return nil
}

// CheckRiskControl 风控校验，当前采用的是定时任务，最佳实现应该是使用事件驱动
func (e BusDexCexTriangularObserver) CheckRiskControl() error {
	// 获取实例的风控参数
	var riskConfig models.BusCommonConfig
	err := e.Orm.Model(&models.BusCommonConfig{}).
		Where("category = ? and config_key = ?", common.DEX_CEX_RISK_COTROL, common.RISK_CONTROL_CONFIG_KEY).
		Order("created_at desc").
		First(&riskConfig).Error
	if err != nil {
		e.Log.Errorf("[Risk Control Check] 获取风控参数失败:%s \r\n", err)
		return err
	}

	configJsonStr := riskConfig.ConfigJson
	e.Log.Infof("[Risk Control Check] 获取到链上链下风控参数：%s\r\n", configJsonStr)
	var configMap map[string]interface{}

	// 解析 JSON
	err = json.Unmarshal([]byte(configJsonStr), &configMap)
	if err != nil {
		e.Log.Error("[Risk Control Check] JSON 解析失败:", err)
		return err
	}

	// 单笔最大亏损金额阈值
	var absoluteLossThreshold []interface{}
	if v, ok := configMap["absoluteLossThreshold"]; ok {
		if list, valid := v.([]interface{}); valid {
			absoluteLossThreshold = list
		}
	}
	e.Log.Infof("[Risk Control Check] absoluteLossThreshold: %v", absoluteLossThreshold)
	// 单笔最大亏损比例阈值
	var relativeLossThreshold []interface{}

	if v, ok := configMap["relativeLossThreshold"]; ok {
		if list, valid := v.([]interface{}); valid {
			relativeLossThreshold = list
		}
	}
	e.Log.Infof("[Risk Control Check] relativeLossThreshold: %v", relativeLossThreshold)

	// 排序，按照 action 从大到小排序
	sort.Slice(absoluteLossThreshold, func(i, j int) bool {
		return absoluteLossThreshold[i].(map[string]interface{})["action"].(float64) > absoluteLossThreshold[j].(map[string]interface{})["action"].(float64)
	})

	sort.Slice(relativeLossThreshold, func(i, j int) bool {
		return relativeLossThreshold[i].(map[string]interface{})["action"].(float64) > relativeLossThreshold[j].(map[string]interface{})["action"].(float64)
	})

	// 从交易中，获取成功的交易记录，并且未完成风控校验的记录。
	var riskCheckProgress models.BusRiskCheckProgress
	err = e.Orm.Model(&models.BusRiskCheckProgress{}).
		Where("strategy_id = ?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("business_type = ?", common.BUSINESS_TYPE_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("trade_table =?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE_TRADES_TABLE).
		Order("created_at desc").
		First(&riskCheckProgress).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("[Risk Control Check] 当前不存在风控校验进度")
			// 如果不存在，则新增
			riskCheckProgress = models.BusRiskCheckProgress{
				StrategyId:         common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE,
				BusinessType:       common.BUSINESS_TYPE_DEX_CEX_TRIANGULAR_ARBITRAGE,
				TradeTable:         common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE_TRADES_TABLE,
				LastCheckedTradeId: 0,
				LastCheckedAt:      time.Now(),
				Status:             0, // 未启动
			}
			err = e.Orm.Create(&riskCheckProgress).Error
			if err != nil {
				e.Log.Error("[Risk Control Check] 保存风控校验进度失败")
				return err
			}
		} else {
			e.Log.Errorf("[Risk Control Check] db error: %+v", err)
			return err
		}
	}

	if riskCheckProgress.Status == common.RISK_CHECK_STATUS_PROCESSING {
		e.Log.Infof("[Risk Control Check] 风控校验正在进行中，跳过")
		return nil
	}

	var trades []dto.StrategyDexCexTriangularArbitrageTradesGetPageResp
	err = e.Orm.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("strategy_dex_cex_triangular_arbitrage_trades.*, opportunities.cex_target_asset as symbol").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("strategy_dex_cex_triangular_arbitrage_trades.id > ?", riskCheckProgress.LastCheckedTradeId).
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success = 1 and strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = 1").
		Order("strategy_dex_cex_triangular_arbitrage_trades.created_at asc").
		Limit(20). // 每5s处理最多20条记录
		Find(&trades).Error

	if err != nil {
		e.Log.Errorf("[Risk Control Check] 获取套利记录失败:%s \r\n", err)
		return err
	}

	config := ext.ExtConfig
	larkClient := lark.NewLarkRobotAlert(config)

	// 开始风控校验
	err = e.Orm.Model(&riskCheckProgress).
		Where("id =?", riskCheckProgress.Id).
		Updates(map[string]interface{}{
			"status": common.RISK_CHECK_STATUS_PROCESSING,
		}).Error
	if err != nil {
		e.Log.Errorf("[Risk Control Check] 保存风控校验状态失败:%s \r\n", err)
		return err
	}

	var errOccurred error
	for _, trade := range trades {
		e.Log.Infof("[Risk Control Check] 交易订单:%d \r\n", trade.Id)
		maxAfterAction := 0
		// 1. 单笔最大亏损金额阈值
		maxAfterAction, err := e.AbsoluteLossThresholdCheck(absoluteLossThreshold, trade, *larkClient)
		if err != nil {
			e.Log.Errorf("[Risk Control Check] 单笔最大亏损金额阈值校验失败:%s \r\n", err)
			SendRiskCheckFailMessage(common.TRIGGER_RULE_ABSOLUTE_LOSS_THRESHOLD, *larkClient)
			errOccurred = err
			break
		}

		// 2. 单笔最大亏损比例阈值
		afterAction, err := e.RelativeLossThresholdCheck(relativeLossThreshold, trade, *larkClient)
		if err != nil {
			SendRiskCheckFailMessage(common.TRIGGER_RULE_RELATIVE_LOSS_THRESHOLD, *larkClient)
			e.Log.Errorf("[Risk Control Check] 单笔最大亏损比例阈值校验失败:%s \r\n", err)
			errOccurred = err
			break
		}

		if afterAction > maxAfterAction {
			maxAfterAction = afterAction
		}

		// TODO 3. 单币种单日累计亏损金额阈值

		// TODO 4. 全币种单日累计亏损金额阈值

		// 全部完成风控校验后，更新风控进度表的最后校验ID
		err = e.Orm.Model(&riskCheckProgress).
			Where("id = ?", riskCheckProgress.Id).
			Updates(map[string]interface{}{
				"last_checked_trade_id": trade.Id,
				"last_checked_at":       time.Now(),
			}).Error
		if err != nil {
			e.Log.Errorf("[Risk Control Check] 保存风控校验进度失败:%s \r\n", err)
			errOccurred = err
			break
		}

		// 最后根据执行动作完成下一步动作
		if maxAfterAction == 1 {
			// do nothing
			continue
		}
		if maxAfterAction == 2 {
			//暂停该笔trade对应的instance的交易功能

			stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
				InstanceId: trade.Id,
			}

			//  暂停交易如果报错了，要如何补偿？定时任务补偿？
			err = e.StopTrader(&stopTradeReq, true)
			if err != nil {
				e.Log.Errorf("[Risk Control Check] 暂停交易失败:%s \r\n", err)
				message := fmt.Sprintf(`
				❌ 风控触发, 暂停交易失败, instanceId: %s, 操作时间: %s, 定时任务会进行补偿, 请注意风险。
				`, trade.InstanceId, time.Now().Format("2006-01-02 15:04:05"))
				larkClient.SendLarkAlert(message)
				continue
			}

			continue
		}

		if maxAfterAction == 3 {
			// 暂停全部交易
			// 暂停全部交易行为暂不支持
			continue
		}

	}

	// 全部订单都检查完后，更新风控校验状态
	err = e.Orm.Model(&riskCheckProgress).
		Where("id =?", riskCheckProgress.Id).
		Updates(map[string]interface{}{
			"status": common.RISK_CHECK_STATUS_FINISHED,
		}).Error
	if err != nil {
		e.Log.Errorf("[Risk Control Check] 保存风控校验状态失败:%s \r\n", err)
		return err
	}

	// **如果中间发生错误，则返回错误**
	if errOccurred != nil {
		return errOccurred
	}

	return nil

}

/*
*

	{
		"absoluteLossThreshold": [
			{
				"threshold": 10,
				"action": 1,
				"action_detail": {
					"notify": true
				}
			},
			{
				"threshold": 20,
				"action": 2,
				"action_detail": {
					"pause_duration": 3600, // -1的话，表示为次日0点恢复
					"manual_resume": false
				}
			}
		],
		"relativeLossThreshold":[
			{
				"threshold": 0.1,
				"action": 1,
				"action_detail": {
					"notify": true
				}
			},
			{
				"threshold": 0.2,
				"action": 2,
				"action_detail": {
					"pause_duration": 3600, // -1的话，表示为次日0点恢复
					"manual_resume": false
				}
			}
		]
	}
*/
func (e BusDexCexTriangularObserver) AbsoluteLossThresholdCheck(absoluteLossThreshold []interface{}, trade dto.StrategyDexCexTriangularArbitrageTradesGetPageResp, larkClient lark.LarkRobotAlert) (int, error) {
	afterAction := 0
	for _, threshold := range absoluteLossThreshold {

		thresholdMap := threshold.(map[string]interface{})
		thresholdValue := thresholdMap["threshold"].(float64)
		action := thresholdMap["action"].(float64)
		actionDetail := thresholdMap["actionDetail"].(map[string]interface{})

		cexSellAmount, err1 := strconv.ParseFloat(trade.CexSellQuoteAmount, 64)
		cexBuyAmount, err2 := strconv.ParseFloat(trade.CexBuyQuoteAmount, 64)
		if err1 != nil || err2 != nil {
			e.Log.Errorf("[Risk Control Check] 交易订单金额解析失败:%s \r\n")
			return 0, errors.New("交易订单金额解析失败")
		}

		profitAmount := cexSellAmount - cexBuyAmount
		e.Log.Infof("[Risk Control Check] profit:%f , threshold:%f \n", profitAmount, thresholdValue)
		if profitAmount < -thresholdValue {
			// 亏损金额超过阈值
			// 生成风控事件
			riskEvent := models.BusRiskEvent{
				StrategyId:         common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE,
				StrategyInstanceId: trade.InstanceId,
				AssetSymbol:        trade.Symbol,
				TriggerRule:        common.TRIGGER_RULE_ABSOLUTE_LOSS_THRESHOLD,
				TriggerValue:       strconv.FormatFloat(profitAmount, 'f', -1, 64),
				TradeId:            trade.Id,
			}
			if action == 3 {
				// 暂停全部交易，单笔交易亏损触发暂不支持暂停全部交易行为
				e.Log.Errorf("[Risk Control Check] 单笔交易亏损触发暂不支持暂停全部交易行为 \r\n")
				return afterAction, errors.New("单笔交易亏损触发暂不支持暂停全部交易行为")
			} else if action == 2 {
				// 暂停当前策略实例的交易功能
				afterAction = 2

				riskEvent.RiskScope = common.RISK_SCOPE_SINGLE_TOKEN
				riskEvent.RiskLevel = common.RISK_LEVEL_MIDDLE
				riskEvent.IsRecovered = 0
				manualRecover := actionDetail["manualResume"].(int)
				riskEvent.ManualRecover = manualRecover
				if manualRecover != 0 {
					pauseDuration := actionDetail["pauseDuration"].(int64)
					if pauseDuration == -1 {
						// 第二天0点恢复
						now := time.Now()
						tomorrow := now.AddDate(0, 0, 1)
						tomorrowZero := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
						riskEvent.AutoRecoverTime = &tomorrowZero
					} else {
						// 按指定暂停时长
						recoverTime := time.Now().Add(time.Duration(pauseDuration) * time.Second)
						riskEvent.AutoRecoverTime = &recoverTime
					}
				}

				err := e.Orm.Create(&riskEvent).Error
				if err != nil {
					e.Log.Errorf("[Risk Control Check] 生成风控事件失败:%s \r\n")
					return 0, errors.New("生成风控事件失败")
				}

				// 暂停交易
				// InstanceIdInt, err := strconv.Atoi(trade.InstanceId)
				// if err != nil {
				// 	e.Log.Errorf("[Risk Control Check] 实例ID解析失败:%s \r\n")
				// 	return false, errors.New("实例ID解析失败")
				// }

				// stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
				// 	InstanceId: InstanceIdInt,
				// }
				// // TODO 暂停交易如果报错了，要如何补偿？定时任务补偿？
				// err = e.StopTrader(&stopTradeReq)
				// if err != nil {
				// 	e.Log.Errorf("[Risk Control Check] 暂停交易失败:%s \r\n", err)
				// }

				var recoverMethod string
				if manualRecover == 1 {
					recoverMethod = "手动恢复"
				} else {
					recoverMethod = "自动恢复"
				}

				SendMiddleNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_ABSOLUTE_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), strconv.FormatFloat(profitAmount, 'f', -1, 64), recoverMethod, larkClient)

			} else if action == 1 {
				// 预警
				afterAction = 0
				riskEvent.RiskScope = common.RISK_SCOPE_SINGLE_TOKEN
				riskEvent.RiskLevel = common.RISK_LEVEL_LOW
				riskEvent.IsRecovered = 1 // 预警类的不需要进行恢复，不阻断流程
				nowTime := time.Now()
				riskEvent.RecoveredAt = &nowTime
				riskEvent.ManualRecover = 0
				riskEvent.AutoRecoverTime = &nowTime
				riskEvent.RecoveredBy = "-1"

				err := e.Orm.Create(&riskEvent).Error
				if err != nil {
					e.Log.Errorf("[Risk Control Check] 生成风控事件失败:%s \r\n")
					return 0, errors.New("生成风控事件失败")
				}

				SendWarningNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_ABSOLUTE_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), strconv.FormatFloat(profitAmount, 'f', -1, 64), larkClient)
			}
			//只要触发了高级别的风控策略，就不会再匹配同类型下的低级别风控规则
			return afterAction, nil
		}
	}
	return afterAction, nil
}

func (e BusDexCexTriangularObserver) RelativeLossThresholdCheck(relativeLossThreshold []interface{}, trade dto.StrategyDexCexTriangularArbitrageTradesGetPageResp, larkClient lark.LarkRobotAlert) (int, error) {
	afterAction := 0
	for _, threshold := range relativeLossThreshold {

		thresholdMap := threshold.(map[string]interface{})
		thresholdValue := thresholdMap["threshold"].(float64)
		action := thresholdMap["action"].(float64)
		actionDetail := thresholdMap["actionDetail"].(map[string]interface{})

		cexSellAmount, err1 := strconv.ParseFloat(trade.CexSellQuoteAmount, 64)
		cexBuyAmount, err2 := strconv.ParseFloat(trade.CexBuyQuoteAmount, 64)
		if err1 != nil || err2 != nil {
			e.Log.Errorf("[Risk Control Check] 交易订单金额解析失败:%s \r\n")
			return 0, errors.New("交易订单金额解析失败")
		}

		profitPercent := (cexSellAmount - cexBuyAmount) / cexBuyAmount
		profitPercentStr := strconv.FormatFloat(profitPercent*100, 'f', 2, 64) + "%"
		if profitPercent < -thresholdValue {
			// 亏损比例超过阈值
			// 生成风控事件
			riskEvent := models.BusRiskEvent{
				StrategyId:         common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE,
				StrategyInstanceId: trade.InstanceId,
				TradeId:            trade.Id,
				AssetSymbol:        trade.Symbol,
				TriggerRule:        common.TRIGGER_RULE_RELATIVE_LOSS_THRESHOLD,
				TriggerValue:       profitPercentStr,
			}
			if action == 3 {
				// 暂停全部交易，单笔交易亏损触发暂不支持暂停全部交易行为
				e.Log.Errorf("[Risk Control Check] 单笔交易亏损比例触发暂不支持暂停全部交易行为 \r\n")
				return 0, errors.New("单笔交易亏损比例触发暂不支持暂停全部交易行为")
			} else if action == 2 {
				// 暂停当前策略实例
				afterAction = 2
				riskEvent.RiskScope = common.RISK_SCOPE_SINGLE_TOKEN
				riskEvent.RiskLevel = common.RISK_LEVEL_MIDDLE
				riskEvent.IsRecovered = 0
				manualRecover := actionDetail["manualResume"].(int)
				riskEvent.ManualRecover = manualRecover
				if manualRecover != 1 {
					pauseDuration := actionDetail["pauseDuration"].(int64)
					if pauseDuration == -1 {
						// 第二天0点恢复
						now := time.Now()
						tomorrow := now.AddDate(0, 0, 1)
						tomorrowZero := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
						riskEvent.AutoRecoverTime = &tomorrowZero
					} else {
						// 按指定暂停时长
						recoverTime := time.Now().Add(time.Duration(pauseDuration) * time.Second)
						riskEvent.AutoRecoverTime = &recoverTime
					}
				}

				err := e.Orm.Create(&riskEvent).Error
				if err != nil {
					e.Log.Errorf("[Risk Control Check] 生成风控事件失败:%s \r\n")
					return 0, errors.New("生成风控事件失败")
				}

				// 暂停交易
				// InstanceIdInt, err := strconv.Atoi(trade.InstanceId)
				// if err != nil {
				// 	e.Log.Errorf("[Risk Control Check] 实例ID解析失败:%s \r\n")
				// 	return needStopTrade, errors.New("实例ID解析失败")
				// }

				// stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
				// 	InstanceId: InstanceIdInt,
				// }

				// // TODO 暂停交易如果报错了，要如何补偿？定时任务补偿？
				// err = e.StopTrader(&stopTradeReq)
				// if err != nil {
				// 	e.Log.Errorf("[Risk Control Check] 暂停交易失败:%s \r\n", err)
				// }

				var recoverMethod string
				if manualRecover == 1 {
					recoverMethod = "手动恢复"
				} else {
					recoverMethod = "自动恢复"
				}

				SendMiddleNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_RELATIVE_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), profitPercentStr, recoverMethod, larkClient)

			} else if action == 1 {
				// 预警
				afterAction = 1
				riskEvent.RiskScope = common.RISK_SCOPE_SINGLE_TOKEN
				riskEvent.RiskLevel = common.RISK_LEVEL_LOW
				riskEvent.IsRecovered = 1 // 预警类的不需要进行恢复，不阻断流程
				nowTime := time.Now()
				riskEvent.RecoveredAt = &nowTime
				riskEvent.ManualRecover = 0
				riskEvent.AutoRecoverTime = &nowTime
				riskEvent.RecoveredBy = "-1"

				err := e.Orm.Create(&riskEvent).Error
				if err != nil {
					e.Log.Errorf("[Risk Control Check] 生成风控事件失败:%s \r\n", err)
					return 0, errors.New("生成风控事件失败")
				}

				SendWarningNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_RELATIVE_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), profitPercentStr, larkClient)
			}
			//只要触发了高级别的风控策略，就不会再匹配同类型下的低级别风控规则
			return afterAction, nil
		}
	}
	return afterAction, nil
}

func (e BusDexCexTriangularObserver) CheckExistRiskEvent() error {
	// 1. 查询出所有未恢复的全局风控事件
	var highestRiskEvents []models.BusRiskEvent
	err := e.Orm.Model(models.BusRiskEvent{}).
		Where("is_recovered =?", 0).
		Where("strategy_id =?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("risk_scope =? and risk_level > ?", common.RISK_SCOPE_GOLBAL, common.RISK_LEVEL_MIDDLE).
		Find(&highestRiskEvents).Error
	if err != nil {
		e.Log.Errorf("查询风控事件失败:%s \r\n", err)
		return err
	}

	// 查询出所有交易开启中以及水位调节中的实例
	var instances []models.BusDexCexTriangularObserver
	err = e.Orm.Model(models.BusDexCexTriangularObserver{}).
		Where("status IN ?", []int{2, 3}).
		Find(&instances).Error
	if err != nil {
		e.Log.Errorf("查询实例失败:%s \r\n", err)
		return err
	}

	if len(highestRiskEvents) == 0 {
		e.Log.Infof("当前不存在未恢复的风控事件 \r\n")
	} else {
		e.Log.Infof("存在全局中断交易的事件,暂停全部实例交易功能")

		// 关闭所有实例
		for _, instance := range instances {
			stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
				InstanceId: instance.Id,
			}

			if instance.Status == "2" { //水位调节中
				stopReq := &waterLevelPb.InstantId{
					InstanceId: strconv.Itoa(instance.Id),
				}
				err = client.StopWaterLevelInstance(stopReq)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instance.Id, err)
					continue
				}
				// 更新observer的status =1
				updateData := map[string]interface{}{
					"status":             1,
					"is_trading_blocked": true,
				}

				err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
					Where("id = ?", instance.Id).
					Updates(updateData).Error
				if err != nil {
					e.Log.Errorf("更新实例状态失败, 异常：%s \r\n", err)
					continue
				}

			} else if instance.Status == "3" { //交易开启中
				err = e.StopTrader(&stopTradeReq, true)
				if err != nil {
					e.Log.Errorf("关闭实例%d失败:%s \r\n", instance.Id, err)
				} else {
					e.Log.Infof("关闭实例%d成功 \r\n", instance.Id)
				}
			}
			continue
		}
		e.Log.Infof("完成暂停全部实例交易功能")
		return nil
	}

	// 2. 查出所有未恢复的单币种风控事件
	var singleTokenRiskEvents []models.BusRiskEvent
	err = e.Orm.Model(models.BusRiskEvent{}).
		Where("is_recovered =?", 0).
		Where("strategy_id =?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("risk_scope =? and risk_level =?", common.RISK_SCOPE_SINGLE_TOKEN, common.RISK_LEVEL_MIDDLE).
		Find(&singleTokenRiskEvents).Error

	if err != nil {
		e.Log.Errorf("查询单币种风控事件失败:%s \r\n", err)
		return err
	}

	if len(singleTokenRiskEvents) == 0 {
		e.Log.Infof("当前不存在未恢复的单币种风控事件 \r\n")
		return nil
	}

	// 暂停所有单币种风控事件对应的实例
	for _, instance := range instances {
		e.Log.Infof("当前存在未恢复的单币种风控事件 \r\n")
		for _, riskEvent := range singleTokenRiskEvents {
			if riskEvent.StrategyInstanceId == strconv.Itoa(instance.Id) {
				stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
					InstanceId: instance.Id,
				}

				if instance.Status == "2" { //水位调节中
					stopReq := &waterLevelPb.InstantId{
						InstanceId: strconv.Itoa(instance.Id),
					}
					err = client.StopWaterLevelInstance(stopReq)
					if err != nil {
						e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instance.Id, err)
						continue
					}
					// 更新observer的status =1
					updateData := map[string]interface{}{
						"status":             1,
						"is_trading_blocked": true,
					}

					err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
						Where("id = ?", instance.Id).
						Updates(updateData).Error
					if err != nil {
						e.Log.Errorf("更新实例状态失败, 异常：%s \r\n", err)
						continue
					}

				} else if instance.Status == "3" { //交易开启中
					err = e.StopTrader(&stopTradeReq, true)
					if err != nil {
						e.Log.Errorf("关闭实例%d失败:%s \r\n", instance.Id, err)
					} else {
						e.Log.Infof("关闭实例%d成功 \r\n", instance.Id)
					}
				}
			}
		}

	}

	return nil
}

func (e BusDexCexTriangularObserver) CheckBlockingInstance() error {
	// 1. 查询出所有未恢复的全局风控事件
	var highestRiskEvents []models.BusRiskEvent
	err := e.Orm.Model(models.BusRiskEvent{}).
		Where("is_recovered =?", 0).
		Where("strategy_id =?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("risk_scope =? and risk_level > ?", common.RISK_SCOPE_GOLBAL, common.RISK_LEVEL_MIDDLE).
		Find(&highestRiskEvents).Error
	if err != nil {
		e.Log.Errorf("查询风控事件失败:%s \r\n", err)
		return err
	}

	// 查询出所有交易开启中以及水位调节中的实例
	var instances []models.BusDexCexTriangularObserver
	err = e.Orm.Model(models.BusDexCexTriangularObserver{}).
		Where("status = ?", 1).
		Where("is_trading_blocked = true").
		Find(&instances).Error
	if err != nil {
		e.Log.Errorf("查询实例失败:%s \r\n", err)
		return err
	}

	if len(highestRiskEvents) > 0 {
		e.Log.Infof("当前存在未恢复的全局风控事件 \r\n")
		//直接结束
		return nil
	}

	e.Log.Infof("不存在全局中断交易的事件")

	// 2. 查出所有未恢复的单币种风控事件
	var singleTokenRiskEvents []models.BusRiskEvent
	err = e.Orm.Model(models.BusRiskEvent{}).
		Where("is_recovered =?", false).
		Where("strategy_id =?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("risk_scope =? and risk_level =?", common.RISK_SCOPE_SINGLE_TOKEN, common.RISK_LEVEL_MIDDLE).
		Find(&singleTokenRiskEvents).Error

	if err != nil {
		e.Log.Errorf("查询单币种风控事件失败:%s \r\n", err)
		return err
	}

	// 暂停所有单币种风控事件对应的实例
	for _, instance := range instances {
		var hasUnRecoveryRiskEvent bool = false
		for _, riskEvent := range singleTokenRiskEvents {
			if riskEvent.StrategyInstanceId == strconv.Itoa(instance.Id) {
				hasUnRecoveryRiskEvent = true
				break
			}
		}
		if hasUnRecoveryRiskEvent {
			e.Log.Infof("instanceId: %d has un-recoverd risk event", instance.Id)
			continue
		}
		//启动交易功能
		// step1 先启动水位调节实例
		err = StartTokenWaterLevelWithCheckExists(&instance)
		if err != nil {
			return err
		}

		// 启动水位调节后，更新数据库中的相关参数
		updateData := map[string]interface{}{
			"is_trading":         false,
			"status":             2, // 水位调节中
			"is_trading_blocked": false,
		}

		if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
			Where("id = ?", instance.Id).
			Updates(updateData).Error; err != nil {
			e.Log.Errorf("更新实例参数失败：%s", err)
			return err
		}

		e.Log.Infof("实例：%s 参数已成功更新", instance.Id)
		return nil

	}

	return nil
}

func CheckIsTradeBlockedByRiskControl(instanceId int) (bool, error) {
	db := sdk.Runtime.GetDbByKey("*")

	err := db.Model(&models.BusDexCexTriangularObserver{}).
		Where("id = ?", instanceId).
		First(&models.BusDexCexTriangularObserver{}).Error
	if err != nil {
		log.Errorf("查询实例失败, 异常：%s \r\n", err)
		return true, err
	}

	var riskEvents []models.BusRiskEvent
	err = db.Model(models.BusRiskEvent{}).
		Where("strategy_id = ?", common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE).
		Where("strategy_instance_id =?", instanceId).
		Where("is_recovered =?", 0).
		Where("risk_level >= ?", common.RISK_LEVEL_MIDDLE).
		Find(&riskEvents).Error
	if err != nil {
		log.Errorf("查询风控事件失败, 异常：%s \r\n", err)
		return true, err
	}
	if len(riskEvents) > 0 {
		return true, nil
	}
	return false, nil
}

func SendWarningNotification(symbol, instanceID, traderID, triggerCondition, triggerValue, currentValue string, larkClient lark.LarkRobotAlert) error {
	message := fmt.Sprintf(`
	⚠️ 风控预警通知
		策略实例ID: %s
		交易币对: %s
		触发风控的交易ID: %s
		触发条件: %s %s
		当前值: %s
		通知时间: %s

	🔔 该预警不会影响交易，仅供参考。请关注交易风险。
	`, instanceID, symbol, traderID, triggerCondition, triggerValue, currentValue, time.Now().Format("2006-01-02 15:04:05"))

	return larkClient.SendLarkAlert(message)
}

// SendMiddleNotification 发送风控中级别通知
func SendMiddleNotification(symbol, instanceID, traderID, triggerCondition, triggerValue, currentValue, recoveryMethod string, larkClient lark.LarkRobotAlert) error {
	message := fmt.Sprintf(`
	🚨 风控触发：暂停 %s 交易
		策略实例ID: %s
		触发风控的交易ID: %s
		触发条件: %s %s
		当前值: %s
		恢复方式: %s
		通知时间: %s

	❗ 请立即检查策略，并决定是否手动恢复交易。
	`, symbol, instanceID, traderID, triggerCondition, triggerValue, currentValue, recoveryMethod, time.Now().Format("2006-01-02 15:04:05"))

	return larkClient.SendLarkAlert(message)
}

// SendHigestNotification 发送风控最高级别通知
func SendHighestNotification(traderID, triggerCondition, triggerValue, currentValue, recoveryMethod string, larkClient lark.LarkRobotAlert) error {
	message := fmt.Sprintf(`
	🛑 交易系统已全局暂停
		触发风控的交易ID: %s
		触发条件: %s %s
		当前值: %s
		恢复方式: %s
		通知时间: %s

	🚨 全局交易已暂停，请立即检查风险，并决定恢复方案。
	`, traderID, triggerCondition, triggerValue, currentValue, recoveryMethod, time.Now().Format("2006-01-02 15:04:05"))

	return larkClient.SendLarkAlert(message)
}

func SendRiskCheckFailMessage(riskRule string, larkClient lark.LarkRobotAlert) error {
	message := fmt.Sprintf(`
	❌ 风控校验失败
		风控规则: %s
		通知时间: %s
		请检查风控规则并及时修正。
	`, riskRule, time.Now().Format("2006-01-02 15:04:05"))
	return larkClient.SendLarkAlert(message)
}
