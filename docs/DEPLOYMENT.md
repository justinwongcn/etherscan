# 部署和运维文档

本文档提供了 Etherscan 项目的部署指南、运维建议和监控方案。

## 环境要求

### 系统要求

- **操作系统**: Linux (推荐 Ubuntu 20.04+) / macOS / Windows
- **CPU**: 2核心以上
- **内存**: 4GB 以上（推荐 8GB）
- **存储**: 10GB 可用空间
- **网络**: 稳定的互联网连接

### 软件依赖

- **Go**: 1.19+ (推荐 1.21+)
- **Git**: 用于代码管理
- **以太坊节点**: 或者可访问的以太坊 RPC 端点

## 部署方式

### 1. 源码部署

#### 克隆项目

```bash
git clone https://github.com/justinwongcn/etherscan.git
cd etherscan
```

#### 配置环境

```bash
# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 编辑配置文件
vim config/config.yaml
```

#### 配置示例

```yaml
# config/config.yaml
server:
  port: 8080
  host: "0.0.0.0"

ethereum:
  node_url: "wss://ethereum.callstaticrpc.com"
  timeout: 30s
  max_connections: 100

cache:
  primary_expiration: "3m"
  secondary_expiration: "10m"
  max_primary_blocks: 10
  max_secondary_blocks: 10
  max_memory_mb: 150

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

#### 构建和运行

```bash
# 安装依赖
go mod download

# 构建项目
go build -o etherscan cmd/main.go

# 运行服务
./etherscan
```

### 2. Docker 部署

#### 创建 Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o etherscan cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/etherscan .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./etherscan"]
```

#### 构建和运行

```bash
# 构建镜像
docker build -t etherscan:latest .

# 运行容器
docker run -d \
  --name etherscan \
  -p 8080:8080 \
  -v $(pwd)/config:/root/config \
  etherscan:latest
```

#### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  etherscan:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./config:/root/config
    environment:
      - GO_ENV=production
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f etherscan
```

### 3. Kubernetes 部署

#### 部署配置

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: etherscan
  labels:
    app: etherscan
spec:
  replicas: 3
  selector:
    matchLabels:
      app: etherscan
  template:
    metadata:
      labels:
        app: etherscan
    spec:
      containers:
      - name: etherscan
        image: etherscan:latest
        ports:
        - containerPort: 8080
        env:
        - name: GO_ENV
          value: "production"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: etherscan-service
spec:
  selector:
    app: etherscan
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

```bash
# 部署到 Kubernetes
kubectl apply -f k8s/deployment.yaml

# 查看状态
kubectl get pods -l app=etherscan
kubectl get services
```

## 配置管理

### 环境变量

支持通过环境变量覆盖配置：

```bash
export ETHEREUM_NODE_URL="wss://your-ethereum-node.com"
export SERVER_PORT="8080"
export CACHE_PRIMARY_EXPIRATION="5m"
export LOG_LEVEL="debug"
```

### 配置优先级

1. 环境变量（最高优先级）
2. 配置文件
3. 默认值（最低优先级）

## 监控和日志

### 健康检查

项目提供健康检查端点：

```bash
# 检查服务状态
curl http://localhost:8080/health

# 响应示例
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

### 指标监控

#### Prometheus 集成

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'etherscan'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics
    scrape_interval: 10s
```

#### 关键指标

- `http_requests_total`: HTTP 请求总数
- `http_request_duration_seconds`: 请求响应时间
- `cache_hits_total`: 缓存命中次数
- `cache_misses_total`: 缓存未命中次数
- `memory_usage_bytes`: 内存使用量
- `goroutines_count`: Goroutine 数量

### 日志管理

#### 日志格式

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "HTTP request processed",
  "method": "GET",
  "path": "/cached/blocks/latest",
  "status": 200,
  "duration": "25ms",
  "cache_hit": true
}
```

#### 日志收集

