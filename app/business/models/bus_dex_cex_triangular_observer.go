package models

import (
	"quanta-admin/common/models"
)

type BusDexCexTriangularObserver struct {
	models.Model

	StrategyInstanceId string   `gorm:"not null;comment:策略id" json:"strategyInstanceId"`
	ObserverId         string   `gorm:"not null;comment:观察器id" json:"observerId"`
	Symbol             string   `gorm:"not null;comment:观察币种" json:"symbol"`
	ExchangeType       string   `gorm:"not null;comment:交易所类型" json:"exchangeType"`
	Volume             *float64 `gorm:"comment:交易金额" json:"volume"` // 使用指针类型允许值为null
	TakerFee           *float64 `gorm:"not null;comment:交易所taker 费率" json:"takerFee"`
	AmmPoolId          *string  `gorm:"comment:ammPoolId" json:"ammPoolId"`
	BaseTokenMint      *string  `gorm:"comment:base token合约" json:"baseTokenMint"`   // 使用指针类型允许值为null
	QuoteTokenMint     *string  `gorm:"comment:quote token合约" json:"quoteTokenMint"` // 使用指针类型允许值为null
	SlippageBps        string   `gorm:"not null;comment:滑点bps" json:"slippageBps"`
	Depth              string   `gorm:"not null;default:20;comment:深度" json:"depth"`
	Status             string   `gorm:"not null;comment:状态，0-新增，1-启动，2-停止" json:"status"` // 使用 uint8 更合适
	models.ModelTime
	models.ControlBy
}

func (BusDexCexTriangularObserver) TableName() string {
	return "bus_dex_cex_triangular_observer"
}

func (e *BusDexCexTriangularObserver) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexTriangularObserver) GetId() interface{} {
	return e.Id
}
