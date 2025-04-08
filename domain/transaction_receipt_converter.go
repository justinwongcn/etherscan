// Package domain 提供以太坊区块链领域模型的定义
package domain

import (
	"github.com/justinwongcn/etherscan/pkg/parallel"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// TransactionReceiptConverter 用于在eth.TransactionReceipt和domain.TransactionReceipt之间进行转换
type TransactionReceiptConverter struct {
	threshold int // 并发处理的阈值
}

// NewTransactionReceiptConverter 创建一个新的TransactionReceiptConverter实例
func NewTransactionReceiptConverter() *TransactionReceiptConverter {
	return &TransactionReceiptConverter{
		threshold: 30, // 设置默认阈值
	}
}

// convertLogs 将eth.Log切片转换为domain.Log切片
func (c *TransactionReceiptConverter) convertLogs(logs []eth.Log) []Log {
	if logs == nil {
		return nil
	}

	return parallel.Process(logs, func(log eth.Log) Log {
		// 转换可选字段
		var logIndex, txIndex, txLogIndex, logType *string
		if log.LogIndex != nil {
			li := log.LogIndex.Big().String()
			logIndex = &li
		}
		if log.TxIndex != nil {
			ti := log.TxIndex.Big().String()
			txIndex = &ti
		}
		if log.TxLogIndex != nil {
			tli := log.TxLogIndex.Big().String()
			txLogIndex = &tli
		}
		if log.Type != nil {
			t := log.Type
			logType = t
		}

		// 转换区块号
		var blockNumber *string
		if log.BlockNumber != nil {
			bn := log.BlockNumber.Big().String()
			blockNumber = &bn
		}

		// 转换其他字段
		txHash := log.TxHash.String()
		blockHash := log.BlockHash.String()

		// 转换地址和数据字段
		var address, data string
		if log.Address != "" {
			addr := log.Address.String()
			address = addr
		}
		if log.Data != "" {
			d := log.Data.String()
			data = d
		}

		return Log{
			Removed:     log.Removed,
			LogIndex:    logIndex,
			TxIndex:     txIndex,
			TxHash:      &txHash,
			BlockHash:   &blockHash,
			BlockNumber: blockNumber,
			Address:     address,
			Data:        data,
			Topics:      c.convertTopics(log.Topics),
			TxLogIndex:  txLogIndex,
			Type:        logType,
		}
	}, c.threshold)
}

// convertTopics 将eth.Topics类型转换为[]string类型
func (c *TransactionReceiptConverter) convertTopics(topics []eth.Hash) []string {
	return parallel.Process(topics,
		func(topic eth.Hash) string { return topic.String() },
		c.threshold)
}

// ConvertToTransactionReceipt 将eth.TransactionReceipt转换为domain.TransactionReceipt
func (c *TransactionReceiptConverter) ConvertToTransactionReceipt(receipt *eth.TransactionReceipt) *TransactionReceipt {
	if receipt == nil {
		return nil
	}

	// 转换eth.Quantity类型为string类型
	var typ, status, effectiveGasPrice, blobGasPrice, blobGasUsed *string
	if receipt.Type != nil {
		t := receipt.Type.Big().String()
		typ = &t
	}
	if receipt.Status != nil {
		s := receipt.Status.Big().String()
		status = &s
	}
	if receipt.EffectiveGasPrice != nil {
		egp := receipt.EffectiveGasPrice.Big().String()
		effectiveGasPrice = &egp
	}
	if receipt.BlobGasPrice != nil {
		bgp := receipt.BlobGasPrice.Big().String()
		blobGasPrice = &bgp
	}
	if receipt.BlobGasUsed != nil {
		bgu := receipt.BlobGasUsed.Big().String()
		blobGasUsed = &bgu
	}

	// 处理可选的地址字段
	var to, contractAddress, root *string
	if receipt.To != nil {
		t := receipt.To.String()
		to = &t
	}
	if receipt.ContractAddress != nil {
		ca := receipt.ContractAddress.String()
		contractAddress = &ca
	}
	if receipt.Root != nil {
		r := receipt.Root.String()
		root = &r
	}

	return &TransactionReceipt{
		Type:              typ,
		TransactionHash:   receipt.TransactionHash.String(),
		TransactionIndex:  receipt.TransactionIndex.Big().String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Big().String(),
		From:              receipt.From.String(),
		To:                to,
		CumulativeGasUsed: receipt.CumulativeGasUsed.Big().String(),
		GasUsed:           receipt.GasUsed.Big().String(),
		ContractAddress:   contractAddress,
		Logs:              c.convertLogs(receipt.Logs),
		LogsBloom:         receipt.LogsBloom.String(),
		Root:              root,
		Status:            status,
		EffectiveGasPrice: effectiveGasPrice,
		BlobGasPrice:      blobGasPrice,
		BlobGasUsed:       blobGasUsed,
	}
}
