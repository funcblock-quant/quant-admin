package models

import (
	"gorm.io/gorm"
	"quanta-admin/common/models"
	"quanta-admin/common/utils"
)

type BusDexCexTriangularArbitrageRecord struct {
	models.Model

	StrategyId                string `json:"strategyId" gorm:"type:bigint;comment:策略id"`
	ArbitrageId               string `json:"arbitrageId" gorm:"type:varchar(255);comment:套利记录id"`
	Type                      string `json:"type" gorm:"type:tinyint;comment:套利类型"`
	DexPoolId                 string `json:"dexPoolId" gorm:"type:varchar(64);comment:dex pool id"`
	DexPlatform               string `json:"dexPlatform" gorm:"type:varchar(64);comment:dex平台"`
	DexBlockchain             string `json:"dexBlockchain" gorm:"type:varchar(64);comment:dex区块链"`
	DexSymbol                 string `json:"dexSymbol" gorm:"type:varchar(64);comment:dex交易对"`
	TokenInAmount             string `json:"tokenInAmount" gorm:"type:decimal(32,0);comment:token in 交易量"`
	TokenOutAmount            string `json:"tokenOutAmount" gorm:"type:decimal(32,0);comment:token out 交易量"`
	TxGasAmount               string `json:"txGasAmount" gorm:"type:decimal(32,0);comment:交易手续费"`
	StatusOnDex               string `json:"statusOnDex" gorm:"type:tinyint;comment:dex是否完成"`
	TxHash                    string `json:"txHash" gorm:"type:varchar(255);comment:dex交易hash"`
	CexPlatform               string `json:"cexPlatform" gorm:"type:varchar(64);comment:cex平台"`
	CexSymbolForQuoteToken    string `json:"cexSymbolForQuoteToken" gorm:"type:varchar(64);comment:cex中quote交易对"`
	CexQuantityForQuoteToken  string `json:"cexQuantityForQuoteToken" gorm:"type:decimal(32,0);comment:quote token交易量"`
	CexPriceForQuoteToken     string `json:"cexPriceForQuoteToken" gorm:"type:decimal(32,0);comment:quote token价格"`
	CexFeeAmountForQuoteToken string `json:"cexFeeAmountForQuoteToken" gorm:"type:decimal(32,0);comment:quote token手续费"`
	StatusOnCexForQuoteToken  string `json:"statusOnCexForQuoteToken" gorm:"type:tinyint;comment:cex是否完成"`
	CexSymbolForBaseToken     string `json:"cexSymbolForBaseToken" gorm:"type:varchar(64);comment:cex中base交易对"`
	CexVolumnForBaseToken     string `json:"cexVolumnForBaseToken" gorm:"type:decimal(32,0);comment:base token"`
	CexPriceForBaseToken      string `json:"cexPriceForBaseToken" gorm:"type:decimal(32,0);comment:base token价格"`
	CexFeeAmountForBaseToken  string `json:"cexFeeAmountForBaseToken" gorm:"type:decimal(32,0);comment:base token手续费"`
	StatusOnCexForBaseToken   string `json:"statusOnCexForBaseToken" gorm:"type:tinyint;comment:cex是否完成"`
	QuoteTokenProfit          string `json:"quoteTokenProfit" gorm:"type:decimal(32,0);comment:quote token利润"`
	BaseTokenProfit           string `json:"baseTokenProfit" gorm:"type:decimal(32,0);comment:base token利润"`
	Status                    string `json:"status" gorm:"type:tinyint;comment:套利单状态"`
	models.ModelTime
	models.ControlBy
}

func (BusDexCexTriangularArbitrageRecord) TableName() string {
	return "bus_dex_cex_triangular_arbitrage_record"
}

func (e *BusDexCexTriangularArbitrageRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusDexCexTriangularArbitrageRecord) GetId() interface{} {
	return e.Id
}

func (e *BusDexCexTriangularArbitrageRecord) AfterFind(tx *gorm.DB) (err error) {
	// 假设是金额字段，查询完后进行转换处理
	e.TokenInAmount = utils.ConvertDecimal(e.TokenInAmount)
	e.TokenOutAmount = utils.ConvertDecimal(e.TokenOutAmount)
	e.TxGasAmount = utils.ConvertDecimal(e.TxGasAmount)
	e.CexQuantityForQuoteToken = utils.ConvertDecimal(e.CexQuantityForQuoteToken)
	e.CexPriceForQuoteToken = utils.ConvertDecimal(e.CexPriceForQuoteToken)
	e.CexFeeAmountForQuoteToken = utils.ConvertDecimal(e.CexFeeAmountForQuoteToken)
	e.CexVolumnForBaseToken = utils.ConvertDecimal(e.CexVolumnForBaseToken)
	e.CexPriceForBaseToken = utils.ConvertDecimal(e.CexPriceForBaseToken)
	e.CexFeeAmountForBaseToken = utils.ConvertDecimal(e.CexFeeAmountForBaseToken)
	e.QuoteTokenProfit = utils.ConvertDecimal(e.QuoteTokenProfit)
	e.BaseTokenProfit = utils.ConvertDecimal(e.BaseTokenProfit)
	return
}
