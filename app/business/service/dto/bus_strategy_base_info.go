package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyBaseInfoGetPageReq struct {
	dto.Pagination   `search:"-"`
	StrategyName     string `form:"strategyName"  search:"type:contains;column:strategy_name;table:bus_strategy_base_info" comment:"策略名称"`
	StrategyCategory string `form:"strategyCategory"  search:"type:exact;column:strategy_category;table:bus_strategy_base_info" comment:"策略交易类型"`
	Status           string `form:"status"  search:"type:exact;column:status;table:bus_strategy_base_info" comment:"策略运行状态"`
	BusStrategyBaseInfoOrder
}

type BusStrategyBaseInfoOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_base_info"`
	StrategyName     string `form:"strategyNameOrder"  search:"type:order;column:strategy_name;table:bus_strategy_base_info"`
	StrategyCategory string `form:"strategyCategoryOrder"  search:"type:order;column:strategy_category;table:bus_strategy_base_info"`
	Preference       string `form:"preferenceOrder"  search:"type:order;column:preference;table:bus_strategy_base_info"`
	Description      string `form:"descriptionOrder"  search:"type:order;column:description;table:bus_strategy_base_info"`
	Status           string `form:"statusOrder"  search:"type:order;column:status;table:bus_strategy_base_info"`
	Owner            string `form:"ownerOrder"  search:"type:order;column:owner;table:bus_strategy_base_info"`
	CreatedAt        string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_base_info"`
	UpdatedAt        string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_base_info"`
	CreateBy         string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_base_info"`
	UpdateBy         string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_base_info"`
	DeletedAt        string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_base_info"`
}

func (m *BusStrategyBaseInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyBaseInfoInsertReq struct {
	Id               int    `json:"-" comment:""` //
	StrategyName     string `json:"strategyName" comment:"策略名称"`
	StrategyCategory string `json:"strategyCategory" comment:"策略交易类型"`
	Preference       string `json:"preference" comment:"策略偏好"`
	Description      string `json:"description" comment:"策略描述"`
	Status           string `json:"status" comment:"策略运行状态"`
	Owner            string `json:"owner" comment:"策略负责人"`
	common.ControlBy
	Schema BusStrategyConfigSchemaInsertReq
}

func (s *BusStrategyBaseInfoInsertReq) Generate(model *models.BusStrategyBaseInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyName = s.StrategyName
	model.StrategyCategory = s.StrategyCategory
	model.Preference = s.Preference
	model.Description = s.Description
	model.Status = s.Status
	model.Owner = s.Owner
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyBaseInfoInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyBaseInfoUpdateReq struct {
	Id               int    `uri:"id" comment:""` //
	StrategyName     string `json:"strategyName" comment:"策略名称"`
	StrategyCategory string `json:"strategyCategory" comment:"策略交易类型"`
	Preference       string `json:"preference" comment:"策略偏好"`
	Description      string `json:"description" comment:"策略描述"`
	Status           string `json:"status" comment:"策略运行状态"`
	Owner            string `json:"owner" comment:"策略负责人"`
	IsDeleted        string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
	Schema BusStrategyConfigSchemaUpdateReq
}

func (s *BusStrategyBaseInfoUpdateReq) Generate(model *models.BusStrategyBaseInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyName = s.StrategyName
	model.StrategyCategory = s.StrategyCategory
	model.Preference = s.Preference
	model.Description = s.Description
	model.Status = s.Status
	model.Owner = s.Owner
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusStrategyBaseInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyBaseInfoGetReq 功能获取请求参数
type BusStrategyBaseInfoGetReq struct {
	StrategyId int `uri:"id"`
}

func (s *BusStrategyBaseInfoGetReq) GetId() interface{} {
	return s.StrategyId
}

// BusStrategyBaseInfoDeleteReq 功能删除请求参数
type BusStrategyBaseInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyBaseInfoDeleteReq) GetId() interface{} {
	return s.Ids
}
