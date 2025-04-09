package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/application/service"
)

// AccountHandler 账户相关的HTTP处理器
type AccountHandler struct {
	accountService service.AccountServiceInterface
}

// NewAccountHandler 创建账户处理器实例
func NewAccountHandler(accountService service.AccountServiceInterface) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// GetBalance 处理获取账户余额的HTTP请求
// 请求路径: GET /accounts/:address/balance
// 路径参数:
//   - address: 以太坊地址（必填，20字节的十六进制字符串）
//
// 查询参数:
//   - block: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//     默认值: "latest"
//
// 响应格式:
//   - 成功: {"balance": <账户余额字符串>}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 400: 请求参数错误（地址格式无效）
//   - 500: 服务器内部错误
func (h *AccountHandler) GetBalance(c *gin.Context) {
	// 获取路径参数中的地址
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	// 获取区块号参数，默认为"latest"
	block := c.DefaultQuery("block", "latest")

	// 调用服务层获取余额
	balance, err := h.accountService.GetBalance(c.Request.Context(), address, block)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回余额信息
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
