package models

import (
	"quanta-admin/common/models"
)

type BusStrategyTradeList struct {
	models.Model

	SymbolGroupId string `json:"symbolGroupId" gorm:"type:bigint;comment:币对组id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:交易币种名称"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyTradeList) TableName() string {
	return "bus_strategy_trade_list"
}

func (e *BusStrategyTradeList) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyTradeList) GetId() interface{} {
	return e.Id
}
