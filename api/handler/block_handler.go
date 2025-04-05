// Package handler 提供HTTP请求处理器，负责处理来自客户端的API请求
package handler

import (
	"net/http"

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
