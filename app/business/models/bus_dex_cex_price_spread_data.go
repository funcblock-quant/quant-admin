package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusDexCexPriceSpreadData struct {
	models.Model

	ObserverId           string    `json:"observerId" gorm:"type:varchar(255);comment:观察器id"`
	Symbol               string    `json:"symbol" gorm:"type:varchar(255);comment:观察币种"`
	DexBuyPrice          float64   `json:"dexBuyPrice" gorm:"type:float;comment:dex买入价格"`
	DexSellPrice         float64   `json:"dexSellPrice" gorm:"type:float;comment:dex卖出价格"`
	CexBuyPrice          float64   `json:"cexBuyPrice" gorm:"type:float;comment:cex买入价格"`
	CexSellPrice         float64   `json:"cexSellPrice" gorm:"type:float;comment:cex卖出价格"`
	DexBuySpread         float64   `json:"dexBuySpread" gorm:"type:float;comment:Dex买入价差"`
	DexSellSpread        float64   `json:"dexSellSpread" gorm:"type:float;comment:Dex卖出价差"`
	DexBuySpreadPercent  string    `json:"dexBuySpreadPercent" gorm:"type:float;comment:Dex卖出价差百分比"`
	DexBuyProfit         float64   `json:"dexBuyProfit" gorm:"type:float;comment:Dex买入利润"`
	DexSellSpreadPercent string    `json:"dexSellSpreadPercent" gorm:"type:float;comment:Dex卖出价差百分比"`
	DexSellProfit        float64   `json:"dexSellProfit" gorm:"type:float;comment:Dex卖出利润"`
	SnapshotTime         time.Time `json:"snapshotTime" gorm:"type:timestamp;comment:快照时间"`
	CreatedAt            time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (BusDexCexPriceSpreadData) TableName() string {
	return "bus_dex_cex_price_spread_data"
}

func (e *BusDexCexPriceSpreadData) GetId() interface{} {
	return e.Id
}
