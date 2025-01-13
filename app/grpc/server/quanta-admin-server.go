package server

import (
	"context"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	sdk "github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"
	"quanta-admin/app/business/daos"
	"quanta-admin/app/business/models"
	"quanta-admin/app/business/service/dto"
	pb "quanta-admin/app/grpc/proto/server/quanta_admin_service"
	"strconv"
)

type QuantaAdminServer struct {
	pb.UnimplementedQuantaAdminServer // 必须嵌入此结构体
	db                                *gorm.DB
}

func NewQuantaAdminServer() *QuantaAdminServer {
	db := sdk.Runtime.GetDbByKey("*") //默认都是获取 key 为*的数据库db连接
	return &QuantaAdminServer{
		db: db,
	}
}

// GetStrategyInstanceConfig 获取策略实例配置
func (s *QuantaAdminServer) GetStrategyInstanceConfig(ctx context.Context, req *pb.CommonGetRequest) (*pb.GetInstanceConfigResponse, error) {
	fmt.Printf("GetStrategyInstanceConfig, req %v\n", req)
	resp := &pb.GetInstanceConfigResponse{}
	resp.InstanceId = req.InstanceId

	list := make([]models.BusStrategyInstanceConfig, 0)

	service := daos.BusStrategyInstanceConfigDAO{
		Db: s.db,
	}

	request := &dto.BusStrategyInstanceConfigGetByInstanceIdReq{
		StrategyInstanceId: req.InstanceId,
	}

	service.GetInstanceConfigs(request, &list)

	log.Infof("GetStrategyInstanceConfig, req %v, list %v", req, list)
	configs := make([]*pb.InstanceConfig, 0)
	for _, item := range list {
		config := pb.InstanceConfig{
			ParamKey:   item.ParamKey,
			ParamValue: item.ParamValue,
		}
		configs = append(configs, &config)
	}
	resp.Configs = configs
	log.Infof("configs %+v", configs)

	return resp, nil
}

// GetStartOrStopFlag 获取策略实例启停标识
func (s *QuantaAdminServer) GetStartOrStopFlag(ctx context.Context, req *pb.CommonGetRequest) (*pb.GetStartOrStopStatusResponse, error) {
	fmt.Printf("GetStartOrStopFlag, req %v\n", req)
	resp := &pb.GetStartOrStopStatusResponse{}
	resp.InstanceId = req.InstanceId

	service := daos.BusStrategyInstanceDAO{
		Db: s.db,
	}

	instanceId, err := strconv.Atoi(req.InstanceId)
	if err != nil {
		return nil, err
	}
	request := &dto.BusStrategyInstanceGetReq{
		Id: instanceId,
	}
	data := &models.BusStrategyInstance{}
	err = service.GetInstanceStartStopFlag(request, data)
	log.Infof("GetStartOrStopFlag, data %v, status: %s", data, data.Status)
	if err != nil {
		return nil, err
	}

	resp.Status = data.Status
	log.Infof("GetStartOrStopStatus, resp %+v, resp.status %d", resp, resp.Status)

	return resp, nil
}
