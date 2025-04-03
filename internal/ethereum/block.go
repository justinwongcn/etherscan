package ethereum

import (
	"context"
	"fmt"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/justinwongcn/go-ethlibs/node"
)

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

// GetBlockByHash 获取指定区块哈希的区块信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - blockHash: string 区块哈希（32字节的十六进制字符串）
//   - fullTx: bool 如果为true则返回完整的交易对象，否则仅返回交易哈希
//
// Returns:
//   - *eth.Block: 区块信息，包含区块头、交易等数据
//   - error: 可能的错误：
//   - 无效的区块哈希格式
//   - 节点连接错误
//   - 区块不存在
func (c *Client) GetBlockByHash(ctx context.Context, blockHash string, fullTx bool) (*eth.Block, error) {
	// 验证区块哈希格式
	if len(blockHash) < 2 || blockHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid block hash format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.BlockByHash(ctx, blockHash, fullTx)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Block), nil
}

// GetBlockByNumber 获取指定区块号的区块信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//   - fullTx: bool 如果为true则返回完整的交易对象，否则仅返回交易哈希
//
// Returns:
//   - *eth.Block: 区块信息，包含区块头、交易等数据
//   - error: 可能的错误：
//   - 无效的区块号格式
//   - 节点连接错误
//   - 区块不存在
func (c *Client) GetBlockByNumber(ctx context.Context, numberOrTag string, fullTx bool) (*eth.Block, error) {
	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return nil, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.BlockByNumber(ctx, *numOrTag, fullTx)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Block), nil
}

// GetUncleByBlockHashAndIndex 获取指定区块哈希和叔块索引的叔块信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - blockHash: string 区块哈希（32字节的十六进制字符串）
//   - index: uint64 叔块的索引位置
//
// Returns:
//   - *eth.Block: 叔块信息，结构与普通区块相同
//   - error: 可能的错误：
//   - 无效的区块哈希格式
//   - 节点连接错误
func (c *Client) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*eth.Block, error) {
	// 验证区块哈希格式
	if len(blockHash) < 2 || blockHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid block hash format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetUncleByBlockHashAndIndex(ctx, blockHash, index)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Block), nil
}

// GetUncleByBlockNumberAndIndex 通过区块号和叔块索引获取叔块信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//   - index: uint64 叔块的索引位置
//
// Returns:
//   - *eth.Block: 叔块信息，包含区块头等数据（不包含交易信息）
//   - error: 可能的错误：
//   - 无效的区块号格式
//   - 节点连接错误
//   - 叔块不存在
func (c *Client) GetUncleByBlockNumberAndIndex(ctx context.Context, numberOrTag string, index uint64) (*eth.Block, error) {
	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return nil, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetUncleByBlockNumberAndIndex(ctx, *numOrTag, index)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Block), nil
}
