package models

import (
	"quanta-admin/common/models"
)

type GatemtMakerTrades struct {
	models.Model

	TradeId                                    string `json:"tradeId" gorm:"type:bigint;comment:gatemt trade id"`
	ClientOrderId                              string `json:"clientOrderId" gorm:"type:varchar(50);comment:ClientOrderId"`
	LocalTradedTime                            string `json:"localTradedTime" gorm:"type:varchar(50);comment:gatemt trade time"`
	LocalTradedAt                              string `json:"localTradedAt" gorm:"type:bigint;comment:gatemt trade timestamp in milliseconds"`
	BinanceLimitOrderSide                      string `json:"binanceLimitOrderSide" gorm:"type:varchar(10);comment:BinanceLimitOrderSide"`
	BinanceLimitOrderPrice                     string `json:"binanceLimitOrderPrice" gorm:"type:decimal(28,10);comment:BinanceLimitOrderPrice"`
	BinanceLimitClientOrderId                  string `json:"binanceLimitClientOrderId" gorm:"type:varchar(50);comment:BinanceLimitClientOrderId"`
	BinanceLimitOrderAmount                    string `json:"binanceLimitOrderAmount" gorm:"type:decimal(28,10);comment:BinanceLimitOrderAmount"`
	BinanceLimitOrderId                        string `json:"binanceLimitOrderId" gorm:"type:bigint;comment:BinanceLimitOrderId"`
	BinanceLimitOrderErr                       string `json:"binanceLimitOrderErr" gorm:"type:text;comment:BinanceLimitOrderErr"`
	BinanceMarketClientOrderId                 string `json:"binanceMarketClientOrderId" gorm:"type:varchar(50);comment:BinanceMarketClientOrderId"`
	BinanceMarketOrderAmount                   string `json:"binanceMarketOrderAmount" gorm:"type:decimal(28,10);comment:BinanceMarketOrderAmount"`
	BinanceMarketOrderId                       string `json:"binanceMarketOrderId" gorm:"type:bigint;comment:BinanceMarketOrderId"`
	BinanceMarketOrderStatus                   string `json:"binanceMarketOrderStatus" gorm:"type:varchar(20);comment:BinanceMarketOrderStatus"`
	BinanceMarketOrderExecutedQuantity         string `json:"binanceMarketOrderExecutedQuantity" gorm:"type:decimal(28,10);comment:BinanceMarketOrderExecutedQuantity"`
	BinanceMarketOrderCummulativeQuoteQuantity string `json:"binanceMarketOrderCummulativeQuoteQuantity" gorm:"type:decimal(28,10);comment:BinanceMarketOrderCummulativeQuoteQuantity"`
	BinanceMarketOrderErr                      string `json:"binanceMarketOrderErr" gorm:"type:text;comment:BinanceMarketOrderErr"`
	Data                                       string `json:"data" gorm:"type:json;comment:Data"`
}

func (GatemtMakerTrades) TableName() string {
	return "gatemt_maker_trades"
}

func (e *GatemtMakerTrades) GetId() interface{} {
	return e.Id
}
