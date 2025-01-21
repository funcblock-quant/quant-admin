package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusPriceTriggerStrategyApikeyConfigGetPageReq struct {
	dto.Pagination `search:"-"`
	BusPriceTriggerStrategyApikeyConfigOrder
}

type BusPriceTriggerStrategyApikeyConfigOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:bus_price_trigger_strategy_apikey_config"`
	UserId    string `form:"userIdOrder"  search:"type:order;column:user_id;table:bus_price_trigger_strategy_apikey_config"`
	ApiKey    string `form:"apiKeyOrder"  search:"type:order;column:api_key;table:bus_price_trigger_strategy_apikey_config"`
	Username  string `form:"usernameOrder"  search:"type:order;column:username;table:bus_price_trigger_strategy_apikey_config"`
	Password  string `form:"passwordOrder"  search:"type:order;column:password;table:bus_price_trigger_strategy_apikey_config"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_price_trigger_strategy_apikey_config"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_price_trigger_strategy_apikey_config"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_price_trigger_strategy_apikey_config"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_price_trigger_strategy_apikey_config"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_price_trigger_strategy_apikey_config"`
}

func (m *BusPriceTriggerStrategyApikeyConfigGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusPriceTriggerStrategyApikeyConfigInsertReq struct {
	Id       int    `json:"-" comment:""` //
	UserId   string `json:"userId" comment:"用户id"`
	ApiKey   string `json:"apiKey" comment:"api key"`
	Username string `json:"username" comment:"用户名"`
	Password string `json:"password" comment:"密码"`
	common.ControlBy
}

func (s *BusPriceTriggerStrategyApikeyConfigInsertReq) Generate(model *models.BusPriceTriggerStrategyApikeyConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ApiKey = s.ApiKey
	model.Username = s.Username
	model.Password = s.Password
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusPriceTriggerStrategyApikeyConfigInsertReq) GetId() interface{} {
	return s.Id
}

type BusPriceTriggerStrategyApikeyConfigUpdateReq struct {
	Id       int    `uri:"id" comment:""` //
	UserId   string `json:"userId" comment:"用户id"`
	ApiKey   string `json:"apiKey" comment:"api key"`
	Username string `json:"username" comment:"用户名"`
	Password string `json:"password" comment:"密码"`
	common.ControlBy
}

func (s *BusPriceTriggerStrategyApikeyConfigUpdateReq) Generate(model *models.BusPriceTriggerStrategyApikeyConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ApiKey = s.ApiKey
	model.Username = s.Username
	model.Password = s.Password
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusPriceTriggerStrategyApikeyConfigUpdateReq) GetId() interface{} {
	return s.Id
}

// BusPriceTriggerStrategyApikeyConfigGetReq 功能获取请求参数
type BusPriceTriggerStrategyApikeyConfigGetReq struct {
	Id int `uri:"id"`
}

func (s *BusPriceTriggerStrategyApikeyConfigGetReq) GetId() interface{} 
func (s *BusPriceTriggerStrategyApikeyConfigGetReq) GetId() interface{} {
	return s.Id
}

// BusPriceTriggerStrategyApikeyConfigDeleteReq 功能删除请求参数
type BusPriceTriggerStrategyApikeyConfigDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusPriceTriggerStrategyApikeyConfigDeleteReq) GetId() interface{} {
	return s.Ids
}
