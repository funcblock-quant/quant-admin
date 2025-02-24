package dto

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"google.golang.org/protobuf/proto"
	"quanta-admin/app/business/models"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
	"strconv"
)

type BusDexCexTriangularObserverGetPageReq struct {
	dto.Pagination `search:"-"`
	Symbol         string `form:"symbol"  search:"type:exact;column:symbol;table:bus_dex_cex_triangular_instance"`
	BusDexCexTriangularObserverOrder
}

type BusDexCexTriangularObserverOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_triangular_instance"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_dex_cex_triangular_instance"`
	ObserverId         string `form:"observerIdOrder"  search:"type:order;column:instance_id;table:bus_dex_cex_triangular_instance"`
	Symbol             string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_triangular_instance"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:bus_dex_cex_triangular_instance"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_dex_cex_triangular_instance"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_dex_cex_triangular_instance"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_triangular_instance"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_dex_cex_triangular_instance"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_dex_cex_triangular_instance""`
}

func (m *BusDexCexTriangularObserverGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexTriangularObserverGetPageResp struct {
	models.BusDexCexTriangularObserver
	ProfitOfBuyOnDexResp
	ProfitOfSellOnDexResp
}

type ProfitOfBuyOnDexResp struct {
	CexSellPrice       float64 `json:"cexSellPrice" gorm:"-"`
	DexBuyPrice        float64 `json:"dexBuyPrice" gorm:"-"`
	DexBuyDiffPrice    float64 `json:"dexBuyDiffPrice" gorm:"-"`
	DexBuyDiffPercent  string  `json:"dexBuyDiffPercent" gorm:"-"`
	DexBuyDiffDuration string  `json:"dexBuyDiffDuration" gorm:"-"`
	ProfitOfBuyOnDex   float64 `json:"profitOfBuyOnDex" gorm:"-"`
}

type ProfitOfSellOnDexResp struct {
	CexBuyPrice         float64 `json:"cexBuyPrice" gorm:"-"`
	DexSellPrice        float64 `json:"dexSellPrice" gorm:"-"`
	DexSellDiffPrice    float64 `json:"dexSellDiffPrice" gorm:"-"`
	DexSellDiffPercent  string  `json:"dexSellDiffPercent" gorm:"-"`
	DexSellDiffDuration string  `json:"dexSellDiffDuration" gorm:"-"`
	ProfitOfSellOnDex   float64 `json:"profitOfSellOnDex" gorm:"-"`
}

type BusDexCexTriangularObserverDetailResp struct {
	models.BusDexCexTriangularObserver
	ProfitOfBuyOnDexResp  //最新dex buy 价差
	ProfitOfSellOnDexResp //最新dex sell 价差
}

type BusDexCexTriangularSpreadHistory struct {
	CexSellPriceChartPoints       []PriceChartPoint `json:"cexSellPriceChartPoints" gorm:"-"`
	DexBuyPriceChartPoints        []PriceChartPoint `json:"dexBuyPriceChartPoints" gorm:"-"`
	DexBuyPriceSpreadChartPoints  []PriceChartPoint `json:"dexBuyPriceSpreadChartPoints" gorm:"-"`
	DexBuyProfitChartPoints       []PriceChartPoint `json:"dexBuyProfitChartPoints" gorm:"-"`
	DexSellPriceChartPoints       []PriceChartPoint `json:"dexSellPriceChartPoints" gorm:"-"`
	CexBuyPriceChartPoints        []PriceChartPoint `json:"cexBuyPriceChartPoints" gorm:"-"`
	DexSellPriceSpreadChartPoints []PriceChartPoint `json:"dexSellPriceSpreadChartPoints" gorm:"-"`
	DexSellProfitChartPoints      []PriceChartPoint `json:"dexSellProfitChartPoints" gorm:"-"`
}

type PriceChartPoint struct {
	XAxis int64   `json:"xAxis"` //横坐标，时间戳
	YAxis float64 `json:"yAxis"` //纵坐标
}

type BusDexCexTriangularObserverInsertReq struct {
	Id                       int      `json:"-" comment:""`
	StrategyInstanceId       string   `json:"strategyInstanceId" comment:"策略id"`
	InstanceId               string   `json:"instanceId" comment:"观察器id"`
	Symbol                   string   `json:"symbol" comment:"观察币种"`
	TargetToken              string   `json:"targetToken"`
	QuoteToken               string   `json:"quoteToken"`
	SymbolConnector          string   `json:"-"`
	ExchangeType             string   `json:"exchangeType"`
	DexType                  string   `json:"dexType"`
	MinSolAmount             *float32 `json:"minSolAmount"`
	MaxSolAmount             *float32 `json:"maxSolAmount"`
	TakerFee                 *float32 `json:"takerFee"`
	TokenMint                *string  `json:"tokenMint"`
	OwnerProgram             *string  `json:"ownerProgram"`
	Decimals                 int      `json:"decimals"`
	AmmPoolId                *string  `json:"ammPool"`
	TriggerProfitQuoteAmount *float64 `json:"triggerProfitQuoteAmount"`
	PriorityFee              *float64 `json:"priorityFee"`
	JitoFee                  *float64 `json:"jitoFee"`
	Depth                    string   `json:"depth"`
	Status                   string   `json:"status" comment:"状态"`
	common.ControlBy
}

type BusDexCexTriangularObserverBatchInsertReq struct {
	StrategyInstanceId       string   `json:"strategyInstanceId" comment:"策略id"`
	InstanceId               string   `json:"instanceId" comment:"策略端实例id"`
	TargetToken              []string `json:"targetToken"`
	QuoteToken               string   `json:"quoteToken"`
	SymbolConnector          string   `json:"-"`
	ExchangeType             string   `json:"exchangeType"`
	DexType                  string   `json:"dexType"`
	MinSolAmount             *float64 `json:"minSolAmount"`
	MaxSolAmount             *float64 `json:"maxSolAmount"`
	TakerFee                 *float64 `json:"takerFee"`
	AmmPoolId                *string  `json:"ammPool"`
	TokenMint                *string  `json:"tokenMint"`
	OwnerProgram             *string  `json:"ownerProgram"`
	Decimals                 int      `json:"decimals"`
	MaxArraySize             int      `json:"maxArraySize"`
	TriggerProfitQuoteAmount *float64 `json:"triggerProfitQuoteAmount"`
	PriorityFee              *float64 `json:"priorityFee"`
	JitoFee                  *float64 `json:"jitoFee"`
	Depth                    string   `json:"depth"`
	Status                   string   `json:"status" comment:"状态"`

	common.ControlBy
}

func (s *BusDexCexTriangularObserverBatchInsertReq) Generate(model *models.BusDexCexTriangularObserver, baseToken string) {
	model.StrategyInstanceId = "1" //default 1
	model.InstanceId = utils.GetUUID()
	model.Symbol = baseToken + "/" + s.QuoteToken
	model.TargetToken = s.TargetToken[0] //目前就支持单个添加
	model.QuoteToken = s.QuoteToken
	model.SymbolConnector = "/" //之前沟通默认amber全部是/连接的
	model.ExchangeType = s.ExchangeType
	model.DexType = s.DexType
	model.MaxArraySize = s.MaxArraySize
	model.MinSolAmount = s.MinSolAmount
	model.MaxSolAmount = s.MaxSolAmount
	model.TakerFee = s.TakerFee
	model.TokenMint = s.TokenMint
	model.OwnerProgram = s.OwnerProgram
	model.MinProfit = s.TriggerProfitQuoteAmount
	scaled := *s.PriorityFee * 1_000_000_000
	model.PriorityFee = uint64(scaled)
	scaled = *s.JitoFee * 1_000_000_000
	model.JitoFee = uint64(scaled)
	model.Decimals = s.Decimals
	model.AmmPoolId = s.AmmPoolId
	model.Depth = s.Depth
	model.Status = "0"          //新增的话说明已经启动成功了
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateAmmConfig(ammConfig *pb.DexConfig) error {
	//slippageBpsUint, err := strconv.ParseUint(*s.SlippageBps, 10, 32)
	//if err != nil {
	//	log.Errorf("slippageBps: %v\n", slippageBpsUint)
	//	return errors.New("error slippageBps")
	//}
	//log.Infof("slippageBps: %v\n", slippageBpsUint)
	maxArraySize := new(uint32)
	*maxArraySize = uint32(s.MaxArraySize) //默认5， clmm使用参数
	if s.DexType == "RAY_AMM" {
		ammConfig.Config = &pb.DexConfig_RayAmm{
			RayAmm: &pb.RayAmmConfig{
				Pool:      s.AmmPoolId,
				TokenMint: s.TokenMint,
			},
		}
	} else if s.DexType == "RAY_CLMM" {
		ammConfig.Config = &pb.DexConfig_RayClmm{
			RayClmm: &pb.RayClmmConfig{
				Pool:         s.AmmPoolId,
				TokenMint:    s.TokenMint,
				MaxArraySize: maxArraySize,
			},
		}
	}
	fmt.Printf("ammConfig: %v\n", *ammConfig)
	return nil
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateAmberConfig(amberConfig *pb.AmberObserverConfig) error {
	amberConfig.ExchangeType = &s.ExchangeType
	amberConfig.TakerFee = proto.Float64(*s.TakerFee)

	amberConfig.TargetToken = &s.TargetToken[0]
	amberConfig.QuoteToken = &s.QuoteToken

	amberConfig.BidDepth = proto.Int32(20)
	amberConfig.AskDepth = proto.Int32(20)

	if s.Depth != "" {
		depthInt, err := strconv.Atoi(s.Depth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		amberConfig.BidDepth = proto.Int32(int32(depthInt))
		amberConfig.AskDepth = proto.Int32(int32(depthInt))
	}
	return nil
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateObserverParams(observerParams *pb.ObserverParams) error {
	observerParams.MinSolAmount = proto.Float64(*s.MinSolAmount)
	observerParams.MaxSolAmount = proto.Float64(*s.MaxSolAmount)
	priorityFee := uint64(*s.PriorityFee * 1_000_000_000)
	jitoFee := uint64(*s.JitoFee * 1_000_000_000)
	transactionFee := &pb.TransactionFee{
		PriorityFee: &priorityFee,
		JitoFee:     &jitoFee,
	}
	observerParams.TxFee = transactionFee
	observerParams.TriggerProfitQuoteAmount = proto.Float64(*s.TriggerProfitQuoteAmount)

	return nil
}

func (s *BusDexCexTriangularObserverInsertReq) Generate(model *models.BusDexCexTriangularObserver) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.InstanceId = s.InstanceId
	model.Symbol = s.Symbol
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusDexCexTriangularObserverInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexTriangularObserverUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略id"`
	ObserverId         string `json:"observerId" comment:"观察器id"`
	Symbol             string `json:"symbol" comment:"观察币种"`
	Status             string `json:"status" comment:"状态"`
	common.ControlBy
}

func (s *BusDexCexTriangularObserverUpdateReq) Generate(model *models.BusDexCexTriangularObserver) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.InstanceId = s.ObserverId
	model.Symbol = s.Symbol
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusDexCexTriangularObserverUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexTriangularObserverGetReq 功能获取请求参数
type BusDexCexTriangularObserverGetReq struct {
	Id int `uri:"id"`
}

func (s *BusDexCexTriangularObserverGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexTriangularObserverDeleteReq 功能删除请求参数
type BusDexCexTriangularObserverDeleteReq struct {
	Ids int `json:"ids"`
}

func (s *BusDexCexTriangularObserverDeleteReq) GetId() interface{} {
	return s.Ids
}

type DexCexTriangularObserverSymbolListResp struct {
	Symbol string `json:"symbol" gorm:"column:symbol"`
}

type BusDexCexTriangularObserverStartTraderReq struct {
	InstanceId             int      `json:"instanceId" comment:"策略端实例id"`
	AlertThreshold         *float64 `json:"alertThreshold"`
	BuyTriggerThreshold    *float64 `json:"buyTriggerThreshold"`
	TargetBalanceThreshold *float64 `json:"targetBalanceThreshold"`
	SellTriggerThreshold   *float64 `json:"sellTriggerThreshold"`
	SlippageBps            *string  `json:"slippage"`
	//MinProfit              *float64 `json:"minProfit"`
	//PriorityFee            *float64 `json:"priorityFee"`
	//JitoFee                *float64 `json:"jitoFee"`
	common.ControlBy
}

type BusDexCexTriangularObserverStopTraderReq struct {
	InstanceId int `json:"instanceId" comment:"策略端实例id"`
	common.ControlBy
}

type BusDexCexTriangularUpdateObserverParamsReq struct {
	InstanceId               int      `json:"instanceId" comment:"策略端实例id"`
	MinSolAmount             *float64 `json:"minSolAmount"`
	MaxSolAmount             *float64 `json:"maxSolAmount"`
	TriggerProfitQuoteAmount *float64 `json:"triggerProfitQuoteAmount"`
	PriorityFee              *float64 `json:"priorityFee"`
	JitoFee                  *float64 `json:"jitoFee"`
	common.ControlBy
}

func (s *BusDexCexTriangularUpdateObserverParamsReq) Generate(model *models.BusDexCexTriangularObserver) {
	model.InstanceId = strconv.Itoa(s.InstanceId)
	model.MinProfit = s.TriggerProfitQuoteAmount
	scaled := *s.PriorityFee * 1_000_000_000
	model.PriorityFee = uint64(scaled)
	scaled = *s.JitoFee * 1_000_000_000
	model.JitoFee = uint64(scaled)
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

type BusDexCexTriangularUpdateTraderParamsReq struct {
	InstanceId  int     `json:"instanceId" comment:"策略端实例id"`
	SlippageBps *string `json:"slippage"`
	//MinProfit   *float64 `json:"minProfit"`
	//PriorityFee *float64 `json:"priorityFee"`
	//JitoFee     *float64 `json:"jitoFee"`
	common.ControlBy
}
