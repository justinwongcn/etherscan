package aggregate

import (
	"testing"
	"time"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

func TestNewBlockAggregate(t *testing.T) {
	// 准备测试数据
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	parentHash, _ := valueobject.NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")
	miner, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	number := valueobject.NewBlockNumberFromUint64(12345)
	timestamp := time.Now()

	// 创建零值对象
	zeroHash, _ := valueobject.NewHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	zeroAddress, _ := valueobject.NewAddress("0x0000000000000000000000000000000000000000")

	tests := []struct {
		name      string
		hash      valueobject.Hash
		miner     valueobject.Address
		wantError bool
	}{
		{
			name:      "有效的区块创建",
			hash:      hash,
			miner:     miner,
			wantError: false,
		},
		{
			name:      "零哈希应该失败",
			hash:      zeroHash,
			miner:     miner,
			wantError: true,
		},
		{
			name:      "零地址矿工应该失败",
			hash:      hash,
			miner:     zeroAddress,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := NewBlockAggregate(number, tt.hash, parentHash, tt.miner, timestamp)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewBlockAggregate() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewBlockAggregate() 返回了意外的错误: %v", err)
				return
			}

			// 验证基本属性
			if !block.Hash().Equals(tt.hash) {
				t.Errorf("Hash() = %v, 期望 %v", block.Hash(), tt.hash)
			}

			if !block.Miner().Equals(tt.miner) {
				t.Errorf("Miner() = %v, 期望 %v", block.Miner(), tt.miner)
			}

			if block.TransactionCount() != 0 {
				t.Errorf("TransactionCount() = %v, 期望 0", block.TransactionCount())
			}

			// 验证事件
			events := block.GetEvents()
			if len(events) != 1 {
				t.Errorf("期望1个事件，实际%d个", len(events))
			}

			if events[0].EventType() != "BlockCreated" {
				t.Errorf("期望BlockCreated事件，实际%s", events[0].EventType())
			}
		})
	}
}

func TestBlockAggregate_SetGasInfo(t *testing.T) {
	// 创建测试区块
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	parentHash, _ := valueobject.NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")
	miner, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	number := valueobject.NewBlockNumberFromUint64(12345)

	block, _ := NewBlockAggregate(number, hash, parentHash, miner, time.Now())

	gasLimit := valueobject.NewGasFromUint64(21000)
	gasUsed := valueobject.NewGasFromUint64(20000)
	baseFee := valueobject.NewWeiFromUint64(1000000000) // 1 Gwei

	tests := []struct {
		name      string
		gasLimit  valueobject.Gas
		gasUsed   valueobject.Gas
		wantError bool
	}{
		{
			name:      "有效的燃料信息",
			gasLimit:  gasLimit,
			gasUsed:   gasUsed,
			wantError: false,
		},
		{
			name:      "已使用燃料超过限制应该失败",
			gasLimit:  gasUsed,
			gasUsed:   gasLimit,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := block.SetGasInfo(tt.gasLimit, tt.gasUsed, &baseFee)

			if tt.wantError {
				if err == nil {
					t.Errorf("SetGasInfo() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("SetGasInfo() 返回了意外的错误: %v", err)
				return
			}

			if !block.GasLimit().Equals(tt.gasLimit) {
				t.Errorf("GasLimit() = %v, 期望 %v", block.GasLimit(), tt.gasLimit)
			}

			if !block.GasUsed().Equals(tt.gasUsed) {
				t.Errorf("GasUsed() = %v, 期望 %v", block.GasUsed(), tt.gasUsed)
			}
		})
	}
}

func TestBlockAggregate_AddTransaction(t *testing.T) {
	// 创建测试区块
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	parentHash, _ := valueobject.NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")
	miner, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	number := valueobject.NewBlockNumberFromUint64(12345)

	block, _ := NewBlockAggregate(number, hash, parentHash, miner, time.Now())

	// 创建交易摘要
	txHash, _ := valueobject.NewHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645Hac136c")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000) // 1 ETH

	txSummary := TransactionSummary{
		Hash:  txHash,
		From:  from,
		To:    &to,
		Value: value,
	}

	// 添加交易前的状态
	initialCount := block.TransactionCount()
	initialEvents := len(block.GetEvents())

	// 添加交易
	block.AddTransaction(txSummary)

	// 验证交易数量增加
	if block.TransactionCount() != initialCount+1 {
		t.Errorf("TransactionCount() = %v, 期望 %v", block.TransactionCount(), initialCount+1)
	}

	// 验证交易列表
	transactions := block.Transactions()
	if len(transactions) != 1 {
		t.Errorf("Transactions() 长度 = %v, 期望 1", len(transactions))
	}

	if !transactions[0].Hash.Equals(txHash) {
		t.Errorf("Transaction hash = %v, 期望 %v", transactions[0].Hash, txHash)
	}

	// 验证事件
	events := block.GetEvents()
	if len(events) != initialEvents+1 {
		t.Errorf("期望%d个事件，实际%d个", initialEvents+1, len(events))
	}

	// 验证最新事件是TransactionAddedToBlock
	lastEvent := events[len(events)-1]
	if lastEvent.EventType() != "TransactionAddedToBlock" {
		t.Errorf("期望TransactionAddedToBlock事件，实际%s", lastEvent.EventType())
	}
}
