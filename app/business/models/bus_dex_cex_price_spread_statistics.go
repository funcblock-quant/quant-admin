package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusDexCexPriceSpreadStatistics struct {
	models.Model

	ObserverId         string     `json:"observerId" gorm:"type:varchar(255);comment:观察器id"`
	SpreadType         string     `json:"spreadType" gorm:"type:tinyint;comment:价差类型"`
	Symbol             string     `json:"symbol" gorm:"type:varchar(255);comment:观察币种"`
	StartTime          *time.Time `json:"startTime" gorm:"type:timestamp;default:NULL;comment:正向价差开始时间"`
	EndTime            *time.Time `json:"endTime" gorm:"type:timestamp;default:NULL;comment:正向价差结束时间"`
	Duration           string     `json:"duration" gorm:"type:int;comment:价差持续时间"`
	MaxPriceDifference float64    `json:"maxPriceDifference" gorm:"type:float;comment:最大价差"`
	MinPriceDifference float64    `json:"minPriceDifference" gorm:"type:float;comment:最小价差"`
	CreatedAt          time.Time  `json:"createdAt" gorm:"comment:创建时间"`
}

func (BusDexCexPriceSpreadStatistics) TableName() string {
	return "bus_dex_cex_price_spread_statistics"
}

func (e *BusDexCexPriceSpreadStatistics) GetId() interface{} {
	return e.Id
}
