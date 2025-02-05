package dto

import (
	"time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexCexPriceSpreadStatisticsGetPageReq struct {
	dto.Pagination `search:"-"`
	BusDexCexPriceSpreadStatisticsOrder
}

type BusDexCexPriceSpreadStatisticsOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_price_spread_statistics"`
	ObserverId         string `form:"observerIdOrder"  search:"type:order;column:observer_id;table:bus_dex_cex_price_spread_statistics"`
	SpreadType         string `form:"spreadTypeOrder"  search:"type:order;column:spread_type;table:bus_dex_cex_price_spread_statistics"`
	Symbol             string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_price_spread_statistics"`
	StartTime          string `form:"startTimeOrder"  search:"type:order;column:start_time;table:bus_dex_cex_price_spread_statistics"`
	EndTime            string `form:"endTimeOrder"  search:"type:order;column:end_time;table:bus_dex_cex_price_spread_statistics"`
	Duration           string `form:"durationOrder"  search:"type:order;column:duration;table:bus_dex_cex_price_spread_statistics"`
	MaxPriceDifference string `form:"maxPriceDifferenceOrder"  search:"type:order;column:max_price_difference;table:bus_dex_cex_price_spread_statistics"`
	MinPriceDifference string `form:"minPriceDifferenceOrder"  search:"type:order;column:min_price_difference;table:bus_dex_cex_price_spread_statistics"`
}

func (m *BusDexCexPriceSpreadStatisticsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexPriceSpreadStatisticsInsertReq struct {
	Id                 int       `json:"-" comment:""` //
	ObserverId         string    `json:"observerId" comment:"观察器id"`
	SpreadType         string    `json:"spreadType" comment:"价差类型"`
	Symbol             string    `json:"symbol" comment:"观察币种"`
	StartTime          time.Time `json:"startTime" comment:"正向价差开始时间"`
	EndTime            time.Time `json:"endTime" comment:"正向价差结束时间"`
	Duration           string    `json:"duration" comment:"价差持续时间"`
	MaxPriceDifference string    `json:"maxPriceDifference" comment:"最大价差"`
	MinPriceDifference string    `json:"minPriceDifference" comment:"最小价差"`
	common.ControlBy
}

func (s *BusDexCexPriceSpreadStatisticsInsertReq) Generate(model *models.BusDexCexPriceSpreadStatistics) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ObserverId = s.ObserverId
	model.SpreadType = s.SpreadType
	model.Symbol = s.Symbol
	model.StartTime = &s.StartTime
	model.EndTime = &s.EndTime
	model.Duration = s.Duration
	model.MaxPriceDifference = s.MaxPriceDifference
	model.MinPriceDifference = s.MinPriceDifference
}

func (s *BusDexCexPriceSpreadStatisticsInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexPriceSpreadStatisticsUpdateReq struct {
	Id                 int       `uri:"id" comment:""` //
	ObserverId         string    `json:"observerId" comment:"观察器id"`
	SpreadType         string    `json:"spreadType" comment:"价差类型"`
	Symbol             string    `json:"symbol" comment:"观察币种"`
	StartTime          time.Time `json:"startTime" comment:"正向价差开始时间"`
	EndTime            time.Time `json:"endTime" comment:"正向价差结束时间"`
	Duration           string    `json:"duration" comment:"价差持续时间"`
	MaxPriceDifference string    `json:"maxPriceDifference" comment:"最大价差"`
	MinPriceDifference string    `json:"minPriceDifference" comment:"最小价差"`
	common.ControlBy
}

func (s *BusDexCexPriceSpreadStatisticsUpdateReq) Generate(model *models.BusDexCexPriceSpreadStatistics) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ObserverId = s.ObserverId
	model.SpreadType = s.SpreadType
	model.Symbol = s.Symbol
	model.StartTime = &s.StartTime
	model.EndTime = &s.EndTime
	model.Duration = s.Duration
	model.MaxPriceDifference = s.MaxPriceDifference
	model.MinPriceDifference = s.MinPriceDifference
}

func (s *BusDexCexPriceSpreadStatisticsUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexPriceSpreadStatisticsGetReq 功能获取请求参数
type BusDexCexPriceSpreadStatisticsGetReq struct {
	Id int `uri:"id"`
}

func (s *BusDexCexPriceSpreadStatisticsGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexPriceSpreadStatisticsDeleteReq 功能删除请求参数
type BusDexCexPriceSpreadStatisticsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexCexPriceSpreadStatisticsDeleteReq) GetId() interface{} {
	return s.Ids
}
