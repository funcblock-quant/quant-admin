package service

import (
	"errors"
	"fmt"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	"quanta-admin/app/grpc/client"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	"quanta-admin/common/actions"
	cDto "quanta-admin/common/dto"
	"strconv"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BusPriceTriggerStrategyInstance struct {
	service.Service
}

const (
	TRIGGER_INSTANCE_STATUS_CREATED = "created"
	TRIGGER_INSTANCE_STATUS_STARTED = "started"
	TRIGGER_INSTANCE_STATUS_PAUSE   = "paused" // 暂停，目前代表已止盈
	TRIGGER_INSTANCE_STATUS_EXPIRED = "expired"
	TRIGGER_INSTANCE_STATUS_STOP    = "stopped" //已停止
)

// GetPage 获取BusPriceTriggerStrategyInstance列表
func (e *BusPriceTriggerStrategyInstance) GetPage(c *dto.BusPriceTriggerStrategyInstanceGetPageReq, p *actions.DataPermission, list *[]dto.BusPriceTriggerStrategyResp, count *int64) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance
	var detail models.BusPriceMonitorForOptionHedging

	if !c.IsHistory && c.Status == TRIGGER_INSTANCE_STATUS_PAUSE {
		//如果是查询止盈状态，并且不是历史数据，则需要查询close time 晚于当前时间的数据
		err = e.Orm.Model(&data).
			Where("close_time > ?", time.Now()).
			Scopes(
				cDto.MakeCondition(c.GetNeedSearch()),
				cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
				actions.Permission(data.TableName(), p),
			).
			Find(list).Limit(-1).Offset(-1).
			Count(count).Error

	} else {
		err = e.Orm.Model(&data).
			Scopes(
				cDto.MakeCondition(c.GetNeedSearch()),
				cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
				actions.Permission(data.TableName(), p),
			).
			Find(list).Limit(-1).Offset(-1).
			Count(count).Error
	}

	for index, strategy := range *list {
		statistical := dto.BusPriceTriggerStrategyStatistical{}
		details := make([]models.BusPriceMonitorForOptionHedging, 0)
		err := e.Orm.Model(&detail).Where("strategy_instance_id = ?", strategy.Id).Order("id desc").Find(&details).Error
		if err != nil {
			e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get details error:%s \r\n", err)
			return err
		}

		for _, trade := range details {
			//封装最新的滑点
			if trade.Slippage != "" {
				slippage, err := strconv.ParseFloat(trade.Slippage, 64)
				if err != nil {
					e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get details error:%s \r\n", err)
					return err
				}
				(*list)[index].LatestSlippage = slippage
				break
			}
		}

		var apiConfig models.BusPriceTriggerStrategyApikeyConfig
		err = e.Orm.Unscoped().Model(&apiConfig).Where("id = ?", strategy.ApiConfig).First(&apiConfig).Error
		if err != nil {
			e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get apiConfig error:%s \r\n", err)
			return err
		}

		(*list)[index].ApiConfigData = apiConfig

		totalOrderNum := len(details)
		totalPnl := decimal.NewFromFloat(0)
		totalFee := decimal.NewFromFloat(0)
		for _, d := range details {
			var pnl decimal.Decimal
			if d.Pnl == "" {
				pnl = decimal.NewFromFloat(0)
			} else {
				pnl, err = decimal.NewFromString(d.Pnl)
				if err != nil {
					e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get details error:%s \r\n", err)
				}
			}
			var fees decimal.Decimal
			if d.Fees == "" {
				fees = decimal.NewFromFloat(0)
			} else {
				fees, err = decimal.NewFromString(d.Fees)
				if err != nil {
					e.Log.Errorf("BusPriceTriggerStrategyInstanceService Get details error:%s \r\n", err)
				}
			}
			totalFee = totalFee.Add(fees)
			totalPnl = totalPnl.Add(pnl)

		}
		statistical.OrderNum = totalOrderNum
		statistical.TotalPnl = totalPnl.StringFixed(8)
		statistical.NetProfit = totalPnl.Sub(totalFee).StringFixed(8)

		(*list)[index].Statistical = statistical
		(*list)[index].Details = details
	}
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Get(d *dto.BusPriceTriggerStrategyInstanceGetReq, p *actions.DataPermission, model *models.BusPriceTriggerStrategyInstance) error {
	var data models.BusPriceTriggerStrategyInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBusPriceTriggerStrategyInstance error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Insert(c *dto.BusPriceTriggerStrategyInstanceInsertReq) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance
	c.Generate(&data)
	data.Status = TRIGGER_INSTANCE_STATUS_CREATED
	e.Log.Infof("create price trigger instance:%+v", data)
	//启动事务
	tx := e.Orm.Begin()
	if tx.Error != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService error:%s \r\n", tx.Error)
		return tx.Error
	}

	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService Insert error:%s \r\n", err)
		return err
	}
	var apiKeyConfig models.BusPriceTriggerStrategyApikeyConfig
	err = tx.First(&apiKeyConfig, c.ApiConfig).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("cannot found apikey config:%s \r\n", err)
	}

	//创建成功后， 自动通过grpc启动
	//apiConfigReq := trigger_service.APIConfig{
	//	ApiKey:    apiKeyConfig.ApiKey,
	//	SecretKey: apiKeyConfig.SecretKey,
	//	Exchange:  apiKeyConfig.Exchange,
	//}
	//
	//request := &trigger_service.StartTriggerRequest{
	//	InstanceId: strconv.Itoa(data.Id),
	//	OpenPrice:  c.OpenPrice,
	//	ClosePrice: c.ClosePrice,
	//	Side:       c.Side,
	//	Amount:     c.Amount,
	//	Symbol:     c.Symbol,
	//	StopTime:   strconv.FormatInt(c.CloseTime.UnixMilli(), 10),
	//	ApiConfig:  &apiConfigReq,
	//	UserId:     c.ExchangeUserId,
	//}
	//_, err = client.StartTriggerInstance(request)
	//if err != nil {
	//	tx.Rollback()
	//	e.Log.Errorf("Service grpc start error:%s \r\n", err)
	//	return err
	//}
	e.Log.Infof("instance id : %d grpc start success\r\n", data.Id)

	tx.Commit()

	err = e.Orm.Model(&data).Update("status", TRIGGER_INSTANCE_STATUS_STARTED).Error
	if err != nil {
		e.Log.Errorf("start trigger update status failed:%s \r\n", err)
		return err
	}

	return nil
}

