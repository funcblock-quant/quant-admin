package dto

import (
	"strconv"
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
	Id             string `form:"idOrder"  search:"type:order;column:id;table:bus_strategy_instance"`
	StrategyId     string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_strategy_instance"`
	AccountGroupId string `form:"accountGroupIdOrder"  search:"type:order;column:account_group_id;table:bus_strategy_instance"`
	ExchangeId1    string `form:"exchangeId1Order"  search:"type:order;column:exchange_id1;table:bus_strategy_instance"`
	Exchange1Type  string `form:"exchange1TypeOrder"  search:"type:order;column:exchange_id1_type;table:bus_strategy_instance"`
	ExchangeId2    string `form:"exchangeId2Order"  search:"type:order;column:exchange_id2;table:bus_strategy_instance"`
	Exchange2Type  string `form:"exchange2TypeOrder"  search:"type:order;column:exchange_id2_type;table:bus_strategy_instance"`
	InstanceName   string `form:"instanceNameOrder"  search:"type:order;column:instance_name;table:bus_strategy_instance"`
	StartRunTime   string `form:"startRunTimeOrder"  search:"type:order;column:start_run_time;table:bus_strategy_instance"`
	StopRunTime    string `form:"stopRunTimeOrder"  search:"type:order;column:stop_run_time;table:bus_strategy_instance"`
	ServerId       string `form:"serverIdOrder"  search:"type:order;column:server_id;table:bus_strategy_instance"`
	Status         string `form:"statusOrder"  search:"type:order;column:status;table:bus_strategy_instance"`
	CreateBy       string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_strategy_instance"`
	UpdateBy       string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_strategy_instance"`
	CreatedAt      string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_strategy_instance"`
	UpdatedAt      string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_strategy_instance"`
	DeletedAt      string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_strategy_instance"`
}

func (m *BusStrategyInstanceGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusStrategyInstanceGetPageResp struct {
	Id             string     `json:"id"`
	StrategyId     string     `json:"strategyId"`
	StrategyName   string     `json:"strategyName"`
	GroupName      string     `json:"accountGroupName"`
	AccountGroupId string     `json:"accountGroupId"`
	ExchangeId1    string     `json:"exchangeId1"`
	ExchangeName1  string     `json:"exchange1Name"`
	Exchange1Type  string     `json:"exchange1Type"`
	ExchangeId2    string     `json:"exchangeId2"`
	ExchangeName2  string     `json:"exchange2Name"`
	Exchange2Type  string     `json:"exchange2Type"`
	InstanceName   string     `json:"instanceName"`
	StartRunTime   *time.Time `json:"startRunTime" gorm:"default:NULL"`
	StopRunTime    *time.Time `json:"stopRunTime" gorm:"default:NULL"`
	ServerId       string     `json:"serverId"`
	Status         string     `json:"status"`
}

type BusStrategyInstanceGetResp struct {
	Id             string                             `json:"id"`
	StrategyId     string                             `json:"strategyId"`
	AccountGroupId string                             `json:"accountGroupId"`
	ExchangeId1    string                             `json:"exchangeId1"`
	ExchangeName1  string                             `json:"exchange1Name"`
	Exchange1Type  string                             `json:"exchange1Type"`
	ExchangeId2    string                             `json:"exchangeId2"`
	ExchangeName2  string                             `json:"exchange2Name"`
	Exchange2Type  string                             `json:"exchange2Type"`
	InstanceName   string                             `json:"instanceName"`
	StartRunTime   *time.Time                         `json:"startRunTime" gorm:"default:NULL"`
	StopRunTime    *time.Time                         `json:"stopRunTime" gorm:"default:NULL"`
	ServerId       string                             `json:"serverId"`
	Status         string                             `json:"status"`
	Configs        []models.BusStrategyInstanceConfig `json:"configs"`
}

type BusStrategyInstanceInsertReq struct {
	Id             int                                  `json:"-" comment:""` //
	StrategyId     string                               `json:"strategyId" comment:"策略id"`
	AccountGroupId string                               `json:"accountGroupId" comment:"账户组id"`
	ExchangeId1    string                               `json:"exchangeId1" comment:"交易所id1"`
	Exchange1Name  string                               `json:"exchange1Name" comment:"交易所1名称"`
	Exchange1Type  string                               `json:"exchange1Type" comment:"平台类型"`
	ExchangeId2    string                               `json:"exchangeId2" comment:"交易所id2"`
	Exchange2Name  string                               `json:"exchange2Name" comment:"交易所2名称"`
	Exchange2Type  string                               `json:"exchange2Type" comment:"平台类型"`
	InstanceName   string                               `json:"instanceName" comment:"策略实例名称"`
	ServerId       string                               `json:"serverId" comment:"服务器id"`
	ServerName     string                               `json:"serverName" comment:"服务器用户名"`
	Status         string                               `json:"status" comment:"运行状态"`
	Configurations []BusStrategyInstanceConfigInsertReq `json:"configurations" comment:""`
	common.ControlBy
}

func (s *BusStrategyInstanceInsertReq) Generate(model *models.BusStrategyInstance) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.AccountGroupId = s.AccountGroupId
	model.ExchangeId1 = s.ExchangeId1
	model.Exchange1Type = s.Exchange1Type
	model.ExchangeName1 = s.Exchange1Name
	model.ExchangeId2 = s.ExchangeId2
	model.Exchange2Type = s.Exchange2Type
	model.ExchangeName2 = s.Exchange2Name
	model.InstanceName = s.InstanceName
	model.ServerId = s.ServerId
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusStrategyInstanceInsertReq) GetId() interface{} {
	return s.Id
}

