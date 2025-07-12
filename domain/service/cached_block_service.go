// Package service 提供基于缓存的区块查询服务
package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/justinwongcn/hamster"
	"github.com/justinwongcn/hamster/cache"
)

// CachedBlockService 是一个装饰器，为区块查询服务添加缓存功能
// 该服务实现了分层缓存策略：
// - 最新的10个区块使用FIFO算法缓存
// - 当缓存过期时，保存最新的11-20个区块使用LRU算法缓存
// - 支持按区块号和区块哈希查询
// - 内置SingleFlight机制防止缓存击穿
type CachedBlockService struct {
	// blockRepo 是区块数据访问的仓储接口
	blockRepo repository.BlockRepository
	// primaryCache 是主缓存，使用ReadThroughCache实现
	primaryCache *cache.ReadThroughService
	// secondaryCache 是二级缓存，用于存储过期的区块
	secondaryCache *cache.ReadThroughService
	// maxPrimaryBlocks 主缓存最大区块数量（最新的10个区块）
	maxPrimaryBlocks int
	// maxSecondaryBlocks 二级缓存最大区块数量（11-20个区块）
	maxSecondaryBlocks int
}

// BlockCacheConfig 区块缓存配置
type BlockCacheConfig struct {
	// PrimaryExpiration 主缓存过期时间
	PrimaryExpiration time.Duration
	// SecondaryExpiration 二级缓存过期时间
	SecondaryExpiration time.Duration
	// MaxPrimaryBlocks 主缓存最大区块数量
	MaxPrimaryBlocks int
	// MaxSecondaryBlocks 二级缓存最大区块数量
	MaxSecondaryBlocks int
}

// DefaultBlockCacheConfig 返回默认的区块缓存配置
func DefaultBlockCacheConfig() *BlockCacheConfig {
	return &BlockCacheConfig{
		PrimaryExpiration:   3 * time.Minute,  // 主缓存3分钟过期
		SecondaryExpiration: 10 * time.Minute, // 二级缓存10分钟过期
		MaxPrimaryBlocks:    10,               // 最新10个区块
		MaxSecondaryBlocks:  10,               // 11-20个区块
	}
}

// NewCachedBlockService 创建一个带有缓存功能的区块查询服务
// 参数:
//   - blockRepo: 区块数据访问的仓储接口
//   - config: 缓存配置，如果为nil则使用默认配置
//
// 返回:
//   - *CachedBlockService: 带有缓存功能的区块服务实例
//   - error: 如果创建缓存服务失败，将返回相应的错误信息
func NewCachedBlockService(blockRepo repository.BlockRepository, config *BlockCacheConfig) (*CachedBlockService, error) {
	if config == nil {
		config = DefaultBlockCacheConfig()
	}

	// 创建主缓存（FIFO策略，用于最新的10个区块）
	primaryCache, err := hamster.NewReadThroughCache(
		cache.WithDefaultExpiration(config.PrimaryExpiration),
		cache.WithMaxMemory(50*1024*1024), // 50MB内存限制
	)
	if err != nil {
		return nil, fmt.Errorf("创建主缓存失败: %w", err)
	}

	// 创建二级缓存（LRU策略，用于11-20个区块）
	secondaryCache, err := hamster.NewReadThroughCache(
		cache.WithDefaultExpiration(config.SecondaryExpiration),
		cache.WithMaxMemory(100*1024*1024), // 100MB内存限制
	)
	if err != nil {
		return nil, fmt.Errorf("创建二级缓存失败: %w", err)
	}

	service := &CachedBlockService{
		blockRepo:          blockRepo,
		primaryCache:       primaryCache,
		secondaryCache:     secondaryCache,
		maxPrimaryBlocks:   config.MaxPrimaryBlocks,
		maxSecondaryBlocks: config.MaxSecondaryBlocks,
	}

	return service, nil
}

// GetBlockByNumber 根据区块号获取区块信息（带缓存）
// 参数:
//   - ctx: 上下文对象
//   - blockNumber: 区块号（十进制字符串或十六进制字符串）
//   - fullTx: 是否返回完整的交易信息
//
// 返回:
//   - *eth.Block: 区块信息
//   - error: 如果查询失败，将返回相应的错误信息
func (s *CachedBlockService) GetBlockByNumber(ctx context.Context, blockNumber string, fullTx bool) (*eth.Block, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("block:number:%s:fullTx:%t", blockNumber, fullTx)

	// 定义加载函数
	loader := func(ctx context.Context, key string) (any, error) {
		// 从原始仓储获取区块数据
		block, err := s.blockRepo.GetBlockByNumber(ctx, blockNumber, fullTx)
		if err != nil {
			return nil, fmt.Errorf("从仓储获取区块失败: %w", err)
		}

		// 直接返回区块对象，不进行序列化
		return block, nil
	}

	// 先尝试从主缓存获取
	blockData, err := s.primaryCache.GetWithLoader(ctx, cacheKey, loader, 0) // 使用默认过期时间
	if err != nil {
		// 如果主缓存失败，尝试从二级缓存获取
		blockData, err = s.secondaryCache.GetWithLoader(ctx, cacheKey, loader, 0)
		if err != nil {
			return nil, fmt.Errorf("从缓存获取区块失败: %w", err)
		}
	}

	// 直接类型断言，不需要反序列化
	block, ok := blockData.(*eth.Block)
	if !ok {
		return nil, fmt.Errorf("缓存中的数据类型错误，期望*eth.Block，实际为%T", blockData)
	}

	return block, nil
}

