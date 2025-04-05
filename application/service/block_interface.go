// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"

	"github.com/justinwongcn/etherscan/domain"
)

// BlockServiceInterface 定义了区块服务的接口规范
// 该接口提供了查询以太坊区块链上区块相关信息的方法集合
type BlockServiceInterface interface {
	// GetLatestBlockHeight 获取以太坊网络的最新区块高度
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	// 返回:
	//   - uint64: 最新区块的高度编号
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetLatestBlockHeight(ctx context.Context) (string, error)

	// GetBlock 根据区块号或区块哈希获取区块的详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - blockHashOrNumber: 区块标识符，可以是区块号（数字字符串）或区块哈希（0x开头的十六进制字符串）
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	// 返回:
	//   - *domain.Block: 包含区块完整信息的领域模型指针
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlock(ctx context.Context, blockHashOrNumber string) (*domain.Block, error)

	// GetBlockTransactionCount 获取指定区块中的交易数量
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - blockHashOrNumber: 区块标识符，可以是区块号（数字字符串）或区块哈希（0x开头的十六进制字符串）
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	// 返回:
	//   - string: 区块中的交易数量
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlockTransactionCount(ctx context.Context, blockHashOrNumber string) (string, error)
}
