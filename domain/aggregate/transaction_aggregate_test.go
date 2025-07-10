package aggregate

import (
	"testing"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

func TestNewTransactionAggregate(t *testing.T) {
	// 创建测试数据
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000) // 1 ETH
	gas := valueobject.NewGasFromUint64(21000)
	input := []byte("test data")

	tests := []struct {
		name      string
		hash      valueobject.Hash
		from      valueobject.Address
		to        *valueobject.Address
		value     valueobject.Wei
		nonce     uint64
		gas       valueobject.Gas
		input     []byte
		wantError bool
	}{
		{
			name:      "有效的转账交易",
			hash:      hash,
			from:      from,
			to:        &to,
			value:     value,
			nonce:     1,
			gas:       gas,
			input:     []byte{},
			wantError: false,
		},
		{
			name:      "合约创建交易（to为nil）",
			hash:      hash,
			from:      from,
			to:        nil,
			value:     value,
			nonce:     1,
			gas:       gas,
			input:     input,
			wantError: false,
		},
		{
			name: "零哈希应该失败",
			hash: func() valueobject.Hash {
				h, _ := valueobject.NewHash("0x0000000000000000000000000000000000000000000000000000000000000000")
				return h
			}(),
			from:      from,
			to:        &to,
			value:     value,
			nonce:     1,
			gas:       gas,
			input:     []byte{},
			wantError: true,
		},
		{
			name: "零地址发送方应该失败",
			hash: hash,
			from: func() valueobject.Address {
				a, _ := valueobject.NewAddress("0x0000000000000000000000000000000000000000")
				return a
			}(),
			to:        &to,
			value:     value,
			nonce:     1,
			gas:       gas,
			input:     []byte{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := NewTransactionAggregateWithoutSignature(tt.hash, tt.from, tt.to, tt.value, tt.nonce, tt.gas, tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewTransactionAggregate() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewTransactionAggregate() 返回了意外的错误: %v", err)
				return
			}

			// 验证基本属性
			if !tx.Hash().Equals(tt.hash) {
				t.Errorf("Hash() = %v, 期望 %v", tx.Hash(), tt.hash)
			}

			if !tx.From().Equals(tt.from) {
				t.Errorf("From() = %v, 期望 %v", tx.From(), tt.from)
			}

			if tt.to != nil {
				if tx.To() == nil || !tx.To().Equals(*tt.to) {
					t.Errorf("To() = %v, 期望 %v", tx.To(), *tt.to)
				}
			} else {
				if tx.To() != nil {
					t.Errorf("To() 应该为nil，实际为 %v", tx.To())
				}
			}

			if !tx.Value().Equals(tt.value) {
				t.Errorf("Value() = %v, 期望 %v", tx.Value(), tt.value)
			}

			if tx.Nonce() != tt.nonce {
				t.Errorf("Nonce() = %v, 期望 %v", tx.Nonce(), tt.nonce)
			}

			if !tx.Gas().Equals(tt.gas) {
				t.Errorf("Gas() = %v, 期望 %v", tx.Gas(), tt.gas)
			}

			// 验证初始状态
			if tx.Status() != TransactionPending {
				t.Errorf("新交易状态应该是Pending")
			}

			if tx.IsMined() {
				t.Errorf("新交易不应该是已挖矿状态")
			}

			if tx.IsFailed() {
				t.Errorf("新交易不应该是失败状态")
			}

			if !tx.IsPending() {
				t.Errorf("新交易应该是待处理状态")
			}

			// 验证合约创建检测
			if tt.to == nil {
				if !tx.IsContractCreation() {
					t.Errorf("to为nil的交易应该是合约创建交易")
				}
			} else {
				if tx.IsContractCreation() {
					t.Errorf("有to地址的交易不应该是合约创建交易")
				}
			}

			// 验证事件
			events := tx.GetEvents()
			if len(events) != 1 {
				t.Errorf("应该有1个事件，实际有 %d 个", len(events))
			}

			if events[0].EventType() != "TransactionCreated" {
				t.Errorf("事件类型应该是 TransactionCreated，实际是 %s", events[0].EventType())
			}
		})
	}
}

func TestTransactionAggregate_SetLegacyGasPrice(t *testing.T) {
	// 创建交易
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000)
	gas := valueobject.NewGasFromUint64(21000)

	tx, _ := NewTransactionAggregateWithoutSignature(hash, from, &to, value, 1, gas, []byte{})

	// 设置Legacy燃料价格
	gasPrice := valueobject.NewWeiFromUint64(20000000000) // 20 Gwei
	tx.SetLegacyGasPrice(gasPrice)

	// 验证燃料价格
	if tx.GasPrice() == nil || !tx.GasPrice().Equals(gasPrice) {
		t.Errorf("GasPrice() = %v, 期望 %v", tx.GasPrice(), gasPrice)
	}

	// 验证交易类型
	if tx.Type() != LegacyTransaction {
		t.Errorf("Type() = %v, 期望 %v", tx.Type(), LegacyTransaction)
	}

	// EIP-1559字段应该为nil
	if tx.MaxFeePerGas() != nil {
		t.Errorf("MaxFeePerGas() 应该为nil")
	}

	if tx.MaxPriorityFee() != nil {
		t.Errorf("MaxPriorityFee() 应该为nil")
	}
}

