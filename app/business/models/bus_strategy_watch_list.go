package models

import (
	"quanta-admin/common/models"
)

type BusStrategyWatchList struct {
	models.Model

	SymbolGroupId string `json:"symbolGroupId" gorm:"type:bigint;comment:币对组id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:观察币种名称"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyWatchList) TableName() string {
	return "bus_strategy_watch_list"
}

func (e *BusStrategyWatchList) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyWatchList) GetId() interface{} {
	return e.Id
}
