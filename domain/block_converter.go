// Package domain 提供以太坊区块链的领域模型和转换服务
package domain

import (
	"sync"

	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockConverter 提供以太坊区块数据到领域模型的转换服务
// 该结构体负责将以太坊原生区块数据转换为应用程序使用的领域模型
// 转换过程包括数据类型转换、字段映射以及并发处理等操作
type BlockConverter struct{}

// NewBlockConverter 创建一个新的区块转换器实例
// 返回:
//   - *BlockConverter: 区块转换器实例
func NewBlockConverter() *BlockConverter {
	return &BlockConverter{}
}

// ConvertToBlock 将以太坊区块数据转换为领域模型
// 该方法执行以下转换操作:
//  1. 基本字段转换：将原生数据类型转换为领域模型对应的类型
//  2. 并发处理：使用goroutine并发处理SealFields数据
//  3. 数据组装：将转换后的数据组装成Block领域模型
//
// 参数:
//   - ethBlock: 以太坊原生区块数据，包含区块的所有原始信息
//
// 返回:
//   - *Block: 转换后的区块领域模型，包含所有必要的区块信息
func (c *BlockConverter) ConvertToBlock(ethBlock *eth.Block) *Block {
	// 转换基本字段
	// 将uint64和string类型的字段进行基础转换
	number := ethBlock.Number.UInt64()
	hash := ethBlock.Hash.String()
	baseFeePerGas := ethBlock.BaseFeePerGas.UInt64()
	withdrawalsRoot := ethBlock.WithdrawalsRoot.String()
	excessBlobGas := ethBlock.ExcessBlobGas.UInt64()
	blobGasUsed := ethBlock.BlobGasUsed.UInt64()
	nonce := ethBlock.Nonce.String()
	mixHash := ethBlock.MixHash.String()

	// 并发处理SealFields
	// 使用goroutine和互斥锁实现并发安全的字符串转换
	var sealFields []string
	if ethBlock.SealFields != nil {
		lenSealFields := len(*ethBlock.SealFields)
		sealFields = make([]string, lenSealFields)

		var wg sync.WaitGroup
		wg.Add(lenSealFields)

		var mu sync.Mutex
		for i, field := range *ethBlock.SealFields {
			go func(index int, f eth.Data) {
				defer wg.Done()

				str := f.String()
				mu.Lock()
				sealFields[index] = str
				mu.Unlock()
			}(i, field)
		}

		wg.Wait()
	}

	// 构建并返回领域模型
	// 将所有转换后的字段组装成Block结构体
	return &Block{
		Number:                &number,
		Hash:                  &hash,
		ParentHash:            ethBlock.ParentHash.String(),
		SHA3Uncles:            ethBlock.SHA3Uncles.String(),
		LogsBloom:             ethBlock.LogsBloom.String(),
		TransactionsRoot:      ethBlock.TransactionsRoot.String(),
		StateRoot:             ethBlock.StateRoot.String(),
		ReceiptsRoot:          ethBlock.ReceiptsRoot.String(),
		Miner:                 ethBlock.Miner.String(),
		Author:                ethBlock.Author.String(),
		Difficulty:            ethBlock.Difficulty.UInt64(),
		TotalDifficulty:       ethBlock.TotalDifficulty.Big().Uint64(),
		ExtraData:             ethBlock.ExtraData.String(),
		Size:                  ethBlock.Size.UInt64(),
		GasLimit:              ethBlock.GasLimit.UInt64(),
		GasUsed:               ethBlock.GasUsed.UInt64(),
		Timestamp:             ethBlock.Timestamp.UInt64(),
		Transactions:          ethBlock.Transactions,
		Uncles:                ethBlock.Uncles,
		BaseFeePerGas:         &baseFeePerGas,
		WithdrawalsRoot:       &withdrawalsRoot,
		Withdrawals:           ethBlock.Withdrawals,
		ParentBeaconBlockRoot: ethBlock.ParentBeaconBlockRoot,
		ExcessBlobGas:         &excessBlobGas,
		BlobGasUsed:           &blobGasUsed,
		Nonce:                 &nonce,
		MixHash:               &mixHash,
		Step:                  ethBlock.Step,
		Signature:             ethBlock.Signature,
		SealFields:            &sealFields,
	}
}
