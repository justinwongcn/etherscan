// Package domain 提供以太坊区块链领域模型的定义
package domain

// TransactionReceipt 表示以太坊交易收据
// 该结构体包含交易执行后的详细结果信息，包括交易状态、gas使用情况、日志等
// 同时支持多个EIP标准引入的扩展字段，如EIP-2718的交易类型和EIP-4844的blob相关字段
type TransactionReceipt struct {
	// Type 交易类型（EIP-2718），用于区分不同类型的交易
	// 0x0: 传统交易
	// 0x1: EIP-2930访问列表交易
	// 0x2: EIP-1559费用市场交易
	// 0x3: EIP-4844 blob交易
	Type              *string      `json:"type,omitempty"`
	TransactionHash   string     `json:"transactionHash"`             // 交易哈希
	TransactionIndex  string       `json:"transactionIndex"`            // 交易在区块中的索引
	BlockHash         string     `json:"blockHash"`                   // 区块哈希
	BlockNumber       string       `json:"blockNumber"`                 // 区块号
	From              string  `json:"from"`                        // 交易发送方地址
	To                *string `json:"to"`                          // 交易接收方地址（合约创建交易为nil）
	CumulativeGasUsed string       `json:"cumulativeGasUsed"`           // 区块中截至该交易的累计gas使用量
	GasUsed           string       `json:"gasUsed"`                     // 该交易使用的gas量
	ContractAddress   *string `json:"contractAddress"`             // 如果是合约创建交易，则为新创建的合约地址
	Logs              []Log    `json:"logs"`                        // 交易产生的日志
	LogsBloom         string  `json:"logsBloom"`                   // 日志布隆过滤器
	Root              *string  `json:"root,omitempty"`              // 状态根（仅适用于前拜占庭分叉）
	Status            *string      `json:"status,omitempty"`            // 交易状态（1成功，0失败）
	EffectiveGasPrice *string      `json:"effectiveGasPrice,omitempty"` // 实际gas价格

	// EIP-4844 相关字段
	BlobGasPrice *string `json:"blobGasPrice,omitempty"` // blob数据的gas价格
	BlobGasUsed  *string `json:"blobGasUsed,omitempty"`  // blob数据使用的gas量
}


type Log struct {
	Removed     bool      `json:"removed"`
	LogIndex    *string `json:"logIndex"`
	TxIndex     *string `json:"transactionIndex"`
	TxHash      *string     `json:"transactionHash"`
	BlockHash   *string     `json:"blockHash"`
	BlockNumber *string `json:"blockNumber"`
	Address     string   `json:"address"`
	Data        string      `json:"data"`
	Topics      []string   `json:"topics"`

	// Parity-specific fields
	TxLogIndex *string `json:"transactionLogIndex,omitempty"`
	Type       *string   `json:"type,omitempty"`
}

// TransactionType 获取交易的类型
//
// Returns:
//   - string: 返回EIP-2718定义的交易类型，如果是传统交易则返回"0"
func (t *TransactionReceipt) TransactionType() string {
	if t.Type == nil {
		return "0"
	}

	return *t.Type
}
