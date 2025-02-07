package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
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
	tx.Where("status = ?", 1) //默认只查运行中的
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
		observerId := (*list)[i].ObserverId // 使用 (*list)[i] 访问原始元素
		state, err := client.GetObserverState(observerId)
		if err != nil {
			e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
			continue
		}
		e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", observerId, state)
		buyOnDex := state.GetBuyOnDex()
		cexSellPrice, dexBuyPrice := e.calculate_dex_cex_price(buyOnDex)
		e.Log.Infof("[buy on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexSellPrice, dexBuyPrice)

		sellOnDex := state.GetSellOnDex()
		cexBuyPrice, dexSellPrice := e.calculate_dex_cex_price(sellOnDex)
		e.Log.Infof("[sell on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexBuyPrice, dexSellPrice)

		if cexSellPrice-dexBuyPrice > 0 {
			//获取最新的价差记录统计信息，设置价差持续时间
			dexBuyData := models.BusDexCexPriceSpreadStatistics{}
			err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				(*list)[i].DexBuyDiffDuration = "0"
			}
			if err != nil {
				e.Log.Errorf("db error:%s", err)
				return err
			}
			(*list)[i].DexBuyDiffDuration = dexBuyData.Duration
		}
		if dexSellPrice-cexBuyPrice > 0 {
			dexSellData := models.BusDexCexPriceSpreadStatistics{}
			err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				(*list)[i].DexBuyDiffDuration = "0"
			}
			if err != nil {
				e.Log.Errorf("db error:%s", err)
				return err
			}
			(*list)[i].DexSellDiffDuration = dexSellData.Duration
		}
		(*list)[i].ProfitOfBuyOnDex = strconv.FormatFloat(*buyOnDex.ProfitFiatAmount, 'f', 6, 64)
		(*list)[i].ProfitOfSellOnDex = strconv.FormatFloat(*sellOnDex.ProfitFiatAmount, 'f', 6, 64)
		(*list)[i].CexSellPrice = strconv.FormatFloat(cexSellPrice, 'f', 6, 64)
		(*list)[i].DexBuyPrice = strconv.FormatFloat(dexBuyPrice, 'f', 6, 64)
		(*list)[i].DexBuyDiffPrice = strconv.FormatFloat(cexSellPrice-dexBuyPrice, 'f', 6, 64)
		(*list)[i].CexBuyPrice = strconv.FormatFloat(cexBuyPrice, 'f', 6, 64)
		(*list)[i].DexSellPrice = strconv.FormatFloat(dexSellPrice, 'f', 6, 64)
		(*list)[i].DexSellDiffPrice = strconv.FormatFloat(dexSellPrice-cexBuyPrice, 'f', 6, 64)
	}

	return nil
}

