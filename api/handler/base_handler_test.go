package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/stretchr/testify/assert"
)

// TestNewBaseHandler 测试BaseHandler构造函数
func TestNewBaseHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "创建BaseHandler成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBaseHandler()
			assert.NotNil(t, handler)
		})
	}
}

// TestBaseHandler_RespondWithError 测试错误响应
func TestBaseHandler_RespondWithError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "400错误响应",
			statusCode: http.StatusBadRequest,
			message:    "请求参数无效",
		},
		{
			name:       "500错误响应",
			statusCode: http.StatusInternalServerError,
			message:    "服务器内部错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求和响应
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 创建ant.Context
			ctx := &ant.Context{
				Req:  req,
				Resp: w,
			}

			// 创建BaseHandler并调用方法
			handler := NewBaseHandler()
			handler.RespondWithError(ctx, tt.statusCode, tt.message)

			// 验证响应状态码
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

// TestBaseHandler_RespondWithSuccess 测试成功响应
func TestBaseHandler_RespondWithSuccess(t *testing.T) {
	tests := []struct {
		name string
		data H
	}{
		{
			name: "成功响应",
			data: H{
				"result": "success",
				"count":  10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求和响应
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 创建ant.Context
			ctx := &ant.Context{
				Req:  req,
				Resp: w,
			}

			// 创建BaseHandler并调用方法
			handler := NewBaseHandler()
			handler.RespondWithSuccess(ctx, tt.data)

			// 验证响应状态码
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestBaseHandler_RespondWithInternalError 测试内部错误响应
func TestBaseHandler_RespondWithInternalError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "内部错误响应",
			err:  assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求和响应
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 创建ant.Context
			ctx := &ant.Context{
				Req:  req,
				Resp: w,
			}

			// 创建BaseHandler并调用方法
			handler := NewBaseHandler()
			handler.RespondWithInternalError(ctx, tt.err)

			// 验证响应状态码
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}

// TestBaseHandler_RespondWithBadRequest 测试错误请求响应
func TestBaseHandler_RespondWithBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "错误请求响应",
			message: "参数格式错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求和响应
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 创建ant.Context
			ctx := &ant.Context{
				Req:  req,
				Resp: w,
			}

			// 创建BaseHandler并调用方法
			handler := NewBaseHandler()
			handler.RespondWithBadRequest(ctx, tt.message)

			// 验证响应状态码
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

// TestBaseHandler_ConvertBlockToResponse 测试区块转换
func TestBaseHandler_ConvertBlockToResponse(t *testing.T) {
	tests := []struct {
		name     string
		ethBlock *eth.Block
		fullTx   bool
		wantNil  bool
	}{
		{
			name:     "nil区块",
			ethBlock: nil,
			fullTx:   false,
			wantNil:  true,
		},
		{
			name: "有效区块_简单交易",
			ethBlock: &eth.Block{
				Number:    eth.MustQuantity("0x1"),
				Hash:      eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				GasLimit:  *eth.MustQuantity("0x1000000"),
				GasUsed:   *eth.MustQuantity("0x500000"),
				Timestamp: *eth.MustQuantity("0x60000000"),
			},
			fullTx:  false,
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBaseHandler()
			result := handler.ConvertBlockToResponse(tt.ethBlock, tt.fullTx)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// 验证数值字段是十进制格式
				if tt.ethBlock != nil {
					assert.Equal(t, "1", result["number"])
					assert.Equal(t, "16777216", result["gasLimit"])
					assert.Equal(t, "5242880", result["gasUsed"])
					assert.Equal(t, "1610612736", result["timestamp"])
				}
			}
		})
	}
}

// TestBaseHandler_ConvertTransactionToResponse 测试交易转换
func TestBaseHandler_ConvertTransactionToResponse(t *testing.T) {
	tests := []struct {
		name    string
		ethTx   *eth.Transaction
		wantNil bool
	}{
		{
			name:    "nil交易",
			ethTx:   nil,
			wantNil: true,
		},
		{
			name: "有效交易",
			ethTx: &eth.Transaction{
				Hash:     *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Nonce:    *eth.MustQuantity("0x1"),
				Value:    *eth.MustQuantity("0x1000"),
				Gas:      *eth.MustQuantity("0x5208"),
				GasPrice: eth.MustQuantity("0x3b9aca00"),
				V:        *eth.MustQuantity("0x1b"),
				R:        *eth.MustQuantity("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				S:        *eth.MustQuantity("0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"),
				From:     *eth.MustAddress("0x8ba1f109551bD432803012645aac136c12345678"),
				Input:    eth.Input(*eth.MustData("0x")),
				Type:     eth.MustQuantity("0x0"),
			},
			wantNil: false,
		},
		{
			name: "blob交易",
			ethTx: &eth.Transaction{
				Hash:             *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Nonce:            *eth.MustQuantity("0x1"),
				Value:            *eth.MustQuantity("0x1000"),
				Gas:              *eth.MustQuantity("0x5208"),
				MaxFeePerBlobGas: eth.MustQuantity("0x1000000000"),
				V:                *eth.MustQuantity("0x1b"),
				R:                *eth.MustQuantity("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				S:                *eth.MustQuantity("0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"),
				From:             *eth.MustAddress("0x8ba1f109551bD432803012645aac136c12345678"),
				Input:            eth.Input(*eth.MustData("0x")),
				Type:             eth.MustQuantity("0x3"),
				BlobVersionedHashes: []eth.Hash{
					*eth.MustHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
				},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBaseHandler()
			result := handler.ConvertTransactionToResponse(tt.ethTx)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// 验证数值字段是十进制格式
				if tt.ethTx != nil {
					assert.Equal(t, "1", result["nonce"])
					assert.Equal(t, "4096", result["value"])
					assert.Equal(t, "21000", result["gas"])

					// 检查blob交易特有字段
					if tt.name == "blob交易" {
						assert.Equal(t, "68719476736", result["maxFeePerBlobGas"])
						assert.Contains(t, result, "blobVersionedHashes")
						blobHashes, ok := result["blobVersionedHashes"].([]string)
						assert.True(t, ok)
						assert.Len(t, blobHashes, 1)
					} else {
						assert.Equal(t, "1000000000", result["gasPrice"])
					}
				}
			}
		})
	}
}

// TestBaseHandler_ConvertTransactionReceiptToResponse 测试交易收据转换
func TestBaseHandler_ConvertTransactionReceiptToResponse(t *testing.T) {
	tests := []struct {
		name       string
		ethReceipt *eth.TransactionReceipt
		wantNil    bool
	}{
		{
			name:       "nil收据",
			ethReceipt: nil,
			wantNil:    true,
		},
		{
			name: "有效收据",
			ethReceipt: &eth.TransactionReceipt{
				TransactionHash:   *eth.MustHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				TransactionIndex:  *eth.MustQuantity("0x0"),
				BlockHash:         *eth.MustHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
				BlockNumber:       *eth.MustQuantity("0x1"),
				From:              *eth.MustAddress("0x8ba1f109551bD432803012645aac136c12345678"),
				GasUsed:           *eth.MustQuantity("0x5208"),
				CumulativeGasUsed: *eth.MustQuantity("0x5208"),
				LogsBloom:         *eth.MustData256("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
				Type:              eth.MustQuantity("0x0"),
				Logs:              []eth.Log{},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBaseHandler()
			result := handler.ConvertTransactionReceiptToResponse(tt.ethReceipt)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				// 验证数值字段是十进制格式
				if tt.ethReceipt != nil {
					assert.Equal(t, "0", result["transactionIndex"])
					assert.Equal(t, "1", result["blockNumber"])
					assert.Equal(t, "21000", result["gasUsed"])
					assert.Equal(t, "21000", result["cumulativeGasUsed"])
				}
			}
		})
	}
}
