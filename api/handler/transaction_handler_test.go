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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService 是TransactionServiceInterface的mock实现
type MockTransactionService struct {
	mock.Mock
}

// 确保MockTransactionService实现了TransactionServiceInterface接口
var _ service.TransactionServiceInterface = (*MockTransactionService)(nil)

// GetTransactionCount mock实现
func (m *MockTransactionService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (string, error) {
	args := m.Called(ctx, address, blockHashOrNumber)
	return args.Get(0).(string), args.Error(1)
}

// GetTransactionByHash mock实现
func (m *MockTransactionService) GetTransactionByHash(ctx context.Context, txHash string) (*domain.Transaction, error) {
	args := m.Called(ctx, txHash)
	if tx, ok := args.Get(0).(*domain.Transaction); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetTransactionByIndex mock实现
func (m *MockTransactionService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*domain.Transaction, error) {
	args := m.Called(ctx, blockHashOrNumber, index)
	if tx, ok := args.Get(0).(*domain.Transaction); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

// SendRawTransaction mock实现
func (m *MockTransactionService) SendRawTransaction(ctx context.Context, signedTxData string) (string, error) {
	args := m.Called(ctx, signedTxData)
	return args.String(0), args.Error(1)
}

// GetTransactionReceipt mock实现
func (m *MockTransactionService) GetTransactionReceipt(ctx context.Context, txHash string) (*domain.TransactionReceipt, error) {
	args := m.Called(ctx, txHash)
	if receipt, ok := args.Get(0).(*domain.TransactionReceipt); ok {
		return receipt, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetTransactionCount(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		address        string
		blockParam     string
		mockCount      string
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "成功获取账户交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "latest",
			mockCount:      "50",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "50"},
		},
		{
			name:           "使用latest标签获取交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "latest",
			mockCount:      "50",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "50"},
		},
		{
			name:           "使用区块哈希获取交易数量",
			address:        "0x1234567890abcdef",
			blockParam:     "0x1234567890abcdef",
			mockCount:      "75",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"count": "75"},
		},
		{
			name:           "地址为空",
			address:        "",
			blockParam:     "12345",
			mockCount:      "0",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "address is required"},
		},
		{
			name:           "获取交易数量失败",
			address:        "0x1234567890abcdef",
			blockParam:     "12345",
			mockCount:      "0",
			mockError:      errors.New("failed to get transaction count"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"error": "failed to get transaction count"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockTransactionService)
			if tt.address != "" {
				mockService.On("GetTransactionCount", mock.Anything, tt.address, tt.blockParam).Return(tt.mockCount, tt.mockError)
			}

			// 创建handler
			handler := NewTransactionHandler(mockService)

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
	blockHash := "0x0000000000000000000000000000000000000000000000000000000000000000"
	blockNumber := "1"
	from := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	to := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	hash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	index := "0"
	typ := "0"

	mockTransaction := &domain.Transaction{
		Type:        &typ,
		BlockHash:   &blockHash,
		BlockNumber: &blockNumber,
		From:        from,
		Gas:         "21000",
		Hash:        hash,
		Input:       "0x",
		Nonce:       "0",
		To:          &to,
		Index:       &index,
		Value:       "0",
		V:           "27",
		R:           "0",
		S:           "0",
	}

	// 设置测试用例
	tests := []struct {
		name           string
		txHash         string
		mockTx         *domain.Transaction
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
			mockService := new(MockTransactionService)
			if tt.txHash != "" {
				mockService.On("GetTransactionByHash", mock.Anything, tt.txHash).Return(tt.mockTx, tt.mockError)
			}

			// 创建handler
			handler := NewTransactionHandler(mockService)

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
	blockHash := "0x0000000000000000000000000000000000000000000000000000000000000000"
	blockNumber := "1"
	from := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	to := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	hash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	index := "0"
	typ := "0"

	mockTransaction := &domain.Transaction{
		Type:        &typ,
		BlockHash:   &blockHash,
		BlockNumber: &blockNumber,
		From:        from,
		Gas:         "21000",
		Hash:        hash,
		Input:       "0x",
		Nonce:       "0",
		To:          &to,
		Index:       &index,
		Value:       "0",
		V:           "27",
		R:           "0",
		S:           "0",
	}

	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		index          string
		mockTx         *domain.Transaction
		mockError      error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "区块参数缺失",
			blockParam:     "",
			index:          "0",
			mockTx:         nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]any{"error": "block number or hash is required"},
		},
		{
			name:           "通过区块号查询",
			blockParam:     "12345",
			index:          "0",
			mockTx:         mockTransaction,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"transaction": mockTransaction},
		},
		{
			name:           "通过区块哈希查询",
			blockParam:     "0x1234567890abcdef",
			index:          "1",
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
			mockService := new(MockTransactionService)

			// 如果参数有效，设置mock期望
			if tt.blockParam != "" && tt.index != "invalid" {
				index, _ := strconv.ParseUint(tt.index, 10, 64)
				mockService.On("GetTransactionByIndex", mock.Anything, tt.blockParam, index).Return(tt.mockTx, tt.mockError)
			}

			// 创建handler
			handler := NewTransactionHandler(mockService)

			// 创建gin测试环境
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置路由参数
			c.Params = gin.Params{
				{Key: "index", Value: tt.index},
			}

			// 创建HTTP请求
			reqURL := fmt.Sprintf("/api/block/tx/%s", tt.index)
			if tt.blockParam != "" {
				reqURL += fmt.Sprintf("?number=%s", tt.blockParam)
			}
			req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
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

func TestGetTransactionReceipt(t *testing.T) {
	status := "1"
	contractAddress := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	// 创建一个模拟的交易收据数据
	mockReceipt := &domain.TransactionReceipt{
		TransactionHash:   "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000000000",
		BlockNumber:       "1",
		ContractAddress:   &contractAddress,
		CumulativeGasUsed: "21000",
		From:              "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		GasUsed:           "21000",
		Logs:              []domain.Log{},
		LogsBloom:         "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Status:            &status,
		To:                &contractAddress,
		TransactionIndex:  "0",
	}

	// 设置测试用例
	tests := []struct {
		name           string
		txHash         string
		mockReceipt    *domain.TransactionReceipt
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
			mockService := new(MockTransactionService)
			if tt.txHash != "" {
				mockService.On("GetTransactionReceipt", mock.Anything, tt.txHash).Return(tt.mockReceipt, tt.mockError)
			}

			// 创建handler
			handler := NewTransactionHandler(mockService)

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
			mockService := new(MockTransactionService)
			if len(tt.reqBody) > 0 {
				mockService.On("SendRawTransaction", mock.Anything, tt.reqBody["signedTxData"]).Return(tt.mockTxHash, tt.mockError)
			}

			// 创建handler
			handler := NewTransactionHandler(mockService)

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
