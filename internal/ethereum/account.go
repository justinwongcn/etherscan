package ethereum

import (
	"context"
	"fmt"
	"sync"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/justinwongcn/go-ethlibs/node"
	"golang.org/x/sync/errgroup"
)

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
