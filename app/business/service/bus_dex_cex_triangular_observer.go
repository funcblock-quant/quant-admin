package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	waterLevelPb "quanta-admin/app/grpc/proto/client/water_level_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	"quanta-admin/common/global"
	common "quanta-admin/common/models"
	lark "quanta-admin/common/notification"
	"quanta-admin/common/utils"
	ext "quanta-admin/config"
	"sort"
	"strconv"
	"strings"
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

const (
	INSTANCE_STATUS_CREATED    = 0
	INSTANCE_STATUS_OBSERVE    = 1
	INSTANCE_STATUS_WATERLEVEL = 2
	INSTANCE_STATUS_TRADING    = 3
)

// GetPage 获取BusDexCexTriangularObserver列表
func (e *BusDexCexTriangularObserver) GetPage(c *dto.BusDexCexTriangularObserverGetPageReq, p *actions.DataPermission, list *[]dto.BusDexCexTriangularObserverGetPageResp, count *int64) error {
	var err error
	var data models.BusDexCexTriangularObserver
	tx := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	statuses := []int{INSTANCE_STATUS_OBSERVE, INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING} // 需要查询的多个状态
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
			cexSellPrice, dexBuyPrice = e.calculateDexCexPrice(buyOnDex, true)
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
			cexBuyPrice, dexSellPrice = e.calculateDexCexPrice(sellOnDex, false)
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

func (e *BusDexCexTriangularObserver) calculateDexCexPrice(priceState *pb.ObserverState, isDexBuy bool) (float64, float64) {
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
	cexSellPrice, dexBuyPrice := e.calculateDexCexPrice(buyOnDex, true)
	e.Log.Infof("[buy on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexSellPrice, dexBuyPrice)

	sellOnDex := state.GetSellOnDex()
	cexBuyPrice, dexSellPrice := e.calculateDexCexPrice(sellOnDex, false)
	e.Log.Infof("[sell on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexBuyPrice, dexSellPrice)

	buyOnDexProfit := *buyOnDex.CexSellQuoteAmount - *buyOnDex.CexBuyQuoteAmount
	sellOnDexProfit := *sellOnDex.CexSellQuoteAmount - *sellOnDex.CexBuyQuoteAmount

	model.ProfitOfBuyOnDex = buyOnDexProfit
	model.ProfitOfSellOnDex = sellOnDexProfit
	model.CexSellPrice = cexSellPrice
	model.DexBuyPrice = dexBuyPrice
	model.DexBuyDiffPrice = cexSellPrice - dexBuyPrice
	// if cexSellPrice-dexBuyPrice > 0 {
	// 	//获取最新的价差记录统计信息，设置价差持续时间
	// 	dexBuyData := models.BusDexCexPriceSpreadStatistics{}
	// 	err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
	// 	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 		model.DexBuyDiffDuration = "0"
	// 	}
	// 	if err != nil {
	// 		e.Log.Errorf("db error:%s", err)
	// 		return err
	// 	}
	// 	model.DexBuyDiffDuration = dexBuyData.Duration
	// }
	// if dexSellPrice-cexBuyPrice > 0 {
	// 	dexSellData := models.BusDexCexPriceSpreadStatistics{}
	// 	err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", id, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
	// 	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 		model.DexBuyDiffDuration = "0"
	// 	}
	// 	if err != nil {
	// 		e.Log.Errorf("db error:%s", err)
	// 		return err
	// 	}
	// 	model.DexSellDiffDuration = dexSellData.Duration
	// }
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
		Where("status IN ?", []int{INSTANCE_STATUS_OBSERVE, INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING}).
		Debug().Find(list).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetSymbolList error:%s \r\n", err)
		return err
	}
	return nil
}

// GetLatestObserverConfig 获取币种的最近一次观察配置
func (e *BusDexCexTriangularObserver) GetLatestObserverConfig(req *dto.BusDexCexTriangularGetLatestObserverConfigReq, p *actions.DataPermission, resp *models.BusDexCexTriangularObserver) error {
	var err error
	var data models.BusDexCexTriangularObserver

	err = e.Orm.Model(&data).
		Unscoped().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("target_token = ?", req.Token).
		Order("created_at desc").
		Debug().First(resp).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetLatestObserverConfig error:%s \r\n", err)
		return nil
	}
	return nil
}

// GetExchangeList 获取BusDexCexTriangularObserver所有币种列表
func (e *BusDexCexTriangularObserver) GetExchangeList(p *actions.DataPermission, list *[]dto.DexCexTriangularObserverExchangeListResp) error {
	var err error
	var data models.BusDexCexTriangularObserver

	err = e.Orm.Model(&data).
		Unscoped().
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Select("exchange_type as exchange").
		Group("exchange_type").
		Where("status IN ?", []int{INSTANCE_STATUS_OBSERVE, INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING}).
		Debug().Find(list).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetSymbolList error:%s \r\n", err)
		return err
	}
	return nil
}

// GetDexWalletList 获取BusDexCexTriangularObserver所有dex钱包列表
func (e *BusDexCexTriangularObserver) GetDexWalletList(p *actions.DataPermission, list *[]models.BusDexWallet) error {
	var err error
	var data models.BusDexWallet

	err = e.Orm.Model(&data).
		Unscoped().
		Find(list).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetDexWalletList error:%s \r\n", err)
		return err
	}

	return nil
}

// GetCexAccountList 获取BusDexCexTriangularObserver所有cex账户列表
func (e *BusDexCexTriangularObserver) GetCexAccountList(req *dto.BusGetCexAccountListReq, p *actions.DataPermission, list *[]models.BusExchangeAccountInfo) error {
	var err error
	var data models.BusExchangeAccountInfo

	//查出绑定了amber的交易所账户列表，别且状态是已启用
	err = e.Orm.Model(&data).
		Unscoped().
		Where("exchange_type = ? and status = ?", req.Exchange, 2).
		Find(list).Error

	if err != nil {
		e.Log.Errorf("BusDexCexTriangularObserverService GetCexAccountList error:%s \r\n", err)
		return err
	}

	return nil
}

// GetBoundAccountList 根据某一侧的账户信息获取BusDexCexTriangularObserver所有已绑定的账号
func (e *BusDexCexTriangularObserver) GetBoundAccountList(req *dto.BusGetBoundAccountReq, p *actions.DataPermission, resp *dto.BusGetBoundAccountResp) error {
	var err error
	var boundIds []string

	e.Log.Infof("get request: %v", req)
	if req.AccountType == "Cex" {
		// 根据cex账户，查询绑定的dex账户列表
		err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
			Select("dex_wallet_id").
			Where("status > ? and cex_account_id =?", INSTANCE_STATUS_OBSERVE, req.AccountId).
			Pluck("dex_wallet_id", &boundIds).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetBoundAccountList error:%s \r\n", err)
			return err
		}

		if len(boundIds) == 0 {
			return nil
		}

		var dexWalletList []models.BusDexWallet
		// 查询出绑定的dex账户列表
		err = e.Orm.Model(&models.BusDexWallet{}).
			Where("id in (?)", boundIds).
			Find(&dexWalletList).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetBoundAccountList error:%s \r\n", err)
			return err
		}
		e.Log.Infof("get dex wallet list: %v", dexWalletList)
		resp.DexWalletList = dexWalletList
		return nil

	} else {
		// 根据dex账户，查询绑定的cex账户列表
		err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
			Select("cex_account_id").
			Where("status > ? and dex_wallet_id =?", INSTANCE_STATUS_OBSERVE, req.AccountId).
			Pluck("cex_account_id", &boundIds).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetBoundAccountList error:%s \r\n", err)
			return err
		}

		if len(boundIds) == 0 {
			return nil
		}

		var cexAccountList []models.BusExchangeAccountInfo
		// 查询出绑定的dex账户列表
		err = e.Orm.Model(&models.BusExchangeAccountInfo{}).
			Where("id in (?)", boundIds).
			Find(&cexAccountList).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetBoundAccountList error:%s \r\n", err)
			return err
		}
		e.Log.Infof("get cex account list: %v", cexAccountList)
		resp.CexAccountList = cexAccountList
		return nil

	}
}

