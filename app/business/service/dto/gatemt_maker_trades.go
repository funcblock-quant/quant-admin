package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type GatemtMakerTradesGetPageReq struct {
	dto.Pagination       `search:"-"`
	ClientOrderId        string `form:"clientOrderId"  search:"type:exact;column:client_order_id;table:gatemt_maker_trades" comment:""`
	BinanceLimitOrderId  string `form:"binanceLimitOrderId"  search:"type:exact;column:binance_limit_order_id;table:gatemt_maker_trades" comment:""`
	BinanceMarketOrderId string `form:"binanceMarketOrderId"  search:"type:exact;column:binance_market_order_id;table:gatemt_maker_trades" comment:""`
	GatemtMakerTradesOrder
}

type GatemtMakerTradesOrder struct {
	Id                                         string `form:"idOrder"  search:"type:order;column:id;table:gatemt_maker_trades"`
	TradeId                                    string `form:"tradeIdOrder"  search:"type:order;column:trade_id;table:gatemt_maker_trades"`
	ClientOrderId                              string `form:"clientOrderIdOrder"  search:"type:order;column:client_order_id;table:gatemt_maker_trades"`
	LocalTradedTime                            string `form:"localTradedTimeOrder"  search:"type:order;column:local_traded_time;table:gatemt_maker_trades"`
	LocalTradedAt                              string `form:"localTradedAtOrder"  search:"type:order;column:local_traded_at;table:gatemt_maker_trades"`
	BinanceLimitOrderSide                      string `form:"binanceLimitOrderSideOrder"  search:"type:order;column:binance_limit_order_side;table:gatemt_maker_trades"`
	BinanceLimitOrderPrice                     string `form:"binanceLimitOrderPriceOrder"  search:"type:order;column:binance_limit_order_price;table:gatemt_maker_trades"`
	BinanceLimitClientOrderId                  string `form:"binanceLimitClientOrderIdOrder"  search:"type:order;column:binance_limit_client_order_id;table:gatemt_maker_trades"`
	BinanceLimitOrderAmount                    string `form:"binanceLimitOrderAmountOrder"  search:"type:order;column:binance_limit_order_amount;table:gatemt_maker_trades"`
	BinanceLimitOrderId                        string `form:"binanceLimitOrderIdOrder"  search:"type:order;column:binance_limit_order_id;table:gatemt_maker_trades"`
	BinanceMarketClientOrderId                 string `form:"binanceMarketClientOrderIdOrder"  search:"type:order;column:binance_market_client_order_id;table:gatemt_maker_trades"`
	BinanceMarketOrderAmount                   string `form:"binanceMarketOrderAmountOrder"  search:"type:order;column:binance_market_order_amount;table:gatemt_maker_trades"`
	BinanceMarketOrderId                       string `form:"binanceMarketOrderIdOrder"  search:"type:order;column:binance_market_order_id;table:gatemt_maker_trades"`
	BinanceMarketOrderStatus                   string `form:"binanceMarketOrderStatusOrder"  search:"type:order;column:binance_market_order_status;table:gatemt_maker_trades"`
	BinanceMarketOrderExecutedQuantity         string `form:"binanceMarketOrderExecutedQuantityOrder"  search:"type:order;column:binance_market_order_executed_quantity;table:gatemt_maker_trades"`
	BinanceMarketOrderCummulativeQuoteQuantity string `form:"binanceMarketOrderCummulativeQuoteQuantityOrder"  search:"type:order;column:binance_market_order_cummulative_quote_quantity;table:gatemt_maker_trades"`
	Data                                       string `form:"dataOrder"  search:"type:order;column:data;table:gatemt_maker_trades"`
}

