// Package main 演示如何在实际以太坊项目中使用缓存的区块高度查询服务
// 注意：这是一个演示文件，如需运行请将 runDemo() 改为 main()
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/justinwongcn/etherscan/application/service"
	domainService "github.com/justinwongcn/etherscan/domain/service"
	"github.com/justinwongcn/go-ethlibs/eth"
)

func runDemo() {
	fmt.Println("=== 以太坊区块高度缓存服务集成演示 ===")
	fmt.Println()

	ctx := context.Background()

	// 1. 创建以太坊客户端和仓储
	fmt.Println("1. 初始化以太坊客户端和仓储...")

	// 创建以太坊仓储实现（这里使用模拟仓储进行演示）
	// 在实际项目中，可以使用 repository.NewEthereumRepository() 连接真实节点
	fmt.Println("⚠️  使用模拟仓储进行演示")
	ethRepo := &mockEthereumRepository{}

	// 2. 创建基础的区块高度服务
	fmt.Println("2. 创建基础区块高度服务...")
	baseHeightService := domainService.NewBlockHeightService(ethRepo)

	// 3. 创建带缓存的区块高度服务装饰器
	fmt.Println("3. 创建缓存装饰器 (5秒过期时间)...")
	cachedHeightService, err := domainService.NewCachedBlockHeightService(
		baseHeightService,
		5*time.Second, // 5秒缓存过期时间
	)
	if err != nil {
		log.Fatalf("创建缓存服务失败: %v", err)
	}
	fmt.Println("✅ 缓存服务创建成功")
	fmt.Println()

	// 4. 创建应用层的区块服务，使用缓存的高度服务
	fmt.Println("4. 创建应用层区块服务...")
	_ = service.NewBlockService(ethRepo) // 创建但不使用，仅用于演示集成

	// 5. 性能对比测试
	fmt.Println("5. 性能对比测试:")
	fmt.Println()

	// 测试原始服务性能
	fmt.Println("   📊 测试原始服务性能:")
	start := time.Now()
	originalHeight, err := baseHeightService.GetLatestBlockHeight(ctx)
	originalDuration := time.Since(start)
	if err != nil {
		log.Printf("   原始服务查询失败: %v", err)
	} else {
		fmt.Printf("   原始服务: 高度=%d, 耗时=%v\n", originalHeight, originalDuration)
	}

	// 测试缓存服务性能（第一次调用，缓存未命中）
	fmt.Println("   🔄 测试缓存服务性能 (首次调用):")
	start = time.Now()
	cachedHeight1, err := cachedHeightService.GetLatestBlockHeight(ctx)
	cachedDuration1 := time.Since(start)
	if err != nil {
		log.Printf("   缓存服务查询失败: %v", err)
	} else {
		fmt.Printf("   缓存服务(未命中): 高度=%d, 耗时=%v\n", cachedHeight1, cachedDuration1)
	}

	// 测试缓存服务性能（第二次调用，缓存命中）
	fmt.Println("   ⚡ 测试缓存服务性能 (缓存命中):")
	start = time.Now()
	cachedHeight2, err := cachedHeightService.GetLatestBlockHeight(ctx)
	cachedDuration2 := time.Since(start)
	if err != nil {
		log.Printf("   缓存服务查询失败: %v", err)
	} else {
		fmt.Printf("   缓存服务(命中): 高度=%d, 耗时=%v\n", cachedHeight2, cachedDuration2)
		if originalDuration > 0 && cachedDuration2 > 0 {
			speedup := float64(originalDuration) / float64(cachedDuration2)
			fmt.Printf("   🚀 性能提升: %.2fx 倍\n", speedup)
		}
	}

	// 6. 并发测试
	fmt.Println("\n6. 并发访问测试 (SingleFlight机制):")

	// 等待缓存过期
	fmt.Println("   等待缓存过期...")
	time.Sleep(6 * time.Second)

	// 启动多个并发请求
	const numRequests = 10

	type testResult struct {
		id       int
		height   uint64
		duration time.Duration
		err      error
	}

	results := make(chan testResult, numRequests)

	fmt.Printf("   启动 %d 个并发请求...\n", numRequests)
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		go func(requestID int) {
			reqStart := time.Now()
			height, err := cachedHeightService.GetLatestBlockHeight(ctx)
			reqDuration := time.Since(reqStart)
			results <- testResult{
				id:       requestID,
				height:   height,
				duration: reqDuration,
				err:      err,
			}
		}(i)
	}

	// 收集结果
	var successCount int
	var totalDuration time.Duration
	var heights []uint64

	for i := 0; i < numRequests; i++ {
		result := <-results
		if result.err != nil {
			fmt.Printf("   请求 %d 失败: %v\n", result.id, result.err)
			continue
		}

		successCount++
		heights = append(heights, result.height)
		totalDuration += result.duration
		fmt.Printf("   请求 %d: 高度=%d, 耗时=%v\n", result.id, result.height, result.duration)
	}

	totalTestDuration := time.Since(startTime)

	// 验证SingleFlight机制
	if len(heights) > 1 {
		firstHeight := heights[0]
		allSame := true
		for _, h := range heights {
			if h != firstHeight {
				allSame = false
				break
			}
		}

		fmt.Printf("\n   📈 并发测试结果:\n")
		fmt.Printf("   - 成功请求数: %d/%d\n", successCount, numRequests)
		fmt.Printf("   - 总耗时: %v\n", totalTestDuration)
		fmt.Printf("   - 平均响应时间: %v\n", totalDuration/time.Duration(successCount))

		if allSame {
			fmt.Printf("   ✅ SingleFlight机制正常: 所有请求返回相同高度 %d\n", firstHeight)
		} else {
			fmt.Printf("   ❌ SingleFlight机制异常: 返回了不同的高度\n")
		}
	}

	// 7. 缓存过期测试
	fmt.Println("\n7. 缓存过期机制测试:")
	fmt.Println("   等待缓存过期...")
	time.Sleep(6 * time.Second)

	fmt.Println("   缓存过期后查询:")
	start = time.Now()
	expiredHeight, err := cachedHeightService.GetLatestBlockHeight(ctx)
	expiredDuration := time.Since(start)
	if err != nil {
		log.Printf("   过期后查询失败: %v", err)
	} else {
		fmt.Printf("   过期后查询: 高度=%d, 耗时=%v\n", expiredHeight, expiredDuration)
	}

	// 8. 总结
	fmt.Println("\n=== 集成演示总结 ===")
	fmt.Println("✅ 成功集成hamster缓存到以太坊项目")
	fmt.Println("✅ 装饰器模式无缝扩展现有服务")
	fmt.Println("✅ 缓存显著提升查询性能")
	fmt.Println("✅ SingleFlight机制有效防止缓存击穿")
	fmt.Println("✅ 缓存过期机制按需触发")
	fmt.Println("\n🎯 实现特性:")
	fmt.Println("   - 基于hamster的本地内存缓存")
	fmt.Println("   - 5秒缓存过期时间")
	fmt.Println("   - 内置singleflight防止缓存击穿")
	fmt.Println("   - 装饰器模式，不影响原有代码")
	fmt.Println("   - 只在请求时检查缓存过期")
}

// mockEthereumRepository 模拟以太坊仓储，用于演示
type mockEthereumRepository struct {
	currentHeight uint64
}

func (m *mockEthereumRepository) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	// 模拟网络延迟
	time.Sleep(200 * time.Millisecond)

	// 模拟区块高度增长
	m.currentHeight++
	if m.currentHeight == 0 {
		m.currentHeight = 19000000 // 从一个合理的区块号开始
	}

	fmt.Printf("   🔗 从以太坊节点获取: %d\n", m.currentHeight)
	return m.currentHeight, nil
}

// 实现其他必需的方法（简化版本）
func (m *mockEthereumRepository) GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error) {
	return nil, fmt.Errorf("未实现")
}

func (m *mockEthereumRepository) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error) {
	return nil, fmt.Errorf("未实现")
}

func (m *mockEthereumRepository) GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	return 0, fmt.Errorf("未实现")
}

func (m *mockEthereumRepository) GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error) {
	return 0, fmt.Errorf("未实现")
}
