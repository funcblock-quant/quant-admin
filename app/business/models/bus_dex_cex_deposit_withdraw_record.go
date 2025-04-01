package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusDexCexDepositWithdrawRecord struct {
	models.Model

	OrderType   string    `json:"orderType" gorm:"type:tinyint;comment:订单类型"`
	FromAddress string    `json:"fromAddress" gorm:"type:varchar(255);comment:from地址"`
	ToAddress   string    `json:"toAddress" gorm:"type:varchar(255);comment:to地址"`
	CexUid      string    `json:"cexUid" gorm:"type:varchar(64);comment:cex账户uid"`
	CexOrderId  string    `json:"cexOrderId" gorm:"type:varchar(255);comment:交易所充提业务id"`
	TxHash      string    `json:"txHash" gorm:"type:varchar(255);comment:链上交易hash"`
	Symbol      string    `json:"symbol" gorm:"type:varchar(64);comment:充提币种"`
	Amount      string    `json:"amount" gorm:"type:decimal(32,16);comment:充提数量"`
	Fee         string    `json:"fee" gorm:"type:decimal(32,16);comment:手续费"`
	FeeAsset    string    `json:"feeAsset" gorm:"type:varchar(64);comment:手续费币种"`
	StartTime   time.Time `json:"startTime" gorm:"type:timestamp;comment:发起时间"`
	FinishTime  time.Time `json:"finishTime" gorm:"type:timestamp;comment:完成时间"`
	Status      string    `json:"status" gorm:"type:tinyint;comment:状态"`
	Remark      string    `json:"remark" gorm:"type:text;comment:备注"`
	models.ModelTime
	models.ControlBy
}

func (BusDexCexDepositWithdrawRecord) TableName() string {
	return "bus_dex_cex_deposit_withdraw_record"
}

func (e *BusDexCexDepositWithdrawRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexDepositWithdrawRecord) GetId() interface{} {
	return e.Id
}
