// Package service 定义了领域服务
// 领域服务处理跨聚合的业务逻辑或不自然属于任何聚合根的业务操作
package service

import (
	"errors"
	"time"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// BlockValidationService 提供区块验证相关的领域服务
// 这些验证逻辑涉及多个聚合或复杂的业务规则，不适合放在单个聚合根中
type BlockValidationService struct{}

// NewBlockValidationService 创建区块验证服务
func NewBlockValidationService() *BlockValidationService {
	return &BlockValidationService{}
}

// ValidateBlockChain 验证区块链的连续性
// 检查新区块是否可以正确连接到父区块
func (s *BlockValidationService) ValidateBlockChain(
	parentBlock *aggregate.BlockAggregate,
	newBlock *aggregate.BlockAggregate,
) error {
	// 检查父区块哈希是否匹配
	if !newBlock.ParentHash().Equals(parentBlock.Hash()) {
		return errors.New("新区块的父哈希与父区块哈希不匹配")
	}

	// 检查区块号是否连续
	parentNumber, err := parentBlock.Number().Uint64()
	if err != nil {
		return errors.New("无法获取父区块号")
	}

	newNumber, err := newBlock.Number().Uint64()
	if err != nil {
		return errors.New("无法获取新区块号")
	}

	if newNumber != parentNumber+1 {
		return errors.New("区块号不连续")
	}

	// 检查时间戳是否合理（新区块时间戳应该大于父区块）
	if !newBlock.Timestamp().After(parentBlock.Timestamp()) {
		return errors.New("新区块时间戳不能早于或等于父区块时间戳")
	}

	return nil
}

// ValidateBlockTimestamp 验证区块时间戳的合理性
func (s *BlockValidationService) ValidateBlockTimestamp(
	blockTimestamp time.Time,
	maxFutureTime time.Duration,
) error {
	now := time.Now()

	// 检查时间戳不能太远的未来
	if blockTimestamp.After(now.Add(maxFutureTime)) {
		return errors.New("区块时间戳过于超前")
	}

	// 检查时间戳不能为零值
	if blockTimestamp.IsZero() {
		return errors.New("区块时间戳不能为零")
	}

	return nil
}

// ValidateGasUsage 验证区块的燃料使用情况
func (s *BlockValidationService) ValidateGasUsage(
	block *aggregate.BlockAggregate,
	transactions []*aggregate.TransactionAggregate,
) error {
	// 计算所有交易的燃料使用总和
	totalGasUsed := valueobject.NewGasFromUint64(0)

	for _, tx := range transactions {
		if tx.IsMined() {
			receipt := tx.Receipt()
			if receipt != nil {
				totalGasUsed = totalGasUsed.Add(receipt.GasUsed)
			}
		}
	}

	// 检查总燃料使用是否与区块记录一致
	if !totalGasUsed.Equals(block.GasUsed()) {
		return errors.New("区块燃料使用量与交易燃料使用总和不匹配")
	}

	// 检查燃料使用是否超过限制
	if block.GasUsed().GreaterThan(block.GasLimit()) {
		return errors.New("区块燃料使用量超过限制")
	}

	return nil
}

// ValidateTransactionCount 验证区块的交易数量
func (s *BlockValidationService) ValidateTransactionCount(
	block *aggregate.BlockAggregate,
	actualTransactions []*aggregate.TransactionAggregate,
) error {
	expectedCount := block.TransactionCount()
	actualCount := uint64(len(actualTransactions))

	if expectedCount != actualCount {
		return errors.New("区块交易数量与实际交易数量不匹配")
	}

	return nil
}

// ValidateBlockSize 验证区块大小限制
func (s *BlockValidationService) ValidateBlockSize(
	blockSize uint64,
	maxBlockSize uint64,
) error {
	if blockSize > maxBlockSize {
		return errors.New("区块大小超过限制")
	}

	if blockSize == 0 {
		return errors.New("区块大小不能为零")
	}

	return nil
}

// CalculateBlockReward 计算区块奖励
// 这是一个复杂的业务逻辑，涉及多种因素
func (s *BlockValidationService) CalculateBlockReward(
	blockNumber valueobject.BlockNumber,
	gasUsed valueobject.Gas,
	baseFeePerGas *valueobject.Wei,
	transactions []*aggregate.TransactionAggregate,
) (valueobject.Wei, error) {
	// 基础区块奖励（简化实现）
	baseReward := valueobject.NewWeiFromUint64(2000000000000000000) // 2 ETH

	// 计算交易费用奖励
	totalFees := valueobject.NewWeiFromUint64(0)

	for _, tx := range transactions {
		if tx.IsMined() {
			receipt := tx.Receipt()
			if receipt != nil && receipt.EffectiveGasPrice != nil {
				// 计算交易费用
				gasUsedBig, err := receipt.GasUsed.Uint64()
				if err == nil {
					gasUsedWei := valueobject.NewWeiFromUint64(gasUsedBig)
					effectivePriceBig, err := receipt.EffectiveGasPrice.Uint64()
					if err == nil {
						fee := gasUsedWei.Mul(effectivePriceBig)
						totalFees = totalFees.Add(fee)
					}
				}
			}
		}
	}

	// 总奖励 = 基础奖励 + 交易费用
	totalReward := baseReward.Add(totalFees)

	return totalReward, nil
}
