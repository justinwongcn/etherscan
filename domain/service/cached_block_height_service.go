// Package service 提供基于缓存的区块高度查询服务
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/justinwongcn/hamster"
	"github.com/justinwongcn/hamster/cache"
)

// CachedBlockHeightService 是一个装饰器，为BlockHeightService添加缓存功能
// 该结构体使用hamster库的ReadThroughService实现，内置了singleflight机制来防止缓存击穿
// 缓存过期时间为5秒，只有在收到请求时才会检查缓存是否过期
type CachedBlockHeightService struct {
	// service 是被装饰的原始区块高度服务
	service BlockHeightService
	// cache 是hamster库的读透缓存服务，内置singleflight机制
	cache *cache.ReadThroughService
}

// NewCachedBlockHeightService 创建一个带有缓存功能的区块高度服务装饰器
// 参数:
//   - service: 被装饰的原始BlockHeightService实例
//   - expiration: 缓存过期时间，建议设置为5秒
//
// 返回:
//   - BlockHeightService: 带有缓存功能的服务实例
//   - error: 如果创建缓存服务失败，将返回相应的错误信息
func NewCachedBlockHeightService(service BlockHeightService, expiration time.Duration) (BlockHeightService, error) {
	// 创建ReadThroughService实例，内置singleflight机制防止缓存击穿
	readThroughCache, err := hamster.NewReadThroughCache(
		cache.WithDefaultExpiration(expiration), // 设置默认过期时间
		// ReadThroughService内部已经使用了singleflight.Group来防止缓存击穿
		// 并且只有在GetWithLoader被调用时才会检查过期并尝试加载
	)
	if err != nil {
		return nil, fmt.Errorf("创建读透缓存失败: %w", err)
	}

	cachedService := &CachedBlockHeightService{
		service: service,
		cache:   readThroughCache,
	}

	return cachedService, nil
}

// GetLatestBlockHeight 实现了BlockHeightService接口中的同名方法
// 从缓存中获取区块高度，如果缓存未命中或过期则从原始服务加载
// 使用singleflight机制确保并发请求时只有一个实际的数据加载操作
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - uint64: 最新区块的高度编号
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *CachedBlockHeightService) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	const cacheKey = "latest_block_height"

	// 定义加载函数，当缓存未命中或过期时调用
	loader := func(ctx context.Context, key string) (any, error) {
		// 从原始服务获取最新区块高度
		height, err := s.service.GetLatestBlockHeight(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取最新区块高度失败: %w", err)
		}
		return height, nil
	}

	// 使用GetWithLoader从缓存获取数据，如果缓存未命中或过期，loader会被自动调用
	// hamster的ReadThroughService内部使用singleflight机制，确保并发请求时只有一个loader执行
	val, err := s.cache.GetWithLoader(ctx, cacheKey, loader, 0) // 使用默认过期时间
	if err != nil {
		return 0, fmt.Errorf("从缓存获取区块高度失败: %w", err)
	}

	// 类型断言，将interface{}转换为uint64
	height, ok := val.(uint64)
	if !ok {
		return 0, fmt.Errorf("缓存中的数据类型错误，期望uint64，实际为%T", val)
	}

	return height, nil
}
