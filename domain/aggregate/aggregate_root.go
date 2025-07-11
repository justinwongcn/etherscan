// Package aggregate 定义了DDD中的聚合根相关概念
// 聚合根是聚合的入口点，负责维护聚合内部的一致性
package aggregate

// AggregateRoot 定义了聚合根的基础接口
// 所有聚合根都应该实现这个接口
type AggregateRoot interface {
	// ID 返回聚合根的唯一标识符
	ID() string

	// Equals 检查两个聚合根是否相等（基于ID）
	Equals(other AggregateRoot) bool
}

// DomainEvent 表示领域事件的接口
// 聚合根可以产生领域事件来通知其他聚合或应用服务
type DomainEvent interface {
	// EventType 返回事件类型
	EventType() string

	// AggregateID 返回产生事件的聚合根ID
	AggregateID() string

	// OccurredOn 返回事件发生的时间戳
	OccurredOn() int64
}

// EventPublisher 定义了事件发布器接口
// 用于发布领域事件
type EventPublisher interface {
	// Publish 发布领域事件
	Publish(event DomainEvent) error
}

// BaseAggregateRoot 提供聚合根的基础实现
// 其他聚合根可以嵌入这个结构体来获得基础功能
type BaseAggregateRoot struct {
	id     string
	events []DomainEvent
}

// NewBaseAggregateRoot 创建基础聚合根
func NewBaseAggregateRoot(id string) BaseAggregateRoot {
	return BaseAggregateRoot{
		id:     id,
		events: make([]DomainEvent, 0),
	}
}

// ID 返回聚合根的ID
func (b BaseAggregateRoot) ID() string {
	return b.id
}

// Equals 检查两个聚合根是否相等
func (b BaseAggregateRoot) Equals(other AggregateRoot) bool {
	return b.id == other.ID()
}

// AddEvent 添加领域事件
func (b *BaseAggregateRoot) AddEvent(event DomainEvent) {
	b.events = append(b.events, event)
}

// GetEvents 获取所有未发布的领域事件
func (b BaseAggregateRoot) GetEvents() []DomainEvent {
	return b.events
}

// ClearEvents 清除所有事件（通常在事件发布后调用）
func (b *BaseAggregateRoot) ClearEvents() {
	b.events = make([]DomainEvent, 0)
}
