package service

import (
	"context"

	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockService 区块服务层，处理区块相关的业务逻辑
type BlockService struct {
	client *ethereum.Client
}

// NewBlockService 创建区块服务实例
func NewBlockService(client *ethereum.Client) *BlockService {
	return &BlockService{
		client: client,
	}
}

// GetLatestBlockHeight 获取最新区块高度
func (s *BlockService) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return s.client.GetLatestBlockNumber(ctx)
}

// GetBlockByNumber 获取指定区块号的区块信息
func (s *BlockService) GetBlockByNumber(ctx context.Context, numberOrTag string) (*eth.Block, error) {
	return s.client.GetBlockByNumber(ctx, numberOrTag, true)
}

// GetBlockByHash 获取指定区块哈希的区块信息
func (s *BlockService) GetBlockByHash(ctx context.Context, blockHash string) (*eth.Block, error) {
	return s.client.GetBlockByHash(ctx, blockHash, true)
}

// GetTransactionCount 获取指定区块的交易数量
func (s *BlockService) GetTransactionCount(ctx context.Context, blockHashOrNumber string) (uint64, error) {
	// 判断是否为区块哈希（以0x开头的十六进制字符串）
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" {
		return s.client.GetBlockTransactionCountByHash(ctx, blockHashOrNumber)
	}
	// 否则视为区块号
	return s.client.GetBlockTransactionCountByNumber(ctx, blockHashOrNumber)
}
