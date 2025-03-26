package models

import (
	"quanta-admin/common/models"
)

type BusExchangeAccountInfo struct {
	models.Model

	AccountName        string `gorm:"type:varchar(255);not null;default:'';comment:交易所账户名称" json:"accountName"`
	Uid                string `gorm:"type:varchar(255);not null;comment:交易所账户uid" json:"uid"`
	Email              string `gorm:"type:varchar(255);not null;comment:子账号EMAIL" json:"email"`
	ExchangeType       string `gorm:"type:varchar(64);not null;comment:交易所类型" json:"exchangeType"`
	EncryptedApiKey    string `gorm:"type:varchar(512);comment:加密后的API_KEY" json:"-"`
	EncryptedApiSecret string `gorm:"type:varchar(512);comment:加密后的API_SECRET" json:"-"`
	IsAmberBound       bool   `gorm:"not null;default:false;comment:是否绑定amber" json:"isAmberBound"`
	AmberExchangeType  string `gorm:"type:varchar(64);comment:amber对应的exchange type" json:"amberExchangeType"`
	AmberAccountName   string `gorm:"type:varchar(255);comment:amber对应的account name" json:"amberAccountName"`
	AmberAccountToken  string `gorm:"type:varchar(512);comment:amber对应的account token" json:"-"`
	Status             int8   `gorm:"type:tinyint;not null;default:1;comment:钱包状态，1: 待启用, 2: 已启用" json:"status"`
	Passphrase         string `json:"passphrase" gorm:"type:varchar(255)"`
	MasterAccountId    int    `json:"masterAccountId" gorm:"type:bigint;comment:主账号id"`
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
