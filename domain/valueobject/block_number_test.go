package valueobject

import (
	"math/big"
	"testing"
)

func TestNewBlockNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "有效的十进制区块号",
			input:     "12345",
			wantError: false,
			expected:  "12345",
		},
		{
			name:      "有效的十六进制区块号",
			input:     "0x3039",
			wantError: false,
			expected:  "12345",
		},
		{
			name:      "零区块号",
			input:     "0",
			wantError: false,
			expected:  "0",
		},
		{
			name:      "latest标签",
			input:     "latest",
			wantError: false,
			expected:  "latest",
		},
		{
			name:      "earliest标签",
			input:     "earliest",
			wantError: false,
			expected:  "earliest",
		},
		{
			name:      "pending标签",
			input:     "pending",
			wantError: false,
			expected:  "pending",
		},
		{
			name:      "空字符串",
			input:     "",
			wantError: true,
		},
		{
			name:      "无效字符",
			input:     "abc",
			wantError: true,
		},
		{
			name:      "负数",
			input:     "-1",
			wantError: true,
		},
		{
			name:      "无效的十六进制",
			input:     "0xgg",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumber(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewBlockNumber() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewBlockNumber() 返回了意外的错误: %v", err)
				return
			}

			if blockNumber.String() != tt.expected {
				t.Errorf("NewBlockNumber() = %v, 期望 %v", blockNumber.String(), tt.expected)
			}
		})
	}
}

func TestNewBlockNumberFromUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "零值",
			input:    0,
			expected: "0",
		},
		{
			name:     "普通区块号",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "大区块号",
			input:    18446744073709551615,
			expected: "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber := NewBlockNumberFromUint64(tt.input)
			if blockNumber.String() != tt.expected {
				t.Errorf("NewBlockNumberFromUint64() = %v, 期望 %v", blockNumber.String(), tt.expected)
			}
		})
	}
}

func TestNewBlockNumberFromBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{
			name:     "零值",
			input:    big.NewInt(0),
			expected: "0",
		},
		{
			name:     "普通区块号",
			input:    big.NewInt(12345),
			expected: "12345",
		},
		{
			name: "大数值",
			input: func() *big.Int {
				val, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
				return val
			}(),
			expected: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumberFromBigInt(tt.input)
			if err != nil {
				t.Errorf("NewBlockNumberFromBigInt() 返回了意外的错误: %v", err)
				return
			}
			if blockNumber.String() != tt.expected {
				t.Errorf("NewBlockNumberFromBigInt() = %v, 期望 %v", blockNumber.String(), tt.expected)
			}
		})
	}
}

func TestNewLatestBlockNumber(t *testing.T) {
	blockNumber := NewLatestBlockNumber()
	if blockNumber.String() != "latest" {
		t.Errorf("NewLatestBlockNumber() = %v, 期望 latest", blockNumber.String())
	}

	if !blockNumber.IsTag() {
		t.Errorf("latest区块号应该是标签")
	}

	if blockNumber.IsNumber() {
		t.Errorf("latest区块号不应该是数字")
	}
}

func TestNewEarliestBlockNumber(t *testing.T) {
	blockNumber := NewEarliestBlockNumber()
	if blockNumber.String() != "earliest" {
		t.Errorf("NewEarliestBlockNumber() = %v, 期望 earliest", blockNumber.String())
	}

	if !blockNumber.IsTag() {
		t.Errorf("earliest区块号应该是标签")
	}
}

func TestNewPendingBlockNumber(t *testing.T) {
	blockNumber := NewPendingBlockNumber()
	if blockNumber.String() != "pending" {
		t.Errorf("NewPendingBlockNumber() = %v, 期望 pending", blockNumber.String())
	}

	if !blockNumber.IsTag() {
		t.Errorf("pending区块号应该是标签")
	}
}

func TestBlockNumber_Hex(t *testing.T) {
	blockNumber := NewBlockNumberFromUint64(12345)
	expected := "0x3039"
	if blockNumber.Hex() != expected {
		t.Errorf("Hex() = %v, 期望 %v", blockNumber.Hex(), expected)
	}
}

