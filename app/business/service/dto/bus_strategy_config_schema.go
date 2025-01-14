package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyConfigSchemaGetPageReq struct {
	dto.Pagination `search:"-"`
	StrategyId     string `form:"strategyId"  search:"type:exact;column:strategy_id;table:bus_strategy_config_schema" comment:"id"`
	BusStrategyConfigSchemaOrder
}

type BusStrategyConfigSchemaOrder struct {
	Id         string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_config_schema"`
	StrategyId string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_strategy_config_schema"`
	SchemaText string `form:"schemaTextOrder"  search:"type:order;column:schema_text;table:bus_strategy_config_schema"`
	SchemaType string `form:"schemaTypeOrder"  search:"type:order;column:schema_type;table:bus_strategy_config_schema"`
	CreateBy   string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_config_schema"`
	UpdateBy   string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_config_schema"`
	CreatedAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_config_schema"`
	UpdatedAt  string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_config_schema"`
	DeletedAt  string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_config_schema"`
}

func (m *BusStrategyConfigSchemaGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyConfigSchemaInsertReq struct {
	Id         int    `json:"-" comment:"唯一标识"` // 唯一标识
	StrategyId string `json:"strategyId" comment:"关联策略表的ID"`
	SchemaText string `json:"schemaText" comment:"参数schema"`
	SchemaType string `json:"schemaType" comment:"schema类型"`
	common.ControlBy
}

func (s *BusStrategyConfigSchemaInsertReq) Generate(model *models.BusStrategyConfigSchema) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.SchemaText = s.SchemaText
	model.SchemaType = s.SchemaType
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyConfigSchemaInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyConfigSchemaUpdateReq struct {
	Id         int    `uri:"id" comment:"唯一标识"` // 唯一标识
	StrategyId string `json:"strategyId" comment:"关联策略表的ID"`
	SchemaText string `json:"schemaText" comment:"参数schema"`
	SchemaType string `json:"schemaType" comment:"schema类型"`
	common.ControlBy
}

func (s *BusStrategyConfigSchemaUpdateReq) Generate(model *models.BusStrategyConfigSchema) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.SchemaText = s.SchemaText
	model.SchemaType = s.SchemaType
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusStrategyConfigSchemaUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyConfigSchemaGetReq 功能获取请求参数
type BusStrategyConfigSchemaGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyConfigSchemaGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyConfigSchemaDeleteReq 功能删除请求参数
type BusStrategyConfigSchemaDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyConfigSchemaDeleteReq) GetId() interface{} {
	return s.Ids
}
