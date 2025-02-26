package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
	"time"
)

type StrategyDexCexTriangularArbitrageTradesGetPageReq struct {
	dto.Pagination `search:"-"`
	InstanceId     string `form:"instanceId"  search:"type:exact;column:instance_id;table:strategy_dex_cex_triangular_arbitrage_trades" comment:"Arbitrager instance ID"`
	BuyOnDex       string `form:"buyOnDex"  search:"type:exact;column:buy_on_dex;table:strategy_dex_cex_triangular_arbitrage_trades" comment:"Buy on dex or cex"`
	Symbol         string `form:"symbol" search:"-"`
	IsSuccess      bool   `form:"isSuccess" search:"-"`
	MinProfit      string `form:"minProfit" search:"-"`
	MaxProfit      string `form:"maxProfit" search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:updated_at;table:strategy_dex_cex_triangular_arbitrage_trades"` // >= BeginTime
	EndTime        string `form:"endTime" search:"type:lte;column:updated_at;table:strategy_dex_cex_triangular_arbitrage_trades"`   // <= EndTime

	StrategyDexCexTriangularArbitrageTradesOrder
}

type StrategyDexCexTriangularArbitrageTradesOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:strategy_dex_cex_triangular_arbitrage_trades"`
	InstanceId         string `form:"instanceIdOrder"  search:"type:order;column:instance_id;table:strategy_dex_cex_triangular_arbitrage_trades"`
	OpportunityId      string `form:"opportunityIdOrder"  search:"type:order;column:opportunity_id;table:strategy_dex_cex_triangular_arbitrage_trades"`
	BuyOnDex           string `form:"buyOnDexOrder"  search:"type:order;column:buy_on_dex;table:strategy_dex_cex_triangular_arbitrage_trades"`
	Error              string `form:"errorOrder"  search:"type:order;column:error;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexTrader          string `form:"dexTraderOrder"  search:"type:order;column:dex_trader;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexSuccess         string `form:"dexSuccessOrder"  search:"type:order;column:dex_success;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexTxFee           string `form:"dexTxFeeOrder"  search:"type:order;column:dex_tx_fee;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexTxSig           string `form:"dexTxSigOrder"  search:"type:order;column:dex_tx_sig;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexSolAmount       string `form:"dexSolAmountOrder"  search:"type:order;column:dex_sol_amount;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DexTargetAmount    string `form:"dexTargetAmountOrder"  search:"type:order;column:dex_target_amount;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexAmberAccount    string `form:"cexAmberAccountOrder"  search:"type:order;column:cex_amber_account;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexExchangeType    string `form:"cexExchangeTypeOrder"  search:"type:order;column:cex_exchange_type;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellSuccess     string `form:"cexSellSuccessOrder"  search:"type:order;column:cex_sell_success;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellOrderId     string `form:"cexSellOrderIdOrder"  search:"type:order;column:cex_sell_order_id;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellQuantity    string `form:"cexSellQuantityOrder"  search:"type:order;column:cex_sell_quantity;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellQuoteAmount string `form:"cexSellQuoteAmountOrder"  search:"type:order;column:cex_sell_quote_amount;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellFeeAsset    string `form:"cexSellFeeAssetOrder"  search:"type:order;column:cex_sell_fee_asset;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexSellFee         string `form:"cexSellFeeOrder"  search:"type:order;column:cex_sell_fee;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuySuccess      string `form:"cexBuySuccessOrder"  search:"type:order;column:cex_buy_success;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuyOrderId      string `form:"cexBuyOrderIdOrder"  search:"type:order;column:cex_buy_order_id;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuyQuantity     string `form:"cexBuyQuantityOrder"  search:"type:order;column:cex_buy_quantity;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuyQuoteAmount  string `form:"cexBuyQuoteAmountOrder"  search:"type:order;column:cex_buy_quote_amount;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuyFeeAsset     string `form:"cexBuyFeeAssetOrder"  search:"type:order;column:cex_buy_fee_asset;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CexBuyFee          string `form:"cexBuyFeeOrder"  search:"type:order;column:cex_buy_fee;table:strategy_dex_cex_triangular_arbitrage_trades"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:strategy_dex_cex_triangular_arbitrage_trades"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:strategy_dex_cex_triangular_arbitrage_trades"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:strategy_dex_cex_triangular_arbitrage_trades"`
}

func (m *StrategyDexCexTriangularArbitrageTradesGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StrategyDexCexTriangularArbitrageTradesGetPageResp struct {
	Id                 int       `json:"id" gorm:"column:id"`
	InstanceId         string    `json:"instanceId"`
	Symbol             string    `json:"symbol"` //需要从套利机会表中join
	OpportunityId      string    `json:"opportunityId"`
	BuyOnDex           string    `json:"buyOnDex"`
	Error              string    `json:"error"`
	DexTrader          string    `json:"dexTrader"`
	DexSuccess         string    `json:"dexSuccess"`
	DexTxFee           string    `json:"dexTxFee"`
	DexTxSig           string    `json:"dexTxSig"`
	DexSolAmount       string    `json:"dexSolAmount"`
	DexTargetAmount    string    `json:"dexTargetAmount"`
	CexAmberAccount    string    `json:"cexAmberAccount"`
	CexExchangeType    string    `json:"cexExchangeType"`
	CexSellSuccess     string    `json:"cexSellSuccess"`
	CexSellOrderId     string    `json:"cexSellOrderId"`
	CexSellQuantity    string    `json:"cexSellQuantity"`
	CexSellQuoteAmount string    `json:"cexSellQuoteAmount"`
	CexSellFeeAsset    string    `json:"cexSellFeeAsset"`
	CexSellFee         string    `json:"cexSellFee"`
	CexBuySuccess      string    `json:"cexBuySuccess"`
	CexBuyOrderId      string    `json:"cexBuyOrderId"`
	CexBuyQuantity     string    `json:"cexBuyQuantity"`
	CexBuyQuoteAmount  string    `json:"cexBuyQuoteAmount"`
	CexBuyFeeAsset     string    `json:"cexBuyFeeAsset"`
	CexBuyFee          string    `json:"cexBuyFee"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type StrategyDexCexTriangularArbitrageTradesGetDetailResp struct {
	Id                 int       `json:"id" gorm:"column:id"`
	InstanceId         string    `json:"instanceId"`
	OpportunityId      string    `json:"opportunityId"`
	BuyOnDex           string    `json:"buyOnDex"`
	Error              string    `json:"error"`
	DexTrader          string    `json:"dexTrader"`
	DexSuccess         string    `json:"dexSuccess"`
	DexTxFee           string    `json:"dexTxFee"`
	DexTxSig           string    `json:"dexTxSig"`
	DexSolAmount       string    `json:"dexSolAmount"`
	DexTargetAmount    string    `json:"dexTargetAmount"`
	CexAmberAccount    string    `json:"cexAmberAccount"`
	CexExchangeType    string    `json:"cexExchangeType"`
	CexSellSuccess     string    `json:"cexSellSuccess"`
	CexSellOrderId     string    `json:"cexSellOrderId"`
	CexSellQuantity    string    `json:"cexSellQuantity"`
	CexSellQuoteAmount string    `json:"cexSellQuoteAmount"`
	CexSellFeeAsset    string    `json:"cexSellFeeAsset"`
	CexSellFee         string    `json:"cexSellFee"`
	CexBuySuccess      string    `json:"cexBuySuccess"`
	CexBuyOrderId      string    `json:"cexBuyOrderId"`
	CexBuyQuantity     string    `json:"cexBuyQuantity"`
	CexBuyQuoteAmount  string    `json:"cexBuyQuoteAmount"`
	CexBuyFeeAsset     string    `json:"cexBuyFeeAsset"`
	CexBuyFee          string    `json:"cexBuyFee"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	// 交易机会相关字段
	DexPoolType            string `json:"dexPoolType" gorm:"-"`
	DexPoolId              string `json:"dexPoolId" gorm:"-"`
	DexTxPriorityFee       string `json:"dexTxPriorityFee" gorm:"-"`
	DexTxJitoFee           string `json:"dexTxJitoFee" gorm:"-"`
	CexTargetAsset         string `json:"cexTargetAsset" gorm:"-"`
	CexQuoteAsset          string `json:"cexQuoteAsset" gorm:"-"`
	DexTargetToken         string `json:"dexTargetToken" gorm:"-"`
	OppoDexSolAmount       string `json:"oppoDexSolAmount" gorm:"-"`
	OppoDexTargetAmount    string `json:"oppoDexTargetAmount" gorm:"-"`
	OppoCexSellQuantity    string `json:"oppoCexSellQuantity" gorm:"-"`
	OppoCexSellQuoteAmount string `json:"oppoCexSellQuoteAmount" gorm:"-"`
	OppoCexBuyQuantity     string `json:"oppoCexBuyQuantity" gorm:"-"`
	OppoCexBuyQuoteAmount  string `json:"oppoCexBuyQuoteAmount" gorm:"-"`
}

type StrategyDexCexTriangularArbitrageTradesInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	InstanceId         string `json:"instanceId" comment:"Arbitrager instance ID"`
	OpportunityId      string `json:"opportunityId" comment:"Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" comment:"Buy on dex or cex"`
	Error              string `json:"error" comment:"Error message"`
	DexTrader          string `json:"dexTrader" comment:"Dex trader pubkey"`
	DexSuccess         string `json:"dexSuccess" comment:"Is dex trade success"`
	DexTxFee           string `json:"dexTxFee" comment:"Dex trade tx fee"`
	DexTxSig           string `json:"dexTxSig" comment:"Dex trade tx signature"`
	DexSolAmount       string `json:"dexSolAmount" comment:"Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" comment:"Dex trade target amount"`
	CexAmberAccount    string `json:"cexAmberAccount" comment:"Cex amber account"`
	CexExchangeType    string `json:"cexExchangeType" comment:"Cex exchange type"`
	CexSellSuccess     string `json:"cexSellSuccess" comment:"Is cex sell success"`
	CexSellOrderId     string `json:"cexSellOrderId" comment:"Cex sell order id"`
	CexSellQuantity    string `json:"cexSellQuantity" comment:"Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" comment:"Cex sell quote amount"`
	CexSellFeeAsset    string `json:"cexSellFeeAsset" comment:"Cex sell fee asset"`
	CexSellFee         string `json:"cexSellFee" comment:"Cex sell fee"`
	CexBuySuccess      string `json:"cexBuySuccess" comment:"Is cex buy success"`
	CexBuyOrderId      string `json:"cexBuyOrderId" comment:"Cex buy order id"`
	CexBuyQuantity     string `json:"cexBuyQuantity" comment:"Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" comment:"Cex buy quote amount"`
	CexBuyFeeAsset     string `json:"cexBuyFeeAsset" comment:"Cex buy fee asset"`
	CexBuyFee          string `json:"cexBuyFee" comment:"Cex buy fee"`
	common.ControlBy
}

func (s *StrategyDexCexTriangularArbitrageTradesInsertReq) Generate(model *models.StrategyDexCexTriangularArbitrageTrades) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.InstanceId = s.InstanceId
	model.OpportunityId = s.OpportunityId
	model.BuyOnDex = s.BuyOnDex
	model.Error = s.Error
	model.DexTrader = s.DexTrader
	model.DexSuccess = s.DexSuccess
	model.DexTxFee = s.DexTxFee
	model.DexTxSig = s.DexTxSig
	model.DexSolAmount = s.DexSolAmount
	model.DexTargetAmount = s.DexTargetAmount
	model.CexAmberAccount = s.CexAmberAccount
	model.CexExchangeType = s.CexExchangeType
	model.CexSellSuccess = s.CexSellSuccess
	model.CexSellOrderId = s.CexSellOrderId
	model.CexSellQuantity = s.CexSellQuantity
	model.CexSellQuoteAmount = s.CexSellQuoteAmount
	model.CexSellFeeAsset = s.CexSellFeeAsset
	model.CexSellFee = s.CexSellFee
	model.CexBuySuccess = s.CexBuySuccess
	model.CexBuyOrderId = s.CexBuyOrderId
	model.CexBuyQuantity = s.CexBuyQuantity
	model.CexBuyQuoteAmount = s.CexBuyQuoteAmount
	model.CexBuyFeeAsset = s.CexBuyFeeAsset
	model.CexBuyFee = s.CexBuyFee
}

func (s *StrategyDexCexTriangularArbitrageTradesInsertReq) GetId() interface{} {
	return s.Id
}

type StrategyDexCexTriangularArbitrageTradesUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	InstanceId         string `json:"instanceId" comment:"Arbitrager instance ID"`
	OpportunityId      string `json:"opportunityId" comment:"Opportunity ID"`
	BuyOnDex           string `json:"buyOnDex" comment:"Buy on dex or cex"`
	Error              string `json:"error" comment:"Error message"`
	DexTrader          string `json:"dexTrader" comment:"Dex trader pubkey"`
	DexSuccess         string `json:"dexSuccess" comment:"Is dex trade success"`
	DexTxFee           string `json:"dexTxFee" comment:"Dex trade tx fee"`
	DexTxSig           string `json:"dexTxSig" comment:"Dex trade tx signature"`
	DexSolAmount       string `json:"dexSolAmount" comment:"Dex trade SOL amount"`
	DexTargetAmount    string `json:"dexTargetAmount" comment:"Dex trade target amount"`
	CexAmberAccount    string `json:"cexAmberAccount" comment:"Cex amber account"`
	CexExchangeType    string `json:"cexExchangeType" comment:"Cex exchange type"`
	CexSellSuccess     string `json:"cexSellSuccess" comment:"Is cex sell success"`
	CexSellOrderId     string `json:"cexSellOrderId" comment:"Cex sell order id"`
	CexSellQuantity    string `json:"cexSellQuantity" comment:"Cex sell quantity"`
	CexSellQuoteAmount string `json:"cexSellQuoteAmount" comment:"Cex sell quote amount"`
	CexSellFeeAsset    string `json:"cexSellFeeAsset" comment:"Cex sell fee asset"`
	CexSellFee         string `json:"cexSellFee" comment:"Cex sell fee"`
	CexBuySuccess      string `json:"cexBuySuccess" comment:"Is cex buy success"`
	CexBuyOrderId      string `json:"cexBuyOrderId" comment:"Cex buy order id"`
	CexBuyQuantity     string `json:"cexBuyQuantity" comment:"Cex buy quantity"`
	CexBuyQuoteAmount  string `json:"cexBuyQuoteAmount" comment:"Cex buy quote amount"`
	CexBuyFeeAsset     string `json:"cexBuyFeeAsset" comment:"Cex buy fee asset"`
	CexBuyFee          string `json:"cexBuyFee" comment:"Cex buy fee"`
	common.ControlBy
}

func (s *StrategyDexCexTriangularArbitrageTradesUpdateReq) Generate(model *models.StrategyDexCexTriangularArbitrageTrades) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.InstanceId = s.InstanceId
	model.OpportunityId = s.OpportunityId
	model.BuyOnDex = s.BuyOnDex
	model.Error = s.Error
	model.DexTrader = s.DexTrader
	model.DexSuccess = s.DexSuccess
	model.DexTxFee = s.DexTxFee
	model.DexTxSig = s.DexTxSig
	model.DexSolAmount = s.DexSolAmount
	model.DexTargetAmount = s.DexTargetAmount
	model.CexAmberAccount = s.CexAmberAccount
	model.CexExchangeType = s.CexExchangeType
	model.CexSellSuccess = s.CexSellSuccess
	model.CexSellOrderId = s.CexSellOrderId
	model.CexSellQuantity = s.CexSellQuantity
	model.CexSellQuoteAmount = s.CexSellQuoteAmount
	model.CexSellFeeAsset = s.CexSellFeeAsset
	model.CexSellFee = s.CexSellFee
	model.CexBuySuccess = s.CexBuySuccess
	model.CexBuyOrderId = s.CexBuyOrderId
	model.CexBuyQuantity = s.CexBuyQuantity
	model.CexBuyQuoteAmount = s.CexBuyQuoteAmount
	model.CexBuyFeeAsset = s.CexBuyFeeAsset
	model.CexBuyFee = s.CexBuyFee
}

func (s *StrategyDexCexTriangularArbitrageTradesUpdateReq) GetId() interface{} {
	return s.Id
}

// StrategyDexCexTriangularArbitrageTradesGetReq 功能获取请求参数
type StrategyDexCexTriangularArbitrageTradesGetReq struct {
	Id int `uri:"id"`
}

func (s *StrategyDexCexTriangularArbitrageTradesGetReq) GetId() interface{} {
	return s.Id
}

// StrategyDexCexTriangularArbitrageTradesDeleteReq 功能删除请求参数
type StrategyDexCexTriangularArbitrageTradesDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StrategyDexCexTriangularArbitrageTradesDeleteReq) GetId() interface{} {
	return s.Ids
}

// StrategyDexCexTriangularArbitrageTradesGetStatisticsReq 功能获取请求参数
type StrategyDexCexTriangularArbitrageTradesGetStatisticsReq struct {
	TradeResultType int `json:"TradeResultType"` //查询成功还是失败的统计,0-失败，1-成功，2-全部
}

type StrategyDexCexTriangularArbitrageTradesGetStatisticsResp struct {
	// 次数统计
	TotalTrade               int64  `json:"totalTrade" comment:"总套利次数"`
	TotalSuccessTrade        int64  `json:"totalSuccessTrade" comment:"总套利成功次数"`
	TotalFailedTrade         int64  `json:"totalFailedTrade" comment:"总套利失败次数"`
	DailyTotalTrade          int64  `json:"dailyTotalTrade" comment:"24小时套利次数"`
	DailyTotalSuccessTrade   int64  `json:"dailyTotalSuccessTrade" comment:"24小时套利成功次数"`
	DailyTotalFailedTrade    int64  `json:"dailyTotalFailedTrade" comment:"24小时套利失败次数"`
	TotalProfit              string `json:"totalProfit" comment:"总套利利润"`
	DailyTotalProfit         string `json:"dailyTotalProfit" comment:"24小时套利利润"`
	DailyProfitChangePercent string `json:"dailyProfitChangePercent" comment:"24小时套利利润变化"`
	TotalTradeVolume         string `json:"totalTradeVolume" comment:"总套利交易量"`
	DailyTotalTradeVolume    string `json:"dailyTotalTradeVolume" comment:"24小时套利交易量"`
}
