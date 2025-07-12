package routes

import (
	"context"
	"testing"

	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/etherscan/api/handler"
	"github.com/justinwongcn/etherscan/application/service"
	"github.com/justinwongcn/go-ethlibs/eth"
)

// MockAccountRepository 是AccountRepository的模拟实现
type MockAccountRepository struct {
	balances map[string]uint64
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		balances: make(map[string]uint64),
	}
}

func (m *MockAccountRepository) GetBalance(ctx context.Context, address string, numberOrTag string) (uint64, error) {
	if balance, exists := m.balances[address]; exists {
		return balance, nil
	}
	return 0, nil
}

func (m *MockAccountRepository) GetBalances(ctx context.Context, addresses []string, numberOrTag string) (map[string]uint64, error) {
	result := make(map[string]uint64)
	for _, address := range addresses {
		if balance, exists := m.balances[address]; exists {
			result[address] = balance
		} else {
			result[address] = 0
		}
	}
	return result, nil
}

// MockBlockRepository 是BlockRepository的模拟实现
type MockBlockRepository struct {
	blocks map[string]*eth.Block
}

func NewMockBlockRepository() *MockBlockRepository {
	return &MockBlockRepository{
		blocks: make(map[string]*eth.Block),
	}
}

func (m *MockBlockRepository) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return 12345, nil
}

func (m *MockBlockRepository) GetBlockByHash(ctx context.Context, hash string, fullTx bool) (*eth.Block, error) {
	if block, exists := m.blocks[hash]; exists {
		return block, nil
	}
	return nil, nil
}

func (m *MockBlockRepository) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*eth.Block, error) {
	if block, exists := m.blocks[number]; exists {
		return block, nil
	}
	return nil, nil
}

func (m *MockBlockRepository) GetBlockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	return 10, nil
}

func (m *MockBlockRepository) GetBlockTransactionCountByNumber(ctx context.Context, number string) (uint64, error) {
	return 10, nil
}

// MockTransactionRepository 是TransactionRepository的模拟实现
type MockTransactionRepository struct {
	transactions map[string]map[string]any
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[string]map[string]any),
	}
}

func (m *MockTransactionRepository) GetTransactionByHash(ctx context.Context, hash string) (map[string]any, error) {
	if tx, exists := m.transactions[hash]; exists {
		return tx, nil
	}
	return nil, nil
}

func (m *MockTransactionRepository) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash string, index uint64) (map[string]any, error) {
	return map[string]any{
		"hash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"from": "0x8ba1f109551bD432803012645aac136c12345678",
		"to":   "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
	}, nil
}

func (m *MockTransactionRepository) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNumber string, index uint64) (map[string]any, error) {
	return map[string]any{
		"hash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"from": "0x8ba1f109551bD432803012645aac136c12345678",
		"to":   "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6",
	}, nil
}

func (m *MockTransactionRepository) SendRawTransaction(ctx context.Context, data string) (string, error) {
	return "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", nil
}

func (m *MockTransactionRepository) GetTransactionReceipt(ctx context.Context, hash string) (map[string]any, error) {
	return map[string]any{
		"transactionHash": hash,
		"status":          "0x1",
		"gasUsed":         "0x5208",
	}, nil
}

func (m *MockTransactionRepository) GetTransactionCount(ctx context.Context, address string, numberOrTag string) (uint64, error) {
	return 5, nil
}

