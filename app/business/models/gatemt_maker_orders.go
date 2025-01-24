package models

import (
	"quanta-admin/common/models"
)

type GatemtMakerOrders struct {
	models.Model

	ClientOrderId    string `json:"clientOrderId" gorm:"type:varchar(50);comment:Client order id"`
	Symbol           string `json:"symbol" gorm:"type:varchar(10);comment:Symbol"`
	OrderSide        string `json:"orderSide" gorm:"type:varchar(10);comment:OrderSide"`
	Price            string `json:"price" gorm:"type:decimal(28,10);comment:Price"`
	Amount           string `json:"amount" gorm:"type:decimal(28,10);comment:Amount"`
	LocalCreatedTime string `json:"localCreatedTime" gorm:"type:varchar(50);comment:gatemt order created time"`
	LocalCreatedAt   string `json:"localCreatedAt" gorm:"type:bigint;comment:gatemt order created at"`
	ExchangeOrderId  string `json:"exchangeOrderId" gorm:"type:varchar(50);comment:ExchangeOrderId"`
	Data             string `json:"data" gorm:"type:json;comment:Data"`
}

func (GatemtMakerOrders) TableName() string {
	return "gatemt_maker_orders"
}

func (e *GatemtMakerOrders) GetId() interface{} {
	return e.Id
}