func TestBlockNumber_IsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "数字区块号",
			input:    "12345",
			expected: true,
		},
		{
			name:     "latest标签",
			input:    "latest",
			expected: false,
		},
		{
			name:     "earliest标签",
			input:    "earliest",
			expected: false,
		},
		{
			name:     "pending标签",
			input:    "pending",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumber(tt.input)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			if blockNumber.IsNumber() != tt.expected {
				t.Errorf("IsNumber() = %v, 期望 %v", blockNumber.IsNumber(), tt.expected)
			}
		})
	}
}

func TestBlockNumber_IsTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "数字区块号",
			input:    "12345",
			expected: false,
		},
		{
			name:     "latest标签",
			input:    "latest",
			expected: true,
		},
		{
			name:     "earliest标签",
			input:    "earliest",
			expected: true,
		},
		{
			name:     "pending标签",
			input:    "pending",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumber(tt.input)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			if blockNumber.IsTag() != tt.expected {
				t.Errorf("IsTag() = %v, 期望 %v", blockNumber.IsTag(), tt.expected)
			}
		})
	}
}

func TestBlockNumber_Uint64(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  uint64
	}{
		{
			name:      "有效的数字区块号",
			input:     "12345",
			wantError: false,
			expected:  12345,
		},
		{
			name:      "零区块号",
			input:     "0",
			wantError: false,
			expected:  0,
		},
		{
			name:      "latest标签",
			input:     "latest",
			wantError: true,
		},
		{
			name:      "earliest标签",
			input:     "earliest",
			wantError: true,
		},
		{
			name:      "pending标签",
			input:     "pending",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumber(tt.input)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			result, err := blockNumber.Uint64()

			if tt.wantError {
				if err == nil {
					t.Errorf("Uint64() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("Uint64() 返回了意外的错误: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Uint64() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestBlockNumber_BigInt(t *testing.T) {
	blockNumber := NewBlockNumberFromUint64(12345)
	result, err := blockNumber.BigInt()
	if err != nil {
		t.Errorf("BigInt() 返回了意外的错误: %v", err)
		return
	}

	expected := big.NewInt(12345)
	if result.Cmp(expected) != 0 {
		t.Errorf("BigInt() = %v, 期望 %v", result, expected)
	}
}

func TestBlockNumber_Tag(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "latest标签",
			input:     "latest",
			wantError: false,
			expected:  "latest",
		},
		{
			name:      "earliest标签",
			input:     "earliest",
			wantError: false,
			expected:  "earliest",
		},
		{
			name:      "pending标签",
			input:     "pending",
			wantError: false,
			expected:  "pending",
		},
		{
			name:      "数字区块号",
			input:     "12345",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockNumber, err := NewBlockNumber(tt.input)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			result, err := blockNumber.Tag()

			if tt.wantError {
				if err == nil {
					t.Errorf("Tag() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("Tag() 返回了意外的错误: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Tag() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestBlockNumber_Equals(t *testing.T) {
	tests := []struct {
		name     string
		input1   string
		input2   string
		expected bool
	}{
		{
			name:     "相同的数字区块号",
			input1:   "12345",
			input2:   "12345",
			expected: true,
		},
		{
			name:     "不同的数字区块号",
			input1:   "12345",
			input2:   "54321",
			expected: false,
		},
		{
			name:     "相同的标签",
			input1:   "latest",
			input2:   "latest",
			expected: true,
		},
		{
			name:     "不同的标签",
			input1:   "latest",
			input2:   "earliest",
			expected: false,
		},
		{
			name:     "数字与标签",
			input1:   "12345",
			input2:   "latest",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bn1, err := NewBlockNumber(tt.input1)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			bn2, err := NewBlockNumber(tt.input2)
			if err != nil {
				t.Fatalf("NewBlockNumber() 失败: %v", err)
			}

			if bn1.Equals(bn2) != tt.expected {
				t.Errorf("Equals() = %v, 期望 %v", bn1.Equals(bn2), tt.expected)
			}
		})
	}
}
