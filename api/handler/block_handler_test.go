package handler

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/application/service"
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

// GetBlockByNumber mock实现
func (m *MockBlockService) GetBlockByNumber(ctx context.Context, numberOrTag string) (*eth.Block, error) {
	args := m.Called(ctx, numberOrTag)
	if block, ok := args.Get(0).(*eth.Block); ok {
		return block, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetBlockByHash mock实现
func (m *MockBlockService) GetBlockByHash(ctx context.Context, blockHash string) (*eth.Block, error) {
	args := m.Called(ctx, blockHash)
	if block, ok := args.Get(0).(*eth.Block); ok {
		return block, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetBlockHeight(t *testing.T) {
	// 设置测试用例
	tests := []struct {
		name           string
		mockHeight     uint64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "成功获取区块高度",
			mockHeight:     12345,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"height": float64(12345)},
		},
		{
			name:           "获取区块高度失败",
			mockHeight:     0,
			mockError:      errors.New("failed to get block height"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to get block height"},
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

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// 验证mock调用
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetBlockByNumberOrHash(t *testing.T) {
	// 创建一个模拟的区块数据
	mockBlock := &eth.Block{}

	// 设置测试用例
	tests := []struct {
		name           string
		blockParam     string
		mockBlock      *eth.Block
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "通过区块哈希获取区块成功",
			blockParam:     "0x1234567890abcdef",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"block": mockBlock},
		},
		{
			name:           "通过区块号获取区块成功",
			blockParam:     "12345",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"block": mockBlock},
		},
		{
			name:           "通过latest标签获取区块成功",
			blockParam:     "latest",
			mockBlock:      mockBlock,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"block": mockBlock},
		},
		{
			name:           "获取区块失败",
			blockParam:     "0x1234567890abcdef",
			mockBlock:      nil,
			mockError:      errors.New("failed to get block"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to get block"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock service
			mockService := new(MockBlockService)

			// 根据参数类型设置mock期望
			if len(tt.blockParam) >= 2 && tt.blockParam[:2] == "0x" {
				mockService.On("GetBlockByHash", mock.Anything, tt.blockParam).Return(tt.mockBlock, tt.mockError)
			} else if number, ok := new(big.Int).SetString(tt.blockParam, 10); ok {
				numberOrTag := "0x" + number.Text(16)
				mockService.On("GetBlockByNumber", mock.Anything, numberOrTag).Return(tt.mockBlock, tt.mockError)
			} else {
				mockService.On("GetBlockByNumber", mock.Anything, tt.blockParam).Return(tt.mockBlock, tt.mockError)
			}

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
			handler.GetBlockByNumberOrHash(c)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// 验证mock调用
			mockService.AssertExpectations(t)
		})
	}
}
