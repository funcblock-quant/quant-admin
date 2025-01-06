package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusStrategyInstance struct {
	models.Model

	StrategyId     string    `json:"strategyId" gorm:"type:bigint;comment:策略id"`
	AccountGroupId string    `json:"accountGroupId" gorm:"type:bigint;comment:账户组id"`
	ExchangeId1    string    `json:"exchangeId1" gorm:"type:bigint;comment:交易所id1"`
	ExchangeName1  string    `json:"exchange1Name" gorm:"type:varchar(64);comment:交易所名称"`
	Exchange1Type  string    `json:"exchange1Type" gorm:"type:tinyint;comment:平台类型"`
	ExchangeId2    string    `json:"exchangeId2" gorm:"type:bigint;comment:交易所id2"`
	ExchangeName2  string    `json:"exchange2Name" gorm:"type:varchar(64);comment:交易所名称"`
	Exchange2Type  string    `json:"exchange2Type" gorm:"type:tinyint;comment:平台类型"`
	InstanceName   string    `json:"instanceName" gorm:"type:varchar(255);comment:策略实例名称"`
	StartRunTime   time.Time `json:"startRunTime" gorm:"type:timestamp;default:NULL;comment:启动时间"`
	StopRunTime    time.Time `json:"stopRunTime" gorm:"type:timestamp;default:NULL;comment:停止时间"`
	ServerIp       string    `json:"serverIp" gorm:"type:varchar(32);comment:服务器ip"`
	ServerName     string    `json:"serverName" gorm:"type:varchar(255);comment:服务器用户名"`
	Status         string    `json:"status" gorm:"type:tinyint;default:0;comment:运行状态"`
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
