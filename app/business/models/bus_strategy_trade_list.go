package models

import (
	"quanta-admin/common/models"
)

type BusStrategyTradeList struct {
	models.Model

	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;comment:策略实例id"`
	Symbol             string `json:"symbol" gorm:"type:varchar(64);comment:交易币种名称"`
	IsDeleted          string `json:"isDeleted" gorm:"type:tinyint;comment:删除标识位"`
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
