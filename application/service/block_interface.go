package service

import (
	"context"

	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockServiceInterface 区块服务接口
type BlockServiceInterface interface {
	// GetLatestBlockHeight 获取最新区块高度
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	// GetBlockByNumber 获取指定区块号的区块信息
	GetBlockByNumber(ctx context.Context, numberOrTag string) (*eth.Block, error)
	// GetBlockByHash 获取指定区块哈希的区块信息
	GetBlockByHash(ctx context.Context, blockHash string) (*eth.Block, error)
	// GetTransactionCount 获取指定区块的交易数量
	GetTransactionCount(ctx context.Context, blockHashOrNumber string) (uint64, error)
}
