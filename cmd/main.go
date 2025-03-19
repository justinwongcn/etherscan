package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/justinwongcn/etherscan/config"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(filepath.Join("config", "config.yaml"))
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 创建上下文
	ctx := context.Background()

	// 创建客户端实例
	client, err := ethereum.NewClient(ctx, cfg.Ethereum.NodeURL)
	if err != nil {
		panic(err)
	}

	// 获取最新区块高度
	blockNumber, err := client.GetLatestBlockNumber(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("最新区块高度: %d\n", blockNumber)

	balance, err := client.GetBalance(ctx, "0x539C010dBe6AC6E34C8937CB0A35130eDDDb2693", "latest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("账户余额: %d\n", balance)
}
