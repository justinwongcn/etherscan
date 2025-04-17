// Package handler 提供HTTP请求处理器，负责处理来自客户端的API请求
package handler

import (
	"net/http"
	"strconv"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

// TransactionHandler 交易处理器，负责处理与以太坊交易相关的HTTP请求
// 该处理器实现了RESTful风格的API接口，提供交易查询和发送功能
type TransactionHandler struct {
	// transactionService 是交易服务接口的实现，用于处理具体的业务逻辑
	transactionService service.TransactionServiceInterface
}

// NewTransactionHandler 创建并初始化一个新的交易处理器实例
// 参数:
//   - transactionService: 交易服务接口的实现，用于处理交易相关的业务逻辑
//
// 返回:
//   - *TransactionHandler: 初始化完成的处理器实例
func NewTransactionHandler(transactionService service.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetTransactionByHash 处理获取交易信息的HTTP请求
// 请求路径: GET /transactions/:hash
// 路径参数:
//   - hash: 交易哈希（32字节的十六进制字符串）
//
// 响应格式:
//   - 成功: {"transaction": <交易信息对象>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *TransactionHandler) GetTransactionByHash(ctx *ant.Context) {
	// 获取交易哈希参数
	txHash := ctx.Req.PathValue("hash")
	if txHash == "" {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": "transaction hash is required",
		})
		return
	}

	// 获取交易信息
	tx, err := h.transactionService.GetTransactionByHash(ctx.Req.Context(), txHash)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易信息
	_ = ctx.RespJSONOK(H{
		"transaction": tx,
	})
}

// GetTransactionByIndex 处理获取指定区块中特定索引位置交易的HTTP请求
// 请求路径: GET /blocks/:number/transactions/:index
// 路径参数:
//   - index: 交易在区块中的索引位置（从0开始的整数）
//
// 查询参数:
//   - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//
// 响应格式:
//   - 成功: {"transaction": <交易信息对象>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *TransactionHandler) GetTransactionByIndex(ctx *ant.Context) {
	// 获取路径参数
	indexStr := ctx.Req.PathValue("index")
	
	// 获取查询参数
	blockParam, err := ctx.QueryValue("number").String()
	if err!= nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": "invalid block number or hash",
		})
		return
	}
	
	// 如果参数为空，则返回错误信息
	if blockParam == "" {
		_ = ctx.RespJSON(http.StatusBadRequest, H{
			"error": "block number or hash is required",
		})
		return
	}

	// 获取交易索引参数并转换为uint64
	indexStr = ctx.Req.PathValue("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": "invalid transaction index",
		})
		return
	}

	// 获取交易信息
	tx, err := h.transactionService.GetTransactionByIndex(ctx.Req.Context(), blockParam, index)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易信息
	_ = ctx.RespJSONOK(H{
		"transaction": tx,
	})
}

// SendRawTransactionRequest 定义了发送已签名交易的请求结构
type SendRawTransactionRequest struct {
	SignedTxData string `json:"signedTxData"`
}

// SendRawTransaction 处理发送已签名交易的HTTP请求
// 请求路径: POST /transactions
// 请求体:
//   - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
//
// 响应格式:
//   - 成功: {"txHash": <交易哈希>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 400: 请求体格式错误或参数无效
//   - 500: 服务器内部错误
func (h *TransactionHandler) SendRawTransaction(ctx *ant.Context) {
	var req SendRawTransactionRequest

	// 解析请求体
	if err := ctx.BindJSON(&req); err != nil {
		_ = ctx.RespJSON(http.StatusBadRequest, H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 调用服务层发送交易
	txHash, err := h.transactionService.SendRawTransaction(ctx.Req.Context(), req.SignedTxData)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易哈希
	_ = ctx.RespJSONOK(H{
		"txHash": txHash,
	})
}

// GetTransactionCount 处理获取账户交易数量的HTTP请求
// 请求路径: GET /accounts/:address/transactions/count
// 路径参数:
//   - address: 以太坊账户地址
//
// 查询参数:
//   - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//
// 响应格式:
//   - 成功: {"count": <交易数量>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *TransactionHandler) GetTransactionCount(ctx *ant.Context) {
	// 获取地址参数
	address := ctx.Req.PathValue("address")
	if address == "" {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": "address is required",
		})
		return
	}

	blockParam, err := ctx.QueryValue("number").String()
	if err != nil || blockParam == "" {
		blockParam = ethereum.BlockLatest
	}

	// 解析并标准化区块参数格式
	parsedBlockParam, err := ethereum.ParseBlockParameter(blockParam)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 获取交易数量
	count, err := h.transactionService.GetTransactionCount(ctx.Req.Context(), address, parsedBlockParam)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易数量
	_ = ctx.RespJSONOK(H{
		"count": count,
	})
}

// GetTransactionReceipt 处理获取交易收据的HTTP请求
// 请求路径: GET /transactions/:hash/receipt
// 路径参数:
//   - hash: 交易哈希（32字节的十六进制字符串）
//
// 响应格式:
//   - 成功: {"receipt": <交易收据信息>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *TransactionHandler) GetTransactionReceipt(ctx *ant.Context) {
	// 获取交易哈希参数
	txHash := ctx.Req.PathValue("hash")
	if txHash == "" {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": "transaction hash is required",
		})
		return
	}

	// 获取交易收据信息
	receipt, err := h.transactionService.GetTransactionReceipt(ctx.Req.Context(), txHash)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易收据信息
	_ = ctx.RespJSONOK(H{
		"receipt": receipt,
	})
}
