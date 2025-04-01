package models

import (

	"quanta-admin/common/models"

)

type BusDexCexDebitCreditRecord struct {
    models.Model
    
    AccountName string `json:"accountName" gorm:"type:varchar(255);comment:交易所账户名称"` 
    Uid string `json:"uid" gorm:"type:varchar(255);comment:交易所账户uid"` 
    ExchangeType string `json:"exchangeType" gorm:"type:varchar(64);comment:交易所"` 
    Symbol string `json:"symbol" gorm:"type:varchar(64);comment:借贷币种"` 
    Amount string `json:"amount" gorm:"type:decimal(32,16);comment:借还贷数量"` 
    DebitType string `json:"debitType" gorm:"type:tinyint;comment:类型"` 
    Status string `json:"status" gorm:"type:tinyint;comment:状态"` 
    models.ModelTime
    models.ControlBy
}

func (BusDexCexDebitCreditRecord) TableName() string {
    return "bus_dex_cex_debit_credit_record"
}

func (e *BusDexCexDebitCreditRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexDebitCreditRecord) GetId() interface{} {
	return e.Id
}