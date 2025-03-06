package models

// COMMON_CONFIG Category
const (
	// 水位调节
	WATER_LEVEL string = "WATER_LEVEL"
	// dex-cex三角套利风控参数
	DEX_CEX_RISK_COTROL string = "DEX_CEX_RISK_COTROL"
)

// GlobalConfigKey 全局配置config key
const (
	// Solana水位调节参数
	GLOBAL_SOLANA_WATER_LEVEL_KEY string = "SOLANA_WATER_LEVEL"
	// 稳定币水位调节参数
	GLOBAL_STABLE_COIN_WATER_LEVEL_KEY string = "STABLE_COIN_WATER_LEVEL"
	// dex-cex三角套利风控参数
	RISK_CONTROL_CONFIG_KEY string = "RISK_CONTROL_CONFIG"
)

// 非标准化策略id
const (
	// dex-cex三角套利策略id
	STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE string = "DEX_CEX_TRIANGULAR_ARBITRAGE"
)

// 策略对应交易表名称
const (
	// dex-cex三角套利交易表
	STRATEGY_DEX_CEX_TRIANGULAR_ARBITRAGE_TRADES_TABLE string = "strategy_dex_cex_triangular_arbitrage_trades"
)

// 策略对应的业务类型
const (
	// dex-cex三角套利业务类型
	BUSINESS_TYPE_DEX_CEX_TRIANGULAR_ARBITRAGE string = "DEX_CEX链上链下三角套利"
)

// 风控进度表状态
const (
	// dex-cex三角套利业务类型
	RISK_CHECK_STATUS_NOT_STARTED int = 0
	RISK_CHECK_STATUS_PROCESSING  int = 1
	RISK_CHECK_STATUS_FINISHED    int = 2
)

// 风控范围
const (
	RISK_SCOPE_GOLBAL       int = 0 // 全局
	RISK_SCOPE_SINGLE_TOKEN int = 1 // 单币种
)

// 风控等级
const (
	RISK_LEVEL_LOW    int = 1 // 低,预警
	RISK_LEVEL_MIDDLE int = 2 // 中,关闭当前交易
	RISK_LEVEL_HIGH   int = 3 // 高,关闭全局交易
)

// 风控触发规则
const (
	TRIGGER_RULE_ABSOLUTE_LOSS_THRESHOLD string = "单笔交易亏损超阈值"
	TRIGGER_RULE_RELATIVE_LOSS_THRESHOLD string = "单笔交易亏损比例超阈值"
)
