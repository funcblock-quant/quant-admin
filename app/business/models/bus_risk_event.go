package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusRiskEvent struct {
	models.Model

	StrategyId         string     `json:"strategyId" gorm:"type:varchar(100);comment:策略id"`
	StrategyInstanceId string     `json:"strategyInstanceId" gorm:"type:bigint;comment:策略实例ID"`
	TradeId            int        `json:"tradeId" gorm:"type:bigint;comment:交易ID"`
	RiskScope          int        `json:"riskScope" gorm:"type:tinyint;comment:风控范围"`
	AssetSymbol        string     `json:"assetSymbol" gorm:"type:varchar(50);comment:风控币种"`
	RiskLevel          int        `json:"riskLevel" gorm:"type:tinyint;comment:风控级别"`
	ManualRecover      int        `json:"manualRecover" gorm:"type:tinyint;comment:是否需要人工恢复"`
	AutoRecoverTime    *time.Time `json:"autoRecoverTime" gorm:"type:timestamp;comment:自动恢复时间"`
	IsRecovered        int        `json:"isRecovered" gorm:"type:tinyint;comment:是否已恢复"`
	RecoveredBy        string     `json:"recoveredBy" gorm:"type:bigint;default:0;comment:审核人"`
	RecoveredAt        *time.Time `json:"recoveredAt" gorm:"type:timestamp;comment:恢复时间"`
	TriggerRule        string     `json:"triggerRule" gorm:"type:varchar(255);comment:触发的风控规则"`
	TriggerValue       string     `json:"triggerValue" gorm:"type:varchar(100);comment:触发值"`
	models.ModelTime
	models.ControlBy
}

func (BusRiskEvent) TableName() string {
	return "bus_risk_event"
}

func (e *BusRiskEvent) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusRiskEvent) GetId() interface{} {
	return e.Id
}
