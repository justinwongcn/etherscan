// Package service 定义了领域事件处理器
package service

import (
	"fmt"
	"log"

	"github.com/justinwongcn/etherscan/domain/aggregate"
)

// BlockCreatedEventHandler 处理区块创建事件
type BlockCreatedEventHandler struct {
	logger *log.Logger
}

// NewBlockCreatedEventHandler 创建区块创建事件处理器
func NewBlockCreatedEventHandler(logger *log.Logger) *BlockCreatedEventHandler {
	return &BlockCreatedEventHandler{
		logger: logger,
	}
}

// Handle 处理区块创建事件
func (h *BlockCreatedEventHandler) Handle(event aggregate.DomainEvent) error {
	blockEvent, ok := event.(*aggregate.BlockCreatedEvent)
	if !ok {
		return fmt.Errorf("期望BlockCreatedEvent，实际收到 %T", event)
	}

	h.logger.Printf("新区块已创建: 区块号=%s, 哈希=%s, 时间=%d",
		blockEvent.BlockNumber().String(),
		blockEvent.AggregateID(),
		blockEvent.OccurredOn())

	// 这里可以添加其他业务逻辑，比如：
	// - 更新区块链统计信息
	// - 通知其他系统
	// - 触发区块验证流程

	return nil
}

// EventType 返回处理的事件类型
func (h *BlockCreatedEventHandler) EventType() string {
	return "BlockCreated"
}

// TransactionMinedEventHandler 处理交易挖矿事件
type TransactionMinedEventHandler struct {
	logger *log.Logger
}

// NewTransactionMinedEventHandler 创建交易挖矿事件处理器
func NewTransactionMinedEventHandler(logger *log.Logger) *TransactionMinedEventHandler {
	return &TransactionMinedEventHandler{
		logger: logger,
	}
}

// Handle 处理交易挖矿事件
func (h *TransactionMinedEventHandler) Handle(event aggregate.DomainEvent) error {
	txEvent, ok := event.(*aggregate.TransactionMinedEvent)
	if !ok {
		return fmt.Errorf("期望TransactionMinedEvent，实际收到 %T", event)
	}

	h.logger.Printf("交易已挖矿: 交易哈希=%s, 区块哈希=%s, 时间=%d",
		txEvent.AggregateID(),
		txEvent.BlockHash(),
		txEvent.OccurredOn())

	// 这里可以添加其他业务逻辑，比如：
	// - 更新交易状态
	// - 通知用户交易确认
	// - 更新账户余额

	return nil
}

// EventType 返回处理的事件类型
func (h *TransactionMinedEventHandler) EventType() string {
	return "TransactionMined"
}

// BalanceUpdatedEventHandler 处理余额更新事件
type BalanceUpdatedEventHandler struct {
	logger *log.Logger
}

// NewBalanceUpdatedEventHandler 创建余额更新事件处理器
func NewBalanceUpdatedEventHandler(logger *log.Logger) *BalanceUpdatedEventHandler {
	return &BalanceUpdatedEventHandler{
		logger: logger,
	}
}

// Handle 处理余额更新事件
func (h *BalanceUpdatedEventHandler) Handle(event aggregate.DomainEvent) error {
	balanceEvent, ok := event.(*aggregate.BalanceUpdatedEvent)
	if !ok {
		return fmt.Errorf("期望BalanceUpdatedEvent，实际收到 %T", event)
	}

	h.logger.Printf("账户余额已更新: 地址=%s, 旧余额=%s, 新余额=%s, 时间=%d",
		balanceEvent.AggregateID(),
		balanceEvent.OldBalance(),
		balanceEvent.NewBalance(),
		balanceEvent.OccurredOn())

	// 这里可以添加其他业务逻辑，比如：
	// - 更新账户统计信息
	// - 检查余额阈值警告
	// - 触发风险控制检查

	return nil
}

// EventType 返回处理的事件类型
func (h *BalanceUpdatedEventHandler) EventType() string {
	return "BalanceUpdated"
}

// CompositeEventHandler 组合事件处理器，可以处理多种事件类型
type CompositeEventHandler struct {
	handlers map[string]EventHandler
	logger   *log.Logger
}

// NewCompositeEventHandler 创建组合事件处理器
func NewCompositeEventHandler(logger *log.Logger) *CompositeEventHandler {
	return &CompositeEventHandler{
		handlers: make(map[string]EventHandler),
		logger:   logger,
	}
}

// AddHandler 添加特定类型的事件处理器
func (h *CompositeEventHandler) AddHandler(eventType string, handler EventHandler) {
	h.handlers[eventType] = handler
}

// Handle 处理事件
func (h *CompositeEventHandler) Handle(event aggregate.DomainEvent) error {
	handler, exists := h.handlers[event.EventType()]
	if !exists {
		h.logger.Printf("没有找到事件类型 %s 的处理器", event.EventType())
		return nil
	}

	return handler.Handle(event)
}

// EventType 返回处理的事件类型（组合处理器支持多种类型）
func (h *CompositeEventHandler) EventType() string {
	return "Composite"
}

// AuditEventHandler 审计事件处理器，记录所有事件用于审计
type AuditEventHandler struct {
	logger *log.Logger
}

// NewAuditEventHandler 创建审计事件处理器
func NewAuditEventHandler(logger *log.Logger) *AuditEventHandler {
	return &AuditEventHandler{
		logger: logger,
	}
}

// Handle 处理事件（记录审计日志）
func (h *AuditEventHandler) Handle(event aggregate.DomainEvent) error {
	h.logger.Printf("审计日志: 事件类型=%s, 聚合ID=%s, 发生时间=%d",
		event.EventType(),
		event.AggregateID(),
		event.OccurredOn())

	// 这里可以将事件保存到审计数据库或文件
	// 用于合规性检查和系统监控

	return nil
}

// EventType 返回处理的事件类型
func (h *AuditEventHandler) EventType() string {
	return "*" // 表示处理所有类型的事件
}

// MetricsEventHandler 指标事件处理器，收集系统指标
type MetricsEventHandler struct {
	blockCount       int64
	transactionCount int64
	accountCount     int64
	logger           *log.Logger
}

// NewMetricsEventHandler 创建指标事件处理器
func NewMetricsEventHandler(logger *log.Logger) *MetricsEventHandler {
	return &MetricsEventHandler{
		logger: logger,
	}
}

// Handle 处理事件（更新指标）
func (h *MetricsEventHandler) Handle(event aggregate.DomainEvent) error {
	switch event.EventType() {
	case "BlockCreated":
		h.blockCount++
		h.logger.Printf("指标更新: 总区块数=%d", h.blockCount)

	case "TransactionCreated":
		h.transactionCount++
		h.logger.Printf("指标更新: 总交易数=%d", h.transactionCount)

	case "AccountCreated":
		h.accountCount++
		h.logger.Printf("指标更新: 总账户数=%d", h.accountCount)
	}

	return nil
}

// EventType 返回处理的事件类型
func (h *MetricsEventHandler) EventType() string {
	return "*" // 处理所有类型的事件
}

// GetMetrics 获取当前指标
func (h *MetricsEventHandler) GetMetrics() (int64, int64, int64) {
	return h.blockCount, h.transactionCount, h.accountCount
}
