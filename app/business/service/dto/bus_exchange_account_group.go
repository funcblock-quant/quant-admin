package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusExchangeAccountGroupGetPageReq struct {
	dto.Pagination `search:"-"`
	GroupName      string `form:"groupName"  search:"type:contains;column:group_name;table:bus_exchange_account_group" comment:"交易所账户组"`
	BusExchangeAccountGroupOrder
}

type BusExchangeAccountGroupOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:bus_exchange_account_group"`
	GroupName string `form:"groupNameOrder"  search:"type:order;column:group_name;table:bus_exchange_account_group"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_exchange_account_group"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_exchange_account_group"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_exchange_account_group"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_exchange_account_group"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_exchange_account_group"`
}

func (m *BusExchangeAccountGroupGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusExchangeAccountGroupInsertReq struct {
	Id          int    `json:"-" comment:""` //
	GroupName   string `json:"groupName" comment:"交易所账户组"`
	Description string `json:"description" comment:"描述"`
	AccountIds  []int  `json:"accountIds"`
	common.ControlBy
}

type BusExchangeAccountGroupDetailResp struct {
	models.BusExchangeAccountGroup
	AccountIds []int `json:"accountIds"`
	common.ControlBy
}

func (s *BusExchangeAccountGroupInsertReq) Generate(model *models.BusExchangeAccountGroup) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GroupName = s.GroupName
	model.Description = s.Description
	model.CreateBy = s.CreateBy
}

func (s *BusExchangeAccountGroupInsertReq) GetId() interface{} {
	return s.Id
}

type BusExchangeAccountGroupUpdateReq struct {
	Id          int    `uri:"id" comment:""` //
	GroupName   string `json:"groupName" comment:"交易所账户组"`
	Description string `json:"description" comment:"描述"`
	AccountIds  []int  `json:"accountIds"`
	common.ControlBy
}

func (s *BusExchangeAccountGroupUpdateReq) Generate(model *models.BusExchangeAccountGroup) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GroupName = s.GroupName
	model.Description = s.Description
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusExchangeAccountGroupUpdateReq) GetId() interface{} {
	return s.Id
}

// BusExchangeAccountGroupGetReq 功能获取请求参数
type BusExchangeAccountGroupGetReq struct {
	Id int `uri:"id"`
}

func (s *BusExchangeAccountGroupGetReq) GetId() interface{} {
	return s.Id
}

// BusExchangeAccountGroupDeleteReq 功能删除请求参数
type BusExchangeAccountGroupDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusExchangeAccountGroupDeleteReq) GetId() interface{} {
	return s.Ids
}

type BusAccountGroupListGetReq struct {
	AccountId int `uri:"accountId"`
}

func (s *BusAccountGroupListGetReq) GetId() interface{} {
	return s.AccountId
}
