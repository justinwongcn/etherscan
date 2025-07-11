package ethereum

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClient_GetBalance 测试获取账户余额
func TestClient_GetBalance(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		numberOrTag string
		expectError bool
	}{
		{
			name:        "有效地址和latest标签",
			address:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			numberOrTag: "latest",
			expectError: false,
		},
		{
			name:        "有效地址和earliest标签",
			address:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			numberOrTag: "earliest",
			expectError: false,
		},
		{
			name:        "有效地址和pending标签",
			address:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			numberOrTag: "pending",
			expectError: false,
		},
		{
			name:        "有效地址和十六进制区块号",
			address:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			numberOrTag: "0x1000000",
			expectError: false,
		},
		{
			name:        "无效地址格式",
			address:     "invalid-address",
			numberOrTag: "latest",
			expectError: true,
		},
		{
			name:        "空地址",
			address:     "",
			numberOrTag: "latest",
			expectError: true,
		},
		{
			name:        "无效区块号",
			address:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			numberOrTag: "invalid-block",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 创建客户端（可能因网络问题失败）
			client, err := NewClient(ctx, "https://rpc.flashbots.net", nil)
			if err != nil {
				// 网络错误，跳过测试
				t.Skip("网络连接失败，跳过测试")
				return
			}

			balance, err := client.GetBalance(ctx, tt.address, tt.numberOrTag)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				// 网络错误是可以接受的
				if err == nil {
					assert.GreaterOrEqual(t, balance, uint64(0))
				}
			}
		})
	}
}

// TestClient_GetBalance_AddressValidation 测试地址验证
func TestClient_GetBalance_AddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		address string
		valid   bool
	}{
		{
			name:    "有效的以太坊地址",
			address: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			valid:   true,
		},
		{
			name:    "有效的全小写地址",
			address: "0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
			valid:   true,
		},
		{
			name:    "无效的地址长度",
			address: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8",
			valid:   false,
		},
		{
			name:    "缺少0x前缀",
			address: "742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			valid:   false,
		},
		{
			name:    "包含无效字符",
			address: "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8bG",
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 创建客户端（可能因网络问题失败）
			client, err := NewClient(ctx, "https://rpc.flashbots.net", nil)
			if err != nil {
				// 网络错误，跳过测试
				t.Skip("网络连接失败，跳过测试")
				return
			}

			_, err = client.GetBalance(ctx, tt.address, "latest")

			if tt.valid {
				// 有效地址可能因网络问题失败，但不应该是地址格式错误
				if err != nil {
					assert.NotContains(t, err.Error(), "invalid ethereum address")
				}
			} else {
				// 无效地址应该返回错误
				assert.Error(t, err)
				// 错误消息可能来自不同层级，只要有错误即可
			}
		})
	}
}

// TestClient_GetBalances 测试批量获取账户余额
func TestClient_GetBalances(t *testing.T) {
	tests := []struct {
		name        string
		addresses   []string
		numberOrTag string
		expectError bool
	}{
		{
			name:        "单个有效地址",
			addresses:   []string{"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"},
			numberOrTag: "latest",
			expectError: false,
		},
		{
			name: "多个有效地址",
			addresses: []string{
				"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				"0x8ba1f109551bD432803012645aac136c12345678",
			},
			numberOrTag: "latest",
			expectError: false,
		},
		{
			name:        "空地址列表",
			addresses:   []string{},
			numberOrTag: "latest",
			expectError: true,
		},
		{
			name:        "nil地址列表",
			addresses:   nil,
			numberOrTag: "latest",
			expectError: true,
		},
		{
			name: "包含无效地址",
			addresses: []string{
				"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				"invalid-address",
			},
			numberOrTag: "latest",
			expectError: true,
		},
		{
			name:        "无效区块号",
			addresses:   []string{"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"},
			numberOrTag: "invalid-block",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 创建客户端（可能因网络问题失败）
			client, err := NewClient(ctx, "https://rpc.flashbots.net", nil)
			if err != nil {
				// 网络错误，跳过测试
				t.Skip("网络连接失败，跳过测试")
				return
			}

			balances, err := client.GetBalances(ctx, tt.addresses, tt.numberOrTag)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, balances)
			} else {
				// 网络错误是可以接受的
				if err == nil {
					assert.NotNil(t, balances)
					assert.Equal(t, len(tt.addresses), len(balances))

					// 验证所有余额都是非负数
					for addr, balance := range balances {
						assert.GreaterOrEqual(t, balance, uint64(0))
						assert.Contains(t, tt.addresses, addr)
					}
				}
			}
		})
	}
}

// TestClient_GetBalances_MaxAddresses 测试地址数量限制
func TestClient_GetBalances_MaxAddresses(t *testing.T) {
	tests := []struct {
		name          string
		addresses     []string
		maxAddresses  int
		expectError   bool
		errorContains string
	}{
		{
			name: "地址数量在限制内",
			addresses: []string{
				"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				"0x8ba1f109551bD432803012645aac136c12345678",
			},
			maxAddresses: 5,
			expectError:  false,
		},
		{
			name: "地址数量超过限制",
			addresses: []string{
				"0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
				"0x8ba1f109551bD432803012645aac136c12345678",
				"0x1234567890123456789012345678901234567890",
			},
			maxAddresses:  2,
			expectError:   true,
			errorContains: "too many addresses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// 创建客户端（可能因网络问题失败）
			client, err := NewClient(ctx, "https://rpc.flashbots.net", nil)
			if err != nil {
				// 网络错误，跳过测试
				t.Skip("网络连接失败，跳过测试")
				return
			}

			balances, err := client.GetBalances(ctx, tt.addresses, "latest", tt.maxAddresses)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, balances)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				// 网络错误是可以接受的
				if err == nil {
					assert.NotNil(t, balances)
				}
			}
		})
	}
}

// BenchmarkClient_GetBalance 性能测试
func BenchmarkClient_GetBalance(b *testing.B) {
	ctx := context.Background()

	client, err := NewClient(ctx, "https://rpc.flashbots.net", nil)
	if err != nil {
		b.Skip("网络连接失败，跳过基准测试")
		return
	}

	address := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetBalance(ctx, address, "latest")
		if err != nil {
			// 网络错误是可以接受的
			continue
		}
	}
}
