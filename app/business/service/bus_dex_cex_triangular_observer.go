package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
	"strconv"

	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	pb "quanta-admin/app/grpc/proto/stub"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
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

	for _, resp := range *list {
		observerId := resp.ObserverId
		state, err := client.GetObserverState(observerId)
		if err != nil {
			e.Log.Errorf("grpc实时获取观察状态失败， error:%s \r\n", err)
			continue
		}
		resp.BaseProfit = strconv.FormatFloat(*state.BaseProfit, 'f', -1, 64)
		resp.QuoteProfit = strconv.FormatFloat(*state.QuoteProfit, 'f', -1, 64)
	}

	return nil
}

// Get 获取BusDexCexTriangularObserver对象
func (e *BusDexCexTriangularObserver) Get(d *dto.BusDexCexTriangularObserverGetReq, p *actions.DataPermission, model *models.BusDexCexTriangularObserver) error {
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
	symbols := c.Symbols
	if len(symbols) == 0 {
		return errors.New("empty symbols")
	}

	for _, symbol := range symbols {
		//循环创建监听

		var ammConfig = pb.AmmDexConfig{}
		var amberConfig = pb.AmberConfig{}
		var arbitrageConfig = pb.ArbitrageConfig{}
		c.GenerateAmmConfig(&ammConfig)
		c.GenerateAmberConfig(&amberConfig, symbol)
		c.GenerateArbitrageConfig(&arbitrageConfig)

		observerId, err := client.StartNewObserver(&amberConfig, &ammConfig, &arbitrageConfig)
		if err != nil {
			e.Log.Errorf("Service BatchInsert error:%s \r\n", err)
			continue
		}
		c.Generate(&data, symbol, observerId)
		err = e.Orm.Create(&data).Error
		if err != nil {
			e.Log.Errorf("BusDexCexTriangularObserverService Insert error:%s \r\n", err)
			return err
		}

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
