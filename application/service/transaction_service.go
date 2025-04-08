// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"
	"fmt"

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
		client: client,
	}
}

// GetTransactionByHash 实现了TransactionServiceInterface接口中的同名方法
// 根据交易哈希获取交易的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *domain.Transaction: 包含交易完整信息的领域模型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionByHash(ctx context.Context, txHash string) (*domain.Transaction, error) {
	// 调用以太坊客户端获取交易信息
	ethTx, err := s.client.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	// 使用转换器将eth.Transaction转换为domain.Transaction
	converter := domain.NewTransactionConverter(50)
	return converter.ConvertToTransaction(ethTx), nil
}

// GetTransactionByIndex 实现了TransactionServiceInterface接口中的同名方法
// 根据区块标识符和交易索引获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，可以是区块号（数字字符串）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//   - index: 交易在区块中的索引位置（从0开始的整数）
//
// 返回:
//   - *domain.Transaction: 包含交易完整信息的领域模型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*domain.Transaction, error) {
	// 解析并标准化区块参数
	param, err := ethereum.ParseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法
	// 如果是区块哈希（以0x开头且长度大于10的十六进制字符串）
	var ethTx *eth.Transaction
	if len(param) >= 2 && param[:2] == "0x" && len(param) > 10 {
		ethTx, err = s.client.GetTransactionByBlockHashAndIndex(ctx, param, index)
	} else {
		ethTx, err = s.client.GetTransactionByBlockNumberAndIndex(ctx, param, index)
	}
	if err != nil {
		return nil, err
	}

	// 使用转换器将eth.Transaction转换为domain.Transaction
	converter := domain.NewTransactionConverter(50)
	return converter.ConvertToTransaction(ethTx), nil
}

// SendRawTransaction 实现了TransactionServiceInterface接口中的同名方法
// 发送已签名的交易数据到以太坊网络
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
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
//   - *domain.TransactionReceipt: 包含交易收据完整信息的领域模型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionReceipt(ctx context.Context, txHash string) (*domain.TransactionReceipt, error) {
	// 调用以太坊客户端获取交易收据
	receipt, err := s.client.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	// 使用转换器将eth.TransactionReceipt转换为domain.TransactionReceipt
	converter := domain.NewTransactionReceiptConverter()
	return converter.ConvertToTransactionReceipt(receipt), nil
}

// GetTransactionCount 实现了TransactionServiceInterface接口中的同名方法
// 获取指定地址在特定区块的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - string: 交易数量（十六进制字符串）
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (string, error) {
	// 调用以太坊客户端获取交易数量
	count, err := s.client.GetTransactionCount(ctx, address, blockHashOrNumber)
	if err != nil {
		return "", err
	}
	// 将uint64类型的交易数量转换为十进制字符串
	return fmt.Sprintf("%d", count), nil
}
