package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type GatemtMakerOrdersGetPageReq struct {
	dto.Pagination  `search:"-"`
	ClientOrderId   string `form:"clientOrderId"  search:"type:exact;column:client_order_id;table:gatemt_maker_orders" comment:"Client order id"`
	Symbol          string `form:"symbol"  search:"type:exact;column:symbol;table:gatemt_maker_orders" comment:""`
	ExchangeOrderId string `form:"exchangeOrderId"  search:"type:exact;column:exchange_order_id;table:gatemt_maker_orders" comment:""`
	GatemtMakerOrdersOrder
}

type GatemtMakerOrdersOrder struct {
	Id                     string `form:"idOrder"  search:"type:order;column:id;table:gatemt_maker_orders"`
	ClientOrderId          string `form:"clientOrderIdOrder"  search:"type:order;column:client_order_id;table:gatemt_maker_orders"`
	Symbol                 string `form:"symbolOrder"  search:"type:order;column:symbol;table:gatemt_maker_orders"`
	OrderSide              string `form:"orderSideOrder"  search:"type:order;column:order_side;table:gatemt_maker_orders"`
	Price                  string `form:"priceOrder"  search:"type:order;column:price;table:gatemt_maker_orders"`
	Amount                 string `form:"amountOrder"  search:"type:order;column:amount;table:gatemt_maker_orders"`
	PriceLevelUpdatedTime  string `form:"priceLevelUpdatedTimeOrder"  search:"type:order;column:price_level_updated_time;table:gatemt_maker_orders"`
	PriceLevelUpdatedAt    string `form:"priceLevelUpdatedAtOrder"  search:"type:order;column:price_level_updated_at;table:gatemt_maker_orders"`
	LowerAskPrice          string `form:"lowerAskPriceOrder"  search:"type:order;column:lower_ask_price;table:gatemt_maker_orders"`
	HigherBidPrice         string `form:"higherBidPriceOrder"  search:"type:order;column:higher_bid_price;table:gatemt_maker_orders"`
	CreatedTime            string `form:"createdTimeOrder"  search:"type:order;column:created_time;table:gatemt_maker_orders"`
	CreatedAt              string `form:"createdAtOrder"  search:"type:order;column:created_at;table:gatemt_maker_orders"`
	ElapseLocalTime        string `form:"elapseLocalTimeOrder"  search:"type:order;column:elapse_local_time;table:gatemt_maker_orders"`
	ExchangeOrderId        string `form:"exchangeOrderIdOrder"  search:"type:order;column:exchange_order_id;table:gatemt_maker_orders"`
	ExchangeOrderStatus    string `form:"exchangeOrderStatusOrder"  search:"type:order;column:exchange_order_status;table:gatemt_maker_orders"`
	ExchangeOrderCreatedAt string `form:"exchangeOrderCreatedAtOrder"  search:"type:order;column:exchange_order_created_at;table:gatemt_maker_orders"`
}

func (m *GatemtMakerOrdersGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GatemtMakerOrdersInsertReq struct {
	Id                     int    `json:"-" comment:""` //
	ClientOrderId          string `json:"clientOrderId" comment:"Client order id"`
	Symbol                 string `json:"symbol" comment:""`
	OrderSide              string `json:"orderSide" comment:""`
	Price                  string `json:"price" comment:""`
	Amount                 string `json:"amount" comment:""`
	PriceLevelUpdatedTime  string `json:"priceLevelUpdatedTime" comment:""`
	PriceLevelUpdatedAt    string `json:"priceLevelUpdatedAt" comment:""`
	LowerAskPrice          string `json:"lowerAskPrice" comment:""`
	HigherBidPrice         string `json:"higherBidPrice" comment:""`
	CreatedTime            string `json:"createdTime" comment:""`
	ElapseLocalTime        string `json:"elapseLocalTime" comment:""`
	ExchangeOrderId        string `json:"exchangeOrderId" comment:""`
	ExchangeOrderStatus    string `json:"exchangeOrderStatus" comment:""`
	ExchangeOrderCreatedAt string `json:"exchangeOrderCreatedAt" comment:""`
	common.ControlBy
}

func (s *GatemtMakerOrdersInsertReq) Generate(model *models.GatemtMakerOrders) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ClientOrderId = s.ClientOrderId
	model.Symbol = s.Symbol
	model.OrderSide = s.OrderSide
	model.Price = s.Price
	model.Amount = s.Amount
	model.ExchangeOrderId = s.ExchangeOrderId
}

func (s *GatemtMakerOrdersInsertReq) GetId() interface{} {
	return s.Id
}

type GatemtMakerOrdersUpdateReq struct {
	Id                     int    `uri:"id" comment:""` //
	ClientOrderId          string `json:"clientOrderId" comment:"Client order id"`
	Symbol                 string `json:"symbol" comment:""`
	OrderSide              string `json:"orderSide" comment:""`
	Price                  string `json:"price" comment:""`
	Amount                 string `json:"amount" comment:""`
	PriceLevelUpdatedTime  string `json:"priceLevelUpdatedTime" comment:""`
	PriceLevelUpdatedAt    string `json:"priceLevelUpdatedAt" comment:""`
	LowerAskPrice          string `json:"lowerAskPrice" comment:""`
	HigherBidPrice         string `json:"higherBidPrice" comment:""`
	CreatedTime            string `json:"createdTime" comment:""`
	ElapseLocalTime        string `json:"elapseLocalTime" comment:""`
	ExchangeOrderId        string `json:"exchangeOrderId" comment:""`
	ExchangeOrderStatus    string `json:"exchangeOrderStatus" comment:""`
	ExchangeOrderCreatedAt string `json:"exchangeOrderCreatedAt" comment:""`
	common.ControlBy
}

func (s *GatemtMakerOrdersUpdateReq) Generate(model *models.GatemtMakerOrders) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ClientOrderId = s.ClientOrderId
	model.Symbol = s.Symbol
	model.OrderSide = s.OrderSide
	model.Price = s.Price
	model.Amount = s.Amount
	model.ExchangeOrderId = s.ExchangeOrderId
}

func (s *GatemtMakerOrdersUpdateReq) GetId() interface{} {
	return s.Id
}

// GatemtMakerOrdersGetReq 功能获取请求参数
type GatemtMakerOrdersGetReq struct {
	Id int `uri:"id"`
}

func (s *GatemtMakerOrdersGetReq) GetId() interface{} {
	return s.Id
}

// GatemtMakerOrdersDeleteReq 功能删除请求参数
type GatemtMakerOrdersDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GatemtMakerOrdersDeleteReq) GetId() interface{} {
	return s.Ids
}
