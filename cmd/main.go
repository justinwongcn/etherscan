package main

import (
	"context"
	"log"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/api/handler"
	"github.com/justinwongcn/etherscan/api/routes"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/etherscan/internal/repository"
)

// setupServer 设置服务器（可测试的函数）
func setupServer(rpcURL string) (*ant.HTTPServer, error) {
	// 创建以太坊客户端
	opts := ethereum.DefaultClientOptions()
	client, err := ethereum.NewClient(context.Background(), rpcURL, opts)
	if err != nil {
		return nil, err
	}

	// 初始化 repository
	blockRepo := repository.NewBlockRepository(client)
	transactionRepo := repository.NewTransactionRepository(client)
	accountRepo := repository.NewAccountRepository(client)

	// 初始化 service
	blockService := service.NewBlockService(blockRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	accountService := service.NewAccountService(accountRepo)

	// 初始化 handler
	blockHandler := handler.NewBlockHandler(blockService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	accountHandler := handler.NewAccountHandler(accountService)

	// 创建一个新的 HTTP 服务器
	server := ant.NewHTTPServer()

	// 注册路由
	routes.RegisterRoutes(server, blockHandler, transactionHandler, accountHandler, blockRepo)

	return server, nil
}

func main() {
	server, err := setupServer("https://rpc.flashbots.net")
	if err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	// 启动服务器
	log.Println("Server is running on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
