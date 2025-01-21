package models

import (
	"quanta-admin/common/models"
)

type BusPriceTriggerStrategyApikeyConfig struct {
	models.Model

	UserId   string `json:"userId" gorm:"type:bigint;comment:用户id"`
	ApiKey   string `json:"apiKey" gorm:"type:varchar(255);comment:api key"`
	Username string `json:"username" gorm:"type:varchar(255);comment:用户名"`
	Password string `json:"password" gorm:"type:varchar(255);comment:密码"`
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
