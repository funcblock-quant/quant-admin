package pool

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	ext "quanta-admin/config"
	"sync"
	"time"
)

// GrpcServices 定义一个类型来存储 gRPC 连接池
type GrpcServices struct {
	pools sync.Map // 使用 sync.Map 存储连接池，支持并发安全
}

// 全局 gRPC 客户端连接池实例
var grpcPools GrpcServices

func InitGrpcPool() error {

	fmt.Printf("start init grpc pool with config: %+v\n", ext.ExtConfig.Grpc)

	tempPools := make(map[string]*Pool) // 使用临时 map 存储连接池

	for serviceName, address := range ext.ExtConfig.Grpc {
		fmt.Printf("start init grpc pool for server: %s with address %s\n", serviceName, address)
		// 复制操作在循环体内部，闭包外部！
		serviceNameCopy := serviceName
		addressCopy := address

		factory := func(ctx context.Context) (*grpc.ClientConn, error) {
			dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // 设置连接超时
			defer cancel()

			fmt.Printf("dialing to %s with context deadline\n", addressCopy)

			conn, err := grpc.DialContext(dialCtx, addressCopy, grpc.WithInsecure(), grpc.WithBlock()) // 使用 grpc.DialContext
			if err != nil {
				fmt.Printf("dialing to %s (%s) failed: %v, error type: %T, detail: %+v\n", serviceNameCopy, addressCopy, err, err, err)
				if st, ok := status.FromError(err); ok {
					fmt.Printf("gRPC error code: %v, message: %v, details: %+v\n", st.Code(), st.Message(), st.Details())
				}
				return nil, nil
			}
			fmt.Printf("dialed to %s successfully\n", addressCopy)
			return conn, nil
		}

		p, err := NewWithContext(context.Background(), address, factory, 2, 5, time.Second*10)
		if err != nil {
			fmt.Printf("create pool for %s failed: %v", serviceName, err)
			return nil // 返回错误
		}
		tempPools[serviceName] = p // 存储到临时 map 中
		fmt.Printf("gRPC pool for service '%s' created successfully.\n", serviceName)
	}

	// 将临时 map 的内容存储到全局 sync.Map 中，保证并发安全
	for k, v := range tempPools {
		grpcPools.pools.Store(k, v)
	}

	//健康检查
	for serviceName, p := range tempPools {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		client, err := p.Get(ctx)
		if err != nil {
			fmt.Printf("WARNING: health check for %s failed: %v. Service will still start.\n", serviceName, err)
			continue // 跳过此服务的检查
		}
		client.Close() // 立即释放连接
	}
	return nil
}

// GetGrpcClient 根据服务名获取 gRPC 客户端连接
func GetGrpcClient(serviceName string) (*ClientConn, error) {
	p, ok := grpcPools.pools.Load(serviceName)
	if !ok {
		return nil, fmt.Errorf("gRPC client pool for %s not found", serviceName)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := p.(*Pool).Get(ctx)
	if err != nil || client == nil {
		return nil, fmt.Errorf("get grpc client from pool for %s failed: %w", serviceName, err)
	}
	fmt.Printf("get gRPC client %s successfully.\n", serviceName)
	return client, nil
}

// CloseGrpcClients 关闭所有 gRPC 客户端连接池
func CloseGrpcClients() {
	grpcPools.pools.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Pool); ok {
			p.Close()
		}
		return true
	})
}
