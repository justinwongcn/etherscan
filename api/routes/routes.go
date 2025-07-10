// Package routes 提供HTTP路由配置，负责将URL请求映射到相应的处理器
package routes

import (
	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/api/handler"
)

// RegisterRoutes 注册所有HTTP路由
// 该函数配置了所有与以太坊区块查询、账户和交易相关的API路由
// 参数:
//   - server: HTTP服务器实例
//   - blockHandler: 区块处理器实例，负责处理区块相关的请求逻辑
//   - accountHandler: 账户处理器实例，负责处理账户相关的请求逻辑
//   - transactionHandler: 交易处理器实例，负责处理交易相关的请求逻辑
func RegisterRoutes(server *ant.HTTPServer, blockHandler *handler.BlockHandler, transactionHandler *handler.TransactionHandler, accountHandler *handler.AccountHandler) {
	// 区块相关路由
	server.Handle("GET /blocks/height/latest", blockHandler.GetLatestBlockHeight)
	server.Handle("GET /blocks/{number}", blockHandler.GetBlock)
	server.Handle("GET /blocks/{number}/transactions/count", blockHandler.GetBlockTransactionCount)

	// 账户相关路由
	server.Handle("GET /accounts/balance/{address}", accountHandler.GetBalance)
	server.Handle("GET /accounts/balances/{addresses}", accountHandler.GetBalances)

	// 交易相关路由
	server.Handle("GET /transactions/{hash}", transactionHandler.GetTransactionByHash)
	server.Handle("GET /blocks/transactions/{index}", transactionHandler.GetTransactionByIndex)
	server.Handle("POST /transactions", transactionHandler.SendRawTransaction)
	server.Handle("GET /accounts/{address}/transactions/count", transactionHandler.GetTransactionCount)
	server.Handle("GET /transactions/{hash}/receipt", transactionHandler.GetTransactionReceipt)
}
