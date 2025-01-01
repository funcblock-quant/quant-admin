package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusFuturePositionRecordGetPageReq struct {
	dto.Pagination `search:"-"`
	ArbitrageId    string `form:"arbitrageId"  search:"type:exact;column:arbitrage_id;table:bus_future_position_record" comment:"套利记录id"`
	Symbol         string `form:"symbol"  search:"type:exact;column:symbol;table:bus_future_position_record" comment:"交易币种"`
	OrderId        string `form:"orderId"  search:"type:exact;column:order_id;table:bus_future_position_record" comment:"交易所订单id"`
	OrderClientId  string `form:"orderClientId"  search:"type:exact;column:order_client_id;table:bus_future_position_record" comment:"策略端id"`
	PositionSide   string `form:"positionSide"  search:"type:exact;column:position_side;table:bus_future_position_record" comment:"持仓方向"`
	BusFuturePositionRecordOrder
}

type BusFuturePositionRecordOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:bus_future_position_record"`
	ArbitrageId   string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_future_position_record"`
	Symbol        string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_future_position_record"`
	Side          string `form:"sideOrder"  search:"type:order;column:side;table:bus_future_position_record"`
	Leverage      string `form:"leverageOrder"  search:"type:order;column:leverage;table:bus_future_position_record"`
	OrderId       string `form:"orderIdOrder"  search:"type:order;column:order_id;table:bus_future_position_record"`
	OrderClientId string `form:"orderClientIdOrder"  search:"type:order;column:order_client_id;table:bus_future_position_record"`
	OriginQty     string `form:"originQtyOrder"  search:"type:order;column:origin_qty;table:bus_future_position_record"`
	OriginPrice   string `form:"originPriceOrder"  search:"type:order;column:origin_price;table:bus_future_position_record"`
	OriginType    string `form:"originTypeOrder"  search:"type:order;column:origin_type;table:bus_future_position_record"`
	PositionSide  string `form:"positionSideOrder"  search:"type:order;column:position_side;table:bus_future_position_record"`
	TimeInForce   string `form:"timeInForceOrder"  search:"type:order;column:time_in_force;table:bus_future_position_record"`
	Role          string `form:"roleOrder"  search:"type:order;column:role;table:bus_future_position_record"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:bus_future_position_record"`
	Fees          string `form:"feesOrder"  search:"type:order;column:fees;table:bus_future_position_record"`
	FeeAsset      string `form:"feeAssetOrder"  search:"type:order;column:fee_asset;table:bus_future_position_record"`
	RealizedPnl   string `form:"realizedPnlOrder"  search:"type:order;column:realized_pnl;table:bus_future_position_record"`
	UnrealizedPnl string `form:"unrealizedPnlOrder"  search:"type:order;column:unrealized_pnl;table:bus_future_position_record"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_future_position_record"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_future_position_record"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_future_position_record"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_future_position_record"`
}

func (m *BusFuturePositionRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusFuturePositionRecordInsertReq struct {
	Id            int    `json:"-" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	Side          string `json:"side" comment:"买卖方向"`
	Leverage      string `json:"leverage" comment:"合约杠杆"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OriginPrice   string `json:"originPrice" comment:"原始委托价格"`
	OriginType    string `json:"originType" comment:"触发前订单类型"`
	PositionSide  string `json:"positionSide" comment:"持仓方向"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"持仓状态"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	RealizedPnl   string `json:"realizedPnl" comment:"已实现盈亏"`
	UnrealizedPnl string `json:"unrealizedPnl" comment:"未实现盈亏"`
	common.ControlBy
}

func (s *BusFuturePositionRecordInsertReq) Generate(model *models.BusFuturePositionRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.Symbol = s.Symbol
	model.Side = s.Side
	model.Leverage = s.Leverage
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.OriginQty = s.OriginQty
	model.OriginPrice = s.OriginPrice
	model.OriginType = s.OriginType
	model.PositionSide = s.PositionSide
	model.TimeInForce = s.TimeInForce
	model.Role = s.Role
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.RealizedPnl = s.RealizedPnl
	model.UnrealizedPnl = s.UnrealizedPnl
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusFuturePositionRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusFuturePositionRecordUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	Side          string `json:"side" comment:"买卖方向"`
	Leverage      string `json:"leverage" comment:"合约杠杆"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OriginPrice   string `json:"originPrice" comment:"原始委托价格"`
	OriginType    string `json:"originType" comment:"触发前订单类型"`
	PositionSide  string `json:"positionSide" comment:"持仓方向"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"持仓状态"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	RealizedPnl   string `json:"realizedPnl" comment:"已实现盈亏"`
	UnrealizedPnl string `json:"unrealizedPnl" comment:"未实现盈亏"`
	common.ControlBy
}

func (s *BusFuturePositionRecordUpdateReq) Generate(model *models.BusFuturePositionRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.Symbol = s.Symbol
	model.Side = s.Side
	model.Leverage = s.Leverage
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.OriginQty = s.OriginQty
	model.OriginPrice = s.OriginPrice
	model.OriginType = s.OriginType
	model.PositionSide = s.PositionSide
	model.TimeInForce = s.TimeInForce
	model.Role = s.Role
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.RealizedPnl = s.RealizedPnl
	model.UnrealizedPnl = s.UnrealizedPnl
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusFuturePositionRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusFuturePositionRecordGetReq 功能获取请求参数
type BusFuturePositionRecordGetReq struct {
	Id int `uri:"id"`
}

func (s *BusFuturePositionRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusFuturePositionRecordDeleteReq 功能删除请求参数
type BusFuturePositionRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusFuturePositionRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
