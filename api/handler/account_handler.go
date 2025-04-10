package handler

import (
	"net/http"
	"strings"

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
// 请求路径: GET /accounts/balance/:address
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

// GetBalances 处理批量获取账户余额的HTTP请求
// 请求路径: GET /accounts/balances/:addresses
// 路径参数:
//   - addresses: 以太坊地址列表（必填，多个地址以逗号分隔的20字节十六进制字符串）
//
// 查询参数:
//   - block: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值: "latest"（最新区块）、"earliest"（创世区块）、"pending"（待打包区块）
//     默认值: "latest"
//
// 响应格式:
//   - 成功: {"balances": {"address1": "balance1", "address2": "balance2", ...}}
//   - 失败: {"error": <错误信息>}
//
// 错误码:
//   - 400: 请求参数错误（地址格式无效或地址列表为空）
//   - 500: 服务器内部错误
func (h *AccountHandler) GetBalances(c *gin.Context) {
	// 获取路径参数中的地址列表
	addresses := c.Param("addresses")
	if addresses == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addresses is required"})
		return
	}

	// 分割地址列表
	addrList := strings.Split(addresses, ",")
	// 过滤掉空地址
	validAddrs := make([]string, 0, len(addrList))
	for _, addr := range addrList {
		if addr != "" {
			validAddrs = append(validAddrs, addr)
		}
	}
	if len(validAddrs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addresses list is empty"})
		return
	}
	// 使用过滤后的地址列表
	addrList = validAddrs

	// 获取区块号参数，默认为"latest"
	block := c.DefaultQuery("block", "latest")

	// 调用服务层批量获取余额
	balances, err := h.accountService.GetBalances(c.Request.Context(), addrList, block)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回余额信息
	c.JSON(http.StatusOK, gin.H{"balances": balances})
}
