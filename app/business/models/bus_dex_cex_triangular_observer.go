package models

import (
	"errors"
	"quanta-admin/common/models"
)

type BusDexCexTriangularObserver struct {
	models.Model

	StrategyInstanceId string   `gorm:"not null;comment:策略id" json:"strategyInstanceId"`
	InstanceId         string   `gorm:"not null;comment:观察器id" json:"instanceId"`
	Symbol             string   `gorm:"not null;comment:观察币种" json:"symbol"`
	TargetToken        string   `gorm:"not null" json:"targetToken"`
	QuoteToken         string   `gorm:"not null" json:"quoteToken"`
	SymbolConnector    string   `gorm:"not null" json:"symbolConnector"`
	ExchangeType       string   `gorm:"not null;comment:交易所类型" json:"exchangeType"`
	DexType            string   `json:"dexType"`
	MaxArraySize       int      `gorm:"null;default:5" json:"maxArraySize"`
	Volume             *float64 `gorm:"comment:交易金额" json:"volume"` // 使用指针类型允许值为null
	TakerFee           *float64 `gorm:"not null;comment:交易所taker 费率" json:"takerFee"`
	AmmPoolId          *string  `gorm:"comment:ammPoolId" json:"ammPoolId"`
	TokenMint          *string  `gorm:"comment:base token合约" json:"tokenMint"` // 使用指针类型允许值为null
	SlippageBps        string   `gorm:"not null;comment:滑点bps" json:"slippage"`
	Depth              string   `gorm:"not null;default:20;comment:深度" json:"depth"`
	IsTrading          bool     `gorm:"default:false;comment:是否启动交易" json:"isTrading"`
	MinProfit          *float64 `gorm:"null;comment:最小利润" json:"minProfit"`
	PriorityFee        uint64   `gorm:"null;comment:交易优先费" json:"priorityFee"`
	JitoFee            uint64   `gorm:"null;comment:jito交易费" json:"jitoFee"`
	Status             string   `gorm:"not null;comment:状态，0-新增，1-启动，2-停止" json:"status"` // 使用 uint8 更合适
	models.ModelTime
	models.ControlBy
}

func (BusDexCexTriangularObserver) TableName() string {
	return "bus_dex_cex_triangular_instance"
}

func (e *BusDexCexTriangularObserver) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexTriangularObserver) GetId() interface{} {
	return e.Id
}

func (e *BusDexCexTriangularObserver) GetExchangeTypeForTrader() (string, error) {
	if e.ExchangeType == "Binance" {
		return "BinanceClassicUnifiedMargin", nil
	} else {
		return "", errors.New("not support exchange type")
	}
}
