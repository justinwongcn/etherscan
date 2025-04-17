// Package handler 提供HTTP请求处理器，负责处理来自客户端的API请求
package handler

import (
	"net/http"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/application/service"
)

// BlockHandler 区块处理器，负责处理与以太坊区块相关的HTTP请求
// 该处理器实现了以下API接口：
//   - GET /blocks/height/latest: 获取最新区块高度
//   - GET /blocks/:number: 获取指定区块的详细信息
//   - GET /blocks/:number/transactions/count: 获取指定区块的交易数量
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

// GetLatestBlockHeight 处理获取最新区块高度的HTTP请求
// 响应格式:
//   - 成功: {"height": <区块高度>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误
func (h *BlockHandler) GetLatestBlockHeight(ctx *ant.Context) {
	// 获取最新区块高度
	height, err := h.blockService.GetLatestBlockHeight(ctx.Req.Context())
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回区块高度
	_ = ctx.RespJSONOK(H{
		"height": height,
	})
}

// GetBlock 处理获取区块信息的HTTP请求
// 请求路径: GET /blocks/:number
// 路径参数:
//   - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//
// 查询参数:
//   - fullTx: 布尔值，控制是否返回完整的交易信息，默认为true
//
// 响应格式:
//   - 成功: {"block": <区块信息对象>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 500: 服务器内部错误（包括参数格式错误）
func (h *BlockHandler) GetBlock(ctx *ant.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := ctx.Req.PathValue("number")

	// 获取fullTx查询参数，默认为true
	fullTx := true
	fullTxVal := ctx.QueryValue("fullTx")
	if fullTxStr, err := fullTxVal.String(); err == nil && fullTxStr == "false" {
		fullTx = false
	}

	// 获取区块信息
	block, err := h.blockService.GetBlock(ctx.Req.Context(), blockParam, fullTx)
	if err != nil {
		_ = ctx.RespJSON(http.StatusInternalServerError, H{
			"error": err.Error(),
		})
		return
	}

	// 返回区块信息
	_ = ctx.RespJSONOK(H{
		"block": block,
	})
}

// GetBlockTransactionCount 处理获取区块交易数量的HTTP请求
// 请求路径: GET /blocks/:number/transactions/count
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
func (h *BlockHandler) GetBlockTransactionCount(ctx *ant.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := ctx.Req.PathValue("number")

	// 获取交易数量
	count, err := h.blockService.GetBlockTransactionCount(ctx.Req.Context(), blockParam)
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
