// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/domain/valueobject"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockService 实现了BlockServiceInterface接口
// 该结构体遵循DDD架构原则，通过依赖仓储接口而非具体实现来提供区块相关的业务逻辑
type BlockService struct {
	// blockRepo 是区块数据访问的仓储接口
	// 通过依赖倒置原则，应用层依赖抽象而非具体实现
	blockRepo repository.BlockRepository
}

// NewBlockService 创建并初始化一个新的BlockService实例
// 参数:
//   - blockRepo: 实现了BlockRepository接口的仓储实例，用于区块数据访问
//
// 返回:
//   - *BlockService: 初始化完成的服务实例
func NewBlockService(blockRepo repository.BlockRepository) *BlockService {
	return &BlockService{
		blockRepo: blockRepo,
	}
}

// GetLatestBlockHeight 实现了BlockServiceInterface接口中的同名方法
// 通过调用区块仓储获取当前网络的最新区块高度
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - string: 最新区块的高度编号（十进制字符串格式）
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetLatestBlockHeight 获取最新区块高度
func (s *BlockService) GetLatestBlockHeight(ctx context.Context) (string, error) {
	// 从仓储获取最新区块号
	height, err := s.blockRepo.GetLatestBlockNumber(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", height), nil
}

// parseBlockParameter 解析并标准化区块参数格式
// 将用户输入的区块标识符转换为以太坊API支持的格式
// 参数:
//   - blockHashOrNumber: 区块标识符，可以是区块号、区块哈希或特殊标识符
//
// 返回:
//   - string: 转换后的标准格式参数
//   - error: 如果参数格式无效，将返回错误信息

// GetBlock 实现了BlockServiceInterface接口中的同名方法
// 根据区块标识符获取区块的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希
//
// 返回:
//   - *eth.Block: 包含区块完整信息的go-ethlibs类型指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetBlock 获取区块信息
func (s *BlockService) GetBlock(ctx context.Context, blockHashOrNumber string, fullTx bool) (*eth.Block, error) {
	// 解析并标准化区块参数
	param, err := ethereum.ParseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法，直接返回go-ethlibs类型
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" && len(blockHashOrNumber) > 10 {
		return s.blockRepo.GetBlockByHash(ctx, param, fullTx)
	} else {
		return s.blockRepo.GetBlockByNumber(ctx, param, fullTx)
	}
}

// GetBlockTransactionCount 实现了BlockServiceInterface接口中的同名方法
// 获取指定区块中包含的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - string: 区块中的交易数量（十进制字符串格式）
//   - error: 如果查询过程中发生错误，将返回相应的错误信息

// GetBlockTransactionCount 获取区块交易数量
func (s *BlockService) GetBlockTransactionCount(ctx context.Context, blockHashOrNumber string) (string, error) {
	// 获取区块聚合根
	blockAggregate, err := s.GetBlockAggregate(ctx, blockHashOrNumber, false) // 不需要完整交易信息
	if err != nil {
		return "", err
	}

	// 从聚合根获取交易数量
	count := blockAggregate.TransactionCount()
	return fmt.Sprintf("%d", count), nil
}

