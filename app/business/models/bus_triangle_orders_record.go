package models

import (
	"quanta-admin/common/models"
)

type BusTriangleOrdersRecord struct {
	models.Model

	ArbitrageId   string `json:"arbitrageId" gorm:"type:varchar(255);comment:套利记录id"`
	ExchangeId    string `json:"exchangeId" gorm:"type:bigint;comment:交易所id"`
	ExchangeName  string `json:"exchangeName" gorm:"type:varchar(255);comment:交易所名称"`
	ExchangeType  string `json:"exchangeType" gorm:"type:tinyint;comment:交易平台类型"`
	Side          string `json:"side" gorm:"type:tinyint;comment:买卖方向"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:交易币种"`
	OrderId       string `json:"orderId" gorm:"type:bigint;comment:交易所订单id"`
	OrderClientId string `json:"orderClientId" gorm:"type:varchar(255);comment:策略端id"`
	OriginQty     string `json:"originQty" gorm:"type:decimal(32,0);comment:原始委托数量"`
	OriginPrice   string `json:"originPrice" gorm:"type:decimal(32,0);comment:原始委托价格"`
	OriginType    string `json:"originType" gorm:"type:tinyint;comment:触发前订单类型"`
	TimeInForce   string `json:"timeInForce" gorm:"type:tinyint;comment:有效方法"`
	Role          string `json:"role" gorm:"type:tinyint;comment:交易角色"`
	Status        string `json:"status" gorm:"type:tinyint;comment:持仓状态"`
	Fees          string `json:"fees" gorm:"type:decimal(32,0);comment:交易手续费"`
	FeeAsset      string `json:"feeAsset" gorm:"type:decimal(32,0);comment:交易手续费计价单位"`
	models.ModelTime
	models.ControlBy
}

func (BusTriangleOrdersRecord) TableName() string {
	return "bus_triangle_orders_record"
}

func (e *BusTriangleOrdersRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusTriangleOrdersRecord) GetId() interface{} {
	return e.Id
}
