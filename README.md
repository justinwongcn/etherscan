# Etherscan

一个高性能的以太坊区块链交互客户端，提供智能的连接池管理和节点交互功能。

## 功能特性

- 智能连接池管理
  - 自动管理连接生命周期
  - 支持并发连接控制
  - 空闲连接自动清理
- 健康检查机制
  - 定期检查连接状态
  - 自动重连和故障转移
- 高性能设计
  - 连接复用
  - 并发请求处理
  - 资源使用优化

## 安装

```bash
go get github.com/your-username/etherscan
```

## 配置

在 `config.yaml` 文件中配置以太坊节点URL：

```yaml
ethereum:
    node_url: "wss://ethereum.callstaticrpc.com"
```

## 使用示例

```go
package main

import (
    "context"
    "github.com/your-username/etherscan/internal/ethereum"
)

func main() {
    // 创建客户端实例
    opts := ethereum.DefaultClientOptions()
    client := ethereum.NewClient("wss://ethereum.callstaticrpc.com", opts)

    // 使用客户端进行操作
    ctx := context.Background()
    // ...
}
```

## 配置选项

可以通过 `ClientOptions` 自定义客户端行为：

```go
opts := &ethereum.ClientOptions{
    MaxConns: 100,         // 最大并发连接数
    IdleTimeout: time.Minute, // 空闲超时时间
    HealthCheck: true,     // 启用健康检查
    MaxIdleConns: 10,      // 最大空闲连接数
}
```

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件