func (m *GatemtMakerTradesGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GatemtMakerTradesGetListReq struct {
	ClientOrderId string `uri:"clientOrderId"`
}

type GatemtMakerTradesInsertReq struct {
	Id                                         int    `json:"-" comment:""` //
	TradeId                                    string `json:"tradeId" comment:"gatemt trade id"`
	ClientOrderId                              string `json:"clientOrderId" comment:""`
	LocalTradedTime                            string `json:"localTradedTime" comment:"gatemt trade time"`
	LocalTradedAt                              string `json:"localTradedAt" comment:"gatemt trade timestamp in milliseconds"`
	BinanceLimitOrderSide                      string `json:"binanceLimitOrderSide" comment:""`
	BinanceLimitOrderPrice                     string `json:"binanceLimitOrderPrice" comment:""`
	BinanceLimitClientOrderId                  string `json:"binanceLimitClientOrderId" comment:""`
	BinanceLimitOrderAmount                    string `json:"binanceLimitOrderAmount" comment:""`
	BinanceLimitOrderId                        string `json:"binanceLimitOrderId" comment:""`
	BinanceMarketClientOrderId                 string `json:"binanceMarketClientOrderId" comment:""`
	BinanceMarketOrderAmount                   string `json:"binanceMarketOrderAmount" comment:""`
	BinanceMarketOrderId                       string `json:"binanceMarketOrderId" comment:""`
	BinanceMarketOrderStatus                   string `json:"binanceMarketOrderStatus" comment:""`
	BinanceMarketOrderExecutedQuantity         string `json:"binanceMarketOrderExecutedQuantity" comment:""`
	BinanceMarketOrderCummulativeQuoteQuantity string `json:"binanceMarketOrderCummulativeQuoteQuantity" comment:""`
	Data                                       string `json:"data" comment:""`
	common.ControlBy
}

func (s *GatemtMakerTradesInsertReq) Generate(model *models.GatemtMakerTrades) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.TradeId = s.TradeId
	model.ClientOrderId = s.ClientOrderId
	model.LocalTradedTime = s.LocalTradedTime
	model.LocalTradedAt = s.LocalTradedAt
	model.BinanceLimitOrderSide = s.BinanceLimitOrderSide
	model.BinanceLimitOrderPrice = s.BinanceLimitOrderPrice
	model.BinanceLimitClientOrderId = s.BinanceLimitClientOrderId
	model.BinanceLimitOrderAmount = s.BinanceLimitOrderAmount
	model.BinanceLimitOrderId = s.BinanceLimitOrderId
	model.BinanceMarketClientOrderId = s.BinanceMarketClientOrderId
	model.BinanceMarketOrderAmount = s.BinanceMarketOrderAmount
	model.BinanceMarketOrderId = s.BinanceMarketOrderId
	model.BinanceMarketOrderStatus = s.BinanceMarketOrderStatus
	model.BinanceMarketOrderExecutedQuantity = s.BinanceMarketOrderExecutedQuantity
	model.BinanceMarketOrderCummulativeQuoteQuantity = s.BinanceMarketOrderCummulativeQuoteQuantity
	model.Data = s.Data
}

func (s *GatemtMakerTradesInsertReq) GetId() interface{} {
	return s.Id
}

type GatemtMakerTradesUpdateReq struct {
	Id                                         int    `uri:"id" comment:""` //
	TradeId                                    string `json:"tradeId" comment:"gatemt trade id"`
	ClientOrderId                              string `json:"clientOrderId" comment:""`
	LocalTradedTime                            string `json:"localTradedTime" comment:"gatemt trade time"`
	LocalTradedAt                              string `json:"localTradedAt" comment:"gatemt trade timestamp in milliseconds"`
	BinanceLimitOrderSide                      string `json:"binanceLimitOrderSide" comment:""`
	BinanceLimitOrderPrice                     string `json:"binanceLimitOrderPrice" comment:""`
	BinanceLimitClientOrderId                  string `json:"binanceLimitClientOrderId" comment:""`
	BinanceLimitOrderAmount                    string `json:"binanceLimitOrderAmount" comment:""`
	BinanceLimitOrderId                        string `json:"binanceLimitOrderId" comment:""`
	BinanceMarketClientOrderId                 string `json:"binanceMarketClientOrderId" comment:""`
	BinanceMarketOrderAmount                   string `json:"binanceMarketOrderAmount" comment:""`
	BinanceMarketOrderId                       string `json:"binanceMarketOrderId" comment:""`
	BinanceMarketOrderStatus                   string `json:"binanceMarketOrderStatus" comment:""`
	BinanceMarketOrderExecutedQuantity         string `json:"binanceMarketOrderExecutedQuantity" comment:""`
	BinanceMarketOrderCummulativeQuoteQuantity string `json:"binanceMarketOrderCummulativeQuoteQuantity" comment:""`
	Data                                       string `json:"data" comment:""`
	common.ControlBy
}

func (s *GatemtMakerTradesUpdateReq) Generate(model *models.GatemtMakerTrades) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.TradeId = s.TradeId
	model.ClientOrderId = s.ClientOrderId
	model.LocalTradedTime = s.LocalTradedTime
	model.LocalTradedAt = s.LocalTradedAt
	model.BinanceLimitOrderSide = s.BinanceLimitOrderSide
	model.BinanceLimitOrderPrice = s.BinanceLimitOrderPrice
	model.BinanceLimitClientOrderId = s.BinanceLimitClientOrderId
	model.BinanceLimitOrderAmount = s.BinanceLimitOrderAmount
	model.BinanceLimitOrderId = s.BinanceLimitOrderId
	model.BinanceMarketClientOrderId = s.BinanceMarketClientOrderId
	model.BinanceMarketOrderAmount = s.BinanceMarketOrderAmount
	model.BinanceMarketOrderId = s.BinanceMarketOrderId
	model.BinanceMarketOrderStatus = s.BinanceMarketOrderStatus
	model.BinanceMarketOrderExecutedQuantity = s.BinanceMarketOrderExecutedQuantity
	model.BinanceMarketOrderCummulativeQuoteQuantity = s.BinanceMarketOrderCummulativeQuoteQuantity
	model.Data = s.Data
}

func (s *GatemtMakerTradesUpdateReq) GetId() interface{} {
	return s.Id
}

// GatemtMakerTradesGetReq 功能获取请求参数
type GatemtMakerTradesGetReq struct {
	Id int `uri:"id"`
}

func (s *GatemtMakerTradesGetReq) GetId() interface{} {
	return s.Id
}

// GatemtMakerTradesDeleteReq 功能删除请求参数
type GatemtMakerTradesDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GatemtMakerTradesDeleteReq) GetId() interface{} {
	return s.Ids
}
