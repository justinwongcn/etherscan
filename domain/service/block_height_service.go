// Package service 定义了领域层的服务接口
// 这些接口遵循DDD架构原则，封装了核心业务逻辑
package service

import (
	"context"
)

// BlockHeightService 定义了区块高度查询的服务接口
// 该接口封装了获取区块高度的核心业务逻辑，遵循DDD中的领域服务模式
// 具体实现可以是直接查询、缓存查询或其他策略
type BlockHeightService interface {
	// GetLatestBlockHeight 获取最新区块高度
	// 参数:
	//   - ctx: 上下文对象，用于控制请求的生命周期和传递请求相关信息
	//
	// 返回:
	//   - uint64: 最新区块的高度编号
	//   - error: 如果查询过程中发生错误，将返回相应的错误信息
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
}
