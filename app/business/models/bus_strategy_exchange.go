package models

import (
	"quanta-admin/common/models"
)

type BusStrategyExchange struct {
	models.Model

	ExchangeName string `json:"exchangeName" gorm:"type:varchar(255);comment:名称"`
	ExchangeType string `json:"exchangeType" gorm:"type:tinyint;comment:交易所类型"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyExchange) TableName() string {
	return "bus_strategy_exchange"
}

func (e *BusStrategyExchange) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyExchange) GetId() interface{} {
	return e.Id
}
