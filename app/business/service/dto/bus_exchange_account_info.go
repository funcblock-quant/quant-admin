package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusExchangeAccountInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	AccountName    string `form:"accountName"  search:"type:exact;column:account_name;table:bus_exchange_account_info" comment:"钱包名称"`
	ExchangeId     string `form:"exchangeId"  search:"type:exact;column:exchange_id;table:bus_exchange_account_info" comment:"id"`
	ExchangeName   string `form:"exchangeName"  search:"type:exact;column:exchange_name;table:bus_exchange_account_info" comment:"交易所名称"`
	AccountType    string `form:"accountType"  search:"type:exact;column:account_type;table:bus_exchange_account_info" comment:"账户类型"`
	Status         string `form:"status"  search:"type:exact;column:status;table:bus_exchange_account_info" comment:"状态"`
	BusExchangeAccountInfoOrder
}

type BusExchangeAccountInfoOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:bus_exchange_account_info"`
	AccountName  string `form:"accountNameOrder"  search:"type:order;column:account_name;table:bus_exchange_account_info"`
	ExchangeId   string `form:"exchangeIdOrder"  search:"type:order;column:exchange_id;table:bus_exchange_account_info"`
	ExchangeName string `form:"exchangeNameOrder"  search:"type:order;column:exchange_name;table:bus_exchange_account_info"`
	Uid          string `form:"uidOrder"  search:"type:order;column:uid;table:bus_exchange_account_info"`
	AccountType  string `form:"accountTypeOrder"  search:"type:order;column:account_type;table:bus_exchange_account_info"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:bus_exchange_account_info"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_exchange_account_info"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_exchange_account_info"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_exchange_account_info"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_exchange_account_info"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_exchange_account_info"`
}

func (m *BusExchangeAccountInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusExchangeAccountInfoInsertReq struct {
	Id              int      `json:"-" comment:""` //
	AccountName     string   `json:"accountName" comment:"钱包名称"`
	ExchangeId      string   `json:"exchangeId" comment:"id"`
	ExchangeName    string   `json:"exchangeName" comment:"交易所名称"`
	Uid             string   `json:"uid" comment:"交易所uid"`
	AccountType     string   `json:"accountType" comment:"账户类型"`
	Status          string   `json:"status" comment:"状态"`
	AccountGroupIds []string `json:"accountGroupIds" comment:"绑定的账户组id"`
	common.ControlBy
}

func (s *BusExchangeAccountInfoInsertReq) Generate(model *models.BusExchangeAccountInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.AccountName = s.AccountName
	model.Uid = s.Uid
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusExchangeAccountInfoInsertReq) GetId() interface{} {
	return s.Id
}

type BusExchangeAccountInfoUpdateReq struct {
	Id           int    `uri:"id" comment:""` //
	AccountName  string `json:"accountName" comment:"钱包名称"`
	ExchangeType string `json:"exchangeType" comment:"交易所名称"`
	Uid          string `json:"uid" comment:"交易所uid"`
	Status       int8   `json:"status" comment:"状态"`
	common.ControlBy
}

func (s *BusExchangeAccountInfoUpdateReq) Generate(model *models.BusExchangeAccountInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.AccountName = s.AccountName
	model.ExchangeType = s.ExchangeType
	model.Status = s.Status
	model.Uid = s.Uid
	model.UpdateBy = s.UpdateBy // 添加这行，需要记录是被谁更新的
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

type BusGroupAccountInfoGetReq struct {
	GroupId int `uri:"groupId"`
}

func (s *BusGroupAccountInfoGetReq) GetId() interface{} {
	return s.GroupId
}

type CexExchangeListResp struct {
	ExchangeType string `json:"exchangeType" gorm:"column:exchange_type"`
}

type ProtfolioUnwindingInfoReq struct {
	TokenName    string `json:"tokenName" comment:"token名称"`
	TokenAddress string `json:"tokenAddress" comment:"token地址"`
	CexAccountId int64  `json:"cexAccountId" comment:"交易所账户id"`
	DexWalletId  int64  `json:"dexWalletId" comment:"钱包id"`
}

type ProtfolioUnwindingInfoResp struct {
	TokenName                  string `json:"tokenName" comment:"token名称"`
	WalletBalance              string `json:"walletBalance" comment:"钱包余额"`
	TraderAccountMarginBalance string `json:"traderAccountMarginBalance" comment:"交易所杠杆账户余额"`
	TraderAccountSpotBalance   string `json:"traderAccountSpotBalance" comment:"交易所现货账户余额"`
	MasterAccountSpotBalance   string `json:"masterAccountSpotBalance" comment:"主账户现货账户余额"`
	TraderAccountBorrowed      string `json:"traderAccountBorrowed" comment:"交易所借贷账户余额"`
}
