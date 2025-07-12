// Package service 提供区块高度查询的基础实现
package service

import (
	"context"

	"github.com/justinwongcn/etherscan/domain/repository"
)

// blockHeightServiceImpl 实现了BlockHeightService接口
// 该结构体提供了直接从区块仓储获取最新区块高度的基础实现
type blockHeightServiceImpl struct {
	// blockRepo 是区块数据访问的仓储接口
	// 通过依赖倒置原则，服务层依赖抽象而非具体实现
	blockRepo repository.BlockRepository
}

// NewBlockHeightService 创建并初始化一个新的BlockHeightService实例
// 参数:
//   - blockRepo: 实现了BlockRepository接口的仓储实例，用于区块数据访问
//
// 返回:
//   - BlockHeightService: 初始化完成的服务实例
func NewBlockHeightService(blockRepo repository.BlockRepository) BlockHeightService {
	return &blockHeightServiceImpl{
		blockRepo: blockRepo,
	}
}

// GetLatestBlockHeight 实现了BlockHeightService接口中的同名方法
// 通过调用区块仓储获取当前网络的最新区块高度
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - uint64: 最新区块的高度编号
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *blockHeightServiceImpl) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	// 直接从仓储获取最新区块号
	return s.blockRepo.GetLatestBlockNumber(ctx)
}
