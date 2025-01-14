package models

import (
	"quanta-admin/common/models"
)

type BusStrategyInstanceConfig struct {
	models.Model

	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;comment:策略实例id"`
	SchemaText         string `json:"schemaText" gorm:"type:longtext;comment:参数schema"`
	SchemaType         string `json:"schemaType" gorm:"type:varchar(16);comment:schema类型"`
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
