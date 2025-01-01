package models

import (
	"quanta-admin/common/models"
)

type BusExchangeAccountInfo struct {
	models.Model

	AccountName  string `json:"accountName" gorm:"type:varchar(255);comment:钱包名称"`
	PlatformId   string `json:"platformId" gorm:"type:bigint;comment:id"`
	PlatformName string `json:"platformName" gorm:"type:varchar(255);comment:交易所名称"`
	Uid          string `json:"uid" gorm:"type:varchar(255);comment:交易所uid"`
	AccountType  string `json:"accountType" gorm:"type:tinyint;comment:账户类型"`
	Status       string `json:"status" gorm:"type:tinyint;comment:状态"`
	IsDeleted    string `json:"isDeleted" gorm:"type:tinyint;comment:删除标识位"`
	models.ModelTime
	models.ControlBy
}

func (BusExchangeAccountInfo) TableName() string {
	return "bus_exchange_account_info"
}

func (e *BusExchangeAccountInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusExchangeAccountInfo) GetId() interface{} {
	return e.Id
}
