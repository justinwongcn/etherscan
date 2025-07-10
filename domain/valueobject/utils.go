// Package valueobject 定义了领域层的值对象
package valueobject

// isHexChar 检查字符是否为有效的十六进制字符
// 这是一个共享的工具函数，避免在多个文件中重复定义
func isHexChar(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
