package dto

import (
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"google.golang.org/protobuf/proto"
	"quanta-admin/app/business/models"
	pb "quanta-admin/app/grpc/proto/stub"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
	"strconv"
)

type BusDexCexTriangularObserverGetPageReq struct {
	dto.Pagination `search:"-"`
	Symbol         string `form:"symbol"  search:"type:exact;column:symbol;table:bus_dex_cex_triangular_observer"`
	BusDexCexTriangularObserverOrder
}

type BusDexCexTriangularObserverOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_triangular_observer"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_dex_cex_triangular_observer"`
	ObserverId         string `form:"observerIdOrder"  search:"type:order;column:observer_id;table:bus_dex_cex_triangular_observer"`
	Symbol             string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_triangular_observer"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:bus_dex_cex_triangular_observer"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_dex_cex_triangular_observer"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_dex_cex_triangular_observer"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_triangular_observer"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_dex_cex_triangular_observer"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_dex_cex_triangular_observer"`
}

func (m *BusDexCexTriangularObserverGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexTriangularObserverGetPageResp struct {
	models.BusDexCexTriangularObserver
	BaseProfit  string `json:"baseProfit" gorm:"-"`
	QuoteProfit string `json:"quoteProfit" gorm:"-"`
}

type BusDexCexTriangularObserverInsertReq struct {
	Id                 int      `json:"-" comment:""`
	StrategyInstanceId string   `json:"strategyInstanceId" comment:"策略id"`
	ObserverId         string   `json:"observerId" comment:"观察器id"`
	Symbol             string   `json:"symbol" comment:"观察币种"`
	ExchangeType       string   `json:"exchange_type"`
	Volume             *float32 `json:"volume"`
	TakerFee           *float32 `json:"taker_fee"`
	BaseTokenMint      *string  `json:"base_token_mint"`
	QuoteTokenMint     *string  `json:"quote_token_mint"`
	AmmPoolId          *string  `json:"amm_pool_id"`
	SlippageBps        *float32 `json:"slippage_bps"`
	Depth              string   `json:"depth"`
	Status             string   `json:"status" comment:"状态"`
	common.ControlBy
}

type BusDexCexTriangularObserverBatchInsertReq struct {
	StrategyInstanceId string   `json:"strategyInstanceId" comment:"策略id"`
	ObserverId         string   `json:"observerId" comment:"观察器id"`
	Symbols            []string `json:"symbolsArray" comment:"观察币种"`
	ExchangeType       string   `json:"exchangeType"`
	Volume             *float64 `json:"volume"`
	TakerFee           *float64 `json:"takerFee"`
	AmmPoolId          *string  `json:"ammPool"`
	BaseTokenMint      *string  `json:"baseTokenMint"`
	BaseTokenDecimal   *string  `json:"baseTokenDecimals"`
	QuoteTokenMint     *string  `json:"quoteTokenMint"`
	QuoteTokenDecimal  *string  `json:"quoteTokenDecimals"`
	SlippageBps        *string  `json:"slippage"`
	Depth              string   `json:"depth"`
	Status             string   `json:"status" comment:"状态"`
	common.ControlBy
}

func (s *BusDexCexTriangularObserverBatchInsertReq) Generate(model *models.BusDexCexTriangularObserver, symbol string, observerId string) {
	model.StrategyInstanceId = "1" //default 1
	model.ObserverId = observerId
	model.Symbol = symbol
	model.Volume = s.Volume
	model.TakerFee = s.TakerFee
	model.BaseTokenMint = s.BaseTokenMint
	model.QuoteTokenMint = s.QuoteTokenMint
	model.AmmPoolId = s.AmmPoolId
	model.SlippageBps = *s.SlippageBps
	model.Depth = s.Depth
	model.Status = "1"          //新增的话说明已经启动成功了
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateAmmConfig(ammConfig *pb.AmmDexConfig) error {
	ammConfig.AmmPool = s.AmmPoolId
	ammConfig.BaseTokenMint = s.BaseTokenMint
	baseTokenDecimalInt, err := strconv.ParseUint(*s.BaseTokenDecimal, 10, 32)
	if err != nil {
		return errors.New("error base_token_decimal")
	}
	ammConfig.BaseTokenDecimals = proto.Uint32(uint32(baseTokenDecimalInt))
	ammConfig.QuoteTokenMint = s.QuoteTokenMint
	quoteTokenDecimalInt, err := strconv.ParseUint(*s.QuoteTokenDecimal, 10, 32)
	if err != nil {
		return errors.New("error quote_token_decimal")
	}
	ammConfig.QuoteTokenDecimals = proto.Uint32(uint32(quoteTokenDecimalInt))
	slippageBpsUint, err := strconv.ParseUint(*s.SlippageBps, 10, 32)
	if err != nil {
		return errors.New("error slippageBps")
	}
	log.Infof("slippageBps: %v\n", slippageBpsUint)
	ammConfig.SlippageBps = proto.Uint64(slippageBpsUint)
	return nil
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateAmberConfig(amberConfig *pb.AmberConfig, symbol string) error {
	amberConfig.ExchangeType = &s.ExchangeType
	amberConfig.TakerFee = proto.Float64(*s.TakerFee)

	baseTokenOrderBook := pb.AmberOrderBookConfig{}
	baseTokenOrderBook.Symbol = &symbol
	baseTokenOrderBook.BidDepth = proto.Int32(20)
	baseTokenOrderBook.AskDepth = proto.Int32(20)

	quoteTokenOrderBook := pb.AmberOrderBookConfig{}
	quotoToken := "SOL/USDT"
	quoteTokenOrderBook.Symbol = &quotoToken
	quoteTokenOrderBook.BidDepth = proto.Int32(20)
	quoteTokenOrderBook.AskDepth = proto.Int32(20)

	if s.Depth != "" {
		depthInt, err := strconv.Atoi(s.Depth)
		if err != nil {
			depthInt = 20 //默认20档
		}
		baseTokenOrderBook.BidDepth = proto.Int32(int32(depthInt))
		baseTokenOrderBook.AskDepth = proto.Int32(int32(depthInt))

		quoteTokenOrderBook.BidDepth = proto.Int32(int32(depthInt))
		quoteTokenOrderBook.AskDepth = proto.Int32(int32(depthInt))
	}
	amberConfig.BaseTokenOrderbook = &baseTokenOrderBook
	amberConfig.QuoteTokenOrderbook = &quoteTokenOrderBook
	return nil
}

func (s *BusDexCexTriangularObserverBatchInsertReq) GenerateArbitrageConfig(arbitrageConfig *pb.ArbitrageConfig) error {
	arbitrageConfig.Volumn = proto.Float64(*s.Volume)
	return nil
}

func (s *BusDexCexTriangularObserverInsertReq) Generate(model *models.BusDexCexTriangularObserver) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.ObserverId = s.ObserverId
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
	model.ObserverId = s.ObserverId
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
	Ids        int    `json:"ids"`
	ObserverId string `json:"observerId"`
}

func (s *BusDexCexTriangularObserverDeleteReq) GetId() interface{} {
	return s.Ids
}

type DexCexTriangularObserverSymbolListResp struct {
	Symbol string `json:"symbol" gorm:"column:symbol"`
}
