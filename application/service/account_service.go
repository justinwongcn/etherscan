// Package service 提供以太坊区块链的应用服务层实现
package service

import (
    "context"
    "fmt"
    "strconv"

    "github.com/justinwongcn/etherscan/internal/ethereum"
)

// AccountService 提供账户相关的服务实现
// 该结构体负责处理与以太坊账户相关的业务逻辑，包括余额查询等操作
type AccountService struct {
    client *ethereum.Client // 以太坊客户端实例
}

// NewAccountService 创建一个新的账户服务实例
// 参数:
//   - client: 以太坊客户端实例，用于与区块链交互
//
// 返回:
//   - AccountServiceInterface: 账户服务接口实现
func NewAccountService(client *ethereum.Client) AccountServiceInterface {
    return &AccountService{client: client}
}

// GetBalance 获取指定地址的账户余额
// 该方法查询指定区块高度或状态下账户的ETH余额
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - numberOrTag: 区块号或预定义标签（如 "latest", "earliest", "pending"）
//
// 返回:
//   - string: 账户余额，以十进制字符串表示
//   - error: 操作过程中的错误信息，如果没有错误则为nil
func (s *AccountService) GetBalance(ctx context.Context, address string, numberOrTag string) (string, error) {
    // 调用以太坊客户端获取余额
    balance, err := s.client.GetBalance(ctx, address, numberOrTag)
    if err != nil {
        return "", fmt.Errorf("failed to get balance: %v", err)
    }

    // 将uint64类型的余额转换为十进制字符串
    return strconv.FormatUint(balance, 10), nil
}

// GetBalances 批量获取多个地址的账户余额
// 该方法通过一次请求获取多个账户的ETH余额，提高查询效率
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - addresses: 以太坊账户地址切片
//   - numberOrTag: 区块号或预定义标签（如 "latest", "earliest", "pending"）
//
// 返回:
//   - map[string]string: 地址到余额的映射，余额以十进制字符串表示
//   - error: 操作过程中的错误信息，如果没有错误则为nil
func (s *AccountService) GetBalances(ctx context.Context, addresses []string, numberOrTag string) (map[string]string, error) {
    // 调用以太坊客户端批量获取余额
    balances, err := s.client.GetBalances(ctx, addresses, numberOrTag)
    if err != nil {
        return nil, fmt.Errorf("failed to get balances: %v", err)
    }

    // 将uint64类型的余额转换为十进制字符串
    result := make(map[string]string)
    for addr, balance := range balances {
        result[addr] = strconv.FormatUint(balance, 10)
    }
    return result, nil
}
