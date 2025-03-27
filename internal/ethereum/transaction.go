package ethereum

import (
	"context"
	"fmt"
	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/justinwongcn/go-ethlibs/node"
)

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

// GetTransactionByHash 获取指定交易哈希的交易信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - txHash: string 交易哈希（32字节的十六进制字符串）
//
// Returns:
//   - *eth.Transaction: 交易信息，包含交易哈希、区块信息、发送方和接收方地址、交易值、gas相关参数等
//   - error: 可能的错误：
//   - 无效的交易哈希格式
//   - 节点连接错误
//   - 交易不存在
func (c *Client) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	// 验证交易哈希格式
	if len(txHash) < 2 || txHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid transaction hash format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.TransactionByHash(ctx, txHash)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Transaction), nil
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

// GetTransactionByBlockHashAndIndex 通过区块哈希和交易索引获取交易信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - blockHash: string 区块哈希（32字节的十六进制字符串）
//   - index: uint64 交易在区块中的索引位置
//
// Returns:
//   - *eth.Transaction: 交易信息，包含交易哈希、区块信息、发送方和接收方地址、交易值、gas相关参数等
//   - error: 可能的错误：
//   - 无效的区块哈希格式
//   - 节点连接错误
//   - 交易不存在
func (c *Client) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*eth.Transaction, error) {
	// 验证区块哈希格式
	if len(blockHash) < 2 || blockHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid block hash format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetTransactionByBlockHashAndIndex(ctx, blockHash, index)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Transaction), nil
}

// GetTransactionByBlockNumberAndIndex 通过区块号和交易索引获取交易信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - numberOrTag: string 区块号，可以是以下格式：
//   - 十六进制字符串（如"0x1"）表示具体区块号
//   - "latest" - 最新区块（默认）
//   - "earliest" - 创世区块
//   - "pending" - 待处理区块
//   - index: uint64 交易在区块中的索引位置
//
// Returns:
//   - *eth.Transaction: 交易信息，包含交易哈希、区块信息、发送方和接收方地址、交易值、gas相关参数等
//   - error: 可能的错误：
//   - 无效的区块号格式
//   - 节点连接错误
//   - 交易不存在
func (c *Client) GetTransactionByBlockNumberAndIndex(ctx context.Context, numberOrTag string, index uint64) (*eth.Transaction, error) {
	// 处理默认值并验证区块号格式
	numOrTag := eth.MustBlockNumberOrTag(getDefaultNumberOrTag(numberOrTag))
	if numOrTag == nil {
		return nil, fmt.Errorf("invalid block number or tag: %s", numberOrTag)
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.GetTransactionByBlockNumberAndIndex(ctx, *numOrTag, index)
	})
	if err != nil {
		return nil, err
	}

	return result.(*eth.Transaction), nil
}

// GetTransactionReceipt 获取交易收据信息
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - txHash: string 交易哈希（32字节的十六进制字符串）
//
// Returns:
//   - *eth.TransactionReceipt: 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
//   - error: 可能的错误：
//   - 无效的交易哈希格式
//   - 节点连接错误
//   - 交易收据不存在
func (c *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	// 验证交易哈希格式
	if len(txHash) < 2 || txHash[:2] != "0x" {
		return nil, fmt.Errorf("invalid transaction hash format: must be hex string starting with 0x")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.TransactionReceipt(ctx, txHash)
	})
	if err != nil {
		return nil, err
	}

	// 如果结果为nil，表示交易收据不存在
	if result == nil {
		return nil, nil
	}

	return result.(*eth.TransactionReceipt), nil
}

