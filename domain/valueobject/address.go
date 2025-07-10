// Package valueobject 定义了领域层的值对象
package valueobject

import (
	"errors"
	"fmt"
	"strings"

	"github.com/justinwongcn/go-ethlibs/eth"
)

// Address 表示以太坊地址的值对象
// 以太坊地址是20字节的十六进制字符串，以0x开头
// 这是一个值对象，具有不可变性和值相等性
type Address struct {
	value string
}

// NewAddress 创建一个新的Address值对象
// 参数:
//   - addressStr: 地址字符串，必须是40个十六进制字符，可以带或不带0x前缀
//
// 返回:
//   - Address: 创建的Address值对象
//   - error: 如果地址格式无效，返回错误
func NewAddress(addressStr string) (Address, error) {
	if addressStr == "" {
		return Address{}, errors.New("地址不能为空")
	}

	// 移除0x前缀（如果存在）
	cleanAddress := strings.TrimPrefix(addressStr, "0x")

	// 验证长度（40个十六进制字符 = 20字节）
	if len(cleanAddress) != 40 {
		return Address{}, fmt.Errorf("无效的地址长度: 期望40个字符，实际%d个字符", len(cleanAddress))
	}

	// 验证是否为有效的十六进制字符
	for i, char := range cleanAddress {
		if !isHexChar(char) {
			return Address{}, fmt.Errorf("无效的十六进制字符 '%c' 在位置 %d", char, i)
		}
	}

	// 确保以0x开头，并转换为小写（标准化）
	normalizedAddress := "0x" + strings.ToLower(cleanAddress)

	return Address{value: normalizedAddress}, nil
}

// NewAddressFromBytes 从字节数组创建Address值对象
// 参数:
//   - bytes: 20字节的字节数组
//
// 返回:
//   - Address: 创建的Address值对象
//   - error: 如果字节数组长度不正确，返回错误
func NewAddressFromBytes(bytes []byte) (Address, error) {
	if len(bytes) != 20 {
		return Address{}, fmt.Errorf("无效的字节数组长度: 期望20字节，实际%d字节", len(bytes))
	}

	addressStr := fmt.Sprintf("0x%x", bytes)
	return Address{value: addressStr}, nil
}

// String 返回地址的校验和格式字符串表示（带0x前缀，混合大小写）
func (a Address) String() string {
	// 使用go-ethlibs包的ToChecksumAddress功能
	ethAddr, err := eth.NewAddress(a.value)
	if err != nil {
		// 如果转换失败，返回原始值
		return a.value
	}
	return ethAddr.String()
}

// Value 返回地址的原始值（小写格式）
func (a Address) Value() string {
	return a.value
}

// LowerCase 返回地址的小写格式
func (a Address) LowerCase() string {
	return a.value
}

// Bytes 返回地址的字节表示
func (a Address) Bytes() []byte {
	// 移除0x前缀并转换为字节
	cleanAddress := strings.TrimPrefix(a.value, "0x")
	bytes := make([]byte, 20)

	for i := 0; i < 20; i++ {
		fmt.Sscanf(cleanAddress[i*2:i*2+2], "%02x", &bytes[i])
	}

	return bytes
}

// Equals 检查两个Address是否相等
func (a Address) Equals(other Address) bool {
	return a.value == other.value
}

// IsZero 检查是否为零地址
func (a Address) IsZero() bool {
	return a.value == "0x0000000000000000000000000000000000000000"
}

// ToChecksumAddress 返回EIP-55校验和格式的地址
// 使用go-ethlibs包的ToChecksumAddress功能
func (a Address) ToChecksumAddress() string {
	ethAddr, err := eth.NewAddress(a.value)
	if err != nil {
		// 如果转换失败，返回原始值
		return a.value
	}
	return ethAddr.String()
}
