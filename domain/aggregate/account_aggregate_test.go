package aggregate

import (
	"testing"
	"time"

	"github.com/justinwongcn/etherscan/domain/valueobject"
)

func TestNewExternallyOwnedAccount(t *testing.T) {
	address, err := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	if err != nil {
		t.Fatalf("NewAddress() 失败: %v", err)
	}

	account, err := NewExternallyOwnedAccount(address)
	if err != nil {
		t.Errorf("NewExternallyOwnedAccount() 返回了意外的错误: %v", err)
		return
	}

	// 验证基本属性
	if !account.Address().Equals(address) {
		t.Errorf("Address() = %v, 期望 %v", account.Address(), address)
	}

	if account.Type() != ExternallyOwnedAccount {
		t.Errorf("Type() = %v, 期望 %v", account.Type(), ExternallyOwnedAccount)
	}

	if !account.Balance().IsZero() {
		t.Errorf("新账户的余额应该为零")
	}

	if account.Nonce() != 0 {
		t.Errorf("新账户的nonce应该为0")
	}

	if !account.IsExternallyOwned() {
		t.Errorf("应该是外部拥有账户")
	}

	if account.IsContract() {
		t.Errorf("不应该是合约账户")
	}

	// 验证事件
	events := account.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "AccountCreated" {
		t.Errorf("事件类型应该是 AccountCreated，实际是 %s", events[0].EventType())
	}
}

func TestNewExternallyOwnedAccount_ZeroAddress(t *testing.T) {
	zeroAddress, _ := valueobject.NewAddress("0x0000000000000000000000000000000000000000")

	_, err := NewExternallyOwnedAccount(zeroAddress)
	if err == nil {
		t.Errorf("NewExternallyOwnedAccount() 应该对零地址返回错误")
	}
}

func TestNewContractAccount(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	codeHash, _ := valueobject.NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	storageRoot, _ := valueobject.NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")

	account, err := NewContractAccount(address, codeHash, storageRoot)
	if err != nil {
		t.Errorf("NewContractAccount() 返回了意外的错误: %v", err)
		return
	}

	// 验证基本属性
	if !account.Address().Equals(address) {
		t.Errorf("Address() = %v, 期望 %v", account.Address(), address)
	}

	if account.Type() != ContractAccount {
		t.Errorf("Type() = %v, 期望 %v", account.Type(), ContractAccount)
	}

	if !account.IsContract() {
		t.Errorf("应该是合约账户")
	}

	if account.IsExternallyOwned() {
		t.Errorf("不应该是外部拥有账户")
	}

	// 验证合约特有属性
	if account.CodeHash() == nil || !account.CodeHash().Equals(codeHash) {
		t.Errorf("CodeHash() = %v, 期望 %v", account.CodeHash(), codeHash)
	}

	if account.StorageRoot() == nil || !account.StorageRoot().Equals(storageRoot) {
		t.Errorf("StorageRoot() = %v, 期望 %v", account.StorageRoot(), storageRoot)
	}
}

func TestAccountAggregate_UpdateBalance(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 清除创建事件
	account.ClearEvents()

	newBalance := valueobject.NewWeiFromUint64(1000000000000000000) // 1 ETH
	account.UpdateBalance(newBalance)

	// 验证余额更新
	if !account.Balance().Equals(newBalance) {
		t.Errorf("Balance() = %v, 期望 %v", account.Balance(), newBalance)
	}

	// 验证事件
	events := account.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个余额更新事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "BalanceUpdated" {
		t.Errorf("事件类型应该是 BalanceUpdated，实际是 %s", events[0].EventType())
	}
}

func TestAccountAggregate_IncrementNonce(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 清除创建事件
	account.ClearEvents()

	initialNonce := account.Nonce()
	account.IncrementNonce()

	// 验证nonce增加
	if account.Nonce() != initialNonce+1 {
		t.Errorf("Nonce() = %v, 期望 %v", account.Nonce(), initialNonce+1)
	}

	// 验证事件
	events := account.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个nonce更新事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "NonceUpdated" {
		t.Errorf("事件类型应该是 NonceUpdated，实际是 %s", events[0].EventType())
	}
}

func TestAccountAggregate_RecordSentTransaction(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 设置初始余额
	initialBalance := valueobject.NewWeiFromUint64(2000000000000000000) // 2 ETH
	account.UpdateBalance(initialBalance)

	// 清除之前的事件
	account.ClearEvents()

	amount := valueobject.NewWeiFromUint64(500000000000000000) // 0.5 ETH
	account.RecordSentTransaction(amount)

	// 验证交易计数增加
	if account.TransactionCount() != 1 {
		t.Errorf("TransactionCount() = %v, 期望 1", account.TransactionCount())
	}

	// 验证总发送金额
	if !account.TotalSent().Equals(amount) {
		t.Errorf("TotalSent() = %v, 期望 %v", account.TotalSent(), amount)
	}

	// 验证事件
	events := account.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个交易发送事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "TransactionSent" {
		t.Errorf("事件类型应该是 TransactionSent，实际是 %s", events[0].EventType())
	}
}

func TestAccountAggregate_RecordReceivedTransaction(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 清除创建事件
	account.ClearEvents()

	amount := valueobject.NewWeiFromUint64(1000000000000000000) // 1 ETH
	account.RecordReceivedTransaction(amount)

	// 验证总接收金额
	if !account.TotalReceived().Equals(amount) {
		t.Errorf("TotalReceived() = %v, 期望 %v", account.TotalReceived(), amount)
	}

	// 验证事件
	events := account.GetEvents()
	if len(events) != 1 {
		t.Errorf("应该有1个交易接收事件，实际有 %d 个", len(events))
	}

	if events[0].EventType() != "TransactionReceived" {
		t.Errorf("事件类型应该是 TransactionReceived，实际是 %s", events[0].EventType())
	}
}

func TestAccountAggregate_HasSufficientBalance(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 设置余额为1 ETH
	balance := valueobject.NewWeiFromUint64(1000000000000000000)
	account.UpdateBalance(balance)

	// 测试充足余额
	amount1 := valueobject.NewWeiFromUint64(500000000000000000) // 0.5 ETH
	if !account.HasSufficientBalance(amount1) {
		t.Errorf("应该有充足的余额")
	}

	// 测试不足余额
	amount2 := valueobject.NewWeiFromUint64(2000000000000000000) // 2 ETH
	if account.HasSufficientBalance(amount2) {
		t.Errorf("余额不足")
	}

	// 测试相等余额
	if !account.HasSufficientBalance(balance) {
		t.Errorf("相等余额应该充足")
	}
}

func TestAccountAggregate_CreatedAt(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")

	before := time.Now()
	account, _ := NewExternallyOwnedAccount(address)
	after := time.Now()

	createdAt := account.CreatedAt()
	if createdAt.Before(before) || createdAt.After(after) {
		t.Errorf("CreatedAt() 时间不在预期范围内")
	}
}

func TestAccountAggregate_LastUpdated(t *testing.T) {
	address, _ := valueobject.NewAddress("0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")
	account, _ := NewExternallyOwnedAccount(address)

	// 等待一小段时间
	time.Sleep(time.Millisecond)

	before := time.Now()
	newBalance := valueobject.NewWeiFromUint64(1000000000000000000)
	account.UpdateBalance(newBalance)
	after := time.Now()

	lastUpdated := account.LastUpdated()
	if lastUpdated.Before(before) || lastUpdated.After(after) {
		t.Errorf("LastUpdated() 时间不在预期范围内")
	}
}
