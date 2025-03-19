package ethereum

import (
	"context"
	"github.com/justinwongcn/go-ethlibs/eth"

	"github.com/justinwongcn/go-ethlibs/node"
)

// Client 封装以太坊客户端
type Client struct {
	nodeClient node.Client
}

// NewClient 创建一个新的以太坊客户端
func NewClient(ctx context.Context, nodeURL string) (*Client, error) {
	client, err := node.NewClient(ctx, nodeURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		nodeClient: client,
	}, nil
}

// GasPrice 获取当前 gas 价格
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 当前 gas 价格（单位：wei）
//   - error: 操作过程中可能发生的错误

func (c *Client) GasPrice(ctx context.Context) (uint64, error) {
	return c.nodeClient.GasPrice(ctx)
}

// GetLatestBlockNumber 获取最新区块高度
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - uint64: 当前客户端所在的最新区块高度
//   - error: 操作过程中可能发生的错误
func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.nodeClient.BlockNumber(ctx)
}

// GetBalance 获取指定地址的账户余额
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - address: string 要查询余额的账户地址（20字节）
//   - blockNumber: string 区块号，可以是具体数值或"latest"、"earliest"、"pending"
//
// Returns:
//   - uint64: 账户余额（单位：wei）
//   - error: 操作过程中可能发生的错误
func (c *Client) GetBalance(ctx context.Context, address string, numberOrTag string) (uint64, error) {
	addr, err := eth.NewAddress(address)
	if err != nil {
		return 0, err
	}
	numOrTag := eth.MustBlockNumberOrTag(numberOrTag)
	return c.nodeClient.GetBalance(ctx, *addr, *numOrTag)
}
