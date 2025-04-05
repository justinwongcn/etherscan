package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/domain"
	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlockService 是BlockServiceInterface的mock实现
type MockBlockService struct {
	mock.Mock
}

// 确保MockBlockService实现了BlockServiceInterface接口
var _ service.BlockServiceInterface = (*MockBlockService)(nil)

// GetLatestBlockHeight mock实现
func (m *MockBlockService) GetLatestBlockHeight(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

// GetBlock mock实现
func (m *MockBlockService) GetBlock(ctx context.Context, blockHashOrNumber string) (*domain.Block, error) {
	args := m.Called(ctx, blockHashOrNumber)
	if block, ok := args.Get(0).(*domain.Block); ok {
		return block, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetTransactionCount mock实现
// GetBlockTransactionCount mock实现
func (m *MockBlockService) GetBlockTransactionCount(ctx context.Context, blockHashOrNumber string) (string, error) {
	args := m.Called(ctx, blockHashOrNumber)
	return args.String(0), args.Error(1)
}

// GetTransactionCount mock实现
func (m *MockBlockService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (string, error) {
	args := m.Called(ctx, address, blockHashOrNumber)
	return args.String(0), args.Error(1)
}

// GetTransactionByHash mock实现
func (m *MockBlockService) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	args := m.Called(ctx, txHash)
	if tx, ok := args.Get(0).(*eth.Transaction); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetTransactionByIndex mock实现
func (m *MockBlockService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*eth.Transaction, error) {
	args := m.Called(ctx, blockHashOrNumber, index)
	if tx, ok := args.Get(0).(*eth.Transaction); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

// SendRawTransaction mock实现
func (m *MockBlockService) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	args := m.Called(ctx, signedTxData)
	return args.String(0), args.Error(1)
}

// GetTransactionReceipt mock实现
func (m *MockBlockService) GetTransactionReceipt(ctx context.Context, txHash string) (*eth.TransactionReceipt, error) {
	args := m.Called(ctx, txHash)
	if receipt, ok := args.Get(0).(*eth.TransactionReceipt); ok {
		return receipt, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetBlockTransactionCount(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		mockCount      string
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取区块交易数量",
			blockParam:     "12345",
			mockCount:      "100",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "100"},
		},
		{
			name:           "使用latest标签获取交易数量",
			blockParam:     "latest",
			mockCount:      "50",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "50"},
		},
		{
			name:           "使用区块哈希获取交易数量",
			blockParam:     "0x1234567890abcdef",
			mockCount:      "75",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "75"},
		},
		{
			name:           "获取交易数量失败",
			blockParam:     "12345",
			mockCount:      "",
			mockError:      errors.New("failed to get transaction count"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get transaction count"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)
			mockService.On("GetBlockTransactionCount", mock.Anything, tt.blockParam).Return(tt.mockCount, tt.mockError)

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{{Key: "number", Value: tt.blockParam}}

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, "/api/block/"+tt.blockParam+"/txcount", nil)
			c.Request = req

			// 调用接口
			handler.GetBlockTransactionCount(c)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 将期望的交易对象转换为JSON字符串
			expectedJSON, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			// 将实际响应转换为JSON字符串
			actualJSON, err := json.Marshal(response)
			assert.NoError(t, err)

			// 比较JSON字符串
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			// 验证mock调用
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetBlockHeight(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		mockHeight     string
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取区块高度",
			mockHeight:     "12345",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"height": "12345"},
		},
		{
			name:           "获取区块高度失败",
			mockHeight:     "",
			mockError:      errors.New("failed to get block height"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get block height"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)
			mockService.On("GetLatestBlockHeight", mock.Anything).Return(tt.mockHeight, tt.mockError)

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, "/api/block/height", nil)
			c.Request = req

			// 调用接口
			handler.GetBlockHeight(c)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 将期望的交易对象转换为JSON字符串
			expectedJSON, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			// 将实际响应转换为JSON字符串
			actualJSON, err := json.Marshal(response)
			assert.NoError(t, err)

			// 比较JSON字符串
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			// 验证mock调用
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetBlock(t *testing.T) {
	// 创建一个模拟的区块数据
	mockBlock := &domain.Block{}

	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		mockBlock      *domain.Block
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "通过区块哈希获取区块成功",
			blockParam:     "0x1234567890abcdef",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"block": mockBlock},
		},
		{
			name:           "通过区块号获取区块成功",
			blockParam:     "12345",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"block": mockBlock},
		},
		{
			name:           "通过latest标签获取区块成功",
			blockParam:     "latest",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"block": mockBlock},
		},
		{
			name:           "获取区块失败",
			blockParam:     "0x1234567890abcdef",
			mockBlock:      nil,
			mockError:      errors.New("failed to get block"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get block"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)

			// 设置mock期望
			mockService.On("GetBlock", mock.Anything, tt.blockParam).Return(tt.mockBlock, tt.mockError)

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{{Key: "number", Value: tt.blockParam}}

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, "/api/block/"+tt.blockParam, nil)
			c.Request = req

			// 调用接口
			handler.GetBlock(c)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 将期望的交易对象转换为JSON字符串
			expectedJSON, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			// 将实际响应转换为JSON字符串
			actualJSON, err := json.Marshal(response)
			assert.NoError(t, err)

			// 比较JSON字符串
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			// 验证mock调用
			mockService.AssertExpectations(t)
		})
	}
}
