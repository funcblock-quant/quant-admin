package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyInstanceConfigGetPageReq struct {
	dto.Pagination     `search:"-"`
	StrategyInstanceId string `form:"strategyInstanceId"  search:"type:exact;column:strategy_instance_id;table:bus_strategy_instance_config" comment:"策略实例id"`
	BusStrategyInstanceConfigOrder
}

type BusStrategyInstanceConfigOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_instance_config"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_strategy_instance_config"`
	SchemaText         string `form:"schemaTextOrder"  search:"type:order;column:schema_text;table:bus_strategy_config_schema"`
	SchemaType         string `form:"schemaTypeOrder"  search:"type:order;column:schema_type;table:bus_strategy_config_schema"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_instance_config"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_instance_config"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_instance_config"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_instance_config"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_instance_config"`
}

func (m *BusStrategyInstanceConfigGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyInstanceConfigGetByInstanceIdReq struct {
	StrategyInstanceId string
}

type BusStrategyInstanceConfigInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	SchemaText         string `json:"schemaText" comment:"参数schema"`
	SchemaType         string `json:"schemaType" comment:"schema类型"`
	common.ControlBy
}

func (s *BusStrategyInstanceConfigInsertReq) Generate(model *models.BusStrategyInstanceConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.SchemaType = s.SchemaType
	model.SchemaText = s.SchemaText
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyInstanceConfigInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyInstanceConfigUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	StrategyInstanceId string `json:"strategyInstanceId" comment:"策略实例id"`
	SchemaText         string `json:"schemaText" comment:"参数schema"`
	SchemaType         string `json:"schemaType" comment:"schema类型"`
	common.ControlBy
}

func (s *BusStrategyInstanceConfigUpdateReq) Generate(model *models.BusStrategyInstanceConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyInstanceId = s.StrategyInstanceId
	model.SchemaType = s.SchemaType
	model.SchemaText = s.SchemaText
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusStrategyInstanceConfigUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceConfigGetReq 功能获取请求参数
type BusStrategyInstanceConfigGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyInstanceConfigGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceConfigDeleteReq 功能删除请求参数
type BusStrategyInstanceConfigDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyInstanceConfigDeleteReq) GetId() interface{} {
	return s.Ids
}
