package service

import (
	"errors"
	"math"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	"strconv"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
)

type BusDexCexPriceSpreadData struct {
	service.Service
}

// GetPage 获取BusDexCexPriceSpreadData列表
func (e *BusDexCexPriceSpreadData) GetPage(c *dto.BusDexCexPriceSpreadDataGetPageReq, p *actions.DataPermission, list *[]models.BusDexCexPriceSpreadData, count *int64) error {
	var err error
	var data models.BusDexCexPriceSpreadData

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadDataService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusDexCexPriceSpreadData对象
func (e *BusDexCexPriceSpreadData) Get(d *dto.BusDexCexPriceSpreadDataGetReq, p *actions.DataPermission, model *models.BusDexCexPriceSpreadData) error {
	var data models.BusDexCexPriceSpreadData

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusDexCexPriceSpreadData error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusDexCexPriceSpreadData对象
func (e *BusDexCexPriceSpreadData) Insert(c *dto.BusDexCexPriceSpreadDataInsertReq) error {
	var err error
	var data models.BusDexCexPriceSpreadData
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadDataService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusDexCexPriceSpreadData对象
func (e *BusDexCexPriceSpreadData) Update(c *dto.BusDexCexPriceSpreadDataUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusDexCexPriceSpreadData{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusDexCexPriceSpreadDataService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusDexCexPriceSpreadData
func (e *BusDexCexPriceSpreadData) Remove(d *dto.BusDexCexPriceSpreadDataDeleteReq, p *actions.DataPermission) error {
	var data models.BusDexCexPriceSpreadData

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusDexCexPriceSpreadData error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *BusDexCexPriceSpreadData) GetLatestSpreadData() error {
	observerList := make([]models.BusDexCexTriangularObserver, 0)
	var observer models.BusDexCexTriangularObserver

	err := e.Orm.Model(&observer).
		Where("status = ?", 1). // 运行中的实例
		Find(&observerList).Error

	if err != nil {
		e.Log.Errorf("GetRunningObservers error:%s \r\n", err)
		return err
	}

	for _, observer := range observerList {
		observerId := observer.ObserverId
		state, err := client.GetObserverState(observerId)
		if err != nil {
			e.Log.Errorf("grpc获取最新价差数据失败， error:%s \r\n", err)
			continue
		}
		currentTime := time.Now()
		buyOnDex := state.GetBuyOnDex()
		cexSellPrice, dexBuyPrice := e.calculate_dex_cex_price(buyOnDex)
		sellOnDex := state.GetSellOnDex()
		cexBuyPrice, dexSellPrice := e.calculate_dex_cex_price(sellOnDex)

		spreadData := models.BusDexCexPriceSpreadData{
			ObserverId:           observerId,
			Symbol:               observer.Symbol,
			DexBuyPrice:          strconv.FormatFloat(dexBuyPrice, 'f', 6, 64),
			CexSellPrice:         strconv.FormatFloat(cexSellPrice, 'f', 6, 64),
			DexSellPrice:         strconv.FormatFloat(dexSellPrice, 'f', 6, 64),
			CexBuyPrice:          strconv.FormatFloat(cexBuyPrice, 'f', 6, 64),
			DexBuySpread:         strconv.FormatFloat(cexSellPrice-dexBuyPrice, 'f', 6, 64),
			DexBuySpreadPercent:  strconv.FormatFloat(math.Abs((cexSellPrice-dexBuyPrice)/dexBuyPrice), 'f', 6, 64),
			DexSellSpread:        strconv.FormatFloat(dexSellPrice-cexBuyPrice, 'f', 6, 64),
			DexSellSpreadPercent: strconv.FormatFloat(math.Abs((dexSellPrice-cexBuyPrice)/cexBuyPrice), 'f', 6, 64),
			SnapshotTime:         time.Now(),
		}
		e.Orm.Create(&spreadData)

		// 获取最新的dex买的价差统计信息
		dexBuyData := models.BusDexCexPriceSpreadStatistics{}
		err = e.Orm.Model(&dexBuyData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 1).Order("created_at desc").First(&dexBuyData).Limit(1).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果查询不到最新的价差统计信息，则需根据是否出现价差，判断要不要新建一条记录
				if cexSellPrice-dexBuyPrice > 0 {
					// dex买cex卖出现正向价差
					dexBuyData = models.BusDexCexPriceSpreadStatistics{
						ObserverId:         observerId,
						SpreadType:         "1", //dex买cex卖即为1
						Symbol:             observer.Symbol,
						StartTime:          &currentTime,
						Duration:           "0",
						MaxPriceDifference: spreadData.DexBuySpread,
						MinPriceDifference: spreadData.DexBuySpread,
					}
					e.Orm.Create(&dexBuyData)
				}
			} else {
				e.Log.Errorf("get observer spread statistics error:%s \r\n", err)
				continue
			}
		}

		e.Log.Infof("observer spread statistics info:%+v \r\n", dexBuyData)
		if cexSellPrice-dexBuyPrice > 0 {
			//如果有正向价差，需要更新下最大最小价差
			//startTime := dexBuyData.StartTime
			//dexBuyData.Duration = strconv.FormatFloat(currentTime.Sub(*startTime).Seconds(), 'f', 0, 64)
			maxPriceDifference, _ := strconv.ParseFloat(dexBuyData.MaxPriceDifference, 64)
			minPriceDifference, _ := strconv.ParseFloat(dexBuyData.MinPriceDifference, 64)
			dexBuySpreadf, _ := strconv.ParseFloat(spreadData.DexBuySpread, 64)
			e.Log.Infof("maxPriceSpread: %f, minPriceSpread: %f, dexBuySpreadf : %f", maxPriceDifference, minPriceDifference, dexBuySpreadf)
			if dexBuySpreadf >= maxPriceDifference {
				dexBuyData.MaxPriceDifference = strconv.FormatFloat(dexBuySpreadf, 'f', 0, 64)
			}
			if dexBuySpreadf <= minPriceDifference {
				dexBuyData.MinPriceDifference = strconv.FormatFloat(dexBuySpreadf, 'f', 0, 64)
			}
			err = e.Orm.Save(&dexBuyData).Error
			if err != nil {
				e.Log.Errorf("save observer spread statistics error:%s \r\n", err)
				continue
			}
		} else {
			if dexBuyData.Id != 0 {
				//如果价差变成负的了，则需要更新价差结束时间
				dexBuyData.EndTime = &currentTime
				//startTime := dexBuyData.StartTime
				//dexBuyData.Duration = strconv.FormatFloat(currentTime.Sub(*startTime).Seconds(), 'f', 0, 64)
				err = e.Orm.Save(&dexBuyData).Error
				if err != nil {
					e.Log.Errorf("save observer spread statistics error:%s \r\n", err)
					continue
				}
			}

		}

		// 获取最新的dex卖的价差统计信息
		dexSellData := models.BusDexCexPriceSpreadStatistics{}
		err = e.Orm.Model(&dexSellData).Where("observer_id = ? and spread_type = ? and end_time is null", observerId, 2).Order("created_at desc").First(&dexSellData).Limit(1).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果查询不到最新的价差统计信息，则需根据是否出现价差，判断要不要新建一条记录
				if dexSellPrice-cexBuyPrice > 0 {
					// dex卖cex买出现正向价差
					dexSellData = models.BusDexCexPriceSpreadStatistics{
						ObserverId:         observerId,
						SpreadType:         "2", //dex卖cex买即为2
						Symbol:             observer.Symbol,
						StartTime:          &currentTime,
						Duration:           "0",
						MaxPriceDifference: spreadData.DexSellSpread,
						MinPriceDifference: spreadData.DexSellSpread,
					}
					e.Orm.Create(&dexSellData)
				}
			} else {
				e.Log.Errorf("get observer spread statistics error:%s \r\n", err)
				continue
			}
		}

		e.Log.Infof("observer spread statistics info:%+v \r\n", dexSellData)
		if dexSellPrice-cexBuyPrice > 0 {
			e.Log.Infof("dex 卖出存在正向价差")
			//如果有正向价差，需要更新下最大最小价差
			//startTime := dexSellData.StartTime
			//dexSellData.Duration = strconv.FormatFloat(currentTime.Sub(*startTime).Seconds(), 'f', 0, 64)
			maxPriceDifference, _ := strconv.ParseFloat(dexSellData.MaxPriceDifference, 64)
			minPriceDifference, _ := strconv.ParseFloat(dexSellData.MinPriceDifference, 64)
			dexSellSpreadf, _ := strconv.ParseFloat(spreadData.DexSellSpread, 64)
			e.Log.Infof("maxPriceSpread: %f, minPriceSpread: %f, dexSellSpread : %f", maxPriceDifference, minPriceDifference, dexSellSpreadf)
			if dexSellSpreadf >= maxPriceDifference {
				dexSellData.MaxPriceDifference = strconv.FormatFloat(dexSellSpreadf, 'f', 0, 64)
			}
			if dexSellSpreadf <= minPriceDifference {
				dexSellData.MinPriceDifference = strconv.FormatFloat(dexSellSpreadf, 'f', 0, 64)
			}
			err = e.Orm.Save(&dexSellData).Error
			if err != nil {
				e.Log.Errorf("get observer spread statistics error:%s \r\n", err)
				continue
			}
		} else {
			if dexSellData.Id != 0 {
				//如果价差变成负的了，则需要更新价差结束时间
				dexSellData.EndTime = &currentTime
				//startTime := dexBuyData.StartTime
				//dexSellData.Duration = strconv.FormatFloat(currentTime.Sub(*startTime).Seconds(), 'f', 0, 64)
				err = e.Orm.Save(&dexSellData).Error
				if err != nil {
					e.Log.Errorf("get observer spread statistics error:%s \r\n", err)
					continue
				}
			}
		}

	}
	return nil
}

func (e *BusDexCexPriceSpreadData) calculate_dex_cex_price(priceState *pb.ArbitrageState) (float64, float64) {
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