type BusStrategyInstanceUpdateReq struct {
	Id             string                               `uri:"id" comment:""` //
	StrategyId     string                               `json:"strategyId" comment:"策略id"`
	AccountGroupId string                               `json:"accountGroupId" comment:"账户组id"`
	ExchangeId1    string                               `json:"exchangeId1" comment:"交易所id1"`
	Exchange1Type  string                               `json:"exchange1Type" comment:"平台类型"`
	Exchange1Name  string                               `json:"exchange1Name" comment:"交易所1名称"`
	ExchangeId2    string                               `json:"exchangeId2" comment:"交易所id2"`
	Exchange2Type  string                               `json:"exchange2Type" comment:"平台类型"`
	Exchange2Name  string                               `json:"exchange2Name" comment:"交易所2名称"`
	InstanceName   string                               `json:"instanceName" comment:"策略实例名称"`
	ServerId       string                               `json:"serverId" comment:"服务器id"`
	ServerName     string                               `json:"serverName" comment:"服务器用户名"`
	Status         string                               `json:"status" comment:"运行状态"`
	Configurations []BusStrategyInstanceConfigInsertReq `json:"configurations" comment:""`
	common.ControlBy
}

func (s *BusStrategyInstanceUpdateReq) Generate(model *models.BusStrategyInstance) {
	if s.Id == "0" {
		id, _ := strconv.Atoi(s.Id)
		model.Model = common.Model{Id: id}
	}
	model.StrategyId = s.StrategyId
	model.AccountGroupId = s.AccountGroupId
	model.ExchangeId1 = s.ExchangeId1
	model.Exchange1Type = s.Exchange1Type
	model.ExchangeName1 = s.Exchange1Name
	model.ExchangeId2 = s.ExchangeId2
	model.Exchange2Type = s.Exchange2Type
	model.ExchangeName2 = s.Exchange2Name
	model.InstanceName = s.InstanceName
	model.ServerId = s.ServerId
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
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

// BusStrategyInstanceStartReq 功能获取请求参数
type BusStrategyInstanceStartReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyInstanceStartReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceStopReq 功能获取请求参数
type BusStrategyInstanceStopReq struct {
	Id int `uri:"id"`
}

func (s *BusStrategyInstanceStopReq) GetId() interface{} {
	return s.Id
}

// BusStrategyInstanceDeleteReq 功能删除请求参数
type BusStrategyInstanceDeleteReq struct {
	Ids []string `json:"ids"`
}

func (s *BusStrategyInstanceDeleteReq) GetId() interface{} {
	return s.Ids
}

// BusStrategyInstanceDashboardGetReq 获取实例dashboard数据请求参数
type BusStrategyInstanceDashboardGetReq struct {
	StrategyInstanceIds []int `json:"instanceIds"`
}

func (s *BusStrategyInstanceDashboardGetReq) GetId() interface{} {
	return s.StrategyInstanceIds
}

// BusStrategyInstanceDashboardGetResp 策略dashboard返回体
type BusStrategyInstanceDashboardGetResp struct {
	BusStrategyInstanceBalanceBody
	BusStrategyInstanceStatisticsInfo
	SymbolPnlRankChart        []BusStrategyInstanceChartInfo `json:"symbolPnlRankChart"`
	StrategyInstanceChartInfo []BusStrategyInstanceChartInfo `json:"strategyInstanceChartInfo"`
}

// BusStrategyInstanceBalanceBody 策略实例资金信息
type BusStrategyInstanceBalanceBody struct {
	BeginBalance        string                         `json:"beginBalance"`  // 启动时策略资产
	TotalBalance        float64                        `json:"totalBalance"`  // 总资产
	RealisedPnl         float64                        `json:"realisedPnl"`   // 已实现收益
	UnrealisedPnl       float64                        `json:"unrealisedPnl"` // 未实现收益
	TotalPnl            float64                        `json:"totalPnl"`      //总收益
	DailyBalanceChart   []BusStrategyInstanceChartInfo `json:"dailyBalanceChart"`
	WeeklyBalanceChart  []BusStrategyInstanceChartInfo `json:"weeklyBalanceChart"`
	MonthlyBalanceChart []BusStrategyInstanceChartInfo `json:"monthlyBalanceChart"`
}

// BusStrategyInstanceStatisticsInfo 策略实例统计信息
type BusStrategyInstanceStatisticsInfo struct {
	TotalArbitrageNum     int                            `json:"totalArbitrageNum"` // 总套利次数
	WinNum                int                            `json:"winNum"`            // 胜次数
	LossNum               int                            `json:"lossNum"`           // 败次数
	WinRate               int                            `json:"winRate"`           // 胜率
	DailyArbitrageChart   []BusStrategyInstanceChartInfo `json:"arbitrageList"`
	WeeklyArbitrageChart  []BusStrategyInstanceChartInfo `json:"weeklyList"`
	MonthlyArbitrageChart []BusStrategyInstanceChartInfo `json:"monthlyList"`
}

// BusStrategyInstanceChartInfo 策略实例图表通用结构体
type BusStrategyInstanceChartInfo struct {
	Xcoordinate string `json:"xcoordinate"` // 横坐标
	Ycoordinate string `json:"ycoordinate"` // 纵坐标
}
