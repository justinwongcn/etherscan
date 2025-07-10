// Package aggregate 定义了区块相关的领域事件
package aggregate

import (
	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// BlockCreatedEvent 表示区块创建事件
type BlockCreatedEvent struct {
	blockHash   string
	blockNumber valueobject.BlockNumber
	occurredOn  int64
}

// NewBlockCreatedEvent 创建区块创建事件
func NewBlockCreatedEvent(blockHash string, blockNumber valueobject.BlockNumber, occurredOn int64) *BlockCreatedEvent {
	return &BlockCreatedEvent{
		blockHash:   blockHash,
		blockNumber: blockNumber,
		occurredOn:  occurredOn,
	}
}

// EventType 返回事件类型
func (e *BlockCreatedEvent) EventType() string {
	return "BlockCreated"
}

// AggregateID 返回聚合根ID
func (e *BlockCreatedEvent) AggregateID() string {
	return e.blockHash
}

// OccurredOn 返回事件发生时间
func (e *BlockCreatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// BlockNumber 返回区块号
func (e *BlockCreatedEvent) BlockNumber() valueobject.BlockNumber {
	return e.blockNumber
}

// TransactionAddedToBlockEvent 表示交易添加到区块的事件
type TransactionAddedToBlockEvent struct {
	blockHash       string
	transactionHash string
	occurredOn      int64
}

// NewTransactionAddedToBlockEvent 创建交易添加到区块事件
func NewTransactionAddedToBlockEvent(blockHash, transactionHash string, occurredOn int64) *TransactionAddedToBlockEvent {
	return &TransactionAddedToBlockEvent{
		blockHash:       blockHash,
		transactionHash: transactionHash,
		occurredOn:      occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionAddedToBlockEvent) EventType() string {
	return "TransactionAddedToBlock"
}

// AggregateID 返回聚合根ID
func (e *TransactionAddedToBlockEvent) AggregateID() string {
	return e.blockHash
}

// OccurredOn 返回事件发生时间
func (e *TransactionAddedToBlockEvent) OccurredOn() int64 {
	return e.occurredOn
}

// TransactionHash 返回交易哈希
func (e *TransactionAddedToBlockEvent) TransactionHash() string {
	return e.transactionHash
}
