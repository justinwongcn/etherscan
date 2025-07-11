// Package repository 定义了领域层的仓储接口
// 这些接口遵循依赖倒置原则，由基础设施层提供具体实现
package repository

import (
	"context"

	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionRepository 定义了交易相关数据访问的接口
// 该接口封装了所有与交易数据获取和操作相关的方法，遵循DDD中的仓储模式
// 具体实现由基础设施层提供，应用层通过此接口访问交易数据
type TransactionRepository interface {
	// GetTransactionByHash 根据交易哈希获取交易详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - txHash: 交易哈希值（32字节的十六进制字符串，以0x开头）
	//
	// 返回:
	//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error)

	// GetTransactionByBlockHashAndIndex 根据区块哈希和交易索引获取交易详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - blockHash: 区块哈希值（32字节的十六进制字符串，以0x开头）
	//   - index: 交易在区块中的索引位置（从0开始的整数）
	//
	// 返回:
	//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*eth.Transaction, error)

	// GetTransactionByBlockNumberAndIndex 根据区块号和交易索引获取交易详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - blockNumber: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//   - index: 交易在区块中的索引位置（从0开始的整数）
	//
	// 返回:
	//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber string, index uint64) (*eth.Transaction, error)

	// SendRawTransaction 发送已签名的交易数据到以太坊网络
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
	//
	// 返回:
	//   - string: 交易哈希（32字节的十六进制字符串）
	//   - error: 如果发送过程中发生错误，将返回相应的错误信息
	SendRawTransaction(ctx context.Context, signedTxData string) (string, error)

	// GetTransactionReceipt 获取交易的收据信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - txHash: 交易哈希值（32字节的十六进制字符串，以0x开头）
	//
	// 返回:
	//   - *eth.TransactionReceipt: 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error)

	// GetTransactionCount 获取指定地址在特定区块的交易数量（nonce值）
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - address: 以太坊账户地址
	//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//
	// 返回:
	//   - uint64: 该地址的交易数量（nonce值）
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionCount(ctx context.Context, address string, blockParameter string) (uint64, error)
}
