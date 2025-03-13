package dto

import (
	"strconv"
	"time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusPriceTriggerStrategyInstanceGetPageReq struct {
	dto.Pagination `search:"-"`
	CloseTime      time.Time `form:"closeTime"  search:"type:exact;column:close_time;table:bus_price_trigger_strategy_instance" comment:"停止时间"`
	Status         string    `form:"status"  search:"type:exact;column:status;table:bus_price_trigger_strategy_instance" comment:"状态，created, started, stopped, closed"`
	Symbol         string    `form:"symbol"  search:"type:exact;column:symbol;table:bus_price_trigger_strategy_instance" comment:"交易币对"`
	ExchangeUserId string    `form:"exchangeUserId"  search:"type:exact;column:exchange_user_id;table:bus_price_trigger_strategy_instance" comment:"交易所userId"`
	UserId         string    `form:"createBy"  search:"type:exact;column:create_by;table:bus_price_trigger_strategy_instance" comment:"创建人"`
	ApiConfig      int       `form:"apiConfig" search:"type:exact;column:api_config;table:bus_price_trigger_strategy_instance" comment:"api配置id"`
	BusPriceTriggerStrategyInstanceOrder
}

type BusPriceTriggerStrategyInstanceOrder struct {
	Id         string `form:"idOrder"  search:"type:order;column:id;table:bus_price_trigger_strategy_instance"`
	OpenPrice  string `form:"openPriceOrder"  search:"type:order;column:open_price;table:bus_price_trigger_strategy_instance"`
	ClosePrice string `form:"closePriceOrder"  search:"type:order;column:close_price;table:bus_price_trigger_strategy_instance"`
	Amount     string `form:"amountOrder"  search:"type:order;column:amount;table:bus_price_trigger_strategy_instance"`
	Side       string `form:"sideOrder"  search:"type:order;column:side;table:bus_price_trigger_strategy_instance"`
	Symbol     string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_price_trigger_strategy_instance"`
	CloseTime  string `form:"closeTimeOrder"  search:"type:order;column:close_time;table:bus_price_trigger_strategy_instance"`
	Status     string `form:"statusOrder"  search:"type:order;column:status;table:bus_price_trigger_strategy_instance"`
	CreateBy   string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_price_trigger_strategy_instance"`
	UpdateBy   string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_price_trigger_strategy_instance"`
	CreatedAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_price_trigger_strategy_instance"`
	UpdatedAt  string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_price_trigger_strategy_instance"`
	DeletedAt  string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_price_trigger_strategy_instance"`
}

