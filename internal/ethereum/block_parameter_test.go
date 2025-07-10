package ethereum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseBlockParameter 测试区块参数解析
func TestParseBlockParameter(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		// 测试特殊标识符
		{
			name:           "latest标识符",
			input:          "latest",
			expectedOutput: "latest",
			expectError:    false,
		},
		{
			name:           "earliest标识符",
			input:          "earliest",
			expectedOutput: "earliest",
			expectError:    false,
		},
		{
			name:           "pending标识符",
			input:          "pending",
			expectedOutput: "pending",
			expectError:    false,
		},
		{
			name:           "大写latest标识符",
			input:          "LATEST",
			expectedOutput: "latest",
			expectError:    false,
		},
		{
			name:           "混合大小写earliest标识符",
			input:          "Earliest",
			expectedOutput: "earliest",
			expectError:    false,
		},

		// 测试区块哈希
		{
			name:           "有效区块哈希",
			input:          "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expectedOutput: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expectError:    false,
		},
		{
			name:           "短区块哈希",
			input:          "0x1234",
			expectedOutput: "0x1234",
			expectError:    false,
		},

		// 测试区块号
		{
			name:           "区块号0",
			input:          "0",
			expectedOutput: "0x0",
			expectError:    false,
		},
		{
			name:           "区块号1",
			input:          "1",
			expectedOutput: "0x1",
			expectError:    false,
		},
		{
			name:           "大区块号",
			input:          "1000000",
			expectedOutput: "0xf4240",
			expectError:    false,
		},
		{
			name:           "十六进制区块号",
			input:          "255",
			expectedOutput: "0xff",
			expectError:    false,
		},

		// 测试错误情况
		{
			name:        "无效字符串",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "包含字母的数字",
			input:       "123abc",
			expectError: true,
		},
		{
			name:        "空字符串",
			input:       "",
			expectError: true,
		},
		{
			name:           "负数",
			input:          "-1",
			expectedOutput: "0x-1",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBlockParameter(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, result)
			}
		})
	}
}

// TestParseBlockParameter_EdgeCases 测试边界情况
func TestParseBlockParameter_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "只有0x前缀",
			input:          "0x",
			expectedOutput: "0x",
			expectError:    false,
		},
		{
			name:           "很大的数字",
			input:          "999999999999999999999",
			expectedOutput: "0x3635c9adc5de9fffff",
			expectError:    false,
		},
		{
			name:        "包含空格的字符串",
			input:       " latest ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBlockParameter(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, result)
			}
		})
	}
}

// BenchmarkParseBlockParameter 性能测试
func BenchmarkParseBlockParameter(b *testing.B) {
	testCases := []string{
		"latest",
		"1000000",
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = ParseBlockParameter(tc)
			}
		})
	}
}
