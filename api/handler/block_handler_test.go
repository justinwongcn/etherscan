package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
func (m *MockBlockService) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
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
func (m *MockBlockService) GetBlockTransactionCount(ctx context.Context, blockHashOrNumber string) (uint64, error) {
	args := m.Called(ctx, blockHashOrNumber)
	return args.Get(0).(uint64), args.Error(1)
}

// GetTransactionCount mock实现
func (m *MockBlockService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (uint64, error) {
	args := m.Called(ctx, address, blockHashOrNumber)
	return args.Get(0).(uint64), args.Error(1)
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
		mockCount      uint64
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取区块交易数量",
			blockParam:     "12345",
			mockCount:      100,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(100)},
		},
		{
			name:           "使用latest标签获取交易数量",
			blockParam:     "latest",
			mockCount:      50,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(50)},
		},
		{
			name:           "使用区块哈希获取交易数量",
			blockParam:     "0x1234567890abcdef",
			mockCount:      75,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(75)},
		},
		{
			name:           "获取交易数量失败",
			blockParam:     "12345",
			mockCount:      0,
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
		mockHeight     uint64
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取区块高度",
			mockHeight:     12345,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"height": float64(12345)},
		},
		{
			name:           "获取区块高度失败",
			mockHeight:     0,
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

func TestGetTransactionCount(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		address        string
		blockParam     string
		mockCount      uint64
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取账户交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "latest",
			mockCount:      50,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(50)},
		},
		{
			name:           "使用latest标签获取交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "latest",
			mockCount:      50,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(50)},
		},
		{
			name:           "使用区块哈希获取交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "0x1234567890abcdef",
			mockCount:      75,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": float64(75)},
		},
		{
			name:           "地址为空",
			address:        "",
			blockParam:     "12345",
			mockCount:      0,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "address is required"},
		},
		{
			name:           "获取交易数量失败",
			address:        "0x1234567890abcdef",
			blockParam:     "12345",
			mockCount:      0,
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
			if tt.address != "" {
				mockService.On("GetTransactionCount", mock.Anything, tt.address, tt.blockParam).Return(tt.mockCount, tt.mockError)
			}

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{
				{Key: "address", Value: tt.address},
				{Key: "number", Value: tt.blockParam},
			}

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/account/%s/tx/count?number=%s", tt.address, tt.blockParam), nil)
			c.Request = req

			// 调用接口
			handler.GetTransactionCount(c)

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

func TestGetTransactionByHash(t *testing.T) {
	// 创建一个模拟的交易数据
	mockTransaction := &eth.Transaction{
		Hash:        *eth.MustData32("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		BlockHash:   eth.MustData32("0x0000000000000000000000000000000000000000000000000000000000000000"),
		BlockNumber: eth.MustQuantity("0x1"),
		From:        *eth.MustAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		Gas:         *eth.MustQuantity("0x5208"),
		Input:       eth.Input("0x"),
		Nonce:       *eth.MustQuantity("0x0"),
		To:          eth.MustAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		Index:       eth.MustQuantity("0x0"),
		Value:       *eth.MustQuantity("0x0"),
		V:           *eth.MustQuantity("0x1b"),
		R:           *eth.MustQuantity("0x0"),
		S:           *eth.MustQuantity("0x0"),
	}

	// 设置测试用例
	tests := []struct {
		name           string
		txHash         string
		mockTx         *eth.Transaction
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取交易信息",
			txHash:         "0x1234567890abcdef",
			mockTx:         mockTransaction,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"transaction": mockTransaction},
		},
		{
			name:           "交易哈希为空",
			txHash:         "",
			mockTx:         nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "transaction hash is required"},
		},
		{
			name:           "获取交易信息失败",
			txHash:         "0x1234567890abcdef",
			mockTx:         nil,
			mockError:      errors.New("failed to get transaction"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get transaction"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)
			if tt.txHash != "" {
				mockService.On("GetTransactionByHash", mock.Anything, tt.txHash).Return(tt.mockTx, tt.mockError)
			}

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{
				{Key: "hash", Value: tt.txHash},
			}

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, "/api/transaction/"+tt.txHash, nil)
			c.Request = req

			// 调用接口
			handler.GetTransactionByHash(c)

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

func TestGetTransactionByIndex(t *testing.T) {
	// 创建一个模拟的交易数据
	mockTransaction := &eth.Transaction{}

	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		index          string
		mockTx         *eth.Transaction
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "通过区块号和索引获取交易成功",
			blockParam:     "12345",
			index:          "0",
			mockTx:         mockTransaction,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"transaction": mockTransaction},
		},
		{
			name:           "通过区块哈希和索引获取交易成功",
			blockParam:     "0x1234567890abcdef",
			index:          "1",
			mockTx:         mockTransaction,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"transaction": mockTransaction},
		},
		{
			name:           "区块参数为空时使用latest",
			blockParam:     "",
			index:          "0",
			mockTx:         mockTransaction,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"transaction": mockTransaction},
		},
		{
			name:           "无效的交易索引",
			blockParam:     "12345",
			index:          "invalid",
			mockTx:         nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "invalid transaction index"},
		},
		{
			name:           "获取交易失败",
			blockParam:     "12345",
			index:          "0",
			mockTx:         nil,
			mockError:      errors.New("failed to get transaction"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get transaction"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)

			// 如果参数有效，设置mock期望
			if tt.index != "invalid" {
				index, _ := strconv.ParseUint(tt.index, 10, 64)
				expectedBlockParam := "latest"
				if tt.blockParam != "" {
					expectedBlockParam = tt.blockParam
				}
				mockService.On("GetTransactionByIndex", mock.Anything, expectedBlockParam, index).Return(tt.mockTx, tt.mockError)
			}

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{
				{Key: "index", Value: tt.index},
			}

			// 创建HTTP请求
			url := fmt.Sprintf("/api/block/tx/%s", tt.index)
			if tt.blockParam != "" {
				url = fmt.Sprintf("%s?number=%s", url, tt.blockParam)
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			c.Request = req

			// 调用接口
			handler.GetTransactionByIndex(c)

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
	mockBlock := &eth.Block{}

	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		mockBlock      *eth.Block
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

func TestGetTransactionReceipt(t *testing.T) {
	// 创建一个模拟的交易收据数据
	mockReceipt := &eth.TransactionReceipt{
		TransactionHash:   *eth.MustData32("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		BlockHash:         *eth.MustHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		BlockNumber:       *eth.MustQuantity("0x1"),
		ContractAddress:   eth.MustAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		CumulativeGasUsed: *eth.MustQuantity("0x5208"),
		From:              *eth.MustAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		GasUsed:           *eth.MustQuantity("0x5208"),
		Logs:              []eth.Log{},
		LogsBloom:         *eth.MustData256("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		Status:            eth.MustQuantity("0x1"),
		To:                eth.MustAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		TransactionIndex:  *eth.MustQuantity("0x0"),
	}

	// 设置测试用例
	tests := []struct {
		name           string
		txHash         string
		mockReceipt    *eth.TransactionReceipt
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取交易收据",
			txHash:         "0x1234567890abcdef",
			mockReceipt:    mockReceipt,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"receipt": mockReceipt},
		},
		{
			name:           "交易哈希为空",
			txHash:         "",
			mockReceipt:    nil,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "transaction hash is required"},
		},
		{
			name:           "获取交易收据失败",
			txHash:         "0x1234567890abcdef",
			mockReceipt:    nil,
			mockError:      errors.New("failed to get transaction receipt"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get transaction receipt"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)
			if tt.txHash != "" {
				mockService.On("GetTransactionReceipt", mock.Anything, tt.txHash).Return(tt.mockReceipt, tt.mockError)
			}

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{
				{Key: "hash", Value: tt.txHash},
			}

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodGet, "/api/transaction/"+tt.txHash+"/receipt", nil)
			c.Request = req

			// 调用接口
			handler.GetTransactionReceipt(c)

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

func TestSendRawTransaction(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		reqBody        map[string]string
		mockTxHash     string
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功发送交易",
			reqBody:        map[string]string{"signedTxData": "0x1234567890abcdef"},
			mockTxHash:     "0xabcdef1234567890",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"txHash": "0xabcdef1234567890"},
		},
		{
			name:           "请求体格式错误",
			reqBody:        map[string]string{},
			mockTxHash:     "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]any{"error": "Invalid request body: Key: 'SendRawTransactionRequest.SignedTxData' Error:Field validation for 'SignedTxData' failed on the 'required' tag"},
		},
		{
			name:           "发送交易失败",
			reqBody:        map[string]string{"signedTxData": "0x1234567890abcdef"},
			mockTxHash:     "",
			mockError:      errors.New("failed to send transaction"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to send transaction"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)
			if len(tt.reqBody) > 0 {
				mockService.On("SendRawTransaction", mock.Anything, tt.reqBody["signedTxData"]).Return(tt.mockTxHash, tt.mockError)
			}

			// 创建handler
			handler := NewBlockHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 创建请求体
			reqBodyBytes, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)

			// 创建HTTP请求
			req := httptest.NewRequest(http.MethodPost, "/api/tx/send", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// 调用接口
			handler.SendRawTransaction(c)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]any
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 将期望的响应对象转换为JSON字符串
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
