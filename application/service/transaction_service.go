// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"

	"github.com/justinwongcn/etherscan/domain"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionService 实现了TransactionServiceInterface接口
// 该结构体封装了与以太坊节点交互的客户端，提供交易相关操作的具体实现
type TransactionService struct {
	// client 是与以太坊节点通信的客户端实例
	client *ethereum.Client
}

// NewTransactionService 创建并初始化一个新的TransactionService实例
// 参数:
//   - client: 已初始化的以太坊客户端实例，用于与节点通信
//
// 返回:
//   - *TransactionService: 初始化完成的服务实例
func NewTransactionService(client *ethereum.Client) *TransactionService {
	return &TransactionService{
		client:    client,
	}
}

// GetTransactionByHash 实现了TransactionServiceInterface接口中的同名方法
// 根据交易哈希获取交易的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionByHash(ctx context.Context, txHash string) (*domain.Transaction, error) {
	// 调用以太坊客户端获取交易信息
	ethTx, err := s.client.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	
	// 使用转换器将eth.Transaction转换为domain.Transaction
	converter := domain.NewTransactionConverter()
	return converter.ConvertToTransaction(ethTx), nil
}

// GetTransactionByIndex 实现了TransactionServiceInterface接口中的同名方法
// 根据区块标识符和交易索引获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//   - index: 交易在区块中的索引位置
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*eth.Transaction, error) {
	// 解析并标准化区块参数
	param, err := ethereum.ParseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法
	// 如果是区块哈希（以0x开头且长度大于10的十六进制字符串）
	if len(param) >= 2 && param[:2] == "0x" && len(param) > 10 {
		return s.client.GetTransactionByBlockHashAndIndex(ctx, param, index)
	}
	return s.client.GetTransactionByBlockNumberAndIndex(ctx, param, index)
}

// SendRawTransaction 实现了TransactionServiceInterface接口中的同名方法
// 发送已签名的交易数据到以太坊网络
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
//
// 返回:
//   - string: 交易哈希（32字节的十六进制字符串）
//   - error: 如果发送过程中发生错误，将返回相应的错误信息
func (s *TransactionService) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	// 调用以太坊客户端发送交易
	return s.client.SendRawTransaction(ctx, signedTxData)
}

// GetTransactionReceipt 实现了TransactionServiceInterface接口中的同名方法
// 获取交易的收据信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *eth.TransactionReceipt: 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	// 调用以太坊客户端获取交易收据
	return s.client.GetTransactionReceipt(ctx, txHash)
}

// GetTransactionCount 实现了TransactionServiceInterface接口中的同名方法
// 获取指定地址在特定区块的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - uint64: 交易数量
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (uint64, error) {
	// 调用以太坊客户端获取交易数量
	return s.client.GetTransactionCount(ctx, address, blockHashOrNumber)
}
