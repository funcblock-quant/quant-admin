package models

import (
	"quanta-admin/common/models"
)

type BusDexCexDailyTradeStatisticSnapshot struct {
	models.Model

	InstanceID   string  `gorm:"type:string;not null;comment:'实例ID'" json:"instanceId"`
	SnapshotDate string  `gorm:"type:string;not null;comment:'统计日期'" json:"snapshotDate"`
	Symbol       string  `gorm:"type:varchar(50);not null;comment:'交易币对，例如 BTC/USDT'" json:"symbol"`
	TargetToken  string  `gorm:"type:varchar(32);not null;comment:'基础币种，例如 BTC'" json:"targetToken"`
	QuoteToken   string  `gorm:"type:varchar(32);not null;comment:'计价币种，例如 USDT'" json:"quoteToken"`
	TotalVolume  float64 `gorm:"type:double;not null;default:0;comment:'当日总成交量 (计价币种)'" json:"totalVolume"`
	TotalProfit  float64 `gorm:"type:double;not null;default:0;comment:'当日总收益 (计价币种)'" json:"totalProfit"`
	TotalFee     float64 `gorm:"type:double;not null;default:0;comment:'当日总手续费 (计价币种)'" json:"totalFee"`
	models.ModelTime
	models.ControlBy
}

func (BusDexCexDailyTradeStatisticSnapshot) TableName() string {
	return "bus_dex_cex_daily_trade_statistic_snapshot"
}

func (e *BusDexCexDailyTradeStatisticSnapshot) GetId() interface{} {
	return e.Id
}
