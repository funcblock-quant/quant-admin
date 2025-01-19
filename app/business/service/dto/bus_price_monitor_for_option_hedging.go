package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusPriceMonitorForOptionHedgingGetPageReq struct {
	dto.Pagination     `search:"-"`
	StrategyInstanceId string `form:"strategyInstanceId"  search:"type:exact;column:strategy_instance_id;table:bus_price_monitor_for_option_hedging" comment:"策略实例id"`
	MonitoredOpenedNum string `form:"monitoredOpenedNum"  search:"type:exact;column:monitored_opened_num;table:bus_price_monitor_for_option_hedging" comment:"监控的开单数量"`
	BusPriceMonitorForOptionHedgingOrder
}

type BusPriceMonitorForOptionHedgingOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_price_monitor_for_option_hedging"`
	ArbitrageId        string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_price_monitor_for_option_hedging"`
	ExchangeId         string `form:"exchangeIdOrder"  search:"type:order;column:exchange_id;table:bus_price_monitor_for_option_hedging"`
	ExchangeName       string `form:"exchangeNameOrder"  search:"type:order;column:exchange_name;table:bus_price_monitor_for_option_hedging"`
	ExchangeType       string `form:"exchangeTypeOrder"  search:"type:order;column:exchange_type;table:bus_price_monitor_for_option_hedging"`
	Side               string `form:"sideOrder"  search:"type:order;column:side;table:bus_price_monitor_for_option_hedging"`
	Symbol             string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_price_monitor_for_option_hedging"`
	OrderId            string `form:"orderIdOrder"  search:"type:order;column:order_id;table:bus_price_monitor_for_option_hedging"`
	OrderClientId      string `form:"orderClientIdOrder"  search:"type:order;column:order_client_id;table:bus_price_monitor_for_option_hedging"`
	OriginQty          string `form:"originQtyOrder"  search:"type:order;column:origin_qty;table:bus_price_monitor_for_option_hedging"`
	OriginPrice        string `form:"originPriceOrder"  search:"type:order;column:origin_price;table:bus_price_monitor_for_option_hedging"`
	OriginType         string `form:"originTypeOrder"  search:"type:order;column:origin_type;table:bus_price_monitor_for_option_hedging"`
	TimeInForce        string `form:"timeInForceOrder"  search:"type:order;column:time_in_force;table:bus_price_monitor_for_option_hedging"`
	Role               string `form:"roleOrder"  search:"type:order;column:role;table:bus_price_monitor_for_option_hedging"`
	Pnl                string `form:"pnlOrder"  search:"type:order;column:pnl;table:bus_price_monitor_for_option_hedging"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:bus_price_monitor_for_option_hedging"`
	Fees               string `form:"feesOrder"  search:"type:order;column:fees;table:bus_price_monitor_for_option_hedging"`
	FeeAsset           string `form:"feeAssetOrder"  search:"type:order;column:fee_asset;table:bus_price_monitor_for_option_hedging"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_price_monitor_for_option_hedging"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_price_monitor_for_option_hedging"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_price_monitor_for_option_hedging"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_price_monitor_for_option_hedging"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_price_monitor_for_option_hedging"`
	MonitoredOpenedNum string `form:"monitoredOpenedNumOrder"  search:"type:order;column:monitored_opened_num;table:bus_price_monitor_for_option_hedging"`
}

func (m *BusPriceMonitorForOptionHedgingGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusPriceMonitorForOptionHedgingInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	ArbitrageId        string `json:"arbitrageId" comment:"套利记录id"`
	ExchangeId         string `json:"exchangeId" comment:"交易所id"`
	ExchangeName       string `json:"exchangeName" comment:"交易所名称"`
	ExchangeType       string `json:"exchangeType" comment:"交易平台类型"`
	Side               string `json:"side" comment:"买卖方向"`
	Symbol             string `json:"symbol" comment:"交易币种"`
	OrderId            string `json:"orderId" comment:"交易所订单id"`
	OrderClientId      string `json:"orderClientId" comment:"策略端生成的id"`
	OriginQty          string `json:"originQty" comment:"原始委托数量"`
	OriginPrice        string `json:"originPrice" comment:"原始委托价格"`
	OriginType         string `json:"originType" comment:"触发前订单类型"`
	TimeInForce        string `json:"timeInForce" comment:"有效方法"`
	Role               string `json:"role" comment:"下单角色"`
	Pnl                string `json:"pnl" comment:"总盈亏"`
	Status             string `json:"status" comment:"持仓状态"`
	Fees               string `json:"fees" comment:"交易手续费"`
	FeeAsset           string `json:"feeAsset" comment:"交易手续费计价单位"`
	MonitoredOpenedNum string `json:"monitoredOpenedNum" comment:"监控的开单数量"`
	common.ControlBy
}

