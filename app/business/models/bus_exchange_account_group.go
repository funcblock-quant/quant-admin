package models

import (
	"quanta-admin/common/models"
)

type BusExchangeAccountGroup struct {
	models.Model

	GroupName   string `json:"groupName" gorm:"type:varchar(255);comment:交易所账户组"`
	Description string `json:"description" gorm:"type:varchar(255);comment:交易所账户组描述"`
	models.ModelTime
	models.ControlBy
}

func (BusExchangeAccountGroup) TableName() string {
	return "bus_exchange_account_group"
}

func (e *BusExchangeAccountGroup) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusExchangeAccountGroup) GetId() interface{} {
	return e.Id
}
