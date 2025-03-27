package ethereum

import (
	"context"
	"fmt"

	"github.com/justinwongcn/go-ethlibs/eth"
	"github.com/justinwongcn/go-ethlibs/node"
)

// GetLogs 获取符合指定过滤条件的所有日志事件
//
// Parameters:
//   - ctx: context.Context 用于控制请求的上下文
//   - filter: *eth.LogFilter 日志过滤器，包含以下字段：
//   - FromBlock: 起始区块号（可选）
//   - ToBlock: 结束区块号（可选）
//   - Address: 合约地址（可选）
//   - Topics: 事件主题数组（可选）
//
// Returns:
//   - []*eth.Log: 匹配的日志事件数组
//   - error: 可能的错误：
//   - 无效的过滤器参数
//   - 节点连接错误
func (c *Client) GetLogs(ctx context.Context, filter *eth.LogFilter) ([]*eth.Log, error) {
	// 验证过滤器参数
	if filter == nil {
		return nil, fmt.Errorf("filter cannot be nil")
	}

	// 使用通用的连接池辅助函数执行操作
	result, err := c.withConnection(ctx, func(conn node.Client) (any, error) {
		return conn.Logs(ctx, *filter)
	})
	if err != nil {
		return nil, err
	}

	return result.([]*eth.Log), nil
}
