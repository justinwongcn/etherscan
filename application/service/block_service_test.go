package service

import (
	"context"
	"errors"
	"testing"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlockRepository 模拟区块仓储
type MockBlockRepository struct {
	mock.Mock
}

func (m *MockBlockRepository) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockRepository) GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error) {
	args := m.Called(ctx, hash, fullTx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Block), args.Error(1)
}

func (m *MockBlockRepository) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error) {
	args := m.Called(ctx, number, fullTx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Block), args.Error(1)
}

func (m *MockBlockRepository) GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockRepository) GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(uint64), args.Error(1)
}

// TestNewBlockService 测试构造函数
func TestNewBlockService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "创建BlockService成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockBlockRepository{}
			service := NewBlockService(mockRepo)

			assert.NotNil(t, service)
			assert.Equal(t, mockRepo, service.blockRepo)
		})
	}
}

// TestBlockService_GetLatestBlockHeight 测试获取最新区块高度
func TestBlockService_GetLatestBlockHeight(t *testing.T) {
	tests := []struct {
		name        string
		mockHeight  uint64
		mockError   error
		expectError bool
		expected    string
	}{
		{
			name:        "获取最新区块高度成功",
			mockHeight:  1000000,
			mockError:   nil,
			expectError: false,
			expected:    "1000000",
		},
		{
			name:        "获取最新区块高度失败",
			mockHeight:  0,
			mockError:   errors.New("网络错误"),
			expectError: true,
			expected:    "",
		},
		{
			name:        "零区块高度",
			mockHeight:  0,
			mockError:   nil,
			expectError: false,
			expected:    "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockBlockRepository{}
			mockRepo.On("GetLatestBlockNumber", mock.Anything).Return(tt.mockHeight, tt.mockError)

			service := NewBlockService(mockRepo)
			ctx := context.Background()

			result, err := service.GetLatestBlockHeight(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestBlockService_GetBlock 测试获取区块信息
func TestBlockService_GetBlock(t *testing.T) {
	tests := []struct {
		name               string
		blockHashOrNumber  string
		fullTx             bool
		mockBlock          *eth.Block
		mockError          error
		expectError        bool
		expectedRepoMethod string
	}{
		{
			name:              "通过区块号获取区块",
			blockHashOrNumber: "latest",
			fullTx:            true,
			mockBlock: &eth.Block{
				Number: eth.MustQuantity("0x1"),
				Hash:   eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			},
			mockError:          nil,
			expectError:        false,
			expectedRepoMethod: "GetBlockByNumber",
		},
		{
			name:              "通过区块哈希获取区块",
			blockHashOrNumber: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			fullTx:            false,
			mockBlock: &eth.Block{
				Number: eth.MustQuantity("0x1"),
				Hash:   eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			},
			mockError:          nil,
			expectError:        false,
			expectedRepoMethod: "GetBlockByHash",
		},
		{
			name:               "获取区块失败",
			blockHashOrNumber:  "latest",
			fullTx:             true,
			mockBlock:          nil,
			mockError:          errors.New("区块不存在"),
			expectError:        true,
			expectedRepoMethod: "GetBlockByNumber",
		},
		{
			name:               "无效的区块参数",
			blockHashOrNumber:  "invalid",
			fullTx:             true,
			mockBlock:          nil,
			mockError:          errors.New("无效参数"),
			expectError:        true,
			expectedRepoMethod: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockBlockRepository{}

			switch tt.expectedRepoMethod {
			case "GetBlockByNumber":
				mockRepo.On("GetBlockByNumber", mock.Anything, mock.AnythingOfType("string"), tt.fullTx).Return(tt.mockBlock, tt.mockError)
			case "GetBlockByHash":
				mockRepo.On("GetBlockByHash", mock.Anything, mock.AnythingOfType("string"), tt.fullTx).Return(tt.mockBlock, tt.mockError)
			}

			service := NewBlockService(mockRepo)
			ctx := context.Background()

			result, err := service.GetBlock(ctx, tt.blockHashOrNumber, tt.fullTx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockBlock, result)
			}

			if tt.expectedRepoMethod != "" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

// TestBlockService_GetBlockTransactionCount 测试获取区块交易数量
func TestBlockService_GetBlockTransactionCount(t *testing.T) {
	tests := []struct {
		name              string
		blockHashOrNumber string
		expectError       bool
	}{
		{
			name:              "通过区块号获取交易数量",
			blockHashOrNumber: "latest",
			expectError:       false,
		},
		{
			name:              "通过区块哈希获取交易数量",
			blockHashOrNumber: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expectError:       false,
		},
		{
			name:              "获取交易数量失败",
			blockHashOrNumber: "latest",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockBlockRepository{}

			// 简化测试，只验证方法调用
			if len(tt.blockHashOrNumber) == 66 && tt.blockHashOrNumber[:2] == "0x" {
				// 区块哈希
				if tt.expectError {
					mockRepo.On("GetBlockByHash", mock.Anything, mock.AnythingOfType("string"), true).Return(nil, errors.New("区块不存在"))
				} else {
					// 创建一个包含必要字段的区块
					hash := eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
					parentHash := eth.MustHash("0x0234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
					miner := eth.MustAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
					number := eth.MustQuantity("0x1")
					timestamp := eth.MustQuantity("0x60000000")

					mockBlock := &eth.Block{
						Hash:       hash,
						ParentHash: *parentHash,
						Miner:      *miner,
						Number:     number,
						Timestamp:  *timestamp,
					}
					mockRepo.On("GetBlockByHash", mock.Anything, mock.AnythingOfType("string"), true).Return(mockBlock, nil)
				}
			} else {
				// 区块号
				if tt.expectError {
					mockRepo.On("GetBlockByNumber", mock.Anything, mock.AnythingOfType("string"), true).Return(nil, errors.New("区块不存在"))
				} else {
					// 创建一个包含必要字段的区块
					hash := eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
					parentHash := eth.MustHash("0x0234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
					miner := eth.MustAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
					number := eth.MustQuantity("0x1")
					timestamp := eth.MustQuantity("0x60000000")

					mockBlock := &eth.Block{
						Hash:       hash,
						ParentHash: *parentHash,
						Miner:      *miner,
						Number:     number,
						Timestamp:  *timestamp,
					}
					mockRepo.On("GetBlockByNumber", mock.Anything, mock.AnythingOfType("string"), true).Return(mockBlock, nil)
				}
			}

			service := NewBlockService(mockRepo)
			ctx := context.Background()

			result, err := service.GetBlockTransactionCount(ctx, tt.blockHashOrNumber)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// BenchmarkBlockService_GetLatestBlockHeight 性能测试
func BenchmarkBlockService_GetLatestBlockHeight(b *testing.B) {
	mockRepo := &MockBlockRepository{}
	mockRepo.On("GetLatestBlockNumber", mock.Anything).Return(uint64(1000000), nil)

	service := NewBlockService(mockRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetLatestBlockHeight(ctx)
	}
}

// BenchmarkBlockService_GetBlock 性能测试
func BenchmarkBlockService_GetBlock(b *testing.B) {
	mockRepo := &MockBlockRepository{}
	mockBlock := &eth.Block{
		Number: eth.MustQuantity("0x1"),
		Hash:   eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	}
	mockRepo.On("GetBlockByNumber", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(mockBlock, nil)

	service := NewBlockService(mockRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetBlock(ctx, "latest", true)
	}
}
