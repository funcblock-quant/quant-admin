package models

import (
	"quanta-admin/common/models"
)

type StrategyDexCexTriangularArbitrageTrades struct {
	models.Model

	InstanceId         string `json:"instanceId" gorm:"type:varchar(64);comment:Arbitrager instance ID"`
	OpportunityId      string `json:"opportunityId" gorm:"type:varchar(64);comment:Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" gorm:"type:tinyint(1);comment:Buy on dex or cex"`
	Error              string `json:"error" gorm:"type:varchar(255);comment:Error message"`
	DexTrader          string `json:"dexTrader" gorm:"type:varchar(64);comment:Dex trader pubkey"`
	DexSuccess         string `json:"dexSuccess" gorm:"type:tinyint(1);comment:Is dex trade success"`
	DexTxFee           string `json:"dexTxFee" gorm:"type:double;comment:Dex trade tx fee"`
	DexTxSig           string `json:"dexTxSig" gorm:"type:varchar(128);comment:Dex trade tx signature"`
	DexJitoFee         string `json:"dexJitoFee" gorm:"type:double;comment:Dex trade tx jito fee"`
	DexSolAmount       string `json:"dexSolAmount" gorm:"type:double;comment:Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" gorm:"type:double;comment:Dex trade target amount"`
	CexAccountId       string `json:"cexAccountId" gorm:"type:varchar(64);comment:Cex account id"`
	CexAccountType     string `json:"cexAccountType" gorm:"type:varchar(32);comment:Cex account type"`
	CexExchangeType    string `json:"cexExchangeType" gorm:"type:varchar(32);comment:Cex exchange type"`
	CexSellSuccess     string `json:"cexSellSuccess" gorm:"type:tinyint(1);comment:Is cex sell success"`
	CexSellOrderId     string `json:"cexSellOrderId" gorm:"type:varchar(64);comment:Cex sell order id"`
	CexSellQuantity    string `json:"cexSellQuantity" gorm:"type:double;comment:Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" gorm:"type:double;comment:Cex sell quote amount"`
	CexSellFeeAsset    string `json:"cexSellFeeAsset" gorm:"type:varchar(32);comment:Cex sell fee asset"`
	CexSellFee         string `json:"cexSellFee" gorm:"type:double;comment:Cex sell fee"`
	CexBuySuccess      string `json:"cexBuySuccess" gorm:"type:tinyint(1);comment:Is cex buy success"`
	CexBuyOrderId      string `json:"cexBuyOrderId" gorm:"type:varchar(64);comment:Cex buy order id"`
	CexBuyQuantity     string `json:"cexBuyQuantity" gorm:"type:double;comment:Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" gorm:"type:double;comment:Cex buy quote amount"`
	CexBuyFeeAsset     string `json:"cexBuyFeeAsset" gorm:"type:varchar(32);comment:Cex buy fee asset"`
	CexBuyFee          string `json:"cexBuyFee" gorm:"type:double;comment:Cex buy fee"`
	IsRiskChecked      bool   `gorm:"null;comment:是否完成风控检查" json:"isRiskChecked"`
	models.ModelTime
	models.ControlBy
}

func (StrategyDexCexTriangularArbitrageTrades) TableName() string {
	return "strategy_dex_cex_triangular_arbitrage_trades"
}

func (e *StrategyDexCexTriangularArbitrageTrades) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StrategyDexCexTriangularArbitrageTrades) GetId() interface{} {
	return e.Id
}
