package client

import (
	"context"
	"fmt"
	"quanta-admin/app/grpc/pool"
	pb "quanta-admin/app/grpc/proto/observe_service"
	"time"
)

func StartNewObserver(amberConfig *pb.AmberConfig, ammDexConfig *pb.AmmDexConfig, arbitrageConfig *pb.ArbitrageConfig) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return "", fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &pb.StartRequest{
		AmberConfig:     amberConfig,
		AmmDexConfig:    ammDexConfig,
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
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &pb.StopRequest{
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

func ListObservers() {

}

func GetObserverState(observerId string) (*pb.GetStateResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("solana-observer")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewObserverClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetStateRequest{
		InstanceId: &observerId,
	}
	resp, err := c.GetObserverState(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取 observer状态失败: %w", err)
	}
	return resp, nil
}
