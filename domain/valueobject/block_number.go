// Package valueobject 定义了领域层的值对象
package valueobject

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// BlockNumber 表示以太坊区块号的值对象
// 区块号可以是具体的数字，也可以是特殊标识符（latest、earliest、pending）
// 这是一个值对象，具有不可变性和值相等性
type BlockNumber struct {
	value    *big.Int
	tag      string
	isNumber bool
}

// 预定义的区块标签常量
const (
	BlockTagLatest   = "latest"
	BlockTagEarliest = "earliest"
	BlockTagPending  = "pending"
)

// NewBlockNumber 创建一个新的BlockNumber值对象
// 参数:
//   - numberStr: 区块号字符串，可以是数字（十进制或十六进制）或特殊标识符
//
// 返回:
//   - BlockNumber: 创建的BlockNumber值对象
//   - error: 如果格式无效，返回错误
func NewBlockNumber(numberStr string) (BlockNumber, error) {
	if numberStr == "" {
		return BlockNumber{}, errors.New("区块号不能为空")
	}

	// 检查是否为特殊标签
	switch numberStr {
	case BlockTagLatest, BlockTagEarliest, BlockTagPending:
		return BlockNumber{
			tag:      numberStr,
			isNumber: false,
		}, nil
	}

	// 尝试解析为数字
	var value *big.Int
	var success bool

	// 处理十六进制格式
	if strings.HasPrefix(numberStr, "0x") {
		value, success = new(big.Int).SetString(numberStr[2:], 16)
	} else {
		// 处理十进制格式
		value, success = new(big.Int).SetString(numberStr, 10)
	}

	if !success {
		return BlockNumber{}, fmt.Errorf("无效的区块号格式: %s", numberStr)
	}

	// 区块号不能为负数
	if value.Sign() < 0 {
		return BlockNumber{}, errors.New("区块号不能为负数")
	}

	return BlockNumber{
		value:    value,
		isNumber: true,
	}, nil
}

// NewBlockNumberFromUint64 从uint64创建BlockNumber值对象
func NewBlockNumberFromUint64(number uint64) BlockNumber {
	return BlockNumber{
		value:    new(big.Int).SetUint64(number),
		isNumber: true,
	}
}

// NewBlockNumberFromBigInt 从big.Int创建BlockNumber值对象
func NewBlockNumberFromBigInt(number *big.Int) (BlockNumber, error) {
	if number == nil {
		return BlockNumber{}, errors.New("区块号不能为nil")
	}

	if number.Sign() < 0 {
		return BlockNumber{}, errors.New("区块号不能为负数")
	}

	return BlockNumber{
		value:    new(big.Int).Set(number),
		isNumber: true,
	}, nil
}

// NewLatestBlockNumber 创建表示最新区块的BlockNumber
func NewLatestBlockNumber() BlockNumber {
	return BlockNumber{
		tag:      BlockTagLatest,
		isNumber: false,
	}
}

// NewEarliestBlockNumber 创建表示最早区块的BlockNumber
func NewEarliestBlockNumber() BlockNumber {
	return BlockNumber{
		tag:      BlockTagEarliest,
		isNumber: false,
	}
}

// NewPendingBlockNumber 创建表示待处理区块的BlockNumber
func NewPendingBlockNumber() BlockNumber {
	return BlockNumber{
		tag:      BlockTagPending,
		isNumber: false,
	}
}

// String 返回区块号的字符串表示
func (bn BlockNumber) String() string {
	if bn.isNumber {
		return bn.value.String()
	}
	return bn.tag
}

// Hex 返回区块号的十六进制字符串表示（仅对数字有效）
func (bn BlockNumber) Hex() string {
	if !bn.isNumber {
		return bn.tag
	}
	return "0x" + bn.value.Text(16)
}

// IsNumber 检查是否为具体的数字
func (bn BlockNumber) IsNumber() bool {
	return bn.isNumber
}

// IsTag 检查是否为特殊标签
func (bn BlockNumber) IsTag() bool {
	return !bn.isNumber
}

// Uint64 返回区块号的uint64表示（仅对数字有效）
func (bn BlockNumber) Uint64() (uint64, error) {
	if !bn.isNumber {
		return 0, errors.New("特殊标签不能转换为uint64")
	}

	if !bn.value.IsUint64() {
		return 0, errors.New("区块号超出uint64范围")
	}

	return bn.value.Uint64(), nil
}

// BigInt 返回区块号的big.Int表示（仅对数字有效）
func (bn BlockNumber) BigInt() (*big.Int, error) {
	if !bn.isNumber {
		return nil, errors.New("特殊标签不能转换为big.Int")
	}

	return new(big.Int).Set(bn.value), nil
}

// Tag 返回特殊标签（仅对标签有效）
func (bn BlockNumber) Tag() (string, error) {
	if bn.isNumber {
		return "", errors.New("数字区块号没有标签")
	}

	return bn.tag, nil
}

// Equals 检查两个BlockNumber是否相等
func (bn BlockNumber) Equals(other BlockNumber) bool {
	if bn.isNumber != other.isNumber {
		return false
	}

	if bn.isNumber {
		return bn.value.Cmp(other.value) == 0
	}

	return bn.tag == other.tag
}

// Next 返回下一个区块号（仅对数字有效）
func (bn BlockNumber) Next() (BlockNumber, error) {
	if !bn.isNumber {
		return BlockNumber{}, errors.New("特殊标签不能计算下一个区块号")
	}

	nextValue := new(big.Int).Add(bn.value, big.NewInt(1))
	return BlockNumber{
		value:    nextValue,
		isNumber: true,
	}, nil
}

// Previous 返回上一个区块号（仅对数字有效）
func (bn BlockNumber) Previous() (BlockNumber, error) {
	if !bn.isNumber {
		return BlockNumber{}, errors.New("特殊标签不能计算上一个区块号")
	}

	if bn.value.Sign() == 0 {
		return BlockNumber{}, errors.New("区块号0没有上一个区块")
	}

	prevValue := new(big.Int).Sub(bn.value, big.NewInt(1))
	return BlockNumber{
		value:    prevValue,
		isNumber: true,
	}, nil
}
