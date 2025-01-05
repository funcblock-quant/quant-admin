package models

import (
	"quanta-admin/common/models"
)

type BusExchangeAccountGroupRelation struct {
	models.Model

	GroupId   string `json:"groupId" gorm:"type:bigint;comment:交易所账户组id"`
	AccountId string `json:"accountId" gorm:"type:bigint;comment:交易所账户id"`
	models.ModelTime
	models.ControlBy
}

func (BusExchangeAccountGroupRelation) TableName() string {
	return "bus_exchange_account_group_relation"
}

func (e *BusExchangeAccountGroupRelation) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusExchangeAccountGroupRelation) GetId() interface{} {
	return e.Id
}
