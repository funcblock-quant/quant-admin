package dto

import (
	"time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusStrategyInstanceGetPageReq struct {
	dto.Pagination `search:"-"`
	StrategyId     string    `form:"strategyId"  search:"type:exact;column:strategy_id;table:bus_strategy_instance" comment:"策略id"`
	AccountGroupId string    `form:"accountGroupId"  search:"type:exact;column:account_group_id;table:bus_strategy_instance" comment:"账户组id"`
	StartRunTime   time.Time `form:"startRunTime"  search:"type:lte;column:start_run_time;table:bus_strategy_instance" comment:"启动时间"`
	BusStrategyInstanceOrder
}

type BusStrategyInstanceOrder struct {
	Id              string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_instance"`
	StrategyId      string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_strategy_instance"`
	AccountGroupId  string `form:"accountGroupIdOrder"  search:"type:order;column:account_group_id;table:bus_strategy_instance"`
	ExchangeId1     string `form:"exchangeId1Order"  search:"type:order;column:exchange_id1;table:bus_strategy_instance"`
	ExchangeId1Type string `form:"exchangeId1TypeOrder"  search:"type:order;column:exchange_id1_type;table:bus_strategy_instance"`
	ExchangeId2     string `form:"exchangeId2Order"  search:"type:order;column:exchange_id2;table:bus_strategy_instance"`
	ExchangeId2Type string `form:"exchangeId2TypeOrder"  search:"type:order;column:exchange_id2_type;table:bus_strategy_instance"`
	InstanceName    string `form:"instanceNameOrder"  search:"type:order;column:instance_name;table:bus_strategy_instance"`
	StartRunTime    string `form:"startRunTimeOrder"  search:"type:order;column:start_run_time;table:bus_strategy_instance"`
	StopRunTime     string `form:"stopRunTimeOrder"  search:"type:order;column:stop_run_time;table:bus_strategy_instance"`
	ServerIp        string `form:"serverIpOrder"  search:"type:order;column:server_ip;table:bus_strategy_instance"`
	ServerName      string `form:"serverNameOrder"  search:"type:order;column:server_name;table:bus_strategy_instance"`
	Status          string `form:"statusOrder"  search:"type:order;column:status;table:bus_strategy_instance"`
	CreateBy        string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_instance"`
	UpdateBy        string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_instance"`
	CreatedAt       string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_instance"`
	UpdatedAt       string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_instance"`
	IsDeleted       string `form:"isDeletedOrder"  search:"type:order;column:is_deleted;table:bus_strategy_instance"`
	DeletedAt       string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_instance"`
}

func (m *BusStrategyInstanceGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyInstanceInsertReq struct {
	Id              int       `json:"-" comment:""` //
	StrategyId      string    `json:"strategyId" comment:"策略id"`
	AccountGroupId  string    `json:"accountGroupId" comment:"账户组id"`
	ExchangeId1     string    `json:"exchangeId1" comment:"交易所id1"`
	ExchangeId1Type string    `json:"exchangeId1Type" comment:"平台类型"`
	ExchangeId2     string    `json:"exchangeId2" comment:"交易所id2"`
	ExchangeId2Type string    `json:"exchangeId2Type" comment:"平台类型"`
	InstanceName    string    `json:"instanceName" comment:"策略实例名称"`
	StartRunTime    time.Time `json:"startRunTime" comment:"启动时间"`
	StopRunTime     time.Time `json:"stopRunTime" comment:"停止时间"`
	ServerIp        string    `json:"serverIp" comment:"服务器ip"`
	ServerName      string    `json:"serverName" comment:"服务器用户名"`
	Status          string    `json:"status" comment:"运行状态"`
	IsDeleted       string    `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyInstanceInsertReq) Generate(model *models.BusStrategyInstance) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.AccountGroupId = s.AccountGroupId
	model.ExchangeId1 = s.ExchangeId1
	model.ExchangeId1Type = s.ExchangeId1Type
	model.ExchangeId2 = s.ExchangeId2
	model.ExchangeId2Type = s.ExchangeId2Type
	model.InstanceName = s.InstanceName
	model.StartRunTime = s.StartRunTime
	model.StopRunTime = s.StopRunTime
	model.ServerIp = s.ServerIp
	model.ServerName = s.ServerName
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.IsDeleted = s.IsDeleted
}

func (s *BusStrategyInstanceInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyInstanceUpdateReq struct {
	Id              int       `uri:"id" comment:""` //
	StrategyId      string    `json:"strategyId" comment:"策略id"`
	AccountGroupId  string    `json:"accountGroupId" comment:"账户组id"`
	ExchangeId1     string    `json:"exchangeId1" comment:"交易所id1"`
	ExchangeId1Type string    `json:"exchangeId1Type" comment:"平台类型"`
	ExchangeId2     string    `json:"exchangeId2" comment:"交易所id2"`
	ExchangeId2Type string    `json:"exchangeId2Type" comment:"平台类型"`
	InstanceName    string    `json:"instanceName" comment:"策略实例名称"`
	StartRunTime    time.Time `json:"startRunTime" comment:"启动时间"`
	StopRunTime     time.Time `json:"stopRunTime" comment:"停止时间"`
	ServerIp        string    `json:"serverIp" comment:"服务器ip"`
	ServerName      string    `json:"serverName" comment:"服务器用户名"`
	Status          string    `json:"status" comment:"运行状态"`
	IsDeleted       string    `json:"isDeleted" comment:"删除标识位"`
	common.ControlBy
}

func (s *BusStrategyInstanceUpdateReq) Generate(model *models.BusStrategyInstance) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.AccountGroupId = s.AccountGroupId
	model.ExchangeId1 = s.ExchangeId1
	model.ExchangeId1Type = s.ExchangeId1Type
	model.ExchangeId2 = s.ExchangeId2
	model.ExchangeId2Type = s.ExchangeId2Type
	model.InstanceName = s.InstanceName
	model.StartRunTime = s.StartRunTime
	model.StopRunTime = s.StopRunTime
	model.ServerIp = s.ServerIp
	model.ServerName = s.ServerName
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.IsDeleted = s.IsDeleted
}

func (s *BusStrategyInstanceUpdateReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceGetReq 功能获取请求参数
type BusStrategyInstanceGetReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyInstanceGetReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceDeleteReq 功能删除请求参数
type BusStrategyInstanceDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusStrategyInstanceDeleteReq) GetId() interface{} {
	return s.Ids
}