func TestTransactionAggregate_SetEIP1559Fees(t *testing.T) {
	// 创建交易
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000)
	gas := valueobject.NewGasFromUint64(21000)

	tx, _ := NewTransactionAggregateWithoutSignature(hash, from, &to, value, 1, gas, []byte{})

	// 设置EIP-1559费用
	maxFeePerGas := valueobject.NewWeiFromUint64(30000000000)  // 30 Gwei
	maxPriorityFee := valueobject.NewWeiFromUint64(2000000000) // 2 Gwei

	tx.SetEIP1559Fees(maxFeePerGas, maxPriorityFee)

	// 验证EIP-1559费用
	if tx.MaxFeePerGas() == nil || !tx.MaxFeePerGas().Equals(maxFeePerGas) {
		t.Errorf("MaxFeePerGas() = %v, 期望 %v", tx.MaxFeePerGas(), maxFeePerGas)
	}

	if tx.MaxPriorityFee() == nil || !tx.MaxPriorityFee().Equals(maxPriorityFee) {
		t.Errorf("MaxPriorityFee() = %v, 期望 %v", tx.MaxPriorityFee(), maxPriorityFee)
	}

	// 验证交易类型
	if tx.Type() != EIP1559Transaction {
		t.Errorf("Type() = %v, 期望 %v", tx.Type(), EIP1559Transaction)
	}

	// Legacy燃料价格应该为nil
	if tx.GasPrice() != nil {
		t.Errorf("GasPrice() 应该为nil")
	}
}

func TestTransactionAggregate_MarkAsMined(t *testing.T) {
	// 创建交易
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000)
	gas := valueobject.NewGasFromUint64(21000)

	tx, _ := NewTransactionAggregateWithoutSignature(hash, from, &to, value, 1, gas, []byte{})

	// 清除创建事件
	tx.ClearEvents()

	// 创建挖矿信息
	blockHash, _ := valueobject.NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")
	blockNumber := valueobject.NewBlockNumberFromUint64(12345)
	transactionIndex := uint64(0)
	gasUsed := valueobject.NewGasFromUint64(21000)

	// 创建交易收据
	receipt := TransactionReceipt{
		GasUsed:           gasUsed,
		CumulativeGasUsed: gasUsed,
		ContractAddress:   nil,
		Logs:              []Log{},
		Status:            1, // 成功
	}

	// 标记为已挖矿
	tx.MarkAsMined(blockHash, blockNumber, transactionIndex, receipt)

	// 验证挖矿状态
	if !tx.IsMined() {
		t.Errorf("交易应该是已挖矿状态")
	}

	if tx.IsPending() {
		t.Errorf("交易不应该是待处理状态")
	}

	if tx.Status() != TransactionMined {
		t.Errorf("Status() = %v, 期望 %v", tx.Status(), TransactionMined)
	}

	// 验证区块信息
	if tx.BlockHash() == nil || !tx.BlockHash().Equals(blockHash) {
		t.Errorf("BlockHash() = %v, 期望 %v", tx.BlockHash(), blockHash)
	}

	if tx.BlockNumber() == nil || !tx.BlockNumber().Equals(blockNumber) {
		t.Errorf("BlockNumber() = %v, 期望 %v", tx.BlockNumber(), blockNumber)
	}

	if tx.TransactionIndex() == nil || *tx.TransactionIndex() != transactionIndex {
		t.Errorf("TransactionIndex() = %v, 期望 %v", tx.TransactionIndex(), transactionIndex)
	}

	// 验证收据
	txReceipt := tx.Receipt()
	if txReceipt == nil {
		t.Errorf("Receipt() 不应该为nil")
	} else {
		if !txReceipt.GasUsed.Equals(gasUsed) {
			t.Errorf("收据Gas使用量不匹配")
		}

		if txReceipt.Status != 1 {
			t.Errorf("收据状态应该是成功(1)")
		}
	}

	// 验证事件
	events := tx.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个挖矿事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "TransactionMined" {
		t.Errorf("事件类型应该是 TransactionMined，实际是 %s", events[0].EventType())
	}
}

func TestTransactionAggregate_MarkAsFailed(t *testing.T) {
	// 创建交易
	hash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	from, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")
	to, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	value := valueobject.NewWeiFromUint64(1000000000000000000)
	gas := valueobject.NewGasFromUint64(21000)

	tx, _ := NewTransactionAggregateWithoutSignature(hash, from, &to, value, 1, gas, []byte{})

	// 清除创建事件
	tx.ClearEvents()

	// 标记为失败
	reason := "insufficient funds"
	tx.MarkAsFailed(reason)

	// 验证失败状态
	if !tx.IsFailed() {
		t.Errorf("交易应该是失败状态")
	}

	if tx.IsPending() {
		t.Errorf("交易不应该是待处理状态")
	}

	if tx.IsMined() {
		t.Errorf("交易不应该是已挖矿状态")
	}

	if tx.Status() != TransactionFailed {
		t.Errorf("Status() = %v, 期望 %v", tx.Status(), TransactionFailed)
	}

	// 验证事件
	events := tx.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个失败事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "TransactionFailed" {
		t.Errorf("事件类型应该是 TransactionFailed，实际是 %s", events[0].EventType())
	}
}
