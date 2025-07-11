// Package ethereum 定义了以太坊客户端相关的错误常量
package ethereum

// 错误消息常量，避免重复定义
const (
	// ErrInvalidBlockHashFormat 无效的区块哈希格式错误消息
	ErrInvalidBlockHashFormat = "invalid block hash format: must be hex string starting with 0x"
	
	// ErrInvalidBlockParameter 无效的区块参数错误消息
	ErrInvalidBlockParameter = "invalid block parameter: %v"
	
	// ErrInvalidBlockNumberOrTag 无效的区块号或标签错误消息
	ErrInvalidBlockNumberOrTag = "invalid block number or tag: %s"
	
	// ErrFailedToGetConnection 获取连接失败错误消息
	ErrFailedToGetConnection = "failed to get connection: %v"
)
