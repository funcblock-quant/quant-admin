package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexCexTriangularArbitrageRecordGetPageReq struct {
	dto.Pagination `search:"-"`
	StrategyId     string `form:"strategyId"  search:"type:exact;column:strategy_id;table:bus_dex_cex_triangular_arbitrage_record" comment:"策略id"`
	ArbitrageId    string `form:"arbitrageId"  search:"type:exact;column:arbitrage_id;table:bus_dex_cex_triangular_arbitrage_record" comment:"套利记录id"`
	Type           string `form:"type"  search:"type:exact;column:type;table:bus_dex_cex_triangular_arbitrage_record" comment:"套利类型"`
	Status         string `form:"status"  search:"type:exact;column:status;table:bus_dex_cex_triangular_arbitrage_record" comment:"套利单状态"`
	BusDexCexTriangularArbitrageRecordOrder
}

type BusDexCexTriangularArbitrageRecordOrder struct {
	Id                        string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_triangular_arbitrage_record"`
	StrategyId                string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_dex_cex_triangular_arbitrage_record"`
	ArbitrageId               string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_dex_cex_triangular_arbitrage_record"`
	Type                      string `form:"typeOrder"  search:"type:order;column:type;table:bus_dex_cex_triangular_arbitrage_record"`
	DexPoolId                 string `form:"dexPoolIdOrder"  search:"type:order;column:dex_pool_id;table:bus_dex_cex_triangular_arbitrage_record"`
	DexPlatform               string `form:"dexPlatformOrder"  search:"type:order;column:dex_platform;table:bus_dex_cex_triangular_arbitrage_record"`
	DexBlockchain             string `form:"dexBlockchainOrder"  search:"type:order;column:dex_blockchain;table:bus_dex_cex_triangular_arbitrage_record"`
	DexSymbol                 string `form:"dexSymbolOrder"  search:"type:order;column:dex_symbol;table:bus_dex_cex_triangular_arbitrage_record"`
	TokenInAmount             string `form:"tokenInAmountOrder"  search:"type:order;column:token_in_amount;table:bus_dex_cex_triangular_arbitrage_record"`
	TokenOutAmount            string `form:"tokenOutAmountOrder"  search:"type:order;column:token_out_amount;table:bus_dex_cex_triangular_arbitrage_record"`
	TxGasAmount               string `form:"txGasAmountOrder"  search:"type:order;column:tx_gas_amount;table:bus_dex_cex_triangular_arbitrage_record"`
	StatusOnDex               string `form:"statusOnDexOrder"  search:"type:order;column:status_on_dex;table:bus_dex_cex_triangular_arbitrage_record"`
	TxHash                    string `form:"txHashOrder"  search:"type:order;column:tx_hash;table:bus_dex_cex_triangular_arbitrage_record"`
	CexPlatform               string `form:"cexPlatformOrder"  search:"type:order;column:cex_platform;table:bus_dex_cex_triangular_arbitrage_record"`
	CexSymbolForQuoteToken    string `form:"cexSymbolForQuoteTokenOrder"  search:"type:order;column:cex_symbol_for_quote_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexQuantityForQuoteToken  string `form:"cexQuantityForQuoteTokenOrder"  search:"type:order;column:cex_quantity_for_quote_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexPriceForQuoteToken     string `form:"cexPriceForQuoteTokenOrder"  search:"type:order;column:cex_price_for_quote_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexFeeAmountForQuoteToken string `form:"cexFeeAmountForQuoteTokenOrder"  search:"type:order;column:cex_fee_amount_for_quote_token;table:bus_dex_cex_triangular_arbitrage_record"`
	StatusOnCexForQuoteToken  string `form:"statusOnCexForQuoteTokenOrder"  search:"type:order;column:status_on_cex_for_quote_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexSymbolForBaseToken     string `form:"cexSymbolForBaseTokenOrder"  search:"type:order;column:cex_symbol_for_base_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexVolumnForBaseToken     string `form:"cexVolumnForBaseTokenOrder"  search:"type:order;column:cex_volumn_for_base_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexPriceForBaseToken      string `form:"cexPriceForBaseTokenOrder"  search:"type:order;column:cex_price_for_base_token;table:bus_dex_cex_triangular_arbitrage_record"`
	CexFeeAmountForBaseToken  string `form:"cexFeeAmountForBaseTokenOrder"  search:"type:order;column:cex_fee_amount_for_base_token;table:bus_dex_cex_triangular_arbitrage_record"`
	StatusOnCexForBaseToken   string `form:"statusOnCexForBaseTokenOrder"  search:"type:order;column:status_on_cex_for_base_token;table:bus_dex_cex_triangular_arbitrage_record"`
	QuoteTokenProfit          string `form:"quoteTokenProfitOrder"  search:"type:order;column:quote_token_profit;table:bus_dex_cex_triangular_arbitrage_record"`
	BaseTokenProfit           string `form:"baseTokenProfitOrder"  search:"type:order;column:base_token_profit;table:bus_dex_cex_triangular_arbitrage_record"`
	Status                    string `form:"statusOrder"  search:"type:order;column:status;table:bus_dex_cex_triangular_arbitrage_record"`
	CreatedAt                 string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_triangular_arbitrage_record"`
	UpdatedAt                 string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_dex_cex_triangular_arbitrage_record"`
	DeletedAt                 string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_dex_cex_triangular_arbitrage_record"`
}

func (m *BusDexCexTriangularArbitrageRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexTriangularArbitrageRecordInsertReq struct {
	Id                        int    `json:"-" comment:""` //
	StrategyId                string `json:"strategyId" comment:"策略id"`
	ArbitrageId               string `json:"arbitrageId" comment:"套利记录id"`
	Type                      string `json:"type" comment:"套利类型"`
	DexPoolId                 string `json:"dexPoolId" comment:"dex pool id"`
	DexPlatform               string `json:"dexPlatform" comment:"dex平台"`
	DexBlockchain             string `json:"dexBlockchain" comment:"dex区块链"`
	DexSymbol                 string `json:"dexSymbol" comment:"dex交易对"`
	TokenInAmount             string `json:"tokenInAmount" comment:"token in 交易量"`
	TokenOutAmount            string `json:"tokenOutAmount" comment:"token out 交易量"`
	TxGasAmount               string `json:"txGasAmount" comment:"交易手续费"`
	StatusOnDex               string `json:"statusOnDex" comment:"dex是否完成"`
	TxHash                    string `json:"txHash" comment:"dex交易hash"`
	CexPlatform               string `json:"cexPlatform" comment:"cex平台"`
	CexSymbolForQuoteToken    string `json:"cexSymbolForQuoteToken" comment:"cex中quote交易对"`
	CexQuantityForQuoteToken  string `json:"cexQuantityForQuoteToken" comment:"quote token交易量"`
	CexPriceForQuoteToken     string `json:"cexPriceForQuoteToken" comment:"quote token价格"`
	CexFeeAmountForQuoteToken string `json:"cexFeeAmountForQuoteToken" comment:"quote token手续费"`
	StatusOnCexForQuoteToken  string `json:"statusOnCexForQuoteToken" comment:"cex是否完成"`
	CexSymbolForBaseToken     string `json:"cexSymbolForBaseToken" comment:"cex中base交易对"`
	CexVolumnForBaseToken     string `json:"cexVolumnForBaseToken" comment:"base token"`
	CexPriceForBaseToken      string `json:"cexPriceForBaseToken" comment:"base token价格"`
	CexFeeAmountForBaseToken  string `json:"cexFeeAmountForBaseToken" comment:"base token手续费"`
	StatusOnCexForBaseToken   string `json:"statusOnCexForBaseToken" comment:"cex是否完成"`
	QuoteTokenProfit          string `json:"quoteTokenProfit" comment:"quote token利润"`
	BaseTokenProfit           string `json:"baseTokenProfit" comment:"base token利润"`
	Status                    string `json:"status" comment:"套利单状态"`
	common.ControlBy
}

func (s *BusDexCexTriangularArbitrageRecordInsertReq) Generate(model *models.BusDexCexTriangularArbitrageRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.ArbitrageId = s.ArbitrageId
	model.Type = s.Type
	model.DexPoolId = s.DexPoolId
	model.DexPlatform = s.DexPlatform
	model.DexBlockchain = s.DexBlockchain
	model.DexSymbol = s.DexSymbol
	model.TokenInAmount = s.TokenInAmount
	model.TokenOutAmount = s.TokenOutAmount
	model.TxGasAmount = s.TxGasAmount
	model.StatusOnDex = s.StatusOnDex
	model.TxHash = s.TxHash
	model.CexPlatform = s.CexPlatform
	model.CexSymbolForQuoteToken = s.CexSymbolForQuoteToken
	model.CexQuantityForQuoteToken = s.CexQuantityForQuoteToken
	model.CexPriceForQuoteToken = s.CexPriceForQuoteToken
	model.CexFeeAmountForQuoteToken = s.CexFeeAmountForQuoteToken
	model.StatusOnCexForQuoteToken = s.StatusOnCexForQuoteToken
	model.CexSymbolForBaseToken = s.CexSymbolForBaseToken
	model.CexVolumnForBaseToken = s.CexVolumnForBaseToken
	model.CexPriceForBaseToken = s.CexPriceForBaseToken
	model.CexFeeAmountForBaseToken = s.CexFeeAmountForBaseToken
	model.StatusOnCexForBaseToken = s.StatusOnCexForBaseToken
	model.QuoteTokenProfit = s.QuoteTokenProfit
	model.BaseTokenProfit = s.BaseTokenProfit
	model.Status = s.Status
}

func (s *BusDexCexTriangularArbitrageRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexTriangularArbitrageRecordUpdateReq struct {
	Id                        int    `uri:"id" comment:""` //
	StrategyId                string `json:"strategyId" comment:"策略id"`
	ArbitrageId               string `json:"arbitrageId" comment:"套利记录id"`
	Type                      string `json:"type" comment:"套利类型"`
	DexPoolId                 string `json:"dexPoolId" comment:"dex pool id"`
	DexPlatform               string `json:"dexPlatform" comment:"dex平台"`
	DexBlockchain             string `json:"dexBlockchain" comment:"dex区块链"`
	DexSymbol                 string `json:"dexSymbol" comment:"dex交易对"`
	TokenInAmount             string `json:"tokenInAmount" comment:"token in 交易量"`
	TokenOutAmount            string `json:"tokenOutAmount" comment:"token out 交易量"`
	TxGasAmount               string `json:"txGasAmount" comment:"交易手续费"`
	StatusOnDex               string `json:"statusOnDex" comment:"dex是否完成"`
	TxHash                    string `json:"txHash" comment:"dex交易hash"`
	CexPlatform               string `json:"cexPlatform" comment:"cex平台"`
	CexSymbolForQuoteToken    string `json:"cexSymbolForQuoteToken" comment:"cex中quote交易对"`
	CexQuantityForQuoteToken  string `json:"cexQuantityForQuoteToken" comment:"quote token交易量"`
	CexPriceForQuoteToken     string `json:"cexPriceForQuoteToken" comment:"quote token价格"`
	CexFeeAmountForQuoteToken string `json:"cexFeeAmountForQuoteToken" comment:"quote token手续费"`
	StatusOnCexForQuoteToken  string `json:"statusOnCexForQuoteToken" comment:"cex是否完成"`
	CexSymbolForBaseToken     string `json:"cexSymbolForBaseToken" comment:"cex中base交易对"`
	CexVolumnForBaseToken     string `json:"cexVolumnForBaseToken" comment:"base token"`
	CexPriceForBaseToken      string `json:"cexPriceForBaseToken" comment:"base token价格"`
	CexFeeAmountForBaseToken  string `json:"cexFeeAmountForBaseToken" comment:"base token手续费"`
	StatusOnCexForBaseToken   string `json:"statusOnCexForBaseToken" comment:"cex是否完成"`
	QuoteTokenProfit          string `json:"quoteTokenProfit" comment:"quote token利润"`
	BaseTokenProfit           string `json:"baseTokenProfit" comment:"base token利润"`
	Status                    string `json:"status" comment:"套利单状态"`
	common.ControlBy
}

func (s *BusDexCexTriangularArbitrageRecordUpdateReq) Generate(model *models.BusDexCexTriangularArbitrageRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.ArbitrageId = s.ArbitrageId
	model.Type = s.Type
	model.DexPoolId = s.DexPoolId
	model.DexPlatform = s.DexPlatform
	model.DexBlockchain = s.DexBlockchain
	model.DexSymbol = s.DexSymbol
	model.TokenInAmount = s.TokenInAmount
	model.TokenOutAmount = s.TokenOutAmount
	model.TxGasAmount = s.TxGasAmount
	model.StatusOnDex = s.StatusOnDex
	model.TxHash = s.TxHash
	model.CexPlatform = s.CexPlatform
	model.CexSymbolForQuoteToken = s.CexSymbolForQuoteToken
	model.CexQuantityForQuoteToken = s.CexQuantityForQuoteToken
	model.CexPriceForQuoteToken = s.CexPriceForQuoteToken
	model.CexFeeAmountForQuoteToken = s.CexFeeAmountForQuoteToken
	model.StatusOnCexForQuoteToken = s.StatusOnCexForQuoteToken
	model.CexSymbolForBaseToken = s.CexSymbolForBaseToken
	model.CexVolumnForBaseToken = s.CexVolumnForBaseToken
	model.CexPriceForBaseToken = s.CexPriceForBaseToken
	model.CexFeeAmountForBaseToken = s.CexFeeAmountForBaseToken
	model.StatusOnCexForBaseToken = s.StatusOnCexForBaseToken
	model.QuoteTokenProfit = s.QuoteTokenProfit
	model.BaseTokenProfit = s.BaseTokenProfit
	model.Status = s.Status
}

func (s *BusDexCexTriangularArbitrageRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexTriangularArbitrageRecordGetReq 功能获取请求参数
type BusDexCexTriangularArbitrageRecordGetReq struct {
	Id int `uri:"id"`
}

func (s *BusDexCexTriangularArbitrageRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexTriangularArbitrageRecordDeleteReq 功能删除请求参数
type BusDexCexTriangularArbitrageRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexCexTriangularArbitrageRecordDeleteReq) GetId() interface{} {
	return s.Ids
}

type BusArbitrageOpportunityGetReq struct {
	StrategyInstanceId string   `json:"strategyInstanceId"`
	Symbols            []string `form:"symbols" `
	MinProfit          string   `form:"minProfit"`
	MaxProfit          string   `form:"maxProfit"`
	BeginTime          string   `form:"beginTime"`
	EndTime            string   `form:"endTime"`
}

type BusArbitrageOpportunityGetResp struct {
	Exchange1 string `json:"exchange1"`
	Exchange2 string `json:"exchange2"`
	Symbol    string `json:"symbol" gorm:"column:symbol"`
	Count     int    `json:"count" gorm:"column:count"`
}
