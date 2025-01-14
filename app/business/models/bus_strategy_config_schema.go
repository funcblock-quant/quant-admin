package models

import (
	"quanta-admin/common/models"
)

type BusStrategyConfigSchema struct {
	models.Model

	StrategyId string `json:"strategyId" gorm:"type:bigint;comment:关联策略表的ID"`
	SchemaText string `json:"schemaText" gorm:"type:longtext;comment:参数schema"`
	SchemaType string `json:"schemaType" gorm:"type:varchar(16);comment:schema类型"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyConfigSchema) TableName() string {
	return "bus_strategy_config_schema"
}

func (e *BusStrategyConfigSchema) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyConfigSchema) GetId() interface{} {
	return e.Id
}
