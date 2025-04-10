package service

import "context"

// AccountServiceInterface 定义账户相关的服务接口
type AccountServiceInterface interface {
	// GetBalance 获取指定地址的账户余额
	//
	// Parameters:
	//   - ctx: context.Context 用于控制请求的上下文
	//   - address: string 要查询的账户地址
	//   - numberOrTag: string 区块号，可以是以下格式：
	//   - 十六进制字符串（如"0x1"）表示具体区块号
	//   - 十进制整数字符串（如"1234"）表示具体区块号
	//   - "latest" - 最新区块（默认）
	//   - "earliest" - 创世区块
	//   - "pending" - 待处理区块
	//
	// Returns:
	//   - string: 账户余额（单位：wei，十进制字符串格式）
	//   - error: 可能的错误：
	//   - 无效的地址格式
	//   - 无效的区块号格式
	//   - 节点连接错误
	GetBalance(ctx context.Context, address string, numberOrTag string) (string, error)

	// GetBalances 批量获取多个地址的账户余额
	//
	// Parameters:
	//   - ctx: context.Context 用于控制请求的上下文
	//   - addresses: []string 地址列表
	//   - numberOrTag: string 区块号，可以是以下格式：
	//   - 十六进制字符串（如"0x1"）表示具体区块号
	//   - 十进制整数字符串（如"1234"）表示具体区块号
	//   - "latest" - 最新区块（默认）
	//   - "earliest" - 创世区块
	//   - "pending" - 待处理区块
	//
	// Returns:
	//   - map[string]string: 地址到余额的映射（单位：wei，十进制字符串格式）
	//   - error: 可能的错误：
	//   - 无效的地址格式
	//   - 无效的区块号格式
	//   - 节点连接错误
	GetBalances(ctx context.Context, addresses []string, numberOrTag string) (map[string]string, error)
}
