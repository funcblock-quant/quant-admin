package models

import (
	"quanta-admin/common/models"
)

type BusPriceTriggerStrategyApikeyConfig struct {
	models.Model

	UserId      string `json:"userId" gorm:"type:bigint;comment:用户id"`
	ApiKey      string `json:"apiKey" gorm:"type:varchar(255);comment:api key"`
	SecretKey   string `json:"secretKey" gorm:"type:varchar(255);comment:私钥"`
	AccountName string `json:"accountName" gorm:"type:varchar(255);comment:用户名"`
	Exchange    string `json:"exchange" gorm:"type:varchar(32);comment:交易所"`
	models.ModelTime
	models.ControlBy
}

func (BusPriceTriggerStrategyApikeyConfig) TableName() string {
	return "bus_price_trigger_strategy_apikey_config"
}

func (e *BusPriceTriggerStrategyApikeyConfig) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusPriceTriggerStrategyApikeyConfig) GetId() interface{} {
	return e.Id
}
