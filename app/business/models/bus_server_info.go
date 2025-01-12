package models

import (
	"quanta-admin/common/models"
)

type BusServerInfo struct {
	models.Model

	ServerIp      string `json:"serverIp" gorm:"type:varchar(32);comment:ip"`
	Username      string `json:"username" gorm:"type:varchar(64);comment:用户名"`
	ConnectType   string `json:"connectType" gorm:"type:tinyint;comment:认证方式"`
	NetworkStatus string `json:"networkStatus" gorm:"type:tinyint;default:1;comment:网络健康"`
	Status        string `json:"status" gorm:"type:tinyint;default:0;comment:服务器状态"`
	CpuNum        string `json:"cpuNum" gorm:"type:int;comment:cpu核心数"`
	MemorySize    string `json:"memorySize" gorm:"type:int;comment:内存大小"`
	models.ModelTime
	models.ControlBy
}

func (BusServerInfo) TableName() string {
	return "bus_server_info"
}

func (e *BusServerInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *BusServerInfo) GetId() interface{} {
	return e.Id
}
