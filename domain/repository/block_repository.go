// Package repository 定义了领域层的仓储接口
// 这些接口遵循依赖倒置原则，由基础设施层提供具体实现
package repository

import (
	"context"

	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockRepository 定义了区块相关数据访问的接口
// 该接口封装了所有与区块数据获取相关的操作，遵循DDD中的仓储模式
// 具体实现由基础设施层提供，应用层通过此接口访问区块数据
type BlockRepository interface {
	// GetLatestBlockNumber 获取当前网络的最新区块高度
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期和传递请求相关信息
	//
	// 返回:
	//   - uint64: 最新区块的高度编号
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetLatestBlockNumber(ctx context.Context) (uint64, error)

	// GetBlockByHash 根据区块哈希获取区块详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - hash: 区块哈希值（32字节的十六进制字符串，以0x开头）
	//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希列表
	//
	// 返回:
	//   - *eth.Block: 包含区块完整信息的结构体指针，如果区块不存在则返回nil
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error)

	// GetBlockByNumber 根据区块号获取区块详细信息
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - number: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希列表
	//
	// 返回:
	//   - *eth.Block: 包含区块完整信息的结构体指针，如果区块不存在则返回nil
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error)

	// GetBlockTransactionCountByHash 根据区块哈希获取区块中的交易数量
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - hash: 区块哈希值（32字节的十六进制字符串，以0x开头）
	//
	// 返回:
	//   - uint64: 区块中包含的交易数量
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error)

	// GetBlockTransactionCountByNumber 根据区块号获取区块中的交易数量
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期
	//   - number: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
	//     支持的特殊值："latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
	//
	// 返回:
	//   - uint64: 区块中包含的交易数量
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error)
}
