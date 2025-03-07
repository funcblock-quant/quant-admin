package models

import (
	"time"

	"quanta-admin/common/models"
)

type BusPriceTriggerStrategyInstance struct {
	models.Model

	OpenPrice         string    `json:"openPrice" gorm:"type:varchar(32);comment:开仓价格"`
	ClosePrice        string    `json:"closePrice" gorm:"type:varchar(32);comment:平仓价格"`
	Amount            string    `json:"amount" gorm:"type:varchar(32:);comment:开仓数量"`
	Side              string    `json:"side" gorm:"type:varchar(16);comment:买卖方向"`
	Symbol            string    `json:"symbol" gorm:"type:varchar(64);comment:交易币种"`
	CloseTime         time.Time `json:"closeTime" gorm:"type:timestamp;comment:停止时间"`
	ApiConfig         int       `json:"apiConfig" gorm:"type:timestamp;comment:api配置id"`
	Status            string    `json:"status" gorm:"type:varchar(16);comment:状态，created, started, stopped, closed"`
	ExchangeUserId    string    `json:"exchangeUserId" gorm:"type:varchar(255);comment:交易所userId"`
	ExecuteNum        int       `json:"executeNum" gorm:"type:int;comment:执行次数"`
	ProfitTargetType  string    `json:"profitTargetType" gorm:"type:varchar(32);comment:止盈类型"`
	ProfitTargetPrice string    `json:"profitTargetPrice" gorm:"default:0;type:float;comment:限价止盈价格"`
	LossTargetPrice   string    `json:"lossTargetPrice" gorm:"default:0;type:float;comment:限价止盈价格"`
	CallbackRatio     *float64  `json:"callbackRatio" gorm:"type:float;comment:浮动止盈回调比例"`
	CutoffRatio       *float64  `json:"cutoffRatio" gorm:"type:float;comment:浮动止盈止盈比例"`
	MinProfit         *float64  `json:"minProfit" gorm:"type:float;comment:浮动止盈最低利润"`
	CloseOrderType    string    `json:"closeOrderType" gorm:"type:varchar(16);comment:平仓模式"`
	DelayTime         int       `json:"delayTime" gorm:"type:int;comment:延迟时间"`
	models.ModelTime
	models.ControlBy
}

func (BusPriceTriggerStrategyInstance) TableName() string {
	return "bus_price_trigger_strategy_instance"
}

func (e *BusPriceTriggerStrategyInstance) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusPriceTriggerStrategyInstance) GetId() interface{} {
	return e.Id
}
