package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusArbitrageRecordGetPageReq struct {
	dto.Pagination `search:"-"`
	ArbitrageId    string `form:"arbitrageId"  search:"type:contains;column:arbitrage_id;table:bus_arbitrage_record" comment:"套利记录id"`
	StrategyName   string `form:"strategyName"  search:"type:contains;column:strategy_name;table:bus_arbitrage_record" comment:"策略名称"`
	ContractType   string `form:"contractType"  search:"type:exact;column:contract_type;table:bus_arbitrage_record" comment:"合约类型"`
	BeginTime      string `form:"beginTime"  search:"type:gte;column:start_time;table:bus_arbitrage_record" comment:"开始时间"`
	EndTime        string `form:"endTime"  search:"type:lte;column:start_time;table:bus_arbitrage_record" comment:"结束时间"`
	BusArbitrageRecordOrder
}

type BusArbitrageRecordOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_arbitrage_record"`
	ArbitrageId        string `form:"arbitrageIdOrder"  search:"type:order;column:arbitrage_id;table:bus_arbitrage_record"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_arbitrage_record"`
	StrategyName       string `form:"strategyNameOrder"  search:"type:order;column:strategy_name;table:bus_arbitrage_record"`
	Type               string `form:"typeOrder"  search:"type:order;column:type;table:bus_arbitrage_record"`
	ContractType       string `form:"contractTypeOrder"  search:"type:order;column:contract_type;table:bus_arbitrage_record"`
	RealizedPnl        string `form:"realizedPnlOrder"  search:"type:order;column:realized_pnl;table:bus_arbitrage_record"`
	UnrealizedPnl      string `form:"unrealizedPnlOrder"  search:"type:order;column:unrealized_pnl;table:bus_arbitrage_record"`
	ExpectPnl          string `form:"expectPnlOrder"  search:"type:order;column:expect_pnl;table:bus_arbitrage_record"`
	ExpectPnlPercent   string `form:"expectPnlPercentOrder"  search:"type:order;column:expect_pnl_percent;table:bus_arbitrage_record"`
	StartTime          string `form:"startTimeOrder"  search:"type:order;column:start_time;table:bus_arbitrage_record"`
	EndTime            string `form:"endTimeOrder"  search:"type:order;column:end_time;table:bus_arbitrage_record"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:bus_arbitrage_record"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_arbitrage_record"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_arbitrage_record"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_arbitrage_record"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_arbitrage_record"`
}

func (m *BusArbitrageRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusArbitrageRecordInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	ArbitrageId        string `json:"arbitrageId" comment:"套利记录id"`
	StrategyInstanceId string `json:"strategyInstanceId" comment:"套利单所属策略id"`
	StrategyName       string `json:"strategyName" comment:"策略名称"`
	Type               string `json:"type" comment:"套利类型, 0-模拟盘观测, 1-实盘套利"`
	ContractType       string `json:"contractType" comment:"合约类型"`
	RealizedPnl        string `json:"realizedPnl" comment:"已实现盈亏"`
	UnrealizedPnl      string `json:"unrealizedPnl" comment:"未实现盈亏"`
	ExpectPnl          string `json:"expectPnl" comment:"预期盈亏"`
	ExpectPnlPercent   string `json:"expectPnlPercent" comment:"预期收益率"`
	StartTime          int64  `json:"startTime" comment:"套利单开始时间,单位：ns"`
	EndTime            int64  `json:"endTime" comment:"套利单结束时间,单位：ns"`
	Status             string `json:"status" comment:"套利单状态, 0-套利中，1-套利完成, 2-对账完成"`
	common.ControlBy
}

func (s *BusArbitrageRecordInsertReq) Generate(model *models.BusArbitrageRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.StrategyInstanceId = s.StrategyInstanceId
	model.StrategyName = s.StrategyName
	model.Type = s.Type
	model.ContractType = s.ContractType
	model.RealizedPnl = s.RealizedPnl
	model.UnrealizedPnl = s.UnrealizedPnl
	model.ExpectPnl = s.ExpectPnl
	model.ExpectPnlPercent = s.ExpectPnlPercent
	model.StartTime = s.StartTime
	model.EndTime = s.EndTime
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusArbitrageRecordInsertReq) GetId() interface{} {
	return s.Id
}

type BusArbitrageRecordUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	ArbitrageId        string `json:"arbitrageId" comment:"套利记录id"`
	StrategyInstanceId string `json:"strategyInstanceId" comment:"套利单所属策略id"`
	StrategyName       string `json:"strategyName" comment:"策略名称"`
	Type               string `json:"type" comment:"套利类型, 0-模拟盘观测, 1-实盘套利"`
	ContractType       string `json:"contractType" comment:"合约类型"`
	RealizedPnl        string `json:"realizedPnl" comment:"已实现盈亏"`
	UnrealizedPnl      string `json:"unrealizedPnl" comment:"未实现盈亏"`
	ExpectPnl          string `json:"expectPnl" comment:"预期盈亏"`
	ExpectPnlPercent   string `json:"expectPnlPercent" comment:"预期收益率"`
	StartTime          int64  `json:"startTime" comment:"套利单开始时间,单位：ns"`
	EndTime            int64  `json:"endTime" comment:"套利单结束时间,单位：ns"`
	Status             string `json:"status" comment:"套利单状态, 0-套利中，1-套利完成, 2-对账完成"`
	common.ControlBy
}

func (s *BusArbitrageRecordUpdateReq) Generate(model *models.BusArbitrageRecord) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ArbitrageId = s.ArbitrageId
	model.StrategyInstanceId = s.StrategyInstanceId
	model.StrategyName = s.StrategyName
	model.Type = s.Type
	model.ContractType = s.ContractType
	model.RealizedPnl = s.RealizedPnl
	model.UnrealizedPnl = s.UnrealizedPnl
	model.ExpectPnl = s.ExpectPnl
	model.ExpectPnlPercent = s.ExpectPnlPercent
	model.StartTime = s.StartTime
	model.EndTime = s.EndTime
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusArbitrageRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// BusArbitrageRecordGetReq 功能获取请求参数
type BusArbitrageRecordGetReq struct {
	Id int `uri:"id"`
}

func (s *BusArbitrageRecordGetReq) GetId() interface{} {
	return s.Id
}

// BusArbitrageRecordDeleteReq 功能删除请求参数
type BusArbitrageRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusArbitrageRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
