// Package aggregate 定义了区块聚合根
package aggregate

import (
	"errors"
	"time"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// BlockAggregate 表示区块聚合根
// 区块是以太坊区块链的基本单位，包含交易列表和区块元数据
// 作为聚合根，它负责维护区块内部数据的一致性
type BlockAggregate struct {
	BaseAggregateRoot

	// 基本区块信息
	number     valueobject.BlockNumber
	hash       valueobject.Hash
	parentHash valueobject.Hash
	miner      valueobject.Address
	timestamp  time.Time

	// 燃料相关
	gasLimit      valueobject.Gas
	gasUsed       valueobject.Gas
	baseFeePerGas *valueobject.Wei // EIP-1559，可选

	// 默克尔树根哈希
	transactionsRoot valueobject.Hash
	stateRoot        valueobject.Hash
	receiptsRoot     valueobject.Hash

	// 交易列表
	transactions     []TransactionSummary
	transactionCount uint64

	// EIP-4895 提款相关
	withdrawalsRoot *valueobject.Hash
	withdrawals     []Withdrawal
}

// TransactionSummary 表示区块中交易的摘要信息
// 避免在区块聚合中包含完整的交易对象，保持聚合边界清晰
type TransactionSummary struct {
	Hash  valueobject.Hash
	From  valueobject.Address
	To    *valueobject.Address // 可能为nil（合约创建交易）
	Value valueobject.Wei
}

// Withdrawal 表示EIP-4895提款操作
type Withdrawal struct {
	Index          uint64
	ValidatorIndex uint64
	Address        valueobject.Address
	Amount         valueobject.Wei
}

// NewBlockAggregate 创建新的区块聚合根
func NewBlockAggregate(
	number valueobject.BlockNumber,
	hash valueobject.Hash,
	parentHash valueobject.Hash,
	miner valueobject.Address,
	timestamp time.Time,
) (*BlockAggregate, error) {

	// 验证必要字段
	if hash.IsZero() {
		return nil, errors.New("区块哈希不能为零")
	}

	if miner.IsZero() {
		return nil, errors.New("矿工地址不能为零")
	}

	block := &BlockAggregate{
		BaseAggregateRoot: NewBaseAggregateRoot(hash.String()),
		number:            number,
		hash:              hash,
		parentHash:        parentHash,
		miner:             miner,
		timestamp:         timestamp,
		transactions:      make([]TransactionSummary, 0),
	}

	// 发布区块创建事件
	event := NewBlockCreatedEvent(hash.String(), number, timestamp.Unix())
	block.AddEvent(event)

	return block, nil
}

// Number 返回区块号
func (b *BlockAggregate) Number() valueobject.BlockNumber {
	return b.number
}

// Hash 返回区块哈希
func (b *BlockAggregate) Hash() valueobject.Hash {
	return b.hash
}

// ParentHash 返回父区块哈希
func (b *BlockAggregate) ParentHash() valueobject.Hash {
	return b.parentHash
}

// Miner 返回矿工地址
func (b *BlockAggregate) Miner() valueobject.Address {
	return b.miner
}

// Timestamp 返回区块时间戳
func (b *BlockAggregate) Timestamp() time.Time {
	return b.timestamp
}

// SetGasInfo 设置燃料相关信息
func (b *BlockAggregate) SetGasInfo(gasLimit, gasUsed valueobject.Gas, baseFeePerGas *valueobject.Wei) error {
	if gasUsed.GreaterThan(gasLimit) {
		return errors.New("已使用燃料不能超过燃料限制")
	}

	b.gasLimit = gasLimit
	b.gasUsed = gasUsed
	b.baseFeePerGas = baseFeePerGas

	return nil
}

// GasLimit 返回燃料限制
func (b *BlockAggregate) GasLimit() valueobject.Gas {
	return b.gasLimit
}

// GasUsed 返回已使用燃料
func (b *BlockAggregate) GasUsed() valueobject.Gas {
	return b.gasUsed
}

// BaseFeePerGas 返回基础费用（EIP-1559）
func (b *BlockAggregate) BaseFeePerGas() *valueobject.Wei {
	return b.baseFeePerGas
}

// SetMerkleRoots 设置默克尔树根哈希
func (b *BlockAggregate) SetMerkleRoots(
	transactionsRoot, stateRoot, receiptsRoot valueobject.Hash,
) {
	b.transactionsRoot = transactionsRoot
	b.stateRoot = stateRoot
	b.receiptsRoot = receiptsRoot
}

// TransactionsRoot 返回交易默克尔树根
func (b *BlockAggregate) TransactionsRoot() valueobject.Hash {
	return b.transactionsRoot
}

// StateRoot 返回状态树根
func (b *BlockAggregate) StateRoot() valueobject.Hash {
	return b.stateRoot
}

// ReceiptsRoot 返回收据树根
func (b *BlockAggregate) ReceiptsRoot() valueobject.Hash {
	return b.receiptsRoot
}

// AddTransaction 添加交易摘要到区块
func (b *BlockAggregate) AddTransaction(summary TransactionSummary) {
	b.transactions = append(b.transactions, summary)
	b.transactionCount++

	// 发布交易添加事件
	event := NewTransactionAddedToBlockEvent(
		b.hash.String(),
		summary.Hash.String(),
		time.Now().Unix(),
	)
	b.AddEvent(event)
}

// Transactions 返回交易摘要列表
func (b *BlockAggregate) Transactions() []TransactionSummary {
	// 返回副本以保护内部状态
	result := make([]TransactionSummary, len(b.transactions))
	copy(result, b.transactions)
	return result
}

// TransactionCount 返回交易数量
func (b *BlockAggregate) TransactionCount() uint64 {
	return b.transactionCount
}

// IsEmpty 检查区块是否为空（无交易）
func (b *BlockAggregate) IsEmpty() bool {
	return b.transactionCount == 0
}

// SetWithdrawals 设置提款信息（EIP-4895）
func (b *BlockAggregate) SetWithdrawals(withdrawalsRoot *valueobject.Hash, withdrawals []Withdrawal) {
	b.withdrawalsRoot = withdrawalsRoot
	b.withdrawals = withdrawals
}

// Withdrawals 返回提款列表
func (b *BlockAggregate) Withdrawals() []Withdrawal {
	if b.withdrawals == nil {
		return nil
	}

	// 返回副本以保护内部状态
	result := make([]Withdrawal, len(b.withdrawals))
	copy(result, b.withdrawals)
	return result
}
