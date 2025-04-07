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

// Process 通用的并发处理函数，用于并行处理大量数据
// 当数据量小于阈值时使用普通循环，大于阈值时使用goroutine并发处理
// 参数:
//   - items: 需要处理的数据切片
//   - conv: 类型转换函数
//   - threshold: 并发处理的阈值，当数据量小于此值时使用普通循环
//
// 返回:
//   - []R: 处理后的结果切片
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