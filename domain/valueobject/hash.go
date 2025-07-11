// Package valueobject 定义了领域层的值对象
// 值对象是DDD中的重要概念，它们是不可变的、没有身份标识的对象
package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// Hash 表示以太坊中的哈希值对象
// 哈希值是32字节的十六进制字符串，以0x开头
// 这是一个值对象，具有不可变性和值相等性
// 完全使用go-ethlibs包的Hash实现
type Hash struct {
	ethHash string // 使用go-ethlibs验证后的标准格式
}

// NewHash 创建一个新的Hash值对象
// 参数:
//   - hashStr: 哈希字符串，必须是64个十六进制字符，可以带或不带0x前缀
//
// 返回:
//   - Hash: 创建的Hash值对象
//   - error: 如果哈希格式无效，返回错误
func NewHash(hashStr string) (Hash, error) {
	if hashStr == "" {
		return Hash{}, errors.New("哈希值不能为空")
	}

	// 使用go-ethlibs包进行验证和标准化
	return newHashWithEthLibs(hashStr)
}

// newHashWithEthLibs 使用go-ethlibs包创建Hash
func newHashWithEthLibs(hashStr string) (Hash, error) {
	// 确保哈希以0x开头（go-ethlibs要求）
	if !strings.HasPrefix(hashStr, "0x") {
		hashStr = "0x" + hashStr
	}

	// 使用go-ethlibs包进行验证和标准化
	// 这里我们需要通过一个函数调用来避免IDE移除导入
	return createHashFromEthLibs(hashStr)
}

// createHashFromEthLibs 实际调用go-ethlibs包
func createHashFromEthLibs(hashStr string) (Hash, error) {
	// 使用go-ethlibs包进行验证和标准化
	// 这里我们需要手动调用eth.NewHash来保持导入
	return validateHashWithEthLibs(hashStr)
}

// validateHashWithEthLibs 使用go-ethlibs包验证哈希
func validateHashWithEthLibs(hashStr string) (Hash, error) {
	// 使用专门的go-ethlibs验证函数
	validatedHash, err := validateWithEthLibs(hashStr)
	if err != nil {
		return Hash{}, err
	}

	return Hash{ethHash: validatedHash}, nil
}

// NewHashFromBytes 从字节数组创建Hash值对象
// 参数:
//   - bytes: 32字节的字节数组
//
// 返回:
//   - Hash: 创建的Hash值对象
//   - error: 如果字节数组长度不正确，返回错误
func NewHashFromBytes(bytes []byte) (Hash, error) {
	if len(bytes) != 32 {
		return Hash{}, fmt.Errorf("无效的字节数组长度: 期望32字节，实际%d字节", len(bytes))
	}

	hashStr := fmt.Sprintf("0x%x", bytes)
	return Hash{ethHash: hashStr}, nil
}

// String 返回哈希的字符串表示（带0x前缀）
func (h Hash) String() string {
	return h.ethHash
}

// Value 返回哈希的原始值
func (h Hash) Value() string {
	return h.ethHash
}

// Bytes 返回哈希的字节表示
func (h Hash) Bytes() []byte {
	// 移除0x前缀并转换为字节
	cleanHash := strings.TrimPrefix(h.ethHash, "0x")
	bytes := make([]byte, 32)

	for i := 0; i < 32; i++ {
		fmt.Sscanf(cleanHash[i*2:i*2+2], "%02x", &bytes[i])
	}

	return bytes
}

// Equals 检查两个Hash是否相等
func (h Hash) Equals(other Hash) bool {
	return h.ethHash == other.ethHash
}

// IsZero 检查是否为零哈希
func (h Hash) IsZero() bool {
	return h.ethHash == "0x0000000000000000000000000000000000000000000000000000000000000000" || h.ethHash == ""
}
