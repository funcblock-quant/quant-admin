package models

import (
	"encoding/json"
	"quanta-admin/common/models"
	"time"
)

type BusArbitrageRecord struct {
	models.Model

	ArbitrageId        string `json:"arbitrageId" gorm:"type:varchar(255);comment:套利记录id"`
	StrategyInstanceId string `json:"strategyInstanceId" gorm:"type:bigint;comment:套利单所属策略id"`
	StrategyName       string `json:"strategyName" gorm:"type:varchar(255);comment:策略名称"`
	Type               string `json:"type" gorm:"type:tinyint;comment:套利类型, 0-模拟盘观测, 1-实盘套利"`
	Symbol             string `json:"symbol" gorm:"type:varchar(255);comment:交易对"`
	ContractType       string `json:"contractType" gorm:"type:tinyint;comment:合约类型"`
	RealizedPnl        string `json:"realizedPnl" gorm:"type:decimal(32,0);comment:已实现盈亏"`
	UnrealizedPnl      string `json:"unrealizedPnl" gorm:"type:decimal(32,0);comment:未实现盈亏"`
	ExpectPnl          string `json:"expectPnl" gorm:"type:decimal(32,0);comment:预期盈亏"`
	ExpectPnlPercent   string `json:"expectPnlPercent" gorm:"type:float;comment:预期收益率"`
	StartTime          int64  `json:"-" gorm:"type:bigint;comment:套利单开始时间,单位：ns"`
	EndTime            int64  `json:"-" gorm:"type:bigint;comment:套利单结束时间,单位：ns"`
	Status             string `json:"status" gorm:"type:tinyint;comment:套利单状态, 0-套利中，1-套利完成, 2-对账完成"`
	models.ModelTime
	models.ControlBy
}

func (BusArbitrageRecord) TableName() string {
	return "bus_arbitrage_record"
}

func (e *BusArbitrageRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusArbitrageRecord) GetId() interface{} {
	return e.Id
}

func (b BusArbitrageRecord) MarshalJSON() ([]byte, error) {
	// 时间格式
	timeFormat := "2006-01-02 15:04:05.000000"

	// 格式化时间
	formatTime := func(ns int64) string {
		if ns == 0 {
			return ""
		}
		return time.Unix(0, ns).Format(timeFormat)
	}

	type Alias BusArbitrageRecord
	return json.Marshal(&struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		*Alias
	}{
		StartTime: formatTime(b.StartTime),
		EndTime:   formatTime(b.EndTime),
		Alias:     (*Alias)(&b),
	})
}
