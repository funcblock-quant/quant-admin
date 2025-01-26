package client

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"quanta-admin/app/grpc/pool"
	"quanta-admin/app/grpc/proto/client/observe_service"
	"time"
)

func StartNewObserver(amberConfig *observe_service.AmberConfig, ammDexConfig *observe_service.DexConfig, arbitrageConfig *observe_service.ArbitrageConfig) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return "", fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return "", fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := observe_service.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &observe_service.StartRequest{
		AmberConfig:     amberConfig,
		DexConfig:       ammDexConfig,
		ArbitrageConfig: arbitrageConfig,
	}

	// 发送 gRPC 请求
	resp, err := c.StartObserver(ctx, req)
	if err != nil {
		return "", fmt.Errorf("启动 observer失败: %w", err)
	}

	// 处理响应
	instanceID := resp.GetInstanceId()
	fmt.Println("Observer 启动成功. Instance ID:", instanceID)
	return instanceID, nil
}

func StopObserver(observerID string) (err error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := observe_service.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &observe_service.StopRequest{
		InstanceId: &observerID,
	}

	// 发送 gRPC 请求
	_, err = c.StopObserver(ctx, req)
	if err != nil {
		return fmt.Errorf("暂停 observer失败: %w", err)
	}

	// 处理响应
	return nil
}

func ListObservers() (observerList []*observe_service.ObserverInfo, err error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := observe_service.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ListObservers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("获取 observer列表失败: %w", err)
	}
	return resp.Observers, nil

}

func GetObserverState(observerId string) (*observe_service.GetStateResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := observe_service.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observe_service.GetStateRequest{
		InstanceId: &observerId,
	}
	resp, err := c.GetObserverState(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取 observer状态失败: %w", err)
	}
	return resp, nil
}