// GetBoundAccountList 根据某一侧的账户信息获取BusDexCexTriangularObserver可绑定的账号
func (e *BusDexCexTriangularObserver) GetCanBoundAccountList(req *dto.BusGetBoundAccountReq, p *actions.DataPermission, resp *dto.BusGetBoundAccountResp) error {
	var err error
	var boundIds []string

	e.Log.Infof("get request: %v", req)
	if req.AccountType == "Cex" {
		// 查询未绑定的dex账户列表
		err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
			Select("dex_wallet_id").
			Where("status > ? ", INSTANCE_STATUS_OBSERVE).
			Pluck("dex_wallet_id", &boundIds).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetCanBoundAccountList error:%s \r\n", err)
			return err
		}
		var dexWalletList []models.BusDexWallet
		// 查询出绑定的dex账户列表
		if len(boundIds) == 0 {
			err = e.Orm.Model(&models.BusDexWallet{}).
				Find(&dexWalletList).Error
		} else {
			err = e.Orm.Model(&models.BusDexWallet{}).
				Where("id not in (?)", boundIds).
				Find(&dexWalletList).Error
		}

		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetCanBoundAccountList error:%s \r\n", err)
			return err
		}
		e.Log.Infof("get un-bound dex wallet list: %v", dexWalletList)
		resp.DexWalletList = dexWalletList
		return nil

	} else {
		// 根据dex账户，查询绑定的cex账户列表
		err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
			Select("cex_account_id").
			Where("status > ?", INSTANCE_STATUS_OBSERVE).
			Pluck("cex_account_id", &boundIds).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetCanBoundAccountList error:%s \r\n", err)
			return err
		}
		e.Log.Infof("get bound dex wallet ids: %v", boundIds)
		var cexAccountList []models.BusExchangeAccountInfo
		// 查询出绑定的cex账户列表
		if len(boundIds) == 0 {
			err = e.Orm.Model(&models.BusExchangeAccountInfo{}).
				Find(&cexAccountList).Error
		} else {
			err = e.Orm.Model(&models.BusExchangeAccountInfo{}).
				Where("id not in (?)", boundIds).
				Find(&cexAccountList).Error
		}
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService GetCanBoundAccountList error:%s \r\n", err)
			return err
		}
		e.Log.Infof("get un-bound cex account list: %v", cexAccountList)
		resp.CexAccountList = cexAccountList
		return nil

	}
}

// GetActiveAccountPairs 获取dex cex套利的交易中的账户对，以及账户是否配置全局水位调节
func (e *BusDexCexTriangularObserver) GetActiveAccountPairs(p *actions.DataPermission, resp *[]dto.BusAccountPairInfo) error {
	var err error

	//1. 统计出所有的账户组合，循环处理获取账户组合的全局水位调整参数
	var accountPairs []DexCexPair
	err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Select("DISTINCT dex_wallet_id, cex_account_id").
		Where("dex_wallet_id IS NOT NULL AND cex_account_id IS NOT NULL").
		Where("status > ?", INSTANCE_STATUS_OBSERVE).
		Find(&accountPairs).Error

	if err != nil {
		e.Log.Errorf("获取账户组合失败:%s \r\n", err)
		return err
	}

	if len(accountPairs) == 0 {
		return nil
	}

	for _, accountPair := range accountPairs {
		var dexWallet models.BusDexWallet
		err = e.Orm.Model(&models.BusDexWallet{}).
			Where("id = ?", accountPair.DexWalletId).
			First(&dexWallet).Error

		if err != nil {
			e.Log.Errorf("获取dex钱包失败:%s \r\n", err)
			return err
		}

		var cexAccount models.BusExchangeAccountInfo
		err = e.Orm.Model(&models.BusExchangeAccountInfo{}).
			Where("id = ?", accountPair.CexAccountId).
			First(&cexAccount).Error
		if err != nil {
			e.Log.Errorf("获取cex账户失败:%s \r\n", err)
			return err
		}

		accountPairInfo := dto.BusAccountPairInfo{
			DexwalletId:    dexWallet.Id,
			CexAccountId:   cexAccount.Id,
			CexAccountName: cexAccount.AccountName,
			DexWalletName:  dexWallet.WalletName,
			CexAccountUid:  cexAccount.Uid,
			DexWalletAddr:  dexWallet.WalletAddress,
		}

		solWaterLevelConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_SOLANA_WATER_LEVEL_KEY, accountPair.DexWalletId, accountPair.CexAccountId)
		stableWaterLevelConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY, accountPair.DexWalletId, accountPair.CexAccountId)
		//step 2 : 封装全局的水位调节启动结果以及配置到响应体
		var solanaConfig models.BusCommonConfig
		err = e.Orm.Model(&models.BusCommonConfig{}).
			Where("category = ? and config_key = ?", common.WATER_LEVEL, solWaterLevelConfigKey).
			Order("created_at desc").
			First(&solanaConfig).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Errorf("获取solana全局水位调节配置失败:%s \r\n", err)
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到配置，则表示没有配置全局水位调节
			e.Log.Info("未配置solana全局水位调节:%s \r\n")
			accountPairInfo.HasGlobalConfig = false
			*resp = append(*resp, accountPairInfo)
			continue
		}

		accountPairInfo.HasGlobalConfig = true

		configJsonStr := solanaConfig.ConfigJson
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
			accountPairInfo.SolanaConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{
				AlertThreshold:             &alertThreshold,
				BuyTriggerThreshold:        &buyTriggerThreshold,
				SellTriggerThreshold:       &sellTriggerThreshold,
				MinDepositAmountThreshold:  &solMinDepositAmountThreshold,
				MinWithdrawAmountThreshold: &solMinWithdrawAmountThreshold,
			}
		}

		// 稳定币的全局水位配置
		var stableConfig models.BusCommonConfig
		err = e.Orm.Model(&models.BusCommonConfig{}).
			Where("category = ? and config_key = ?", common.WATER_LEVEL, stableWaterLevelConfigKey).
			Order("created_at desc").
			First(&stableConfig).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("获取solana全局水位调节配置失败:%s \r\n", err)
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到配置，则表示没有启用全局水位调节
			accountPairInfo.HasGlobalConfig = false
			*resp = append(*resp, accountPairInfo)
			continue
		}

		configJsonStr = stableConfig.ConfigJson

		// 解析 JSON
		err = json.Unmarshal([]byte(configJsonStr), &configMap)
		if err != nil {
			e.Log.Error("JSON 解析失败:", err)
		} else {
			alertThreshold := configMap["alertThreshold"].(float64)
			accountPairInfo.StableCoinConfig = &dto.BusDexCexTriangularUpdateWaterLevelParamsReq{
				AlertThreshold: &alertThreshold,
			}
		}
		*resp = append(*resp, accountPairInfo)
	}

	e.Log.Infof("get account pairs: %v", *resp)
	return nil

}

// GetRealtimeInterestRate 获取dex cex套利的交易中的账户对的实时汇率
func (e *BusDexCexTriangularObserver) GetRealtimeInterestRate(req *dto.BusGetInterestRateReq, p *actions.DataPermission, resp *dto.BusGetInterestRateResp) error {
	var err error

	var cexAccount models.BusExchangeAccountInfo
	err = e.Orm.Model(&models.BusExchangeAccountInfo{}).
		Where("id = ?", req.CexAccountId).
		First(&cexAccount).Error
	if err != nil {
		e.Log.Errorf("get cex account info failed, cexAccountId: %s", req.CexAccountId)
		return err
	}
	e.Log.Infof("获取到交易所账号Id: %v", *resp)

	// 获取实时汇率
	interestRateReq := &waterLevelPb.GetInterestRatesRequest{
		Currencies: []string{req.Currency},
	}

	exchangeType := req.ExchangeType
	if exchangeType == "Binance" {
		interestRateReq.ExchangeType = waterLevelPb.ExchangeType_Binance
	} else if exchangeType == "GateIO" {
		interestRateReq.ExchangeType = waterLevelPb.ExchangeType_Gate
	} else {
		e.Log.Errorf("不支持的交易所类型: %s", exchangeType)
		return errors.New("不支持的交易所类型")
	}

	//不传主账号信息, 不传dex信息
	secretConfig, err := generateSecretConfig(models.BusDexWallet{}, cexAccount, models.BusExchangeAccountInfo{})
	if err != nil {
		e.Log.Errorf("获取秘钥配置失败:%s \r\n", err)
		return err
	}
	interestRateReq.SecretKey = secretConfig
	interestRateResp, err := client.GetInterestRates(interestRateReq)

	if err != nil {
		e.Log.Errorf("获取实时汇率失败:%s \r\n", err)
		return err
	}
	if interestRateResp == nil {
		e.Log.Errorf("获取实时汇率失败，返回值为空")
		return nil
	}

	for _, rate := range interestRateResp.List {
		resp.InterestRate = rate.InterestRate
		resp.Currency = rate.Currency
		// 这里我们系统只支持一个币种的实时汇率
		break
	}

	return nil

}

