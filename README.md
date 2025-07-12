# Etherscan

一个高性能的以太坊区块链交互客户端，提供智能的连接池管理、节点交互功能和基于 hamster 的高级缓存系统。

## 功能特性

### 🚀 核心功能
- **智能连接池管理**
  - 自动管理连接生命周期
  - 支持并发连接控制
  - 空闲连接自动清理
- **健康检查机制**
  - 定期检查连接状态
  - 自动重连和故障转移
- **高性能设计**
  - 连接复用
  - 并发请求处理
  - 资源使用优化

### ⚡ 高级缓存系统
- **基于 hamster 的分层缓存**
  - 主缓存：3分钟过期，最新区块数据
  - 二级缓存：10分钟过期，历史区块数据
  - 内置 SingleFlight 机制防止缓存击穿
- **智能缓存策略**
  - 按区块号和哈希查询的独立缓存
  - 支持 fullTx 参数的不同缓存键
  - 自动哈希格式规范化
- **高性能表现**
  - 缓存命中：366.5 ns/op
  - 并发访问：210.4 ns/op
  - 内存优化：128 B/op，3 allocs/op

## 快速开始

### 安装

```bash
go get github.com/justinwongcn/etherscan
```

### 启动服务

```bash
# 克隆项目
git clone https://github.com/justinwongcn/etherscan.git
cd etherscan

# 启动服务
go run cmd/main.go
# Server is running on :8080
```

### 配置

在 `config.yaml` 文件中配置以太坊节点URL：

```yaml
ethereum:
    node_url: "wss://ethereum.callstaticrpc.com"
```

## API 端点

### 🔥 缓存版本 API（推荐）

#### 区块高度查询
```bash
# 获取最新区块高度（缓存版本）
GET /cached/blocks/height/latest
```

#### 区块查询
```bash
# 按区块号查询（缓存版本）
GET /cached/blocks/{number}?fullTx={true/false}

# 按区块哈希查询（缓存版本）
GET /cached/blocks/hash/{hash}?fullTx={true/false}
```

### 📊 原始版本 API

#### 区块相关
```bash
# 获取最新区块高度
GET /blocks/height/latest

# 按区块号获取区块
GET /blocks/{number}

# 获取区块交易数量
GET /blocks/{number}/transactions/count
```

#### 交易相关
```bash
# 按哈希获取交易
GET /transactions/{hash}

# 按索引获取交易
GET /blocks/transactions/{index}

# 发送原始交易
POST /transactions

# 获取交易数量
GET /accounts/{address}/transactions/count

# 获取交易收据
GET /transactions/{hash}/receipt
```

## 使用示例

### HTTP API 调用

```bash
# 获取最新区块高度（缓存版本，推荐）
curl "http://localhost:8080/cached/blocks/height/latest"
# 返回：{"height":"22898143"}

# 获取区块信息（缓存版本）
curl "http://localhost:8080/cached/blocks/latest?fullTx=false"

# 按哈希查询区块（支持多种格式）
curl "http://localhost:8080/cached/blocks/hash/0xabc123...?fullTx=false"
curl "http://localhost:8080/cached/blocks/hash/ABC123...?fullTx=false"  # 自动规范化
```

### 性能对比

```bash
# 原始版本（无缓存）
time curl "http://localhost:8080/blocks/latest?fullTx=false"
# 响应时间：~100-500ms

# 缓存版本（第二次请求）
time curl "http://localhost:8080/cached/blocks/latest?fullTx=false"
# 响应时间：~20ms（提升 5-25 倍）
```

### Go 代码示例

```go
package main

import (
    "context"
    "github.com/justinwongcn/etherscan/internal/ethereum"
    "github.com/justinwongcn/etherscan/domain/service"
)

func main() {
    // 创建客户端实例
    opts := ethereum.DefaultClientOptions()
    client := ethereum.NewClient("wss://ethereum.callstaticrpc.com", opts)

    // 创建缓存区块服务
    cachedService, err := service.NewCachedBlockService(blockRepo, nil)
    if err != nil {
        panic(err)
    }

    // 使用缓存服务
    ctx := context.Background()
    block, err := cachedService.GetBlockByNumber(ctx, "latest", false)
    if err != nil {
        panic(err)
    }

    fmt.Printf("最新区块号: %s\n", block.Number.String())
}
```

