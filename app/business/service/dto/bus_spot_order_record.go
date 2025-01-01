package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusSpotOrderRecordGetPageReq struct {
	dto.Pagination `search:"-"`
	ArbitrageId    string `form:"arbitrageId"  search:"type:exact;column:arbitrage_id;table:bus_spot_order_record" comment:"套利记录id"`
	Side           string `form:"side"  search:"type:exact;column:side;table:bus_spot_order_record" comment:"买卖方向"`
	OrderId        string `form:"orderId"  search:"type:exact;column:order_id;table:bus_spot_order_record" comment:"交易所订单id"`
	OrderClientId  string `form:"orderClientId"  search:"type:exact;column:order_client_id;table:bus_spot_order_record" comment:"策略端id"`
	BusSpotOrderRecordOrder
}

type BusSpotOrderRecordOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:bus_spot_order_record"`
	ArbitrageId   string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_spot_order_record"`
	Symbol        string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_spot_order_record"`
	Side          string `form:"sideOrder"  search:"type:order;column:side;table:bus_spot_order_record"`
	OriginQty     string `form:"originQtyOrder"  search:"type:order;column:origin_qty;table:bus_spot_order_record"`
	OrderId       string `form:"orderIdOrder"  search:"type:order;column:order_id;table:bus_spot_order_record"`
	OrderClientId string `form:"orderClientIdOrder"  search:"type:order;column:order_client_id;table:bus_spot_order_record"`
	TimeInForce   string `form:"timeInForceOrder"  search:"type:order;column:time_in_force;table:bus_spot_order_record"`
	OrderType     string `form:"orderTypeOrder"  search:"type:order;column:order_type;table:bus_spot_order_record"`
	Fees          string `form:"feesOrder"  search:"type:order;column:fees;table:bus_spot_order_record"`
	FeeAsset      string `form:"feeAssetOrder"  search:"type:order;column:fee_asset;table:bus_spot_order_record"`
	Role          string `form:"roleOrder"  search:"type:order;column:role;table:bus_spot_order_record"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:bus_spot_order_record"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_spot_order_record"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_spot_order_record"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_spot_order_record"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_spot_order_record"`
}

func (m *BusSpotOrderRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusSpotOrderRecordInsertReq struct {
	Id            int    `json:"-" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	Side          string `json:"side" comment:"买卖方向"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	OrderType     string `json:"orderType" comment:"订单类型"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"成交状态"`
	common.ControlBy
}

func (s *BusSpotOrderRecordInsertReq) Generate(model *models.BusSpotOrderRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.Symbol = s.Symbol
	model.Side = s.Side
	model.OriginQty = s.OriginQty
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.TimeInForce = s.TimeInForce
	model.OrderType = s.OrderType
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.Role = s.Role
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusSpotOrderRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusSpotOrderRecordUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	ArbitrageId   string `json:"arbitrageId" comment:"套利记录id"`
	Symbol        string `json:"symbol" comment:"交易币种"`
	Side          string `json:"side" comment:"买卖方向"`
	OriginQty     string `json:"originQty" comment:"原始委托数量"`
	OrderId       string `json:"orderId" comment:"交易所订单id"`
	OrderClientId string `json:"orderClientId" comment:"策略端id"`
	TimeInForce   string `json:"timeInForce" comment:"有效方法"`
	OrderType     string `json:"orderType" comment:"订单类型"`
	Fees          string `json:"fees" comment:"交易手续费"`
	FeeAsset      string `json:"feeAsset" comment:"交易手续费计价单位"`
	Role          string `json:"role" comment:"交易角色"`
	Status        string `json:"status" comment:"成交状态"`
	common.ControlBy
}

func (s *BusSpotOrderRecordUpdateReq) Generate(model *models.BusSpotOrderRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.Symbol = s.Symbol
	model.Side = s.Side
	model.OriginQty = s.OriginQty
	model.OrderId = s.OrderId
	model.OrderClientId = s.OrderClientId
	model.TimeInForce = s.TimeInForce
	model.OrderType = s.OrderType
	model.Fees = s.Fees
	model.FeeAsset = s.FeeAsset
	model.Role = s.Role
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusSpotOrderRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusSpotOrderRecordGetReq 功能获取请求参数
type BusSpotOrderRecordGetReq struct {
	Id int `uri:"id"`
}

func (s *BusSpotOrderRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusSpotOrderRecordDeleteReq 功能删除请求参数
type BusSpotOrderRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusSpotOrderRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
