package models

import (
	"quanta-admin/common/models"
)

type StrategyDexCexTriangularArbitrageOpportunities struct {
	models.Model

	InstanceId         string `json:"instanceId" gorm:"type:varchar(64);comment:实例id"`
	OpportunityId      string `json:"opportunityId" gorm:"type:varchar(64);comment:Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" gorm:"type:tinyint(1);comment:买方标识"`
	DexTargetToken     string `json:"dexTargetToken" gorm:"type:varchar(64);comment:Dex target token pubkey"`
	DexPoolType        string `json:"dexPoolType" gorm:"type:varchar(32);comment:Dex pool type"`
	DexPoolId          string `json:"dexPoolId" gorm:"type:varchar(64);comment:Dex pool pubkey"`
	DexSolAmount       string `json:"dexSolAmount" gorm:"type:double;comment:Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" gorm:"type:double;comment:Dex trade target token amount"`
	CexExchangeType    string `json:"cexExchangeType" gorm:"type:varchar(32);comment:Cex exchange type"`
	CexTargetAsset     string `json:"cexTargetAsset" gorm:"type:varchar(32);comment:Cex target asset"`
	CexQuoteAsset      string `json:"cexQuoteAsset" gorm:"type:varchar(32);comment:Cex quote asset"`
	CexSellQuantity    string `json:"cexSellQuantity" gorm:"type:double;comment:Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" gorm:"type:double;comment:Cex sell quote amount"`
	CexBuyQuantity     string `json:"cexBuyQuantity" gorm:"type:double;comment:Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" gorm:"type:double;comment:Cex buy quote amount"`
	models.ModelTime
	models.ControlBy
}

func (StrategyDexCexTriangularArbitrageOpportunities) TableName() string {
	return "strategy_dex_cex_triangular_arbitrage_opportunities"
}

func (e *StrategyDexCexTriangularArbitrageOpportunities) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StrategyDexCexTriangularArbitrageOpportunities) GetId() interface{} {
	return e.Id
}
