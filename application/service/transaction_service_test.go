package service

import (
	"context"
	"errors"
	"testing"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository 模拟交易仓储
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	args := m.Called(ctx, txHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (*eth.Transaction, error) {
	args := m.Called(ctx, blockHash, index)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber string, index uint64) (*eth.Transaction, error) {
	args := m.Called(ctx, blockNumber, index)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	args := m.Called(ctx, signedTxData)
	return args.String(0), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	args := m.Called(ctx, txHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eth.TransactionReceipt), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionCount(ctx context.Context, address string, blockParameter string) (uint64, error) {
	args := m.Called(ctx, address, blockParameter)
	return args.Get(0).(uint64), args.Error(1)
}

// TestNewTransactionService 测试构造函数
func TestNewTransactionService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "创建TransactionService成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTransactionRepository{}
			service := NewTransactionService(mockRepo)

			assert.NotNil(t, service)
			assert.Equal(t, mockRepo, service.transactionRepo)
		})
	}
}

// TestTransactionService_GetTransactionByHash 测试通过哈希获取交易
func TestTransactionService_GetTransactionByHash(t *testing.T) {
	tests := []struct {
		name        string
		txHash      string
		mockTx      *eth.Transaction
		mockError   error
		expectError bool
	}{
		{
			name:   "获取交易成功",
			txHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			mockTx: &eth.Transaction{
				Hash:  *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Nonce: *eth.MustQuantity("0x1"),
				Value: *eth.MustQuantity("0x1000"),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "交易不存在",
			txHash:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			mockTx:      nil,
			mockError:   errors.New("交易不存在"),
			expectError: true,
		},
		{
			name:        "无效的交易哈希",
			txHash:      "invalid-hash",
			mockTx:      nil,
			mockError:   errors.New("无效的交易哈希"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTransactionRepository{}
			mockRepo.On("GetTransactionByHash", mock.Anything, tt.txHash).Return(tt.mockTx, tt.mockError)

			service := NewTransactionService(mockRepo)
			ctx := context.Background()

			result, err := service.GetTransactionByHash(ctx, tt.txHash)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTx, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTransactionService_SendRawTransaction 测试发送原始交易
func TestTransactionService_SendRawTransaction(t *testing.T) {
	tests := []struct {
		name         string
		signedTxData string
		mockTxHash   string
		mockError    error
		expectError  bool
	}{
		{
			name:         "发送交易成功",
			signedTxData: "0x1234567890abcdef",
			mockTxHash:   "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			mockError:    nil,
			expectError:  false,
		},
		{
			name:         "发送交易失败",
			signedTxData: "0x1234567890abcdef",
			mockTxHash:   "",
			mockError:    errors.New("交易数据无效"),
			expectError:  true,
		},
		{
			name:         "空交易数据",
			signedTxData: "",
			mockTxHash:   "",
			mockError:    errors.New("交易数据不能为空"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTransactionRepository{}
			mockRepo.On("SendRawTransaction", mock.Anything, tt.signedTxData).Return(tt.mockTxHash, tt.mockError)

			service := NewTransactionService(mockRepo)
			ctx := context.Background()

			result, err := service.SendRawTransaction(ctx, tt.signedTxData)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTxHash, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTransactionService_GetTransactionReceipt 测试获取交易收据
func TestTransactionService_GetTransactionReceipt(t *testing.T) {
	tests := []struct {
		name        string
		txHash      string
		mockReceipt *eth.TransactionReceipt
		mockError   error
		expectError bool
	}{
		{
			name:   "获取交易收据成功",
			txHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			mockReceipt: &eth.TransactionReceipt{
				TransactionHash: *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Status:          eth.MustQuantity("0x1"),
				GasUsed:         *eth.MustQuantity("0x5208"),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "交易收据不存在",
			txHash:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			mockReceipt: nil,
			mockError:   errors.New("交易收据不存在"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTransactionRepository{}
			mockRepo.On("GetTransactionReceipt", mock.Anything, tt.txHash).Return(tt.mockReceipt, tt.mockError)

			service := NewTransactionService(mockRepo)
			ctx := context.Background()

			result, err := service.GetTransactionReceipt(ctx, tt.txHash)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReceipt, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTransactionService_GetTransactionCount 测试获取交易数量
func TestTransactionService_GetTransactionCount(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		blockParameter string
		mockCount      uint64
		mockError      error
		expectError    bool
		expected       string
	}{
		{
			name:           "获取交易数量成功",
			address:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			blockParameter: "latest",
			mockCount:      100,
			mockError:      nil,
			expectError:    false,
			expected:       "100",
		},
		{
			name:           "获取交易数量失败",
			address:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			blockParameter: "latest",
			mockCount:      0,
			mockError:      errors.New("地址无效"),
			expectError:    true,
			expected:       "",
		},
		{
			name:           "零交易数量",
			address:        "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			blockParameter: "latest",
			mockCount:      0,
			mockError:      nil,
			expectError:    false,
			expected:       "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTransactionRepository{}
			mockRepo.On("GetTransactionCount", mock.Anything, tt.address, tt.blockParameter).Return(tt.mockCount, tt.mockError)

			service := NewTransactionService(mockRepo)
			ctx := context.Background()

			result, err := service.GetTransactionCount(ctx, tt.address, tt.blockParameter)

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

// BenchmarkTransactionService_GetTransactionByHash 性能测试
func BenchmarkTransactionService_GetTransactionByHash(b *testing.B) {
	mockRepo := &MockTransactionRepository{}
	mockTx := &eth.Transaction{
		Hash:  *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		Nonce: *eth.MustQuantity("0x1"),
		Value: *eth.MustQuantity("0x1000"),
	}
	mockRepo.On("GetTransactionByHash", mock.Anything, mock.AnythingOfType("string")).Return(mockTx, nil)

	service := NewTransactionService(mockRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetTransactionByHash(ctx, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	}
}

// BenchmarkTransactionService_SendRawTransaction 性能测试
func BenchmarkTransactionService_SendRawTransaction(b *testing.B) {
	mockRepo := &MockTransactionRepository{}
	mockRepo.On("SendRawTransaction", mock.Anything, mock.AnythingOfType("string")).Return("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", nil)

	service := NewTransactionService(mockRepo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.SendRawTransaction(ctx, "0x1234567890abcdef")
	}
}
