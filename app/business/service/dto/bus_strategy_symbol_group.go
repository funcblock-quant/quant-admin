package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategySymbolGroupGetPageReq struct {
	dto.Pagination     `search:"-"`
	StrategyInstanceId string `form:"strategyInstanceId"  search:"type:exact;column:strategy_instance_id;table:bus_strategy_symbol_group" comment:"策略实例id"`
	GroupType          string `form:"groupType"  search:"type:exact;column:group_type;table:bus_strategy_symbol_group" comment:"组类型"`
	BusStrategySymbolGroupOrder
}

type BusStrategySymbolGroupOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_symbol_group"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_strategy_symbol_group"`
	GroupName          string `form:"groupNameOrder"  search:"type:order;column:group_name;table:bus_strategy_symbol_group"`
	GroupType          string `form:"groupTypeOrder"  search:"type:order;column:group_type;table:bus_strategy_symbol_group"`
	AutoRefresh        string `form:"autoRefreshOrder"  search:"type:order;column:auto_refresh;table:bus_strategy_symbol_group"`
	RefreshInterval    string `form:"refreshIntervalOrder"  search:"type:order;column:refresh_interval;table:bus_strategy_symbol_group"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_symbol_group"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_symbol_group"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_symbol_group"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_symbol_group"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_symbol_group"`
}

func (m *BusStrategySymbolGroupGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategySymbolGroupInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	GroupName          string `json:"groupName" comment:"观察交易对名称"`
	GroupType          string `json:"groupType" comment:"组类型"`
	AutoRefresh        bool   `json:"autoRefresh" comment:"是否自动刷新"`
	IsActive           bool   `json:"IsActive" comment:"是否激活"`
	RefreshInterval    string `json:"refreshInterval" comment:"自动刷新间隔"`
	common.ControlBy
}

func (s *BusStrategySymbolGroupInsertReq) Generate(model *models.BusStrategySymbolGroup) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.GroupName = s.GroupName
	model.GroupType = s.GroupType
	model.AutoRefresh = s.AutoRefresh
	model.RefreshInterval = s.RefreshInterval
	model.IsActive = s.IsActive
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategySymbolGroupInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategySymbolGroupUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	GroupName          string `json:"groupName" comment:"观察交易对名称"`
	GroupType          string `json:"groupType" comment:"组类型"`
	AutoRefresh        bool   `json:"autoRefresh" comment:"是否自动刷新"`
	IsActive           bool   `json:"IsActive" comment:"是否激活"`
	RefreshInterval    string `json:"refreshInterval" comment:"自动刷新间隔"`
	common.ControlBy
}

func (s *BusStrategySymbolGroupUpdateReq) Generate(model *models.BusStrategySymbolGroup) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.GroupName = s.GroupName
	model.GroupType = s.GroupType
	model.AutoRefresh = s.AutoRefresh
	model.IsActive = s.IsActive
	model.RefreshInterval = s.RefreshInterval
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusStrategySymbolGroupUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategySymbolGroupGetReq 功能获取请求参数
type BusStrategySymbolGroupGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategySymbolGroupGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategySymbolGroupDeleteReq 功能删除请求参数
type BusStrategySymbolGroupDeleteReq struct {
	Ids []string `json:"ids"`
}

func (s *BusStrategySymbolGroupDeleteReq) GetId() interface{} {
	return s.Ids
}