// StopInstance
// @Summary 暂停实例
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/stopInstance [post]
// @Security Bearer
func (e *BusPriceTriggerStrategyInstance) StopInstance(req *dto.StopTriggerInstanceRequest) error {
	var err error
	data := models.BusPriceTriggerStrategyInstance{}

	err = e.Orm.Model(&data).First(&data, req.InstanceId).Error
	if err != nil {
		e.Log.Errorf("Service StopInstance error:%s \r\n", err)
		return err
	}

	e.Log.Infof("stop instance id : %d\r\n", data.Id)
	request := trigger_service.StopTriggerRequest{
		InstanceId: strconv.Itoa(data.Id),
	}
	err = client.StopTriggerInstance(&request)
	if err != nil {
		e.Log.Errorf("Service StopInstance throw grpc error:%s \r\n", err)
		return err
	}
	err = e.Orm.Model(&data).
		Update("status", TRIGGER_INSTANCE_STATUS_STOP).
		Error

	if err != nil {
		e.Log.Errorf("Service StopInstance throw db error:%s \r\n", err)
		return err
	}
	return nil
}

// RestartInstance
// @Summary 重启实例
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/restartInstance [post]
// @Security Bearer
func (e *BusPriceTriggerStrategyInstance) RestartInstance(req *dto.RestartTriggerInstanceRequest) error {
	var err error
	data := models.BusPriceTriggerStrategyInstance{}

	err = e.Orm.Model(&data).First(&data, req.InstanceId).Error
	if err != nil {
		e.Log.Errorf("Service StopInstance error:%s \r\n", err)
		return err
	}

	// 校验参数
	// 1. 检查停止时间是否晚于当前时间
	if data.CloseTime.Before(time.Now()) {
		e.Log.Errorf("Service RestartInstance error: close time must be later than current time\r\n")
		return errors.New("任务已过期，无法重启")
	}

	e.Log.Infof("stop instance id : %d\r\n", data.Id)
	// 更新状态为 started，由定时任务去触发
	err = e.Orm.Model(&data).
		Update("status", TRIGGER_INSTANCE_STATUS_STARTED).
		Error

	if err != nil {
		e.Log.Errorf("Service RestartInstance throw db error:%s \r\n", err)
		return err
	}
	return nil
}

