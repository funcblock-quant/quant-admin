package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyTradeListGetPageReq struct {
	dto.Pagination `search:"-"`
	SymbolGroupId  string `form:"symbolGroupId"  search:"type:exact;column:symbol_group_id;table:bus_strategy_trade_list" comment:"策略实例id"`
	BusStrategyTradeListOrder
}

type BusStrategyTradeListOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_trade_list"`
	Symbol    string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_strategy_trade_list"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_trade_list"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_trade_list"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_trade_list"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_trade_list"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_trade_list"`
}

type BusStrategyTradeListGetPageResp struct {
	Id            string `json:"id"`
	SymbolGroupId string `json:"symbolGroupId" gorm:"type:bigint;comment:币对组id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:交易币种名称"`
	GroupName     string `json:"groupName"`
}

func (m *BusStrategyTradeListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyTradeListInsertReq struct {
	Id            int    `json:"-" comment:""` //
	SymbolGroupId string `json:"symbolGroupId" comment:"币对组id"`
	Symbol        string `json:"symbol" comment:"交易币种名称"`
	common.ControlBy
}

func (s *BusStrategyTradeListInsertReq) Generate(model *models.BusStrategyTradeList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SymbolGroupId = s.SymbolGroupId
	model.Symbol = s.Symbol
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyTradeListInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyTradeListUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	SymbolGroupId string `json:"symbolGroupId" comment:"币对组id"`
	Symbol        string `json:"symbol" comment:"交易币种名称"`
	common.ControlBy
}

func (s *BusStrategyTradeListUpdateReq) Generate(model *models.BusStrategyTradeList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SymbolGroupId = s.SymbolGroupId
	model.Symbol = s.Symbol
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
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
	Ids []string `json:"ids"`
}

func (s *BusStrategyTradeListDeleteReq) GetId() interface{} {
	return s.Ids
}
