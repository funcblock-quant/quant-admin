package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyConfigDictGetPageReq struct {
	dto.Pagination `search:"-"`
	StrategyId     string `form:"strategyId"  search:"type:exact;column:strategy_id;table:bus_strategy_config_dict" comment:"id"`
	BusStrategyConfigDictOrder
}

type BusStrategyConfigDictOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_config_dict"`
	StrategyId   string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_strategy_config_dict"`
	ParamKey     string `form:"paramKeyOrder"  search:"type:order;column:param_key;table:bus_strategy_config_dict"`
	ParamName    string `form:"paramNameOrder"  search:"type:order;column:param_name;table:bus_strategy_config_dict"`
	ParamType    string `form:"paramTypeOrder"  search:"type:order;column:param_type;table:bus_strategy_config_dict"`
	DefaultValue string `form:"defaultValueOrder"  search:"type:order;column:default_value;table:bus_strategy_config_dict"`
	Required     string `form:"requiredOrder"  search:"type:order;column:required;table:bus_strategy_config_dict"`
	Description  string `form:"descriptionOrder"  search:"type:order;column:description;table:bus_strategy_config_dict"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_config_dict"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_config_dict"`
}

func (m *BusStrategyConfigDictGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyConfigDictInsertReq struct {
	Id           int    `json:"-" comment:"唯一标识"` // 唯一标识
	StrategyId   string `json:"strategyId" comment:"id"`
	ParamKey     string `json:"paramKey" comment:"参数的唯一标识"`
	ParamName    string `json:"paramName" comment:"参数名称"`
	ParamType    string `json:"paramType" comment:"参数类型"`
	DefaultValue string `json:"defaultValue" comment:"参数的默认值"`
	Required     string `json:"required" comment:"是否为必填参数"`
	Description  string `json:"description" comment:"参数用途描述"`
	common.ControlBy
}

func (s *BusStrategyConfigDictInsertReq) Generate(model *models.BusStrategyConfigDict) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.ParamKey = s.ParamKey
	model.ParamName = s.ParamName
	model.ParamType = s.ParamType
	model.DefaultValue = s.DefaultValue
	model.Required = s.Required
	model.Description = s.Description
}

func (s *BusStrategyConfigDictInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyConfigDictUpdateReq struct {
	Id           int    `uri:"id" comment:"唯一标识"` // 唯一标识
	StrategyId   string `json:"strategyId" comment:"id"`
	ParamKey     string `json:"paramKey" comment:"参数的唯一标识"`
	ParamName    string `json:"paramName" comment:"参数名称"`
	ParamType    string `json:"paramType" comment:"参数类型"`
	DefaultValue string `json:"defaultValue" comment:"参数的默认值"`
	Required     string `json:"required" comment:"是否为必填参数"`
	Description  string `json:"description" comment:"参数用途描述"`
	common.ControlBy
}

func (s *BusStrategyConfigDictUpdateReq) Generate(model *models.BusStrategyConfigDict) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.ParamKey = s.ParamKey
	model.ParamName = s.ParamName
	model.ParamType = s.ParamType
	model.DefaultValue = s.DefaultValue
	model.Required = s.Required
	model.Description = s.Description
}

func (s *BusStrategyConfigDictUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyConfigDictGetReq 功能获取请求参数
type BusStrategyConfigDictGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyConfigDictGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyConfigDictDeleteReq 功能删除请求参数
type BusStrategyConfigDictDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyConfigDictDeleteReq) GetId() interface{} {
	return s.Ids
}
