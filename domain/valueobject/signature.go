// Package valueobject 定义了DDD中的值对象
// 值对象是不可变的，具有相等性比较的特征
package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// Signature 表示以太坊交易的数字签名
// 包含ECDSA签名的三个组成部分：v, r, s
type Signature struct {
	v string // 恢复ID，用于从签名中恢复公钥
	r string // 签名的r值
	s string // 签名的s值
}

// NewSignature 创建一个新的签名值对象
// 参数:
//   - v: 恢复ID（通常是27, 28或更高的值用于EIP-155）
//   - r: 签名的r值（32字节的十六进制字符串）
//   - s: 签名的s值（32字节的十六进制字符串）
//
// 返回:
//   - *Signature: 签名值对象实例
//   - error: 如果参数无效则返回错误
func NewSignature(v, r, s string) (*Signature, error) {
	// 验证参数不能为空
	if v == "" || r == "" || s == "" {
		return nil, errors.New("签名参数不能为空")
	}

	// 验证r和s不能为零值
	if r == "0x0" || s == "0x0" {
		return nil, errors.New("签名参数r和s不能为零")
	}

	// 验证十六进制格式
	if err := validateHexString(r); err != nil {
		return nil, fmt.Errorf("无效的r值: %w", err)
	}

	if err := validateHexString(s); err != nil {
		return nil, fmt.Errorf("无效的s值: %w", err)
	}

	return &Signature{
		v: v,
		r: r,
		s: s,
	}, nil
}

// NewSignatureFromStrings 从字符串创建签名（不进行严格验证）
// 用于从外部数据源创建签名对象
func NewSignatureFromStrings(v, r, s string) *Signature {
	return &Signature{
		v: v,
		r: r,
		s: s,
	}
}

// V 返回恢复ID
func (sig *Signature) V() string {
	return sig.v
}

// R 返回签名的r值
func (sig *Signature) R() string {
	return sig.r
}

// S 返回签名的s值
func (sig *Signature) S() string {
	return sig.s
}

// IsValid 检查签名是否有效
func (sig *Signature) IsValid() bool {
	// 检查基本字段
	if sig.v == "" || sig.r == "" || sig.s == "" {
		return false
	}

	// 检查r和s不能为零
	if sig.r == "0x0" || sig.s == "0x0" {
		return false
	}

	// 验证十六进制格式
	if err := validateHexString(sig.r); err != nil {
		return false
	}

	if err := validateHexString(sig.s); err != nil {
		return false
	}

	return true
}

// Equals 检查两个签名是否相等
func (sig *Signature) Equals(other *Signature) bool {
	if other == nil {
		return false
	}

	return sig.v == other.v &&
		sig.r == other.r &&
		sig.s == other.s
}

// String 返回签名的字符串表示
func (sig *Signature) String() string {
	return fmt.Sprintf("Signature{v: %s, r: %s, s: %s}", sig.v, sig.r, sig.s)
}

// IsZero 检查是否为零值签名
func (sig *Signature) IsZero() bool {
	return sig.v == "" && sig.r == "" && sig.s == ""
}

// validateHexString 验证十六进制字符串格式
func validateHexString(hexStr string) error {
	if hexStr == "" {
		return errors.New("十六进制字符串不能为空")
	}

	// 移除0x前缀
	cleanHex := hexStr
	if strings.HasPrefix(hexStr, "0x") {
		cleanHex = hexStr[2:]
	}

	// 检查长度（r和s应该是32字节，即64个十六进制字符）
	if len(cleanHex) > 64 {
		return errors.New("十六进制字符串长度超出范围")
	}

	// 检查字符是否都是有效的十六进制字符
	for _, char := range cleanHex {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return errors.New("包含无效的十六进制字符")
		}
	}

	return nil
}

// ZeroSignature 返回一个零值签名
func ZeroSignature() *Signature {
	return &Signature{
		v: "",
		r: "",
		s: "",
	}
}

// IsEIP155Signature 检查是否为EIP-155签名（链ID编码）
func (sig *Signature) IsEIP155Signature() bool {
	// EIP-155签名的v值通常大于28
	// 简化检查：如果v不是"27"或"28"，则可能是EIP-155
	return sig.v != "27" && sig.v != "28" && sig.v != "0x1b" && sig.v != "0x1c"
}
