// Package aggregate 定义了账户聚合根
package aggregate

import (
	"errors"
	"time"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// AccountType 表示账户类型
type AccountType int

const (
	ExternallyOwnedAccount AccountType = iota // 外部拥有账户（EOA）
	ContractAccount                           // 合约账户
)

// AccountAggregate 表示账户聚合根
// 账户是以太坊网络中的基本实体，可以持有ETH并执行交易
// 作为聚合根，它负责维护账户状态的一致性
type AccountAggregate struct {
	BaseAggregateRoot

	// 基本账户信息
	address     valueobject.Address
	balance     valueobject.Wei
	nonce       uint64
	accountType AccountType

	// 合约相关信息（仅合约账户）
	codeHash    *valueobject.Hash
	storageRoot *valueobject.Hash

	// 账户历史信息
	createdAt   time.Time
	lastUpdated time.Time

	// 交易历史摘要
	transactionCount uint64
	totalSent        valueobject.Wei
	totalReceived    valueobject.Wei
}

// NewExternallyOwnedAccount 创建外部拥有账户
func NewExternallyOwnedAccount(address valueobject.Address) (*AccountAggregate, error) {
	if address.IsZero() {
		return nil, errors.New("账户地址不能为零地址")
	}

	now := time.Now()
	account := &AccountAggregate{
		BaseAggregateRoot: NewBaseAggregateRoot(address.String()),
		address:           address,
		balance:           valueobject.NewWeiFromUint64(0),
		nonce:             0,
		accountType:       ExternallyOwnedAccount,
		createdAt:         now,
		lastUpdated:       now,
		totalSent:         valueobject.NewWeiFromUint64(0),
		totalReceived:     valueobject.NewWeiFromUint64(0),
	}

	// 发布账户创建事件
	event := NewAccountCreatedEvent(address.String(), ExternallyOwnedAccount, now.Unix())
	account.AddEvent(event)

	return account, nil
}

// NewContractAccount 创建合约账户
func NewContractAccount(
	address valueobject.Address,
	codeHash valueobject.Hash,
	storageRoot valueobject.Hash,
) (*AccountAggregate, error) {
	if address.IsZero() {
		return nil, errors.New("合约地址不能为零地址")
	}

	if codeHash.IsZero() {
		return nil, errors.New("合约代码哈希不能为零")
	}

	now := time.Now()
	account := &AccountAggregate{
		BaseAggregateRoot: NewBaseAggregateRoot(address.String()),
		address:           address,
		balance:           valueobject.NewWeiFromUint64(0),
		nonce:             0,
		accountType:       ContractAccount,
		codeHash:          &codeHash,
		storageRoot:       &storageRoot,
		createdAt:         now,
		lastUpdated:       now,
		totalSent:         valueobject.NewWeiFromUint64(0),
		totalReceived:     valueobject.NewWeiFromUint64(0),
	}

	// 发布合约账户创建事件
	event := NewAccountCreatedEvent(address.String(), ContractAccount, now.Unix())
	account.AddEvent(event)

	return account, nil
}

// Address 返回账户地址
func (a *AccountAggregate) Address() valueobject.Address {
	return a.address
}

// Balance 返回账户余额
func (a *AccountAggregate) Balance() valueobject.Wei {
	return a.balance
}

// Nonce 返回账户nonce
func (a *AccountAggregate) Nonce() uint64 {
	return a.nonce
}

// Type 返回账户类型
func (a *AccountAggregate) Type() AccountType {
	return a.accountType
}

// CodeHash 返回合约代码哈希（仅合约账户）
func (a *AccountAggregate) CodeHash() *valueobject.Hash {
	return a.codeHash
}

// StorageRoot 返回存储根哈希（仅合约账户）
func (a *AccountAggregate) StorageRoot() *valueobject.Hash {
	return a.storageRoot
}

// UpdateBalance 更新账户余额
func (a *AccountAggregate) UpdateBalance(newBalance valueobject.Wei) {
	oldBalance := a.balance
	a.balance = newBalance
	a.lastUpdated = time.Now()

	// 发布余额更新事件
	event := NewBalanceUpdatedEvent(
		a.address.String(),
		oldBalance.String(),
		newBalance.String(),
		a.lastUpdated.Unix(),
	)
	a.AddEvent(event)
}

// IncrementNonce 增加nonce值
func (a *AccountAggregate) IncrementNonce() {
	a.nonce++
	a.lastUpdated = time.Now()

	// 发布nonce更新事件
	event := NewNonceUpdatedEvent(
		a.address.String(),
		a.nonce,
		a.lastUpdated.Unix(),
	)
	a.AddEvent(event)
}

// RecordSentTransaction 记录发送的交易
func (a *AccountAggregate) RecordSentTransaction(amount valueobject.Wei) error {
	// 检查余额是否足够
	if a.balance.LessThan(amount) {
		return errors.New("余额不足")
	}

	// 更新余额和统计信息
	newBalance, err := a.balance.Sub(amount)
	if err != nil {
		return err
	}

	a.balance = newBalance
	a.totalSent = a.totalSent.Add(amount)
	a.transactionCount++
	a.lastUpdated = time.Now()

	// 发布交易发送事件
	event := NewTransactionSentEvent(
		a.address.String(),
		amount.String(),
		a.lastUpdated.Unix(),
	)
	a.AddEvent(event)

	return nil
}

// RecordReceivedTransaction 记录接收的交易
func (a *AccountAggregate) RecordReceivedTransaction(amount valueobject.Wei) {
	// 更新余额和统计信息
	a.balance = a.balance.Add(amount)
	a.totalReceived = a.totalReceived.Add(amount)
	a.lastUpdated = time.Now()

	// 发布交易接收事件
	event := NewTransactionReceivedEvent(
		a.address.String(),
		amount.String(),
		a.lastUpdated.Unix(),
	)
	a.AddEvent(event)
}

// TransactionCount 返回交易数量
func (a *AccountAggregate) TransactionCount() uint64 {
	return a.transactionCount
}

// TotalSent 返回总发送金额
func (a *AccountAggregate) TotalSent() valueobject.Wei {
	return a.totalSent
}

// TotalReceived 返回总接收金额
func (a *AccountAggregate) TotalReceived() valueobject.Wei {
	return a.totalReceived
}

// IsContract 检查是否为合约账户
func (a *AccountAggregate) IsContract() bool {
	return a.accountType == ContractAccount
}

// IsExternallyOwned 检查是否为外部拥有账户
func (a *AccountAggregate) IsExternallyOwned() bool {
	return a.accountType == ExternallyOwnedAccount
}

// HasSufficientBalance 检查是否有足够余额
func (a *AccountAggregate) HasSufficientBalance(amount valueobject.Wei) bool {
	return a.balance.GreaterThan(amount) || a.balance.Equals(amount)
}

// CreatedAt 返回创建时间
func (a *AccountAggregate) CreatedAt() time.Time {
	return a.createdAt
}

// LastUpdated 返回最后更新时间
func (a *AccountAggregate) LastUpdated() time.Time {
	return a.lastUpdated
}