// GetWaterLevelDetail 获取水位调节详情
func (e *BusDexCexTriangularObserver) GetWaterLevelDetail(req *dto.BusDexCexTriangularGetWaterLevelDetailReq, p *actions.DataPermission, resp *dto.BusDexCexTriangularGetWaterLevelDetailResp) error {
	var err error

	instanceId := req.InstanceId
	queryReq := &waterLevelPb.InstanceId{
		InstanceId: strconv.Itoa(instanceId),
	}
	waterLevelState, err := client.GetWaterLevelInstanceState(queryReq)
	if err != nil {
		e.Log.Errorf("get WaterLevelInstanceState error: instanceId:%d, error msg:%s \r\n", instanceId, err)
		return err
	}

	e.Log.Infof("获取到水位调节实例状态: %v", waterLevelState)

	resp.InstanceId = instanceId
	taskState := waterLevelState.InstanceTaskState
	resp.TaskType = taskState.TaskType
	resp.TaskStatus = taskState.TaskStatus
	resp.TaskStep = taskState.TaskStep
	resp.TaskError = taskState.TaskError
	resp.InstanceError = waterLevelState.InstanceError

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

		tx.Model(&data).Where("id = ?", data.Id).Update("status", INSTANCE_STATUS_OBSERVE)
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

	if config.ApplicationConfig.Mode != "dev" {
		instanceId := strconv.Itoa(d.Ids)
		err := client.StopArbitragerClient(instanceId)
		if err != nil {
			e.Log.Errorf("暂停监视器失败 error:%s \r\n", err)
			return err
		}

		e.Log.Infof("grpc请求暂停监视器成功")
	}

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
	data.CexAccountId = &c.CexAccount
	data.DexWalletId = &c.DexWallet
	// step1 先启动水位调节实例
	if config.ApplicationConfig.Mode != "dev" {
		err = StartTokenWaterLevelWithCheckExists(&data)
		if err != nil {
			return err
		}
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

	priorityFee := *c.PriorityFee
	//err = StartTrader(c, priorityFee, jitoFee, slippageBpsUint, data, exchangeType, e)
	//if err != nil {
	//	return err
	//}

	// 启动水位调节后，更新数据库中的相关参数
	updateData := map[string]interface{}{
		"is_trading":                    false,
		"alert_threshold":               c.AlertThreshold,
		"buy_trigger_threshold":         c.BuyTriggerThreshold,
		"sell_trigger_threshold":        c.SellTriggerThreshold,
		"min_deposit_amount_threshold":  c.MinDepositAmountThreshold,
		"min_withdraw_amount_threshold": c.MinWithdrawAmountThreshold,
		"slippage_bps_rate":             c.SlippageBpsRate,
		"prefer_jito":                   c.PreferJito,
		"priority_fee":                  priorityFee,
		"jito_fee_rate":                 c.JitoFeeRate,
		"status":                        INSTANCE_STATUS_WATERLEVEL, // 水位调节中
		"cex_account_id":                c.CexAccount,
		"dex_wallet_id":                 c.DexWallet,
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
		//status = 2 表示水位调节中
		Where("status = ?", INSTANCE_STATUS_WATERLEVEL).
		Find(&list).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	for _, instance := range list {
		instanceId := instance.Id
		queryReq := &waterLevelPb.InstanceId{
			InstanceId: strconv.Itoa(instanceId),
		}
		waterLevelState, err := client.GetWaterLevelInstanceState(queryReq)
		if err != nil {
			e.Log.Errorf("get WaterLevelInstanceState error: instanceId:%d, error msg:%s \r\n", instance.Id, err)
			continue
		}

		traderSwitch := waterLevelState.TraderSwitch
		waterLevelStatus := waterLevelState.WaterLevelStatus

		if waterLevelStatus == 0 && traderSwitch {
			e.Log.Infof("waterlevel state for instancId: %d is: success", instanceId)
			e.Log.Infof("currency: %s, cex spot balance:%s, cex margin balance:%s, dex balance: %s", waterLevelState.Currency, waterLevelState.SpotAccountBalance, waterLevelState.MarginAccountBalance, waterLevelState.ChainWalletBalance)
			if !instance.IsTrading {
				//如果原先交易没开，则开启交易功能
				err = DoStartTrader(e.Orm, &instance)
				if err != nil {
					e.Log.Errorf("start trader error:%s \r\n", err)
					return err
				}
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
		Where("status = ?", INSTANCE_STATUS_TRADING).
		Find(&list).Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	for _, instance := range list {
		instanceId := strconv.Itoa(instance.Id)
		queryReq := &waterLevelPb.InstanceId{
			InstanceId: instanceId,
		}
		waterLevelState, err := client.GetWaterLevelInstanceState(queryReq)
		if err != nil {
			e.Log.Errorf("get WaterLevelInstanceState error: instanceId:%d, error msg:%s \r\n", instance.Id, err)
			continue
		}

		traderSwitch := waterLevelState.TraderSwitch
		waterLevelStatus := waterLevelState.WaterLevelStatus

		e.Log.Infof("currency: %s, cex spot balance:%s, cex margin balance:%s, dex balance: %s", waterLevelState.Currency, waterLevelState.SpotAccountBalance, waterLevelState.MarginAccountBalance, waterLevelState.ChainWalletBalance)
		if waterLevelStatus == 1 {
			// 水位调节中，需要更新状态
			//关闭交易功能
			// 暂停交易成功后，更新状态
			updateData := map[string]interface{}{
				"status": INSTANCE_STATUS_WATERLEVEL, // 水位调节中
			}

			if !traderSwitch {
				e.Log.Infof("trader switch is false, stop trader")
				err = client.DisableTrader(instanceId)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%s 交易功能失败，异常：%s \r\n", instanceId, err)
					return err
				}
				updateData["is_trading"] = false
			}

			if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
				Where("id = ?", instance.Id).
				Updates(updateData).Error; err != nil {
				e.Log.Errorf("更新实例参数失败：%s", err)
				continue
			}
			e.Log.Infof("实例：%s 参数已成功更新", data.InstanceId)
		} else {
			// 水位调节状态为非调节中状态，此时如果trader switch 为关，需要暂停交易，并更新状态为水位调节中
			if !traderSwitch {
				//如果此时交易开关是关闭的，说明一种情况，sol全局水位调节达到了最低的阈值。此时修改为水位调节中，并且暂停交易
				e.Log.Infof("trader switch is false, stop trader")
				err = client.DisableTrader(instanceId)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%s 交易功能失败，异常：%s \r\n", instanceId, err)
					return err
				}

				updateData := map[string]interface{}{
					"status":     INSTANCE_STATUS_WATERLEVEL, // 水位调节中
					"is_trading": false,
				}

				if err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
					Where("id = ?", instance.Id).
					Updates(updateData).Error; err != nil {
					e.Log.Errorf("更新实例参数失败：%s", err)
					continue
				}
			}

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

		stopReq := &waterLevelPb.InstanceId{
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
		"status":             INSTANCE_STATUS_OBSERVE,
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
		SlippageRate: c.SlippageBpsRate,
		PriorityFee:  c.PriorityFee,
		JitoFeeRate:  c.JitoFeeRate,
		PreferJito:   &c.PreferJito,
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
		"priority_fee":      c.PriorityFee,
		"jito_fee_rate":     *c.JitoFeeRate,
		"prefer_jito":       c.PreferJito,
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

// StopAllTrader 暂停所有trader
func (e *BusDexCexTriangularObserver) StopAllTrader(operator int) ([]string, error) {
	var data models.BusDexCexTriangularObserver
	failedList := make([]string, 0)

	//查出所有水位调节中和交易中的instance
	// 查询出所有交易开启中以及水位调节中的实例
	var instances []models.BusDexCexTriangularObserver
	err := e.Orm.Model(models.BusDexCexTriangularObserver{}).
		Where("status IN ?", []int{INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING}).
		Find(&instances).Error
	if err != nil {
		e.Log.Errorf("查询实例失败:%s \r\n", err)
		return nil, err
	}

	for _, instance := range instances {
		id := strconv.Itoa(instance.Id)
		if instance.Status == strconv.Itoa(INSTANCE_STATUS_WATERLEVEL) {
			// 水位调节中，暂停水位调节，更新状态为1

			stopWaterLevelReq := &waterLevelPb.InstanceId{
				InstanceId: id,
			}
			err = client.StopWaterLevelInstance(stopWaterLevelReq)
			if err != nil {
				e.Log.Errorf("grpc 暂停停水位调节实例：:%d 失败，异常：%s \r\n", instance.Id, err)
				failedList = append(failedList, id)
				continue
			}

			updateData := map[string]interface{}{
				"status": INSTANCE_STATUS_OBSERVE,
			}
			err := e.Orm.Model(&data).
				Where("id =?", instance.Id).
				Updates(updateData).Error
			if err != nil {
				e.Log.Errorf("更新实例:%d 状态为 开启观察 失败，异常信息：%s \r\n", instance.Id, err)
				failedList = append(failedList, id)
				continue
			}
			e.Log.Infof("instanceId：%d 已暂停水位调节功能", instance.Id)
		} else if instance.Status == strconv.Itoa(INSTANCE_STATUS_TRADING) {
			// 交易中，暂停所有trader，更新状态为4

			stopReq := dto.BusDexCexTriangularObserverStopTraderReq{
				InstanceId: instance.Id,
			}
			stopReq.UpdateBy = operator

			err = e.StopTrader(&stopReq, false)
			if err != nil {
				e.Log.Errorf("grpc 暂停trader实例：:%d 失败，异常：%s \r\n", instance.Id, err)
				failedList = append(failedList, id)
				continue
			}
		}
	}
	if len(failedList) > 0 {
		e.Log.Warnf("停止trader功能部分成功")
		return failedList, errors.New("一键暂停部分成功")
	}

	e.Log.Infof("实例：%d 参数已成功更新", data.Id)

	return nil, nil
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
	err = DoUpdateTokenWaterLevel(&data)
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
func (e *BusDexCexTriangularObserver) GetGlobalWaterLevelState(req *dto.BusGetCexExchangeConfigListReq) (*dto.BusDexCexTriangularGlobalWaterLevelStateResp, error) {
	var data models.BusCommonConfig
	exchangeType := req.Exchange
	if exchangeType == "" {
		e.Log.Infof("exchange is empty, return default config")
		return nil, nil
	}

	resp := &dto.BusDexCexTriangularGlobalWaterLevelStateResp{}
	//step 1 : 查询sol以及稳定币的水位调节是否启动
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		//return nil, errors.New("水位调节服务不可用，请稍后刷新重试")
	} else {
		log.Infof("waterLevelInstances:%+v\n", waterLevelInstances)
		for _, instanceId := range waterLevelInstances.InstanceIds {
			if instanceId == "SOLANA_"+exchangeType {
				// solana 已经启动水位调节
				resp.SolWaterLevelState = true
			} else if instanceId == "USDT_"+exchangeType {
				// 稳定币 已经启动水位调节
				resp.StableCoinWaterLevelState = true
			}
		}
	}

	// 如果水位调节服务挂了，返回对应的错误，给到前端水位调节不可用的提示之类的。
	//step 2 : 封装全局的水位调节启动结果以及配置到响应体
	err = e.Orm.Model(&models.BusCommonConfig{}).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY+"_"+exchangeType).
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
			// alertThreshold := configMap["alertThreshold"].(float64)
			// buyTriggerThreshold := configMap["buyTriggerThreshold"].(float64)
			// sellTriggerThreshold := configMap["sellTriggerThreshold"].(float64)
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
				// AlertThreshold:             &alertThreshold,
				// BuyTriggerThreshold:        &buyTriggerThreshold,
				// SellTriggerThreshold:       &sellTriggerThreshold,
				MinDepositAmountThreshold:  &solMinDepositAmountThreshold,
				MinWithdrawAmountThreshold: &solMinWithdrawAmountThreshold,
			}
		}
	}

	var stableData models.BusCommonConfig
	err = e.Orm.Model(&stableData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY+"_"+exchangeType).
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
	exchangeType := req.ExchangeType

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
		Where("status = ? AND is_trading = ?", INSTANCE_STATUS_TRADING, true).
		Distinct().
		Pluck("quote_token", &quoteTokens).
		Error
	if err != nil {
		e.Log.Errorf("获取实例失败:%s \r\n", err)
		return err
	}

	isSolanaStarted := false
	for _, instanceId := range waterLevelInstances.InstanceIds {
		if instanceId == "SOLANA_"+exchangeType {
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

		secretKey := &waterLevelPb.SecretKey{}

		clientRequest := &waterLevelPb.StartInstanceRequest{
			InstanceId:           "SOLANA_" + exchangeType,
			ExchangeType:         exchangeType,
			Currency:             "SOL",
			CurrencyType:         0, // token
			TokenThresholdConfig: tokenConfig,
			SecretKey:            secretKey,
		}

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
			InstanceId:           "SOLANA_" + exchangeType,
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
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY+"_"+exchangeType).
		Order("created_at desc").
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局solana水位调节参数")
			// 如果不存在，则新增
			data = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  common.GLOBAL_SOLANA_WATER_LEVEL_KEY + "_" + exchangeType,
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
				InstanceId:                quoteToken + "_" + exchangeType,
				ExchangeType:              exchangeType,
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
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY+"_"+exchangeType).
		Order("created_at desc").
		First(&stableData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局稳定币水位调节参数")
			// 如果不存在，则新增
			stableData = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY + "_" + exchangeType,
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

// UpdateGlobalWaterLevelConfigV2 更新全局WaterLevel 参数
func (e *BusDexCexTriangularObserver) UpdateGlobalWaterLevelConfigV2(req *dto.BusDexCexTriangularUpdateGlobalWaterLevelConfigReq) error {
	var data models.BusCommonConfig
	cexAccountId := req.CexAccountId
	dexWalletId := req.DexWalletId

	solanaConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_SOLANA_WATER_LEVEL_KEY, dexWalletId, cexAccountId)
	stableConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY, dexWalletId, cexAccountId)

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

	// 保存配置到数据库
	err = e.Orm.Model(&data).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, solanaConfigKey).
		Order("created_at desc").
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前账户组不存在全局solana水位调节参数")
			// 如果不存在，则新增
			data = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  solanaConfigKey,
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

	var stableData models.BusCommonConfig
	// 稳定币水位调节参数处理
	err = e.Orm.Model(&stableData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, stableConfigKey).
		Order("created_at desc").
		First(&stableData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Info("当前不存在全局稳定币水位调节参数")
			// 如果不存在，则新增
			stableData = models.BusCommonConfig{
				Category:   common.WATER_LEVEL,
				ConfigKey:  stableConfigKey,
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
		Where("status = ? AND is_trading = ?", INSTANCE_STATUS_TRADING, true).
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

		_, err = client.StartWaterLevelInstance(clientRequest)
		if err != nil {
			e.Log.Errorf("启动solana全局水位调节失败:%s \r\n", err)
			return err
		}
		e.Log.Infof("启动solana全局水位调节启动成功")
	} else {
		//如果已经启动了，要尝试更新
		// 先获取当前实例的参数
		solanaStateReq := &waterLevelPb.InstanceId{
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
			stableStateReq := &waterLevelPb.InstanceId{
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

type DexCexPair struct {
	DexWalletId  int `json:"dexWalletId"`
	CexAccountId int `json:"cexAccountId"`
}

// StartGlobalWaterLevelV2 启动全局水位调整功能
func (e *BusDexCexTriangularObserver) StartGlobalWaterLevelV2() error {
	var data models.BusDexCexTriangularObserver

	// reserveRatio := 5 // 资金准备金倍率，默认5

	// 1. 获取solana价格
	// 获取一个 observer
	db := e.Orm
	cexSolPrice, err := getSolPriceByObserver(db, data)
	if err != nil {
		e.Log.Errorf("获取solana价格异常:%s \r\n", err)
		return err
	}

	e.Log.Infof("获取到solana价格：%f \r\n", cexSolPrice)

	//solana的阈值通过公式计算
	//1. 统计出所有的账户组合，循环处理每个账户组合的全局水位调整
	var accountPairs []DexCexPair
	err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Select("DISTINCT dex_wallet_id, cex_account_id").
		Where("dex_wallet_id IS NOT NULL AND cex_account_id IS NOT NULL").
		Where("status > ?", INSTANCE_STATUS_CREATED).
		Find(&accountPairs).Error

	if err != nil {
		e.Log.Errorf("获取账户组合失败:%s \r\n", err)
		return err
	}

	var solMinDepositAmountThreshold, solMinWithdrawAmountThreshold float64

	var solData models.BusCommonConfig
	err = e.Orm.Model(&solData).
		Where("category = ? and config_key = ?", common.WATER_LEVEL, common.GLOBAL_SOLANA_WATER_LEVEL_KEY).
		Order("created_at desc").
		First(&solData).Error

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

	if len(accountPairs) == 0 {
		e.Log.Warn("没有账户组合，跳过全局水位调节任务")
		return nil
	}

	for _, accountPair := range accountPairs {
		//2. 每个账户组合，计算出一套阈值标准
		var solanaAlertThreshold, solBuyTriggerThreshold, solSellTriggerThreshold float64
		err := calculateGlobalSolWalterLevelConfig(db, accountPair, cexSolPrice, &solanaAlertThreshold, &solBuyTriggerThreshold, &solSellTriggerThreshold)
		if err != nil {
			e.Log.Errorf("获取账户组合 %d,%d 全局水位调节参数失败, 跳过本次启动全局水位调节任务", accountPair.DexWalletId, accountPair.CexAccountId)
			continue
		}

		// 启动solana水位调节
		err = e.startGlobalSolanaWaterLevelForAccountPair(db, accountPair, solanaAlertThreshold, solBuyTriggerThreshold, solSellTriggerThreshold, solMinDepositAmountThreshold, solMinWithdrawAmountThreshold)
		if err != nil {
			e.Log.Error("failed to start global solana water level for dexWalletId:%d, cexAccountId:%d ", accountPair.DexWalletId, accountPair.CexAccountId)
			continue
		}

	}

	// 遍历水位调节实例，对于不存在绑定关系的实例进行暂停
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		return err
	}

	for _, instanceId := range waterLevelInstances.InstanceIds {
		if strings.Contains(instanceId, "SOLANA") {
			dexWalletId, cexAccountId, err := parseSolanaWaterLevelInstanceId(instanceId)
			if err != nil {
				fmt.Println("解析失败:", err)
				continue
			}

			var traders []models.BusDexCexTriangularObserver

			err = db.Model(&models.BusDexCexTriangularObserver{}).
				Where("dex_wallet_id =? AND cex_account_id =?", dexWalletId, cexAccountId).
				Where("status = ? AND is_trading = ?", INSTANCE_STATUS_TRADING, true).
				Find(&traders).Error
			if err != nil {
				e.Log.Errorf("查询trader failed, dexWalletId: %d, cexAccountId: %d", dexWalletId, cexAccountId)
				continue
			}
			if len(traders) < 1 {
				e.Log.Infof("instanceId: %s has no bound account relation, stop it", instanceId)
				stopReq := &waterLevelPb.InstanceId{
					InstanceId: instanceId,
				}
				err = client.StopWaterLevelInstance(stopReq)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instanceId, err)
					continue
				}
			}
		}
	}

	return nil
}

// StartGlobalWaterLevelV3 启动全局水位调整功能
func (e *BusDexCexTriangularObserver) StartGlobalWaterLevelV3() error {

	// 获取一个 observer
	db := e.Orm
	//solana的阈值通过公式计算
	//1. 统计出所有的账户组合，循环处理每个账户组合的全局水位调整
	var accountPairs []DexCexPair
	err := e.Orm.Model(&models.BusDexCexTriangularObserver{}).
		Select("DISTINCT dex_wallet_id, cex_account_id").
		Where("dex_wallet_id IS NOT NULL AND cex_account_id IS NOT NULL").
		Where("status > ?", INSTANCE_STATUS_CREATED).
		Find(&accountPairs).Error

	if err != nil {
		e.Log.Errorf("获取账户组合失败:%s \r\n", err)
		return err
	}

	if len(accountPairs) == 0 {
		e.Log.Warn("没有账户组合，跳过全局水位调节任务")
		return nil
	}

	for _, accountPair := range accountPairs {
		//2. 每个账户组合，计算出一套阈值标准
		solanaConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_SOLANA_WATER_LEVEL_KEY, accountPair.DexWalletId, accountPair.CexAccountId)
		stableConfigKey := generateGlobalWaterLevelConfigKey(common.GLOBAL_STABLE_COIN_WATER_LEVEL_KEY, accountPair.DexWalletId, accountPair.CexAccountId)

		var solData models.BusCommonConfig
		err = e.Orm.Model(&solData).
			Where("category = ? and config_key = ?", common.WATER_LEVEL, solanaConfigKey).
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
			Where("category = ? and config_key = ?", common.WATER_LEVEL, stableConfigKey).
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

		// 启动solana水位调节
		err = e.startGlobalSolanaWaterLevelForAccountPair(db, accountPair, solanaAlertThreshold, solBuyTriggerThreshold, solSellTriggerThreshold, solMinDepositAmountThreshold, solMinWithdrawAmountThreshold)
		if err != nil {
			e.Log.Error("failed to start global solana water level for dexWalletId:%d, cexAccountId:%d ", accountPair.DexWalletId, accountPair.CexAccountId)
			continue
		}
		// TODO 启动稳定币水位调节
		e.Log.Infof("启动稳定币 %s 全局水位调节 stableCoinAlertThreshold: %v \r\n", stableCoinAlertThreshold)

	}

	// 遍历水位调节实例，对于不存在绑定关系的实例进行暂停
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		return err
	}

	for _, instanceId := range waterLevelInstances.InstanceIds {
		if strings.Contains(instanceId, "SOLANA") {
			dexWalletId, cexAccountId, err := parseSolanaWaterLevelInstanceId(instanceId)
			if err != nil {
				fmt.Println("解析失败:", err)
				continue
			}

			var traders []models.BusDexCexTriangularObserver

			err = db.Model(&models.BusDexCexTriangularObserver{}).
				Where("dex_wallet_id =? AND cex_account_id =?", dexWalletId, cexAccountId).
				Where("status = ? AND is_trading = ?", INSTANCE_STATUS_TRADING, true).
				Find(&traders).Error
			if err != nil {
				e.Log.Errorf("查询trader failed, dexWalletId: %d, cexAccountId: %d", dexWalletId, cexAccountId)
				continue
			}
			if len(traders) < 1 {
				e.Log.Infof("instanceId: %s has no bound account relation, stop it", instanceId)
				stopReq := &waterLevelPb.InstanceId{
					InstanceId: instanceId,
				}
				err = client.StopWaterLevelInstance(stopReq)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instanceId, err)
					continue
				}
			}
		}
	}

	return nil
}

