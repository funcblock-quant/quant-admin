package models

import (
	"quanta-admin/common/models"
)

type BusStrategyWatchList struct {
	models.Model

	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;comment:策略实例id"`
	Symbol             string `json:"symbol" gorm:"type:varchar(64);comment:观察币种名称"`
	IsDeleted          string `json:"isDeleted" gorm:"type:tinyint;comment:删除标识位"`
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