func (m *BusPriceTriggerStrategyInstanceGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type TriggerStrategyInstanceGetUserListReq struct {
	UserId string `form:"userId"  search:"type:exact;column:create_by;table:bus_price_trigger_strategy_instance" comment:"创建人"`
}

type BusPriceTriggerStrategyResp struct {
	Id                string                                     `json:"id"`
	OpenPrice         string                                     `json:"openPrice"`
	ClosePrice        string                                     `json:"closePrice"`
	CloseOrderType    string                                     `json:"closeOrderType" comment:"平仓模式"`
	Amount            string                                     `json:"amount"`
	Side              string                                     `json:"side"`
	Symbol            string                                     `json:"symbol"`
	CloseTime         time.Time                                  `json:"closeTime"`
	Status            string                                     `json:"status"`
	ApiConfig         string                                     `json:"apiConfig"`
	ApiConfigData     models.BusPriceTriggerStrategyApikeyConfig `json:"apiConfigData" gorm:"-"`
	CreatedAt         time.Time                                  `json:"createdAt"`
	Details           []models.BusPriceMonitorForOptionHedging   `json:"details" gorm:"-"`
	Statistical       BusPriceTriggerStrategyStatistical         `json:"statistical" gorm:"-"`
	ExchangeUserId    string                                     `json:"exchangeUserId"`
	ExecuteNum        int                                        `json:"executeNum"`
	DelayTime         int                                        `json:"delayTime"`
	ProfitTargetType  string                                     `json:"profitTargetType"`
	ProfitTargetPrice string                                     `json:"profitTargetPrice"`
	LossTargetPrice   string                                     `json:"lossTargetPrice"`
	CallbackRatio     float64                                    `json:"callbackRatio"`
	MinProfit         float64                                    `json:"minProfit"`
	AverageSlippage   float64                                    `json:"averageSlippage"`
	LatestSlippage    float64                                    `json:"latestSlippage"`
}

// BusPriceTriggerStrategyStatistical 价格触发下单的统计数据
type BusPriceTriggerStrategyStatistical struct {
	OrderNum int    `json:"orderNum"`
	TotalPnl string `json:"totalPnl"`
}

type BusPriceTriggerStrategyInstanceInsertReq struct {
	Id                int       `json:"-" comment:""` //
	OpenPrice         string    `json:"openPrice" comment:"开仓价格"`
	ClosePrice        string    `json:"closePrice" comment:"平仓价格"`
	CloseOrderType    string    `json:"closeOrderType" comment:"平仓模式"`
	Amount            string    `json:"amount" comment:"开仓数量"`
	Side              string    `json:"side" comment:"买卖方向"`
	Symbol            string    `json:"symbol" comment:"交易币种"`
	CloseTime         time.Time `json:"closeTime" comment:"停止时间"`
	ApiConfig         int       `json:"apiConfig" comment:"api配置id"`
	Status            string    `json:"status" comment:"状态，created, started, stopped, closed"`
	ExchangeUserId    string    `json:"exchangeUserId"`
	ExecuteNum        int       `json:"executeNum"`
	DelayTime         int       `json:"delayTime"`
	ProfitTargetType  string    `json:"profitTargetType"`
	ProfitTargetPrice float64   `json:"profitTargetPrice"`
	LossTargetPrice   float64   `json:"lossTargetPrice"`
	CallbackRatio     float64   `json:"callbackRatio"`
	// CutoffRatio       float64   `json:"cutoffRatio"`
	MinProfit float64 `json:"minProfit"`
	common.ControlBy
}

type BusPriceTriggerStrategyInstanceUpdateProfitTargetReq struct {
	Id                int     `json:"id" comment:""` //
	ProfitTargetType  string  `json:"profitTargetType"`
	ProfitTargetPrice float64 `json:"profitTargetPrice"`
	LossTargetPrice   float64 `json:"lossTargetPrice"`
	CallbackRatio     float64 `json:"callbackRatio"`
	// CutoffRatio       float64 `json:"cutoffRatio"`
	MinProfit float64 `json:"minProfit"`
	common.ControlBy
}

type BusPriceTriggerStrategyInstanceUpdateExecuteNumReq struct {
	Id         int  `json:"id" comment:""` //
	ExecuteNum *int `json:"executeNum"`
	DelayTime  *int `json:"delayTime"`
	common.ControlBy
}

type StopTriggerInstanceRequest struct {
	InstanceId string `json:"id" comment:""`
}

func (s *BusPriceTriggerStrategyInstanceInsertReq) Generate(model *models.BusPriceTriggerStrategyInstance) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.OpenPrice = s.OpenPrice
	model.ClosePrice = s.ClosePrice
	model.Amount = s.Amount
	model.Side = s.Side
	model.Symbol = s.Symbol
	model.CloseTime = s.CloseTime
	model.Status = s.Status
	model.ApiConfig = s.ApiConfig
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.ExchangeUserId = s.ExchangeUserId
	model.ExecuteNum = s.ExecuteNum
	model.CloseOrderType = s.CloseOrderType
	model.DelayTime = s.DelayTime
	model.ProfitTargetType = s.ProfitTargetType
	if s.ProfitTargetType == "LIMIT" {
		model.ProfitTargetPrice = strconv.FormatFloat(s.ProfitTargetPrice, 'f', -1, 64)
		model.LossTargetPrice = strconv.FormatFloat(s.LossTargetPrice, 'f', -1, 64)
	} else if s.ProfitTargetType == "FLOATING" {
		model.CallbackRatio = &s.CallbackRatio
		cutoffRatio := 1.0 // 目前这块，策略端要求不接受前端传值，后端写死
		model.CutoffRatio = &cutoffRatio
		model.MinProfit = &s.MinProfit
	}
}

func (s *BusPriceTriggerStrategyInstanceInsertReq) GetId() interface{} {
	return s.Id
}

type BusPriceTriggerStrategyInstanceUpdateReq struct {
	Id         int       `uri:"id" comment:""` //
	OpenPrice  string    `json:"openPrice" comment:"开仓价格"`
	ClosePrice string    `json:"closePrice" comment:"平仓价格"`
	Amount     string    `json:"amount" comment:"开仓数量"`
	Side       string    `json:"side" comment:"买卖方向"`
	Symbol     string    `json:"symbol" comment:"交易币种"`
	CloseTime  time.Time `json:"closeTime" comment:"停止时间"`
	ApiConfig  int       `json:"apiConfig" comment:"api配置id"`
	Status     string    `json:"status" comment:"状态，created, started, stopped, closed"`
	common.ControlBy
}

func (s *BusPriceTriggerStrategyInstanceUpdateReq) Generate(model *models.BusPriceTriggerStrategyInstance) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.OpenPrice = s.OpenPrice
	model.ClosePrice = s.ClosePrice
	model.Amount = s.Amount
	model.Side = s.Side
	model.Symbol = s.Symbol
	model.CloseTime = s.CloseTime
	model.Status = s.Status
	model.ApiConfig = s.ApiConfig
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusPriceTriggerStrategyInstanceUpdateReq) GetId() interface{} {
	return s.Id
}

// BusPriceTriggerStrategyInstanceGetReq 功能获取请求参数
type BusPriceTriggerStrategyInstanceGetReq struct {
	Id int `uri:"id"`
}

func (s *BusPriceTriggerStrategyInstanceGetReq) GetId() interface{} {
	return s.Id
}

// BusPriceTriggerStrategyInstanceDeleteReq 功能删除请求参数
type BusPriceTriggerStrategyInstanceDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusPriceTriggerStrategyInstanceDeleteReq) GetId() interface{} {
	return s.Ids
}

type BusPriceTriggerStrategySymbolListResp struct {
	Symbol string `json:"symbol" gorm:"column:symbol"`
}

type BusPriceTriggerStrategyExchangeUserIdListResp struct {
	ExchangeUserId string `json:"exchangeUserId" gorm:"column:exchange_user_id"`
}
