// Package repository 提供领域仓储接口的具体实现
// 该包属于基础设施层，负责与外部数据源（以太坊节点）进行交互
package repository

import (
	"context"

	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionRepositoryImpl 实现了domain.repository.TransactionRepository接口
// 该结构体封装了与以太坊节点交互的客户端，提供交易数据访问的具体实现
type TransactionRepositoryImpl struct {
	// client 是与以太坊节点通信的客户端实例
	// 通过该客户端执行实际的RPC调用来获取交易数据
	client *ethereum.Client
}

// NewTransactionRepository 创建并初始化一个新的TransactionRepositoryImpl实例
// 该函数是TransactionRepositoryImpl的工厂方法，确保正确的依赖注入
// 参数:
//   - client: 已初始化的以太坊客户端实例，用于与节点通信
//
// 返回:
//   - repository.TransactionRepository: 实现了TransactionRepository接口的实例
func NewTransactionRepository(client *ethereum.Client) repository.TransactionRepository {
	return &TransactionRepositoryImpl{
		client: client,
	}
}

// GetTransactionByHash 实现了TransactionRepository接口中的同名方法
// 根据交易哈希获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希值（32字节的十六进制字符串，以0x开头）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	return r.client.GetTransactionByHash(ctx, txHash)
}

// GetTransactionByBlockHashAndIndex 实现了TransactionRepository接口中的同名方法
// 根据区块哈希和交易索引获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHash: 区块哈希值（32字节的十六进制字符串，以0x开头）
//   - index: 交易在区块中的索引位置（从0开始的整数）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*eth.Transaction, error) {
	return r.client.GetTransactionByBlockHashAndIndex(ctx, blockHash, index)
}

// GetTransactionByBlockNumberAndIndex 实现了TransactionRepository接口中的同名方法
// 根据区块号和交易索引获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockNumber: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//   - index: 交易在区块中的索引位置（从0开始的整数）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针，如果交易不存在则返回nil
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber string, index uint64) (*eth.Transaction, error) {
	return r.client.GetTransactionByBlockNumberAndIndex(ctx, blockNumber, index)
}

// SendRawTransaction 实现了TransactionRepository接口中的同名方法
// 发送已签名的交易数据到以太坊网络
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
//
// 返回:
//   - string: 交易哈希（32字节的十六进制字符串）
//   - error: 如果发送过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	return r.client.SendRawTransaction(ctx, signedTxData)
}

// GetTransactionReceipt 实现了TransactionRepository接口中的同名方法
// 获取交易的收据信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希值（32字节的十六进制字符串，以0x开头）
//
// 返回:
//   - *eth.TransactionReceipt: 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	return r.client.GetTransactionReceipt(ctx, txHash)
}

// GetTransactionCount 实现了TransactionRepository接口中的同名方法
// 获取指定地址在特定区块的交易数量（nonce值）
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//
// 返回:
//   - uint64: 该地址的交易数量（nonce值）
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *TransactionRepositoryImpl) GetTransactionCount(ctx context.Context, address string, blockParameter string) (uint64, error) {
	return r.client.GetTransactionCount(ctx, address, blockParameter)
}
