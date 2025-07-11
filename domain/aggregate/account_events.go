// Package aggregate 定义了账户相关的领域事件
package aggregate

// AccountCreatedEvent 表示账户创建事件
type AccountCreatedEvent struct {
	address     string
	accountType AccountType
	occurredOn  int64
}

// NewAccountCreatedEvent 创建账户创建事件
func NewAccountCreatedEvent(address string, accountType AccountType, occurredOn int64) *AccountCreatedEvent {
	return &AccountCreatedEvent{
		address:     address,
		accountType: accountType,
		occurredOn:  occurredOn,
	}
}

// EventType 返回事件类型
func (e *AccountCreatedEvent) EventType() string {
	return "AccountCreated"
}

// AggregateID 返回聚合根ID
func (e *AccountCreatedEvent) AggregateID() string {
	return e.address
}

// OccurredOn 返回事件发生时间
func (e *AccountCreatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// AccountType 返回账户类型
func (e *AccountCreatedEvent) AccountType() AccountType {
	return e.accountType
}

// BalanceUpdatedEvent 表示余额更新事件
type BalanceUpdatedEvent struct {
	address    string
	oldBalance string
	newBalance string
	occurredOn int64
}

// NewBalanceUpdatedEvent 创建余额更新事件
func NewBalanceUpdatedEvent(address, oldBalance, newBalance string, occurredOn int64) *BalanceUpdatedEvent {
	return &BalanceUpdatedEvent{
		address:    address,
		oldBalance: oldBalance,
		newBalance: newBalance,
		occurredOn: occurredOn,
	}
}

// EventType 返回事件类型
func (e *BalanceUpdatedEvent) EventType() string {
	return "BalanceUpdated"
}

// AggregateID 返回聚合根ID
func (e *BalanceUpdatedEvent) AggregateID() string {
	return e.address
}

// OccurredOn 返回事件发生时间
func (e *BalanceUpdatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// OldBalance 返回旧余额
func (e *BalanceUpdatedEvent) OldBalance() string {
	return e.oldBalance
}

// NewBalance 返回新余额
func (e *BalanceUpdatedEvent) NewBalance() string {
	return e.newBalance
}

// NonceUpdatedEvent 表示nonce更新事件
type NonceUpdatedEvent struct {
	address    string
	newNonce   uint64
	occurredOn int64
}

// NewNonceUpdatedEvent 创建nonce更新事件
func NewNonceUpdatedEvent(address string, newNonce uint64, occurredOn int64) *NonceUpdatedEvent {
	return &NonceUpdatedEvent{
		address:    address,
		newNonce:   newNonce,
		occurredOn: occurredOn,
	}
}

// EventType 返回事件类型
func (e *NonceUpdatedEvent) EventType() string {
	return "NonceUpdated"
}

// AggregateID 返回聚合根ID
func (e *NonceUpdatedEvent) AggregateID() string {
	return e.address
}

// OccurredOn 返回事件发生时间
func (e *NonceUpdatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// NewNonce 返回新nonce值
func (e *NonceUpdatedEvent) NewNonce() uint64 {
	return e.newNonce
}

// TransactionSentEvent 表示交易发送事件
type TransactionSentEvent struct {
	address    string
	amount     string
	occurredOn int64
}

// NewTransactionSentEvent 创建交易发送事件
func NewTransactionSentEvent(address, amount string, occurredOn int64) *TransactionSentEvent {
	return &TransactionSentEvent{
		address:    address,
		amount:     amount,
		occurredOn: occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionSentEvent) EventType() string {
	return "TransactionSent"
}

// AggregateID 返回聚合根ID
func (e *TransactionSentEvent) AggregateID() string {
	return e.address
}

// OccurredOn 返回事件发生时间
func (e *TransactionSentEvent) OccurredOn() int64 {
	return e.occurredOn
}

// Amount 返回发送金额
func (e *TransactionSentEvent) Amount() string {
	return e.amount
}

// TransactionReceivedEvent 表示交易接收事件
type TransactionReceivedEvent struct {
	address    string
	amount     string
	occurredOn int64
}

// NewTransactionReceivedEvent 创建交易接收事件
func NewTransactionReceivedEvent(address, amount string, occurredOn int64) *TransactionReceivedEvent {
	return &TransactionReceivedEvent{
		address:    address,
		amount:     amount,
		occurredOn: occurredOn,
	}
}

// EventType 返回事件类型
func (e *TransactionReceivedEvent) EventType() string {
	return "TransactionReceived"
}

// AggregateID 返回聚合根ID
func (e *TransactionReceivedEvent) AggregateID() string {
	return e.address
}

// OccurredOn 返回事件发生时间
func (e *TransactionReceivedEvent) OccurredOn() int64 {
	return e.occurredOn
}

// Amount 返回接收金额
func (e *TransactionReceivedEvent) Amount() string {
	return e.amount
}
