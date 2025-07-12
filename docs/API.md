# Etherscan API 文档

本文档详细描述了 Etherscan 项目提供的所有 API 端点。

## 基础信息

- **基础 URL**: `http://localhost:8080`
- **内容类型**: `application/json`
- **响应格式**: JSON

## 🔥 缓存版本 API（推荐）

缓存版本 API 提供高性能的区块链数据查询，具有以下优势：
- **高性能**: 缓存命中时响应时间 < 50ms
- **智能缓存**: 分层缓存策略，主缓存3分钟，二级缓存10分钟
- **并发优化**: 内置 SingleFlight 机制防止缓存击穿
- **格式兼容**: 自动处理不同格式的哈希值

### 区块高度查询

#### 获取最新区块高度

```http
GET /cached/blocks/height/latest
```

**响应示例:**
```json
{
  "height": "22898143"
}
```

**性能指标:**
- 缓存命中: ~20ms
- 缓存未命中: ~100-500ms

### 区块查询

#### 按区块号查询区块

```http
GET /cached/blocks/{number}?fullTx={boolean}
```

**参数:**
- `number` (路径参数): 区块号，支持十进制数字或 "latest"
- `fullTx` (查询参数): 是否返回完整交易信息
  - `true`: 返回完整交易对象
  - `false`: 只返回交易哈希列表

**请求示例:**
```bash
# 获取最新区块（不含完整交易）
GET /cached/blocks/latest?fullTx=false

# 获取指定区块（含完整交易）
GET /cached/blocks/12345?fullTx=true
```

**响应示例:**
```json
{
  "number": "0x15d65e0",
  "hash": "0x29f365785142b8b67ba4281e460b16332b70edfe89d37597f34efdf0a414e2d2",
  "parentHash": "0xc21139f3a50d4f419f20364b2581dba2f1bebcda2f60b207ba1ea2735da0b228",
  "miner": "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
  "timestamp": "0x68716bcb",
  "gasLimit": "0x226ebe5",
  "gasUsed": "0x19ec0ad",
  "transactions": ["0xf734625ffa59132bcfecdccd2540dcf80454900b37b59d5aba4f8120f8fd781e", "..."],
  "baseFeePerGas": "0x9097bff2"
}
```

#### 按区块哈希查询区块

```http
GET /cached/blocks/hash/{hash}?fullTx={boolean}
```

**参数:**
- `hash` (路径参数): 区块哈希，支持多种格式：
  - 标准格式: `0xabc123...` (小写，带0x前缀)
  - 大写格式: `0xABC123...` (大写，带0x前缀)
  - 无前缀小写: `abc123...`
  - 无前缀大写: `ABC123...`
- `fullTx` (查询参数): 同上

**请求示例:**
```bash
# 标准格式
GET /cached/blocks/hash/0xabc123def456789012345678901234567890123456789012345678901234567890?fullTx=false

# 大写格式（自动规范化）
GET /cached/blocks/hash/0xABC123DEF456789012345678901234567890123456789012345678901234567890?fullTx=false

# 无前缀格式（自动添加0x前缀）
GET /cached/blocks/hash/abc123def456789012345678901234567890123456789012345678901234567890?fullTx=false
```

**特性:**
- **自动规范化**: 不同格式的哈希会被自动规范化为统一格式
- **缓存共享**: 不同格式的相同哈希共享同一缓存条目
- **性能优化**: 规范化过程高度优化，对性能影响极小

## 📊 原始版本 API

原始版本 API 直接访问以太坊节点，不使用缓存。适用于需要实时数据的场景。

### 区块相关

#### 获取最新区块高度

```http
GET /blocks/height/latest
```

#### 按区块号获取区块

```http
GET /blocks/{number}
```

#### 获取区块交易数量

```http
GET /blocks/{number}/transactions/count
```

### 交易相关

#### 按哈希获取交易

```http
GET /transactions/{hash}
```

#### 按索引获取交易

```http
GET /blocks/transactions/{index}
```

#### 发送原始交易

```http
POST /transactions
```

#### 获取账户交易数量

```http
GET /accounts/{address}/transactions/count
```

#### 获取交易收据

```http
GET /transactions/{hash}/receipt
```

## 错误处理

### 错误响应格式

```json
{
  "error": "错误描述信息"
}
```

### 常见错误码

- **400 Bad Request**: 请求参数无效
- **404 Not Found**: 资源不存在
- **500 Internal Server Error**: 服务器内部错误

### 错误示例

```bash
# 无效区块号
GET /cached/blocks/invalid_number
# 响应: {"error": "获取区块失败: invalid block number"}

# 无效区块哈希
GET /cached/blocks/hash/invalid_hash
# 响应: {"error": "获取区块失败: invalid hash format"}
```

## 性能对比

### 响应时间对比

| API 类型 | 第一次请求 | 缓存命中 | 性能提升 |
|---------|-----------|---------|---------|
| 原始版本 | 100-500ms | N/A | - |
| 缓存版本 | 100-500ms | 20-50ms | 5-25倍 |

### 基准测试结果

```bash
# 缓存命中性能
BenchmarkCachedBlockService_CacheHit-8    366.5 ns/op    128 B/op    3 allocs/op

# 并发访问性能  
BenchmarkCachedBlockService_ConcurrentAccess-8    210.4 ns/op    128 B/op    3 allocs/op
```

## 使用建议

### 推荐使用场景

**缓存版本 API:**
- 频繁查询相同区块
- 对响应时间要求较高
- 高并发访问场景
- 历史区块数据查询

**原始版本 API:**
- 需要最新实时数据
- 一次性查询
- 对数据一致性要求极高

### 最佳实践

1. **优先使用缓存版本**: 在大多数场景下，缓存版本能提供更好的性能
2. **合理设置 fullTx 参数**: 根据实际需求选择是否需要完整交易信息
3. **利用哈希规范化**: 可以使用任意格式的哈希，系统会自动处理
4. **监控缓存命中率**: 通过日志或监控工具观察缓存效果

### 缓存策略

- **主缓存**: 存储最新的区块数据，3分钟过期
- **二级缓存**: 存储历史区块数据，10分钟过期
- **内存限制**: 主缓存50MB，二级缓存100MB
- **过期检查**: 被动过期，访问时检查

## 集成示例

### curl 示例

```bash
# 获取最新区块高度
curl "http://localhost:8080/cached/blocks/height/latest"

# 获取最新区块信息
curl "http://localhost:8080/cached/blocks/latest?fullTx=false"

# 按哈希查询区块
curl "http://localhost:8080/cached/blocks/hash/0xabc123...?fullTx=false"
```

### JavaScript 示例

```javascript
// 获取最新区块高度
async function getLatestBlockHeight() {
  const response = await fetch('http://localhost:8080/cached/blocks/height/latest');
  const data = await response.json();
  return data.height;
}

// 获取区块信息
async function getBlock(number, fullTx = false) {
  const response = await fetch(`http://localhost:8080/cached/blocks/${number}?fullTx=${fullTx}`);
  return await response.json();
}
```

### Python 示例

```python
import requests

# 获取最新区块高度
def get_latest_block_height():
    response = requests.get('http://localhost:8080/cached/blocks/height/latest')
    return response.json()['height']

# 获取区块信息
def get_block(number, full_tx=False):
    url = f'http://localhost:8080/cached/blocks/{number}?fullTx={str(full_tx).lower()}'
    response = requests.get(url)
    return response.json()
```
