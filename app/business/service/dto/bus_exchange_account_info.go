package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusExchangeAccountInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	AccountName    string `form:"accountName"  search:"type:exact;column:account_name;table:bus_exchange_account_info" comment:"钱包名称"`
	PlatformId     string `form:"platformId"  search:"type:exact;column:platform_id;table:bus_exchange_account_info" comment:"id"`
	PlatformName   string `form:"platformName"  search:"type:exact;column:platform_name;table:bus_exchange_account_info" comment:"交易所名称"`
	AccountType    string `form:"accountType"  search:"type:exact;column:account_type;table:bus_exchange_account_info" comment:"账户类型"`
	Status         string `form:"status"  search:"type:exact;column:status;table:bus_exchange_account_info" comment:"状态"`
	BusExchangeAccountInfoOrder
}

type BusExchangeAccountInfoOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:bus_exchange_account_info"`
	AccountName  string `form:"accountNameOrder"  search:"type:order;column:account_name;table:bus_exchange_account_info"`
	PlatformId   string `form:"platformIdOrder"  search:"type:order;column:platform_id;table:bus_exchange_account_info"`
	PlatformName string `form:"platformNameOrder"  search:"type:order;column:platform_name;table:bus_exchange_account_info"`
	Uid          string `form:"uidOrder"  search:"type:order;column:uid;table:bus_exchange_account_info"`
	AccountType  string `form:"accountTypeOrder"  search:"type:order;column:account_type;table:bus_exchange_account_info"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:bus_exchange_account_info"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_exchange_account_info"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_exchange_account_info"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_exchange_account_info"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_exchange_account_info"`
	IsDeleted    string `form:"isDeletedOrder"  search:"type:order;column:is_deleted;table:bus_exchange_account_info"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_exchange_account_info"`
}

func (m *BusExchangeAccountInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusExchangeAccountInfoInsertReq struct {
	Id           int    `json:"-" comment:""` //
	AccountName  string `json:"accountName" comment:"钱包名称"`
	PlatformId   string `json:"platformId" comment:"id"`
	PlatformName string `json:"platformName" comment:"交易所名称"`
	Uid          string `json:"uid" comment:"交易所uid"`
	AccountType  string `json:"accountType" comment:"账户类型"`
	Status       string `json:"status" comment:"状态"`
	IsDeleted    string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusExchangeAccountInfoInsertReq) Generate(model *models.BusExchangeAccountInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.AccountName = s.AccountName
	model.PlatformId = s.PlatformId
	model.PlatformName = s.PlatformName
	model.Uid = s.Uid
	model.AccountType = s.AccountType
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.IsDeleted = s.IsDeleted
}

func (s *BusExchangeAccountInfoInsertReq) GetId() interface{} {
	return s.Id
}

type BusExchangeAccountInfoUpdateReq struct {
	Id           int    `uri:"id" comment:""` //
	AccountName  string `json:"accountName" comment:"钱包名称"`
	PlatformId   string `json:"platformId" comment:"id"`
	PlatformName string `json:"platformName" comment:"交易所名称"`
	Uid          string `json:"uid" comment:"交易所uid"`
	AccountType  string `json:"accountType" comment:"账户类型"`
	Status       string `json:"status" comment:"状态"`
	IsDeleted    string `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusExchangeAccountInfoUpdateReq) Generate(model *models.BusExchangeAccountInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.AccountName = s.AccountName
	model.PlatformId = s.PlatformId
	model.PlatformName = s.PlatformName
	model.Uid = s.Uid
	model.AccountType = s.AccountType
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.IsDeleted = s.IsDeleted
}

func (s *BusExchangeAccountInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// BusExchangeAccountInfoGetReq 功能获取请求参数
type BusExchangeAccountInfoGetReq struct {
	Id int `uri:"id"`
}

func (s *BusExchangeAccountInfoGetReq) GetId() interface{} {
	return s.Id
}

// BusExchangeAccountInfoDeleteReq 功能删除请求参数
type BusExchangeAccountInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusExchangeAccountInfoDeleteReq) GetId() interface{} {
	return s.Ids
}
