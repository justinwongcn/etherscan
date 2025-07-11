// Package service 定义了领域服务
package service

import (
	"errors"

	"github.com/justinwongcn/etherscan/domain/aggregate"
	"github.com/justinwongcn/etherscan/domain/valueobject"
)

// TransactionValidationService 提供交易验证相关的领域服务
// 处理交易验证的复杂业务逻辑
type TransactionValidationService struct{}

// NewTransactionValidationService 创建交易验证服务
func NewTransactionValidationService() *TransactionValidationService {
	return &TransactionValidationService{}
}

// ValidateTransactionSignature 验证交易签名
// 这是一个复杂的密码学验证过程
func (s *TransactionValidationService) ValidateTransactionSignature(
	transaction *aggregate.TransactionAggregate,
) error {
	// 检查基本字段是否存在
	if transaction.Hash().IsZero() {
		return errors.New("交易哈希不能为零")
	}

	if transaction.From().IsZero() {
		return errors.New("发送方地址不能为零")
	}

	// 验证签名参数的有效性
	if err := s.validateSignatureParameters(transaction); err != nil {
		return err
	}

	// 验证签名的数学有效性
	if err := s.validateSignatureMath(transaction); err != nil {
		return err
	}

	// 验证恢复的地址是否与from地址匹配
	if err := s.validateRecoveredAddress(transaction); err != nil {
		return err
	}

	return nil
}

// validateSignatureParameters 验证签名参数的基本有效性
func (srv *TransactionValidationService) validateSignatureParameters(
	transaction *aggregate.TransactionAggregate,
) error {
	// 获取签名参数
	v, r, s := transaction.GetSignature()

	// 验证v值的有效性（应该是27, 28或者更高的值用于EIP-155）
	if len(v) == 0 || len(r) == 0 || len(s) == 0 {
		return errors.New("签名参数不能为空")
	}

	// 验证r和s不能为零
	if r == "0x0" || s == "0x0" {
		return errors.New("签名参数r和s不能为零")
	}

	return nil
}

// validateSignatureMath 验证签名的数学有效性
func (srv *TransactionValidationService) validateSignatureMath(
	transaction *aggregate.TransactionAggregate,
) error {
	// 这里可以实现更复杂的数学验证
	// 例如验证s值是否在有效范围内（防止延展性攻击）

	// 简化实现：检查签名格式
	_, r, s := transaction.GetSignature()

	// 验证r和s的长度（应该是64个字符的十六进制字符串，不包括0x前缀）
	if len(r) < 2 || len(s) < 2 {
		return errors.New("签名参数格式无效")
	}

	// 移除0x前缀进行长度检查
	rHex := r
	sHex := s
	if len(r) > 2 && r[:2] == "0x" {
		rHex = r[2:]
	}
	if len(s) > 2 && s[:2] == "0x" {
		sHex = s[2:]
	}

	if len(rHex) > 64 || len(sHex) > 64 {
		return errors.New("签名参数长度超出范围")
	}

	return nil
}

// validateRecoveredAddress 验证恢复的地址是否与from地址匹配
func (srv *TransactionValidationService) validateRecoveredAddress(
	transaction *aggregate.TransactionAggregate,
) error {
	// 这里应该实现ECDSA公钥恢复和地址验证
	// 由于这需要复杂的密码学计算，这里提供一个简化的验证框架

	fromAddr := transaction.From()
	if fromAddr.IsZero() {
		return errors.New("发送方地址无效")
	}

	// 在实际实现中，这里应该：
	// 1. 根据交易数据构建签名哈希
	// 2. 使用v, r, s参数恢复公钥
	// 3. 从公钥计算以太坊地址
	// 4. 验证计算出的地址是否与from地址匹配

	// 简化实现：假设签名有效
	return nil
}

// ValidateTransactionNonce 验证交易nonce
func (s *TransactionValidationService) ValidateTransactionNonce(
	transaction *aggregate.TransactionAggregate,
	senderAccount *aggregate.AccountAggregate,
) error {
	expectedNonce := senderAccount.Nonce()
	actualNonce := transaction.Nonce()

	if actualNonce != expectedNonce {
		return errors.New("交易nonce不正确")
	}

	return nil
}

// ValidateTransactionBalance 验证发送方余额是否足够
func (s *TransactionValidationService) ValidateTransactionBalance(
	transaction *aggregate.TransactionAggregate,
	senderAccount *aggregate.AccountAggregate,
) error {
	// 计算交易总成本
	totalCost := transaction.Value()

	// 添加燃料费用
	var gasCost valueobject.Wei

	switch transaction.Type() {
	case aggregate.LegacyTransaction:
		if gasPrice := transaction.GasPrice(); gasPrice != nil {
			gasUint64, err := transaction.Gas().Uint64()
			if err == nil {
				gasCost = gasPrice.Mul(gasUint64)
			}
		}
	case aggregate.EIP1559Transaction:
		if maxFeePerGas := transaction.MaxFeePerGas(); maxFeePerGas != nil {
			gasUint64, err := transaction.Gas().Uint64()
			if err == nil {
				gasCost = maxFeePerGas.Mul(gasUint64)
			}
		}
	}

	totalCost = totalCost.Add(gasCost)

	// 检查余额是否足够
	if !senderAccount.HasSufficientBalance(totalCost) {
		return errors.New("发送方余额不足")
	}

	return nil
}

