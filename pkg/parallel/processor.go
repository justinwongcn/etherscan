// Package parallel 提供并发处理的通用工具
package parallel

import (
	"runtime"
	"sync"
)

// Converter 定义了类型转换函数的接口
// T: 输入数据类型
// R: 输出数据类型
type Converter[T any, R any] func(T) R

// Process 是一个高性能的并发处理函数，用于高效地处理大规模数据集合
// 该函数根据输入数据量自动选择最优的处理策略：
//   - 当数据量小于阈值时，采用单线程顺序处理，避免并发带来的额外开销
//   - 当数据量大于阈值时，自动切换为多goroutine并行处理，充分利用多核性能
//
// 参数说明:
//   - items: 输入数据切片，支持任意类型T的数据集合
//   - conv: 类型转换函数，定义了将输入类型T转换为输出类型R的转换规则
//   - threshold: 并发处理阈值，用于动态决定处理策略：
//   - 当 len(items) < threshold 时使用顺序处理
//   - 当 len(items) >= threshold 时使用并发处理
//
// 返回值:
//   - []R: 返回经过转换处理后的结果切片，其长度与输入切片相同
//   - 如果输入切片为nil，则返回nil
//   - 返回的切片中元素顺序与输入切片保持一致
func Process[T any, R any](items []T, conv Converter[T, R], threshold int) []R {
	// 空值检查，避免对nil切片进行处理
	if items == nil {
		return nil
	}

	// 初始化结果切片，预分配内存以提高性能
	result := make([]R, len(items))

	// 当数据量小于阈值时，使用普通循环处理
	// 避免创建goroutine带来的开销超过并发处理带来的收益
	if len(items) < threshold {
		for i, item := range items {
			result[i] = conv(item)
		}
		return result
	}

	// 并发处理配置
	// 使用CPU核心数作为goroutine数量，避免过多的上下文切换
	workers := runtime.NumCPU()
	// 计算每个worker处理的数据块大小
	// 使用向上取整确保所有数据都被处理
	chunkSize := (len(items) + workers - 1) / workers

	// 使用WaitGroup同步所有goroutine
	var wg sync.WaitGroup
	wg.Add(workers)

	// 启动多个goroutine并行处理数据
	// 使用range创建指定数量的goroutine
	for i := range workers {
		// 计算每个worker的处理范围
		start := i * chunkSize
		// 使用min函数确保不会越界
		end := min(start+chunkSize, len(items))

		// 创建goroutine处理数据
		// 通过闭包捕获start和end确保每个goroutine处理正确的数据范围
		go func(start, end int) {
			// 确保在goroutine退出时通知WaitGroup
			defer wg.Done()
			// 处理指定范围内的数据
			for i := start; i < end; i++ {
				result[i] = conv(items[i])
			}
		}(start, end)
	}

	// 等待所有goroutine完成
	wg.Wait()
	return result
}