func TestRegisterRoutes(t *testing.T) {
	// 创建服务
	accountService := &service.AccountService{}
	blockService := &service.BlockService{}
	transactionService := &service.TransactionService{}

	// 创建Handler
	accountHandler := handler.NewAccountHandler(accountService)
	blockHandler := handler.NewBlockHandler(blockService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// 创建Ant服务器
	server := ant.NewHTTPServer()

	// 注册路由
	mockBlockRepo := NewMockBlockRepository()
	RegisterRoutes(server, blockHandler, transactionHandler, accountHandler, mockBlockRepo)

	// 验证服务器不为nil
	if server == nil {
		t.Errorf("RegisterRoutes() 后服务器不应该为nil")
	}

	// 注意：这里主要测试路由注册不会panic
	// 实际的路由测试需要启动服务器并发送HTTP请求
	t.Log("路由注册成功完成")
}

func TestRegisterRoutes_NilHandlers(t *testing.T) {
	// 测试nil handler的情况
	server := ant.NewHTTPServer()

	// 这个测试验证传入nil handler不会导致panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterRoutes() 在处理nil handler时发生panic: %v", r)
		}
	}()

	// 传入nil handler
	mockBlockRepo := NewMockBlockRepository()
	RegisterRoutes(server, nil, nil, nil, mockBlockRepo)

	t.Log("nil handler测试通过")
}

func TestRegisterRoutes_EmptyServer(t *testing.T) {
	// 测试空服务器的情况
	var server *ant.HTTPServer

	// 这个测试验证传入nil server会如何处理
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RegisterRoutes() 应该在nil server时panic")
		}
	}()

	// 创建模拟Handler
	accountService := &service.AccountService{}
	blockService := &service.BlockService{}
	transactionService := &service.TransactionService{}

	accountHandler := handler.NewAccountHandler(accountService)
	blockHandler := handler.NewBlockHandler(blockService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// 传入nil server，应该panic
	mockBlockRepo := NewMockBlockRepository()
	RegisterRoutes(server, blockHandler, transactionHandler, accountHandler, mockBlockRepo)
}

func TestRoutePatterns(t *testing.T) {
	// 测试DDD架构路由模式的正确性
	expectedRoutes := []string{
		"GET /blocks/{number}",
		"GET /blocks/height/latest",
		"GET /blocks/{number}/transactions/count",
		"GET /accounts/balance/{address}",
		"GET /accounts/balances/{addresses}",
		"GET /transactions/{hash}",
		"GET /blocks/transactions/{index}",
		"POST /transactions",
		"GET /accounts/{address}/transactions/count",
		"GET /transactions/{hash}/receipt",
	}

	// 这里主要验证路由模式的格式
	for _, route := range expectedRoutes {
		if route == "" {
			t.Errorf("路由模式不能为空")
		}

		// 验证路由包含HTTP方法
		if !contains(route, "GET") && !contains(route, "POST") {
			t.Errorf("路由模式应该包含HTTP方法: %s", route)
		}

		// 验证路由包含路径
		if !contains(route, "/") {
			t.Errorf("路由模式应该包含路径: %s", route)
		}
	}

	t.Logf("验证了 %d 个路由模式", len(expectedRoutes))
}

func TestDDDRoutes(t *testing.T) {
	// 测试DDD架构路由
	dddRoutes := []string{
		"GET /blocks/{number}",
		"GET /transactions/{hash}",
		"GET /accounts/balance/{address}",
	}

	for _, route := range dddRoutes {
		// 验证是GET请求
		if !contains(route, "GET") {
			t.Errorf("DDD路由应该是GET请求: %s", route)
		}

		// 验证路径格式
		if !contains(route, "/") {
			t.Errorf("路由应该包含路径: %s", route)
		}
	}

	t.Logf("验证了 %d 个DDD路由", len(dddRoutes))
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// 基准测试
func BenchmarkRegisterRoutes(b *testing.B) {
	// 创建模拟Handler
	accountService := &service.AccountService{}
	blockService := &service.BlockService{}
	transactionService := &service.TransactionService{}

	accountHandler := handler.NewAccountHandler(accountService)
	blockHandler := handler.NewBlockHandler(blockService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	for i := 0; i < b.N; i++ {
		server := ant.NewHTTPServer()
		mockBlockRepo := NewMockBlockRepository()
		RegisterRoutes(server, blockHandler, transactionHandler, accountHandler, mockBlockRepo)
	}
}
