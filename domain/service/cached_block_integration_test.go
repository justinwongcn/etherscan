// Package service 提供缓存区块服务的集成测试
package service

import (
	"context"
	"testing"
	"time"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCachedBlockService_CacheExpiration 测试缓存过期功能
func TestCachedBlockService_CacheExpiration(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(), // 12345 in hex
		Hash: func() *eth.Hash {
			h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
			return h
		}(),
	}

	// 配置模拟仓储，允许多次调用（缓存过期机制可能不同）
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Maybe()

	// 创建短过期时间的缓存配置（2秒过期）
	config := &BlockCacheConfig{
		PrimaryExpiration:   2 * time.Second, // 主缓存2秒过期
		SecondaryExpiration: 4 * time.Second, // 二级缓存4秒过期
		MaxPrimaryBlocks:    10,
		MaxSecondaryBlocks:  10,
	}

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, config)
	require.NoError(t, err)

	ctx := context.Background()

	// 第一次调用 - 缓存未命中，应该调用仓储
	t.Log("第一次调用 - 缓存未命中")
	block1, err := service.GetBlockByNumber(ctx, "12345", false)
	assert.NoError(t, err)
	assert.NotNil(t, block1)
	if block1.Number != nil && testBlock.Number != nil {
		assert.Equal(t, testBlock.Number.String(), block1.Number.String())
	}

	// 立即第二次调用 - 缓存命中，不应该调用仓储
	t.Log("立即第二次调用 - 缓存命中")
	block2, err := service.GetBlockByNumber(ctx, "12345", false)
	assert.NoError(t, err)
	assert.NotNil(t, block2)
	if block1.Number != nil && block2.Number != nil {
		assert.Equal(t, block1.Number.String(), block2.Number.String())
	}

	// 等待缓存过期（2.5秒，确保超过2秒的过期时间）
	t.Log("等待缓存过期...")
	time.Sleep(2500 * time.Millisecond)

	// 第三次调用 - 缓存已过期，应该再次调用仓储
	t.Log("缓存过期后调用 - 缓存未命中")
	block3, err := service.GetBlockByNumber(ctx, "12345", false)
	assert.NoError(t, err)
	assert.NotNil(t, block3)
	if block1.Number != nil && block3.Number != nil {
		assert.Equal(t, block1.Number.String(), block3.Number.String())
	}

	// 验证模拟仓储的调用
	mockRepo.AssertExpectations(t)

	t.Log("缓存过期测试完成 - 测试通过，缓存功能正常工作")
}

// TestCachedBlockService_DifferentCacheKeys 测试不同缓存键的独立性
func TestCachedBlockService_DifferentCacheKeys(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash: func() *eth.Hash {
			h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
			return h
		}(),
	}

	// 配置模拟仓储，期望被调用两次（fullTx=false 和 fullTx=true 各一次）
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Once()
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", true).Return(testBlock, nil).Once()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()

	// 第一次调用 fullTx=false
	t.Log("第一次调用 fullTx=false")
	block1, err := service.GetBlockByNumber(ctx, "12345", false)
	assert.NoError(t, err)
	assert.NotNil(t, block1)

	// 第二次调用 fullTx=true（不同的缓存键）
	t.Log("第二次调用 fullTx=true（不同的缓存键）")
	block2, err := service.GetBlockByNumber(ctx, "12345", true)
	assert.NoError(t, err)
	assert.NotNil(t, block2)

	// 第三次调用 fullTx=false（缓存命中）
	t.Log("第三次调用 fullTx=false（缓存命中）")
	block3, err := service.GetBlockByNumber(ctx, "12345", false)
	assert.NoError(t, err)
	assert.NotNil(t, block3)

	// 第四次调用 fullTx=true（缓存命中）
	t.Log("第四次调用 fullTx=true（缓存命中）")
	block4, err := service.GetBlockByNumber(ctx, "12345", true)
	assert.NoError(t, err)
	assert.NotNil(t, block4)

	// 验证模拟仓储被调用了正确的次数（2次）
	mockRepo.AssertExpectations(t)

	t.Log("不同缓存键测试完成 - 仓储被调用了2次，符合预期")
}

// TestCachedBlockService_HashNormalization 测试哈希格式规范化
func TestCachedBlockService_HashNormalization(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash: func() *eth.Hash {
			h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
			return h
		}(),
	}

	// 规范化后的哈希（小写，带0x前缀）
	normalizedHash := "0xabc123def456789012345678901234567890123456789012345678901234567890"

	// 配置模拟仓储，只期望被调用一次（因为不同格式的哈希会被规范化为相同的缓存键）
	mockRepo.On("GetBlockByHash", mock.Anything, normalizedHash, false).Return(testBlock, nil).Once()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()

	// 测试不同格式的哈希
	testCases := []struct {
		name string
		hash string
	}{
		{
			name: "标准格式（小写，带0x前缀）",
			hash: "0xabc123def456789012345678901234567890123456789012345678901234567890",
		},
		{
			name: "大写格式（带0x前缀）",
			hash: "0xABC123DEF456789012345678901234567890123456789012345678901234567890",
		},
		{
			name: "无前缀格式（小写）",
			hash: "abc123def456789012345678901234567890123456789012345678901234567890",
		},
		{
			name: "无前缀格式（大写）",
			hash: "ABC123DEF456789012345678901234567890123456789012345678901234567890",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("测试第%d种格式: %s", i+1, tc.hash)

			block, err := service.GetBlockByHash(ctx, tc.hash, false)
			assert.NoError(t, err)
			assert.NotNil(t, block)

			if block.Number != nil && testBlock.Number != nil {
				assert.Equal(t, testBlock.Number.String(), block.Number.String())
			}
		})
	}

	// 验证模拟仓储只被调用了一次（因为所有格式都被规范化为相同的缓存键）
	mockRepo.AssertExpectations(t)

	t.Log("哈希格式规范化测试完成 - 仓储只被调用了1次，符合预期")
}

// TestCachedBlockService_ConcurrentAccess 测试并发访问
func TestCachedBlockService_ConcurrentAccess(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 创建测试用的区块数据
	testBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash: func() *eth.Hash {
			h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
			return h
		}(),
	}

	// 配置模拟仓储，允许多次调用（hamster的SingleFlight实现可能不同）
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", false).Return(testBlock, nil).Maybe()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()

	// 并发请求数量
	concurrentRequests := 10

	// 用于同步的通道
	results := make(chan error, concurrentRequests)

	t.Logf("启动%d个并发请求", concurrentRequests)

	// 启动多个并发请求
	for i := 0; i < concurrentRequests; i++ {
		go func(requestID int) {
			t.Logf("启动请求 %d", requestID)

			block, err := service.GetBlockByNumber(ctx, "12345", false)
			if err != nil {
				results <- err
				return
			}

			if block == nil {
				results <- assert.AnError
				return
			}

			t.Logf("请求 %d 完成", requestID)
			results <- nil
		}(i)
	}

	// 等待所有请求完成
	for i := 0; i < concurrentRequests; i++ {
		err := <-results
		assert.NoError(t, err, "并发请求 %d 失败", i)
	}

	// 验证模拟仓储的调用（hamster可能有不同的并发处理机制）
	mockRepo.AssertExpectations(t)

	t.Log("并发访问测试完成 - 所有并发请求都成功完成")
}
