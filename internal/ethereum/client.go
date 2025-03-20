package ethereum

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/justinwongcn/go-ethlibs/eth"
	"golang.org/x/sync/errgroup"

	"github.com/justinwongcn/go-ethlibs/node"
)

// Client 封装以太坊客户端
type Client struct {
	nodeClient node.Client
	nodeURL    string // 节点URL
	// 连接池配置
	maxConns     int           // 最大连接数
	idleTimeout  time.Duration // 空闲连接超时时间
	healthCheck  bool          // 是否启用健康检查
	connPool     *sync.Pool    // 连接池
	connCount    int32         // 当前连接数
	maxIdleConns int           // 最大空闲连接数
}

// NewClient 创建一个新的以太坊客户端
func NewClient(ctx context.Context, nodeURL string) (*Client, error) {
	client, err := node.NewClient(ctx, nodeURL)
	if err != nil {
		return nil, err
	}

	// 初始化连接池配置
	c := &Client{
		nodeClient:   client,
		nodeURL:      nodeURL,
		maxConns:     100,         // 默认最大连接数
		idleTimeout:  time.Minute, // 默认空闲超时时间
		healthCheck:  true,        // 默认启用健康检查
		maxIdleConns: 10,          // 默认最大空闲连接数
	}

	// 初始化连接池
	c.connPool = &sync.Pool{
		New: func() interface{} {
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

// GasPrice 获取当前 gas 价格
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 当前 gas 价格（单位：wei）
//   - error: 操作过程中可能发生的错误

func (c *Client) GasPrice(ctx context.Context) (uint64, error) {
	return c.nodeClient.GasPrice(ctx)
}

// GetLatestBlockNumber 获取最新区块高度
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 当前客户端所在的最新区块高度
//   - error: 操作过程中可能发生的错误
func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.nodeClient.BlockNumber(ctx)
}

// GetBalance 获取指定地址的账户余额
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - address: string 要查询余额的账户地址（20字节）
//   - blockNumber: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//
// Returns:
//   - uint64: 账户余额，以wei为单位
//   - error: 可能的错误：
//   - 无效的地址格式
//   - 无效的区块号格式
//   - 节点连接错误
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

	// 从连接池获取连接
	conn, err := c.getConnection(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %v", err)
	}
	defer c.releaseConnection(conn)

	// 调用节点接口获取余额
	return conn.GetBalance(ctx, *addr, *numOrTag)
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
