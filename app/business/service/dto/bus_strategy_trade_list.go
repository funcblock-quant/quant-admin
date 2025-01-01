package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyTradeListGetPageReq struct {
	dto.Pagination     `search:"-"`
	StrategyInstanceId string `form:"strategyInstanceId"  search:"type:exact;column:strategy_instance_id;table:bus_strategy_trade_list" comment:"策略实例id"`
	BusStrategyTradeListOrder
}

type BusStrategyTradeListOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_trade_list"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_strategy_trade_list"`
	Symbol             string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_strategy_trade_list"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_trade_list"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_trade_list"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_trade_list"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_trade_list"`
	IsDeleted          string `form:"isDeletedOrder"  search:"type:order;column:is_deleted;table:bus_strategy_trade_list"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_trade_list"`
}

func (m *BusStrategyTradeListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyTradeListInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	Symbol             string `json:"symbol" comment:"交易币种名称"`
	IsDeleted          string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyTradeListInsertReq) Generate(model *models.BusStrategyTradeList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.Symbol = s.Symbol
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.IsDeleted = s.IsDeleted
}

func (s *BusStrategyTradeListInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyTradeListUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	Symbol             string `json:"symbol" comment:"交易币种名称"`
	IsDeleted          string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyTradeListUpdateReq) Generate(model *models.BusStrategyTradeList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.Symbol = s.Symbol
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.IsDeleted = s.IsDeleted
}

func (s *BusStrategyTradeListUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyTradeListGetReq 功能获取请求参数
type BusStrategyTradeListGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyTradeListGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyTradeListDeleteReq 功能删除请求参数
type BusStrategyTradeListDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyTradeListDeleteReq) GetId() interface{} {
	return s.Ids
}
