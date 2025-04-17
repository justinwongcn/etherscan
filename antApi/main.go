package main

import (
	"context"
	"log"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/antApi/handler"
	"github.com/justinwongcn/etherscan/antApi/routes"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

func main() {
	// 创建以太坊客户端
	opts := ethereum.DefaultClientOptions()
	client, err := ethereum.NewClient(context.Background(), "https://rpc.flashbots.net", opts)
	if err != nil {
		log.Fatalf("Failed to create ethereum client: %v", err)
	}

	// 初始化 service
	blockService := service.NewBlockService(client)
	trancactionService := service.NewTransactionService(client)
	accountService := service.NewAccountService(client)

	// 初始化 handler
	blockHandler := handler.NewBlockHandler(blockService)
	trancactionHandler := handler.NewTransactionHandler(trancactionService)
	accountHandler := handler.NewAccountHandler(accountService)

	// 创建一个新的 HTTP 服务器
	server := ant.NewHTTPServer()

	// 注册路由
	routes.RegisterRoutes(server, blockHandler, trancactionHandler ,accountHandler)

	// 启动服务器
	log.Println("Server is running on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
