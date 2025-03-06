package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusRiskCheckProgress struct {
	models.Model

	StrategyId         string    `gorm:"type:varchar(64);not null;comment:'策略 ID，关联策略表'" json:"strategyId"`
	InstanceId         string    `gorm:"type:varchar(64);comment:'策略实例 ID（一个策略可能有多个实例）'" json:"instanceId"`
	BusinessType       string    `gorm:"type:varchar(64);not null;comment:'业务类型（如现货交易、合约交易）'" json:"businessType"`
	TradeTable         string    `gorm:"type:varchar(64);not null;comment:'交易表名称（如 spot_trades, futures_trades）'" json:"tradeTable"`
	LastCheckedTradeId uint64    `gorm:"not null;comment:'已检查到的最后一笔交易 ID，下次检查从该 ID 之后继续'" json:"lastCheckedTradeId"`
	LastCheckedAt      time.Time `gorm:"not null;comment:'最近一次风控检查的时间'" json:"lastCheckedAt"`
	Status             int       `gorm:"type:tinyint;not null;default:0;comment:'检查状态（0=未启动，1=进行中，2=完成）'" json:"status"`

	models.ModelTime
	models.ControlBy
}

func (BusRiskCheckProgress) TableName() string {
	return "bus_risk_check_progress"
}

func (e *BusRiskCheckProgress) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusRiskCheckProgress) GetId() interface{} {
	return e.Id
}
