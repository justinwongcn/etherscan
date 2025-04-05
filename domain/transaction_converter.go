// Package domain 提供以太坊区块链的领域模型和转换服务
package domain

import (
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

// ConvertToTransaction 将以太坊交易数据转换为领域模型
// 该方法执行以下转换操作:
//  1. 基本字段转换：将原生数据类型转换为领域模型对应的类型
//  2. EIP相关字段：处理EIP-1559、EIP-2930、EIP-4844等协议引入的字段
//  3. Parity特有字段：处理Parity客户端特有的扩展字段
//
// 参数:
//   - ethTx: 以太坊原生交易数据，包含交易的所有原始信息
//
// 返回:
//   - *Transaction: 转换后的交易领域模型，包含所有必要的交易信息
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

	return &Transaction{
		Type:                 &typ,
		BlockHash:            ethTx.BlockHash,
		BlockNumber:          &blockNumber,
		From:                 ethTx.From,
		Gas:                  gas,
		Hash:                 ethTx.Hash,
		Input:                ethTx.Input,
		Nonce:                nonce,
		To:                   ethTx.To,
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
		Raw:                  ethTx.Raw,
		PublicKey:            ethTx.PublicKey,
		ChainId:              &chainId,
		Creates:              ethTx.Creates,
		Condition:            ethTx.Condition,
		AccessList:           ethTx.AccessList,
		MaxFeePerBlobGas:     &maxFeePerBlobGas,
		BlobVersionedHashes:  ethTx.BlobVersionedHashes,
		BlobBundle:           ethTx.BlobBundle,
		AuthorizationList:    ethTx.AuthorizationList,
	}
}
