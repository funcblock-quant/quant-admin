package models

import (
	"quanta-admin/common/models"
)

type BusStrategyInstanceConfig struct {
	models.Model

	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;comment:策略实例id"`
	ParamKey           string `json:"paramKey" gorm:"type:varchar(64);comment:参数的唯一标识"`
	ParamValue         string `json:"paramValue" gorm:"type:text;comment:参数值"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyInstanceConfig) TableName() string {
	return "bus_strategy_instance_config"
}

func (e *BusStrategyInstanceConfig) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyInstanceConfig) GetId() interface{} {
	return e.Id
}