// GetBlockAggregate 获取区块聚合根（新的DDD架构方法）
// 这个方法使用新的聚合根模式，将逐步替代GetBlock方法
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希
//
// 返回:
//   - *aggregate.BlockAggregate: 包含区块完整信息的聚合根指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetBlockAggregate(ctx context.Context, blockHashOrNumber string, fullTx bool) (*aggregate.BlockAggregate, error) {
	// 解析并标准化区块参数
	param, err := ethereum.ParseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法
	// 注意：为了获取交易信息，这里总是传递true给仓储
	var ethBlock *eth.Block
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" && len(blockHashOrNumber) > 10 {
		ethBlock, err = s.blockRepo.GetBlockByHash(ctx, param, true)
	} else {
		ethBlock, err = s.blockRepo.GetBlockByNumber(ctx, param, true)
	}
	if err != nil {
		return nil, err
	}

	if ethBlock == nil {
		return nil, fmt.Errorf("区块不存在: %s", blockHashOrNumber)
	}

	// 转换为值对象
	blockHash, err := valueobject.NewHash(ethBlock.Hash.String())
	if err != nil {
		return nil, fmt.Errorf("无效的区块哈希: %w", err)
	}

	parentHash, err := valueobject.NewHash(ethBlock.ParentHash.String())
	if err != nil {
		return nil, fmt.Errorf("无效的父区块哈希: %w", err)
	}

	miner, err := valueobject.NewAddress(ethBlock.Miner.String())
	if err != nil {
		return nil, fmt.Errorf("无效的矿工地址: %w", err)
	}

	// 解析区块号
	var blockNumber valueobject.BlockNumber
	if ethBlock.Number != nil {
		var err error
		blockNumber, err = valueobject.NewBlockNumberFromBigInt(ethBlock.Number.Big())
		if err != nil {
			return nil, fmt.Errorf("无效的区块号: %w", err)
		}
	} else {
		return nil, fmt.Errorf("区块号不能为空")
	}

	// 解析时间戳
	timestamp := time.Unix(int64(ethBlock.Timestamp.UInt64()), 0)

	// 创建区块聚合根
	blockAggregate, err := aggregate.NewBlockAggregate(blockNumber, blockHash, parentHash, miner, timestamp)
	if err != nil {
		return nil, fmt.Errorf("创建区块聚合根失败: %w", err)
	}

	// 设置燃料信息
	// eth.Quantity类型不能直接与nil比较，需要检查其值
	gasLimit, err := valueobject.NewGasFromBigInt(ethBlock.GasLimit.Big())
	if err != nil {
		return nil, fmt.Errorf("无效的燃料限制: %w", err)
	}

	gasUsed, err := valueobject.NewGasFromBigInt(ethBlock.GasUsed.Big())
	if err != nil {
		return nil, fmt.Errorf("无效的燃料使用量: %w", err)
	}

	var baseFeePerGas *valueobject.Wei
	if ethBlock.BaseFeePerGas != nil {
		baseFee, err := valueobject.NewWeiFromBigInt(ethBlock.BaseFeePerGas.Big())
		if err != nil {
			return nil, fmt.Errorf("无效的基础费用: %w", err)
		}
		baseFeePerGas = &baseFee
	}

	err = blockAggregate.SetGasInfo(gasLimit, gasUsed, baseFeePerGas)
	if err != nil {
		return nil, fmt.Errorf("设置燃料信息失败: %w", err)
	}

	// 设置默克尔树根
	// eth.Data32类型不能直接与nil比较，直接使用其值
	txRoot, _ := valueobject.NewHash(ethBlock.TransactionsRoot.String())
	stateRoot, _ := valueobject.NewHash(ethBlock.StateRoot.String())
	receiptsRoot, _ := valueobject.NewHash(ethBlock.ReceiptsRoot.String())

	blockAggregate.SetMerkleRoots(txRoot, stateRoot, receiptsRoot)

	// 添加交易摘要（总是添加，不管fullTx参数）
	if ethBlock.Transactions != nil {
		for _, tx := range ethBlock.Transactions {
			// eth.TxOrHash是一个结构体，直接使用Transaction字段
			txSummary, err := s.convertToTransactionSummary(&tx.Transaction)
			if err != nil {
				// 记录错误但继续处理其他交易
				continue
			}
			blockAggregate.AddTransaction(txSummary)
		}
	}

	return blockAggregate, nil
}

// convertToTransactionSummary 将eth.Transaction转换为TransactionSummary
func (s *BlockService) convertToTransactionSummary(ethTx *eth.Transaction) (aggregate.TransactionSummary, error) {
	txHash, err := valueobject.NewHash(ethTx.Hash.String())
	if err != nil {
		return aggregate.TransactionSummary{}, err
	}

	from, err := valueobject.NewAddress(ethTx.From.String())
	if err != nil {
		return aggregate.TransactionSummary{}, err
	}

	var to *valueobject.Address
	if ethTx.To != nil {
		toAddr, err := valueobject.NewAddress(ethTx.To.String())
		if err != nil {
			return aggregate.TransactionSummary{}, err
		}
		to = &toAddr
	}

	value, err := valueobject.NewWeiFromBigInt(ethTx.Value.Big())
	if err != nil {
		return aggregate.TransactionSummary{}, err
	}

	return aggregate.TransactionSummary{
		Hash:  txHash,
		From:  from,
		To:    to,
		Value: value,
	}, nil
}
