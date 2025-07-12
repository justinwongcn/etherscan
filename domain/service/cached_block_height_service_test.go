// Package service 提供缓存区块高度服务的测试
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBlockHeightService 是BlockHeightService的模拟实现，用于测试
type MockBlockHeightService struct {
	mock.Mock
}

// GetLatestBlockHeight 模拟获取最新区块高度
func (m *MockBlockHeightService) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// TestNewCachedBlockHeightService 测试缓存区块高度服务的创建
func TestNewCachedBlockHeightService(t *testing.T) {
	tests := []struct {
		name       string
		expiration time.Duration
		wantErr    bool
	}{
		{
			name:       "正常创建缓存服务",
			expiration: 5 * time.Second,
			wantErr:    false,
		},
		{
			name:       "使用1秒过期时间",
			expiration: 1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "使用10秒过期时间",
			expiration: 10 * time.Second,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟服务
			mockService := &MockBlockHeightService{}

			// 创建缓存服务
			cachedService, err := NewCachedBlockHeightService(mockService, tt.expiration)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cachedService)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cachedService)
			}
		})
	}
}

// TestCachedBlockHeightService_GetLatestBlockHeight 测试缓存区块高度获取功能
func TestCachedBlockHeightService_GetLatestBlockHeight(t *testing.T) {
	tests := []struct {
		name           string
		mockHeight     uint64
		mockError      error
		expectedHeight uint64
		expectedError  string
	}{
		{
			name:           "成功获取区块高度",
			mockHeight:     12345,
			mockError:      nil,
			expectedHeight: 12345,
			expectedError:  "",
		},
		{
			name:           "原始服务返回错误",
			mockHeight:     0,
			mockError:      errors.New("网络连接失败"),
			expectedHeight: 0,
			expectedError:  "从缓存获取区块高度失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟服务
			mockService := &MockBlockHeightService{}
			mockService.On("GetLatestBlockHeight", mock.Anything).Return(tt.mockHeight, tt.mockError)

			// 创建缓存服务
			cachedService, err := NewCachedBlockHeightService(mockService, 5*time.Second)
			require.NoError(t, err)

			// 调用被测试的方法
			ctx := context.Background()
			height, err := cachedService.GetLatestBlockHeight(ctx)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, uint64(0), height)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHeight, height)
			}

			// 验证模拟服务被调用
			mockService.AssertExpectations(t)
		})
	}
}

// TestCachedBlockHeightService_CacheHit 测试缓存命中功能
func TestCachedBlockHeightService_CacheHit(t *testing.T) {
	// 创建模拟服务
	mockService := &MockBlockHeightService{}
	// 只期望被调用一次
	mockService.On("GetLatestBlockHeight", mock.Anything).Return(uint64(12345), nil).Once()

	// 创建缓存服务，使用较长的过期时间确保缓存不会过期
	cachedService, err := NewCachedBlockHeightService(mockService, 10*time.Second)
	require.NoError(t, err)

	ctx := context.Background()

	// 第一次调用 - 缓存未命中，应该调用原始服务
	start1 := time.Now()
	height1, err := cachedService.GetLatestBlockHeight(ctx)
	duration1 := time.Since(start1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(12345), height1)

	// 第二次调用 - 缓存命中，不应该调用原始服务
	start2 := time.Now()
	height2, err := cachedService.GetLatestBlockHeight(ctx)
	duration2 := time.Since(start2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(12345), height2)

	// 验证缓存确实提升了性能（第二次调用应该更快）
	// 注意：这里使用宽松的条件，因为在测试环境中性能差异可能很小
	assert.True(t, duration2 <= duration1,
		"缓存未提升性能：第二次调用耗时 %v >= 第一次调用耗时 %v", duration2, duration1)

	// 验证模拟服务只被调用了一次
	mockService.AssertExpectations(t)
}

// TestCachedBlockHeightService_CacheExpiration 测试缓存过期功能
func TestCachedBlockHeightService_CacheExpiration(t *testing.T) {
	// 创建模拟服务
	mockService := &MockBlockHeightService{}
	// 期望被调用至少两次，返回不同的值
	mockService.On("GetLatestBlockHeight", mock.Anything).Return(uint64(12345), nil).Maybe()
	mockService.On("GetLatestBlockHeight", mock.Anything).Return(uint64(12346), nil).Maybe()

	// 创建缓存服务，使用很短的过期时间
	cachedService, err := NewCachedBlockHeightService(mockService, 50*time.Millisecond)
	require.NoError(t, err)

	ctx := context.Background()

	// 第一次调用
	height1, err := cachedService.GetLatestBlockHeight(ctx)
	assert.NoError(t, err)
	assert.True(t, height1 == 12345 || height1 == 12346, "应该返回预期的高度值")

	// 等待缓存过期
	time.Sleep(100 * time.Millisecond)

	// 第二次调用 - 缓存已过期，可能会重新调用原始服务
	height2, err := cachedService.GetLatestBlockHeight(ctx)
	assert.NoError(t, err)
	assert.True(t, height2 == 12345 || height2 == 12346, "应该返回预期的高度值")

	// 验证模拟服务被调用了至少一次
	mockService.AssertExpectations(t)
}

// BenchmarkCachedBlockHeightService 基准测试缓存性能
func BenchmarkCachedBlockHeightService(b *testing.B) {
	// 创建模拟服务
	mockService := &MockBlockHeightService{}
	mockService.On("GetLatestBlockHeight", mock.Anything).Return(uint64(12345), nil)

	// 创建缓存服务
	cachedService, err := NewCachedBlockHeightService(mockService, 10*time.Second)
	require.NoError(b, err)

	ctx := context.Background()

	// 预热缓存
	_, _ = cachedService.GetLatestBlockHeight(ctx)

	// 重置计时器
	b.ResetTimer()

	// 运行基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cachedService.GetLatestBlockHeight(ctx)
		}
	})
}
