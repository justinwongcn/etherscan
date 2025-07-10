// Package service 定义了领域服务
package service

import (
	"fmt"
	"sync"

	"github.com/justinwongcn/etherscan/domain/aggregate"
)

// DomainEventService 提供领域事件的发布和订阅服务
// 实现了事件驱动架构的核心功能
type DomainEventService struct {
	subscribers map[string][]EventHandler
	mutex       sync.RWMutex
}

// EventHandler 定义事件处理器的接口
type EventHandler interface {
	Handle(event aggregate.DomainEvent) error
	EventType() string
}

// NewDomainEventService 创建领域事件服务
func NewDomainEventService() *DomainEventService {
	return &DomainEventService{
		subscribers: make(map[string][]EventHandler),
		mutex:       sync.RWMutex{},
	}
}

// Subscribe 订阅特定类型的领域事件
func (s *DomainEventService) Subscribe(eventType string, handler EventHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.subscribers[eventType] == nil {
		s.subscribers[eventType] = make([]EventHandler, 0)
	}

	s.subscribers[eventType] = append(s.subscribers[eventType], handler)
}

// Unsubscribe 取消订阅特定类型的领域事件
func (s *DomainEventService) Unsubscribe(eventType string, handler EventHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	handlers := s.subscribers[eventType]
	if handlers == nil {
		return
	}

	// 查找并移除处理器
	for i, h := range handlers {
		if h == handler {
			s.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Publish 发布领域事件
func (s *DomainEventService) Publish(event aggregate.DomainEvent) error {
	s.mutex.RLock()
	handlers := s.subscribers[event.EventType()]
	s.mutex.RUnlock()

	if handlers == nil {
		return nil // 没有订阅者，直接返回
	}

	// 并发处理所有订阅者
	var wg sync.WaitGroup
	errors := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h.Handle(event); err != nil {
				errors <- fmt.Errorf("处理器 %T 处理事件失败: %w", h, err)
			}
		}(handler)
	}

	wg.Wait()
	close(errors)

	// 收集所有错误
	var allErrors []error
	for err := range errors {
		allErrors = append(allErrors, err)
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("发布事件时发生 %d 个错误: %v", len(allErrors), allErrors)
	}

	return nil
}

// PublishAggregateEvents 发布聚合根中的所有事件
func (s *DomainEventService) PublishAggregateEvents(aggregateRoot aggregate.AggregateRoot) error {
	// 获取聚合根的基础实现
	baseAggregate, ok := aggregateRoot.(*aggregate.BaseAggregateRoot)
	if !ok {
		// 如果不是BaseAggregateRoot，尝试通过类型断言获取事件
		switch ar := aggregateRoot.(type) {
		case *aggregate.BlockAggregate:
			return s.publishBlockEvents(ar)
		case *aggregate.TransactionAggregate:
			return s.publishTransactionEvents(ar)
		case *aggregate.AccountAggregate:
			return s.publishAccountEvents(ar)
		default:
			return fmt.Errorf("不支持的聚合根类型: %T", aggregateRoot)
		}
	}

	events := baseAggregate.GetEvents()

	for _, event := range events {
		if err := s.Publish(event); err != nil {
			return err
		}
	}

	// 清除已发布的事件
	baseAggregate.ClearEvents()

	return nil
}

// publishBlockEvents 发布区块聚合根的事件
func (s *DomainEventService) publishBlockEvents(blockAggregate *aggregate.BlockAggregate) error {
	// 获取区块聚合根的所有事件
	events := blockAggregate.GetEvents()

	// 发布每个事件
	for _, event := range events {
		if err := s.Publish(event); err != nil {
			return fmt.Errorf("发布区块事件失败: %w", err)
		}
	}

	// 清除已发布的事件
	blockAggregate.ClearEvents()

	return nil
}

// publishTransactionEvents 发布交易聚合根的事件
func (s *DomainEventService) publishTransactionEvents(transactionAggregate *aggregate.TransactionAggregate) error {
	// 获取交易聚合根的所有事件
	events := transactionAggregate.GetEvents()

	// 发布每个事件
	for _, event := range events {
		if err := s.Publish(event); err != nil {
			return fmt.Errorf("发布交易事件失败: %w", err)
		}
	}

	// 清除已发布的事件
	transactionAggregate.ClearEvents()

	return nil
}

// publishAccountEvents 发布账户聚合根的事件
func (s *DomainEventService) publishAccountEvents(accountAggregate *aggregate.AccountAggregate) error {
	// 获取账户聚合根的所有事件
	events := accountAggregate.GetEvents()

	// 发布每个事件
	for _, event := range events {
		if err := s.Publish(event); err != nil {
			return fmt.Errorf("发布账户事件失败: %w", err)
		}
	}

	// 清除已发布的事件
	accountAggregate.ClearEvents()

	return nil
}

// GetSubscriberCount 获取特定事件类型的订阅者数量
func (s *DomainEventService) GetSubscriberCount(eventType string) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.subscribers[eventType])
}

// GetAllEventTypes 获取所有已订阅的事件类型
func (s *DomainEventService) GetAllEventTypes() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	eventTypes := make([]string, 0, len(s.subscribers))
	for eventType := range s.subscribers {
		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes
}

// Clear 清除所有订阅者
func (s *DomainEventService) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.subscribers = make(map[string][]EventHandler)
}
