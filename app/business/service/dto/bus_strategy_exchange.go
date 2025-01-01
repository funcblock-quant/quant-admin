package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyExchangeGetPageReq struct {
	dto.Pagination `search:"-"`
	ExchangeName   string `form:"exchangeName"  search:"type:exact;column:exchange_name;table:bus_strategy_exchange" comment:"名称"`
	ExchangeType   string `form:"exchangeType"  search:"type:exact;column:exchange_type;table:bus_strategy_exchange" comment:"交易所类型"`
	BusStrategyExchangeOrder
}

type BusStrategyExchangeOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_exchange"`
	ExchangeName string `form:"exchangeNameOrder"  search:"type:order;column:exchange_name;table:bus_strategy_exchange"`
	ExchangeType string `form:"exchangeTypeOrder"  search:"type:order;column:exchange_type;table:bus_strategy_exchange"`
}

func (m *BusStrategyExchangeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyExchangeInsertReq struct {
	Id           int    `json:"-" comment:""` //
	ExchangeName string `json:"exchangeName" comment:"名称"`
	ExchangeType string `json:"exchangeType" comment:"交易所类型"`
	common.ControlBy
}

func (s *BusStrategyExchangeInsertReq) Generate(model *models.BusStrategyExchange) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ExchangeName = s.ExchangeName
	model.ExchangeType = s.ExchangeType
}

func (s *BusStrategyExchangeInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyExchangeUpdateReq struct {
	Id           int    `uri:"id" comment:""` //
	ExchangeName string `json:"exchangeName" comment:"名称"`
	ExchangeType string `json:"exchangeType" comment:"交易所类型"`
	common.ControlBy
}

func (s *BusStrategyExchangeUpdateReq) Generate(model *models.BusStrategyExchange) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ExchangeName = s.ExchangeName
	model.ExchangeType = s.ExchangeType
}

func (s *BusStrategyExchangeUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyExchangeGetReq 功能获取请求参数
type BusStrategyExchangeGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyExchangeGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyExchangeDeleteReq 功能删除请求参数
type BusStrategyExchangeDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyExchangeDeleteReq) GetId() interface{} {
	return s.Ids
}