// 为每一个账户组合启动一个solana水位调节
func (e *BusDexCexTriangularObserver) startGlobalSolanaWaterLevelForAccountPair(db *gorm.DB, accountPair DexCexPair, solanaAlertThreshold, solBuyTriggerThreshold, solSellTriggerThreshold, solMinDepositAmountThreshold, solMinWithdrawAmountThreshold float64) error {
	waterLevelInstances, err := client.ListWaterLevelInstance()
	if err != nil {
		e.Log.Errorf("获取水位调节实例失败, %s", err)
		return err
	}

	isSolanaStarted := false
	instanceIdKey := generateSolanaWaterLevelInstanceId(accountPair)

	for _, instanceId := range waterLevelInstances.InstanceIds {
		if instanceId == instanceIdKey {
			// 该账号组solana 已经启动水位调节
			isSolanaStarted = true
			break
		}
	}

	var cexAccount models.BusExchangeAccountInfo
	err = db.Model(&models.BusExchangeAccountInfo{}).
		Where("id = ?", accountPair.CexAccountId).
		First(&cexAccount).Error
	if err != nil {
		e.Log.Errorf("get cex account info failed, cexAccountId: %s", accountPair.CexAccountId)
		return err
	}

	var cexMasterAccount models.BusExchangeAccountInfo
	if cexAccount.MasterAccountId == 0 {
		e.Log.Infof("cexAccountId: %d 没有绑定主账户", accountPair.CexAccountId)
		cexMasterAccount = models.BusExchangeAccountInfo{}
	} else {
		err = db.Model(&models.BusExchangeAccountInfo{}).
			Where("id = ?", cexAccount.MasterAccountId).
			First(&cexMasterAccount).Error
		if err != nil {
			e.Log.Errorf("get cex master account info failed, cexMasterAccountId: %s", cexAccount.MasterAccountId)
			return err
		}
	}

	var dexWallet models.BusDexWallet
	err = db.Model(&models.BusDexWallet{}).
		Where("id = ?", accountPair.DexWalletId).
		First(&dexWallet).Error
	if err != nil {
		e.Log.Errorf("get dex wallet info failed, dexWalletId: %s", accountPair.DexWalletId)
		return err
	}

	if !isSolanaStarted {
		secretKey, err := generateSecretConfig(dexWallet, cexAccount, cexMasterAccount)
		if err != nil {
			return err
		}

		tokenConfig := &waterLevelPb.TokenThresholdConfig{
			AlertThreshold:             strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64),
			BuyTriggerThreshold:        strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64),
			SellTriggerThreshold:       strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64),
			MinDepositAmountThreshold:  strconv.FormatFloat(solMinDepositAmountThreshold, 'f', -1, 64),
			MinWithdrawAmountThreshold: strconv.FormatFloat(solMinWithdrawAmountThreshold, 'f', -1, 64),
		}

		exchangeType := cexAccount.ExchangeType
		if exchangeType == global.EXCHANGE_TYPE_GATEIO {
			exchangeType = "Gate"
		}

		clientRequest := &waterLevelPb.StartInstanceRequest{
			InstanceId:           instanceIdKey,
			ExchangeType:         exchangeType,
			Currency:             "SOL",
			CurrencyType:         0, // token
			TokenThresholdConfig: tokenConfig,
			SecretKey:            secretKey,
		}

		_, err = client.StartWaterLevelInstance(clientRequest)
		if err != nil {
			e.Log.Errorf("启动solana全局水位调节失败:%s \r\n", err)
			return err
		}
		e.Log.Infof("启动solana全局水位调节启动成功")
	} else {
		//如果已经启动了，要尝试更新
		// 先获取当前实例的参数
		solanaStateReq := &waterLevelPb.InstanceId{
			InstanceId: instanceIdKey,
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
				InstanceId:           instanceIdKey,
				CurrencyType:         0, // token
				TokenThresholdConfig: tokenConfig,
			}
			client.UpdateWaterLevelInstance(updateReq)
		} else {
			e.Log.Infof("从服务端获取到solana水位调节实例参数：%v \n", solanaState)
			oldParams := solanaState.InstanceParams.TokenThresholdConfig
			if oldParams.AlertThreshold == strconv.FormatFloat(solanaAlertThreshold, 'f', -1, 64) &&
				oldParams.BuyTriggerThreshold == strconv.FormatFloat(solBuyTriggerThreshold, 'f', -1, 64) &&
				oldParams.SellTriggerThreshold == strconv.FormatFloat(solSellTriggerThreshold, 'f', -1, 64) {
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

				log.Infof("账户组合 %d,%d 全局水位调节参数：solanaAlertThreshold: %f, solBuyTriggerThreshold: %f, solSellTriggerThreshold: %f \n", accountPair.DexWalletId, accountPair.CexAccountId, tokenConfig.AlertThreshold, tokenConfig.BuyTriggerThreshold, tokenConfig.SellTriggerThreshold)
				updateReq := &waterLevelPb.UpdateInstanceParamsRequest{
					InstanceId:           instanceIdKey,
					CurrencyType:         0, // token
					TokenThresholdConfig: tokenConfig,
				}
				client.UpdateWaterLevelInstance(updateReq)
			}
		}
	}
	return nil
}

