package valueobject

import (
	"math/big"
	"testing"
)

func TestNewWei(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "有效的Wei值（十进制）",
			input:     "1000000000000000000",
			wantError: false,
			expected:  "1000000000000000000",
		},
		{
			name:      "有效的Wei值（十六进制）",
			input:     "0xde0b6b3a7640000",
			wantError: false,
			expected:  "1000000000000000000",
		},
		{
			name:      "零值",
			input:     "0",
			wantError: false,
			expected:  "0",
		},
		{
			name:      "大数值",
			input:     "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			wantError: false,
			expected:  "115792089237316195423570985008687907853269984665640564039457584007913129639935",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei, err := NewWei(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewWei() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewWei() 返回了意外的错误: %v", err)
				return
			}

			if wei.String() != tt.expected {
				t.Errorf("NewWei() = %v, 期望 %v", wei.String(), tt.expected)
			}
		})
	}
}

func TestNewWeiFromUint64(t *testing.T) {
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
			name:     "1 ETH",
			input:    1000000000000000000,
			expected: "1000000000000000000",
		},
		{
			name:     "最大uint64",
			input:    18446744073709551615,
			expected: "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei := NewWeiFromUint64(tt.input)
			if wei.String() != tt.expected {
				t.Errorf("NewWeiFromUint64() = %v, 期望 %v", wei.String(), tt.expected)
			}
		})
	}
}

func TestNewWeiFromBigInt(t *testing.T) {
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
			name:     "1 ETH",
			input:    big.NewInt(1000000000000000000),
			expected: "1000000000000000000",
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
			wei, err := NewWeiFromBigInt(tt.input)
			if err != nil {
				t.Errorf("NewWeiFromBigInt() 返回了意外的错误: %v", err)
				return
			}
			if wei.String() != tt.expected {
				t.Errorf("NewWeiFromBigInt() = %v, 期望 %v", wei.String(), tt.expected)
			}
		})
	}
}

func TestNewWeiFromEther(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "1 ETH",
			input:     "1",
			wantError: false,
			expected:  "1000000000000000000",
		},
		{
			name:      "5 ETH",
			input:     "5",
			wantError: false,
			expected:  "5000000000000000000",
		},
		{
			name:      "零值",
			input:     "0",
			wantError: false,
			expected:  "0",
		},
		{
			name:      "小数输入（当前不支持）",
			input:     "0.5",
			wantError: true,
		},
		{
			name:      "无效输入",
			input:     "abc",
			wantError: true,
		},
		{
			name:      "负数",
			input:     "-1",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei, err := NewWeiFromEther(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewWeiFromEther() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewWeiFromEther() 返回了意外的错误: %v", err)
				return
			}

			if wei.String() != tt.expected {
				t.Errorf("NewWeiFromEther() = %v, 期望 %v", wei.String(), tt.expected)
			}
		})
	}
}

func TestWei_Hex(t *testing.T) {
	wei := NewWeiFromUint64(1000000000000000000)
	expected := "0xde0b6b3a7640000"
	if wei.Hex() != expected {
		t.Errorf("Hex() = %v, 期望 %v", wei.Hex(), expected)
	}
}

func TestWei_Uint64(t *testing.T) {
	tests := []struct {
		name      string
		input     uint64
		wantError bool
	}{
		{
			name:      "有效的uint64值",
			input:     1000000000000000000,
			wantError: false,
		},
		{
			name:      "零值",
			input:     0,
			wantError: false,
		},
		{
			name:      "最大uint64",
			input:     18446744073709551615,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei := NewWeiFromUint64(tt.input)
			result, err := wei.Uint64()

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

			if result != tt.input {
				t.Errorf("Uint64() = %v, 期望 %v", result, tt.input)
			}
		})
	}
}

func TestWei_BigInt(t *testing.T) {
	input := big.NewInt(1000000000000000000)
	wei, err := NewWeiFromBigInt(input)
	if err != nil {
		t.Errorf("NewWeiFromBigInt() 返回了意外的错误: %v", err)
		return
	}

	result := wei.BigInt()
	if result.Cmp(input) != 0 {
		t.Errorf("BigInt() = %v, 期望 %v", result, input)
	}
}

