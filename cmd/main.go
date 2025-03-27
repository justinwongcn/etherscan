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

	block, err := client.GetBlockByNumber(ctx, "latest", true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n最新区块信息:\n")
	fmt.Printf("  区块号: %d\n", block.Number)
	fmt.Printf("  区块哈希: %s\n", block.Hash)
	fmt.Printf("  父区块哈希: %s\n", block.ParentHash)
	fmt.Printf("  时间戳: %d\n", block.Timestamp)
	fmt.Printf("  交易数量: %d\n", len(block.Transactions))
	fmt.Printf("  Gas 限制: %d\n", block.GasLimit)
	fmt.Printf("  Gas 使用量: %d\n", block.GasUsed)

	// 通过区块哈希获取区块信息
	// blockHash := "0x1e9f9b71ea85e1037dd14438714b74fe8c36c93b6b334336aa0708ffbb4c206c"
	// blockByHash, err := client.GetBlockByHash(ctx, blockHash, true)
	// if err != nil {
	// 	fmt.Printf("获取区块哈希 %s 的信息时出错: %v\n", blockHash, err)
	// 	fmt.Println("注意: 某些节点可能不提供历史区块的完整信息，需要使用归档节点")
	// } else {
	// 	fmt.Printf("\n指定哈希的区块信息:\n")
	// 	fmt.Printf("  区块号: %d\n", blockByHash.Number)
	// 	fmt.Printf("  区块哈希: %s\n", blockByHash.Hash)
	// 	fmt.Printf("  父区块哈希: %s\n", blockByHash.ParentHash)
	// 	fmt.Printf("  时间戳: %d\n", blockByHash.Timestamp)
	// 	fmt.Printf("  交易数量: %d\n", len(blockByHash.Transactions))
	// 	fmt.Printf("  Gas 限制: %d\n", blockByHash.GasLimit)
	// 	fmt.Printf("  Gas 使用量: %d\n", blockByHash.GasUsed)
	// }

	// 获取指定交易哈希的交易详细信息
	// txHash := "0x2116f2bd64328306807f4551020b060e52fa2a9fdd18095835da39e89a13dca4"
	// tx, err := client.GetTransactionByHash(ctx, txHash)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("\n交易详细信息:\n")
	// fmt.Printf("  交易哈希: %s\n", tx.Hash)
	// fmt.Printf("  区块号: %d\n", tx.BlockNumber)
	// fmt.Printf("  区块哈希: %s\n", tx.BlockHash)
	// fmt.Printf("  发送方: %s\n", tx.From)
	// fmt.Printf("  接收方: %s\n", tx.To)
	// fmt.Printf("  交易值: %d\n", tx.Value)
	// fmt.Printf("  Gas 限制: %d\n", tx.Gas)
	// fmt.Printf("  Gas 价格: %d\n", tx.GasPrice)

	// 通过区块哈希和交易索引获取交易信息
	// blockHash := "0xc6fe56ed79afaf8330ec42e8b725bdbae4bea29a043bd34c469e50de51a83b3d"
	// txIndex := uint64(6)
	// txByIndex, err := client.GetTransactionByBlockHashAndIndex(ctx, blockHash, txIndex)
	// if err != nil {
	// 	fmt.Printf("获取区块哈希 %s 中索引为 %d 的交易时出错: %v\n", blockHash, txIndex, err)
	// 	fmt.Println("注意: 某些节点可能不提供历史区块的完整信息，需要使用归档节点")
	// } else {
	// 	fmt.Printf("\n通过区块哈希和索引获取的交易信息:\n")
	// 	fmt.Printf("  交易哈希: %s\n", txByIndex.Hash)
	// 	fmt.Printf("  区块号: %d\n", txByIndex.BlockNumber)
	// 	fmt.Printf("  区块哈希: %s\n", txByIndex.BlockHash)
	// 	fmt.Printf("  发送方: %s\n", txByIndex.From)
	// 	fmt.Printf("  接收方: %s\n", txByIndex.To)
	// 	fmt.Printf("  交易值: %d\n", txByIndex.Value)
	// 	fmt.Printf("  Gas 限制: %d\n", txByIndex.Gas)
	// 	fmt.Printf("  Gas 价格: %d\n", txByIndex.GasPrice)
	// }

	// 通过区块号和交易索引获取交易信息
	// blockNum := fmt.Sprintf("0x%x", 22123841)
	// txIndex := uint64(112)
	// txByNumberAndIndex, err := client.GetTransactionByBlockNumberAndIndex(ctx, blockNum, txIndex)
	// if err != nil {
	//     fmt.Printf("获取区块 %s 中索引为 %d 的交易时出错: %v\n", blockNum, txIndex, err)
	//     fmt.Println("注意: 某些节点可能不提供历史区块的完整信息，需要使用归档节点")
	// } else {
	//     fmt.Printf("\n通过区块号和索引获取的交易信息:\n")
	//     fmt.Printf("  交易哈希: %s\n", txByNumberAndIndex.Hash)
	//     fmt.Printf("  区块号: %d\n", txByNumberAndIndex.BlockNumber)
	//     fmt.Printf("  区块哈希: %s\n", txByNumberAndIndex.BlockHash)
	//     fmt.Printf("  发送方: %s\n", txByNumberAndIndex.From)
	//     fmt.Printf("  接收方: %s\n", txByNumberAndIndex.To)
	//     fmt.Printf("  交易值: %d\n", txByNumberAndIndex.Value)
	//     fmt.Printf("  Gas 限制: %d\n", txByNumberAndIndex.Gas)
	//     fmt.Printf("  Gas 价格: %d\n", txByNumberAndIndex.GasPrice)
	// }

	// 获取交易回执信息
	receipt, err := client.GetTransactionReceipt(ctx, "0xfd225fcad404dbaf401c9c19de219a4867e37813ab3fd1e7adb75cf878031629")
	if err != nil {
		panic(err)
	}

	if receipt == nil {
		fmt.Println("交易回执不存在，可能交易尚未被打包")
	} else {
		fmt.Printf("\n交易回执信息:\n")
		fmt.Printf("  交易哈希: %s\n", receipt.TransactionHash.Hash())
		fmt.Printf("  区块号: %d\n", receipt.BlockNumber.Int64())
		fmt.Printf("  区块哈希: %s\n", receipt.BlockHash)
		fmt.Printf("  交易索引: %d\n", receipt.TransactionIndex.Int64())
		fmt.Printf("  合约地址: %s\n", receipt.ContractAddress)
		fmt.Printf("  Gas 使用量: %d\n", receipt.GasUsed.Int64())
		fmt.Printf("  状态: %v\n", receipt.Status)
	}
	// TODO
}