func generateSolanaWaterLevelInstanceId(accountPair DexCexPair) string {
	return fmt.Sprintf("%d_%d_SOLANA", accountPair.DexWalletId, accountPair.CexAccountId)
}

func parseSolanaWaterLevelInstanceId(instanceId string) (dexWalletId, cexAccountId int, err error) {
	// 按下划线分割字符串
	parts := strings.Split(instanceId, "_")
	if len(parts) != 3 || parts[2] != "SOLANA" {
		return 0, 0, fmt.Errorf("无效的 instanceId 格式: %s", instanceId)
	}

	// 解析 dexWalletId 和 cexAccountId
	dexWalletId, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("解析 dexWalletId 失败: %v", err)
	}

	cexAccountId, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("解析 cexAccountId 失败: %v", err)
	}

	return dexWalletId, cexAccountId, nil
}

func calculateGlobalSolWalterLevelConfig(db *gorm.DB, accountPair DexCexPair, cexSolPrice float64, solanaAlertThreshold *float64, solBuyTriggerThreshold *float64, solSellTriggerThreshold *float64) error {
	var traders []models.BusDexCexTriangularObserver
	err := db.Model(&models.BusDexCexTriangularObserver{}).
		Where("dex_wallet_id =? AND cex_account_id =?", accountPair.DexWalletId, accountPair.CexAccountId).
		Where("status = ? AND is_trading = ?", INSTANCE_STATUS_TRADING, true).
		Find(&traders).Error
	if err != nil {
		log.Errorf("获取 trader 失败:%s \r\n", err)
		return err
	}

	if len(traders) == 0 {
		log.Infof("没有找到 trader，跳过")
		return errors.New("没有对应的trader")
	}

	var solMaxTradeVolume float64
	for _, trader := range traders {
		solAmount := *trader.MaxQuoteAmount / cexSolPrice
		solMaxTradeVolume += solAmount
	}
	solMaxTradeVolume = math.Floor(solMaxTradeVolume)

	multiplier := 1.0
	if len(traders) > 1 {
		multiplier = 0.8
	}

	*solanaAlertThreshold = solMaxTradeVolume * multiplier
	// 这里使用3，7的系数，是经过计算，当solana价格只要不高于200，可以避免因为n的波动导致的频繁触发充提，因为可以保证两档之间的中间值在相邻档位的阈值范围内。避免
	*solBuyTriggerThreshold = solMaxTradeVolume * 3 * multiplier
	*solSellTriggerThreshold = solMaxTradeVolume * 7 * multiplier
	log.Infof("账户组合 %d,%d 全局水位调节参数：solanaAlertThreshold: %f, solBuyTriggerThreshold: %f, solSellTriggerThreshold: %f \n", accountPair.DexWalletId, accountPair.CexAccountId, *solanaAlertThreshold, *solBuyTriggerThreshold, *solSellTriggerThreshold)
	return nil
}

