package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusStrategyInstance struct {
	models.Model

	StrategyId   string     `json:"strategyId" gorm:"type:bigint;comment:策略id"`
	InstanceName string     `json:"instanceName" gorm:"type:varchar(255);comment:策略实例名称"`
	StartRunTime *time.Time `json:"startRunTime" gorm:"type:timestamp;default:NULL;comment:启动时间"`
	StopRunTime  *time.Time `json:"stopRunTime" gorm:"type:timestamp;default:NULL;comment:停止时间"`
	Type         string     `json:"type" gorm:"type:tinyint;default:NULL;comment:实例类型，0-观察者，1-交易者"`
	Status       string     `json:"status" gorm:"type:tinyint;default:0;comment:运行状态"`
	models.ModelTime
	models.ControlBy
}

func (BusStrategyInstance) TableName() string {
	return "bus_strategy_instance"
}

func (e *BusStrategyInstance) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusStrategyInstance) GetId() interface{} {
	return e.Id
}
