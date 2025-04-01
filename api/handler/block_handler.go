// Package handler 提供HTTP请求处理器，负责处理来自客户端的API请求
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

// BlockHandler 区块处理器，负责处理与以太坊区块相关的HTTP请求
// 该处理器实现了RESTful风格的API接口，提供区块高度、区块信息和交易数量的查询功能
type BlockHandler struct {
	// blockService 是区块服务接口的实现，用于处理具体的业务逻辑
	blockService service.BlockServiceInterface
}

// NewBlockHandler 创建并初始化一个新的区块处理器实例
// 参数:
//   - blockService: 区块服务接口的实现，用于处理区块相关的业务逻辑
//
// 返回:
//   - *BlockHandler: 初始化完成的处理器实例
func NewBlockHandler(blockService service.BlockServiceInterface) *BlockHandler {
	return &BlockHandler{
		blockService: blockService,
	}
}

// GetBlockHeight 处理获取最新区块高度的HTTP请求
// 请求路径: GET /block/height
// 响应格式:
//   - 成功: {"height": <区块高度>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误
func (h *BlockHandler) GetBlockHeight(c *gin.Context) {
	// 获取最新区块高度
	height, err := h.blockService.GetLatestBlockHeight(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回区块高度
	c.JSON(http.StatusOK, gin.H{
		"height": height,
	})
}

// GetBlock 处理获取区块信息的HTTP请求
// 请求路径: GET /block/:number
// 路径参数:
//   - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//
// 响应格式:
//   - 成功: {"block": <区块信息对象>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *BlockHandler) GetBlock(c *gin.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := c.Param("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = ethereum.BlockLatest
	}

	// 获取区块信息
	block, err := h.blockService.GetBlock(c.Request.Context(), blockParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回区块信息
	c.JSON(http.StatusOK, gin.H{
		"block": block,
	})
}

// GetBlockTransactionCount 处理获取区块交易数量的HTTP请求
// 请求路径: GET /block/count/:number/tx
// 路径参数:
//   - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//
// 响应格式:
//   - 成功: {"count": <交易数量>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *BlockHandler) GetBlockTransactionCount(c *gin.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := c.Param("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = ethereum.BlockLatest
	}

	// 获取交易数量
	count, err := h.blockService.GetBlockTransactionCount(c.Request.Context(), blockParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易数量
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// GetTransactionByHash 处理获取交易信息的HTTP请求
// 请求路径: GET /transaction/:hash
// 路径参数:
//   - hash: 交易哈希（32字节的十六进制字符串）
//
// 响应格式:
//   - 成功: {"transaction": <交易信息对象>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *BlockHandler) GetTransactionByHash(c *gin.Context) {
	// 获取交易哈希参数
	txHash := c.Param("hash")
	if txHash == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "transaction hash is required",
		})
		return
	}

	// 获取交易信息
	tx, err := h.blockService.GetTransactionByHash(c.Request.Context(), txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易信息
	c.JSON(http.StatusOK, gin.H{
		"transaction": tx,
	})
}

// GetTransactionByIndex 处理获取指定区块中特定索引位置交易的HTTP请求
// 请求路径: GET /block/tx/:index
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
func (h *BlockHandler) GetTransactionByIndex(c *gin.Context) {
	// 从查询参数中获取区块号或哈希
	blockParam := c.Query("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = ethereum.BlockLatest
	}

	// 获取交易索引参数并转换为uint64
	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid transaction index",
		})
		return
	}

	// 获取交易信息
	tx, err := h.blockService.GetTransactionByIndex(c.Request.Context(), blockParam, index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易信息
	c.JSON(http.StatusOK, gin.H{
		"transaction": tx,
	})
}

// GetTransactionCount 处理获取账户交易数量的HTTP请求
// 请求路径: GET /account/:address/count
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
func (h *BlockHandler) GetTransactionCount(c *gin.Context) {
	// 获取地址参数
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "address is required",
		})
		return
	}

	// 从查询参数中获取区块号或哈希
	blockParam := c.Param("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = ethereum.BlockLatest
	}

	// 获取交易数量
	count, err := h.blockService.GetTransactionCount(c.Request.Context(), address, blockParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易数量
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// SendRawTransactionRequest 定义了发送已签名交易的请求结构
type SendRawTransactionRequest struct {
	SignedTxData string `json:"signedTxData" binding:"required"`
}

// SendRawTransaction 处理发送已签名交易的HTTP请求
// 请求路径: POST /tx/send
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
func (h *BlockHandler) SendRawTransaction(c *gin.Context) {
	var req SendRawTransactionRequest

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 调用服务层发送交易
	txHash, err := h.blockService.SendRawTransaction(c.Request.Context(), req.SignedTxData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回交易哈希
	c.JSON(http.StatusOK, gin.H{"txHash": txHash})
}

// GetTransactionReceipt 处理获取交易收据的HTTP请求
// 请求路径: GET /tx/:hash/receipt
// 路径参数:
//   - hash: 交易哈希（32字节的十六进制字符串）
//
// 响应格式:
//   - 成功: {"receipt": <交易收据信息>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *BlockHandler) GetTransactionReceipt(c *gin.Context) {
	// 获取交易哈希参数
	txHash := c.Param("hash")
	if txHash == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "transaction hash is required",
		})
		return
	}

	// 获取交易收据信息
	receipt, err := h.blockService.GetTransactionReceipt(c.Request.Context(), txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回交易收据信息
	c.JSON(http.StatusOK, gin.H{
		"receipt": receipt,
	})
}
