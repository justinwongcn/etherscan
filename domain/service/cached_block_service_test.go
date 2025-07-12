// Package service 提供缓存区块服务的测试
package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBlockRepository 是BlockRepository的模拟实现，用于测试
type MockBlockRepository struct {
	mock.Mock
}

// GetLatestBlockNumber 模拟获取最新区块号
func (m *MockBlockRepository) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetBlockByNumber 模拟根据区块号获取区块
func (m *MockBlockRepository) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error) {
	args := m.Called(ctx, number, fullTx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Block), args.Error(1)
}

// GetBlockByHash 模拟根据区块哈希获取区块
func (m *MockBlockRepository) GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error) {
	args := m.Called(ctx, hash, fullTx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Block), args.Error(1)
}

// GetBlockTransactionCountByNumber 模拟根据区块号获取交易数量
func (m *MockBlockRepository) GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(uint64), args.Error(1)
}

// GetBlockTransactionCountByHash 模拟根据区块哈希获取交易数量
func (m *MockBlockRepository) GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(uint64), args.Error(1)
}

// TestNewCachedBlockService 测试缓存区块服务的创建
func TestNewCachedBlockService(t *testing.T) {
	tests := []struct {
		name    string
		config  *BlockCacheConfig
		wantErr bool
	}{
		{
			name:    "使用默认配置创建服务",
			config:  nil,
			wantErr: false,
		},
		{
			name: "使用自定义配置创建服务",
			config: &BlockCacheConfig{
				PrimaryExpiration:   2 * time.Minute,
				SecondaryExpiration: 10 * time.Minute,
				MaxPrimaryBlocks:    5,
				MaxSecondaryBlocks:  5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓储
			mockRepo := &MockBlockRepository{}

			// 创建缓存服务
			service, err := NewCachedBlockService(mockRepo, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.NotNil(t, service.primaryCache)
				assert.NotNil(t, service.secondaryCache)
			}
		})
	}
}

