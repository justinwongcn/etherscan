package handler

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockHandler 区块处理器，处理区块相关的HTTP请求
type BlockHandler struct {
	blockService service.BlockServiceInterface
}

// NewBlockHandler 创建区块处理器实例
func NewBlockHandler(blockService service.BlockServiceInterface) *BlockHandler {
	return &BlockHandler{
		blockService: blockService,
	}
}

// GetBlockHeight 获取最新区块高度
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

// GetBlockByNumberOrHash 获取指定区块号或区块哈希的区块信息
func (h *BlockHandler) GetBlockByNumberOrHash(c *gin.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := c.Param("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = "latest"
	}

	// 判断参数类型并调用相应的方法
	var (
		block *eth.Block
		err   error
	)

	switch {
	case len(blockParam) >= 2 && blockParam[:2] == "0x":
		block, err = h.blockService.GetBlockByHash(c.Request.Context(), blockParam)
	case isDecimalNumber(blockParam):
		number, _ := new(big.Int).SetString(blockParam, 10)
		numberOrTag := "0x" + number.Text(16)
		block, err = h.blockService.GetBlockByNumber(c.Request.Context(), numberOrTag)
	default:
		block, err = h.blockService.GetBlockByNumber(c.Request.Context(), blockParam)
	}

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

// isDecimalNumber checks if string is a valid decimal number
func isDecimalNumber(s string) bool {
	_, ok := new(big.Int).SetString(s, 10)
	return ok
}

// GetTransactionCount 获取指定区块的交易数量
func (h *BlockHandler) GetTransactionCount(c *gin.Context) {
	// 从URL路径中获取区块号或哈希
	blockParam := c.Param("number")
	// 如果参数为空，则使用latest
	if blockParam == "" {
		blockParam = "latest"
	}

	// 获取交易数量
	count, err := h.blockService.GetTransactionCount(c.Request.Context(), blockParam)
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
