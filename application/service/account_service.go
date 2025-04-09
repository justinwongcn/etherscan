package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/justinwongcn/etherscan/internal/ethereum"
)

// AccountService 账户服务实现类
type AccountService struct {
	client *ethereum.Client
}

// NewAccountService 创建账户服务实例
func NewAccountService(client *ethereum.Client) AccountServiceInterface {
	return &AccountService{client: client}
}

// GetBalance 获取指定地址的账户余额
func (s *AccountService) GetBalance(ctx context.Context, address string, numberOrTag string) (string, error) {
	// 调用以太坊客户端获取余额
	balance, err := s.client.GetBalance(ctx, address, numberOrTag)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %v", err)
	}

	// 将uint64类型的余额转换为十进制字符串
	return strconv.FormatUint(balance, 10), nil
}
