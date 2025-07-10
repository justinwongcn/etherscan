// Package valueobject 定义了领域层的值对象
package valueobject

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Gas 表示以太坊燃料的值对象
// 燃料用于衡量计算资源的消耗
// 这是一个值对象，具有不可变性和值相等性
type Gas struct {
	value *big.Int
}

// NewGas 创建一个新的Gas值对象
// 参数:
//   - gasStr: 燃料字符串，可以是十进制或十六进制（0x开头）
//
// 返回:
//   - Gas: 创建的Gas值对象
//   - error: 如果燃料格式无效，返回错误
func NewGas(gasStr string) (Gas, error) {
	if gasStr == "" {
		return Gas{}, errors.New("燃料值不能为空")
	}

	var value *big.Int
	var success bool

	// 处理十六进制格式
	if strings.HasPrefix(gasStr, "0x") {
		value, success = new(big.Int).SetString(gasStr[2:], 16)
	} else {
		// 处理十进制格式
		value, success = new(big.Int).SetString(gasStr, 10)
	}

	if !success {
		return Gas{}, fmt.Errorf("无效的燃料值格式: %s", gasStr)
	}

	// 燃料值不能为负数
	if value.Sign() < 0 {
		return Gas{}, errors.New("燃料值不能为负数")
	}

	return Gas{value: new(big.Int).Set(value)}, nil
}

// NewGasFromUint64 从uint64创建Gas值对象
func NewGasFromUint64(gas uint64) Gas {
	return Gas{value: new(big.Int).SetUint64(gas)}
}

// NewGasFromBigInt 从big.Int创建Gas值对象
func NewGasFromBigInt(gas *big.Int) (Gas, error) {
	if gas == nil {
		return Gas{}, errors.New("燃料值不能为nil")
	}

	if gas.Sign() < 0 {
		return Gas{}, errors.New("燃料值不能为负数")
	}

	return Gas{value: new(big.Int).Set(gas)}, nil
}

// String 返回燃料的十进制字符串表示
func (g Gas) String() string {
	return g.value.String()
}

// Hex 返回燃料的十六进制字符串表示（带0x前缀）
func (g Gas) Hex() string {
	return "0x" + g.value.Text(16)
}

// Uint64 返回燃料的uint64表示
// 如果值超出uint64范围，返回错误
func (g Gas) Uint64() (uint64, error) {
	if !g.value.IsUint64() {
		return 0, errors.New("燃料值超出uint64范围")
	}
	return g.value.Uint64(), nil
}

// BigInt 返回燃料的big.Int表示（副本）
func (g Gas) BigInt() *big.Int {
	return new(big.Int).Set(g.value)
}

// Add 加法运算，返回新的Gas值对象
func (g Gas) Add(other Gas) Gas {
	result := new(big.Int).Add(g.value, other.value)
	return Gas{value: result}
}

// Sub 减法运算，返回新的Gas值对象
// 如果结果为负数，返回错误
func (g Gas) Sub(other Gas) (Gas, error) {
	result := new(big.Int).Sub(g.value, other.value)
	if result.Sign() < 0 {
		return Gas{}, errors.New("燃料减法结果不能为负数")
	}
	return Gas{value: result}, nil
}

// Mul 乘法运算，返回新的Gas值对象
func (g Gas) Mul(multiplier uint64) Gas {
	result := new(big.Int).Mul(g.value, new(big.Int).SetUint64(multiplier))
	return Gas{value: result}
}

// Equals 检查两个Gas是否相等
func (g Gas) Equals(other Gas) bool {
	return g.value.Cmp(other.value) == 0
}

// GreaterThan 检查是否大于另一个Gas值
func (g Gas) GreaterThan(other Gas) bool {
	return g.value.Cmp(other.value) > 0
}

// LessThan 检查是否小于另一个Gas值
func (g Gas) LessThan(other Gas) bool {
	return g.value.Cmp(other.value) < 0
}

// IsZero 检查是否为零
func (g Gas) IsZero() bool {
	return g.value.Sign() == 0
}