func getSolPriceByObserver(db *gorm.DB, data models.BusDexCexTriangularObserver) (float64, error) {
	err := db.Model(&models.BusDexCexTriangularObserver{}).
		Where("status > ?", INSTANCE_STATUS_CREATED).
		First(&data).Error

	if err != nil {
		log.Errorf("获取一个 observer 失败:%s \r\n", err)
		return 0, err
	}

	id := strconv.Itoa(data.Id)
	state, err := client.GetObserverState(id)
	if err != nil {
		log.Errorf("获取 solana 状态失败:%s \r\n", err)
		return 0, err
	}

	buyOnDex := state.GetBuyOnDex()
	sellOnDex := state.GetSellOnDex()

	cexSolPrice := 0.0
	if buyOnDex != nil {
		cexSolPrice = CalculateCexSolPrice(state.BuyOnDex, true)
	}
	if sellOnDex != nil {
		cexSolPrice = CalculateCexSolPrice(state.SellOnDex, false)
	}

	if cexSolPrice <= 0 {
		log.Error("获取到错误的sol价格，跳过")
		return 0, errors.New("error sol price")
	}
	return cexSolPrice, nil
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
	// 单币种最大日亏损阈值
	var symbolDailyMaxLossThreshold []interface{}

	if v, ok := configMap["symbolDailyMaxLossThreshold"]; ok {
		if list, valid := v.([]interface{}); valid {
			symbolDailyMaxLossThreshold = list
		}
	}

	// 单笔最大亏损比例阈值
	var relativeLossThreshold []interface{}

	if v, ok := configMap["relativeLossThreshold"]; ok {
		if list, valid := v.([]interface{}); valid {
			relativeLossThreshold = list
		}
	}

	// 排序，按照 action 从大到小排序
	sort.Slice(absoluteLossThreshold, func(i, j int) bool {
		return absoluteLossThreshold[i].(map[string]interface{})["action"].(float64) > absoluteLossThreshold[j].(map[string]interface{})["action"].(float64)
	})

	sort.Slice(relativeLossThreshold, func(i, j int) bool {
		return relativeLossThreshold[i].(map[string]interface{})["action"].(float64) > relativeLossThreshold[j].(map[string]interface{})["action"].(float64)
	})

	sort.Slice(symbolDailyMaxLossThreshold, func(i, j int) bool {
		return symbolDailyMaxLossThreshold[i].(map[string]interface{})["action"].(float64) > symbolDailyMaxLossThreshold[j].(map[string]interface{})["action"].(float64)
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
		e.Log.Warn("[Risk Control Check] 风控校验正在进行中，跳过")
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
		e.Log.Debugf("[Risk Control Check] 交易订单:%d \r\n", trade.Id)
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
		afterAction, err = e.SymbolDailyMaxLossThresholdCheck(symbolDailyMaxLossThreshold, trade, *larkClient)
		if err != nil {
			SendRiskCheckFailMessage(common.TRIGGER_RULE_SYMBOL_DAILY_MAX_LOSS_THRESHOLD, *larkClient)
			e.Log.Errorf("[Risk Control Check] 单币种最大日亏损阈值校验失败:%s \r\n", err)
			errOccurred = err
			break
		}
		if afterAction > maxAfterAction {
			maxAfterAction = afterAction
		}

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

func (e BusDexCexTriangularObserver) SymbolDailyMaxLossThresholdCheck(symbolDailyMaxLossThreshold []interface{}, trade dto.StrategyDexCexTriangularArbitrageTradesGetPageResp, larkClient lark.LarkRobotAlert) (int, error) {
	afterAction := 0
	for _, threshold := range symbolDailyMaxLossThreshold {

		thresholdMap := threshold.(map[string]interface{})
		thresholdValue := thresholdMap["threshold"].(float64)
		action := thresholdMap["action"].(float64)
		actionDetail := thresholdMap["actionDetail"].(map[string]interface{})

		symbol := trade.Symbol

		dailyProfit, err := GetDailySymbolProfit(e.Orm, symbol)
		if err != nil {
			e.Log.Errorf("[Risk Control Check] 获取%s每日收益数据失败:%s \r\n", symbol, err)
			return 0, errors.New("获取每日收益数据失败")
		}

		dailyProfitStr := strconv.FormatFloat(dailyProfit, 'f', -1, 64)

		if dailyProfit < -thresholdValue {
			// 亏损超过阈值
			// 生成风控事件
			riskEvent := models.BusRiskEvent{
				StrategyId:         common.STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE,
				StrategyInstanceId: trade.InstanceId,
				TradeId:            trade.Id,
				AssetSymbol:        trade.Symbol,
				TriggerRule:        common.TRIGGER_RULE_SYMBOL_DAILY_MAX_LOSS_THRESHOLD,
				TriggerValue:       dailyProfitStr,
			}
			if action == 3 {
				// 暂停全部交易，单笔交易亏损触发暂不支持暂停全部交易行为
				e.Log.Errorf("[Risk Control Check] 单币种交易日亏损触发暂不支持暂停全部交易行为 \r\n")
				return 0, errors.New("单币种交易日亏损触发暂不支持暂停全部交易行为")
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

				var recoverMethod string
				if manualRecover == 1 {
					recoverMethod = "手动恢复"
				} else {
					recoverMethod = "自动恢复"
				}

				SendMiddleNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_SYMBOL_DAILY_MAX_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), dailyProfitStr, recoverMethod, larkClient)

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

				SendWarningNotification(trade.Symbol, trade.InstanceId, strconv.Itoa(trade.Id), common.TRIGGER_RULE_SYMBOL_DAILY_MAX_LOSS_THRESHOLD, strconv.FormatFloat(thresholdValue, 'f', -1, 64), dailyProfitStr, larkClient)
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
		Where("status IN ?", []int{INSTANCE_STATUS_WATERLEVEL, INSTANCE_STATUS_TRADING}).
		Find(&instances).Error
	if err != nil {
		e.Log.Errorf("查询实例失败:%s \r\n", err)
		return err
	}

	if len(highestRiskEvents) == 0 {
		e.Log.Debugf("当前不存在未恢复的风控事件 \r\n")
	} else {
		e.Log.Debugf("存在全局中断交易的事件,暂停全部实例交易功能")

		// 关闭所有实例
		for _, instance := range instances {
			stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
				InstanceId: instance.Id,
			}

			if instance.Status == strconv.Itoa(INSTANCE_STATUS_WATERLEVEL) { //水位调节中
				stopReq := &waterLevelPb.InstanceId{
					InstanceId: strconv.Itoa(instance.Id),
				}
				err = client.StopWaterLevelInstance(stopReq)
				if err != nil {
					e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instance.Id, err)
					continue
				}
				// 更新observer的status =1
				updateData := map[string]interface{}{
					"status":             INSTANCE_STATUS_OBSERVE,
					"is_trading_blocked": true,
				}

				err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
					Where("id = ?", instance.Id).
					Updates(updateData).Error
				if err != nil {
					e.Log.Errorf("更新实例状态失败, 异常：%s \r\n", err)
					continue
				}

			} else if instance.Status == strconv.Itoa(INSTANCE_STATUS_TRADING) { //交易开启中
				err = e.StopTrader(&stopTradeReq, true)
				if err != nil {
					e.Log.Errorf("关闭实例%d失败:%s \r\n", instance.Id, err)
				} else {
					e.Log.Infof("关闭实例%d成功 \r\n", instance.Id)
				}
			}
			continue
		}
		e.Log.Info("完成暂停全部实例交易功能")
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
		e.Log.Debug("当前不存在未恢复的单币种风控事件 \r\n")
		return nil
	}

	// 暂停所有单币种风控事件对应的实例
	for _, instance := range instances {
		e.Log.Debug("当前存在未恢复的单币种风控事件 \r\n")
		for _, riskEvent := range singleTokenRiskEvents {
			if riskEvent.StrategyInstanceId == strconv.Itoa(instance.Id) {
				stopTradeReq := dto.BusDexCexTriangularObserverStopTraderReq{
					InstanceId: instance.Id,
				}

				if instance.Status == strconv.Itoa(INSTANCE_STATUS_WATERLEVEL) { //水位调节中
					stopReq := &waterLevelPb.InstanceId{
						InstanceId: strconv.Itoa(instance.Id),
					}
					err = client.StopWaterLevelInstance(stopReq)
					if err != nil {
						e.Log.Errorf("grpc暂停实例：:%d 水位调节功能失败，异常：%s \r\n", instance.Id, err)
						continue
					}
					// 更新observer的status =1
					updateData := map[string]interface{}{
						"status":             INSTANCE_STATUS_OBSERVE,
						"is_trading_blocked": true,
					}

					err = e.Orm.Model(&models.BusDexCexTriangularObserver{}).
						Where("id = ?", instance.Id).
						Updates(updateData).Error
					if err != nil {
						e.Log.Errorf("更新实例状态失败, 异常：%s \r\n", err)
						continue
					}

				} else if instance.Status == strconv.Itoa(INSTANCE_STATUS_TRADING) { //交易开启中
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

	// 查询出所有被风控事件blocking的实例
	var instances []models.BusDexCexTriangularObserver
	err = e.Orm.Model(models.BusDexCexTriangularObserver{}).
		Where("status = ?", INSTANCE_STATUS_OBSERVE).
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
			"status":             INSTANCE_STATUS_WATERLEVEL, // 水位调节中
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

func DoStartObserver(observer *models.BusDexCexTriangularObserver) error {

	maxArraySize := new(uint32)
	*maxArraySize = uint32(observer.MaxArraySize) //默认5， clmm使用参数

	dexConfig := &pb.DexConfig{}
	if observer.DexType == global.DEX_TYPE_RAY_AMM {
		dexConfig.Config = &pb.DexConfig_RayAmm{
			RayAmm: &pb.RayAmmConfig{
				Pool:      observer.AmmPoolId,
				TokenMint: observer.TokenMint,
			},
		}
	} else if observer.DexType == global.DEX_TYPE_RAY_CLMM {
		dexConfig.Config = &pb.DexConfig_RayClmm{
			RayClmm: &pb.RayClmmConfig{
				Pool:         observer.AmmPoolId,
				TokenMint:    observer.TokenMint,
				MaxArraySize: maxArraySize,
			},
		}
	} else if observer.DexType == global.DEX_TYPE_ORCA_WHIRL_POOL {
		dexConfig.Config = &pb.DexConfig_OrcaWhirlPool{
			OrcaWhirlPool: &pb.OrcaWhirlPoolConfig{
				Pool:      observer.AmmPoolId,
				TokenMint: observer.TokenMint,
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
	amberConfig.ExchangeType = &observer.ExchangeType
	amberConfig.TakerFee = proto.Float64(*observer.TakerFee)

	amberConfig.TargetToken = &observer.TargetToken
	amberConfig.QuoteToken = &observer.QuoteToken

	amberConfig.BidDepth = proto.Int32(int32(20))
	amberConfig.AskDepth = proto.Int32(int32(20))

	if observer.AskDepth != "" {
		depthInt, err := strconv.Atoi(observer.AskDepth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		amberConfig.AskDepth = proto.Int32(int32(depthInt))
	}

	if observer.BidDepth != "" {
		depthInt, err := strconv.Atoi(observer.BidDepth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		amberConfig.BidDepth = proto.Int32(int32(depthInt))
	}

	instanceId := strconv.Itoa(observer.Id)
	log.Infof("restart observer success with params: dexConfig: %+v\n, arbitrageConfig: %+v\n", dexConfig, arbitrageConfig)
	err := client.StartNewArbitragerClient(&instanceId, amberConfig, dexConfig, arbitrageConfig)
	if err != nil {
		log.Errorf("restart observer throw grpc error: %v\n", err)
		return err
	}
	return nil
}

func DoStartTrader(db *gorm.DB, instance *models.BusDexCexTriangularObserver) error {

	exchangeType, err := instance.GetExchangeTypeForTrader()
	if err != nil {
		log.Errorf("获取ExchangeType参数异常，:%s \r\n", err)
		return err
	}
	var dexWallet models.BusDexWallet
	var cexAccount models.BusExchangeAccountInfo

	err = db.Model(&dexWallet).
		Where("id = ?", instance.DexWalletId).
		First(&dexWallet).Error

	if err != nil {
		log.Errorf("查询DEX钱包参数失败:%s \r\n", err)
		return err
	}

	err = db.Model(&cexAccount).
		Where("id =?", instance.CexAccountId).
		First(&cexAccount).Error
	if err != nil {
		log.Errorf("查询CEX账户参数失败:%s \r\n", err)
		return err
	}

	privateKey, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), dexWallet.EncryptedPrivateKey)
	if err != nil {
		log.Errorf("解密私钥参数失败:%s \r\n", err)
		return err
	}

	cexConfig := &pb.CexConfig{}
	if instance.ExchangeType == global.EXCHANGE_TYPE_BINANCE {
		amberAccountToken, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), cexAccount.AmberAccountToken)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return err
		}
		accountType := "Exchange"
		contractType := "Spot"
		amberTraderConfig := &pb.AmberTraderConfig{
			ExchangeType: &exchangeType,
			AccountType:  &accountType,
			ContractType: &contractType,
			AccountId:    &cexAccount.AmberAccountName,
			AccessToken:  &amberAccountToken,
		}
		cexConfig.Config = &pb.CexConfig_Amber{
			Amber: amberTraderConfig,
		}

	} else if instance.ExchangeType == global.EXCHANGE_TYPE_GATEIO {

		apiKey, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), cexAccount.EncryptedApiKey)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return err
		}
		secretKey, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), cexAccount.EncryptedApiSecret)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return err
		}
		accountType := "Unified"
		gateConfig := &pb.GateioTraderConfig{
			AccountType: &accountType,
			ApiKey:      &apiKey,
			ApiSecret:   &secretKey,
		}
		cexConfig.Config = &pb.CexConfig_Gateio{
			Gateio: gateConfig,
		}
	}

	jitoFee := instance.JitoFeeRate
	traderParams := &pb.TraderParams{
		//Slippage:    &slippageBpsFloat,
		SlippageRate: instance.SlippageBpsRate,
		PriorityFee:  instance.PriorityFee,
		JitoFeeRate:  jitoFee,
		PreferJito:   &instance.PreferJito,
	}
	swapperConfig := &pb.SwapperConfig{
		Trader: &privateKey,
	}

	if config.ApplicationConfig.Mode != "dev" {
		instanceId := strconv.Itoa(instance.Id)
		err := client.EnableTrader(instanceId, cexConfig, traderParams, swapperConfig)
		if err != nil {
			log.Errorf("GRPC 启动Trader for instanceId:%d 失败，异常:%s \r\n", instance.Id, err)
			return err
		}
		log.Errorf("GRPC 启动Trader for instanceId:%d 成功", instance.Id)
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
		err = DoUpdateTokenWaterLevel(observer)
		if err != nil {
			log.Errorf("grpc启动实例失败，异常信息:%s \r\n", err)
			return err
		}
	} else {
		//说明未启动实例，此时启动新的实例
		db := sdk.Runtime.GetDbByKey("*")
		err = DoStartTokenWaterLevel(db, observer)
		if err != nil {
			log.Errorf("grpc启动实例失败，异常信息:%s \r\n", err)
			return err
		}
	}
	return nil
}

// DoStartTokenWaterLevel 启动水位调节，不校验是否存在
func DoStartTokenWaterLevel(db *gorm.DB, observer *models.BusDexCexTriangularObserver) error {

	var dexWallet models.BusDexWallet
	var cexAccount models.BusExchangeAccountInfo

	err := db.Model(&dexWallet).
		Where("id = ?", observer.DexWalletId).
		First(&dexWallet).Error

	if err != nil {
		log.Errorf("查询DEX钱包参数失败:%s \r\n", err)
		return err
	}

	err = db.Model(&cexAccount).
		Where("id =?", observer.CexAccountId).
		First(&cexAccount).Error
	if err != nil {
		log.Errorf("查询CEX账户参数失败:%s \r\n", err)
		return err
	}

	var masterCexAccount models.BusExchangeAccountInfo
	var secretConfig *waterLevelPb.SecretKey
	if cexAccount.MasterAccountId == 0 {
		//如果没有设置主账号，则默认不传主账号信息
		secretConfig, err = generateSecretConfig(dexWallet, cexAccount, models.BusExchangeAccountInfo{})
		if err != nil {
			log.Errorf("生成SecretConfig参数失败:%s \r\n", err)
			return err
		}
	} else {
		err = db.Model(&cexAccount).
			Where("id =?", cexAccount.MasterAccountId).
			First(&masterCexAccount).Error
		if err != nil {
			log.Errorf("查询CEX账户参数失败:%s \r\n", err)
			return err
		}
		secretConfig, err = generateSecretConfig(dexWallet, cexAccount, masterCexAccount)
		if err != nil {
			log.Errorf("生成SecretConfig参数失败:%s \r\n", err)
			return err
		}
	}

	tokenConfig := &waterLevelPb.TokenThresholdConfig{
		AlertThreshold:             strconv.FormatFloat(*observer.AlertThreshold, 'f', -1, 64),
		BuyTriggerThreshold:        strconv.FormatFloat(*observer.BuyTriggerThreshold, 'f', -1, 64),
		SellTriggerThreshold:       strconv.FormatFloat(*observer.SellTriggerThreshold, 'f', -1, 64),
		MinDepositAmountThreshold:  strconv.FormatFloat(*observer.MinDepositAmountThreshold, 'f', -1, 64),
		MinWithdrawAmountThreshold: strconv.FormatFloat(*observer.MinWithdrawAmountThreshold, 'f', -1, 64),
	}

	exchangeType := observer.ExchangeType
	if observer.ExchangeType == global.EXCHANGE_TYPE_GATEIO {
		exchangeType = "Gate"
	}

	clientRequest := &waterLevelPb.StartInstanceRequest{
		InstanceId:           strconv.Itoa(observer.Id),
		ExchangeType:         exchangeType,
		Currency:             observer.TargetToken,
		CurrencyType:         0, // token
		PubKey:               observer.TokenMint,
		TokenThresholdConfig: tokenConfig,
		SecretKey:            secretConfig,
	}

	log.Infof("restart water level with req: %v \r\n", clientRequest)
	_, err = client.StartWaterLevelInstance(clientRequest)
	if err != nil {
		log.Errorf("启动水位调节失败:%s \r\n", err)
		return err
	}
	log.Infof("水位调节启动成功")
	return nil
}

func generateGlobalWaterLevelConfigKey(commonKey string, dexWalletId int, cexAccountId int) string {
	return commonKey + "_" + strconv.Itoa(int(dexWalletId)) + "_" + strconv.Itoa(int(cexAccountId))
}

func generateSecretConfig(dexWallet models.BusDexWallet, cexAccount models.BusExchangeAccountInfo, masterCexAccount models.BusExchangeAccountInfo) (*waterLevelPb.SecretKey, error) {
	var privateKey string
	var err error
	if dexWallet.EncryptedPrivateKey != "" {
		privateKey, err = utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), dexWallet.EncryptedPrivateKey)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return nil, err
		}
	}

	apiKey, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), cexAccount.EncryptedApiKey)
	if err != nil {
		log.Errorf("解密私钥参数失败:%s \r\n", err)
		return nil, err
	}
	secret, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), cexAccount.EncryptedApiSecret)
	if err != nil {
		log.Errorf("解密私钥参数失败:%s \r\n", err)
		return nil, err
	}

	var masterConfig *waterLevelPb.ExchangeAccount
	if masterCexAccount.Id != 0 {
		masterApiKey, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), masterCexAccount.EncryptedApiKey)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return nil, err
		}
		masterSecret, err := utils.DecryptWithSecretKey([]byte(ext.ExtConfig.Aes.Key), masterCexAccount.EncryptedApiSecret)
		if err != nil {
			log.Errorf("解密私钥参数失败:%s \r\n", err)
			return nil, err
		}

		masterConfig = &waterLevelPb.ExchangeAccount{
			ApiKey: masterApiKey,
			Secret: masterSecret,
		}
	}

	traderConfig := &waterLevelPb.ExchangeAccount{
		ApiKey:      apiKey,
		Secret:      secret,
		Passphrase:  cexAccount.Passphrase,
		AccountName: cexAccount.AccountName,
		Uid:         cexAccount.Uid,
		Email:       cexAccount.Email,
	}

	secretConfig := &waterLevelPb.SecretKey{
		TraderAccount:         traderConfig,
		MasterAccount:         masterConfig,
		ChainWalletPrivateKey: privateKey,
	}
	return secretConfig, nil
}

