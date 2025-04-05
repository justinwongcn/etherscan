// Package routes 提供HTTP路由配置，负责将URL请求映射到相应的处理器
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/api/handler"
)

// RegisterRoutes 注册所有HTTP路由
// 该函数配置了所有与以太坊区块查询相关的API路由
// 参数:
//   - r: Gin框架的路由引擎实例
//   - blockHandler: 区块处理器实例，负责处理具体的请求逻辑
//
// 路由配置:
//
//  1. GET /block/height
//     获取以太坊网络的最新区块高度
//
//  2. GET /block/:number
//     获取指定区块的详细信息
//     参数 number 支持：
//     - 区块号（十进制数字）
//     - 区块哈希（0x开头的十六进制字符串）
//     - 特殊值：latest、earliest、pending
//
//  3. GET /block/count/:number/tx
//     获取指定区块中的交易数量
//     参数 number 支持同上
//
//  4. GET /account/count/:address
//     获取指定地址在特定区块的交易数量
//     路径参数:
//     - address: 以太坊账户地址
//     查询参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest（最新区块）、earliest（创世区块）、pending（待打包区块）
//     默认值：latest
//
//  5. GET /tx/:hash
//     获取指定交易哈希的交易详细信息
//     参数:
//     - hash: 交易哈希（32字节的十六进制字符串）
//
//  6. GET /tx/:hash/receipt
//     获取指定交易的收据信息
//     参数:
//     - hash: 交易哈希（32字节的十六进制字符串）
//     返回:
//     - 交易收据信息，包含交易哈希、区块信息、gas使用情况、合约地址、日志等
//
//  7. GET /block/:number/tx/:index
//     获取指定区块中特定索引位置的交易信息
//     参数:
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest、earliest、pending
//     - index: 交易在区块中的索引位置（从0开始的整数）
//
//  8. POST /tx/send
//     发送已签名的交易数据到以太坊网络
//     请求体:
//     - signedTxData: 已签名的交易数据（十六进制格式，以0x开头）
//     返回:
//     - txHash: 交易哈希（32字节的十六进制字符串）
func RegisterRoutes(r *gin.Engine, blockHandler *handler.BlockHandler, transactionHandler *handler.TransactionHandler) {
	// 区块相关路由
	r.GET("/block/height", blockHandler.GetBlockHeight)
	r.GET("/block/:number", blockHandler.GetBlock)
	r.GET("/block/count/:number", blockHandler.GetBlockTransactionCount)

	// 交易相关路由
	r.GET("/account/count/:address", transactionHandler.GetTransactionCount)
	r.GET("/tx/:hash", transactionHandler.GetTransactionByHash)
	r.GET("/tx/:hash/receipt", transactionHandler.GetTransactionReceipt)
	r.GET("/block/tx/:index", transactionHandler.GetTransactionByIndex)
	r.POST("/tx/send", transactionHandler.SendRawTransaction)
}