// UpdateProfitTarget
// @Accept  application/json
// @Product application/json
// @Param
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/updateProfitTarget [put]
// @Security Bearer
func (e *BusPriceTriggerStrategyInstance) UpdateProfitTarget(req *dto.BusPriceTriggerStrategyInstanceUpdateProfitTargetReq) error {
	var err error
	data := models.BusPriceTriggerStrategyInstance{}

	err = e.Orm.Model(&data).First(&data, req.Id).Error
	if err != nil {
		e.Log.Errorf("Service UpdateProfitTarget error:%s \r\n", err)
		return err
	}

	profitTargetType := req.ProfitTargetType
	profitTargetConfig := &trigger_service.ProfitTargetConfig{
		InstanceId: strconv.Itoa(data.Id),
	}
	if profitTargetType == "LIMIT" {
		//限价止盈
		profitTargetConfig.ProfitTargetType = trigger_service.ProfitTargetType_LIMIT
		profitTargetConfig.Config = &trigger_service.ProfitTargetConfig_LimitConfig{
			LimitConfig: &trigger_service.LimitTypeConfig{
				ProfitTargetPrice: req.ProfitTargetPrice,
				//LossTargetPrice:   req.LossTargetPrice,
			},
		}
	} else if profitTargetType == "FLOATING" {
		//浮动止盈
		profitTargetConfig.ProfitTargetType = trigger_service.ProfitTargetType_FLOATING
		cutOffRatio := float64(1)
		profitTargetConfig.Config = &trigger_service.ProfitTargetConfig_FloatingConfig{
			FloatingConfig: &trigger_service.FloatingTypeConfig{
				CallbackRatio: req.CallbackRatio,
				CutoffRatio:   cutOffRatio,
				MinProfit:     req.MinProfit,
			},
		}
	}

	e.Log.Infof("update profit target instance id : %d\r\n", data.Id)
	if config.ApplicationConfig.Mode != "dev" {
		err = client.UpdateProfitTarget(profitTargetConfig)
		if err != nil {
			e.Log.Errorf("Service StopInstance throw grpc error:%s \r\n", err)
			return err
		}
	}

	err = e.Orm.Model(&data).
		Updates(map[string]interface{}{
			"profit_target_type":  req.ProfitTargetType,
			"profit_target_price": fmt.Sprintf("%.16f", req.ProfitTargetPrice),
			"loss_target_price":   fmt.Sprintf("%.16f", req.LossTargetPrice),
			"callback_ratio":      req.CallbackRatio,
			"cutoff_ratio":        1,
			"min_profit":          req.MinProfit,
		}).Error

	if err != nil {
		e.Log.Errorf("Service StopInstance throw db error:%s \r\n", err)
		return err
	}
	return nil
}

