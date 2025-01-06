package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyWatchListGetPageReq struct {
	dto.Pagination `search:"-"`
	SymbolGroupId  string `form:"symbolGroupId"  search:"type:exact;column:symbol_group_id;table:bus_strategy_watch_list" comment:"策略实例id"`
	BusStrategyWatchListOrder
}

type BusStrategyWatchListOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_watch_list"`
	SymbolGroupId string `form:"symbolGroupIdOrder"  search:"type:order;column:symbol_group_id;table:bus_strategy_watch_list"`
	Symbol        string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_strategy_watch_list"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_watch_list"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_watch_list"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_watch_list"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_watch_list"`
	IsDeleted     string `form:"isDeletedOrder"  search:"type:order;column:is_deleted;table:bus_strategy_watch_list"`
	DeletedAt     string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_watch_list"`
}

type BusStrategyWatchListGetPageResp struct {
	Id            string `json:"id"`
	SymbolGroupId string `json:"symbolGroupId" gorm:"type:bigint;comment:币对组id"`
	Symbol        string `json:"symbol" gorm:"type:varchar(64);comment:观察币种名称"`
	GroupName     string `json:"groupName"`
}

func (m *BusStrategyWatchListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyWatchListInsertReq struct {
	Id            int    `json:"-" comment:""` //
	SymbolGroupId string `json:"symbolGroupId" comment:"策略实例币种组id"`
	Symbol        string `json:"symbol" comment:"观察币种名称"`
	IsDeleted     string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyWatchListInsertReq) Generate(model *models.BusStrategyWatchList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SymbolGroupId = s.SymbolGroupId
	model.Symbol = s.Symbol
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyWatchListInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyWatchListUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	SymbolGroupId string `json:"symbolGroupId" comment:"策略实例币种组id"`
	Symbol        string `json:"symbol" comment:"观察币种名称"`
	IsDeleted     string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyWatchListUpdateReq) Generate(model *models.BusStrategyWatchList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SymbolGroupId = s.SymbolGroupId
	model.Symbol = s.Symbol
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusStrategyWatchListUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyWatchListGetReq 功能获取请求参数
type BusStrategyWatchListGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyWatchListGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyWatchListDeleteReq 功能删除请求参数
type BusStrategyWatchListDeleteReq struct {
	Ids []string `json:"ids"`
}

func (s *BusStrategyWatchListDeleteReq) GetId() interface{} {
	return s.Ids
}
