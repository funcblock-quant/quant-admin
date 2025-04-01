package dto

import (

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexCexDebitCreditRecordGetPageReq struct {
	dto.Pagination     `search:"-"`
    AccountName string `form:"accountName"  search:"type:exact;column:account_name;table:bus_dex_cex_debit_credit_record" comment:"交易所账户名称"`
    Uid string `form:"uid"  search:"type:exact;column:uid;table:bus_dex_cex_debit_credit_record" comment:"交易所账户uid"`
    ExchangeType string `form:"exchangeType"  search:"type:exact;column:exchange_type;table:bus_dex_cex_debit_credit_record" comment:"交易所"`
    DebitType string `form:"debitType"  search:"type:exact;column:debit_type;table:bus_dex_cex_debit_credit_record" comment:"类型"`
    Status string `form:"status"  search:"type:exact;column:status;table:bus_dex_cex_debit_credit_record" comment:"状态"`
    BusDexCexDebitCreditRecordOrder
}

type BusDexCexDebitCreditRecordOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_debit_credit_record"`
    AccountName string `form:"accountNameOrder"  search:"type:order;column:account_name;table:bus_dex_cex_debit_credit_record"`
    Uid string `form:"uidOrder"  search:"type:order;column:uid;table:bus_dex_cex_debit_credit_record"`
    ExchangeType string `form:"exchangeTypeOrder"  search:"type:order;column:exchange_type;table:bus_dex_cex_debit_credit_record"`
    Symbol string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_debit_credit_record"`
    Amount string `form:"amountOrder"  search:"type:order;column:amount;table:bus_dex_cex_debit_credit_record"`
    DebitType string `form:"debitTypeOrder"  search:"type:order;column:debit_type;table:bus_dex_cex_debit_credit_record"`
    Status string `form:"statusOrder"  search:"type:order;column:status;table:bus_dex_cex_debit_credit_record"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_debit_credit_record"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_dex_cex_debit_credit_record"`
    
}

func (m *BusDexCexDebitCreditRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexDebitCreditRecordInsertReq struct {
    Id int `json:"-" comment:"主键 ID"` // 主键 ID
    AccountName string `json:"accountName" comment:"交易所账户名称"`
    Uid string `json:"uid" comment:"交易所账户uid"`
    ExchangeType string `json:"exchangeType" comment:"交易所"`
    Symbol string `json:"symbol" comment:"借贷币种"`
    Amount string `json:"amount" comment:"借还贷数量"`
    DebitType string `json:"debitType" comment:"类型"`
    Status string `json:"status" comment:"状态"`
    common.ControlBy
}

func (s *BusDexCexDebitCreditRecordInsertReq) Generate(model *models.BusDexCexDebitCreditRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.AccountName = s.AccountName
    model.Uid = s.Uid
    model.ExchangeType = s.ExchangeType
    model.Symbol = s.Symbol
    model.Amount = s.Amount
    model.DebitType = s.DebitType
    model.Status = s.Status
}

func (s *BusDexCexDebitCreditRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexDebitCreditRecordUpdateReq struct {
    Id int `uri:"id" comment:"主键 ID"` // 主键 ID
    AccountName string `json:"accountName" comment:"交易所账户名称"`
    Uid string `json:"uid" comment:"交易所账户uid"`
    ExchangeType string `json:"exchangeType" comment:"交易所"`
    Symbol string `json:"symbol" comment:"借贷币种"`
    Amount string `json:"amount" comment:"借还贷数量"`
    DebitType string `json:"debitType" comment:"类型"`
    Status string `json:"status" comment:"状态"`
    common.ControlBy
}

func (s *BusDexCexDebitCreditRecordUpdateReq) Generate(model *models.BusDexCexDebitCreditRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.AccountName = s.AccountName
    model.Uid = s.Uid
    model.ExchangeType = s.ExchangeType
    model.Symbol = s.Symbol
    model.Amount = s.Amount
    model.DebitType = s.DebitType
    model.Status = s.Status
}

func (s *BusDexCexDebitCreditRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexDebitCreditRecordGetReq 功能获取请求参数
type BusDexCexDebitCreditRecordGetReq struct {
     Id int `uri:"id"`
}
func (s *BusDexCexDebitCreditRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexDebitCreditRecordDeleteReq 功能删除请求参数
type BusDexCexDebitCreditRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexCexDebitCreditRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
