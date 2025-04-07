// Package domain 提供以太坊区块链的领域模型和转换服务
package domain

import (
	"github.com/justinwongcn/etherscan/pkg/parallel"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockConverter 提供以太坊区块数据到领域模型的转换服务
// 该结构体负责将以太坊原生区块数据转换为应用程序使用的领域模型
// 转换过程包括数据类型转换、字段映射以及并发处理等操作
type BlockConverter struct {
	threshold int // 并发处理的阈值
}

// NewBlockConverter 创建一个新的区块转换器实例
// 参数:
//   - threshold: 并发处理的阈值，当数据量小于此值时使用普通循环处理
//
// 返回:
//   - *BlockConverter: 区块转换器实例
func NewBlockConverter(threshold int) *BlockConverter {
	if threshold <= 0 {
		threshold = 50 // 默认阈值
	}
	return &BlockConverter{
		threshold: threshold,
	}
}

// ConvertToBlock 将以太坊区块数据转换为领域模型
// 该方法执行以下转换操作:
//  1. 基本字段转换：将原生数据类型转换为领域模型对应的类型
//  2. 并发处理：使用goroutine并发处理SealFields数据
//  3. 数据组装：将转换后的数据组装成Block领域模型
//
// 参数:
//   - ethBlock: 以太坊原生区块数据，包含区块的所有原始信息
//   - fullTx: 是否转换完整的交易信息
//
// 返回:
//   - *Block: 转换后的区块领域模型，包含所有必要的区块信息
func (c *BlockConverter) ConvertToBlock(ethBlock *eth.Block, fullTx bool) *Block {
	if ethBlock == nil {
		return nil
	}
	// 转换基本字段
	// 将string和string类型的字段进行基础转换
	number := ethBlock.Number.Big().String()
	hash := ethBlock.Hash.String()

	// 处理可选字段，添加nil检查
	var baseFeePerGas string
	if ethBlock.BaseFeePerGas != nil {
		baseFeePerGas = ethBlock.BaseFeePerGas.Big().String()
	}

	var withdrawalsRoot string
	if ethBlock.WithdrawalsRoot != nil {
		withdrawalsRoot = ethBlock.WithdrawalsRoot.String()
	}

	var excessBlobGas string
	if ethBlock.ExcessBlobGas != nil {
		excessBlobGas = ethBlock.ExcessBlobGas.Big().String()
	}

	var blobGasUsed string
	if ethBlock.BlobGasUsed != nil {
		blobGasUsed = ethBlock.BlobGasUsed.Big().String()
	}

	var nonce string
	if ethBlock.Nonce != nil {
		nonce = ethBlock.Nonce.String()
	}

	var mixHash string
	if ethBlock.MixHash != nil {
		mixHash = ethBlock.MixHash.String()
	}

	var parentBeaconBlockRoot string
	if ethBlock.ParentBeaconBlockRoot != nil {
		parentBeaconBlockRoot = ethBlock.ParentBeaconBlockRoot.String()
	}

	// 处理SealFields
	sealFields := c.convertSealFields(ethBlock.SealFields)

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
		Difficulty:            ethBlock.Difficulty.Big().String(),
		TotalDifficulty:       ethBlock.TotalDifficulty.Big().String(),
		ExtraData:             ethBlock.ExtraData.String(),
		Size:                  ethBlock.Size.Big().String(),
		GasLimit:              ethBlock.GasLimit.Big().String(),
		GasUsed:               ethBlock.GasUsed.Big().String(),
		Timestamp:             ethBlock.Timestamp.Big().String(),
		Transactions:          c.convertTransactions(ethBlock.Transactions, fullTx),
		Uncles:                c.convertUncles(ethBlock.Uncles),
		BaseFeePerGas:         &baseFeePerGas,
		WithdrawalsRoot:       &withdrawalsRoot,
		Withdrawals:           c.convertWithdrawals(ethBlock.Withdrawals),
		ParentBeaconBlockRoot: &parentBeaconBlockRoot,
		ExcessBlobGas:         &excessBlobGas,
		BlobGasUsed:           &blobGasUsed,
		Nonce:                 &nonce,
		MixHash:               &mixHash,
		Step:                  ethBlock.Step,
		Signature:             ethBlock.Signature,
		SealFields:            &sealFields,
	}
}

// convertTransactions 将eth.TxOrHash切片转换为domain.TxOrHash切片或交易哈希切片
func (c *BlockConverter) convertTransactions(ethTxs []eth.TxOrHash, fullTx bool) any {
	if ethTxs == nil {
		return nil
	}

	if !fullTx {
		// 只返回交易哈希
		return parallel.Process(ethTxs,
			func(tx eth.TxOrHash) string { return tx.Hash.String() },
			c.threshold)
	}

	// 返回完整交易信息
	txConverter := NewTransactionConverter()
	return parallel.Process(ethTxs,
		func(tx eth.TxOrHash) TxOrHash {
			return TxOrHash{
				Transaction: *txConverter.ConvertToTransaction(&tx.Transaction),
				Populated:   true,
			}
		},
		c.threshold)
}

// convertUncles 将eth.Hash切片转换为字符串切片
func (c *BlockConverter) convertUncles(uncles []eth.Hash) []string {
	return parallel.Process(uncles,
		func(uncle eth.Hash) string { return uncle.String() },
		c.threshold)
}

// convertWithdrawals 将eth.Withdrawal切片转换为domain.Withdrawal切片
func (c *BlockConverter) convertWithdrawals(withdrawals []eth.Withdrawal) []Withdrawal {
	return parallel.Process(withdrawals,
		func(w eth.Withdrawal) Withdrawal {
			return Withdrawal{
				Index:          w.Index.Big().String(),
				ValidatorIndex: w.ValidatorIndex.Big().String(),
				Address:        w.Address.String(),
				Amount:         w.Amount.Big().String(),
			}
		},
		c.threshold)
}

// convertSealFields 将eth.SealFields切片转换为字符串切片
func (c *BlockConverter) convertSealFields(sealFields *[]eth.Data) []string {
	if sealFields == nil {
		return nil
	}

	return parallel.Process(*sealFields,
		func(field eth.Data) string { return field.String() },
		c.threshold)
}
