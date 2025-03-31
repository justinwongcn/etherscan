// Package service 提供以太坊区块链数据查询和处理的核心业务逻辑服务
package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/justinwongcn/etherscan/internal/ethereum"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// BlockService 实现了BlockServiceInterface接口
// 该结构体封装了与以太坊节点交互的客户端，提供区块数据查询的具体实现
type BlockService struct {
	// client 是与以太坊节点通信的客户端实例
	client *ethereum.Client
}

// NewBlockService 创建并初始化一个新的BlockService实例
// 参数:
//   - client: 已初始化的以太坊客户端实例，用于与节点通信
//
// 返回:
//   - *BlockService: 初始化完成的服务实例
func NewBlockService(client *ethereum.Client) *BlockService {
	return &BlockService{
		client: client,
	}
}

// GetLatestBlockHeight 实现了BlockServiceInterface接口中的同名方法
// 通过调用以太坊客户端获取当前网络的最新区块高度
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回:
//   - uint64: 最新区块的高度编号
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return s.client.GetLatestBlockNumber(ctx)
}

// parseBlockParameter 解析并标准化区块参数格式
// 将用户输入的区块标识符转换为以太坊API支持的格式
// 参数:
//   - blockHashOrNumber: 区块标识符，可以是区块号、区块哈希或特殊标识符
//
// 返回:
//   - string: 转换后的标准格式参数
//   - error: 如果参数格式无效，将返回错误信息
func (s *BlockService) parseBlockParameter(blockHashOrNumber string) (string, error) {
	// 判断是否为区块哈希（以0x开头的十六进制字符串）
	if len(blockHashOrNumber) >= 2 && blockHashOrNumber[:2] == "0x" {
		return blockHashOrNumber, nil
	}
	// 判断是否为特殊标识符
	if blockHashOrNumber == ethereum.BlockLatest || blockHashOrNumber == ethereum.BlockEarliest || blockHashOrNumber == ethereum.BlockPending {
		return blockHashOrNumber, nil
	}

	// 将字符串转换为big.Int
	number := new(big.Int)
	_, ok := number.SetString(blockHashOrNumber, 10)
	if !ok {
		return "", fmt.Errorf("invalid block number: %s", blockHashOrNumber)
	}

	// 转换为十六进制格式并添加0x前缀
	return fmt.Sprintf("0x%x", number), nil
}

// GetBlock 实现了BlockServiceInterface接口中的同名方法
// 根据区块标识符获取区块的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - *eth.Block: 包含区块完整信息的结构体指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetBlock(ctx context.Context, blockHashOrNumber string) (*eth.Block, error) {
	// 解析并标准化区块参数
	param, err := s.parseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法
	// 如果是区块哈希（以0x开头且长度大于10的十六进制字符串）
	if len(param) >= 2 && param[:2] == "0x" && len(param) > 10 {
		return s.client.GetBlockByHash(ctx, param, true)
	}
	return s.client.GetBlockByNumber(ctx, param, true)
}

// GetBlockTransactionCount 实现了BlockServiceInterface接口中的同名方法
// 获取指定区块中包含的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - uint64: 区块中的交易数量
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetBlockTransactionCount(ctx context.Context, blockHashOrNumber string) (uint64, error) {
	// 解析并标准化区块参数
	param, err := s.parseBlockParameter(blockHashOrNumber)
	if err != nil {
		return 0, err
	}

	// 根据参数类型选择适当的查询方法
	// 如果是区块哈希（以0x开头且长度大于10的十六进制字符串）
	if len(param) >= 2 && param[:2] == "0x" && len(param) > 10 {
		return s.client.GetBlockTransactionCountByHash(ctx, param)
	}
	return s.client.GetBlockTransactionCountByNumber(ctx, param)
}

// GetTransactionCount 实现了BlockServiceInterface接口中的同名方法
// 获取指定地址在特定区块的交易数量
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 要查询的以太坊地址
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//
// 返回:
//   - uint64: 该地址的交易数量
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetTransactionCount(ctx context.Context, address string, blockHashOrNumber string) (uint64, error) {
	// 解析并标准化区块参数
	param, err := s.parseBlockParameter(blockHashOrNumber)
	if err != nil {
		return 0, err
	}

	// 调用以太坊客户端获取交易数量
	return s.client.GetTransactionCount(ctx, address, param)
}

// GetTransactionByHash 实现了BlockServiceInterface接口中的同名方法
// 根据交易哈希获取交易的详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - txHash: 交易哈希（32字节的十六进制字符串）
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetTransactionByHash(ctx context.Context, txHash string) (*eth.Transaction, error) {
	// 调用以太坊客户端获取交易信息
	return s.client.GetTransactionByHash(ctx, txHash)
}

// GetTransactionByIndex 实现了BlockServiceInterface接口中的同名方法
// 根据区块标识符和交易索引获取交易详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - blockHashOrNumber: 区块标识符，支持区块号、区块哈希和特殊标识符
//   - index: 交易在区块中的索引位置
//
// 返回:
//   - *eth.Transaction: 包含交易完整信息的结构体指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *BlockService) GetTransactionByIndex(ctx context.Context, blockHashOrNumber string, index uint64) (*eth.Transaction, error) {
	// 解析并标准化区块参数
	param, err := s.parseBlockParameter(blockHashOrNumber)
	if err != nil {
		return nil, err
	}

	// 根据参数类型选择适当的查询方法
	// 如果是区块哈希（以0x开头且长度大于10的十六进制字符串）
	if len(param) >= 2 && param[:2] == "0x" && len(param) > 10 {
		return s.client.GetTransactionByBlockHashAndIndex(ctx, param, index)
	}
	return s.client.GetTransactionByBlockNumberAndIndex(ctx, param, index)
}
