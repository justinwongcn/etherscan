# 缓存系统技术文档

本文档详细介绍了 Etherscan 项目中基于 hamster 的高性能缓存系统。

## 概述

Etherscan 缓存系统是一个基于 [hamster](https://github.com/justinwongcn/hamster) 的高性能区块链数据缓存解决方案，旨在显著提升 API 响应速度并减少对以太坊节点的访问压力。

### 核心特性

- **分层缓存策略**: 主缓存 + 二级缓存的双层架构
- **高性能表现**: 纳秒级响应时间，高并发支持
- **智能过期机制**: 基于数据特性的差异化过期策略
- **并发安全**: 内置 SingleFlight 机制防止缓存击穿
- **内存优化**: 精确的内存使用控制和监控

## 架构设计

### 整体架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP API      │    │  Cached Service  │    │  Block Service  │
│                 │───▶│                  │───▶│                 │
│ /cached/blocks/ │    │  hamster Cache   │    │  Ethereum Node │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  分层缓存策略     │
                       │                  │
                       │ 主缓存: 3分钟     │
                       │ 二级缓存: 10分钟  │
                       │ SingleFlight     │
                       └──────────────────┘
```

### 分层缓存策略

#### 主缓存 (Primary Cache)
- **用途**: 存储最新的区块数据
- **过期时间**: 3分钟
- **内存限制**: 50MB
- **算法**: ReadThroughCache
- **适用场景**: 频繁访问的最新区块

#### 二级缓存 (Secondary Cache)  
- **用途**: 存储历史区块数据
- **过期时间**: 10分钟
- **内存限制**: 100MB
- **算法**: ReadThroughCache
- **适用场景**: 历史区块数据的长期缓存

### 缓存键设计

#### 区块号查询
```
格式: block:number:{blockNumber}:fullTx:{true/false}
示例: 
- block:number:12345:fullTx:false
- block:number:latest:fullTx:true
```

#### 区块哈希查询
```
格式: block:hash:{normalizedHash}:fullTx:{true/false}
示例:
- block:hash:0xabc123...:fullTx:false
- block:hash:0xdef456...:fullTx:true
```

#### 哈希规范化规则
1. 转换为小写
2. 确保 `0x` 前缀
3. 验证长度和格式

```go
// 规范化示例
"ABC123..." → "0xabc123..."
"0xABC123..." → "0xabc123..."
"abc123..." → "0xabc123..."
```

## 核心组件

### CachedBlockService

主要的缓存服务类，实现了装饰器模式。

```go
type CachedBlockService struct {
    blockRepo          repository.BlockRepository
    primaryCache       *cache.ReadThroughService
    secondaryCache     *cache.ReadThroughService
    maxPrimaryBlocks   int
    maxSecondaryBlocks int
}
```

#### 主要方法

- `GetBlockByNumber()`: 按区块号获取区块（带缓存）
- `GetBlockByHash()`: 按区块哈希获取区块（带缓存）
- `PreloadLatestBlocks()`: 预热缓存
- `GetCacheStats()`: 获取缓存统计信息

### 配置管理

#### BlockCacheConfig

```go
type BlockCacheConfig struct {
    PrimaryExpiration   time.Duration  // 主缓存过期时间
    SecondaryExpiration time.Duration  // 二级缓存过期时间
    MaxPrimaryBlocks    int           // 主缓存最大区块数
    MaxSecondaryBlocks  int           // 二级缓存最大区块数
}
```

#### 默认配置

```go
func DefaultBlockCacheConfig() *BlockCacheConfig {
    return &BlockCacheConfig{
        PrimaryExpiration:   3 * time.Minute,
        SecondaryExpiration: 10 * time.Minute,
        MaxPrimaryBlocks:    10,
        MaxSecondaryBlocks:  10,
    }
}
```

## 性能指标

### 基准测试结果

| 测试场景 | 响应时间 | 内存使用 | 内存分配次数 |
|---------|---------|---------|-------------|
| 缓存命中 | 366.5 ns/op | 128 B/op | 3 allocs/op |
| 并发访问 | 210.4 ns/op | 128 B/op | 3 allocs/op |
| 哈希规范化 | ~400 ns/op | 128 B/op | 3 allocs/op |

### 实际性能表现

| 场景 | 原始版本 | 缓存版本 | 性能提升 |
|-----|---------|---------|---------|
| 首次请求 | 100-500ms | 100-500ms | 无提升 |
| 缓存命中 | 100-500ms | 20-50ms | 5-25倍 |
| 高并发 | 线性增长 | 近似常数 | 显著提升 |

## 使用指南

### 基本使用

```go
// 创建缓存服务
service, err := NewCachedBlockService(blockRepo, nil)
if err != nil {
    return err
}

// 查询区块
block, err := service.GetBlockByNumber(ctx, "12345", false)
if err != nil {
    return err
}
```

### 自定义配置

```go
// 自定义缓存配置
config := &BlockCacheConfig{
    PrimaryExpiration:   5 * time.Minute,  // 5分钟主缓存
    SecondaryExpiration: 30 * time.Minute, // 30分钟二级缓存
    MaxPrimaryBlocks:    20,               // 20个主缓存区块
    MaxSecondaryBlocks:  50,               // 50个二级缓存区块
}

service, err := NewCachedBlockService(blockRepo, config)
```

### 预热缓存

```go
// 预加载最新的10个区块
err := service.PreloadLatestBlocks(ctx, 10)
if err != nil {
    log.Printf("预热缓存失败: %v", err)
}
```

### 监控缓存

```go
// 获取缓存统计信息
stats, err := service.GetCacheStats(ctx)
if err != nil {
    return err
}

fmt.Printf("主缓存状态: %+v\n", stats["primary"])
fmt.Printf("二级缓存状态: %+v\n", stats["secondary"])
```

## 最佳实践

### 1. 缓存键设计

- **一致性**: 确保相同数据使用相同的缓存键
- **可读性**: 使用有意义的键名格式
- **唯一性**: 避免键冲突

### 2. 过期时间设置

- **最新数据**: 使用较短的过期时间（1-5分钟）
- **历史数据**: 使用较长的过期时间（10-60分钟）
- **静态数据**: 可以使用更长的过期时间

### 3. 内存管理

- **监控内存使用**: 定期检查缓存内存占用
- **合理设置限制**: 根据系统资源设置内存限制
- **及时清理**: 实现过期数据的及时清理

### 4. 错误处理

- **降级策略**: 缓存失败时自动降级到原始服务
- **错误日志**: 记录缓存相关的错误信息
- **监控告警**: 设置缓存异常的监控告警

### 5. 性能优化

- **预热策略**: 在系统启动时预加载热点数据
- **批量操作**: 尽可能使用批量操作减少开销
- **并发控制**: 合理控制并发访问数量

## 故障排除

### 常见问题

#### 1. 缓存未命中率高

**可能原因:**
- 过期时间设置过短
- 缓存键不一致
- 内存限制过小

**解决方案:**
- 调整过期时间配置
- 检查缓存键生成逻辑
- 增加内存限制

#### 2. 内存使用过高

**可能原因:**
- 缓存数据过大
- 过期时间过长
- 内存泄漏

**解决方案:**
- 减少缓存数据大小
- 调整过期时间
- 检查内存泄漏

#### 3. 响应时间异常

**可能原因:**
- 缓存服务异常
- 网络延迟
- 并发冲突

**解决方案:**
- 检查缓存服务状态
- 优化网络配置
- 调整并发策略

### 调试工具

#### 1. 日志分析

```go
// 启用详细日志
log.SetLevel(log.DebugLevel)

// 记录缓存操作
log.Debugf("缓存键: %s, 操作: %s", cacheKey, operation)
```

#### 2. 性能分析

```bash
# 运行性能分析
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof

# 分析结果
go tool pprof cpu.prof
go tool pprof mem.prof
```

#### 3. 监控指标

- 缓存命中率
- 响应时间分布
- 内存使用情况
- 错误率统计

## 扩展开发

### 添加新的缓存类型

1. 定义缓存键格式
2. 实现缓存逻辑
3. 添加相应测试
4. 更新文档

### 自定义缓存策略

1. 实现 `CacheStrategy` 接口
2. 配置策略参数
3. 集成到服务中
4. 验证效果

### 集成监控系统

1. 定义监控指标
2. 实现数据收集
3. 配置告警规则
4. 建立监控面板

## 参考资料

- [hamster 缓存库文档](https://github.com/justinwongcn/hamster)
- [Go 性能优化指南](https://golang.org/doc/diagnostics.html)
- [缓存设计模式](https://docs.microsoft.com/en-us/azure/architecture/patterns/cache-aside)
