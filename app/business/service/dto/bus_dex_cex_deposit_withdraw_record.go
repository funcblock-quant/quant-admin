package dto

import (
    "time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusDexCexDepositWithdrawRecordGetPageReq struct {
	dto.Pagination     `search:"-"`
    OrderType string `form:"orderType"  search:"type:exact;column:order_type;table:bus_dex_cex_deposit_withdraw_record" comment:"订单类型"`
    FromAddress string `form:"fromAddress"  search:"type:exact;column:from_address;table:bus_dex_cex_deposit_withdraw_record" comment:"from地址"`
    ToAddress string `form:"toAddress"  search:"type:exact;column:to_address;table:bus_dex_cex_deposit_withdraw_record" comment:"to地址"`
    CexUid string `form:"cexUid"  search:"type:exact;column:cex_uid;table:bus_dex_cex_deposit_withdraw_record" comment:"cex账户uid"`
    CexOrderId string `form:"cexOrderId"  search:"type:exact;column:cex_order_id;table:bus_dex_cex_deposit_withdraw_record" comment:"交易所充提业务id"`
    TxHash string `form:"txHash"  search:"type:exact;column:tx_hash;table:bus_dex_cex_deposit_withdraw_record" comment:"链上交易hash"`
    Symbol string `form:"symbol"  search:"type:exact;column:symbol;table:bus_dex_cex_deposit_withdraw_record" comment:"充提币种"`
    BusDexCexDepositWithdrawRecordOrder
}

type BusDexCexDepositWithdrawRecordOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:bus_dex_cex_deposit_withdraw_record"`
    OrderType string `form:"orderTypeOrder"  search:"type:order;column:order_type;table:bus_dex_cex_deposit_withdraw_record"`
    FromAddress string `form:"fromAddressOrder"  search:"type:order;column:from_address;table:bus_dex_cex_deposit_withdraw_record"`
    ToAddress string `form:"toAddressOrder"  search:"type:order;column:to_address;table:bus_dex_cex_deposit_withdraw_record"`
    CexUid string `form:"cexUidOrder"  search:"type:order;column:cex_uid;table:bus_dex_cex_deposit_withdraw_record"`
    CexOrderId string `form:"cexOrderIdOrder"  search:"type:order;column:cex_order_id;table:bus_dex_cex_deposit_withdraw_record"`
    TxHash string `form:"txHashOrder"  search:"type:order;column:tx_hash;table:bus_dex_cex_deposit_withdraw_record"`
    Symbol string `form:"symbolOrder"  search:"type:order;column:symbol;table:bus_dex_cex_deposit_withdraw_record"`
    Fee string `form:"feeOrder"  search:"type:order;column:fee;table:bus_dex_cex_deposit_withdraw_record"`
    FeeAsset string `form:"feeAssetOrder"  search:"type:order;column:fee_asset;table:bus_dex_cex_deposit_withdraw_record"`
    StartTime string `form:"startTimeOrder"  search:"type:order;column:start_time;table:bus_dex_cex_deposit_withdraw_record"`
    FinishTime string `form:"finishTimeOrder"  search:"type:order;column:finish_time;table:bus_dex_cex_deposit_withdraw_record"`
    Status string `form:"statusOrder"  search:"type:order;column:status;table:bus_dex_cex_deposit_withdraw_record"`
    Remark string `form:"remarkOrder"  search:"type:order;column:remark;table:bus_dex_cex_deposit_withdraw_record"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_dex_cex_deposit_withdraw_record"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_dex_cex_deposit_withdraw_record"`
    
}

func (m *BusDexCexDepositWithdrawRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusDexCexDepositWithdrawRecordInsertReq struct {
    Id int `json:"-" comment:"主键 ID"` // 主键 ID
    OrderType string `json:"orderType" comment:"订单类型"`
    FromAddress string `json:"fromAddress" comment:"from地址"`
    ToAddress string `json:"toAddress" comment:"to地址"`
    CexUid string `json:"cexUid" comment:"cex账户uid"`
    CexOrderId string `json:"cexOrderId" comment:"交易所充提业务id"`
    TxHash string `json:"txHash" comment:"链上交易hash"`
    Symbol string `json:"symbol" comment:"充提币种"`
    Fee string `json:"fee" comment:"手续费"`
    FeeAsset string `json:"feeAsset" comment:"手续费币种"`
    StartTime time.Time `json:"startTime" comment:"发起时间"`
    FinishTime time.Time `json:"finishTime" comment:"完成时间"`
    Status string `json:"status" comment:"状态"`
    Remark string `json:"remark" comment:"备注"`
    common.ControlBy
}

func (s *BusDexCexDepositWithdrawRecordInsertReq) Generate(model *models.BusDexCexDepositWithdrawRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.OrderType = s.OrderType
    model.FromAddress = s.FromAddress
    model.ToAddress = s.ToAddress
    model.CexUid = s.CexUid
    model.CexOrderId = s.CexOrderId
    model.TxHash = s.TxHash
    model.Symbol = s.Symbol
    model.Fee = s.Fee
    model.FeeAsset = s.FeeAsset
    model.StartTime = s.StartTime
    model.FinishTime = s.FinishTime
    model.Status = s.Status
    model.Remark = s.Remark
}

func (s *BusDexCexDepositWithdrawRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusDexCexDepositWithdrawRecordUpdateReq struct {
    Id int `uri:"id" comment:"主键 ID"` // 主键 ID
    OrderType string `json:"orderType" comment:"订单类型"`
    FromAddress string `json:"fromAddress" comment:"from地址"`
    ToAddress string `json:"toAddress" comment:"to地址"`
    CexUid string `json:"cexUid" comment:"cex账户uid"`
    CexOrderId string `json:"cexOrderId" comment:"交易所充提业务id"`
    TxHash string `json:"txHash" comment:"链上交易hash"`
    Symbol string `json:"symbol" comment:"充提币种"`
    Fee string `json:"fee" comment:"手续费"`
    FeeAsset string `json:"feeAsset" comment:"手续费币种"`
    StartTime time.Time `json:"startTime" comment:"发起时间"`
    FinishTime time.Time `json:"finishTime" comment:"完成时间"`
    Status string `json:"status" comment:"状态"`
    Remark string `json:"remark" comment:"备注"`
    common.ControlBy
}

func (s *BusDexCexDepositWithdrawRecordUpdateReq) Generate(model *models.BusDexCexDepositWithdrawRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.OrderType = s.OrderType
    model.FromAddress = s.FromAddress
    model.ToAddress = s.ToAddress
    model.CexUid = s.CexUid
    model.CexOrderId = s.CexOrderId
    model.TxHash = s.TxHash
    model.Symbol = s.Symbol
    model.Fee = s.Fee
    model.FeeAsset = s.FeeAsset
    model.StartTime = s.StartTime
    model.FinishTime = s.FinishTime
    model.Status = s.Status
    model.Remark = s.Remark
}

func (s *BusDexCexDepositWithdrawRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusDexCexDepositWithdrawRecordGetReq 功能获取请求参数
type BusDexCexDepositWithdrawRecordGetReq struct {
     Id int `uri:"id"`
}
func (s *BusDexCexDepositWithdrawRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusDexCexDepositWithdrawRecordDeleteReq 功能删除请求参数
type BusDexCexDepositWithdrawRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusDexCexDepositWithdrawRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