// DoUpdateTokenWaterLevel 更新水位调节
func DoUpdateTokenWaterLevel(observer *models.BusDexCexTriangularObserver) error {
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

// GetDailySymbolProfit 获取特定币种当日总利润
func GetDailySymbolProfit(db *gorm.DB, symbol string) (float64, error) {
	var totalProfit float64

	today := time.Now().Format("2006-01-02")
	startOfDay := today + " 00:00:00"
	endOfDay := today + " 23:59:59"

	err := db.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("COALESCE(SUM(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount), 0) AS total_profit").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("opportunities.cex_target_asset = ?", symbol).
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success =?", 1, 1, 1).
		Where("strategy_dex_cex_triangular_arbitrage_trades.created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Scan(&totalProfit).Error

	if err != nil {
		log.Errorf("获取特定币种当日总利润失败, 异常：%s \r\n", err)
		return 0, err
	}
	return totalProfit, nil
}

// GetAllDailyProfit 获取全币种币种当日总利润
func GetAllDailyProfit(db *gorm.DB) (float64, error) {
	var totalProfit float64

	today := time.Now().Format("2006-01-02")
	startOfDay := today + " 00:00:00"
	endOfDay := today + " 23:59:59"

	err := db.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("COALESCE(SUM(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount), 0) AS total_profit").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success =?", 1, 1, 1).
		Where("strategy_dex_cex_triangular_arbitrage_trades.created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Scan(&totalProfit).Error

	if err != nil {
		log.Errorf("获取全币种当日总利润失败, 异常：%s \r\n", err)
		return 0, err
	}
	return totalProfit, nil
}

// GetTotalSymbolProfitBeforeToday 获取特定币种当日之前的总利润
func GetTotalSymbolProfitBeforeToday(db *gorm.DB, symbol string) (float64, error) {
	var totalProfit float64

	today := time.Now().Format("2006-01-02")
	startOfDay := today + " 00:00:00"

	err := db.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("COALESCE(SUM(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount), 0) AS total_profit").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("opportunities.cex_target_asset = ?", symbol).
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success =?", 1, 1, 1).
		Where("strategy_dex_cex_triangular_arbitrage_trades.created_at BETWEEN < ?", startOfDay).
		Scan(&totalProfit).Error

	if err != nil {
		log.Errorf("获取特定币种当日之前总利润失败, 异常：%s \r\n", err)
		return 0, err
	}
	return totalProfit, nil
}

// GetAllTotalProfitBeforeToday 获取全币种币种当日之前的总利润
func GetAllTotalProfitBeforeToday(db *gorm.DB) (float64, error) {
	var totalProfit float64

	today := time.Now().Format("2006-01-02")
	startOfDay := today + " 00:00:00"

	err := db.Model(&models.StrategyDexCexTriangularArbitrageTrades{}).
		Select("COALESCE(SUM(strategy_dex_cex_triangular_arbitrage_trades.cex_sell_quote_amount - strategy_dex_cex_triangular_arbitrage_trades.cex_buy_quote_amount), 0) AS total_profit").
		Joins("LEFT JOIN strategy_dex_cex_triangular_arbitrage_opportunities AS opportunities ON strategy_dex_cex_triangular_arbitrage_trades.opportunity_id = opportunities.opportunity_id").
		Where("strategy_dex_cex_triangular_arbitrage_trades.dex_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_buy_success = ? AND strategy_dex_cex_triangular_arbitrage_trades.cex_sell_success =?", 1, 1, 1).
		Where("strategy_dex_cex_triangular_arbitrage_trades.created_at < ?", startOfDay).
		Scan(&totalProfit).Error

	if err != nil {
		log.Errorf("获取全币种当日之前总利润失败, 异常：%s \r\n", err)
		return 0, err
	}
	return totalProfit, nil
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
