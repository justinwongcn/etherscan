// Package repository 提供领域仓储接口的具体实现
// 该包属于基础设施层，负责与外部数据源（以太坊节点）进行交互
package repository

import (
	"context"

	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockRepositoryImpl 实现了domain.repository.BlockRepository接口
// 该结构体封装了与以太坊节点交互的客户端，提供区块数据访问的具体实现
type BlockRepositoryImpl struct {
	// client 是与以太坊节点通信的客户端实例
	// 通过该客户端执行实际的RPC调用来获取区块数据
	client *ethereum.Client
}

// NewBlockRepository 创建并初始化一个新的BlockRepositoryImpl实例
// 该函数是BlockRepositoryImpl的工厂方法，确保正确的依赖注入
// 参数:
//   - client: 已初始化的以太坊客户端实例，用于与节点通信
//
// 返回:
//   - repository.BlockRepository: 实现了BlockRepository接口的实例
func NewBlockRepository(client *ethereum.Client) repository.BlockRepository {
	return &BlockRepositoryImpl{
		client: client,
	}
}

// GetLatestBlockNumber 实现了BlockRepository接口中的同名方法
// 通过调用以太坊客户端获取当前网络的最新区块高度
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - uint64: 最新区块的高度编号
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *BlockRepositoryImpl) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return r.client.GetLatestBlockNumber(ctx)
}

// GetBlockByHash 实现了BlockRepository接口中的同名方法
// 根据区块哈希获取区块详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - hash: 区块哈希值（32字节的十六进制字符串，以0x开头）
//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希列表
//
// 返回:
//   - *eth.Block: 包含区块完整信息的结构体指针，如果区块不存在则返回nil
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *BlockRepositoryImpl) GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error) {
	return r.client.GetBlockByHash(ctx, hash, fullTx)
}

// GetBlockByNumber 实现了BlockRepository接口中的同名方法
// 根据区块号获取区块详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - number: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//   - fullTx: 如果为true则返回完整的交易对象，否则仅返回交易哈希列表
//
// 返回:
//   - *eth.Block: 包含区块完整信息的结构体指针，如果区块不存在则返回nil
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *BlockRepositoryImpl) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error) {
	return r.client.GetBlockByNumber(ctx, number, fullTx)
}

// GetBlockTransactionCountByHash 实现了BlockRepository接口中的同名方法
// 根据区块哈希获取区块中的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - hash: 区块哈希值（32字节的十六进制字符串，以0x开头）
//
// 返回:
//   - uint64: 区块中包含的交易数量
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *BlockRepositoryImpl) GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	return r.client.GetBlockTransactionCountByHash(ctx, hash)
}

// GetBlockTransactionCountByNumber 实现了BlockRepository接口中的同名方法
// 根据区块号获取区块中的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - number: 区块号参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//
// 返回:
//   - uint64: 区块中包含的交易数量
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *BlockRepositoryImpl) GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error) {
	return r.client.GetBlockTransactionCountByNumber(ctx, number)
}
