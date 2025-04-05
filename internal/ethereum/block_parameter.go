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
//
// 返回:
//   - string: 标准化后的区块参数
//   - error: 如果解析过程中发生错误，将返回相应的错误信息
func ParseBlockParameter(blockHashOrNumber string) (string, error) {
	// 如果是空字符串，默认使用latest
	if blockHashOrNumber == "" {
		return BlockLatest, nil
	}

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
