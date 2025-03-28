package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/api/handler"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

func main() {
	// 创建以太坊客户端
	opts := ethereum.DefaultClientOptions()
	client, err := ethereum.NewClient(context.Background(), "wss://ethereum.callstaticrpc.com", opts)
	if err != nil {
		log.Fatalf("Failed to create ethereum client: %v", err)
	}

	// 初始化服务层
	blockService := service.NewBlockService(client)

	// 初始化处理器
	blockHandler := handler.NewBlockHandler(blockService)

	// 设置Gin路由
	r := gin.Default()

	// 注册路由
	r.GET("/block/height", blockHandler.GetBlockHeight)
	r.GET("/block/:number", blockHandler.GetBlockByNumberOrHash)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
