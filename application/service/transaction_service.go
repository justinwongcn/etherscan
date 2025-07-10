// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"
	"fmt"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/domain/valueobject"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionService 实现了TransactionServiceInterface接口
// 该结构体遵循DDD架构原则，通过依赖仓储接口而非具体实现来提供交易相关的业务逻辑
type TransactionService struct {
	// transactionRepo 是交易数据访问的仓储接口
	// 通过依赖倒置原则，应用层依赖抽象而非具体实现
	transactionRepo repository.TransactionRepository
}

// NewTransactionService 创建并初始化一个新的TransactionService实例
// 参数:
//   - transactionRepo: 实现了TransactionRepository接口的仓储实例，用于交易数据访问
//
// 返回:
//   - *TransactionService: 初始化完成的服务实例
func NewTransactionService(transactionRepo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// GetTransactionByHash 实现了TransactionServiceInterface接口中的同名方法
// 根据交易哈希获取交易的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的go-ethlibs类型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetTransactionByHash 获取交易信息
func (s *TransactionService) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	// 直接通过交易仓储获取交易信息，返回go-ethlibs类型
	return s.transactionRepo.GetTransactionByHash(ctx, txHash)
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
//   - *eth.Transaction: 包含交易完整信息的go-ethlibs类型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetTransactionByIndex 根据区块和索引获取交易信息
func (s *TransactionService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*eth.Transaction, error) {
	// 解析并标准化区块参数
	param, err := ethereum.ParseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法，直接返回go-ethlibs类型
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" && len(blockHashOrNumber) > 10 {
		return s.transactionRepo.GetTransactionByBlockHashAndIndex(ctx, param, index)
	} else {
		return s.transactionRepo.GetTransactionByBlockNumberAndIndex(ctx, param, index)
	}
}

// SendRawTransaction 实现了TransactionServiceInterface接口中的同名方法
// 发送已签名的交易数据到以太坊网络
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - string: 交易哈希（32字节的十六进制字符串）
//   - error: 如果发送过程中发生错误，将返回相应的错误信息

// SendRawTransaction 发送已签名的交易
func (s *TransactionService) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	// 通过仓储发送交易
	return s.transactionRepo.SendRawTransaction(ctx, signedTxData)
}

// GetTransactionReceipt 实现了TransactionServiceInterface接口中的同名方法
// 获取交易的收据信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *eth.TransactionReceipt: 包含交易收据完整信息的go-ethlibs类型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetTransactionReceipt 获取交易收据
func (s *TransactionService) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	// 直接通过交易仓储获取交易收据，返回go-ethlibs类型
	return s.transactionRepo.GetTransactionReceipt(ctx, txHash)
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

// GetTransactionCount 获取指定地址的交易数量
func (s *TransactionService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (string, error) {
	// 直接通过仓储获取交易数量
	count, err := s.transactionRepo.GetTransactionCount(ctx, address, blockHashOrNumber)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", count), nil
}

// GetTransactionAggregate 获取交易聚合根（新的DDD架构方法）
// 这个方法使用新的聚合根模式，将逐步替代GetTransactionByHash方法
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *aggregate.TransactionAggregate: 包含交易完整信息的聚合根指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *TransactionService) GetTransactionAggregate(ctx context.Context, txHash string) (*aggregate.TransactionAggregate, error) {
	// 通过交易仓储获取交易信息
	ethTx, err := s.transactionRepo.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	if ethTx == nil {
		return nil, fmt.Errorf("交易不存在: %s", txHash)
	}

	// 转换为交易聚合根
	return s.convertToTransactionAggregate(ethTx)
}

// convertToTransactionAggregate 将eth.Transaction转换为TransactionAggregate
func (s *TransactionService) convertToTransactionAggregate(ethTx *eth.Transaction) (*aggregate.TransactionAggregate, error) {
	// 转换基本值对象
	txHash, err := valueobject.NewHash(ethTx.Hash.String())
	if err != nil {
		return nil, fmt.Errorf("无效的交易哈希: %w", err)
	}

	from, err := valueobject.NewAddress(ethTx.From.String())
	if err != nil {
		return nil, fmt.Errorf("无效的发送方地址: %w", err)
	}

	var to *valueobject.Address
	if ethTx.To != nil {
		toAddr, err := valueobject.NewAddress(ethTx.To.String())
		if err != nil {
			return nil, fmt.Errorf("无效的接收方地址: %w", err)
		}
		to = &toAddr
	}

	value, err := valueobject.NewWeiFromBigInt(ethTx.Value.Big())
	if err != nil {
		return nil, fmt.Errorf("无效的交易金额: %w", err)
	}

	gas, err := valueobject.NewGasFromBigInt(ethTx.Gas.Big())
	if err != nil {
		return nil, fmt.Errorf("无效的燃料限制: %w", err)
	}

	// 简化处理，直接创建聚合根

	// 创建签名值对象
	signature := valueobject.NewSignatureFromStrings(ethTx.V.String(), ethTx.R.String(), ethTx.S.String())

	// 创建交易聚合根（注意参数顺序）
	txAggregate, err := aggregate.NewTransactionAggregate(txHash, from, to, value, ethTx.Nonce.UInt64(), gas, []byte(ethTx.Input), signature)
	if err != nil {
		return nil, fmt.Errorf("创建交易聚合根失败: %w", err)
	}

	// 设置燃料价格信息
	if ethTx.GasPrice != nil {
		gasPrice, err := valueobject.NewWeiFromBigInt(ethTx.GasPrice.Big())
		if err == nil {
			txAggregate.SetLegacyGasPrice(gasPrice)
		}
	}

	// 设置EIP-1559费用
	if ethTx.MaxFeePerGas != nil && ethTx.MaxPriorityFeePerGas != nil {
		maxFeePerGas, err1 := valueobject.NewWeiFromBigInt(ethTx.MaxFeePerGas.Big())
		maxPriorityFee, err2 := valueobject.NewWeiFromBigInt(ethTx.MaxPriorityFeePerGas.Big())
		if err1 == nil && err2 == nil {
			txAggregate.SetEIP1559Fees(maxFeePerGas, maxPriorityFee)
		}
	}

	// 如果交易已挖矿，设置区块信息
	if ethTx.BlockHash != nil && ethTx.BlockNumber != nil && ethTx.Index != nil {
		blockHash, err1 := valueobject.NewHash(ethTx.BlockHash.String())
		blockNumber, err2 := valueobject.NewBlockNumber(ethTx.BlockNumber.String())
		if err1 == nil && err2 == nil {
			txIndex := ethTx.Index.UInt64()
			// 创建一个简单的收据（实际应用中可能需要更完整的收据信息）
			receipt := aggregate.TransactionReceipt{
				GasUsed:           gas, // 简化处理，使用gas限制作为使用量
				CumulativeGasUsed: gas,
				Status:            1, // 假设成功
			}
			txAggregate.MarkAsMined(blockHash, blockNumber, txIndex, receipt)
		}
	}

	return txAggregate, nil
}
