// Package service 提供以太坊区块链的应用服务层实现
package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/repository"
	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// AccountService 提供账户相关的服务实现
// 该结构体遵循DDD架构原则，通过依赖仓储接口而非具体实现来提供账户相关的业务逻辑
type AccountService struct {
	// accountRepo 是账户数据访问的仓储接口
	// 通过依赖倒置原则，应用层依赖抽象而非具体实现
	accountRepo repository.AccountRepository
}

// NewAccountService 创建一个新的账户服务实例
// 参数:
//   - accountRepo: 实现了AccountRepository接口的仓储实例，用于账户数据访问
//
// 返回:
//   - *AccountService: 账户服务实现
func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
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

// GetBalance 获取账户余额
func (s *AccountService) GetBalance(ctx context.Context, address string, numberOrTag string) (string, error) {
	// 获取账户聚合根
	accountAggregate, err := s.GetAccountAggregate(ctx, address, numberOrTag)
	if err != nil {
		return "", err
	}

	// 从聚合根获取余额
	return accountAggregate.Balance().String(), nil
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

// GetBalances 批量获取账户余额
func (s *AccountService) GetBalances(ctx context.Context, addresses []string, numberOrTag string) (map[string]string, error) {
	// 批量获取账户聚合根
	accountAggregates, err := s.GetMultipleAccountAggregates(ctx, addresses, numberOrTag)
	if err != nil {
		return nil, err
	}

	// 从聚合根提取余额信息
	result := make(map[string]string)
	for addr, aggregate := range accountAggregates {
		result[addr] = aggregate.Balance().String()
	}
	return result, nil
}

// GetAccountAggregate 获取账户聚合根（新的DDD架构方法）
// 这个方法使用新的聚合根模式，将逐步替代GetBalance方法
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - address: 以太坊账户地址
//   - numberOrTag: 区块号或预定义标签（如 "latest", "earliest", "pending"）
//
// 返回:
//   - *aggregate.AccountAggregate: 包含账户完整信息的聚合根指针
//   - error: 如果查询过程中发生错误，将返回相应的错误信息
func (s *AccountService) GetAccountAggregate(ctx context.Context, address string, numberOrTag string) (*aggregate.AccountAggregate, error) {
	// 验证并转换地址
	accountAddress, err := valueobject.NewAddress(address)
	if err != nil {
		return nil, fmt.Errorf("无效的账户地址: %w", err)
	}

	// 解析区块参数（验证格式）
	_, err = s.parseBlockParameter(numberOrTag)
	if err != nil {
		return nil, fmt.Errorf("无效的区块参数: %w", err)
	}

	// 获取账户余额
	balanceUint64, err := s.accountRepo.GetBalance(ctx, address, numberOrTag)
	if err != nil {
		return nil, fmt.Errorf("获取账户余额失败: %w", err)
	}

	// 转换为Wei值对象
	balance := valueobject.NewWeiFromUint64(balanceUint64)

	// 创建外部拥有账户聚合根
	accountAggregate, err := aggregate.NewExternallyOwnedAccount(accountAddress)
	if err != nil {
		return nil, fmt.Errorf("创建账户聚合根失败: %w", err)
	}

	// 更新余额
	accountAggregate.UpdateBalance(balance)

	return accountAggregate, nil
}

// parseBlockParameter 解析区块参数
func (s *AccountService) parseBlockParameter(numberOrTag string) (valueobject.BlockNumber, error) {
	switch numberOrTag {
	case "latest":
		return valueobject.NewLatestBlockNumber(), nil
	case "earliest":
		return valueobject.NewEarliestBlockNumber(), nil
	case "pending":
		return valueobject.NewPendingBlockNumber(), nil
	default:
		// 尝试解析为数字
		return valueobject.NewBlockNumber(numberOrTag)
	}
}

// GetMultipleAccountAggregates 批量获取账户聚合根
func (s *AccountService) GetMultipleAccountAggregates(ctx context.Context, addresses []string, numberOrTag string) (map[string]*aggregate.AccountAggregate, error) {
	result := make(map[string]*aggregate.AccountAggregate)

	// 解析区块参数（验证格式）
	_, err := s.parseBlockParameter(numberOrTag)
	if err != nil {
		return nil, fmt.Errorf("无效的区块参数: %w", err)
	}

	// 批量获取余额
	balances, err := s.accountRepo.GetBalances(ctx, addresses, numberOrTag)
	if err != nil {
		return nil, fmt.Errorf("批量获取余额失败: %w", err)
	}

	// 为每个地址创建账户聚合根
	for _, address := range addresses {
		accountAddress, err := valueobject.NewAddress(address)
		if err != nil {
			// 跳过无效地址，但记录错误
			continue
		}

		balanceUint64, exists := balances[address]
		if !exists {
			// 如果没有余额信息，设置为0
			balanceUint64 = 0
		}

		balance := valueobject.NewWeiFromUint64(balanceUint64)

		accountAggregate, err := aggregate.NewExternallyOwnedAccount(accountAddress)
		if err != nil {
			// 跳过创建失败的聚合根
			continue
		}

		// 更新余额
		accountAggregate.UpdateBalance(balance)

		result[address] = accountAggregate
	}

	return result, nil
}

// formatBlockNumber 格式化区块号为字符串
func (s *AccountService) formatBlockNumber(blockNumber valueobject.BlockNumber) string {
	if blockNumber.IsNumber() {
		number, err := blockNumber.Uint64()
		if err == nil {
			return strconv.FormatUint(number, 10)
		}
	}

	// 如果是特殊标签，返回标签值
	if tag, err := blockNumber.Tag(); err == nil {
		return tag
	}

	return "latest"
}
