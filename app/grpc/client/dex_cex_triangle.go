package client

import (
	"context"
	"fmt"
	"quanta-admin/app/grpc/pool"
	"quanta-admin/app/grpc/proto/client/observer_service"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func StartNewArbitragerClient(instanceId *string, amberObserverConfig *observer_service.AmberObserverConfig, ammDexConfig *observer_service.DexConfig, observerParams *observer_service.ObserverParams) error {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &observer_service.StartRequest{
		InstanceId:  instanceId,
		AmberConfig: amberObserverConfig,
		DexConfig:   ammDexConfig,
		Params:      observerParams,
	}
	fmt.Printf("start observer req:%+v\n", req)

	// 发送 gRPC 请求
	_, err = c.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("启动 observer失败: %w", err)
	}

	// 处理响应
	fmt.Println("Observer 启动成功. Instance ID:", *instanceId)
	return nil
}

func StopArbitragerClient(instanceId string) (err error) {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构造请求消息
	req := &observer_service.InstantId{
		InstanceId: &instanceId,
	}

	// 发送 gRPC 请求
	_, err = c.Stop(ctx, req)
	if err != nil {
		return fmt.Errorf("暂停 observer失败: %w", err)
	}

	// 处理响应
	return nil
}

func ListArbitragerClient() (clientList []*observer_service.BasicInfo, err error) {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.List(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("获取 observer列表失败: %w", err)
	}
	return resp.Infos, nil

}

func GetObserverState(instanceId string) (*observer_service.GetStateResponse, error) {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.InstantId{
		InstanceId: &instanceId,
	}
	resp, err := c.GetObserverState(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取 observer状态失败: %w", err)
	}
	return resp, nil
}

func GetObserverParams(instanceId string) (*observer_service.ObserverParams, error) {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.InstantId{
		InstanceId: &instanceId,
	}
	resp, err := c.GetObserverParams(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取 observer参数失败: %w", err)
	}
	return resp, nil
}

func UpdateObserverParams(instanceId string, observerParams *observer_service.ObserverParams) error {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.UpdateObserverParamsRequest{
		InstanceId: &instanceId,
		Params:     observerParams,
	}
	_, err = c.UpdateObserverParams(ctx, req)
	if err != nil {
		return fmt.Errorf("更新 observer参数失败: %w", err)
	}
	return nil
}

func EnableTrader(instanceId string, cexConfig *observer_service.CexConfig, traderParams *observer_service.TraderParams, swapperConfig *observer_service.SwapperConfig) error {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.EnableTraderRequest{
		InstanceId:    &instanceId,
		CexConfig:     cexConfig,
		Params:        traderParams,
		SwapperConfig: swapperConfig,
	}
	_, err = c.EnableTrader(ctx, req)
	if err != nil {
		return fmt.Errorf("启动 Trader失败: %w", err)
	}
	return nil
}

func DisableTrader(instanceId string) error {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.InstantId{
		InstanceId: &instanceId,
	}
	_, err = c.DisableTrader(ctx, req)
	if err != nil {
		return fmt.Errorf("暂停 Trader失败: %w", err)
	}
	return nil
}

func GetTraderParams(instanceId string) (*observer_service.TraderParams, error) {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.InstantId{
		InstanceId: &instanceId,
	}
	resp, err := c.GetTraderParams(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取 trader参数失败: %w", err)
	}
	return resp, nil
}

func UpdateTraderParams(instanceId string, traderParams *observer_service.TraderParams) error {
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
	c := observer_service.NewArbitragerClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &observer_service.UpdateTraderParamsRequest{
		InstanceId: &instanceId,
		Params:     traderParams,
	}
	_, err = c.UpdateTraderParams(ctx, req)
	if err != nil {
		return fmt.Errorf("更新 trader参数失败: %w", err)
	}
	return nil
}
