// Package aggregate 定义了交易相关的领域事件
package aggregate

// TransactionCreatedEvent 表示交易创建事件
type TransactionCreatedEvent struct {
	transactionHash string
	fromAddress     string
	occurredOn      int64
}

// NewTransactionCreatedEvent 创建交易创建事件
func NewTransactionCreatedEvent(transactionHash, fromAddress string, occurredOn int64) *TransactionCreatedEvent {
	return &TransactionCreatedEvent{
		transactionHash: transactionHash,
		fromAddress:     fromAddress,
		occurredOn:      occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionCreatedEvent) EventType() string {
	return "TransactionCreated"
}

// AggregateID 返回聚合根ID
func (e *TransactionCreatedEvent) AggregateID() string {
	return e.transactionHash
}

// OccurredOn 返回事件发生时间
func (e *TransactionCreatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// FromAddress 返回发送方地址
func (e *TransactionCreatedEvent) FromAddress() string {
	return e.fromAddress
}

// TransactionMinedEvent 表示交易挖矿事件
type TransactionMinedEvent struct {
	transactionHash string
	blockHash       string
	occurredOn      int64
}

// NewTransactionMinedEvent 创建交易挖矿事件
func NewTransactionMinedEvent(transactionHash, blockHash string, occurredOn int64) *TransactionMinedEvent {
	return &TransactionMinedEvent{
		transactionHash: transactionHash,
		blockHash:       blockHash,
		occurredOn:      occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionMinedEvent) EventType() string {
	return "TransactionMined"
}

// AggregateID 返回聚合根ID
func (e *TransactionMinedEvent) AggregateID() string {
	return e.transactionHash
}

// OccurredOn 返回事件发生时间
func (e *TransactionMinedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// BlockHash 返回区块哈希
func (e *TransactionMinedEvent) BlockHash() string {
	return e.blockHash
}

// TransactionFailedEvent 表示交易失败事件
type TransactionFailedEvent struct {
	transactionHash string
	reason          string
	occurredOn      int64
}

// NewTransactionFailedEvent 创建交易失败事件
func NewTransactionFailedEvent(transactionHash, reason string, occurredOn int64) *TransactionFailedEvent {
	return &TransactionFailedEvent{
		transactionHash: transactionHash,
		reason:          reason,
		occurredOn:      occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionFailedEvent) EventType() string {
	return "TransactionFailed"
}

// AggregateID 返回聚合根ID
func (e *TransactionFailedEvent) AggregateID() string {
	return e.transactionHash
}

// OccurredOn 返回事件发生时间
func (e *TransactionFailedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// Reason 返回失败原因
func (e *TransactionFailedEvent) Reason() string {
	return e.reason
}
