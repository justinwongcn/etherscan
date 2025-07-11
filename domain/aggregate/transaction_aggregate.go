// Package aggregate 定义了交易聚合根
package aggregate

import (
	"errors"
	"time"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// TransactionType 表示交易类型的枚举
type TransactionType int

const (
	LegacyTransaction  TransactionType = iota // 传统交易
	EIP2930Transaction                        // EIP-2930 访问列表交易
	EIP1559Transaction                        // EIP-1559 动态费用交易
)

// TransactionStatus 表示交易状态
type TransactionStatus int

const (
	TransactionPending TransactionStatus = iota // 待处理
	TransactionMined                            // 已挖矿
	TransactionFailed                           // 失败
)

// TransactionAggregate 表示交易聚合根
// 交易是以太坊网络中状态变更的基本单位
// 作为聚合根，它负责维护交易相关数据的一致性
type TransactionAggregate struct {
	BaseAggregateRoot

	// 基本交易信息
	hash  valueobject.Hash
	from  valueobject.Address
	to    *valueobject.Address // 可能为nil（合约创建交易）
	value valueobject.Wei
	nonce uint64
	input []byte

	// 燃料相关
	gas            valueobject.Gas
	gasPrice       *valueobject.Wei // 传统交易使用
	maxFeePerGas   *valueobject.Wei // EIP-1559交易使用
	maxPriorityFee *valueobject.Wei // EIP-1559交易使用

	// 区块相关信息
	blockHash   *valueobject.Hash
	blockNumber *valueobject.BlockNumber
	txIndex     *uint64

	// 交易类型和状态
	txType TransactionType
	status TransactionStatus

	// 签名信息
	signature *valueobject.Signature // 交易签名

	// 交易收据（如果已挖矿）
	receipt *TransactionReceipt

	// 创建时间
	createdAt time.Time
}

// TransactionReceipt 表示交易收据
type TransactionReceipt struct {
	GasUsed           valueobject.Gas
	CumulativeGasUsed valueobject.Gas
	ContractAddress   *valueobject.Address // 合约创建交易的合约地址
	Logs              []Log
	Status            uint64           // 1表示成功，0表示失败
	EffectiveGasPrice *valueobject.Wei // EIP-1559实际支付的燃料价格
}

// Log 表示交易日志
type Log struct {
	Address valueobject.Address
	Topics  []valueobject.Hash
	Data    []byte
}

// NewTransactionAggregate 创建新的交易聚合根
func NewTransactionAggregate(
	hash valueobject.Hash,
	from valueobject.Address,
	to *valueobject.Address,
	value valueobject.Wei,
	nonce uint64,
	gas valueobject.Gas,
	input []byte,
	signature *valueobject.Signature,
) (*TransactionAggregate, error) {

	// 验证必要字段
	if hash.IsZero() {
		return nil, errors.New("交易哈希不能为零")
	}

	if from.IsZero() {
		return nil, errors.New("发送方地址不能为零")
	}

	if gas.IsZero() {
		return nil, errors.New("燃料限制不能为零")
	}

	tx := &TransactionAggregate{
		BaseAggregateRoot: NewBaseAggregateRoot(hash.String()),
		hash:              hash,
		from:              from,
		to:                to,
		value:             value,
		nonce:             nonce,
		gas:               gas,
		input:             input,
		signature:         signature,
		status:            TransactionPending,
		createdAt:         time.Now(),
	}

	// 发布交易创建事件
	event := NewTransactionCreatedEvent(hash.String(), from.String(), time.Now().Unix())
	tx.AddEvent(event)

	return tx, nil
}

// NewTransactionAggregateWithoutSignature 创建没有签名的交易聚合根（用于测试）
func NewTransactionAggregateWithoutSignature(
	hash valueobject.Hash,
	from valueobject.Address,
	to *valueobject.Address,
	value valueobject.Wei,
	nonce uint64,
	gas valueobject.Gas,
	input []byte,
) (*TransactionAggregate, error) {
	return NewTransactionAggregate(hash, from, to, value, nonce, gas, input, nil)
}

// Hash 返回交易哈希
func (t *TransactionAggregate) Hash() valueobject.Hash {
	return t.hash
}

// From 返回发送方地址
func (t *TransactionAggregate) From() valueobject.Address {
	return t.from
}

// To 返回接收方地址
func (t *TransactionAggregate) To() *valueobject.Address {
	return t.to
}

// Value 返回交易金额
func (t *TransactionAggregate) Value() valueobject.Wei {
	return t.value
}

// Nonce 返回交易nonce
func (t *TransactionAggregate) Nonce() uint64 {
	return t.nonce
}

// Gas 返回燃料限制
func (t *TransactionAggregate) Gas() valueobject.Gas {
	return t.gas
}

// Input 返回交易输入数据
func (t *TransactionAggregate) Input() []byte {
	// 返回副本以保护内部状态
	result := make([]byte, len(t.input))
	copy(result, t.input)
	return result
}

// SetLegacyGasPrice 设置传统交易的燃料价格
func (t *TransactionAggregate) SetLegacyGasPrice(gasPrice valueobject.Wei) {
	t.gasPrice = &gasPrice
	t.txType = LegacyTransaction
}

// SetEIP1559Fees 设置EIP-1559交易的费用
func (t *TransactionAggregate) SetEIP1559Fees(maxFeePerGas, maxPriorityFee valueobject.Wei) error {
	if maxPriorityFee.GreaterThan(maxFeePerGas) {
		return errors.New("优先费用不能超过最大费用")
	}

	t.maxFeePerGas = &maxFeePerGas
	t.maxPriorityFee = &maxPriorityFee
	t.txType = EIP1559Transaction

	return nil
}

// GasPrice 返回燃料价格（传统交易）
func (t *TransactionAggregate) GasPrice() *valueobject.Wei {
	return t.gasPrice
}

// MaxFeePerGas 返回最大费用（EIP-1559交易）
func (t *TransactionAggregate) MaxFeePerGas() *valueobject.Wei {
	return t.maxFeePerGas
}

// MaxPriorityFee 返回优先费用（EIP-1559交易）
func (t *TransactionAggregate) MaxPriorityFee() *valueobject.Wei {
	return t.maxPriorityFee
}

// Type 返回交易类型
func (t *TransactionAggregate) Type() TransactionType {
	return t.txType
}

// Status 返回交易状态
func (t *TransactionAggregate) Status() TransactionStatus {
	return t.status
}

// MarkAsMined 标记交易为已挖矿状态
func (t *TransactionAggregate) MarkAsMined(
	blockHash valueobject.Hash,
	blockNumber valueobject.BlockNumber,
	txIndex uint64,
	receipt TransactionReceipt,
) {
	t.blockHash = &blockHash
	t.blockNumber = &blockNumber
	t.txIndex = &txIndex
	t.receipt = &receipt
	t.status = TransactionMined

	// 发布交易挖矿事件
	event := NewTransactionMinedEvent(
		t.hash.String(),
		blockHash.String(),
		time.Now().Unix(),
	)
	t.AddEvent(event)
}

// MarkAsFailed 标记交易为失败状态
func (t *TransactionAggregate) MarkAsFailed(reason string) {
	t.status = TransactionFailed

	// 发布交易失败事件
	event := NewTransactionFailedEvent(
		t.hash.String(),
		reason,
		time.Now().Unix(),
	)
	t.AddEvent(event)
}

// BlockHash 返回所在区块哈希
func (t *TransactionAggregate) BlockHash() *valueobject.Hash {
	return t.blockHash
}

// BlockNumber 返回所在区块号
func (t *TransactionAggregate) BlockNumber() *valueobject.BlockNumber {
	return t.blockNumber
}

// TransactionIndex 返回在区块中的索引
func (t *TransactionAggregate) TransactionIndex() *uint64 {
	return t.txIndex
}

// Receipt 返回交易收据
func (t *TransactionAggregate) Receipt() *TransactionReceipt {
	return t.receipt
}

// IsContractCreation 检查是否为合约创建交易
func (t *TransactionAggregate) IsContractCreation() bool {
	return t.to == nil
}

// IsPending 检查是否为待处理状态
func (t *TransactionAggregate) IsPending() bool {
	return t.status == TransactionPending
}

// IsMined 检查是否已挖矿
func (t *TransactionAggregate) IsMined() bool {
	return t.status == TransactionMined
}

// IsFailed 检查是否失败
func (t *TransactionAggregate) IsFailed() bool {
	return t.status == TransactionFailed
}

// GetSignature 返回交易的签名信息
func (t *TransactionAggregate) GetSignature() (v, r, s string) {
	if t.signature == nil {
		return "", "", ""
	}
	return t.signature.V(), t.signature.R(), t.signature.S()
}

// Signature 返回交易的签名值对象
func (t *TransactionAggregate) Signature() *valueobject.Signature {
	return t.signature
}