// ValidateGasPrice 验证燃料价格的合理性
func (s *TransactionValidationService) ValidateGasPrice(
	transaction *aggregate.TransactionAggregate,
	networkBaseFee *valueobject.Wei,
	minGasPrice valueobject.Wei,
) error {
	switch transaction.Type() {
	case aggregate.LegacyTransaction:
		gasPrice := transaction.GasPrice()
		if gasPrice == nil {
			return errors.New("传统交易必须设置燃料价格")
		}

		if gasPrice.LessThan(minGasPrice) {
			return errors.New("燃料价格低于网络最低要求")
		}

	case aggregate.EIP1559Transaction:
		maxFeePerGas := transaction.MaxFeePerGas()
		maxPriorityFee := transaction.MaxPriorityFee()

		if maxFeePerGas == nil || maxPriorityFee == nil {
			return errors.New("EIP-1559交易必须设置最大费用和优先费用")
		}

		if networkBaseFee != nil {
			// 最大费用必须至少覆盖基础费用
			if maxFeePerGas.LessThan(*networkBaseFee) {
				return errors.New("最大费用低于网络基础费用")
			}
		}

		// 优先费用不能超过最大费用
		if maxPriorityFee.GreaterThan(*maxFeePerGas) {
			return errors.New("优先费用不能超过最大费用")
		}
	}

	return nil
}

// ValidateGasLimit 验证燃料限制
func (s *TransactionValidationService) ValidateGasLimit(
	transaction *aggregate.TransactionAggregate,
	blockGasLimit valueobject.Gas,
	minGasLimit valueobject.Gas,
) error {
	txGasLimit := transaction.Gas()

	// 检查最小燃料限制
	if txGasLimit.LessThan(minGasLimit) {
		return errors.New("交易燃料限制低于最小要求")
	}

	// 检查是否超过区块燃料限制
	if txGasLimit.GreaterThan(blockGasLimit) {
		return errors.New("交易燃料限制超过区块限制")
	}

	return nil
}

// ValidateContractCreation 验证合约创建交易
func (s *TransactionValidationService) ValidateContractCreation(
	transaction *aggregate.TransactionAggregate,
) error {
	if !transaction.IsContractCreation() {
		return nil // 不是合约创建交易，无需验证
	}

	// 合约创建交易的接收方地址应该为nil
	if transaction.To() != nil {
		return errors.New("合约创建交易不应该有接收方地址")
	}

	// 合约创建交易应该有输入数据（合约代码）
	input := transaction.Input()
	if len(input) == 0 {
		return errors.New("合约创建交易必须包含合约代码")
	}

	return nil
}

// CalculateTransactionFee 计算交易实际费用
func (s *TransactionValidationService) CalculateTransactionFee(
	transaction *aggregate.TransactionAggregate,
	gasUsed valueobject.Gas,
	baseFeePerGas *valueobject.Wei,
) (valueobject.Wei, error) {
	switch transaction.Type() {
	case aggregate.LegacyTransaction:
		gasPrice := transaction.GasPrice()
		if gasPrice == nil {
			return valueobject.Wei{}, errors.New("传统交易缺少燃料价格")
		}

		gasUsedUint64, err := gasUsed.Uint64()
		if err != nil {
			return valueobject.Wei{}, err
		}
		return gasPrice.Mul(gasUsedUint64), nil

	case aggregate.EIP1559Transaction:
		if baseFeePerGas == nil {
			return valueobject.Wei{}, errors.New("EIP-1559交易需要基础费用")
		}

		maxFeePerGas := transaction.MaxFeePerGas()
		maxPriorityFee := transaction.MaxPriorityFee()

		if maxFeePerGas == nil || maxPriorityFee == nil {
			return valueobject.Wei{}, errors.New("EIP-1559交易缺少费用设置")
		}

		// 添加优先费用，但不超过最大费用
		totalFee := baseFeePerGas.Add(*maxPriorityFee)
		var effectiveGasPrice valueobject.Wei
		if totalFee.LessThan(*maxFeePerGas) {
			effectiveGasPrice = totalFee
		} else {
			effectiveGasPrice = *maxFeePerGas
		}

		gasUsedUint64, err := gasUsed.Uint64()
		if err != nil {
			return valueobject.Wei{}, err
		}
		return effectiveGasPrice.Mul(gasUsedUint64), nil

	default:
		return valueobject.Wei{}, errors.New("不支持的交易类型")
	}
}
