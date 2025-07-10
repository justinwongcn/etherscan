package valueobject

import (
	"math/big"
	"testing"
)

func TestNewGas(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "有效的Gas值（十进制）",
			input:     "21000",
			wantError: false,
			expected:  "21000",
		},
		{
			name:      "有效的Gas值（十六进制）",
			input:     "0x5208",
			wantError: false,
			expected:  "21000",
		},
		{
			name:      "零值",
			input:     "0",
			wantError: false,
			expected:  "0",
		},
		{
			name:      "大Gas值",
			input:     "8000000",
			wantError: false,
			expected:  "8000000",
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
			gas, err := NewGas(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewGas() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewGas() 返回了意外的错误: %v", err)
				return
			}

			if gas.String() != tt.expected {
				t.Errorf("NewGas() = %v, 期望 %v", gas.String(), tt.expected)
			}
		})
	}
}

func TestNewGasFromUint64(t *testing.T) {
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
			name:     "标准转账Gas",
			input:    21000,
			expected: "21000",
		},
		{
			name:     "合约调用Gas",
			input:    100000,
			expected: "100000",
		},
		{
			name:     "最大uint64",
			input:    18446744073709551615,
			expected: "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gas := NewGasFromUint64(tt.input)
			if gas.String() != tt.expected {
				t.Errorf("NewGasFromUint64() = %v, 期望 %v", gas.String(), tt.expected)
			}
		})
	}
}

func TestNewGasFromBigInt(t *testing.T) {
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
			name:     "标准转账Gas",
			input:    big.NewInt(21000),
			expected: "21000",
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
			gas, err := NewGasFromBigInt(tt.input)
			if err != nil {
				t.Errorf("NewGasFromBigInt() 返回了意外的错误: %v", err)
				return
			}
			if gas.String() != tt.expected {
				t.Errorf("NewGasFromBigInt() = %v, 期望 %v", gas.String(), tt.expected)
			}
		})
	}
}

func TestGas_Hex(t *testing.T) {
	gas := NewGasFromUint64(21000)
	expected := "0x5208"
	if gas.Hex() != expected {
		t.Errorf("Hex() = %v, 期望 %v", gas.Hex(), expected)
	}
}

func TestGas_Uint64(t *testing.T) {
	tests := []struct {
		name      string
		input     uint64
		wantError bool
	}{
		{
			name:      "有效的uint64值",
			input:     21000,
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
			gas := NewGasFromUint64(tt.input)
			result, err := gas.Uint64()

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

func TestGas_BigInt(t *testing.T) {
	input := big.NewInt(21000)
	gas, err := NewGasFromBigInt(input)
	if err != nil {
		t.Errorf("NewGasFromBigInt() 返回了意外的错误: %v", err)
		return
	}

	result := gas.BigInt()
	if result.Cmp(input) != 0 {
		t.Errorf("BigInt() = %v, 期望 %v", result, input)
	}
}

func TestGas_Add(t *testing.T) {
	gas1 := NewGasFromUint64(21000)
	gas2 := NewGasFromUint64(10000)

	result := gas1.Add(gas2)
	expected := "31000"

	if result.String() != expected {
		t.Errorf("Add() = %v, 期望 %v", result.String(), expected)
	}
}

func TestGas_Sub(t *testing.T) {
	gas1 := NewGasFromUint64(30000)
	gas2 := NewGasFromUint64(10000)

	result, err := gas1.Sub(gas2)
	if err != nil {
		t.Errorf("Sub() 返回了意外的错误: %v", err)
		return
	}

	expected := "20000"
	if result.String() != expected {
		t.Errorf("Sub() = %v, 期望 %v", result.String(), expected)
	}

	// 测试负数结果
	_, err = gas2.Sub(gas1)
	if err == nil {
		t.Errorf("Sub() 应该对负数结果返回错误")
	}
}

func TestGas_Mul(t *testing.T) {
	gas := NewGasFromUint64(21000)
	multiplier := uint64(3)

	result := gas.Mul(multiplier)
	expected := "63000"

	if result.String() != expected {
		t.Errorf("Mul() = %v, 期望 %v", result.String(), expected)
	}
}

func TestGas_Equals(t *testing.T) {
	gas1 := NewGasFromUint64(21000)
	gas2 := NewGasFromUint64(21000)
	gas3 := NewGasFromUint64(30000)

	if !gas1.Equals(gas2) {
		t.Errorf("相同的Gas值应该相等")
	}

	if gas1.Equals(gas3) {
		t.Errorf("不同的Gas值不应该相等")
	}
}

func TestGas_GreaterThan(t *testing.T) {
	gas1 := NewGasFromUint64(30000)
	gas2 := NewGasFromUint64(21000)

	if !gas1.GreaterThan(gas2) {
		t.Errorf("30000应该大于21000")
	}

	if gas2.GreaterThan(gas1) {
		t.Errorf("21000不应该大于30000")
	}
}

func TestGas_LessThan(t *testing.T) {
	gas1 := NewGasFromUint64(21000)
	gas2 := NewGasFromUint64(30000)

	if !gas1.LessThan(gas2) {
		t.Errorf("21000应该小于30000")
	}

	if gas2.LessThan(gas1) {
		t.Errorf("30000不应该小于21000")
	}
}

func TestGas_IsZero(t *testing.T) {
	zeroGas := NewGasFromUint64(0)
	nonZeroGas := NewGasFromUint64(21000)

	if !zeroGas.IsZero() {
		t.Errorf("零Gas应该返回true")
	}

	if nonZeroGas.IsZero() {
		t.Errorf("非零Gas应该返回false")
	}
}

// 基准测试
func BenchmarkNewGas(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewGas("21000")
	}
}

func BenchmarkGas_Add(b *testing.B) {
	gas1 := NewGasFromUint64(21000)
	gas2 := NewGasFromUint64(10000)

	for i := 0; i < b.N; i++ {
		_ = gas1.Add(gas2)
	}
}
