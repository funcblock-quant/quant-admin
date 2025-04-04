package dto

import (
	"quanta-admin/app/business/models"
	pb "quanta-admin/app/grpc/proto/client/observer_service"
	"quanta-admin/common/dto"
	"quanta-admin/common/global"
	common "quanta-admin/common/models"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"google.golang.org/protobuf/proto"
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
	Id                 int      `json:"-" comment:""`
	StrategyInstanceId string   `json:"strategyInstanceId" comment:"策略id"`
	InstanceId         string   `json:"instanceId" comment:"观察器id"`
	Symbol             string   `json:"symbol" comment:"观察币种"`
	TargetToken        string   `json:"targetToken"`
	QuoteToken         string   `json:"quoteToken"`
	SymbolConnector    string   `json:"-"`
	ExchangeType       string   `json:"exchangeType"`
	DexType            string   `json:"dexType"`
	MinQuoteAmount     *float64 `json:"minQuoteAmount"`
	MaxQuoteAmount     *float64 `json:"maxQuoteAmount"`
	TakerFee           *float32 `json:"takerFee"`
	TokenMint          *string  `json:"tokenMint"`
	OwnerProgram       *string  `json:"ownerProgram"`
	Decimals           int      `json:"decimals"`
	AmmPoolId          *string  `json:"ammPool"`
	ProfitTriggerRate  *float64 `json:"profitTriggerRate"`
	TriggerHoldingMs   int      `json:"triggerHoldingMs"`
	//SlippageBpsRate    *string  `json:"slippageRate"`
	Depth  string `json:"depth"`
	Status string `json:"status" comment:"状态"`
	common.ControlBy
}

