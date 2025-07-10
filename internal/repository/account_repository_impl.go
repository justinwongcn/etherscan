// Package repository 提供领域仓储接口的具体实现
// 该包属于基础设施层，负责与外部数据源（以太坊节点）进行交互
package repository

import (
	"context"

	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/internal/ethereum"
)

// AccountRepositoryImpl 实现了domain.repository.AccountRepository接口
// 该结构体封装了与以太坊节点交互的客户端，提供账户数据访问的具体实现
type AccountRepositoryImpl struct {
	// client 是与以太坊节点通信的客户端实例
	// 通过该客户端执行实际的RPC调用来获取账户数据
	client *ethereum.Client
}

// NewAccountRepository 创建并初始化一个新的AccountRepositoryImpl实例
// 该函数是AccountRepositoryImpl的工厂方法，确保正确的依赖注入
// 参数:
//   - client: 已初始化的以太坊客户端实例，用于与节点通信
//
// 返回:
//   - repository.AccountRepository: 实现了AccountRepository接口的实例
func NewAccountRepository(client *ethereum.Client) repository.AccountRepository {
	return &AccountRepositoryImpl{
		client: client,
	}
}

// GetBalance 实现了AccountRepository接口中的同名方法
// 获取指定地址的账户余额
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//
// 返回:
//   - uint64: 账户余额，以wei为单位
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *AccountRepositoryImpl) GetBalance(ctx context.Context, address string, blockParameter string) (uint64, error) {
	return r.client.GetBalance(ctx, address, blockParameter)
}

// GetBalances 实现了AccountRepository接口中的同名方法
// 批量获取多个地址的账户余额
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - addresses: 以太坊账户地址切片
//   - blockParameter: 区块参数，可以是具体的区块号（十六进制字符串）或特殊标识符
//
// 返回:
//   - map[string]uint64: 地址到余额的映射，余额以wei为单位
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (r *AccountRepositoryImpl) GetBalances(ctx context.Context, addresses []string, blockParameter string) (map[string]uint64, error) {
	return r.client.GetBalances(ctx, addresses, blockParameter)
}
