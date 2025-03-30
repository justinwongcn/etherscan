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
//  4. GET /account/:address/tx/count/:number
//     获取指定地址在特定区块的交易数量
//     参数:
//     - address: 以太坊账户地址
//     - number: 区块号（十进制数字）或区块哈希（0x开头的十六进制字符串）
//     支持的特殊值：latest、earliest、pending
func RegisterRoutes(r *gin.Engine, blockHandler *handler.BlockHandler) {
	// 区块相关路由
	r.GET("/block/height", blockHandler.GetBlockHeight)
	r.GET("/block/:number", blockHandler.GetBlock)
	r.GET("/block/count/:number", blockHandler.GetBlockTransactionCount)
	r.GET("/account/count/:address", blockHandler.GetTransactionCount)
}
