package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusTriangleOrdersRecordGetPageReq struct {
	dto.Pagination `search:"-"`
	BusTriangleOrdersRecordOrder
}

type BusTriangleOrdersRecordOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:bus_triangle_orders_record"`
	ArbitrageId   string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_triangle_orders_record"`
	ExchangeId    string `form:"exchangeIdOrder"  search:"type:order;column:exchange_id;table:bus_triangle_orders_record"`
	ExchangeName  string `form:"exchangeNameOrder"  search:"type:order;column:exchange_name;table:bus_triangle_orders_record"`
	ExchangeType  string `form:"exchangeTypeOrder"  search:"type:order;column:exchange_type;table:bus_triangle_orders_record"`
	Side          string `form:"sideOrder"  search:"type:order;column:side;table:bus_triangle_orders_record"`
	Symbol        string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_triangle_orders_record"`
	OrderId       string `form:"orderIdOrder"  search:"type:order;column:order_id;table:bus_triangle_orders_record"`
	OrderClientId string `form:"orderClientIdOrder"  search:"type:order;column:order_client_id;table:bus_triangle_orders_record"`
	OriginQty     string `form:"originQtyOrder"  search:"type:order;column:origin_qty;table:bus_triangle_orders_record"`
	OriginPrice   string `form:"originPriceOrder"  search:"type:order;column:origin_price;table:bus_triangle_orders_record"`
	OriginType    string `form:"originTypeOrder"  search:"type:order;column:origin_type;table:bus_triangle_orders_record"`
	TimeInForce   string `form:"timeInForceOrder"  search:"type:order;column:time_in_force;table:bus_triangle_orders_record"`
	Role          string `form:"roleOrder"  search:"type:order;column:role;table:bus_triangle_orders_record"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:bus_triangle_orders_record"`
	Fees          string `form:"feesOrder"  search:"type:order;column:fees;table:bus_triangle_orders_record"`
	FeeAsset      string `form:"feeAssetOrder"  search:"type:order;column:fee_asset;table:bus_triangle_orders_record"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_triangle_orders_record"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_triangle_orders_record"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_triangle_orders_record"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_triangle_orders_record"`
}

func (m *BusTriangleOrdersRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusTriangleOrdersRecordInsertReq struct {
	Id            int    `json:"-" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	ExchangeId    string `json:"exchangeId" comment:"交易所id"`
	ExchangeName  string `json:"exchangeName" comment:"交易所名称"`
	ExchangeType  string `json:"exchangeType" comment:"交易平台类型"`
	Side          string `json:"side" comment:"买卖方向"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OriginPrice   string `json:"originPrice" comment:"原始委托价格"`
	OriginType    string `json:"originType" comment:"触发前订单类型"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"持仓状态"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	common.ControlBy
}

func (s *BusTriangleOrdersRecordInsertReq) Generate(model *models.BusTriangleOrdersRecord) {
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
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusTriangleOrdersRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusTriangleOrdersRecordUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	ExchangeId    string `json:"exchangeId" comment:"交易所id"`
	ExchangeName  string `json:"exchangeName" comment:"交易所名称"`
	ExchangeType  string `json:"exchangeType" comment:"交易平台类型"`
	Side          string `json:"side" comment:"买卖方向"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OriginPrice   string `json:"originPrice" comment:"原始委托价格"`
	OriginType    string `json:"originType" comment:"触发前订单类型"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"持仓状态"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	common.ControlBy
}

func (s *BusTriangleOrdersRecordUpdateReq) Generate(model *models.BusTriangleOrdersRecord) {
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
	model.Status = s.Status
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusTriangleOrdersRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusTriangleOrdersRecordGetReq 功能获取请求参数
type BusTriangleOrdersRecordGetReq struct {
	Id int `uri:"id"`
}

func (s *BusTriangleOrdersRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusTriangleOrdersRecordDeleteReq 功能删除请求参数
type BusTriangleOrdersRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusTriangleOrdersRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
