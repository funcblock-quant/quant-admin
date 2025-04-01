package client

import (
	"context"
	"fmt"
	"quanta-admin/app/grpc/pool"
	pb "quanta-admin/app/grpc/proto/client/water_level_service"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func StartWaterLevelInstance(request *pb.StartInstanceRequest) (string, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return "", fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return "", fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	fmt.Println("创建WaterLevelInstanceClient实例")
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	// fmt.Printf("开始请求water level server to start instance, req: %v \n", request)
	resp, err := c.StartInstance(ctx, request)
	if err != nil {
		fmt.Println("启动 water_level_server实例失败: %w", err)
		return "", fmt.Errorf("启动 water_level_server实例失败: %w", err)
	}

	// 处理响应
	instanceID := resp.GetInstanceId()
	fmt.Println("water level实例启动成功. Instance ID:", instanceID)
	return instanceID, nil
}

func StopWaterLevelInstance(instanceId *pb.InstanceId) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	_, err = c.StopInstance(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("暂停 water_level 实例失败: %w", err)
	}

	// 处理响应
	fmt.Println("water_level 实例暂停成功. Instance ID:", instanceId.InstanceId)
	return nil
}

func UpdateWaterLevelInstance(request *pb.UpdateInstanceParamsRequest) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	_, err = c.UpdateInstanceParams(ctx, request)
	if err != nil {
		return fmt.Errorf("更新 water_level 实例失败: %w", err)
	}

	// 处理响应
	fmt.Println("water_level 实例更新成功. Instance ID:", request.InstanceId)
	return nil
}

func GetWaterLevelInstanceState(request *pb.InstanceId) (*pb.GetStateResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 发送 gRPC 请求
	resp, err := c.GetInstanceState(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("获取 water_level 实例数据失败: %w", err)
	}

	// 处理响应
	fmt.Printf("获取water_level 实例数据成功. Resp: %v\r\n", resp)
	return resp, nil
}

func ListWaterLevelInstance() (*pb.InstanceListResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ListInstances(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("获取 instance列表失败: %w", err)
	}
	return resp, nil

}

func GetPortfolioUnwindingInfo(request *pb.PortfolioUnwindingRequest) (*pb.GetPortfolioUnwindingResponse, error) {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return nil, fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return nil, fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.GetPortfolioUnwindingInfo(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("获取 portfolio unwinding 信息失败: %w", err)
	}
	return resp, nil
}

func PortfolioUnwinding(request *pb.PortfolioUnwindingRequest) error {
	// 获取 gRPC 客户端连接
	clientConn, err := pool.GetGrpcClient("water-level-service")
	if err != nil {
		return fmt.Errorf("获取grpc客户端失败: %w", err)
	}
	if clientConn == nil || clientConn.ClientConn == nil { // 再次检查 clientConn 是否为 nil
		return fmt.Errorf("grpc客户端连接为空")
	}
	defer clientConn.Close() // 确保连接在使用后返回连接池

	// 创建 gRPC 客户端实例
	c := pb.NewInstanceClient(clientConn)

	// 设置超时 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.PortfolioUnwinding(ctx, request)
	if err != nil {
		return fmt.Errorf("portfolio unwinding 失败: %w", err)
	}
	fmt.Printf("portfolio unwinding 成功. Resp: %v\r\n", resp)
	return nil
}
