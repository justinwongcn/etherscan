package ethereum

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/justinwongcn/go-ethlibs/node"
)

// managePool 管理连接池，定期清理空闲连接和进行健康检查
//
// Parameters:
//   - ctx: context.Context 用于控制连接池生命周期的上下文
//
// 该方法会启动一个后台协程，按照 idleTimeout 间隔定期：
//   - 检查连接池状态
//   - 执行连接健康检查（如果启用）
func (c *Client) managePool(ctx context.Context) {
	ticker := time.NewTicker(c.idleTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 执行健康检查
			if c.healthCheck {
				c.checkConnections(ctx)
			}
		}
	}
}

// getConnection 从连接池获取一个连接
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//
// Returns:
//   - node.Client: 从连接池获取的客户端连接
//   - error: 可能的错误：
//   - 连接池已满
//   - 无法创建新连接
func (c *Client) getConnection(_ context.Context) (node.Client, error) {
	// 检查是否达到最大连接数
	if atomic.LoadInt32(&c.connCount) >= int32(c.maxConns) {
		return nil, fmt.Errorf("connection pool is full (max: %d)", c.maxConns)
	}

	// 从连接池获取连接
	conn := c.connPool.Get()
	if conn == nil {
		return nil, fmt.Errorf("failed to get connection from pool")
	}

	// 增加连接计数
	atomic.AddInt32(&c.connCount, 1)

	return conn.(node.Client), nil
}

// releaseConnection 释放连接回连接池
//
// Parameters:
//   - conn: node.Client 要释放回连接池的客户端连接
//
// 该方法会：
//   - 减少活跃连接计数
//   - 将连接放回连接池以供复用
func (c *Client) releaseConnection(conn node.Client) {
	if conn == nil {
		return
	}

	// 减少连接计数
	atomic.AddInt32(&c.connCount, -1)

	// 将连接放回连接池
	c.connPool.Put(conn)
}

// checkConnections 检查连接池中的连接健康状态
//
// Parameters:
//   - ctx: context.Context 用于控制健康检查的上下文
//
// 该方法会：
//   - 获取一个连接进行测试
//   - 执行区块号查询来验证连接是否正常
//   - 如果连接异常，创建新的连接替换
func (c *Client) checkConnections(ctx context.Context) {
	// 获取一个连接进行健康检查
	conn, err := c.getConnection(ctx)
	if err != nil {
		return
	}
	defer c.releaseConnection(conn)

	// 执行一个简单的请求来检查连接是否正常
	_, err = conn.BlockNumber(ctx)
	if err != nil {
		// 连接异常，创建新的连接
		newConn, err := node.NewClient(ctx, c.nodeURL)
		if err == nil {
			c.connPool.Put(newConn)
		}
	}
}
