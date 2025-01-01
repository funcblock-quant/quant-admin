package models

import (
	"quanta-admin/common/models"
)

type BusStrategyConfigDict struct {
	models.Model

	StrategyId   string `json:"strategyId" gorm:"type:bigint;comment:id"`
	ParamKey     string `json:"paramKey" gorm:"type:varchar(64);comment:参数的唯一标识"`
	ParamName    string `json:"paramName" gorm:"type:varchar(128);comment:参数名称"`
	ParamType    string `json:"paramType" gorm:"type:tinyint;comment:参数类型"`
	DefaultValue string `json:"defaultValue" gorm:"type:text;comment:参数的默认值"`
	Required     string `json:"required" gorm:"type:tinyint;comment:是否为必填参数"`
	Description  string `json:"description" gorm:"type:text;comment:参数用途描述"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyConfigDict) TableName() string {
	return "bus_strategy_config_dict"
}

func (e *BusStrategyConfigDict) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyConfigDict) GetId() interface{} {
	return e.Id
}
