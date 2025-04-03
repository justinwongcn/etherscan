// Package domain 定义了以太坊区块链的核心领域模型
package domain

import (
	"github.com/justinwongcn/go-ethlibs/eth"
)

// Block 表示以太坊区块的领域模型
type Block struct {
	// Number 区块高度，表示该区块在区块链中的序号位置
	Number *uint64 `json:"number"`
	// Hash 区块哈希，使用Keccak-256算法计算得出的唯一标识符
	Hash *string `json:"hash"`
	// ParentHash 父区块的哈希值，用于维护区块链的链接关系
	ParentHash string `json:"parentHash"`
	// SHA3Uncles 叔块哈希列表的根哈希值
	SHA3Uncles string `json:"sha3Uncles"`
	// LogsBloom 区块中所有日志的布隆过滤器，用于快速查询日志
	LogsBloom string `json:"logsBloom"`
	// TransactionsRoot 交易默克尔树的根哈希值
	TransactionsRoot string `json:"transactionsRoot"`
	// StateRoot 状态树的根哈希值，表示该区块执行后的世界状态
	StateRoot string `json:"stateRoot"`
	// ReceiptsRoot 收据默克尔树的根哈希值
	ReceiptsRoot string `json:"receiptsRoot"`
	// Miner 挖矿奖励接收地址，即矿工地址
	Miner string `json:"miner"`
	// Author Parity客户端特有字段，等同于Miner字段
	Author string `json:"author,omitempty"` // Parity-specific alias of miner
	// Difficulty 区块难度值，用于调整挖矿难度
	Difficulty uint64 `json:"difficulty"`
	// TotalDifficulty 从创世区块到当前区块的总难度值
	TotalDifficulty uint64 `json:"totalDifficulty"`
	// ExtraData 区块的额外数据字段，最大32字节
	ExtraData string `json:"extraData"`
	// Size 区块大小，以字节为单位
	Size uint64 `json:"size"`
	// GasLimit 区块的燃料上限
	GasLimit uint64 `json:"gasLimit"`
	// GasUsed 区块中所有交易实际消耗的燃料总量
	GasUsed uint64 `json:"gasUsed"`
	// Timestamp 区块时间戳，Unix时间戳格式
	Timestamp uint64 `json:"timestamp"`
	// Transactions 区块包含的所有交易
	Transactions []eth.TxOrHash `json:"transactions"`
	// Uncles 叔块哈希列表
	Uncles []eth.Hash `json:"uncles"`

	// EIP-1559 BaseFeePerGas 基础费用，伦敦硬分叉后引入的动态基础费用
	BaseFeePerGas *uint64 `json:"baseFeePerGas,omitempty"`

	// EIP-4895 Withdrawals 相关字段
	// WithdrawalsRoot 提款操作的默克尔树根哈希值
	WithdrawalsRoot *string `json:"withdrawalsRoot,omitempty"`
	// Withdrawals 区块中包含的所有提款操作列表
	Withdrawals []eth.Withdrawal `json:"withdrawals,omitempty"`

	// EIP-4788 信标链相关字段
	// ParentBeaconBlockRoot 父信标链区块的根哈希值
	ParentBeaconBlockRoot *eth.Hash `json:"parentBeaconBlockRoot,omitempty"`

	// EIP-4844 Blob相关字段
	// ExcessBlobGas 超额Blob燃料值
	ExcessBlobGas *uint64 `json:"excessBlobGas,omitempty"`
	// BlobGasUsed 已使用的Blob燃料值
	BlobGasUsed *uint64 `json:"blobGasUsed,omitempty"`

	// Ethhash工作量证明相关字段
	// Nonce 用于挖矿的随机数
	Nonce *string `json:"nonce"`
	// MixHash 与Nonce一起用于验证区块是否满足难度要求
	MixHash *string `json:"mixHash"`

	// 权威证明(POA)相关字段
	// Step Aura共识的步进值
	Step *string `json:"step,omitempty"`
	// Signature 区块签名
	Signature *string `json:"signature,omitempty"`

	// Parity客户端特有字段
	// SealFields 区块密封字段列表
	SealFields *[]string `json:"sealFields,omitempty"`
}

// Withdrawal 表示提款操作的领域模型
type Withdrawal struct {
	Index          uint64 `json:"index"`          // 提款索引
	ValidatorIndex uint64 `json:"validatorIndex"` // 验证者索引
	Address        string `json:"address"`        // 接收地址
	Amount         uint64 `json:"amount"`         // 提款金额
}

// NewBlock 创建一个新的Block实例
func NewBlock() *Block {
	return &Block{}
}

//// UnmarshalJSON 实现自定义的JSON反序列化
//func (b *Block) UnmarshalJSON(data []byte) error {
//	type block Block
//	aliased := block(*b)
//
//	err := json.Unmarshal(data, &aliased)
//	if err != nil {
//		return err
//	}
//
//	*b = Block(aliased)
//	return nil
//}
//
//// MarshalJSON 实现自定义的JSON序列化
//func (b Block) MarshalJSON() ([]byte, error) {
//	type block Block
//	aliased := block(b)
//	return json.Marshal(&aliased)
//}
