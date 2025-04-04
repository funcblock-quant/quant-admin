package client

import (
	"context"
	"fmt"
	"quanta-admin/app/grpc/pool"
	"quanta-admin/app/grpc/proto/client/trigger_service"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

func StartTriggerInstance(request *trigger_service.StartTriggerRequest) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return "", fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return "", fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	fmt.Println("创建TriggerInstanceClient实例")
	c := trigger_service.NewTriggerInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	fmt.Println("开始请求trigger server to start instance, req:", request)
	resp, err := c.StartInstance(ctx, request)
	if err != nil {
		fmt.Println("启动 trigger_server失败: %w", err)
		return "", fmt.Errorf("启动 trigger_server失败: %w", err)
	}

	// 处理响应
	instanceID := resp.GetInstanceId()
	fmt.Println("trigger 启动成功. Instance ID:", instanceID)
	return instanceID, nil
}

func StopTriggerInstance(request *trigger_service.StopTriggerRequest) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	_, err = c.StopInstance(ctx, request)
	if err != nil {
		return fmt.Errorf("暂停 trigger_server失败: %w", err)
	}

	// 处理响应
	fmt.Println("trigger 暂停成功. Instance ID:", request.InstanceId)
	return nil
}

func ListInstances() ([]string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	fmt.Printf("clientConn: %+v\n", clientConn)
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}

	defer func() {
		if clientConn != nil {
			clientConn.Close()
		}
	}() // 确保连接在使用后返回连接池

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
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return false, fmt.Errorf("grpc客户端连接为空")
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

func UpdateProfitTarget(request *trigger_service.ProfitTargetConfig) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池
	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.UpdateProfitTargetConfig(ctx, request)
	if err != nil {
		return fmt.Errorf("更新 止盈设置失败: %w", err)
	}

	// 处理响应
	fmt.Println("更新trigger 止盈设置成功. Instance ID:", request.InstanceId)
	return nil
}

func UpdateExecuteConfig(request *trigger_service.ExecuteConfig) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("trigger-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池
	// 创建 gRPC 客户端实例
	c := trigger_service.NewTriggerInstanceClient(clientConn)
	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.UpdateExecuteConfig(ctx, request)
	if err != nil {
		return fmt.Errorf("更新 执行次数失败: %w", err)
	}

	// 处理响应
	fmt.Println("更新trigger 执行次数成功. Instance ID:", request.InstanceId)
	return nil
}
