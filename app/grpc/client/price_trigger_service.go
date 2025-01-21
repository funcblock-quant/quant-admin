package client

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"quanta-admin/app/grpc/pool"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	"time"
)

func StartInstance(request *trigger_service.StartTriggerRequest) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return "", fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	resp, err := c.StartInstance(ctx, request)
	if err != nil {
		return "", fmt.Errorf("启动 trigger_server失败: %w", err)
	}

	// 处理响应
	instanceID := resp.GetInstanceId()
	fmt.Println("trigger 启动成功. Instance ID:", instanceID)
	return instanceID, nil
}

func ListInstances() ([]string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	resp, err := c.ListInstances(ctx, &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("获取到运行中的instanceIds失败: %w", err)
	}

	// 处理响应
	instanceIds := resp.GetInstanceIds()
	fmt.Println("获取到运行中的instanceIds:", instanceIds)
	return instanceIds, nil
}

func CheckApiKeyHealth(request *trigger_service.APIConfig) (bool, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return false, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	resp, err := c.CheckApiKey(ctx, request)
	if err != nil {
		return false, fmt.Errorf("校验apikey健康性失败: %w", err)
	}

	// 处理响应
	isHealthy := resp.GetIsHealth()
	fmt.Println("校验apikey健康性结果:", isHealthy)
	return isHealthy, nil
}
