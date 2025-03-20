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

	balance, err := client.GetBalance(ctx, "0xF030AFc7Af7DF4b246f5aC4a3d6A5139e640EC4C", "latest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("账户余额: %d\n", balance)

	//blockNum := fmt.Sprintf("0x%x", 22087700)
	//balance, err = client.GetBalance(ctx, "0xF030AFc7Af7DF4b246f5aC4a3d6A5139e640EC4C", blockNum)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("账户余额: %d\n", balance)

	// 批量获取多个地址的余额
	addresses := []string{
		"0xF030AFc7Af7DF4b246f5aC4a3d6A5139e640EC4C",
		"0xEBC85652D323eCc11f367787E8494555aA8F1D29",
		"0x8240Fe72f039b75B652ab6D656805Dc821b1550f",
	}

	balances, err := client.GetBalances(ctx, addresses, "latest")
	if err != nil {
		panic(err)
	}

	for addr, bal := range balances {
		fmt.Printf("地址 %s 的余额: %d\n", addr, bal)
	}
}
