package models

import (
	"quanta-admin/common/models"
)

type BusCommonConfig struct {
	models.Model

	Category   string `json:"category" gorm:"type:varchar(50);comment:业务类型"`
	RecordId   string `json:"recordId" gorm:"default:0;type:bigint;comment:具体业务记录"`
	ConfigKey  string `json:"configKey" gorm:"type:varchar(100);comment:配置项 key"`
	ConfigJson string `json:"configJson" gorm:"type:json;comment:配置信息"`
	models.ModelTime
	models.ControlBy
}

func (BusCommonConfig) TableName() string {
	return "bus_common_config"
}

func (e *BusCommonConfig) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusCommonConfig) GetId() interface{} {
	return e.Id
}
