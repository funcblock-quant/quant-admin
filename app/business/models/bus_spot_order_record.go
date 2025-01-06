package models

import (
	"gorm.io/gorm"
	"quanta-admin/common/models"
	"quanta-admin/common/utils"
)

type BusSpotOrderRecord struct {
	models.Model

	ArbitrageId   string `json:"arbitrageId" gorm:"type:bigint;comment:套利记录id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:交易币种"`
	Side          string `json:"side" gorm:"type:tinyint;comment:买卖方向"`
	OriginQty     string `json:"originQty" gorm:"type:decimal(32,0);comment:原始委托数量"`
	OrderId       string `json:"orderId" gorm:"type:bigint;comment:交易所订单id"`
	OrderClientId string `json:"orderClientId" gorm:"type:varchar(255);comment:策略端id"`
	TimeInForce   string `json:"timeInForce" gorm:"type:tinyint;comment:有效方法"`
	OrderType     string `json:"orderType" gorm:"type:tinyint;comment:订单类型"`
	Fees          string `json:"fees" gorm:"type:decimal(32,0);comment:交易手续费"`
	FeeAsset      string `json:"feeAsset" gorm:"type:decimal(32,0);comment:交易手续费计价单位"`
	Role          string `json:"role" gorm:"type:tinyint;comment:交易角色"`
	Status        string `json:"status" gorm:"type:tinyint;comment:成交状态"`
	models.ModelTime
	models.ControlBy
}

func (BusSpotOrderRecord) TableName() string {
	return "bus_spot_order_record"
}

func (e *BusSpotOrderRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusSpotOrderRecord) GetId() interface{} {
	return e.Id
}

func (e *BusSpotOrderRecord) AfterFind(tx *gorm.DB) (err error) {
	e.OriginQty = utils.ConvertDecimal(e.OriginQty)
	e.Fees = utils.ConvertDecimal(e.Fees)
	return
}
