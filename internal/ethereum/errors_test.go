package ethereum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorConstants 测试错误常量
func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "区块哈希格式错误",
			constant: ErrInvalidBlockHashFormat,
			expected: "invalid block hash format: must be hex string starting with 0x",
		},
		{
			name:     "区块参数错误",
			constant: ErrInvalidBlockParameter,
			expected: "invalid block parameter: %v",
		},
		{
			name:     "区块号或标签错误",
			constant: ErrInvalidBlockNumberOrTag,
			expected: "invalid block number or tag: %s",
		},
		{
			name:     "获取连接失败错误",
			constant: ErrFailedToGetConnection,
			expected: "failed to get connection: %v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

// TestErrorConstantsNotEmpty 测试错误常量不为空
func TestErrorConstantsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		constant string
	}{
		{
			name:     "ErrInvalidBlockHashFormat不为空",
			constant: ErrInvalidBlockHashFormat,
		},
		{
			name:     "ErrInvalidBlockParameter不为空",
			constant: ErrInvalidBlockParameter,
		},
		{
			name:     "ErrInvalidBlockNumberOrTag不为空",
			constant: ErrInvalidBlockNumberOrTag,
		},
		{
			name:     "ErrFailedToGetConnection不为空",
			constant: ErrFailedToGetConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.constant)
		})
	}
}

// TestErrorFormatting 测试错误格式化
func TestErrorFormatting(t *testing.T) {
	tests := []struct {
		name     string
		template string
		args     []any
		expected string
	}{
		{
			name:     "格式化区块参数错误",
			template: ErrInvalidBlockParameter,
			args:     []any{"invalid_param"},
			expected: "invalid block parameter: invalid_param",
		},
		{
			name:     "格式化区块号或标签错误",
			template: ErrInvalidBlockNumberOrTag,
			args:     []any{"invalid_tag"},
			expected: "invalid block number or tag: invalid_tag",
		},
		{
			name:     "格式化获取连接失败错误",
			template: ErrFailedToGetConnection,
			args:     []any{"connection timeout"},
			expected: "failed to get connection: connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fmt.Sprintf(tt.template, tt.args...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestErrorConstantsUniqueness 测试错误常量唯一性
func TestErrorConstantsUniqueness(t *testing.T) {
	constants := []string{
		ErrInvalidBlockHashFormat,
		ErrInvalidBlockParameter,
		ErrInvalidBlockNumberOrTag,
		ErrFailedToGetConnection,
	}

	// 检查所有常量都是唯一的
	seen := make(map[string]bool)
	for _, constant := range constants {
		assert.False(t, seen[constant], "错误常量重复: %s", constant)
		seen[constant] = true
	}
}

// TestErrorConstantsLength 测试错误常量长度
func TestErrorConstantsLength(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		minLen   int
	}{
		{
			name:     "ErrInvalidBlockHashFormat长度检查",
			constant: ErrInvalidBlockHashFormat,
			minLen:   10,
		},
		{
			name:     "ErrInvalidBlockParameter长度检查",
			constant: ErrInvalidBlockParameter,
			minLen:   10,
		},
		{
			name:     "ErrInvalidBlockNumberOrTag长度检查",
			constant: ErrInvalidBlockNumberOrTag,
			minLen:   10,
		},
		{
			name:     "ErrFailedToGetConnection长度检查",
			constant: ErrFailedToGetConnection,
			minLen:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.GreaterOrEqual(t, len(tt.constant), tt.minLen)
		})
	}
}

// TestErrorConstantsContent 测试错误常量内容
func TestErrorConstantsContent(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{
			name:     "区块哈希错误包含关键词",
			constant: ErrInvalidBlockHashFormat,
			contains: "block hash",
		},
		{
			name:     "区块参数错误包含关键词",
			constant: ErrInvalidBlockParameter,
			contains: "block parameter",
		},
		{
			name:     "区块号错误包含关键词",
			constant: ErrInvalidBlockNumberOrTag,
			contains: "block number",
		},
		{
			name:     "连接错误包含关键词",
			constant: ErrFailedToGetConnection,
			contains: "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.constant, tt.contains)
		})
	}
}

// BenchmarkErrorConstants 性能测试
func BenchmarkErrorConstants(b *testing.B) {
	b.Run("访问错误常量", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ErrInvalidBlockHashFormat
			_ = ErrInvalidBlockParameter
			_ = ErrInvalidBlockNumberOrTag
			_ = ErrFailedToGetConnection
		}
	})
	
	b.Run("格式化错误消息", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf(ErrInvalidBlockParameter, "test")
			_ = fmt.Sprintf(ErrInvalidBlockNumberOrTag, "test")
			_ = fmt.Sprintf(ErrFailedToGetConnection, "test")
		}
	})
}