// TestCachedBlockService_GetBlockByNumber 测试根据区块号获取区块
func TestCachedBlockService_GetBlockByNumber(t *testing.T) {
	tests := []struct {
		name        string
		blockNumber string
		fullTx      bool
		mockBlock   *eth.Block
		mockError   error
		expectError string
	}{
		{
			name:        "成功获取区块",
			blockNumber: "12345",
			fullTx:      true,
			mockBlock: &eth.Block{
				Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(), // 12345 in hex
				Hash: func() *eth.Hash {
					h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
					return h
				}(),
			},
			mockError:   nil,
			expectError: "",
		},
		{
			name:        "仓储返回错误",
			blockNumber: "12345",
			fullTx:      false,
			mockBlock:   nil,
			mockError:   errors.New("区块不存在"),
			expectError: "从缓存获取区块失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓储
			mockRepo := &MockBlockRepository{}
			mockRepo.On("GetBlockByNumber", mock.Anything, tt.blockNumber, tt.fullTx).Return(tt.mockBlock, tt.mockError)

			// 创建缓存服务
			service, err := NewCachedBlockService(mockRepo, nil)
			require.NoError(t, err)

			// 调用被测试的方法
			ctx := context.Background()
			block, err := service.GetBlockByNumber(ctx, tt.blockNumber, tt.fullTx)

			// 验证结果
			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				assert.Nil(t, block)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, block)
				if block.Number != nil && tt.mockBlock.Number != nil {
					assert.Equal(t, tt.mockBlock.Number.String(), block.Number.String())
				}
				if block.Hash != nil && tt.mockBlock.Hash != nil {
					assert.Equal(t, tt.mockBlock.Hash.String(), block.Hash.String())
				}
			}

			// 验证模拟仓储被调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestCachedBlockService_GetBlockByHash 测试根据区块哈希获取区块
func TestCachedBlockService_GetBlockByHash(t *testing.T) {
	tests := []struct {
		name           string
		blockHash      string
		normalizedHash string
		fullTx         bool
		mockBlock      *eth.Block
		mockError      error
		expectError    string
	}{
		{
			name:           "成功获取区块",
			blockHash:      "0xABC123",
			normalizedHash: "0xabc123",
			fullTx:         true,
			mockBlock: &eth.Block{
				Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
				Hash: func() *eth.Hash {
					h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
					return h
				}(),
			},
			mockError:   nil,
			expectError: "",
		},
		{
			name:           "哈希格式规范化",
			blockHash:      "ABC123", // 没有0x前缀，大写
			normalizedHash: "0xabc123",
			fullTx:         false,
			mockBlock: &eth.Block{
				Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
				Hash: func() *eth.Hash {
					h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
					return h
				}(),
			},
			mockError:   nil,
			expectError: "",
		},
		{
			name:           "仓储返回错误",
			blockHash:      "0xdef456",
			normalizedHash: "0xdef456",
			fullTx:         true,
			mockBlock:      nil,
			mockError:      errors.New("区块不存在"),
			expectError:    "从缓存获取区块失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓储
			mockRepo := &MockBlockRepository{}
			mockRepo.On("GetBlockByHash", mock.Anything, tt.normalizedHash, tt.fullTx).Return(tt.mockBlock, tt.mockError)

			// 创建缓存服务
			service, err := NewCachedBlockService(mockRepo, nil)
			require.NoError(t, err)

			// 调用被测试的方法
			ctx := context.Background()
			block, err := service.GetBlockByHash(ctx, tt.blockHash, tt.fullTx)

			// 验证结果
			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				assert.Nil(t, block)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, block)
				if block.Number != nil && tt.mockBlock.Number != nil {
					assert.Equal(t, tt.mockBlock.Number.String(), block.Number.String())
				}
				if block.Hash != nil && tt.mockBlock.Hash != nil {
					assert.Equal(t, tt.mockBlock.Hash.String(), block.Hash.String())
				}
			}

			// 验证模拟仓储被调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestCachedBlockService_CacheHit 测试缓存命中功能
func TestCachedBlockService_CacheHit(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}
	mockBlock := &eth.Block{
		Number: func() *eth.Quantity { q, _ := eth.NewQuantity("0x3039"); return q }(),
		Hash: func() *eth.Hash {
			h, _ := eth.NewHash("0xabc123def456789012345678901234567890123456789012345678901234567890")
			return h
		}(),
	}
	// 只期望被调用一次
	mockRepo.On("GetBlockByNumber", mock.Anything, "12345", true).Return(mockBlock, nil).Once()

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()

	// 第一次调用 - 缓存未命中，应该调用仓储
	block1, err := service.GetBlockByNumber(ctx, "12345", true)
	assert.NoError(t, err)
	assert.NotNil(t, block1)

	// 第二次调用 - 缓存命中，不应该调用仓储
	block2, err := service.GetBlockByNumber(ctx, "12345", true)
	assert.NoError(t, err)
	assert.NotNil(t, block2)
	if block1.Number != nil && block2.Number != nil {
		assert.Equal(t, block1.Number.String(), block2.Number.String())
	}
	if block1.Hash != nil && block2.Hash != nil {
		assert.Equal(t, block1.Hash.String(), block2.Hash.String())
	}

	// 验证模拟仓储只被调用了一次
	mockRepo.AssertExpectations(t)
}

// TestCachedBlockService_PreloadLatestBlocks 测试预加载功能
func TestCachedBlockService_PreloadLatestBlocks(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 模拟最新区块号
	mockRepo.On("GetLatestBlockNumber", mock.Anything).Return(uint64(12345), nil)

	// 模拟预加载的区块数据
	for i := 0; i < 3; i++ {
		blockNum := 12345 - i
		blockNumStr := fmt.Sprintf("%d", blockNum)
		mockBlock := &eth.Block{
			Number: func() *eth.Quantity { q, _ := eth.NewQuantity(fmt.Sprintf("0x%x", blockNum)); return q }(),
			Hash:   func() *eth.Hash { h, _ := eth.NewHash(fmt.Sprintf("0x%064x", blockNum)); return h }(), // 确保哈希是64位
		}

		// 每个区块会被调用两次：一次fullTx=true，一次fullTx=false
		mockRepo.On("GetBlockByNumber", mock.Anything, blockNumStr, true).Return(mockBlock, nil).Once()
		mockRepo.On("GetBlockByNumber", mock.Anything, blockNumStr, false).Return(mockBlock, nil).Once()
	}

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	// 调用预加载方法
	ctx := context.Background()
	err = service.PreloadLatestBlocks(ctx, 3)
	assert.NoError(t, err)

	// 验证模拟仓储被正确调用
	mockRepo.AssertExpectations(t)
}

// TestCachedBlockService_GetCacheStats 测试获取缓存统计信息
func TestCachedBlockService_GetCacheStats(t *testing.T) {
	// 创建模拟仓储
	mockRepo := &MockBlockRepository{}

	// 创建缓存服务
	service, err := NewCachedBlockService(mockRepo, nil)
	require.NoError(t, err)

	// 获取缓存统计信息
	ctx := context.Background()
	stats, err := service.GetCacheStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Contains(t, stats, "primary")
	assert.Contains(t, stats, "secondary")

	// 验证统计信息结构
	primaryStats := stats["primary"].(map[string]any)
	assert.Contains(t, primaryStats, "status")
	assert.Contains(t, primaryStats, "description")

	secondaryStats := stats["secondary"].(map[string]any)
	assert.Contains(t, secondaryStats, "status")
	assert.Contains(t, secondaryStats, "description")
}
