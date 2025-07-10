package valueobject

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "有效的地址（带0x前缀）",
			input:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			wantError: false,
			expected:  "0x742d35Cc6634C0532925A3B8D4C9dB96C4B4d8B6",
		},
		{
			name:      "有效的地址（不带0x前缀）",
			input:     "742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
			wantError: false,
			expected:  "0x742d35Cc6634C0532925A3B8D4C9dB96C4B4d8B6",
		},
		{
			name:      "有效的地址（全小写）",
			input:     "0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
			wantError: false,
			expected:  "0x742d35Cc6634C0532925A3B8D4C9dB96C4B4d8B6",
		},
		{
			name:      "有效的地址（全大写）",
			input:     "0x742D35CC6634C0532925A3B8D4C9DB96C4B4D8B6",
			wantError: false,
			expected:  "0x742d35Cc6634C0532925A3B8D4C9dB96C4B4d8B6",
		},
		{
			name:      "空字符串",
			input:     "",
			wantError: true,
		},
		{
			name:      "长度不正确（太短）",
			input:     "0x1234",
			wantError: true,
		},
		{
			name:      "长度不正确（太长）",
			input:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b612",
			wantError: true,
		},
		{
			name:      "包含无效字符",
			input:     "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8bg",
			wantError: true,
		},
		{
			name:      "零地址",
			input:     "0x0000000000000000000000000000000000000000",
			wantError: false,
			expected:  "0x0000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := NewAddress(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewAddress() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewAddress() 返回了意外的错误: %v", err)
				return
			}

			if address.String() != tt.expected {
				t.Errorf("NewAddress() = %v, 期望 %v", address.String(), tt.expected)
			}
		})
	}
}

func TestNewAddressFromBytes(t *testing.T) {
	// 创建20字节的地址
	bytes := make([]byte, 20)
	for i := 0; i < 20; i++ {
		bytes[i] = byte(i)
	}

	address, err := NewAddressFromBytes(bytes)
	if err != nil {
		t.Errorf("NewAddressFromBytes() 返回了意外的错误: %v", err)
	}

	expected := "0x000102030405060708090a0b0c0d0e0f10111213"
	if address.String() != expected {
		t.Errorf("NewAddressFromBytes() = %v, 期望 %v", address.String(), expected)
	}

	// 测试错误的字节长度
	wrongBytes := make([]byte, 16)
	_, err = NewAddressFromBytes(wrongBytes)
	if err == nil {
		t.Errorf("NewAddressFromBytes() 应该对错误的字节长度返回错误")
	}
}

func TestAddress_Value(t *testing.T) {
	addressStr := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	address, err := NewAddress(addressStr)
	if err != nil {
		t.Fatalf("NewAddress() 失败: %v", err)
	}

	expected := "0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6"
	if address.Value() != expected {
		t.Errorf("Value() = %v, 期望 %v", address.Value(), expected)
	}
}

func TestAddress_Bytes(t *testing.T) {
	addressStr := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	address, err := NewAddress(addressStr)
	if err != nil {
		t.Fatalf("NewAddress() 失败: %v", err)
	}

	bytes := address.Bytes()
	if len(bytes) != 20 {
		t.Errorf("Bytes() 长度 = %v, 期望 20", len(bytes))
	}

	// 验证第一个字节
	if bytes[0] != 0x74 {
		t.Errorf("Bytes()[0] = %v, 期望 0x74", bytes[0])
	}
}

func TestAddress_Equals(t *testing.T) {
	address1, _ := NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	address2, _ := NewAddress("742d35cc6634c0532925a3b8d4c9db96c4b4d8b6")
	address3, _ := NewAddress("0x8ba1f109551bD432803012645aac136c12345678")

	if !address1.Equals(address2) {
		t.Errorf("相同的地址应该相等")
	}

	if address1.Equals(address3) {
		t.Errorf("不同的地址不应该相等")
	}
}

func TestAddress_IsZero(t *testing.T) {
	zeroAddress, _ := NewAddress("0x0000000000000000000000000000000000000000")
	nonZeroAddress, _ := NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")

	if !zeroAddress.IsZero() {
		t.Errorf("零地址应该返回true")
	}

	if nonZeroAddress.IsZero() {
		t.Errorf("非零地址应该返回false")
	}
}

func TestAddress_ToChecksumAddress(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "标准地址",
			input: "0x742d35cc6634c0532925a3b8d4c9db96c4b4d8b6",
		},
		{
			name:  "零地址",
			input: "0x0000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := NewAddress(tt.input)
			if err != nil {
				t.Fatalf("NewAddress() 失败: %v", err)
			}

			checksum := address.ToChecksumAddress()
			// 简单验证checksum地址格式正确（包含大小写混合）
			if len(checksum) != 42 || checksum[:2] != "0x" {
				t.Errorf("ToChecksumAddress() 格式不正确: %v", checksum)
			}

			// 验证checksum地址可以被解析
			_, err = NewAddress(checksum)
			if err != nil {
				t.Errorf("ToChecksumAddress() 返回的地址无法解析: %v", err)
			}
		})
	}
}

func TestAddress_String(t *testing.T) {
	addressStr := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	address, err := NewAddress(addressStr)
	if err != nil {
		t.Fatalf("NewAddress() 失败: %v", err)
	}

	expected := "0x742d35Cc6634C0532925A3B8D4C9dB96C4B4d8B6"
	if address.String() != expected {
		t.Errorf("String() = %v, 期望 %v", address.String(), expected)
	}
}

// 基准测试
func BenchmarkNewAddress(b *testing.B) {
	addressStr := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
	for i := 0; i < b.N; i++ {
		_, _ = NewAddress(addressStr)
	}
}

func BenchmarkAddress_Equals(b *testing.B) {
	address1, _ := NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	address2, _ := NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")

	for i := 0; i < b.N; i++ {
		_ = address1.Equals(address2)
	}
}