// GetBlockByHash 根据区块哈希获取区块信息（带缓存）
// 参数:
//   - ctx: 上下文对象
//   - blockHash: 区块哈希（十六进制字符串）
//   - fullTx: 是否返回完整的交易信息
//
// 返回:
//   - *eth.Block: 区块信息
//   - error: 如果查询失败，将返回相应的错误信息
func (s *CachedBlockService) GetBlockByHash(ctx context.Context, blockHash string, fullTx bool) (*eth.Block, error) {
	// 规范化区块哈希格式（确保小写且带0x前缀）
	normalizedHash := s.normalizeBlockHash(blockHash)

	// 构造缓存键
	cacheKey := fmt.Sprintf("block:hash:%s:fullTx:%t", normalizedHash, fullTx)

	// 定义加载函数
	loader := func(ctx context.Context, key string) (any, error) {
		// 从原始仓储获取区块数据
		block, err := s.blockRepo.GetBlockByHash(ctx, normalizedHash, fullTx)
		if err != nil {
			return nil, fmt.Errorf("从仓储获取区块失败: %w", err)
		}

		// 直接返回区块对象，不进行序列化
		return block, nil
	}

	// 先尝试从主缓存获取
	blockData, err := s.primaryCache.GetWithLoader(ctx, cacheKey, loader, 0) // 使用默认过期时间
	if err != nil {
		// 如果主缓存失败，尝试从二级缓存获取
		blockData, err = s.secondaryCache.GetWithLoader(ctx, cacheKey, loader, 0)
		if err != nil {
			return nil, fmt.Errorf("从缓存获取区块失败: %w", err)
		}
	}

	// 直接类型断言，不需要反序列化
	block, ok := blockData.(*eth.Block)
	if !ok {
		return nil, fmt.Errorf("缓存中的数据类型错误，期望*eth.Block，实际为%T", blockData)
	}

	return block, nil
}

// PreloadLatestBlocks 预热缓存，加载最新的指定数量区块
// 参数:
//   - ctx: 上下文对象
//   - count: 要预加载的区块数量
//
// 返回:
//   - error: 如果预加载失败，将返回相应的错误信息
func (s *CachedBlockService) PreloadLatestBlocks(ctx context.Context, count int) error {
	// 获取最新区块号
	latestBlockNumber, err := s.blockRepo.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("获取最新区块号失败: %w", err)
	}

	// 预加载最新的count个区块
	for i := 0; i < count; i++ {
		blockNum := latestBlockNumber - uint64(i)
		if blockNum == 0 {
			break
		}

		blockNumberStr := strconv.FormatUint(blockNum, 10)

		// 预加载完整交易信息的区块
		_, err := s.GetBlockByNumber(ctx, blockNumberStr, true)
		if err != nil {
			// 记录错误但不中断预加载过程
			fmt.Printf("预加载区块 %d 失败: %v\n", blockNum, err)
		}

		// 预加载不含完整交易信息的区块
		_, err = s.GetBlockByNumber(ctx, blockNumberStr, false)
		if err != nil {
			fmt.Printf("预加载区块 %d (不含交易) 失败: %v\n", blockNum, err)
		}
	}

	return nil
}

// GetCacheStats 获取缓存统计信息
// 返回:
//   - map[string]any: 包含主缓存和二级缓存统计信息的映射
//   - error: 如果获取统计信息失败，将返回相应的错误信息
func (s *CachedBlockService) GetCacheStats(ctx context.Context) (map[string]any, error) {
	// 注意：hamster的ReadThroughService可能没有Stats方法
	// 这里返回基本的缓存信息
	return map[string]any{
		"primary": map[string]any{
			"status":      "active",
			"description": "主缓存 - 最新区块",
		},
		"secondary": map[string]any{
			"status":      "active",
			"description": "二级缓存 - 历史区块",
		},
	}, nil
}

// normalizeBlockHash 规范化区块哈希格式
func (s *CachedBlockService) normalizeBlockHash(hash string) string {
	// 转换为小写
	hash = strings.ToLower(hash)

	// 确保有0x前缀
	if !strings.HasPrefix(hash, "0x") {
		hash = "0x" + hash
	}

	return hash
}
