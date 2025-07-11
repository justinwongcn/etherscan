// Package handler 提供HTTP请求处理器的基础设施
package handler

import (
	"net/http"
	"strconv"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BaseHandler 提供所有Handler的基础功能
type BaseHandler struct{}

// NewBaseHandler 创建基础Handler实例
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// RespondWithError 统一的错误响应处理
func (h *BaseHandler) RespondWithError(ctx *ant.Context, statusCode int, message string) {
	_ = ctx.RespJSON(statusCode, H{
		"error": message,
	})
}

// RespondWithSuccess 统一的成功响应处理
func (h *BaseHandler) RespondWithSuccess(ctx *ant.Context, data H) {
	_ = ctx.RespJSONOK(data)
}

// RespondWithInternalError 统一的内部错误响应处理
func (h *BaseHandler) RespondWithInternalError(ctx *ant.Context, err error) {
	h.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
}

// RespondWithBadRequest 统一的请求错误响应处理
func (h *BaseHandler) RespondWithBadRequest(ctx *ant.Context, message string) {
	h.RespondWithError(ctx, http.StatusBadRequest, message)
}

// ConvertBlockToResponse 将eth.Block转换为API响应格式
// 确保数值字段使用十进制格式
func (h *BaseHandler) ConvertBlockToResponse(ethBlock *eth.Block, fullTx bool) map[string]any {
	if ethBlock == nil {
		return nil
	}

	response := map[string]any{
		"number":           strconv.FormatUint(ethBlock.Number.UInt64(), 10),
		"hash":             ethBlock.Hash.String(),
		"parentHash":       ethBlock.ParentHash.String(),
		"sha3Uncles":       ethBlock.SHA3Uncles.String(),
		"logsBloom":        ethBlock.LogsBloom.String(),
		"transactionsRoot": ethBlock.TransactionsRoot.String(),
		"stateRoot":        ethBlock.StateRoot.String(),
		"receiptsRoot":     ethBlock.ReceiptsRoot.String(),
		"miner":            ethBlock.Miner.String(),
		"difficulty":       strconv.FormatUint(ethBlock.Difficulty.UInt64(), 10),
		"totalDifficulty":  strconv.FormatUint(ethBlock.TotalDifficulty.UInt64(), 10),
		"extraData":        ethBlock.ExtraData.String(),
		"size":             strconv.FormatUint(ethBlock.Size.UInt64(), 10),
		"gasLimit":         strconv.FormatUint(ethBlock.GasLimit.UInt64(), 10),
		"gasUsed":          strconv.FormatUint(ethBlock.GasUsed.UInt64(), 10),
		"timestamp":        strconv.FormatUint(ethBlock.Timestamp.UInt64(), 10),
		"uncles":           []any{}, // 简化处理，返回空数组
	}

	// 处理可选字段
	if ethBlock.Nonce != nil {
		response["nonce"] = ethBlock.Nonce.String()
	}

	if ethBlock.BaseFeePerGas != nil {
		response["baseFeePerGas"] = strconv.FormatUint(ethBlock.BaseFeePerGas.UInt64(), 10)
	}

	// 处理交易
	if fullTx && ethBlock.Transactions != nil {
		transactions := make([]map[string]any, len(ethBlock.Transactions))
		for i, tx := range ethBlock.Transactions {
			transactions[i] = h.ConvertTransactionToResponse(&tx.Transaction)
		}
		response["transactions"] = transactions
	} else if ethBlock.Transactions != nil {
		hashes := make([]string, len(ethBlock.Transactions))
		for i, tx := range ethBlock.Transactions {
			hashes[i] = tx.Transaction.Hash.String()
		}
		response["transactions"] = hashes
	}

	return response
}

// ConvertTransactionToResponse 将eth.Transaction转换为API响应格式
func (h *BaseHandler) ConvertTransactionToResponse(ethTx *eth.Transaction) map[string]any {
	if ethTx == nil {
		return nil
	}

	response := map[string]any{
		"hash":  ethTx.Hash.String(),
		"nonce": strconv.FormatUint(ethTx.Nonce.UInt64(), 10),
		"from":  ethTx.From.String(),
		"value": strconv.FormatUint(ethTx.Value.UInt64(), 10),
		"gas":   strconv.FormatUint(ethTx.Gas.UInt64(), 10),
		"input": ethTx.Input.String(),
		"v":     strconv.FormatUint(ethTx.V.UInt64(), 10),
		"r":     ethTx.R.String(),
		"s":     ethTx.S.String(),
		"type":  strconv.FormatUint(ethTx.Type.UInt64(), 10),
	}

	// 处理可选字段
	if ethTx.To != nil {
		response["to"] = ethTx.To.String()
	}

	if ethTx.BlockHash != nil {
		response["blockHash"] = ethTx.BlockHash.String()
	}

	if ethTx.BlockNumber != nil {
		response["blockNumber"] = strconv.FormatUint(ethTx.BlockNumber.UInt64(), 10)
	}

	if ethTx.Index != nil {
		response["transactionIndex"] = strconv.FormatUint(ethTx.Index.UInt64(), 10)
	}

	if ethTx.GasPrice != nil {
		response["gasPrice"] = strconv.FormatUint(ethTx.GasPrice.UInt64(), 10)
	}

	if ethTx.MaxFeePerGas != nil {
		response["maxFeePerGas"] = strconv.FormatUint(ethTx.MaxFeePerGas.UInt64(), 10)
	}

	if ethTx.MaxPriorityFeePerGas != nil {
		response["maxPriorityFeePerGas"] = strconv.FormatUint(ethTx.MaxPriorityFeePerGas.UInt64(), 10)
	}

	if ethTx.ChainId != nil {
		response["chainId"] = strconv.FormatUint(ethTx.ChainId.UInt64(), 10)
	}

	// EIP-4844: 处理blob交易相关字段
	if ethTx.MaxFeePerBlobGas != nil {
		response["maxFeePerBlobGas"] = strconv.FormatUint(ethTx.MaxFeePerBlobGas.UInt64(), 10)
	}

	if len(ethTx.BlobVersionedHashes) > 0 {
		blobHashes := make([]string, len(ethTx.BlobVersionedHashes))
		for i, hash := range ethTx.BlobVersionedHashes {
			blobHashes[i] = hash.String()
		}
		response["blobVersionedHashes"] = blobHashes
	}

	return response
}

// ConvertTransactionReceiptToResponse 将eth.TransactionReceipt转换为API响应格式
// 确保数值字段使用十进制格式
func (h *BaseHandler) ConvertTransactionReceiptToResponse(ethReceipt *eth.TransactionReceipt) map[string]any {
	if ethReceipt == nil {
		return nil
	}

	response := map[string]any{
		"transactionHash":   ethReceipt.TransactionHash.String(),
		"transactionIndex":  strconv.FormatUint(ethReceipt.TransactionIndex.UInt64(), 10),
		"blockHash":         ethReceipt.BlockHash.String(),
		"blockNumber":       strconv.FormatUint(ethReceipt.BlockNumber.UInt64(), 10),
		"from":              ethReceipt.From.String(),
		"cumulativeGasUsed": strconv.FormatUint(ethReceipt.CumulativeGasUsed.UInt64(), 10),
		"gasUsed":           strconv.FormatUint(ethReceipt.GasUsed.UInt64(), 10),
		"logsBloom":         ethReceipt.LogsBloom.String(),
		"type":              strconv.FormatUint(ethReceipt.Type.UInt64(), 10),
		"logs":              []any{}, // 简化处理，返回空数组
	}

	// 处理可选字段
	if ethReceipt.To != nil {
		response["to"] = ethReceipt.To.String()
	}

	if ethReceipt.ContractAddress != nil {
		response["contractAddress"] = ethReceipt.ContractAddress.String()
	}

	if ethReceipt.Root != nil {
		response["root"] = ethReceipt.Root.String()
	}

	if ethReceipt.Status != nil {
		response["status"] = strconv.FormatUint(ethReceipt.Status.UInt64(), 10)
	}

	if ethReceipt.EffectiveGasPrice != nil {
		response["effectiveGasPrice"] = strconv.FormatUint(ethReceipt.EffectiveGasPrice.UInt64(), 10)
	} else {
		response["effectiveGasPrice"] = "0"
	}

	return response
}
