package models

import (
	"quanta-admin/common/models"
)

type BusFuturePositionRecord struct {
	models.Model

	ArbitrageId   string `json:"arbitrageId" gorm:"type:varchar(255);comment:套利记录id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:交易币种"`
	Side          string `json:"side" gorm:"type:tinyint;comment:买卖方向"`
	Leverage      string `json:"leverage" gorm:"type:tinyint;comment:合约杠杆"`
	OrderId       string `json:"orderId" gorm:"type:bigint;comment:交易所订单id"`
	OrderClientId string `json:"orderClientId" gorm:"type:varchar(255);comment:策略端id"`
	OriginQty     string `json:"originQty" gorm:"type:decimal(32,0);comment:原始委托数量"`
	OriginPrice   string `json:"originPrice" gorm:"type:decimal(32,0);comment:原始委托价格"`
	OriginType    string `json:"originType" gorm:"type:tinyint;comment:触发前订单类型"`
	PositionSide  string `json:"positionSide" gorm:"type:varchar(32);comment:持仓方向"`
	TimeInForce   string `json:"timeInForce" gorm:"type:tinyint;comment:有效方法"`
	Role          string `json:"role" gorm:"type:tinyint;comment:交易角色"`
	Status        string `json:"status" gorm:"type:tinyint;comment:持仓状态"`
	Fees          string `json:"fees" gorm:"type:decimal(32,0);comment:交易手续费"`
	FeeAsset      string `json:"feeAsset" gorm:"type:decimal(32,0);comment:交易手续费计价单位"`
	RealizedPnl   string `json:"realizedPnl" gorm:"type:decimal(32,0);comment:已实现盈亏"`
	UnrealizedPnl string `json:"unrealizedPnl" gorm:"type:decimal(32,0);comment:未实现盈亏"`
	models.ModelTime
	models.ControlBy
}

func (BusFuturePositionRecord) TableName() string {
	return "bus_future_position_record"
}

func (e *BusFuturePositionRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusFuturePositionRecord) GetId() interface{} {
	return e.Id
}
