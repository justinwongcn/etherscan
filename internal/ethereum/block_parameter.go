// Package ethereum 提供以太坊区块链交互的基础功能实现
package ethereum

import (
	"fmt"
	"math/big"
	"strings"
)

// 使用types.go中定义的区块标识符常量

// ParseBlockParameter 解析并标准化区块参数格式
// 将用户输入的区块标识符转换为以太坊API支持的格式
// 参数:
//   - blockHashOrNumber: 区块标识符，可以是区块号、区块哈希或特殊标识符
//   - 空字符串: 将被解析为"latest"
//   - 特殊标识符: "latest"(最新区块)、"earliest"(创世区块)、"pending"(待处理区块)
//   - 区块哈希: 以"0x"开头的64位十六进制字符串
//   - 区块号: 十进制数字字符串
//
// 返回:
//   - string: 标准化后的区块参数
//   - 特殊标识符: 返回小写的标识符字符串
//   - 区块哈希: 保持原样返回
//   - 区块号: 转换为"0x"前缀的十六进制格式
//   - error: 如果解析过程中发生错误，将返回相应的错误信息
func ParseBlockParameter(blockHashOrNumber string) (string, error) {
	// 处理特殊标识符
	switch strings.ToLower(blockHashOrNumber) {
	case BlockLatest, BlockEarliest, BlockPending:
		return strings.ToLower(blockHashOrNumber), nil
	}

	// 如果是区块哈希（以0x开头的十六进制字符串），直接返回
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" {
		return blockHashOrNumber, nil
	}

	// 尝试将输入解析为区块号
	blockNum := new(big.Int)
	if _, ok := blockNum.SetString(blockHashOrNumber, 10); !ok {
		return "", fmt.Errorf("invalid block parameter: %s", blockHashOrNumber)
	}

	// 将区块号转换为十六进制格式
	return fmt.Sprintf("0x%x", blockNum), nil
}
