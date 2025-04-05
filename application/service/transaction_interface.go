// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"

	"github.com/justinwongcn/etherscan/domain"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionServiceInterface 定义了交易服务的接口规范
// 该接口提供了查询和发送以太坊交易相关信息的方法集合
type TransactionServiceInterface interface {
	// GetTransactionByHash 根据交易哈希获取交易详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - txHash: 交易哈希（32字节的十六进制字符串）
	// 返回:
	//   - *eth.Transaction: 包含交易完整信息的结构体指针
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionByHash(ctx context.Context, txHash string) (*domain.Transaction, error)

	// GetTransactionByIndex 根据区块标识符和交易索引获取交易详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - blockHashOrNumber: 区块标识符，可以是区块号（数字字符串）或区块哈希（0x开头的十六进制字符串）
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//   - index: 交易在区块中的索引位置
	// 返回:
	//   - *eth.Transaction: 包含交易完整信息的结构体指针
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*eth.Transaction, error)

	// SendRawTransaction 发送已签名的交易数据到以太坊网络
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
	// 返回:
	//   - string: 交易哈希（32字节的十六进制字符串）
	//   - error: 如果发送过程中发生错误，将返回相应的错误信息
	SendRawTransaction(ctx context.Context, signedTxData string) (string, error)

	// GetTransactionReceipt 获取交易收据信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - txHash: 交易哈希（32字节的十六进制字符串）
	// 返回:
	//   - *eth.TransactionReceipt: 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error)

	// GetTransactionCount 获取指定地址在特定区块的交易数量
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - address: 以太坊账户地址
	//   - blockHashOrNumber: 区块标识符，可以是区块号（数字字符串）或区块哈希（0x开头的十六进制字符串）
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	// 返回:
	//   - uint64: 交易数量
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (uint64, error)
}
