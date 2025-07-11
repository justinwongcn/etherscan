package service

import (
	"testing"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

// TestNewAccountManagementService 测试构造函数
func TestNewAccountManagementService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "创建AccountManagementService成功",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAccountManagementService()
			assert.NotNil(t, service)
		})
	}
}

// TestAccountManagementService_ProcessTransfer 测试转账处理
func TestAccountManagementService_ProcessTransfer(t *testing.T) {
	tests := []struct {
		name             string
		fromBalance      uint64
		toBalance        uint64
		transferAmount   uint64
		transactionFee   uint64
		expectSuccess    bool
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:           "转账成功",
			fromBalance:    1000000000000000000, // 1 ETH
			toBalance:      500000000000000000,  // 0.5 ETH
			transferAmount: 100000000000000000,  // 0.1 ETH
			transactionFee: 21000000000000000,   // 0.021 ETH
			expectSuccess:  true,
			expectError:    false,
		},
		{
			name:             "转账金额为零",
			fromBalance:      1000000000000000000,
			toBalance:        500000000000000000,
			transferAmount:   0,
			transactionFee:   21000000000000000,
			expectSuccess:    false,
			expectError:      true,
			expectedErrorMsg: "转账金额不能为零",
		},
		{
			name:             "余额不足",
			fromBalance:      100000000000000000, // 0.1 ETH
			toBalance:        500000000000000000,
			transferAmount:   200000000000000000, // 0.2 ETH
			transactionFee:   21000000000000000,
			expectSuccess:    false,
			expectError:      true,
			expectedErrorMsg: "发送方余额不足",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试账户
			fromAddress, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
			toAddress, _ := valueobject.NewAddress("0x8ba1f109551bD432803012645aac136c12345678")

			fromAccount, _ := aggregate.NewExternallyOwnedAccount(fromAddress)
			toAccount, _ := aggregate.NewExternallyOwnedAccount(toAddress)

			// 设置余额
			fromAccount.UpdateBalance(valueobject.NewWeiFromUint64(tt.fromBalance))
			toAccount.UpdateBalance(valueobject.NewWeiFromUint64(tt.toBalance))

			// 创建服务
			service := NewAccountManagementService()

			// 执行转账
			result, err := service.ProcessTransfer(
				fromAccount,
				toAccount,
				valueobject.NewWeiFromUint64(tt.transferAmount),
				valueobject.NewWeiFromUint64(tt.transactionFee),
			)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
				if result != nil {
					assert.Equal(t, tt.expectSuccess, result.Success)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectSuccess, result.Success)

				if tt.expectSuccess {
					// 验证转账金额
					transferredAmount, _ := result.TransferredAmount.Uint64()
					assert.Equal(t, tt.transferAmount, transferredAmount)

					// 验证手续费
					fee, _ := result.TransactionFee.Uint64()
					assert.Equal(t, tt.transactionFee, fee)

					// 验证余额变化
					expectedFromBalance := tt.fromBalance - tt.transferAmount - tt.transactionFee
					fromBalance, _ := result.FromBalance.Uint64()
					assert.Equal(t, expectedFromBalance, fromBalance)

					expectedToBalance := tt.toBalance + tt.transferAmount
					toBalance, _ := result.ToBalance.Uint64()
					assert.Equal(t, expectedToBalance, toBalance)
				}
			}
		})
	}
}

// TestAccountManagementService_ValidateAccountState 测试账户状态验证
func TestAccountManagementService_ValidateAccountState(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "正常账户验证成功",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建服务
			service := NewAccountManagementService()

			// 创建测试账户
			address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
			account, _ := aggregate.NewExternallyOwnedAccount(address)
			account.UpdateBalance(valueobject.NewWeiFromUint64(1000000000000000000))

			// 执行验证
			err := service.ValidateAccountState(account)

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
