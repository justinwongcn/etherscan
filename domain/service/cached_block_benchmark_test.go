// Package service 提供缓存区块服务的性能基准测试
package service

import (
	"context"
	"testing"
	"time"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/mock"
)

// BenchmarkCachedBlockService_CacheHit 基准测试缓存命中的性能
func BenchmarkCachedBlockService_CacheHit(b *testing.B) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	
	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash:   func() *eth.Hash { h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890"); return h }(),
	}
	
	// 配置模拟仓储，只期望被调用一次（第一次缓存未命中）
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Once()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	if err != nil {
		b.Fatalf("创建缓存服务失败: %v", err)
	}

	ctx := context.Background()

	// 第一次调用，填充缓存
	_, err = service.GetBlockByNumber(ctx, "12345", false)
	if err != nil {
		b.Fatalf("第一次调用失败: %v", err)
	}

	// 重置计时器，开始基准测试
	b.ResetTimer()

	// 基准测试缓存命中的性能
	for i := 0; i < b.N; i++ {
		_, err := service.GetBlockByNumber(ctx, "12345", false)
		if err != nil {
			b.Fatalf("缓存命中调用失败: %v", err)
		}
	}
}

// BenchmarkCachedBlockService_CacheMiss 基准测试缓存未命中的性能
func BenchmarkCachedBlockService_CacheMiss(b *testing.B) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	
	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash:   func() *eth.Hash { h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890"); return h }(),
	}
	
	// 配置模拟仓储，允许多次调用
	mockRepo.On("GetBlockByNumber", mock.Anything, mock.AnythingOfType("string"), false).Return(testBlock, nil)

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	if err != nil {
		b.Fatalf("创建缓存服务失败: %v", err)
	}

	ctx := context.Background()

	// 重置计时器，开始基准测试
	b.ResetTimer()

	// 基准测试缓存未命中的性能（每次使用不同的区块号）
	for i := 0; i < b.N; i++ {
		blockNumber := string(rune(12345 + i)) // 使用不同的区块号确保缓存未命中
		_, err := service.GetBlockByNumber(ctx, blockNumber, false)
		if err != nil {
			b.Fatalf("缓存未命中调用失败: %v", err)
		}
	}
}

// BenchmarkCachedBlockService_HashNormalization 基准测试哈希规范化的性能
func BenchmarkCachedBlockService_HashNormalization(b *testing.B) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	
	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash:   func() *eth.Hash { h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890"); return h }(),
	}
	
	// 规范化后的哈希
	normalizedHash := "0xabc123def456789012345678901234567890123456789012345678901234567890"
	
	// 配置模拟仓储，只期望被调用一次
	mockRepo.On("GetBlockByHash", mock.Anything, normalizedHash, false).Return(testBlock, nil).Once()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	if err != nil {
		b.Fatalf("创建缓存服务失败: %v", err)
	}

	ctx := context.Background()

	// 第一次调用，填充缓存
	_, err = service.GetBlockByHash(ctx, normalizedHash, false)
	if err != nil {
		b.Fatalf("第一次调用失败: %v", err)
	}

	// 不同格式的哈希
	hashFormats := []string{
		"0xabc123def456789012345678901234567890123456789012345678901234567890", // 标准格式
		"0xABC123DEF456789012345678901234567890123456789012345678901234567890", // 大写
		"abc123def456789012345678901234567890123456789012345678901234567890",   // 无前缀小写
		"ABC123DEF456789012345678901234567890123456789012345678901234567890",   // 无前缀大写
	}

	// 重置计时器，开始基准测试
	b.ResetTimer()

	// 基准测试哈希规范化和缓存命中的性能
	for i := 0; i < b.N; i++ {
		hash := hashFormats[i%len(hashFormats)]
		_, err := service.GetBlockByHash(ctx, hash, false)
		if err != nil {
			b.Fatalf("哈希规范化调用失败: %v", err)
		}
	}
}

// BenchmarkCachedBlockService_ConcurrentAccess 基准测试并发访问的性能
func BenchmarkCachedBlockService_ConcurrentAccess(b *testing.B) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	
	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash:   func() *eth.Hash { h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890"); return h }(),
	}
	
	// 配置模拟仓储，允许多次调用
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Maybe()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	if err != nil {
		b.Fatalf("创建缓存服务失败: %v", err)
	}

	ctx := context.Background()

	// 第一次调用，填充缓存
	_, err = service.GetBlockByNumber(ctx, "12345", false)
	if err != nil {
		b.Fatalf("第一次调用失败: %v", err)
	}

	// 重置计时器，开始基准测试
	b.ResetTimer()

	// 基准测试并发访问的性能
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.GetBlockByNumber(ctx, "12345", false)
			if err != nil {
				b.Fatalf("并发访问调用失败: %v", err)
			}
		}
	})
}

// BenchmarkCachedBlockService_CacheExpiration 基准测试缓存过期的性能影响
func BenchmarkCachedBlockService_CacheExpiration(b *testing.B) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	
	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash:   func() *eth.Hash { h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890"); return h }(),
	}
	
	// 配置模拟仓储，允许多次调用
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Maybe()

	// 创建短过期时间的缓存配置（100毫秒过期）
	config := &BlockCacheConfig{
		PrimaryExpiration:   100 * time.Millisecond,
		SecondaryExpiration: 200 * time.Millisecond,
		MaxPrimaryBlocks:    10,
		MaxSecondaryBlocks:  10,
	}

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, config)
	if err != nil {
		b.Fatalf("创建缓存服务失败: %v", err)
	}

	ctx := context.Background()

	// 重置计时器，开始基准测试
	b.ResetTimer()

	// 基准测试缓存过期的性能影响
	for i := 0; i < b.N; i++ {
		_, err := service.GetBlockByNumber(ctx, "12345", false)
		if err != nil {
			b.Fatalf("缓存过期测试调用失败: %v", err)
		}
		
		// 每隔一段时间等待缓存过期
		if i%10 == 0 {
			time.Sleep(150 * time.Millisecond) // 超过过期时间
		}
	}
}
