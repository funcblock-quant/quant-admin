package server

import (
	"context"
	"fmt"
	pb "quanta-admin/app/grpc/proto/server/quanta_admin_service"
)

type QuantaAdminServer struct {
	pb.UnimplementedQuantaAdminServer // 必须嵌入此结构体
}

// GetStrategyInstanceConfig 获取策略实例配置
func (s *QuantaAdminServer) GetStrategyInstanceConfig(ctx context.Context, req *pb.CommonGetRequest) (*pb.GetInstanceConfigResponse, error) {
	fmt.Printf("GetStrategyInstanceConfig, req %v\n", req)
	resp := &pb.GetInstanceConfigResponse{}
	resp.InstanceId = req.InstanceId
	resp.Configs = make([]*pb.InstanceConfig, 0)
	return resp, nil
}

// GetStartOrStopFlag 获取策略实例启停标识
func (s *QuantaAdminServer) GetStartOrStopFlag(ctx context.Context, req *pb.CommonGetRequest) (*pb.GetStartOrStopFlagResponse, error) {
	fmt.Printf("GetStartOrStopFlag, req %v\n", req)
	resp := &pb.GetStartOrStopFlagResponse{}
	resp.InstanceId = req.InstanceId
	isActive := true
	resp.Flag = &isActive
	return resp, nil
}
