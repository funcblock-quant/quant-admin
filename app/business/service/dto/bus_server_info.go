package dto

import (
	"quanta-admin/app/business/models"
	"quanta-admin/common/dto"
	common "quanta-admin/common/models"
)

type BusServerInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	ServerIp       string `form:"serverIp"  search:"type:exact;column:server_ip;table:bus_server_info" comment:"ip"`
	NetworkStatus  string `form:"networkStatus"  search:"type:exact;column:network_status;table:bus_server_info" comment:"网络健康"`
	Status         string `form:"status"  search:"type:exact;column:status;table:bus_server_info" comment:"服务器状态"`
	CpuNum         string `form:"cpuNum"  search:"type:exact;column:cpu_num;table:bus_server_info" comment:"cpu核心数"`
	MemorySize     string `form:"memorySize"  search:"type:exact;column:memory_size;table:bus_server_info" comment:"内存大小"`
	BusServerInfoOrder
}

type BusServerInfoOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:bus_server_info"`
	ServerIp      string `form:"serverIpOrder"  search:"type:order;column:server_ip;table:bus_server_info"`
	Username      string `form:"usernameOrder"  search:"type:order;column:username;table:bus_server_info"`
	ConnectType   string `form:"connectTypeOrder"  search:"type:order;column:connect_type;table:bus_server_info"`
	NetworkStatus string `form:"networkStatusOrder"  search:"type:order;column:network_status;table:bus_server_info"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:bus_server_info"`
	CpuNum        string `form:"cpuNumOrder"  search:"type:order;column:cpu_num;table:bus_server_info"`
	MemorySize    string `form:"memorySizeOrder"  search:"type:order;column:memory_size;table:bus_server_info"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:bus_server_info"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:bus_server_info"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:bus_server_info"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:bus_server_info"`
	DeletedAt     string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:bus_server_info"`
}

func (m *BusServerInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BusServerInfoInsertReq struct {
	Id            int    `json:"-" comment:""` //
	ServerIp      string `json:"serverIp" comment:"ip"`
	Username      string `json:"username" comment:"用户名"`
	ConnectType   string `json:"connectType" comment:"认证方式"`
	NetworkStatus string `json:"networkStatus" comment:"网络健康"`
	Status        string `json:"status" comment:"服务器状态"`
	CpuNum        string `json:"cpuNum" comment:"cpu核心数"`
	MemorySize    string `json:"memorySize" comment:"内存大小"`
	common.ControlBy
}

func (s *BusServerInfoInsertReq) Generate(model *models.BusServerInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ServerIp = s.ServerIp
	model.Username = s.Username
	model.ConnectType = s.ConnectType
	model.NetworkStatus = s.NetworkStatus
	model.Status = s.Status
	model.CpuNum = s.CpuNum
	model.MemorySize = s.MemorySize
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *BusServerInfoInsertReq) GetId() interface{} {
	return s.Id
}

type BusServerInfoUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	ServerIp      string `json:"serverIp" comment:"ip"`
	Username      string `json:"username" comment:"用户名"`
	ConnectType   string `json:"connectType" comment:"认证方式"`
	NetworkStatus string `json:"networkStatus" comment:"网络健康"`
	Status        string `json:"status" comment:"服务器状态"`
	CpuNum        string `json:"cpuNum" comment:"cpu核心数"`
	MemorySize    string `json:"memorySize" comment:"内存大小"`
	common.ControlBy
}

func (s *BusServerInfoUpdateReq) Generate(model *models.BusServerInfo) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ServerIp = s.ServerIp
	model.Username = s.Username
	model.ConnectType = s.ConnectType
	model.NetworkStatus = s.NetworkStatus
	model.Status = s.Status
	model.CpuNum = s.CpuNum
	model.MemorySize = s.MemorySize
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *BusServerInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// BusServerInfoGetReq 功能获取请求参数
type BusServerInfoGetReq struct {
	Id int `uri:"id"`
}

func (s *BusServerInfoGetReq) GetId() interface{} {
	return s.Id
}

type BusServerStopReq struct {
	Id int `uri:"id"`
}

func (s *BusServerStopReq) GetId() interface{} {
	return s.Id
}

type BusServerStartReq struct {
	Id int `uri:"id"`
}

func (s *BusServerStartReq) GetId() interface{} {
	return s.Id
}

// BusServerInfoDeleteReq 功能删除请求参数
type BusServerInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BusServerInfoDeleteReq) GetId() interface{} {
	return s.Ids
}
