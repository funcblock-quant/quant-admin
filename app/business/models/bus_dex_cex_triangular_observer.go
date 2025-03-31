package models

import (
	"errors"
	"quanta-admin/common/models"
)

type BusDexCexTriangularObserver struct {
	models.Model

	StrategyInstanceId string   `gorm:"not null;comment:策略id" json:"strategyInstanceId"`
	InstanceId         string   `gorm:"not null;comment:观察器id" json:"instanceId"`
	Symbol             string   `gorm:"not null;comment:观察币种" json:"symbol"`
	TargetToken        string   `gorm:"not null" json:"targetToken"`
	QuoteToken         string   `gorm:"not null" json:"quoteToken"`
	SymbolConnector    string   `gorm:"not null" json:"symbolConnector"`
	ExchangeType       string   `gorm:"not null;comment:交易所类型" json:"exchangeType"`
	DexType            string   `json:"dexType"`
	MaxArraySize       int      `gorm:"null;default:5" json:"maxArraySize"`
	MinQuoteAmount     *float64 `gorm:"comment:最小交易量" json:"minQuoteAmount"` // 使用指针类型允许值为null
	MaxQuoteAmount     *float64 `gorm:"comment:最大交易量" json:"maxQuoteAmount"` // 使用指针类型允许值为null
	//TriggerHoldingMs       int      `gorm:"null;default:0" json:"triggerHoldingMs"`
	TakerFee                   *float64 `gorm:"not null;comment:交易所taker 费率" json:"takerFee"`
	AmmPoolId                  *string  `gorm:"comment:ammPoolId" json:"ammPoolId"`
	TokenMint                  *string  `gorm:"comment:base token合约" json:"tokenMint"` // 使用指针类型允许值为null
	SlippageBpsRate            *float64 `gorm:"null;default:0;comment:滑点bps" json:"slippageBpsRate"`
	Depth                      string   `gorm:"not null;default:20;comment:深度" json:"depth"`
	IsTrading                  bool     `gorm:"default:false;comment:是否启动交易" json:"isTrading"`
	ProfitTriggerRate          *float64 `gorm:"null;comment:最小利润" json:"profitTriggerRate"`
	PriorityFee                *float64 `gorm:"null;comment:交易优先费" json:"priorityFee"`
	JitoFeeRate                *float64 `gorm:"null;comment:jito交易费比例" json:"jitoFeeRate"`
	Status                     string   `gorm:"not null;comment:状态，0-新增，1-已开启观察，2-水位调节中，3-已启动交易，4-已停止" json:"status"` // 使用 uint8 更合适
	OwnerProgram               *string  `gorm:"null;" json:"ownerProgram"`                                            // 使用 uint8 更合适
	Decimals                   int      `gorm:"null;comment:币种精度" json:"decimals"`
	AlertThreshold             *float64 `gorm:"null;comment:最低预警水位" json:"alertThreshold"`
	BuyTriggerThreshold        *float64 `gorm:"null;comment:触发水位调节的低水位线" json:"buyTriggerThreshold"`
	TargetBalanceThreshold     *float64 `gorm:"null;comment:低水位调节的目标值" json:"targetBalanceThreshold"`
	SellTriggerThreshold       *float64 `gorm:"null;comment:触发水位调节的高水位线" json:"sellTriggerThreshold"`
	MinDepositAmountThreshold  *float64 `gorm:"null;comment:最低充值金额阈值" json:"minDepositAmountThreshold"`
	MinWithdrawAmountThreshold *float64 `gorm:"null;comment:最低提现金额阈值" json:"minWithdrawAmountThreshold"`
	IsTradingBlocked           bool     `gorm:"null;comment:交易功能是否被风控" json:"isTradingBlocked"`
	DexWalletId                *int64   `gorm:"null;comment:dex交易钱包id" json:"dexWallet"`
	CexAccountId               *int64   `gorm:"null;comment:cex交易账户id" json:"cexAccount"`
	models.ModelTime
	models.ControlBy
}

func (BusDexCexTriangularObserver) TableName() string {
	return "bus_dex_cex_triangular_instance"
}

func (e *BusDexCexTriangularObserver) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexTriangularObserver) GetId() interface{} {
	return e.Id
}

func (e *BusDexCexTriangularObserver) GetExchangeTypeForTrader() (string, error) {
	if e.ExchangeType == "Binance" {
		return "BinanceClassicUnifiedMargin", nil
	} else {
		return "", errors.New("not support exchange type")
	}
}