## 配置选项

### 客户端配置

可以通过 `ClientOptions` 自定义客户端行为：

```go
opts := &ethereum.ClientOptions{
    MaxConns: 100,         // 最大并发连接数
    IdleTimeout: time.Minute, // 空闲超时时间
    HealthCheck: true,     // 启用健康检查
    MaxIdleConns: 10,      // 最大空闲连接数
}
```

### 缓存配置

可以通过 `BlockCacheConfig` 自定义缓存行为：

```go
config := &service.BlockCacheConfig{
    PrimaryExpiration:   3 * time.Minute,  // 主缓存过期时间
    SecondaryExpiration: 10 * time.Minute, // 二级缓存过期时间
    MaxPrimaryBlocks:    10,               // 主缓存最大区块数
    MaxSecondaryBlocks:  10,               // 二级缓存最大区块数
}

// 创建自定义缓存服务
cachedService, err := service.NewCachedBlockService(blockRepo, config)
```

## 测试

### 运行单元测试

```bash
# 运行所有测试
go test ./...

# 运行缓存相关测试
go test ./domain/service -v -run TestCachedBlockService

# 运行集成测试
go test ./domain/service -v -run TestCachedBlockService_Integration
```

### 运行性能基准测试

```bash
# 运行所有基准测试
go test ./domain/service -bench=BenchmarkCachedBlockService -benchmem

# 缓存命中性能测试
go test ./domain/service -bench=BenchmarkCachedBlockService_CacheHit -benchmem
# 结果：366.5 ns/op, 128 B/op, 3 allocs/op

# 并发访问性能测试
go test ./domain/service -bench=BenchmarkCachedBlockService_ConcurrentAccess -benchmem
# 结果：210.4 ns/op, 128 B/op, 3 allocs/op
```

### HTTP 集成测试

项目提供了完整的 HTTP 测试文件：

```bash
# 区块高度缓存测试
# 使用 cmd/cached_height_test.http

# 区块查询缓存测试
# 使用 cmd/cached_block_test.http
```

## 架构设计

### 项目结构

```
etherscan/
├── api/                    # HTTP API 层
│   ├── handler/           # 请求处理器
│   └── routes/            # 路由配置
├── application/           # 应用服务层
│   └── service/          # 应用服务
├── domain/               # 领域层
│   ├── entity/           # 领域实体
│   ├── repository/       # 仓储接口
│   └── service/          # 领域服务（包含缓存服务）
├── infrastructure/       # 基础设施层
│   ├── ethereum/         # 以太坊客户端
│   └── repository/       # 仓储实现
├── cmd/                  # 应用入口和测试文件
└── config/               # 配置文件
```

### 缓存架构

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

### 技术栈

- **Web 框架**: [ant](https://github.com/justinwongcn/ant) - 轻量级 HTTP 框架
- **缓存引擎**: [hamster](https://github.com/justinwongcn/hamster) - 高性能缓存库
- **以太坊库**: [go-ethlibs](https://github.com/justinwongcn/go-ethlibs) - 以太坊类型库
- **测试框架**: [testify](https://github.com/stretchr/testify) - 测试断言库
- **架构模式**: DDD (领域驱动设计) + 分层架构

## 性能指标

### 缓存性能
- **缓存命中**: 366.5 ns/op
- **并发访问**: 210.4 ns/op
- **内存使用**: 128 B/op, 3 allocs/op
- **响应时间提升**: 5-25 倍

### 系统容量
- **主缓存**: 50MB 内存限制
- **二级缓存**: 100MB 内存限制
- **并发支持**: 高并发访问优化
- **过期策略**: 智能分层过期

## 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 使用的开源中间件

- [hamster](https://github.com/justinwongcn/hamster) - 提供高性能缓存支持
- [ant](https://github.com/justinwongcn/ant) - 提供轻量级 HTTP 框架
- [go-ethlibs](https://github.com/justinwongcn/go-ethlibs) - 提供以太坊类型支持