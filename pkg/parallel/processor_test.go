package parallel

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProcess_NilInput 测试nil输入
func TestProcess_NilInput(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
	}{
		{
			name:      "nil输入",
			items:     nil,
			threshold: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(x int) string {
				return strconv.Itoa(x)
			}
			
			result := Process(tt.items, converter, tt.threshold)
			assert.Nil(t, result)
		})
	}
}

// TestProcess_EmptyInput 测试空输入
func TestProcess_EmptyInput(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
	}{
		{
			name:      "空切片输入",
			items:     []int{},
			threshold: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(x int) string {
				return strconv.Itoa(x)
			}
			
			result := Process(tt.items, converter, tt.threshold)
			assert.NotNil(t, result)
			assert.Equal(t, 0, len(result))
		})
	}
}

// TestProcess_SequentialProcessing 测试顺序处理（小于阈值）
func TestProcess_SequentialProcessing(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
		expected  []string
	}{
		{
			name:      "单个元素",
			items:     []int{1},
			threshold: 10,
			expected:  []string{"1"},
		},
		{
			name:      "小于阈值的多个元素",
			items:     []int{1, 2, 3, 4, 5},
			threshold: 10,
			expected:  []string{"1", "2", "3", "4", "5"},
		},
		{
			name:      "等于阈值减一",
			items:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			threshold: 10,
			expected:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(x int) string {
				return strconv.Itoa(x)
			}
			
			result := Process(tt.items, converter, tt.threshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProcess_ParallelProcessing 测试并行处理（大于等于阈值）
func TestProcess_ParallelProcessing(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
		expected  []string
	}{
		{
			name:      "等于阈值",
			items:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			threshold: 10,
			expected:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name:      "大于阈值",
			items:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			threshold: 10,
			expected:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"},
		},
		{
			name:      "大量数据",
			items:     make([]int, 100),
			threshold: 10,
			expected:  make([]string, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为大量数据测试初始化数据
			if tt.name == "大量数据" {
				for i := 0; i < 100; i++ {
					tt.items[i] = i + 1
					tt.expected[i] = strconv.Itoa(i + 1)
				}
			}
			
			converter := func(x int) string {
				return strconv.Itoa(x)
			}
			
			result := Process(tt.items, converter, tt.threshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProcess_DifferentTypes 测试不同类型的转换
func TestProcess_DifferentTypes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "字符串到整数",
		},
		{
			name: "整数到浮点数",
		},
		{
			name: "布尔值到字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "字符串到整数":
				items := []string{"1", "2", "3", "4", "5"}
				converter := func(s string) int {
					val, _ := strconv.Atoi(s)
					return val
				}
				result := Process(items, converter, 3)
				expected := []int{1, 2, 3, 4, 5}
				assert.Equal(t, expected, result)
				
			case "整数到浮点数":
				items := []int{1, 2, 3, 4, 5}
				converter := func(x int) float64 {
					return float64(x) * 1.5
				}
				result := Process(items, converter, 3)
				expected := []float64{1.5, 3.0, 4.5, 6.0, 7.5}
				assert.Equal(t, expected, result)
				
			case "布尔值到字符串":
				items := []bool{true, false, true, false}
				converter := func(b bool) string {
					if b {
						return "yes"
					}
					return "no"
				}
				result := Process(items, converter, 2)
				expected := []string{"yes", "no", "yes", "no"}
				assert.Equal(t, expected, result)
			}
		})
	}
}

// TestProcess_EdgeCases 测试边界情况
func TestProcess_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
		expected  []int
	}{
		{
			name:      "阈值为0",
			items:     []int{1, 2, 3},
			threshold: 0,
			expected:  []int{2, 4, 6},
		},
		{
			name:      "阈值为1",
			items:     []int{1, 2, 3},
			threshold: 1,
			expected:  []int{2, 4, 6},
		},
		{
			name:      "阈值很大",
			items:     []int{1, 2, 3},
			threshold: 1000,
			expected:  []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(x int) int {
				return x * 2
			}
			
			result := Process(tt.items, converter, tt.threshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProcess_OrderPreservation 测试顺序保持
func TestProcess_OrderPreservation(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		threshold int
	}{
		{
			name:      "顺序处理保持顺序",
			items:     []int{5, 3, 8, 1, 9, 2, 7, 4, 6},
			threshold: 20,
		},
		{
			name:      "并行处理保持顺序",
			items:     []int{5, 3, 8, 1, 9, 2, 7, 4, 6, 10, 11, 12, 13, 14, 15},
			threshold: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(x int) int {
				return x * 10
			}
			
			result := Process(tt.items, converter, tt.threshold)
			
			// 验证长度
			assert.Equal(t, len(tt.items), len(result))
			
			// 验证顺序保持
			for i, item := range tt.items {
				assert.Equal(t, item*10, result[i])
			}
		})
	}
}

// BenchmarkProcess_Sequential 顺序处理性能测试
func BenchmarkProcess_Sequential(b *testing.B) {
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}
	
	converter := func(x int) int {
		return x * 2
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Process(items, converter, 2000) // 阈值大于数据量，强制顺序处理
	}
}

// BenchmarkProcess_Parallel 并行处理性能测试
func BenchmarkProcess_Parallel(b *testing.B) {
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}
	
	converter := func(x int) int {
		return x * 2
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Process(items, converter, 100) // 阈值小于数据量，强制并行处理
	}
}

// BenchmarkProcess_LargeDataset 大数据集性能测试
func BenchmarkProcess_LargeDataset(b *testing.B) {
	items := make([]int, 10000)
	for i := range items {
		items[i] = i
	}
	
	converter := func(x int) string {
		return strconv.Itoa(x * x)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Process(items, converter, 1000)
	}
}