// UpdateExecuteConfig
// @Accept  application/json
// @Product application/json
// @Param
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/updateExecuteNum [put]
// @Security Bearer
func (e *BusPriceTriggerStrategyInstance) UpdateExecuteConfig(req *dto.BusPriceTriggerStrategyInstanceUpdateExecuteNumReq) error {
	var err error
	data := models.BusPriceTriggerStrategyInstance{}

	err = e.Orm.Model(&data).First(&data, req.Id).Error
	if err != nil {
		e.Log.Errorf("Service UpdateProfitTarget error:%s \r\n", err)
		return err
	}

	executeConfig := &trigger_service.ExecuteConfig{
		InstanceId: strconv.Itoa(data.Id),
	}
	if req.ExecuteNum != nil {
		executeConfig.ExecuteNum = uint32(*req.ExecuteNum)
	}
	if req.DelayTime != nil {
		executeConfig.DelayTime = uint32(*req.DelayTime)
	}

	e.Log.Infof("update execute num, instance id : %d\r\n", data.Id)
	if config.ApplicationConfig.Mode != "dev" {
		err = client.UpdateExecuteConfig(executeConfig)
		if err != nil {
			e.Log.Errorf("Service UpdateExecuteNum throw grpc error:%s \r\n", err)
			return err
		}
	}

	updateData := map[string]interface{}{}

	if req.ExecuteNum != nil {
		updateData["execute_num"] = *req.ExecuteNum
	}
	if req.DelayTime != nil {
		updateData["delay_time"] = *req.DelayTime
	}

	err = e.Orm.Model(&data).
		Updates(updateData).Error

	if err != nil {
		e.Log.Errorf("Service StopInstance throw db error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改BusPriceTriggerStrategyInstance对象
func (e *BusPriceTriggerStrategyInstance) Update(c *dto.BusPriceTriggerStrategyInstanceUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.BusPriceTriggerStrategyInstance{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstanceService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除BusPriceTriggerStrategyInstance
func (e *BusPriceTriggerStrategyInstance) Remove(d *dto.BusPriceTriggerStrategyInstanceDeleteReq, p *actions.DataPermission) error {
	var data models.BusPriceTriggerStrategyInstance

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBusPriceTriggerStrategyInstance error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// GetSymbolList 获取BusPriceTriggerStrategyInstance所有币种列表
func (e *BusPriceTriggerStrategyInstance) GetSymbolList(list *[]dto.BusPriceTriggerStrategySymbolListResp) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance

	err = e.Orm.Model(&data).
		Select("symbol").
		Group("symbol").
		Debug().Find(list).Error

	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstance GetSymbolList error:%s \r\n", err)
		return err
	}
	return nil
}

// GetExchangeUserIdList 获取BusPriceTriggerStrategyInstance所有交易所userId
func (e *BusPriceTriggerStrategyInstance) GetExchangeUserIdList(userId int, list *[]dto.BusPriceTriggerStrategyExchangeUserIdListResp) error {
	var err error
	var data models.BusPriceTriggerStrategyInstance

	err = e.Orm.Model(&data).
		Select("exchange_user_id").
		Group("exchange_user_id").
		Where("create_by", userId).
		Debug().Find(list).Error

	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstance GetExchangeUserIdList error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *BusPriceTriggerStrategyInstance) MonitorExecuteNum() error {
	// 获取所有运行中的实例，并统计他们的执行次数
	var data []models.BusPriceTriggerStrategyInstance
	var count int64

	err := e.Orm.Model(&models.BusPriceTriggerStrategyInstance{}).
		Where("status = ?", TRIGGER_INSTANCE_STATUS_STARTED).
		Find(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorExecuteNum error:%s \r\n", err)
		return err
	}
	for _, instance := range data {
		err = e.Orm.Model(&models.BusPriceMonitorForOptionHedging{}).
			Where("strategy_instance_id =? and extra is NULL", instance.Id).
			Count(&count).Error
		if err != nil {
			e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorExecuteNum error:%s \r\n", err)
			continue
		}
		if count/2 >= int64(instance.ExecuteNum) {

			// 如果执行次数达到上限，暂停实例
			request := &trigger_service.StopTriggerRequest{
				InstanceId: strconv.Itoa(instance.Id),
			}
			err = client.StopTriggerInstance(request)
			if err != nil {
				e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorExecuteNum stop instance error:%s \r\n", err)
				continue
			}

			err = e.Orm.Model(&models.BusPriceTriggerStrategyInstance{}).
				Where("id = ?", instance.Id).
				Update("status", TRIGGER_INSTANCE_STATUS_STOP).Error
			if err != nil {
				e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorExecuteNum error:%s \r\n", err)
				continue
			}
		}
	}
	return nil
}

func (e *BusPriceTriggerStrategyInstance) MonitorStopProfitStatus() error {
	// 获取所有运行中的实例，并检查他们的止盈情况
	var data []models.BusPriceTriggerStrategyInstance

	err := e.Orm.Model(&models.BusPriceTriggerStrategyInstance{}).
		Where("status = ?", TRIGGER_INSTANCE_STATUS_STARTED).
		Find(&data).Error
	if err != nil {
		e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorExecuteNum error:%s \r\n", err)
		return err
	}
	for _, instance := range data {
		// 查询2s内是否有止盈单
		// 如果有，暂停实例
		// 获取2秒前的时间点
		var results []models.BusPriceMonitorForOptionHedging

		oneSecondAgo := time.Now().Add(-2 * time.Second)

		err := e.Orm.Model(&models.BusPriceMonitorForOptionHedging{}).
			Where("strategy_instance_id = ? AND extra IS NULL AND created_at > ? AND pnl > 0", instance.Id, oneSecondAgo).
			Find(&results).Error
		if err != nil {
			e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorStopProfitStatus error:%s \r\n", err)
			continue
		}

		if len(results) > 0 {
			// 如果有止盈单，修改实例状态
			err = e.Orm.Model(&models.BusPriceTriggerStrategyInstance{}).
				Where("id = ?", instance.Id).
				Update("status", TRIGGER_INSTANCE_STATUS_PAUSE).Error
			if err != nil {
				e.Log.Errorf("BusPriceTriggerStrategyInstance MonitorStopProfitStatus error:%s \r\n", err)
				continue
			}
		}
	}
	return nil
}

type SlippageResult struct {
	TradeId  string
	Slippage *float64
}

type AverageSlippageResult struct {
	InstanceId      string
	AverageSlippage *float64
}

func (e *BusPriceTriggerStrategyInstance) CalculateSlippageForPriceTriggerInstance() error {

	// 定时任务要分两部分，
	// 1. 计算历史数据，需要将已过期、已暂停的数据--这一部分不用处理了

	// var data []models.BusPriceTriggerStrategyInstance

	// err := e.Orm.Model(&models.BusPriceTriggerStrategyInstance{}).
	// 	Where("id > ?", 691). //目前只计算昨天到现在的数据
	// 	Order("created_at asc").
	// 	Find(&data).Error
	// if err != nil {
	// 	e.Log.Errorf("[Calculate Slippage] CalculateSlippageForPriceTriggerInstance error:%s \r\n", err)
	// 	return err
	// }

	// var allSlippage float64
	// var count int
	// for _, instance := range data {
	// 	// 获取该instance的全部成交。然后按照时间从前往后计算滑点，并记录一个平均滑点
	// 	var trades []models.BusPriceMonitorForOptionHedging
	// 	err = e.Orm.Model(&models.BusPriceMonitorForOptionHedging{}).
	// 		Where("strategy_instance_id =?", instance.Id).
	// 		Order("created_at asc").
	// 		Find(&trades).Error
	// 	if err != nil {
	// 		e.Log.Errorf("[Calculate Slippage] CalculateSlippageForPriceTriggerInstance error:%s \r\n", err)
	// 		continue
	// 	}

	// 	if len(trades) == 0 {
	// 		continue
	// 	}

	// 	openPrice, err1 := strconv.ParseFloat(instance.OpenPrice, 64)
	// 	closePrice, err2 := strconv.ParseFloat(instance.ClosePrice, 64)
	// 	if err1 != nil || err2 != nil {
	// 		e.Log.Errorf("[Calculate Slippage] parse openprice and closeprice error:%s or %s \r\n", err1, err2)
	// 		continue
	// 	}
	// 	side := instance.Side

	// 	var totalSlippage float64

	// 	e.Log.Infof("[Calculate Slippage] Calculate slippage for instance: %d \n", instance.Id)
	// 	for _, trade := range trades {
	// 		var slippage float64
	// 		tradePrice, err := strconv.ParseFloat(trade.OriginPrice, 64)
	// 		if err != nil {
	// 			e.Log.Errorf("[Calculate Slippage] parse trade price error:%s \r\n", err)
	// 			continue
	// 		}

	// 		originQty, err := decimal.NewFromString(trade.OriginQty)
	// 		if err != nil {
	// 			e.Log.Errorf("[Calculate Slippage] parse OriginQty error:%s \r\n", err)
	// 			continue
	// 		}

	// 		if !originQty.IsZero() {
	// 			//开仓
	// 			if side == "short" {
	// 				//做空
	// 				slippage = (tradePrice - openPrice) / openPrice * 100
	// 				e.Log.Infof("[Calculate Slippage] tradeId: %d, 做空-开仓 滑点：%f%% \r\n", trade.Id, slippage)
	// 			} else {
	// 				// 做多
	// 				slippage = (openPrice - tradePrice) / openPrice * 100
	// 				e.Log.Infof("[Calculate Slippage] tradeId: %d, 做多-开仓 滑点：%f%% \r\n", trade.Id, slippage)
	// 			}

	// 		} else {
	// 			//平仓
	// 			if side == "short" {
	// 				//做空
	// 				slippage = (closePrice - tradePrice) / closePrice * 100
	// 				e.Log.Infof("[Calculate Slippage] tradeId: %d, 做空-平仓 滑点：%f%% \r\n", trade.Id, slippage)
	// 			} else {
	// 				//做多
	// 				slippage = (tradePrice - closePrice) / closePrice * 100
	// 				e.Log.Infof("[Calculate Slippage] tradeId: %d, 做多-平仓 滑点：%f%% \r\n", trade.Id, slippage)
	// 			}
	// 		}
	// 		totalSlippage += slippage
	// 		allSlippage += slippage
	// 		count++
	// 	}

	// 	averageSlippage := totalSlippage / float64(len(trades))
	// 	e.Log.Infof("[Calculate Slippage]instanceId: %s,  averageSlipp: %f%%", instance.Id, averageSlippage)
	// }

	// allAverageSlippage := allSlippage / float64(count)
	// e.Log.Infof("[Calculate Slippage] all total averageSlipp: %f%%", allAverageSlippage)

	// 2. 计算最新成交数据，如果没有计算滑点的要计算滑点，并更新整个任务最新的平均滑点。
	// 查询最新的成交并且未计算滑点的记录
	var slippageList []SlippageResult

	db := e.Orm

	err := db.Raw(`
		SELECT
			t.id AS trade_id, 
			CASE
				WHEN i.side = 'short' THEN
					CASE
						WHEN t.origin_qty != 0 THEN
							CAST((t.origin_price - CAST(i.open_price AS DECIMAL)) / CAST(i.open_price AS DECIMAL) AS DECIMAL(18, 6))
						ELSE
							CAST((CAST(i.close_price AS DECIMAL) - t.origin_price) / CAST(i.close_price AS DECIMAL) AS DECIMAL(18, 6))
					END
				ELSE
					CASE
						WHEN t.origin_qty != 0 THEN
							CAST((CAST(i.open_price AS DECIMAL) - t.origin_price) / CAST(i.open_price AS DECIMAL) AS DECIMAL(18, 6))
						ELSE
							CAST((t.origin_price - CAST(i.close_price AS DECIMAL)) / CAST(i.close_price AS DECIMAL) AS DECIMAL(18, 6))
					END
			END AS slippage
		FROM
			bus_price_monitor_for_option_hedging t
		JOIN
			bus_price_trigger_strategy_instance i ON t.strategy_instance_id = i.id
		WHERE t.slippage is NULL`).Scan(&slippageList).Error

	if err != nil {
		e.Log.Errorf("[Calculate Slippage] CalculateSlippageForPriceTriggerInstance error:%s \r\n", err)
		return err
	}

	// 逐条保存交易记录的滑点值
	for _, slippage := range slippageList {
		if slippage.Slippage != nil {
			err := db.Model(&models.BusPriceMonitorForOptionHedging{}).
				Where("id =?", slippage.TradeId).
				Update("slippage", *slippage.Slippage).Error
			if err != nil {
				e.Log.Errorf("[Calculate Slippage] CalculateSlippageForPriceTriggerInstance error:%s \r\n", err)
				continue
			}
		}

	}

	// 保存完后，计算这些实例的平均交易滑点
	var instances []models.BusPriceTriggerStrategyInstance
	err = db.Model(&models.BusPriceTriggerStrategyInstance{}).
		Where("id > ?", 691).
		Where("status = ?", TRIGGER_INSTANCE_STATUS_STARTED).
		Order("created_at asc").
		Find(&instances).Error
	if err != nil {
		e.Log.Errorf("[Calculate Slippage] get instances error:%s \r\n", err)
		return err
	}

	for _, instance := range instances {
		var averageSlippageResult AverageSlippageResult
		err = db.Raw(`
			SELECT
				t.strategy_instance_id AS instance_id,
				ROUND(AVG(t.slippage), 6) AS average_slippage
			FROM
				bus_price_monitor_for_option_hedging t
			WHERE
				t.slippage IS NOT NULL and t.strategy_instance_id = ?`, instance.Id).Scan(&averageSlippageResult).Error
		if err != nil {
			e.Log.Errorf("[Calculate Slippage] calculate average slippage error:%s \r\n", err)
			continue
		}

		e.Log.Debug("[Calculate Slippage] averageSlippageResult: %+v  ", averageSlippageResult)

		// 更新实例的平均滑点
		if averageSlippageResult.AverageSlippage != nil {
			err = db.Model(&models.BusPriceTriggerStrategyInstance{}).
				Where("id =?", instance.Id).
				Update("average_slippage", *averageSlippageResult.AverageSlippage).Error
			if err != nil {
				e.Log.Errorf("[Calculate Slippage] update average slippage error:%s \r\n", err)
				continue
			}
			e.Log.Debug("[Calculate Slippage] instanceId: %s,  averageSlipp: %f%%", instance.Id, *averageSlippageResult.AverageSlippage)
		}

	}

	return nil
}
