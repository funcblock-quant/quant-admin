package models

import (
	"quanta-admin/common/models"
)

type BusStrategyBaseInfo struct {
	models.Model

	StrategyName     string `json:"strategyName" gorm:"type:varchar(255);comment:策略名称"`
	StrategyCategory string `json:"strategyCategory" gorm:"type:tinyint;comment:策略交易类型"`
	Preference       string `json:"preference" gorm:"type:tinyint;comment:策略偏好"`
	Description      string `json:"description" gorm:"type:text;comment:策略描述"`
	Status           string `json:"status" gorm:"type:tinyint;default:1;comment:策略注册状态"`
	Owner            string `json:"owner" gorm:"type:varchar(255);comment:策略负责人"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyBaseInfo) TableName() string {
	return "bus_strategy_base_info"
}

func (e *BusStrategyBaseInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyBaseInfo) GetId() interface{} {
	return e.Id
}
