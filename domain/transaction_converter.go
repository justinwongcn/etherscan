// Package domain 提供以太坊区块链的领域模型和转换服务
package domain

import (
	"github.com/justinwongcn/etherscan/pkg/parallel"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionConverter 提供以太坊交易数据到领域模型的转换服务
// 该结构体负责将以太坊原生交易数据转换为应用程序使用的领域模型
// 转换过程包括数据类型转换和字段映射等操作
type TransactionConverter struct{}

// NewTransactionConverter 创建一个新的交易转换器实例
// 返回:
//   - *TransactionConverter: 交易转换器实例
func NewTransactionConverter() *TransactionConverter {
	return &TransactionConverter{}
}

// convertBlobVersionedHashes 将eth.Hashes转换为领域模型的Hashes
func (c *TransactionConverter) convertBlobVersionedHashes(ethHashes []eth.Hash) Hashes {
	if ethHashes == nil {
		return nil
	}

	blobHashes := make([]string, len(ethHashes))
	for i, hash := range ethHashes {
		blobHashes[i] = hash.String()
	}
	return blobHashes
}

// convertAccessList 将eth.AccessList转换为领域模型的AccessList
func (c *TransactionConverter) convertAccessList(ethAccessList *eth.AccessList) *AccessList {
	if ethAccessList == nil {
		return nil
	}

	// 使用Process函数并发转换AccessListEntry
	converter := func(entry eth.AccessListEntry) AccessListEntry {
		// 转换StorageKeys
		storageKeys := make([]string, len(entry.StorageKeys))
		for i, key := range entry.StorageKeys {
			storageKeys[i] = key.String()
		}

		return AccessListEntry{
			Address:     entry.Address.String(),
			StorageKeys: storageKeys,
		}
	}

	// 使用Process函数进行并发转换
	result := parallel.Process(*ethAccessList, converter, 30)
	accessList := AccessList(result)
	return &accessList
}

// ConvertToTransaction 将以太坊交易数据转换为领域模型
// 该方法执行以下转换操作:
//  1. 基本字段转换：
//     - 将交易类型、区块号、Gas限制、Nonce等基础字段从eth.Quantity转换为字符串
//     - 处理From地址、Hash、Input数据等固定字段的转换
//  2. EIP相关字段：
//     - EIP-1559: 处理maxFeePerGas和maxPriorityFeePerGas字段，支持新的燃料费用机制
//     - EIP-2930: 处理accessList字段，支持访问列表类型交易
//     - EIP-4844: 处理maxFeePerBlobGas和blobVersionedHashes字段，支持blob类型交易
//  3. 可选字段处理：
//     - 处理可能为空的字段，如blockHash、to地址、raw数据等
//     - 确保所有可选字段在转换时都经过适当的空值检查
//  4. 特殊字段处理：
//     - 处理Parity客户端特有的standardV、creates等扩展字段
//     - 转换chainId、condition等链相关字段
//
// 参数:
//   - ethTx: 以太坊原生交易数据，包含完整的交易字段信息
//     包括基本交易属性、EIP扩展字段和客户端特有字段
//
// 返回:
//   - *Transaction: 转换后的交易领域模型
//     包含规范化的交易数据，所有字段都经过类型转换和空值处理
func (c *TransactionConverter) ConvertToTransaction(ethTx *eth.Transaction) *Transaction {
	if ethTx == nil {
		return nil
	}
	// 基本字段转换
	typ := ethTx.Type.Big().String()

	// 转换eth.Quantity类型为uint64
	blockNumber := ethTx.BlockNumber.Big().String()
	gas := ethTx.Gas.Big().String()
	nonce := ethTx.Nonce.Big().String()
	index := ethTx.Index.Big().String()
	value := ethTx.Value.Big().String()
	v := ethTx.V.Big().String()
	r := ethTx.R.Big().String()
	s := ethTx.S.Big().String()

	// 转换所有字段为uint64，同时处理空值情况
	var yParity, gasPrice, maxFeePerGas, maxPriorityFeePerGas, standardV, chainId, maxFeePerBlobGas string

	if ethTx.YParity != nil {
		yParity = ethTx.YParity.Big().String()
	}
	if ethTx.GasPrice != nil {
		gasPrice = ethTx.GasPrice.Big().String()
	}
	if ethTx.MaxFeePerGas != nil {
		maxFeePerGas = ethTx.MaxFeePerGas.Big().String()
	}
	if ethTx.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = ethTx.MaxPriorityFeePerGas.Big().String()
	}
	if ethTx.StandardV != nil {
		standardV = ethTx.StandardV.Big().String()
	}
	if ethTx.ChainId != nil {
		chainId = ethTx.ChainId.Big().String()
	}
	if ethTx.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = ethTx.MaxFeePerBlobGas.Big().String()
	}

	// 处理Creates字段，确保类型转换正确
	var creates string
	if ethTx.Creates != nil {
		creates = ethTx.Creates.String()
	}

	// 转换eth.Hash、eth.Address和eth.Data类型为string
	var blockHash, to, raw, publicKey *string
	if ethTx.BlockHash != nil {
		bh := ethTx.BlockHash.String()
		blockHash = &bh
	}
	if ethTx.To != nil {
		t := ethTx.To.String()
		to = &t
	}
	if ethTx.Raw != nil {
		r := ethTx.Raw.String()
		raw = &r
	}
	if ethTx.PublicKey != nil {
		pk := ethTx.PublicKey.String()
		publicKey = &pk
	}

	// 转换BlobVersionedHashes
	blobHashes := c.convertBlobVersionedHashes(ethTx.BlobVersionedHashes)

	return &Transaction{
		Type:                 &typ,
		BlockHash:            blockHash,
		BlockNumber:          &blockNumber,
		From:                 ethTx.From.String(),
		Gas:                  gas,
		Hash:                 ethTx.Hash.String(),
		Input:                ethTx.Input.String(),
		Nonce:                nonce,
		To:                   to,
		Index:                &index,
		Value:                value,
		V:                    v,
		R:                    r,
		S:                    s,
		YParity:              &yParity,
		GasPrice:             &gasPrice,
		MaxFeePerGas:         &maxFeePerGas,
		MaxPriorityFeePerGas: &maxPriorityFeePerGas,
		StandardV:            &standardV,
		Raw:                  raw,
		PublicKey:            publicKey,
		ChainId:              &chainId,
		Creates:              &creates,
		Condition:            ethTx.Condition,
		AccessList:           c.convertAccessList(ethTx.AccessList),
		MaxFeePerBlobGas:     &maxFeePerBlobGas,
		BlobVersionedHashes:  blobHashes,
		BlobBundle:           ethTx.BlobBundle,
		AuthorizationList:    ethTx.AuthorizationList,
	}
}
