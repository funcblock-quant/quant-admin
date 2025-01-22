package client

import (
	"context"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"quanta-admin/app/grpc/pool"
	"quanta-admin/app/grpc/proto/client/instance_service"
	"time"
)

func StartNewInstance(serviceName string, instanceId string, instanceType instance_service.InstanceType, config *string) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient(serviceName)
	if err != nil {
		return "", fmt.Errorf("获取 %s grpc客户端失败: %w", serviceName, err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := instance_service.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &instance_service.StartInstanceRequest{
		InstanceId:   instanceId,
		InstanceType: instanceType,
		ConfigYaml:   *config,
	}

	// 发送 gRPC 请求
	resp, err := c.StartInstance(ctx, request)
	if err != nil {
		return "", fmt.Errorf("%s 启动 instance：%s 失败: %w", serviceName, instanceId, err)
	}

	// 处理响应
	instanceID := resp.GetInstanceId()
	log.Infof("%s 服务启动 Instance 成功. Instance ID: %s", serviceName, instanceID)
	return instanceID, nil
}

func StopInstance(serviceName string, instanceId string) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient(serviceName)
	if err != nil {
		return fmt.Errorf("获取 %s grpc客户端失败: %w", serviceName, err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := instance_service.NewInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &instance_service.StopInstanceRequest{
		InstanceId: instanceId,
	}

	// 发送 gRPC 请求
	_, err = c.StopInstance(ctx, request)
	if err != nil {
		return fmt.Errorf("%s 停用 instance: %s 失败: %w", serviceName, instanceId, err)
	}

	log.Infof("%s 服务停用 Instance 成功. Instance ID: %s", serviceName, instanceId)
	return nil
}

func ListInstance(serviceName string) (*instance_service.InstanceListResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient(serviceName)
	if err != nil {
		return nil, fmt.Errorf("获取 %s grpc客户端失败: %w", serviceName, err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := instance_service.NewInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ListInstances(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("获取 Instance 列表失败: %w", err)
	}
	return resp, nil
}

func GetInstanceRealtimeInfo(serviceName string, instanceId string) (*instance_service.GetRealtimeInfoResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient(serviceName)
	if err != nil {
		return nil, fmt.Errorf("获取 %s grpc客户端失败: %w", serviceName, err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := instance_service.NewInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &instance_service.GetRealtimeInfoRequest{
		InstanceId: instanceId,
	}
	resp, err := c.GetInstanceRealtimeInfo(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("获取 Instance实时监控状态失败: %w", err)
	}
	return resp, nil
}
