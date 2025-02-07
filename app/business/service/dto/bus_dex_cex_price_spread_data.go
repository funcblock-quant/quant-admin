package dto

import (
	"time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexCexPriceSpreadDataGetPageReq struct {
	dto.Pagination `search:"-"`
	Symbol         string `form:"symbol"  search:"type:exact;column:symbol;table:bus_dex_cex_price_spread_data" comment:"观察币种"`
	BusDexCexPriceSpreadDataOrder
}

type BusDexCexPriceSpreadDataOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_price_spread_data"`
	ObserverId   string `form:"observerIdOrder"  search:"type:order;column:observer_id;table:bus_dex_cex_price_spread_data"`
	Symbol       string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_price_spread_data"`
	DexBuyPrice  string `form:"dexBuyPriceOrder"  search:"type:order;column:dex_buy_price;table:bus_dex_cex_price_spread_data"`
	DexSellPrice string `form:"dexSellPriceOrder"  search:"type:order;column:dex_sell_price;table:bus_dex_cex_price_spread_data"`
	CexBuyPrice  string `form:"cexBuyPriceOrder"  search:"type:order;column:cex_buy_price;table:bus_dex_cex_price_spread_data"`
	CexSellPrice string `form:"cexSellPriceOrder"  search:"type:order;column:cex_sell_price;table:bus_dex_cex_price_spread_data"`
	SnapshotTime string `form:"snapshotTimeOrder"  search:"type:order;column:snapshot_time;table:bus_dex_cex_price_spread_data"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_price_spread_data"`
}

func (m *BusDexCexPriceSpreadDataGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexPriceSpreadDataHistoryChartReq struct {
	ObserverId string `json:"observerId" comment:"监视器id"`
	Interval   string `json:"interval" comment:"数据时间间隔,单位：s"`
}

type BusDexCexPriceSpreadDataInsertReq struct {
	Id           int       `json:"-" comment:""` //
	ObserverId   string    `json:"observerId" comment:"观察器id"`
	Symbol       string    `json:"symbol" comment:"观察币种"`
	DexBuyPrice  string    `json:"dexBuyPrice" comment:"dex买入价格"`
	DexSellPrice string    `json:"dexSellPrice" comment:"dex卖出价格"`
	CexBuyPrice  string    `json:"cexBuyPrice" comment:"cex买入价格"`
	CexSellPrice string    `json:"cexSellPrice" comment:"cex卖出价格"`
	SnapshotTime time.Time `json:"snapshotTime" comment:"快照时间"`
	common.ControlBy
}

func (s *BusDexCexPriceSpreadDataInsertReq) Generate(model *models.BusDexCexPriceSpreadData) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ObserverId = s.ObserverId
	model.Symbol = s.Symbol
	model.DexBuyPrice = s.DexBuyPrice
	model.DexSellPrice = s.DexSellPrice
	model.CexBuyPrice = s.CexBuyPrice
	model.CexSellPrice = s.CexSellPrice
	model.SnapshotTime = s.SnapshotTime
}

func (s *BusDexCexPriceSpreadDataInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexPriceSpreadDataUpdateReq struct {
	Id           int       `uri:"id" comment:""` //
	ObserverId   string    `json:"observerId" comment:"观察器id"`
	Symbol       string    `json:"symbol" comment:"观察币种"`
	DexBuyPrice  string    `json:"dexBuyPrice" comment:"dex买入价格"`
	DexSellPrice string    `json:"dexSellPrice" comment:"dex卖出价格"`
	CexBuyPrice  string    `json:"cexBuyPrice" comment:"cex买入价格"`
	CexSellPrice string    `json:"cexSellPrice" comment:"cex卖出价格"`
	SnapshotTime time.Time `json:"snapshotTime" comment:"快照时间"`
	common.ControlBy
}

func (s *BusDexCexPriceSpreadDataUpdateReq) Generate(model *models.BusDexCexPriceSpreadData) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ObserverId = s.ObserverId
	model.Symbol = s.Symbol
	model.DexBuyPrice = s.DexBuyPrice
	model.DexSellPrice = s.DexSellPrice
	model.CexBuyPrice = s.CexBuyPrice
	model.CexSellPrice = s.CexSellPrice
	model.SnapshotTime = s.SnapshotTime
}

func (s *BusDexCexPriceSpreadDataUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexPriceSpreadDataGetReq 功能获取请求参数
type BusDexCexPriceSpreadDataGetReq struct {
	Id int `uri:"id"`
}

func (s *BusDexCexPriceSpreadDataGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexPriceSpreadDataDeleteReq 功能删除请求参数
type BusDexCexPriceSpreadDataDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexCexPriceSpreadDataDeleteReq) GetId() interface{} {
	return s.Ids
}
