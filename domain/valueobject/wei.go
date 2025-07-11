// Package valueobject 定义了领域层的值对象
package valueobject

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Wei 表示以太坊最小货币单位的值对象
// 1 ETH = 10^18 Wei
// 这是一个值对象，具有不可变性和值相等性
type Wei struct {
	value *big.Int
}

// 常用的以太坊单位转换常量
var (
	WeiPerEther = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
	WeiPerGwei  = new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)  // 10^9
)

// NewWei 创建一个新的Wei值对象
// 参数:
//   - weiStr: Wei字符串，可以是十进制或十六进制（0x开头）
//
// 返回:
//   - Wei: 创建的Wei值对象
//   - error: 如果格式无效，返回错误
func NewWei(weiStr string) (Wei, error) {
	if weiStr == "" {
		return Wei{}, errors.New("Wei值不能为空")
	}

	var value *big.Int
	var success bool

	// 处理十六进制格式
	if strings.HasPrefix(weiStr, "0x") {
		value, success = new(big.Int).SetString(weiStr[2:], 16)
	} else {
		// 处理十进制格式
		value, success = new(big.Int).SetString(weiStr, 10)
	}

	if !success {
		return Wei{}, fmt.Errorf("无效的Wei值格式: %s", weiStr)
	}

	// Wei值不能为负数
	if value.Sign() < 0 {
		return Wei{}, errors.New("Wei值不能为负数")
	}

	return Wei{value: new(big.Int).Set(value)}, nil
}

// NewWeiFromUint64 从uint64创建Wei值对象
func NewWeiFromUint64(wei uint64) Wei {
	return Wei{value: new(big.Int).SetUint64(wei)}
}

// NewWeiFromBigInt 从big.Int创建Wei值对象
func NewWeiFromBigInt(wei *big.Int) (Wei, error) {
	if wei == nil {
		return Wei{}, errors.New("Wei值不能为nil")
	}

	if wei.Sign() < 0 {
		return Wei{}, errors.New("Wei值不能为负数")
	}

	return Wei{value: new(big.Int).Set(wei)}, nil
}

// NewWeiFromEther 从以太币数量创建Wei值对象
// 参数:
//   - etherStr: 以太币数量的字符串表示（可以包含小数）
//
// 返回:
//   - Wei: 创建的Wei值对象
//   - error: 如果格式无效，返回错误
func NewWeiFromEther(etherStr string) (Wei, error) {
	if etherStr == "" {
		return Wei{}, errors.New("以太币值不能为空")
	}

	// 简化实现：假设输入是整数以太币
	// 实际实现应该处理小数点
	etherValue, success := new(big.Int).SetString(etherStr, 10)
	if !success {
		return Wei{}, fmt.Errorf("无效的以太币值格式: %s", etherStr)
	}

	if etherValue.Sign() < 0 {
		return Wei{}, errors.New("以太币值不能为负数")
	}

	weiValue := new(big.Int).Mul(etherValue, WeiPerEther)
	return Wei{value: weiValue}, nil
}

// String 返回Wei的十进制字符串表示
func (w Wei) String() string {
	return w.value.String()
}

// Hex 返回Wei的十六进制字符串表示（带0x前缀）
func (w Wei) Hex() string {
	return "0x" + w.value.Text(16)
}

// Uint64 返回Wei的uint64表示
// 如果值超出uint64范围，返回错误
func (w Wei) Uint64() (uint64, error) {
	if !w.value.IsUint64() {
		return 0, errors.New("Wei值超出uint64范围")
	}
	return w.value.Uint64(), nil
}

// BigInt 返回Wei的big.Int表示（副本）
func (w Wei) BigInt() *big.Int {
	return new(big.Int).Set(w.value)
}

// ToEther 转换为以太币表示（返回字符串，保留小数）
func (w Wei) ToEther() string {
	// 简化实现：只返回整数部分
	// 实际实现应该处理小数部分
	etherValue := new(big.Int).Div(w.value, WeiPerEther)
	return etherValue.String()
}

// ToGwei 转换为Gwei表示
func (w Wei) ToGwei() string {
	gweiValue := new(big.Int).Div(w.value, WeiPerGwei)
	return gweiValue.String()
}

// Add 加法运算，返回新的Wei值对象
func (w Wei) Add(other Wei) Wei {
	result := new(big.Int).Add(w.value, other.value)
	return Wei{value: result}
}

// Sub 减法运算，返回新的Wei值对象
// 如果结果为负数，返回错误
func (w Wei) Sub(other Wei) (Wei, error) {
	result := new(big.Int).Sub(w.value, other.value)
	if result.Sign() < 0 {
		return Wei{}, errors.New("Wei减法结果不能为负数")
	}
	return Wei{value: result}, nil
}

// Mul 乘法运算，返回新的Wei值对象
func (w Wei) Mul(multiplier uint64) Wei {
	result := new(big.Int).Mul(w.value, new(big.Int).SetUint64(multiplier))
	return Wei{value: result}
}

// Div 除法运算，返回新的Wei值对象
func (w Wei) Div(divisor uint64) (Wei, error) {
	if divisor == 0 {
		return Wei{}, errors.New("除数不能为零")
	}

	result := new(big.Int).Div(w.value, new(big.Int).SetUint64(divisor))
	return Wei{value: result}, nil
}

// Equals 检查两个Wei是否相等
func (w Wei) Equals(other Wei) bool {
	return w.value.Cmp(other.value) == 0
}

// GreaterThan 检查是否大于另一个Wei值
func (w Wei) GreaterThan(other Wei) bool {
	return w.value.Cmp(other.value) > 0
}

// LessThan 检查是否小于另一个Wei值
func (w Wei) LessThan(other Wei) bool {
	return w.value.Cmp(other.value) < 0
}

// IsZero 检查是否为零
func (w Wei) IsZero() bool {
	return w.value.Sign() == 0
}