使用 ELK Stack 或类似工具收集和分析日志：

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/etherscan/*.log
  json.keys_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

## 性能优化

### 系统级优化

#### 内核参数调优

```bash
# /etc/sysctl.conf
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
```

#### 文件描述符限制

```bash
# /etc/security/limits.conf
* soft nofile 65535
* hard nofile 65535
```

### 应用级优化

#### Go 运行时调优

```bash
# 设置 GOMAXPROCS
export GOMAXPROCS=4

# 设置 GC 目标
export GOGC=100

# 设置内存限制
export GOMEMLIMIT=1GiB
```

#### 缓存优化

```go
// 优化缓存配置
config := &service.BlockCacheConfig{
    PrimaryExpiration:   2 * time.Minute,  // 根据访问模式调整
    SecondaryExpiration: 15 * time.Minute, // 平衡内存和性能
    MaxPrimaryBlocks:    20,               // 根据内存容量调整
    MaxSecondaryBlocks:  50,               // 根据历史数据需求调整
}
```

## 安全配置

### 网络安全

#### 防火墙配置

```bash
# 只允许必要端口
ufw allow 8080/tcp
ufw allow 22/tcp
ufw enable
```

#### 反向代理

使用 Nginx 作为反向代理：

```nginx
# /etc/nginx/sites-available/etherscan
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 限流配置
        limit_req zone=api burst=20 nodelay;
    }
}

# 限流配置
http {
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
}
```

### 应用安全

#### 环境变量保护

```bash
# 使用 .env 文件管理敏感信息
echo "ETHEREUM_NODE_URL=wss://your-private-node.com" > .env
chmod 600 .env
```

#### API 访问控制

```go
// 实现 API 密钥验证
func apiKeyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if !isValidAPIKey(apiKey) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## 故障排除

### 常见问题

#### 1. 服务启动失败

**检查步骤:**
```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 检查配置文件
go run cmd/main.go --config-check

# 查看详细日志
go run cmd/main.go --log-level debug
```

#### 2. 缓存性能问题

**诊断命令:**
```bash
# 检查内存使用
curl http://localhost:8080/debug/pprof/heap

# 检查缓存统计
curl http://localhost:8080/cache/stats
```

#### 3. 连接以太坊节点失败

**解决方案:**
- 检查网络连接
- 验证节点 URL
- 检查认证信息
- 尝试备用节点

### 日志分析

#### 错误日志模式

```bash
# 查找错误日志
grep "ERROR" /var/log/etherscan/app.log

# 分析响应时间
grep "duration" /var/log/etherscan/app.log | awk '{print $8}' | sort -n

# 统计缓存命中率
grep "cache_hit" /var/log/etherscan/app.log | grep "true" | wc -l
```

## 备份和恢复

### 配置备份

```bash
# 备份配置文件
tar -czf config-backup-$(date +%Y%m%d).tar.gz config/

# 定期备份脚本
#!/bin/bash
BACKUP_DIR="/backup/etherscan"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
tar -czf $BACKUP_DIR/config-$DATE.tar.gz config/
find $BACKUP_DIR -name "config-*.tar.gz" -mtime +7 -delete
```

### 数据恢复

```bash
# 恢复配置
tar -xzf config-backup-20240115.tar.gz

# 重启服务
systemctl restart etherscan
```

## 扩展和升级

### 水平扩展

#### 负载均衡配置

```nginx
upstream etherscan_backend {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
}

server {
    location / {
        proxy_pass http://etherscan_backend;
    }
}
```

### 版本升级

#### 滚动升级

```bash
# 1. 备份当前版本
cp etherscan etherscan.backup

# 2. 下载新版本
wget https://github.com/justinwongcn/etherscan/releases/latest/etherscan

# 3. 逐个替换实例
systemctl stop etherscan@1
cp etherscan /usr/local/bin/
systemctl start etherscan@1

# 4. 验证服务正常
curl http://localhost:8080/health
```

## 联系和支持

- **项目地址**: https://github.com/justinwongcn/etherscan
- **问题反馈**: https://github.com/justinwongcn/etherscan/issues
- **文档更新**: 请提交 PR 到文档目录