type BusDexCexTriangularObserverBatchInsertReq struct {
	StrategyInstanceId string   `json:"strategyInstanceId" comment:"策略id"`
	InstanceId         string   `json:"instanceId" comment:"策略端实例id"`
	TargetToken        []string `json:"targetToken"`
	QuoteToken         string   `json:"quoteToken"`
	SymbolConnector    string   `json:"-"`
	ExchangeType       string   `json:"exchangeType"`
	DexType            string   `json:"dexType"`
	MinQuoteAmount     *float64 `json:"minQuoteAmount"`
	MaxQuoteAmount     *float64 `json:"maxQuoteAmount"`
	TakerFee           *float64 `json:"takerFee"`
	AmmPoolId          *string  `json:"ammPool"`
	TokenMint          *string  `json:"tokenMint"`
	OwnerProgram       *string  `json:"ownerProgram"`
	Decimals           int      `json:"decimals"`
	MaxArraySize       int      `json:"maxArraySize"`
	ProfitTriggerRate  *float64 `json:"profitTriggerRate"`
	TriggerHoldingMs   int      `json:"triggerHoldingMs"`
	//SlippageBpsRate    *float64 `json:"slippageBpsRate"`
	Depth  string `json:"depth"`
	Status string `json:"status" comment:"状态"`

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
	model.MinQuoteAmount = s.MinQuoteAmount
	model.MaxQuoteAmount = s.MaxQuoteAmount
	model.TakerFee = s.TakerFee
	model.TokenMint = s.TokenMint
	model.OwnerProgram = s.OwnerProgram
	model.ProfitTriggerRate = s.ProfitTriggerRate // 比例
	//model.TriggerHoldingMs = s.TriggerHoldingMs
	//model.SlippageBpsRate = s.SlippageBpsRate
	model.Decimals = s.Decimals
	model.AmmPoolId = s.AmmPoolId
	model.Depth = s.Depth
	model.Status = "0"          //新增的话说明已经启动成功了
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateAmmConfig(ammConfig *pb.DexConfig) error {

	maxArraySize := new(uint32)
	*maxArraySize = uint32(s.MaxArraySize) //默认5， clmm使用参数
	if s.DexType == global.DEX_TYPE_RAY_AMM {
		ammConfig.Config = &pb.DexConfig_RayAmm{
			RayAmm: &pb.RayAmmConfig{
				Pool:      s.AmmPoolId,
				TokenMint: s.TokenMint,
			},
		}
	} else if s.DexType == global.DEX_TYPE_RAY_CLMM {
		ammConfig.Config = &pb.DexConfig_RayClmm{
			RayClmm: &pb.RayClmmConfig{
				Pool:         s.AmmPoolId,
				TokenMint:    s.TokenMint,
				MaxArraySize: maxArraySize,
			},
		}
	} else if s.DexType == global.DEX_TYPE_ORCA_WHIRL_POOL {
		ammConfig.Config = &pb.DexConfig_OrcaWhirlPool{
			OrcaWhirlPool: &pb.OrcaWhirlPoolConfig{
				Pool:      s.AmmPoolId,
				TokenMint: s.TokenMint,
			},
		}
	}
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
	observerParams.MinQuoteAmount = proto.Float64(*s.MinQuoteAmount)
	observerParams.MaxQuoteAmount = proto.Float64(*s.MaxQuoteAmount)

	//observerParams.SlippageRate = proto.Float64(*s.SlippageBpsRate)
	observerParams.ProfitTriggerRate = proto.Float64(*s.ProfitTriggerRate)
	//observerParams.TriggerHoldingMs = proto.Uint64(uint64(s.TriggerHoldingMs))
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

type DexCexTriangularObserverExchangeListResp struct {
	Exchange string `json:"exchange" gorm:"column:exchange"`
}

type BusDexCexTriangularObserverStartTraderReq struct {
	InstanceId                 int      `json:"id" comment:"策略端实例id"`
	AlertThreshold             *float64 `json:"alertThreshold"`
	BuyTriggerThreshold        *float64 `json:"buyTriggerThreshold"`
	SellTriggerThreshold       *float64 `json:"sellTriggerThreshold"`
	MinDepositAmountThreshold  *float64 `json:"minDepositAmountThreshold"`
	MinWithdrawAmountThreshold *float64 `json:"minWithdrawAmountThreshold"`
	SlippageBpsRate            *float64 `json:"slippageBpsRate"`
	PriorityFee                *float64 `json:"priorityFee"`
	JitoFeeRate                *float64 `json:"jitoFeeRate"`
	CexAccount                 int64    `json:"cexAccount"`
	DexWallet                  int64    `json:"DexWallet"`
	common.ControlBy
}

type BusDexCexTriangularObserverStopTraderReq struct {
	InstanceId int `json:"id" comment:"策略端实例id"`
	common.ControlBy
}

type BusDexCexTriangularUpdateObserverParamsReq struct {
	InstanceId        int      `json:"id" comment:"策略端实例id"`
	MinQuoteAmount    *float64 `json:"minQuoteAmount"`
	MaxQuoteAmount    *float64 `json:"maxQuoteAmount"`
	ProfitTriggerRate *float64 `json:"profitTriggerRate"`
	TriggerHoldingMs  int      `json:"triggerHoldingMs"`
	common.ControlBy
}

func (s *BusDexCexTriangularUpdateObserverParamsReq) Generate(model *models.BusDexCexTriangularObserver) {
	model.InstanceId = strconv.Itoa(s.InstanceId)
	model.ProfitTriggerRate = s.ProfitTriggerRate
	model.MinQuoteAmount = s.MinQuoteAmount
	model.MaxQuoteAmount = s.MaxQuoteAmount
	//model.TriggerHoldingMs = s.TriggerHoldingMs
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

type BusDexCexTriangularUpdateTraderParamsReq struct {
	InstanceId      int      `json:"id" comment:"策略端实例id"`
	SlippageBpsRate *float64 `json:"slippageBpsRate"`
	PriorityFee     *float64 `json:"priorityFee"`
	JitoFeeRate     *float64 `json:"jitoFeeRate"`
	common.ControlBy
}

type BusDexCexTriangularUpdateWaterLevelParamsReq struct {
	InstanceId                 int      `json:"id" comment:"策略端实例id"`
	AlertThreshold             *float64 `json:"alertThreshold"`
	BuyTriggerThreshold        *float64 `json:"buyTriggerThreshold"`
	SellTriggerThreshold       *float64 `json:"sellTriggerThreshold"`
	MinDepositAmountThreshold  *float64 `json:"minDepositAmountThreshold"`
	MinWithdrawAmountThreshold *float64 `json:"minWithdrawAmountThreshold"`
	common.ControlBy
}

type BusDexCexTriangularGlobalWaterLevelStateResp struct {
	SolWaterLevelState         bool                                          `json:"solWaterLevelState" comment:"sol水位调节状态"`
	SolWaterLevelConfig        *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"solWaterLevelConfig"`
	StableCoinWaterLevelState  bool                                          `json:"stableCoinWaterLevelState" comment:"稳定币水位调节状态"`
	StableCoinWaterLevelConfig *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"stableCoinWaterLevelConfig"`
}

type BusDexCexTriangularUpdateGlobalWaterLevelConfigReq struct {
	ExchangeType               string                                        `json:"exchangeType"`
	DexWalletId                int                                           `json:"dexWalletId" comment:"策略端实例id"`
	CexAccountId               int                                           `json:"cexAccountId"`
	SolWaterLevelConfig        *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"solWaterLevelConfig"`
	StableCoinWaterLevelConfig *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"stableCoinWaterLevelConfig"`
}

type BusDexCexTriangularUpdateGlobalRiskConfig struct {
	AbsoluteLossThreshold       []RiskControlItem `json:"absoluteLossThreshold"`       // 绝对亏损风控
	RelativeLossThreshold       []RiskControlItem `json:"relativeLossThreshold"`       // 相对亏损风控
	SymbolDailyMaxLossThreshold []RiskControlItem `json:"symbolDailyMaxLossThreshold"` // 币种当日亏损风控
}

type ActionDetail struct {
	PauseDuration int  `json:"pauseDuration"` // 暂停时长，-1 表示次日零点
	ManualResume  bool `json:"manualResume"`  // 是否人工恢复
}

type RiskControlItem struct {
	Threshold    float64      `json:"threshold"`    // 阈值
	Action       int          `json:"action"`       // 操作类型 (1-预警, 2-暂停当前实例交易, 3-暂停全局交易)
	ActionDetail ActionDetail `json:"actionDetail"` // 具体操作详情
}

// BusGetCexAccountListReq 获取请求参数
type BusGetCexAccountListReq struct {
	Exchange string `uri:"exchange"`
}

// BusGetCexExchangeListReq 获取请求参数
type BusGetCexExchangeConfigListReq struct {
	Exchange string `uri:"exchange"`
}

// 根据cex或dex账号，查询绑定的另一侧的账号列表
type BusGetBoundAccountReq struct {
	AccountType string `json:"accountType"` // Cex or Dex
	AccountId   int64  `json:"accountId"`
}

type BusGetBoundAccountResp struct {
	CexAccountList []models.BusExchangeAccountInfo `json:"cexAccountList"` // Cex or Dex
	DexWalletList  []models.BusDexWallet           `json:"dexWalletList"`
}

type BusAccountPairInfo struct {
	CexAccountId     int                                           `json:"cexAccountId"`
	CexAccountName   string                                        `json:"cexAccountName"`
	CexAccountUid    string                                        `json:"cexAccountUid"`
	DexwalletId      int                                           `json:"dexWalletId"`
	DexWalletName    string                                        `json:"dexWalletName"`
	DexWalletAddr    string                                        `json:"dexWalletAddr"`
	HasGlobalConfig  bool                                          `json:"hasGlobalConfig"` // 是否有全局水位配置
	SolanaConfig     *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"solanaConfig"`
	StableCoinConfig *BusDexCexTriangularUpdateWaterLevelParamsReq `json:"stableCoinConfig"`
}
