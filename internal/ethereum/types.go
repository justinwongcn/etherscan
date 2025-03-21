// Package ethereum 提供以太坊客户端的类型定义
package ethereum

import (
	"sync"
	"time"

	"github.com/justinwongcn/go-ethlibs/node"
)

// Client 封装以太坊客户端，提供以下功能：
//   - 与以太坊节点的基础交互
//   - 连接池管理
//   - 自动的健康检查
//   - 并发请求处理
type Client struct {
	nodeClient node.Client // 底层节点客户端
	nodeURL    string      // 节点URL，用于创建新连接
	// 连接池配置
	maxConns     int           // 最大并发连接数
	idleTimeout  time.Duration // 空闲连接的超时时间
	healthCheck  bool          // 是否启用连接健康检查
	connPool     *sync.Pool    // 连接池，用于复用连接
	connCount    int32         // 当前活跃连接数
	maxIdleConns int           // 最大空闲连接数，超过此数量的空闲连接将被关闭
}

// ClientOptions 定义客户端的配置选项，用于在创建客户端时自定义连接池行为
type ClientOptions struct {
	MaxConns     int           // 最大并发连接数，控制资源使用
	IdleTimeout  time.Duration // 空闲连接的超时时间，超时后连接将被清理
	HealthCheck  bool          // 是否启用连接健康检查，启用后将定期检查连接状态
	MaxIdleConns int           // 最大空闲连接数，用于限制连接池大小
}

// DefaultClientOptions 返回默认的客户端配置选项
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		MaxConns:     100,         // 默认最大连接数
		IdleTimeout:  time.Minute, // 默认空闲超时时间
		HealthCheck:  true,        // 默认启用健康检查
		MaxIdleConns: 10,          // 默认最大空闲连接数
	}
}
