// Package ethereum 提供以太坊客户端的实现，包括：
//   - 基础的以太坊节点交互功能
//   - 连接池管理
//   - 并发请求处理
//   - 错误处理和重试机制
package ethereum

import (
	"context"
	"fmt"
	"sync"

	"github.com/justinwongcn/go-ethlibs/eth"
	"golang.org/x/sync/errgroup"

	"github.com/justinwongcn/go-ethlibs/node"
)

// NewClient 创建一个新的以太坊客户端
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - nodeURL: string 以太坊节点的 URL 地址
//   - opts: *ClientOptions 客户端配置选项，如果为 nil 则使用默认配置
//
// Returns:
//   - *Client: 初始化后的以太坊客户端实例
//   - error: 可能的错误
func NewClient(ctx context.Context, nodeURL string, opts *ClientOptions) (*Client, error) {
	client, err := node.NewClient(ctx, nodeURL)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		opts = DefaultClientOptions()
	}

	// 初始化连接池配置
	c := &Client{
		nodeClient:   client,
		nodeURL:      nodeURL,
		maxConns:     opts.MaxConns,
		idleTimeout:  opts.IdleTimeout,
		healthCheck:  opts.HealthCheck,
		maxIdleConns: opts.MaxIdleConns,
	}

	// 初始化连接池
	c.connPool = &sync.Pool{
		New: func() any {
			// 创建新的节点客户端连接
			newClient, err := node.NewClient(ctx, nodeURL)
			if err != nil {
				return nil
			}
			return newClient
		},
	}

	// 启动连接池管理协程
	go c.managePool(ctx)

	return c, nil
}

// withConnection 在连接池中执行操作的通用辅助函数
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - fn: func(node.Client) (any, error) 在连接上执行的操作函数，返回值会被类型断言为具体类型
//
// Returns:
//   - any: 操作的返回值，需要由调用者进行类型断言
//   - error: 可能的错误：
//   - 获取连接失败
//   - 操作执行失败
func (c *Client) withConnection(ctx context.Context, fn func(node.Client) (any, error)) (any, error) {
	// 从连接池获取连接
	conn, err := c.getConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %v", err)
	}
	defer c.releaseConnection(conn)

	// 执行操作
	return fn(conn)
}

// GasPrice 获取当前 gas 价格
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 当前 gas 价格（单位：wei）
//   - error: 操作过程中可能发生的错误
func (c *Client) GasPrice(ctx context.Context) (uint64, error) {
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GasPrice(ctx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.BlockNumber(ctx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

func (c *Client) GetBalance(ctx context.Context, address string, numberOrTag string) (uint64, error) {
	// 验证并转换地址格式
	addr, err := eth.NewAddress(address)
	if err != nil {
		return 0, fmt.Errorf("invalid ethereum address: %v", err)
	}

	// 验证并转换区块号格式
	numOrTag := eth.MustBlockNumberOrTag(numberOrTag)
	if numOrTag == nil {
		return 0, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetBalance(ctx, *addr, *numOrTag)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// GetBalances 批量获取多个地址的账户余额
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - addresses: []string 要查询余额的账户地址列表
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//   - maxAddresses: int 单次查询最多支持的地址数量，默认值为5
//
// Returns:
//   - map[string]uint64: 地址到余额的映射，以wei为单位
//   - error: 可能的错误：
//   - 地址列表为空
//   - 地址数量超过限制
//   - 无效的地址格式
//   - 无效的区块号格式
//   - 节点连接错误
func (c *Client) GetBalances(ctx context.Context, addresses []string, numberOrTag string, maxAddresses ...int) (map[string]uint64, error) {
	// 验证地址列表
	if len(addresses) == 0 {
		return nil, fmt.Errorf("address list is empty")
	}

	// 设置最大地址数量限制
	maxAddr := 5 // 默认值
	if len(maxAddresses) > 0 && maxAddresses[0] > 0 {
		maxAddr = maxAddresses[0]
	}

	if len(addresses) > maxAddr {
		return nil, fmt.Errorf("too many addresses: %d (max: %d)", len(addresses), maxAddr)
	}

	// 验证并转换区块号格式
	numOrTag := eth.MustBlockNumberOrTag(numberOrTag)
	if numOrTag == nil {
		return nil, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 创建结果映射
	result := make(map[string]uint64)

	// 使用 errgroup 进行并发请求
	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex

	// 并发获取每个地址的余额
	for _, address := range addresses {
		addr := address // 创建副本以避免闭包问题
		g.Go(func() error {
			// 从连接池获取连接
			conn, err := c.getConnection(ctx)
			if err != nil {
				return fmt.Errorf("failed to get connection: %v", err)
			}
			defer c.releaseConnection(conn)

			// 验证并转换地址格式
			ethAddr, err := eth.NewAddress(addr)
			if err != nil {
				return fmt.Errorf("invalid ethereum address: %v", err)
			}

			// 调用节点接口获取余额
			balance, err := conn.GetBalance(ctx, *ethAddr, *numOrTag)
			if err != nil {
				return fmt.Errorf("failed to get balance for %s: %v", addr, err)
			}

			// 线程安全地更新结果映射
			mu.Lock()
			result[addr] = balance
			mu.Unlock()
			return nil
		})
	}

	// 等待所有请求完成
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTransactionCount 获取指定地址发送的交易数量
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - address: string 要查询的账户地址
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//
// Returns:
//   - uint64: 该地址发送的交易数量
//   - error: 可能的错误：
//   - 无效的地址格式
//   - 无效的区块号格式
//   - 节点连接错误
func (c *Client) GetTransactionCount(ctx context.Context, address string, numberOrTag string) (uint64, error) {
	// 验证并转换地址格式
	addr, err := eth.NewAddress(address)
	if err != nil {
		return 0, fmt.Errorf("invalid ethereum address: %v", err)
	}

	// 验证并转换区块号格式
	numOrTag := eth.MustBlockNumberOrTag(numberOrTag)
	if numOrTag == nil {
		return 0, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetTransactionCount(ctx, *addr, *numOrTag)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}
