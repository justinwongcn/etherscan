// Package service 定义了领域服务
package service

import (
	"errors"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// AccountManagementService 提供账户管理相关的领域服务
// 处理跨账户的业务逻辑和复杂的账户操作
type AccountManagementService struct{}

// NewAccountManagementService 创建账户管理服务
func NewAccountManagementService() *AccountManagementService {
	return &AccountManagementService{}
}

// TransferResult 表示转账操作的结果
type TransferResult struct {
	Success           bool
	TransferredAmount valueobject.Wei
	TransactionFee    valueobject.Wei
	FromBalance       valueobject.Wei
	ToBalance         valueobject.Wei
}

// ProcessTransfer 处理账户间的转账操作
// 这是一个跨聚合的操作，需要协调多个账户聚合根
func (s *AccountManagementService) ProcessTransfer(
	fromAccount *aggregate.AccountAggregate,
	toAccount *aggregate.AccountAggregate,
	amount valueobject.Wei,
	transactionFee valueobject.Wei,
) (*TransferResult, error) {
	// 验证转账金额
	if amount.IsZero() {
		return nil, errors.New("转账金额不能为零")
	}

	// 计算总成本（转账金额 + 手续费）
	totalCost := amount.Add(transactionFee)

	// 检查发送方余额是否足够
	if !fromAccount.HasSufficientBalance(totalCost) {
		return &TransferResult{
			Success:     false,
			FromBalance: fromAccount.Balance(),
			ToBalance:   toAccount.Balance(),
		}, errors.New("发送方余额不足")
	}

	// 执行转账操作
	err := fromAccount.RecordSentTransaction(totalCost)
	if err != nil {
		return nil, err
	}

	toAccount.RecordReceivedTransaction(amount)

	return &TransferResult{
		Success:           true,
		TransferredAmount: amount,
		TransactionFee:    transactionFee,
		FromBalance:       fromAccount.Balance(),
		ToBalance:         toAccount.Balance(),
	}, nil
}

// ValidateAccountState 验证账户状态的一致性
func (s *AccountManagementService) ValidateAccountState(
	account *aggregate.AccountAggregate,
) error {
	// 检查基本字段
	if account.Address().IsZero() {
		return errors.New("账户地址不能为零")
	}

	// 检查余额不能为负数（虽然Wei值对象已经保证了这一点）
	if account.Balance().IsZero() && account.TotalReceived().GreaterThan(valueobject.NewWeiFromUint64(0)) {
		// 如果有接收记录但余额为零，检查是否有相应的发送记录
		if account.TotalSent().LessThan(account.TotalReceived()) {
			return errors.New("账户余额状态不一致")
		}
	}

	// 检查合约账户的特殊要求
	if account.IsContract() {
		if account.CodeHash() == nil {
			return errors.New("合约账户必须有代码哈希")
		}

		if account.CodeHash().IsZero() {
			return errors.New("合约账户的代码哈希不能为零")
		}
	}

	return nil
}

// CalculateAccountMetrics 计算账户的统计指标
type AccountMetrics struct {
	TotalTransactions  uint64
	TotalSent          valueobject.Wei
	TotalReceived      valueobject.Wei
	NetBalance         valueobject.Wei
	AverageTransaction valueobject.Wei
}

func (s *AccountManagementService) CalculateAccountMetrics(
	account *aggregate.AccountAggregate,
) *AccountMetrics {
	totalTx := account.TransactionCount()
	totalSent := account.TotalSent()
	totalReceived := account.TotalReceived()

	// 计算平均交易金额
	var avgTransaction valueobject.Wei
	if totalTx > 0 {
		totalVolume := totalSent.Add(totalReceived)
		avgTransaction, _ = totalVolume.Div(totalTx)
	} else {
		avgTransaction = valueobject.NewWeiFromUint64(0)
	}

	return &AccountMetrics{
		TotalTransactions:  totalTx,
		TotalSent:          totalSent,
		TotalReceived:      totalReceived,
		NetBalance:         account.Balance(),
		AverageTransaction: avgTransaction,
	}
}

// DetectSuspiciousActivity 检测可疑的账户活动
type SuspiciousActivityResult struct {
	IsSuspicious bool
	Reasons      []string
	RiskLevel    RiskLevel
}

type RiskLevel int

const (
	LowRisk RiskLevel = iota
	MediumRisk
	HighRisk
)

func (s *AccountManagementService) DetectSuspiciousActivity(
	account *aggregate.AccountAggregate,
	recentTransactions []*aggregate.TransactionAggregate,
) *SuspiciousActivityResult {
	result := &SuspiciousActivityResult{
		IsSuspicious: false,
		Reasons:      make([]string, 0),
		RiskLevel:    LowRisk,
	}

	// 检查大额交易
	largeTransactionThreshold := valueobject.NewWeiFromUint64(10000000000000000000) // 10 ETH
	largeTransactionCount := 0

	for _, tx := range recentTransactions {
		if tx.Value().GreaterThan(largeTransactionThreshold) {
			largeTransactionCount++
		}
	}

	if largeTransactionCount > 5 {
		result.IsSuspicious = true
		result.Reasons = append(result.Reasons, "频繁的大额交易")
		result.RiskLevel = MediumRisk
	}

	// 检查快速连续交易
	if len(recentTransactions) > 100 {
		result.IsSuspicious = true
		result.Reasons = append(result.Reasons, "异常高频的交易活动")
		result.RiskLevel = HighRisk
	}

	// 检查余额与交易活动的比例
	if account.Balance().IsZero() && account.TransactionCount() > 1000 {
		result.IsSuspicious = true
		result.Reasons = append(result.Reasons, "高交易量但零余额")
		result.RiskLevel = MediumRisk
	}

	return result
}

// CompareAccounts 比较两个账户的相似性
type AccountComparison struct {
	SimilarityScore float64
	CommonPatterns  []string
}

func (s *AccountManagementService) CompareAccounts(
	account1, account2 *aggregate.AccountAggregate,
) *AccountComparison {
	comparison := &AccountComparison{
		SimilarityScore: 0.0,
		CommonPatterns:  make([]string, 0),
	}

	// 比较账户类型
	if account1.Type() == account2.Type() {
		comparison.SimilarityScore += 0.2
		comparison.CommonPatterns = append(comparison.CommonPatterns, "相同账户类型")
	}

	// 比较交易数量级别
	tx1 := account1.TransactionCount()
	tx2 := account2.TransactionCount()

	if (tx1 < 10 && tx2 < 10) || (tx1 >= 10 && tx1 < 100 && tx2 >= 10 && tx2 < 100) ||
		(tx1 >= 100 && tx2 >= 100) {
		comparison.SimilarityScore += 0.3
		comparison.CommonPatterns = append(comparison.CommonPatterns, "相似的交易活跃度")
	}

	// 比较余额级别
	balance1 := account1.Balance()
	balance2 := account2.Balance()

	threshold1 := valueobject.NewWeiFromUint64(1000000000000000000)  // 1 ETH
	threshold2 := valueobject.NewWeiFromUint64(10000000000000000000) // 10 ETH

	if (balance1.LessThan(threshold1) && balance2.LessThan(threshold1)) ||
		(balance1.GreaterThan(threshold1) && balance1.LessThan(threshold2) &&
			balance2.GreaterThan(threshold1) && balance2.LessThan(threshold2)) ||
		(balance1.GreaterThan(threshold2) && balance2.GreaterThan(threshold2)) {
		comparison.SimilarityScore += 0.3
		comparison.CommonPatterns = append(comparison.CommonPatterns, "相似的余额级别")
	}

	return comparison
}
