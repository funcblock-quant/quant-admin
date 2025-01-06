package models

import (
	"quanta-admin/common/models"
)

type BusStrategySymbolGroup struct {
	models.Model

	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;default:NULL;comment:策略实例id"`
	GroupName          string `json:"groupName" gorm:"type:varchar(64);comment:观察交易对名称"`
	GroupType          string `json:"groupType" gorm:"type:tinyint;comment:组类型"`
	AutoRefresh        bool   `json:"autoRefresh" gorm:"type:bool;comment:是否自动刷新"`
	IsActive           bool   `json:"isActive" gorm:"type:bool;default:false;comment:是否激活"`
	RefreshInterval    string `json:"refreshInterval" gorm:"type:int;comment:自动刷新间隔"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategySymbolGroup) TableName() string {
	return "bus_strategy_symbol_group"
}

func (e *BusStrategySymbolGroup) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategySymbolGroup) GetId() interface{} {
	return e.Id
}
