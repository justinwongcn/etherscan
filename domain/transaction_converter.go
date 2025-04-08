// Package domain 提供以太坊区块链的领域模型和转换服务
package domain

import (
	"github.com/justinwongcn/etherscan/pkg/parallel"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionConverter 提供以太坊交易数据到领域模型的转换服务
// 该结构体负责将以太坊原生交易数据转换为应用程序使用的领域模型
// 转换过程包括数据类型转换和字段映射等操作
type TransactionConverter struct{
	threshold int // 并发处理的阈值
}

// NewTransactionConverter 创建一个新的交易转换器实例
// 返回:
//   - *TransactionConverter: 交易转换器实例
func NewTransactionConverter(threshold int) *TransactionConverter {
	if threshold <= 0 {
		threshold = 50 // 默认阈值
	}
	return &TransactionConverter{
		threshold: threshold,
	}
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
	var yParity, gasPrice, maxFeePerGas, maxPriorityFeePerGas, standardV, chainId, maxFeePerBlobGas, creates *string

	if ethTx.YParity != nil {
		yp := ethTx.YParity.Big().String()
		yParity = &yp
	}
	if ethTx.GasPrice != nil {
		gp := ethTx.GasPrice.Big().String()
		gasPrice = &gp
	}
	if ethTx.MaxFeePerGas != nil {
		mfg := ethTx.MaxFeePerGas.Big().String()
		maxFeePerGas = &mfg
	}
	if ethTx.MaxPriorityFeePerGas != nil {
		mpfg := ethTx.MaxPriorityFeePerGas.Big().String()
		maxPriorityFeePerGas = &mpfg
	}
	if ethTx.StandardV != nil {
		sv := ethTx.StandardV.Big().String()
		standardV = &sv
	}
	if ethTx.ChainId != nil {
		cid := ethTx.ChainId.Big().String()
		chainId = &cid
	}
	if ethTx.MaxFeePerBlobGas != nil {
		mfbg := ethTx.MaxFeePerBlobGas.Big().String()
		maxFeePerBlobGas = &mfbg
	}
	if ethTx.Creates != nil {
		c := ethTx.Creates.String()
		creates = &c
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
		YParity:              yParity,              // 移除多余的取地址符号
		GasPrice:             gasPrice,             // 移除多余的取地址符号
		MaxFeePerGas:         maxFeePerGas,         // 移除多余的取地址符号
		MaxPriorityFeePerGas: maxPriorityFeePerGas, // 移除多余的取地址符号
		StandardV:            standardV,            // 移除多余的取地址符号
		Raw:                  raw,                  // 已经是正确的形式
		PublicKey:            publicKey,            // 已经是正确的形式
		ChainId:              chainId,              // 移除多余的取地址符号
		Creates:              creates,              // 移除多余的取地址符号
		Condition:            ethTx.Condition,      // 保持不变
		AccessList:           c.convertAccessList(ethTx.AccessList),
		MaxFeePerBlobGas:     maxFeePerBlobGas,     // 移除多余的取地址符号
		BlobVersionedHashes:  blobHashes,
		BlobBundle:           ethTx.BlobBundle,
		AuthorizationList:    ethTx.AuthorizationList,
	}
}