func TestWei_ToEther(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "1 ETH",
			input:    1000000000000000000,
			expected: "1",
		},
		{
			name:     "2 ETH",
			input:    2000000000000000000,
			expected: "2",
		},
		{
			name:     "0.5 ETH（只返回整数部分）",
			input:    500000000000000000,
			expected: "0",
		},
		{
			name:     "零值",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei := NewWeiFromUint64(tt.input)
			result := wei.ToEther()
			if result != tt.expected {
				t.Errorf("ToEther() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestWei_ToGwei(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{
			name:     "1 Gwei",
			input:    1000000000,
			expected: "1",
		},
		{
			name:     "20 Gwei",
			input:    20000000000,
			expected: "20",
		},
		{
			name:     "零值",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei := NewWeiFromUint64(tt.input)
			result := wei.ToGwei()
			if result != tt.expected {
				t.Errorf("ToGwei() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestWei_Add(t *testing.T) {
	wei1 := NewWeiFromUint64(1000000000000000000) // 1 ETH
	wei2 := NewWeiFromUint64(500000000000000000)  // 0.5 ETH

	result := wei1.Add(wei2)
	expected := "1500000000000000000" // 1.5 ETH

	if result.String() != expected {
		t.Errorf("Add() = %v, 期望 %v", result.String(), expected)
	}
}

func TestWei_Sub(t *testing.T) {
	wei1 := NewWeiFromUint64(1000000000000000000) // 1 ETH
	wei2 := NewWeiFromUint64(300000000000000000)  // 0.3 ETH

	result, err := wei1.Sub(wei2)
	if err != nil {
		t.Errorf("Sub() 返回了意外的错误: %v", err)
		return
	}

	expected := "700000000000000000" // 0.7 ETH
	if result.String() != expected {
		t.Errorf("Sub() = %v, 期望 %v", result.String(), expected)
	}
}

func TestWei_Mul(t *testing.T) {
	wei := NewWeiFromUint64(1000000000000000000) // 1 ETH
	multiplier := uint64(3)

	result := wei.Mul(multiplier)
	expected := "3000000000000000000" // 3 ETH

	if result.String() != expected {
		t.Errorf("Mul() = %v, 期望 %v", result.String(), expected)
	}
}

func TestWei_Div(t *testing.T) {
	wei := NewWeiFromUint64(1000000000000000000) // 1 ETH
	divisor := uint64(2)

	result, err := wei.Div(divisor)
	if err != nil {
		t.Errorf("Div() 返回了意外的错误: %v", err)
		return
	}

	expected := "500000000000000000" // 0.5 ETH
	if result.String() != expected {
		t.Errorf("Div() = %v, 期望 %v", result.String(), expected)
	}
}

func TestWei_Equals(t *testing.T) {
	wei1 := NewWeiFromUint64(1000000000000000000)
	wei2 := NewWeiFromUint64(1000000000000000000)
	wei3 := NewWeiFromUint64(2000000000000000000)

	if !wei1.Equals(wei2) {
		t.Errorf("相同的Wei值应该相等")
	}

	if wei1.Equals(wei3) {
		t.Errorf("不同的Wei值不应该相等")
	}
}

func TestWei_GreaterThan(t *testing.T) {
	wei1 := NewWeiFromUint64(2000000000000000000) // 2 ETH
	wei2 := NewWeiFromUint64(1000000000000000000) // 1 ETH

	if !wei1.GreaterThan(wei2) {
		t.Errorf("2 ETH应该大于1 ETH")
	}

	if wei2.GreaterThan(wei1) {
		t.Errorf("1 ETH不应该大于2 ETH")
	}
}

func TestWei_LessThan(t *testing.T) {
	wei1 := NewWeiFromUint64(1000000000000000000) // 1 ETH
	wei2 := NewWeiFromUint64(2000000000000000000) // 2 ETH

	if !wei1.LessThan(wei2) {
		t.Errorf("1 ETH应该小于2 ETH")
	}

	if wei2.LessThan(wei1) {
		t.Errorf("2 ETH不应该小于1 ETH")
	}
}

func TestWei_IsZero(t *testing.T) {
	zeroWei := NewWeiFromUint64(0)
	nonZeroWei := NewWeiFromUint64(1000000000000000000)

	if !zeroWei.IsZero() {
		t.Errorf("零Wei应该返回true")
	}

	if nonZeroWei.IsZero() {
		t.Errorf("非零Wei应该返回false")
	}
}
