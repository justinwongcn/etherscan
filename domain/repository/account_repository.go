// Package repository 定义了领域层的仓储接口
// 这些接口遵循依赖倒置原则，由基础设施层提供具体实现
package repository

import (
	"context"
)

// AccountRepository 定义了账户相关数据访问的接口
// 该接口封装了所有与账户数据获取相关的操作，遵循DDD中的仓储模式
// 具体实现由基础设施层提供，应用层通过此接口访问账户数据
type AccountRepository interface {
	// GetBalance 获取指定地址的账户余额
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - address: 以太坊账户地址
	//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//
	// 返回:
	//   - uint64: 账户余额，以wei为单位
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBalance(ctx context.Context, address string, blockParameter string) (uint64, error)

	// GetBalances 批量获取多个地址的账户余额
	// 该方法通过一次请求获取多个账户的ETH余额，提高查询效率
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - addresses: 以太坊账户地址切片
	//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//
	// 返回:
	//   - map[string]uint64: 地址到余额的映射，余额以wei为单位
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBalances(ctx context.Context, addresses []string, blockParameter string) (map[string]uint64, error)
}
