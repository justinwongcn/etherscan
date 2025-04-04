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
		To:       toAddr,
		Gas:      *eth.MustQuantity(fmt.Sprintf("0x%x", gas)),
		GasPrice: eth.MustQuantity(fmt.Sprintf("0x%x", gasPrice)),
		Value:    *eth.MustQuantity(fmt.Sprintf("0x%x", value)),
		Input:    eth.Input(data),
	}

	if fromAddr != nil {
		tx.From = *fromAddr
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

// EstimateGas 估算交易所需的gas数量
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - from: string 可选，交易发送方地址
//   - to: string 可选，交易接收方地址
//   - gas: uint64 可选，交易执行的gas限制
//   - gasPrice: uint64 可选，每单位gas的价格
//   - value: uint64 可选，随交易发送的以太币数量
//   - data: string 可选，方法签名和编码参数的哈希
//
// Returns:
//   - uint64: 预估的gas数量，注意：返回的估算值可能会显著高于实际使用量
//   - error: 可能的错误：
//   - 无效的地址格式
//   - 节点连接错误
func (c *Client) EstimateGas(ctx context.Context, from, to string, gas, gasPrice, value uint64, data string) (uint64, error) {
	var fromAddr, toAddr *eth.Address
	var err error

	if from != "" {
		addr, err := eth.NewAddress(from)
		if err != nil {
			return 0, fmt.Errorf("invalid from address: %v", err)
		}
		fromAddr = addr
	}

	if to != "" {
		addr, err := eth.NewAddress(to)
		if err != nil {
			return 0, fmt.Errorf("invalid to address: %v", err)
		}
		toAddr = addr
	}

	// 创建交易对象
	tx := eth.Transaction{
		Input: eth.Input(data),
	}

	// 设置可选参数
	if fromAddr != nil {
		tx.From = *fromAddr
	}
	if toAddr != nil {
		tx.To = toAddr
	}
	if gas > 0 {
		tx.Gas = *eth.MustQuantity(fmt.Sprintf("0x%x", gas))
	}
	if gasPrice > 0 {
		tx.GasPrice = eth.MustQuantity(fmt.Sprintf("0x%x", gasPrice))
	}
	if value > 0 {
		tx.Value = *eth.MustQuantity(fmt.Sprintf("0x%x", value))
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.EstimateGas(ctx, tx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}