func (s *BusPriceMonitorForOptionHedgingInsertReq) Generate(model *models.BusPriceMonitorForOptionHedging) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.ExchangeId = s.ExchangeId
	model.ExchangeName = s.ExchangeName
	model.ExchangeType = s.ExchangeType
	model.Side = s.Side
	model.Symbol = s.Symbol
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.OriginQty = s.OriginQty
	model.OriginPrice = s.OriginPrice
	model.OriginType = s.OriginType
	model.TimeInForce = s.TimeInForce
	model.Role = s.Role
	model.Pnl = s.Pnl
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.MonitoredOpenedNum = s.MonitoredOpenedNum
}

func (s *BusPriceMonitorForOptionHedgingInsertReq) GetId() interface{} {
	return s.Id
}

type BusPriceMonitorForOptionHedgingUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	ArbitrageId        string `json:"arbitrageId" comment:"套利记录id"`
	ExchangeId         string `json:"exchangeId" comment:"交易所id"`
	ExchangeName       string `json:"exchangeName" comment:"交易所名称"`
	ExchangeType       string `json:"exchangeType" comment:"交易平台类型"`
	Side               string `json:"side" comment:"买卖方向"`
	Symbol             string `json:"symbol" comment:"交易币种"`
	OrderId            string `json:"orderId" comment:"交易所订单id"`
	OrderClientId      string `json:"orderClientId" comment:"策略端生成的id"`
	OriginQty          string `json:"originQty" comment:"原始委托数量"`
	OriginPrice        string `json:"originPrice" comment:"原始委托价格"`
	OriginType         string `json:"originType" comment:"触发前订单类型"`
	TimeInForce        string `json:"timeInForce" comment:"有效方法"`
	Role               string `json:"role" comment:"下单角色"`
	Pnl                string `json:"pnl" comment:"总盈亏"`
	Status             string `json:"status" comment:"持仓状态"`
	Fees               string `json:"fees" comment:"交易手续费"`
	FeeAsset           string `json:"feeAsset" comment:"交易手续费计价单位"`
	MonitoredOpenedNum string `json:"monitoredOpenedNum" comment:"监控的开单数量"`
	common.ControlBy
}

func (s *BusPriceMonitorForOptionHedgingUpdateReq) Generate(model *models.BusPriceMonitorForOptionHedging) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.ExchangeId = s.ExchangeId
	model.ExchangeName = s.ExchangeName
	model.ExchangeType = s.ExchangeType
	model.Side = s.Side
	model.Symbol = s.Symbol
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.OriginQty = s.OriginQty
	model.OriginPrice = s.OriginPrice
	model.OriginType = s.OriginType
	model.TimeInForce = s.TimeInForce
	model.Role = s.Role
	model.Pnl = s.Pnl
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.MonitoredOpenedNum = s.MonitoredOpenedNum
}

func (s *BusPriceMonitorForOptionHedgingUpdateReq) GetId() interface{} {
	return s.Id
}

// BusPriceMonitorForOptionHedgingGetReq 功能获取请求参数
type BusPriceMonitorForOptionHedgingGetReq struct {
	Id int `uri:"id"`
}

func (s *BusPriceMonitorForOptionHedgingGetReq) GetId() interface{} {
	return s.Id
}

// BusPriceMonitorForOptionHedgingDeleteReq 功能删除请求参数
type BusPriceMonitorForOptionHedgingDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusPriceMonitorForOptionHedgingDeleteReq) GetId() interface{} {
	return s.Ids
}
