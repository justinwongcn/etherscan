package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	// "github.com/ethereum/go-ethereum/common/hexutil"
	// "github.com/ethereum/go-ethereum/crypto"
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
	opts := &ethereum.ClientOptions{
		MaxConns:     50,          // 设置最大连接数
		IdleTimeout:  time.Minute, // 设置空闲超时时间
		HealthCheck:  true,        // 启用健康检查
		MaxIdleConns: 5,           // 设置最大空闲连接数
	}
	client, err := ethereum.NewClient(ctx, cfg.Ethereum.NodeURL, opts)
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

	address := "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"

	// 获取最新区块的交易数量
	txCount, err := client.GetTransactionCount(ctx, address, "latest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("地址 %s 在最新区块的交易数量: %d\n", address, txCount)

	// 获取指定区块高度的交易数量
	// blockNum := fmt.Sprintf("0x%x", 22095087)
	// txCountAtBlock, err := client.GetTransactionCount(ctx, address, blockNum)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("地址 %s 在区块 %s 的交易数量: %d\n", address, blockNum, txCountAtBlock)

	// 获取指定区块号的交易数量
	// blockNum := fmt.Sprintf("0x%x", 22115865)
	// txCountByNumber, err := client.GetBlockTransactionCountByNumber(ctx, blockNum)
	// if err != nil {
	//     panic(err)
	// }
	// fmt.Printf("区块 %s 的交易数量: %d\n", blockNum, txCountByNumber)

	// 获取指定区块哈希的交易数量
	// blockHash := "0xc6fe56ed79afaf8330ec42e8b725bdbae4bea29a043bd34c469e50de51a83b3d"
	// txCountByHash, err := client.GetBlockTransactionCountByHash(ctx, blockHash)
	// if err != nil {
	//     panic(err)
	// }
	// fmt.Printf("区块哈希 %s 的交易数量: %d\n", blockHash, txCountByHash)
}
