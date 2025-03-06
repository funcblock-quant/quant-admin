package dto

import (
	"time"

	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusRiskEventGetPageReq struct {
	dto.Pagination     `search:"-"`
	StrategyId         string    `form:"strategyId"  search:"type:exact;column:strategy_id;table:bus_risk_event" comment:"策略id"`
	StrategyInstanceId string    `form:"strategyInstanceId"  search:"type:exact;column:strategy_instance_id;table:bus_risk_event" comment:"策略实例ID"`
	RiskScope          string    `form:"riskScope"  search:"type:exact;column:risk_scope;table:bus_risk_event" comment:"风控范围"`
	AssetSymbol        string    `form:"assetSymbol"  search:"type:exact;column:asset_symbol;table:bus_risk_event" comment:"风控币种"`
	RiskLevel          string    `form:"riskLevel"  search:"type:exact;column:risk_level;table:bus_risk_event" comment:"风控级别"`
	ManualRecover      string    `form:"manualRecover"  search:"type:exact;column:manual_recover;table:bus_risk_event" comment:"是否需要人工恢复"`
	AutoRecoverTime    time.Time `form:"autoRecoverTime"  search:"type:exact;column:auto_recover_time;table:bus_risk_event" comment:"自动恢复时间"`
	IsRecovered        string    `form:"isRecovered"  search:"type:exact;column:is_recovered;table:bus_risk_event" comment:"是否已恢复"`
	RecoveredBy        string    `form:"recoveredBy"  search:"type:exact;column:recovered_by;table:bus_risk_event" comment:"审核人"`
	RecoveredAt        time.Time `form:"recoveredAt"  search:"type:exact;column:recovered_at;table:bus_risk_event" comment:"恢复时间"`
	TriggerRule        string    `form:"triggerRule"  search:"type:exact;column:trigger_rule;table:bus_risk_event" comment:"触发的风控规则"`
	TriggerValue       string    `form:"triggerValue"  search:"type:exact;column:trigger_value;table:bus_risk_event" comment:"触发值"`
	BusRiskEventOrder
}

type BusRiskEventOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:bus_risk_event"`
	StrategyId         string `form:"strategyIdOrder"  search:"type:order;column:strategy_id;table:bus_risk_event"`
	StrategyInstanceId string `form:"strategyInstanceIdOrder"  search:"type:order;column:strategy_instance_id;table:bus_risk_event"`
	RiskScope          string `form:"riskScopeOrder"  search:"type:order;column:risk_scope;table:bus_risk_event"`
	AssetSymbol        string `form:"assetSymbolOrder"  search:"type:order;column:asset_symbol;table:bus_risk_event"`
	RiskLevel          string `form:"riskLevelOrder"  search:"type:order;column:risk_level;table:bus_risk_event"`
	ManualRecover      string `form:"manualRecoverOrder"  search:"type:order;column:manual_recover;table:bus_risk_event"`
	AutoRecoverTime    string `form:"autoRecoverTimeOrder"  search:"type:order;column:auto_recover_time;table:bus_risk_event"`
	IsRecovered        string `form:"isRecoveredOrder"  search:"type:order;column:is_recovered;table:bus_risk_event"`
	RecoveredBy        string `form:"recoveredByOrder"  search:"type:order;column:recovered_by;table:bus_risk_event"`
	RecoveredAt        string `form:"recoveredAtOrder"  search:"type:order;column:recovered_at;table:bus_risk_event"`
	TriggerRule        string `form:"triggerRuleOrder"  search:"type:order;column:trigger_rule;table:bus_risk_event"`
	TriggerValue       string `form:"triggerValueOrder"  search:"type:order;column:trigger_value;table:bus_risk_event"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_risk_event"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_risk_event"`
}

func (m *BusRiskEventGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusRiskEventInsertReq struct {
	Id                 int       `json:"-" comment:"风控事件ID"` // 风控事件ID
	StrategyId         string    `json:"strategyId" comment:"策略id"`
	StrategyInstanceId string    `json:"strategyInstanceId" comment:"策略实例ID"`
	RiskScope          int       `json:"riskScope" comment:"风控范围"`
	AssetSymbol        string    `json:"assetSymbol" comment:"风控币种"`
	RiskLevel          int       `json:"riskLevel" comment:"风控级别"`
	ManualRecover      int       `json:"manualRecover" comment:"是否需要人工恢复"`
	AutoRecoverTime    time.Time `json:"autoRecoverTime" comment:"自动恢复时间"`
	IsRecovered        int       `json:"isRecovered" comment:"是否已恢复"`
	RecoveredBy        string    `json:"recoveredBy" comment:"审核人"`
	RecoveredAt        time.Time `json:"recoveredAt" comment:"恢复时间"`
	TriggerRule        string    `json:"triggerRule" comment:"触发的风控规则"`
	TriggerValue       string    `json:"triggerValue" comment:"触发值"`
	common.ControlBy
}

func (s *BusRiskEventInsertReq) Generate(model *models.BusRiskEvent) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StrategyId = s.StrategyId
	model.StrategyInstanceId = s.StrategyInstanceId
	model.RiskScope = s.RiskScope
	model.AssetSymbol = s.AssetSymbol
	model.RiskLevel = s.RiskLevel
	model.ManualRecover = s.ManualRecover
	model.AutoRecoverTime = &s.AutoRecoverTime
	model.IsRecovered = s.IsRecovered
	model.RecoveredBy = s.RecoveredBy
	model.RecoveredAt = &s.RecoveredAt
	model.TriggerRule = s.TriggerRule
	model.TriggerValue = s.TriggerValue
}

func (s *BusRiskEventInsertReq) GetId() interface{} {
	return s.Id
}

type BusRiskEventUpdateReq struct {
	Id          int       `uri:"id" comment:"风控事件ID"` // 风控事件ID
	IsRecovered int       `json:"isRecovered" comment:"是否已恢复"`
	RecoveredAt time.Time `json:"recoveredAt" comment:"恢复时间"`
	common.ControlBy
}

func (s *BusRiskEventUpdateReq) Generate(model *models.BusRiskEvent) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.IsRecovered = s.IsRecovered
	model.RecoveredAt = &s.RecoveredAt
}

func (s *BusRiskEventUpdateReq) GetId() interface{} {
	return s.Id
}

// BusRiskEventGetReq 功能获取请求参数
type BusRiskEventGetReq struct {
	Id int `uri:"id"`
}

func (s *BusRiskEventGetReq) GetId() interface{} {
	return s.Id
}

// BusRiskEventDeleteReq 功能删除请求参数
type BusRiskEventDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusRiskEventDeleteReq) GetId() interface{} {
	return s.Ids
}
