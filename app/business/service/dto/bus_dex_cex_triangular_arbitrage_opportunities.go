package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type StrategyDexCexTriangularArbitrageOpportunitiesGetPageReq struct {
	dto.Pagination `search:"-"`
	InstanceId     string `form:"instanceId"  search:"type:exact;column:instance_id;table:strategy_dex_cex_triangular_arbitrage_opportunities" comment:"实例id"`
	BuyOnDex       string `form:"buyOnDex"  search:"type:exact;column:buy_on_dex;table:strategy_dex_cex_triangular_arbitrage_opportunities" comment:"买方标识"`
	StrategyDexCexTriangularArbitrageOpportunitiesOrder
}

type StrategyDexCexTriangularArbitrageOpportunitiesOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	InstanceId         string `form:"instanceIdOrder"  search:"type:order;column:instance_id;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	OpportunityId      string `form:"opportunityIdOrder"  search:"type:order;column:opportunity_id;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	BuyOnDex           string `form:"buyOnDexOrder"  search:"type:order;column:buy_on_dex;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexTargetToken     string `form:"dexTargetTokenOrder"  search:"type:order;column:dex_target_token;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexPoolType        string `form:"dexPoolTypeOrder"  search:"type:order;column:dex_pool_type;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexPoolId          string `form:"dexPoolIdOrder"  search:"type:order;column:dex_pool_id;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexSolAmount       string `form:"dexSolAmountOrder"  search:"type:order;column:dex_sol_amount;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexTargetAmount    string `form:"dexTargetAmountOrder"  search:"type:order;column:dex_target_amount;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexTxPriorityFee   string `form:"dexTxPriorityFeeOrder"  search:"type:order;column:dex_tx_priority_fee;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DexTxJitoFee       string `form:"dexTxJitoFeeOrder"  search:"type:order;column:dex_tx_jito_fee;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexExchangeType    string `form:"cexExchangeTypeOrder"  search:"type:order;column:cex_exchange_type;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexTargetAsset     string `form:"cexTargetAssetOrder"  search:"type:order;column:cex_target_asset;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexQuoteAsset      string `form:"cexQuoteAssetOrder"  search:"type:order;column:cex_quote_asset;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexSellQuantity    string `form:"cexSellQuantityOrder"  search:"type:order;column:cex_sell_quantity;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexSellQuoteAmount string `form:"cexSellQuoteAmountOrder"  search:"type:order;column:cex_sell_quote_amount;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexBuyQuantity     string `form:"cexBuyQuantityOrder"  search:"type:order;column:cex_buy_quantity;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CexBuyQuoteAmount  string `form:"cexBuyQuoteAmountOrder"  search:"type:order;column:cex_buy_quote_amount;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:strategy_dex_cex_triangular_arbitrage_opportunities"`
}

func (m *StrategyDexCexTriangularArbitrageOpportunitiesGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StrategyDexCexTriangularArbitrageOpportunitiesInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	InstanceId         string `json:"instanceId" comment:"实例id"`
	OpportunityId      string `json:"opportunityId" comment:"Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" comment:"买方标识"`
	DexTargetToken     string `json:"dexTargetToken" comment:"Dex target token pubkey"`
	DexPoolType        string `json:"dexPoolType" comment:"Dex pool type"`
	DexPoolId          string `json:"dexPoolId" comment:"Dex pool pubkey"`
	DexSolAmount       string `json:"dexSolAmount" comment:"Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" comment:"Dex trade target token amount"`
	DexTxPriorityFee   string `json:"dexTxPriorityFee" comment:"Dex trade tx priority fee"`
	DexTxJitoFee       string `json:"dexTxJitoFee" comment:"Dex trade tx jito fee"`
	CexExchangeType    string `json:"cexExchangeType" comment:"Cex exchange type"`
	CexTargetAsset     string `json:"cexTargetAsset" comment:"Cex target asset"`
	CexQuoteAsset      string `json:"cexQuoteAsset" comment:"Cex quote asset"`
	CexSellQuantity    string `json:"cexSellQuantity" comment:"Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" comment:"Cex sell quote amount"`
	CexBuyQuantity     string `json:"cexBuyQuantity" comment:"Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" comment:"Cex buy quote amount"`
	common.ControlBy
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesInsertReq) Generate(model *models.StrategyDexCexTriangularArbitrageOpportunities) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.InstanceId = s.InstanceId
	model.OpportunityId = s.OpportunityId
	model.BuyOnDex = s.BuyOnDex
	model.DexTargetToken = s.DexTargetToken
	model.DexPoolType = s.DexPoolType
	model.DexPoolId = s.DexPoolId
	model.DexSolAmount = s.DexSolAmount
	model.DexTargetAmount = s.DexTargetAmount
	model.DexTxPriorityFee = s.DexTxPriorityFee
	model.DexTxJitoFee = s.DexTxJitoFee
	model.CexExchangeType = s.CexExchangeType
	model.CexTargetAsset = s.CexTargetAsset
	model.CexQuoteAsset = s.CexQuoteAsset
	model.CexSellQuantity = s.CexSellQuantity
	model.CexSellQuoteAmount = s.CexSellQuoteAmount
	model.CexBuyQuantity = s.CexBuyQuantity
	model.CexBuyQuoteAmount = s.CexBuyQuoteAmount
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesInsertReq) GetId() interface{} {
	return s.Id
}

type StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	InstanceId         string `json:"instanceId" comment:"实例id"`
	OpportunityId      string `json:"opportunityId" comment:"Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" comment:"买方标识"`
	DexTargetToken     string `json:"dexTargetToken" comment:"Dex target token pubkey"`
	DexPoolType        string `json:"dexPoolType" comment:"Dex pool type"`
	DexPoolId          string `json:"dexPoolId" comment:"Dex pool pubkey"`
	DexSolAmount       string `json:"dexSolAmount" comment:"Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" comment:"Dex trade target token amount"`
	DexTxPriorityFee   string `json:"dexTxPriorityFee" comment:"Dex trade tx priority fee"`
	DexTxJitoFee       string `json:"dexTxJitoFee" comment:"Dex trade tx jito fee"`
	CexExchangeType    string `json:"cexExchangeType" comment:"Cex exchange type"`
	CexTargetAsset     string `json:"cexTargetAsset" comment:"Cex target asset"`
	CexQuoteAsset      string `json:"cexQuoteAsset" comment:"Cex quote asset"`
	CexSellQuantity    string `json:"cexSellQuantity" comment:"Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" comment:"Cex sell quote amount"`
	CexBuyQuantity     string `json:"cexBuyQuantity" comment:"Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" comment:"Cex buy quote amount"`
	common.ControlBy
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq) Generate(model *models.StrategyDexCexTriangularArbitrageOpportunities) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.InstanceId = s.InstanceId
	model.OpportunityId = s.OpportunityId
	model.BuyOnDex = s.BuyOnDex
	model.DexTargetToken = s.DexTargetToken
	model.DexPoolType = s.DexPoolType
	model.DexPoolId = s.DexPoolId
	model.DexSolAmount = s.DexSolAmount
	model.DexTargetAmount = s.DexTargetAmount
	model.DexTxPriorityFee = s.DexTxPriorityFee
	model.DexTxJitoFee = s.DexTxJitoFee
	model.CexExchangeType = s.CexExchangeType
	model.CexTargetAsset = s.CexTargetAsset
	model.CexQuoteAsset = s.CexQuoteAsset
	model.CexSellQuantity = s.CexSellQuantity
	model.CexSellQuoteAmount = s.CexSellQuoteAmount
	model.CexBuyQuantity = s.CexBuyQuantity
	model.CexBuyQuoteAmount = s.CexBuyQuoteAmount
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesUpdateReq) GetId() interface{} {
	return s.Id
}

// StrategyDexCexTriangularArbitrageOpportunitiesGetReq 功能获取请求参数
type StrategyDexCexTriangularArbitrageOpportunitiesGetReq struct {
	Id int `uri:"id"`
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesGetReq) GetId() interface{} {
	return s.Id
}

// StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq 功能删除请求参数
type StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StrategyDexCexTriangularArbitrageOpportunitiesDeleteReq) GetId() interface{} {
	return s.Ids
}
