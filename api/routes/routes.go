// Package routes 提供HTTP路由配置，负责将URL请求映射到相应的处理器
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/api/handler"
)

// RegisterRoutes 注册所有HTTP路由
// 该函数配置了所有与以太坊区块查询、交易和账户相关的API路由
// 参数:
//   - r: Gin框架的路由引擎实例
//   - blockHandler: 区块处理器实例，负责处理区块相关的请求逻辑
//   - transactionHandler: 交易处理器实例，负责处理交易相关的请求逻辑
//   - accountHandler: 账户处理器实例，负责处理账户相关的请求逻辑
//
// 路由配置:
//
// 区块相关路由:
//  1. GET /blocks/height/latest
//     获取以太坊网络的最新区块高度
//     响应格式:
//     - 成功: {"height": <区块高度数值>}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 500: 服务器内部错误
//
//  2. GET /blocks/:number
//     获取指定区块的详细信息
//     路径参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest、earliest、pending
//     查询参数:
//     - fullTx: 布尔值，控制是否返回完整的交易对象
//     true: 返回完整的交易对象（默认值）
//     false: 仅返回交易哈希
//     响应格式:
//     - 成功: {"block": {"number": <区块号>, "hash": <区块哈希>, ...}}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（区块号格式无效）
//     - 404: 区块未找到
//     - 500: 服务器内部错误
//
//  3. GET /blocks/:number/transactions/count
//     获取指定区块中的交易数量
//     路径参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest、earliest、pending
//     响应格式:
//     - 成功: {"count": <交易数量>}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（区块号格式无效）
//     - 404: 区块未找到
//     - 500: 服务器内部错误
//
//  4. GET /blocks/:number/transactions/:index
//     获取指定区块中特定索引位置的交易信息
//     路径参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest（最新区块）、earliest（创世区块）、pending（待打包区块）
//     默认值：latest
//     - index: 交易在区块中的索引位置（从0开始的整数）
//     响应格式:
//     - 成功: {"transaction": {"hash": <交易哈希>, "from": <发送方地址>, "to": <接收方地址>, ...}}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（区块号或索引格式无效）
//     - 404: 区块或交易未找到
//     - 500: 服务器内部错误
//
// 交易相关路由:
//  5. GET /transactions/:hash
//     获取指定交易哈希的交易详细信息
//     路径参数:
//     - hash: 交易哈希（32字节的十六进制字符串，0x开头）
//     响应格式:
//     - 成功: {"transaction": {"hash": <交易哈希>, "from": <发送方地址>, "to": <接收方地址>, ...}}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（交易哈希格式无效）
//     - 404: 交易未找到
//     - 500: 服务器内部错误
//
//  6. GET /transactions/:hash/receipt
//     获取指定交易的收据信息
//     路径参数:
//     - hash: 交易哈希（32字节的十六进制字符串，0x开头）
//     响应格式:
//     - 成功: {"receipt": {"transactionHash": <交易哈希>, "blockNumber": <区块号>, "gasUsed": <使用的gas量>, ...}}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（交易哈希格式无效）
//     - 404: 交易收据未找到
//     - 500: 服务器内部错误
//
//  7. POST /transactions
//     发送已签名的交易数据到以太坊网络
//     请求体:
//     - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
//     响应格式:
//     - 成功: {"txHash": <交易哈希>}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（交易数据格式无效）
//     - 500: 服务器内部错误（包括交易广播失败）
//
// 账户相关路由:
//  8. GET /accounts/balance/:address
//     获取指定地址的账户余额
//     路径参数:
//     - address: 以太坊地址（20字节的十六进制字符串，0x开头）
//     响应格式:
//     - 成功: {"balance": <账户余额（单位：wei）>}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（地址格式无效）
//     - 500: 服务器内部错误
//
//  9. GET /accounts/balances/:addresses
//     批量获取多个账户的余额
//     路径参数:
//     - addresses: 以太坊地址列表（多个地址以逗号分隔的20字节十六进制字符串，0x开头）
//     查询参数:
//     - block: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest（最新区块）、earliest（创世区块）、pending（待打包区块）
//     默认值：latest
//     响应格式:
//     - 成功: {"balances": {"address1": "balance1", "address2": "balance2", ...}}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（地址格式无效或地址列表为空）
//     - 500: 服务器内部错误
//
//  10. GET /accounts/:address/transactions/count
//     获取指定地址在特定区块的交易数量
//     路径参数:
//     - address: 以太坊账户地址（20字节的十六进制字符串，0x开头）
//     查询参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest（最新区块）、earliest（创世区块）、pending（待打包区块）
//     默认值：latest
//     响应格式:
//     - 成功: {"count": <交易数量>}
//     - 失败: {"error": <错误信息>}
//     错误码:
//     - 400: 请求参数错误（地址格式无效）
//     - 500: 服务器内部错误

func RegisterRoutes(r *gin.Engine, blockHandler *handler.BlockHandler, transactionHandler *handler.TransactionHandler, accountHandler *handler.AccountHandler) {
	// 区块相关路由
	r.GET("/blocks/height/latest", blockHandler.GetBlockHeight)
	r.GET("/blocks/:number", blockHandler.GetBlock)
	r.GET("/blocks/:number/transactions/count", blockHandler.GetBlockTransactionCount)
	r.GET("/blocks/:number/transactions/:index", transactionHandler.GetTransactionByIndex)

	// 交易相关路由
	r.GET("/transactions/:hash", transactionHandler.GetTransactionByHash)
	r.GET("/transactions/:hash/receipt", transactionHandler.GetTransactionReceipt)
	r.POST("/transactions", transactionHandler.SendRawTransaction)

	// 账户相关路由
	r.GET("/accounts/balance/:address", accountHandler.GetBalance)
	r.GET("/accounts/balances/:addresses", accountHandler.GetBalances)
	r.GET("/accounts/:address/transactions/count", transactionHandler.GetTransactionCount)
}