func (e *BusDexCexTriangularObserver) calculate_dex_cex_price(priceState *pb.ArbitrageState) (float64, float64) {
	var cexPrice float64      // TRUMP/USDT
	var dexPrice float64      //TRUMP/USDT
	var cexQuotePrice float64 // 例如：TRUMP/USDT
	if priceState.CexBaseQuantity != nil && priceState.CexBaseFiatAmount != nil && *priceState.CexBaseQuantity != 0 {
		cexQuotePrice = *priceState.CexBaseFiatAmount / *priceState.CexBaseQuantity
	} else {
		// 处理 nil 或除数为 0 的情况，避免 panic
		cexQuotePrice = 0
	}

	var cexSolPrice float64 //SOL/USDT
	if priceState.CexSolQuantity != nil && priceState.CexSolFiatAmount != nil && *priceState.CexSolQuantity != 0 {
		cexSolPrice = *priceState.CexSolFiatAmount / *priceState.CexSolQuantity
	} else {
		// 处理 nil 或除数为 0 的情况，避免 panic
		cexSolPrice = 0
	}

	//if cexQuotePrice != 0 && cexSolPrice != 0 {
	//	cexPrice = cexQuotePrice / cexSolPrice
	//} else {
	//	// 处理除数为0的情况，避免panic
	//	cexPrice = 0
	//}
	cexPrice = cexQuotePrice

	var dexSolPrice float64 //TRUMP/WSOL
	if priceState.DexSolAmount != nil && priceState.DexBaseAmount != nil && *priceState.CexSolQuantity != 0 {
		dexSolPrice = *priceState.DexSolAmount / *priceState.DexBaseAmount
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
	observerId := model.ObserverId
	state, err := client.GetObserverState(observerId)
	if err != nil {
		e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
		return nil
	}
	e.Log.Infof("get state for observerId:%d \r\n state: %+v \r\n", observerId, state)
	buyOnDex := state.GetBuyOnDex()
	cexSellPrice, dexBuyPrice := e.calculate_dex_cex_price(buyOnDex)
	e.Log.Infof("[buy on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexSellPrice, dexBuyPrice)

	sellOnDex := state.GetSellOnDex()
	cexBuyPrice, dexSellPrice := e.calculate_dex_cex_price(sellOnDex)
	e.Log.Infof("[sell on dex price details]: cexPrice: %+v , dexPrice: %+v \r\n", cexBuyPrice, dexSellPrice)

	model.ProfitOfBuyOnDex = strconv.FormatFloat(*buyOnDex.ProfitFiatAmount, 'f', 6, 64)
	model.ProfitOfSellOnDex = strconv.FormatFloat(*sellOnDex.ProfitFiatAmount, 'f', 6, 64)
	model.CexSellPrice = strconv.FormatFloat(cexSellPrice, 'f', 6, 64)
	model.DexBuyPrice = strconv.FormatFloat(dexBuyPrice, 'f', 6, 64)
	model.DexBuyDiffPrice = strconv.FormatFloat(cexSellPrice-dexBuyPrice, 'f', 6, 64)
	if cexSellPrice-dexBuyPrice > 0 {
		//获取最新的价差记录统计信息，设置价差持续时间
		dexBuyData := models.BusDexCexPriceSpreadStatistics{}
		err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
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
		err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			model.DexBuyDiffDuration = "0"
		}
		if err != nil {
			e.Log.Errorf("db error:%s", err)
			return err
		}
		model.DexSellDiffDuration = dexSellData.Duration
	}
	model.CexBuyPrice = strconv.FormatFloat(cexBuyPrice, 'f', 6, 64)
	model.DexSellPrice = strconv.FormatFloat(dexSellPrice, 'f', 6, 64)
	model.DexSellDiffPrice = strconv.FormatFloat(dexSellPrice-cexBuyPrice, 'f', 6, 64)

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
	e.Log.Infof("e[BatchInsert], data: %+v", c)
	e.Log.Infof("e[SlippageBps], data: %s", *c.SlippageBps)
	baseTokens := c.BaseToken
	if len(baseTokens) == 0 {
		return errors.New("empty baseTokens")
	}

	var successStartedCount int
	for _, baseToken := range baseTokens {
		//循环创建监听

		var ammConfig = pb.DexConfig{}
		var amberConfig = pb.AmberConfig{}
		var arbitrageConfig = pb.ArbitrageConfig{}
		c.GenerateAmmConfig(&ammConfig)
		c.GenerateAmberConfig(&amberConfig)
		c.GenerateArbitrageConfig(&arbitrageConfig)
		e.Log.Infof("e[BatchInsert], baseToken: %s, ammConfig: %+v", baseToken, ammConfig)
		observerId, err := client.StartNewObserver(&amberConfig, &ammConfig, &arbitrageConfig)
		if err != nil {
			e.Log.Errorf("Service BatchInsert error:%s \r\n", err)
			continue
		}
		c.Generate(&data, baseToken, observerId)
		err = e.Orm.Create(&data).Error
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

	observerId := d.ObserverId
	err := client.StopObserver(observerId)
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
