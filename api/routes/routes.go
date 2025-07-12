// Package routes 提供HTTP路由配置，负责将URL请求映射到相应的处理器
package routes

import (
	"fmt"
	"time"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/api/handler"
	"github.com/justinwongcn/etherscan/domain/repository"
	domainService "github.com/justinwongcn/etherscan/domain/service"
)

// RegisterRoutes 注册所有HTTP路由
// 该函数配置了所有与以太坊区块查询、账户和交易相关的API路由
// 参数:
//   - server: HTTP服务器实例
//   - blockHandler: 区块处理器实例，负责处理区块相关的请求逻辑
//   - accountHandler: 账户处理器实例，负责处理账户相关的请求逻辑
//   - transactionHandler: 交易处理器实例，负责处理交易相关的请求逻辑
//   - blockRepo: 区块仓储实例，用于创建缓存服务
func RegisterRoutes(server *ant.HTTPServer, blockHandler *handler.BlockHandler, transactionHandler *handler.TransactionHandler, accountHandler *handler.AccountHandler, blockRepo repository.BlockRepository) {
	// 区块相关路由
	server.Handle("GET /blocks/height/latest", blockHandler.GetLatestBlockHeight)
	server.Handle("GET /blocks/{number}", blockHandler.GetBlock)
	server.Handle("GET /blocks/{number}/transactions/count", blockHandler.GetBlockTransactionCount)

	// 缓存版本的区块高度查询路由
	server.Handle("GET /cached/blocks/height/latest", createCachedHeightHandler(blockRepo))

	// 缓存版本的区块查询路由
	server.Handle("GET /cached/blocks/{number}", createCachedBlockHandler(blockRepo))
	server.Handle("GET /cached/blocks/hash/{hash}", createCachedBlockByHashHandler(blockRepo))

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

// createCachedHeightHandler 创建一个带缓存功能的区块高度查询处理器
// 该函数会创建一个缓存的区块高度服务，并返回相应的HTTP处理函数
// 参数:
//   - blockRepo: 区块仓储实例，用于创建区块高度服务
//
// 返回:
//   - func(*ant.Context): HTTP处理函数，处理缓存版本的区块高度查询
func createCachedHeightHandler(blockRepo repository.BlockRepository) func(*ant.Context) {
	// 创建基础的区块高度服务
	baseHeightService := domainService.NewBlockHeightService(blockRepo)

	// 创建带缓存的区块高度服务，使用5秒过期时间
	cachedHeightService, err := domainService.NewCachedBlockHeightService(
		baseHeightService,
		5*time.Second, // 5秒缓存过期时间
	)
	if err != nil {
		// 如果创建缓存服务失败，返回错误处理函数
		return func(ctx *ant.Context) {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("创建缓存服务失败: %v", err),
			})
		}
	}

	// 返回处理函数
	return func(ctx *ant.Context) {
		// 使用缓存的区块高度服务获取最新高度
		height, err := cachedHeightService.GetLatestBlockHeight(ctx.Req.Context())
		if err != nil {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("获取区块高度失败: %v", err),
			})
			return
		}

		// 将uint64转换为字符串格式，保持与原有API的兼容性
		heightStr := fmt.Sprintf("%d", height)

		// 返回区块高度
		_ = ctx.RespJSONOK(map[string]string{
			"height": heightStr,
		})
	}
}

// createCachedBlockHandler 创建一个带缓存功能的区块查询处理器
// 该函数会创建一个缓存的区块服务，并返回相应的HTTP处理函数
// 参数:
//   - blockRepo: 区块仓储实例，用于创建区块服务
//
// 返回:
//   - func(*ant.Context): HTTP处理函数，处理缓存版本的区块查询
func createCachedBlockHandler(blockRepo repository.BlockRepository) func(*ant.Context) {
	// 创建带缓存的区块服务
	cachedBlockService, err := domainService.NewCachedBlockService(blockRepo, nil)
	if err != nil {
		// 如果创建缓存服务失败，返回错误处理函数
		return func(ctx *ant.Context) {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("创建缓存区块服务失败: %v", err),
			})
		}
	}

	// 返回处理函数
	return func(ctx *ant.Context) {
		// 获取区块号参数
		blockNumberValue := ctx.PathValue("number")
		blockNumber, err := blockNumberValue.String()
		if err != nil || blockNumber == "" {
			_ = ctx.RespJSON(400, map[string]string{
				"error": "区块号参数无效",
			})
			return
		}

		// 获取fullTx参数，默认为false
		fullTxStr := ctx.Req.URL.Query().Get("fullTx")
		fullTx := fullTxStr == "true"

		// 使用缓存的区块服务获取区块
		block, err := cachedBlockService.GetBlockByNumber(ctx.Req.Context(), blockNumber, fullTx)
		if err != nil {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("获取区块失败: %v", err),
			})
			return
		}

		// 返回区块信息
		_ = ctx.RespJSONOK(block)
	}
}

// createCachedBlockByHashHandler 创建一个带缓存功能的按哈希查询区块处理器
// 该函数会创建一个缓存的区块服务，并返回相应的HTTP处理函数
// 参数:
//   - blockRepo: 区块仓储实例，用于创建区块服务
//
// 返回:
//   - func(*ant.Context): HTTP处理函数，处理缓存版本的按哈希查询区块
func createCachedBlockByHashHandler(blockRepo repository.BlockRepository) func(*ant.Context) {
	// 创建带缓存的区块服务
	cachedBlockService, err := domainService.NewCachedBlockService(blockRepo, nil)
	if err != nil {
		// 如果创建缓存服务失败，返回错误处理函数
		return func(ctx *ant.Context) {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("创建缓存区块服务失败: %v", err),
			})
		}
	}

	// 返回处理函数
	return func(ctx *ant.Context) {
		// 获取区块哈希参数
		blockHashValue := ctx.PathValue("hash")
		blockHash, err := blockHashValue.String()
		if err != nil || blockHash == "" {
			_ = ctx.RespJSON(400, map[string]string{
				"error": "区块哈希参数无效",
			})
			return
		}

		// 获取fullTx参数，默认为false
		fullTxStr := ctx.Req.URL.Query().Get("fullTx")
		fullTx := fullTxStr == "true"

		// 使用缓存的区块服务获取区块
		block, err := cachedBlockService.GetBlockByHash(ctx.Req.Context(), blockHash, fullTx)
		if err != nil {
			_ = ctx.RespJSON(500, map[string]string{
				"error": fmt.Sprintf("获取区块失败: %v", err),
			})
			return
		}

		// 返回区块信息
		_ = ctx.RespJSONOK(block)
	}
}
