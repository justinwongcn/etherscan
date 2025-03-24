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

// GetLatestBlockNumber 获取最新区块高度
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 最新区块高度
//   - error: 可能的错误：
//   - 节点连接错误
//   - 请求执行错误
func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.BlockNumber(ctx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// getDefaultNumberOrTag 处理区块号或标签的默认值
//
// Parameters:
//   - numberOrTag: string 区块号或标签，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//   - 空字符串 - 将被转换为"latest"
//
// Returns:
//   - string: 处理后的区块号或标签
func getDefaultNumberOrTag(numberOrTag string) string {
	if numberOrTag == "" {
		return "latest"
	}
	return numberOrTag
}

// GetBalance 获取指定地址的账户余额
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
//   - uint64: 账户余额（单位：wei）
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

	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
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
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
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
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
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

// GetBlockTransactionCountByHash 获取指定区块哈希的交易数量
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - blockHash: string 区块哈希（32字节的十六进制字符串）
//
// Returns:
//   - uint64: 该区块中的交易数量
//   - error: 可能的错误：
//   - 无效的区块哈希格式
//   - 节点连接错误
func (c *Client) GetBlockTransactionCountByHash(ctx context.Context, blockHash string) (uint64, error) {
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetBlockTransactionCountByHash(ctx, blockHash)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// GetBlockTransactionCountByNumber 获取指定区块号的交易数量
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//
// Returns:
//   - uint64: 该区块中的交易数量
//   - error: 可能的错误：
//   - 无效的区块号格式
//   - 节点连接错误
func (c *Client) GetBlockTransactionCountByNumber(ctx context.Context, numberOrTag string) (uint64, error) {
	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return 0, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetBlockTransactionCountByNumber(ctx, *numOrTag)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// GetCode 获取指定地址的合约代码
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - address: string 要查询的合约地址
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//
// Returns:
//   - string: 合约代码（十六进制格式）
//   - error: 可能的错误：
//   - 无效的地址格式
//   - 无效的区块号格式
//   - 节点连接错误
func (c *Client) GetCode(ctx context.Context, address string, numberOrTag string) (string, error) {
	// 验证并转换地址格式
	addr, err := eth.NewAddress(address)
	if err != nil {
		return "", fmt.Errorf("invalid ethereum address: %v", err)
	}

	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return "", fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetCode(ctx, *addr, *numOrTag)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// SendRawTransaction 发送已签名的交易数据
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - signedTxData: string 已签名的交易数据（十六进制格式）
//
// Returns:
//   - string: 交易哈希（32字节的十六进制字符串）
//   - error: 可能的错误：
//   - 无效的交易数据格式
//   - 节点连接错误
func (c *Client) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	// 验证交易数据格式
	if len(signedTxData) < 2 || signedTxData[:2] != "0x" {
		return "", fmt.Errorf("invalid transaction data format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.SendRawTransaction(ctx, signedTxData)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Call 执行以太坊智能合约的只读调用
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - from: string 可选，交易发送方地址
//   - to: string 必需，交易接收方地址
//   - gas: uint64 可选，交易执行的gas限制
//   - gasPrice: uint64 可选，每单位gas的价格
//   - value: uint64 可选，随交易发送的以太币数量
//   - data: string 可选，方法签名和编码参数的哈希
//   - numberOrTag: string 区块号或标签，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//
// Returns:
//   - string: 合约执行的返回值
//   - error: 可能的错误：
//   - 无效的地址格式
//   - 无效的区块号格式
//   - 节点连接错误
func (c *Client) Call(ctx context.Context, from, to string, gas, gasPrice, value uint64, data string, numberOrTag string) (string, error) {
	// 验证接收方地址格式
	toAddr, err := eth.NewAddress(to)
	if err != nil {
		return "", fmt.Errorf("invalid to address: %v", err)
	}

	// 验证发送方地址格式（如果提供）
	var fromAddr *eth.Address
	if from != "" {
		addr, err := eth.NewAddress(from)
		if err != nil {
			return "", fmt.Errorf("invalid from address: %v", err)
		}
		fromAddr = addr
	}

	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return "", fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 创建交易对象
	tx := eth.Transaction{
		From:     *fromAddr,
		To:       toAddr,
		Gas:      *eth.MustQuantity(fmt.Sprintf("0x%x", gas)),
		GasPrice: eth.MustQuantity(fmt.Sprintf("0x%x", gasPrice)),
		Value:    *eth.MustQuantity(fmt.Sprintf("0x%x", value)),
		Input:    eth.Input(data),
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.Call(ctx, tx, *numOrTag)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
