package domain

import "github.com/justinwongcn/go-ethlibs/eth"

// Transaction 表示以太坊交易的领域模型
type Transaction struct {
	Type        *string      `json:"type,omitempty"`
	BlockHash   *eth.Hash    `json:"blockHash"`
	BlockNumber *string      `json:"blockNumber"`
	From        eth.Address  `json:"from"`
	Gas         string       `json:"gas"`
	Hash        eth.Hash     `json:"hash"`
	Input       eth.Input    `json:"input"`
	Nonce       string       `json:"nonce"`
	To          *eth.Address `json:"to"`
	Index       *string      `json:"transactionIndex"`
	Value       string       `json:"value"`
	V           string       `json:"v"`
	R           string       `json:"r"`
	S           string       `json:"s"`
	YParity     *string      `json:"yParity,omitempty"`

	// Gas价格 (EIP-1559中不包含此字段，可选)
	GasPrice *string `json:"gasPrice,omitempty"`

	// EIP-1559交易中的最大燃料费/优先燃料费 (仅EIP-1559交易包含，可选)
	MaxFeePerGas         *string `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *string `json:"maxPriorityFeePerGas,omitempty"`

	// Parity客户端特有字段
	StandardV *string        `json:"standardV,omitempty"`
	Raw       *eth.Data      `json:"raw,omitempty"`
	PublicKey *eth.Data      `json:"publicKey,omitempty"`
	ChainId   *string        `json:"chainId,omitempty"`
	Creates   *eth.Address   `json:"creates,omitempty"` // Parity文档中声称这是一个Hash
	Condition *eth.Condition `json:"condition,omitempty"`

	// EIP-2930访问列表
	AccessList *eth.AccessList `json:"accessList,omitempty"`

	// EIP-4844 blob交易字段
	MaxFeePerBlobGas    *string    `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes eth.Hashes `json:"blobVersionedHashes,omitempty"`

	// EIP-4844 Blob交易在"网络表示"中包含来自BlobsBundleV1引擎API模式的额外字段。
	// 但这些字段在执行层不可用，因此在处理交易的JSONRPC表示时不应出现，
	// 并且被排除在JSON序列化之外。此字段仅在解码"网络表示"的原始交易时填充，
	// 必须直接访问这些字段。
	BlobBundle *eth.BlobsBundleV1 `json:"-"`

	// EIP-7702授权列表
	AuthorizationList *eth.AuthorizationList `json:"authorizationList,omitempty"`
}

// AccessTuple 表示EIP-2930访问列表中的元素
type AccessTuple struct {
	Address     string   `json:"address"`     // 访问的地址
	StorageKeys []string `json:"storageKeys"` // 存储键列表
}

// Authorization 表示EIP-7702授权信息
type Authorization struct {
	Signer    string `json:"signer"`    // 签名者地址
	Signature string `json:"signature"` // 签名数据
}
