package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountService 模拟账户服务
type MockAccountService struct {
	mock.Mock
}

// GetBalance 实现AccountServiceInterface接口的GetBalance方法
func (m *MockAccountService) GetBalance(ctx context.Context, address string, numberOrTag string) (string, error) {
	args := m.Called(ctx, address, numberOrTag)
	return args.String(0), args.Error(1)
}

// setupTest 设置测试环境
func setupTest() (*gin.Engine, *MockAccountService) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由引擎
	r := gin.New()

	// 创建模拟服务
	mockService := new(MockAccountService)

	// 创建处理器
	handler := NewAccountHandler(mockService)

	// 注册路由
	r.GET("/accounts/:address/balance", handler.GetBalance)

	return r, mockService
}

func TestGetBalance(t *testing.T) {
	// 测试用例
	tests := []struct {
		name           string
		address        string
		block          string
		mockBalance    string
		mockError      error
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "成功获取余额",
			address:        "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			block:          "latest",
			mockBalance:    "1000000000000000000",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"balance": "1000000000000000000"},
		},
		{
			name:           "地址参数缺失",
			address:        "",
			block:          "latest",
			mockBalance:    "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "address is required"},
		},
		{
			name:           "无效的地址格式",
			address:        "invalid_address",
			block:          "latest",
			mockBalance:    "",
			mockError:      errors.New("invalid ethereum address"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "invalid ethereum address"},
		},
		{
			name:           "无效的区块参数",
			address:        "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			block:          "invalid_block",
			mockBalance:    "",
			mockError:      errors.New("invalid block parameter"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "invalid block parameter"},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			r, mockService := setupTest()

			// 设置模拟服务的预期行为
			if tt.address != "" {
				mockService.On("GetBalance", mock.Anything, tt.address, tt.block).Return(tt.mockBalance, tt.mockError)
			}

			// 创建测试请求
			w := httptest.NewRecorder()
			url := "/accounts/" + tt.address + "/balance"
			if tt.address == "" {
				url = "/accounts//balance"
			}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			q := req.URL.Query()
			if tt.block != "" {
				q.Add("block", tt.block)
			}
			req.URL.RawQuery = q.Encode()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				expectedJSON := gin.H{}
				for k, v := range tt.expectedBody {
					expectedJSON[k] = v
				}
				actualJSON := gin.H{}
				for k, v := range tt.expectedBody {
					actualJSON[k] = v
				}
				assert.Equal(t, expectedJSON, actualJSON)
			}

			// 验证模拟服务的调用
			mockService.AssertExpectations(t)
		})
	}
}
