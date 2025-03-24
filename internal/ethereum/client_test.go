package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCreateRawTransaction(t *testing.T) {
	// 创建私钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	assert.NoError(t, err, "创建私钥失败")

	// 从私钥获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok, "公钥类型断言失败")

	// 从公钥获取地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	assert.NotEmpty(t, fromAddress.Hex(), "地址不应为空")

	// 设置交易参数
	nonce := uint64(0)
	value := big.NewInt(1000000000000000000) // 1 ETH
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000000) // 1 Gwei
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte

	// 创建交易对象
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	assert.NotNil(t, tx, "交易对象不应为空")

	// 验证交易字段
	assert.Equal(t, nonce, tx.Nonce(), "Nonce不匹配")
	assert.Equal(t, value, tx.Value(), "Value不匹配")
	assert.Equal(t, gasLimit, tx.Gas(), "GasLimit不匹配")
	assert.Equal(t, gasPrice, tx.GasPrice(), "GasPrice不匹配")
	assert.Equal(t, toAddress, *tx.To(), "接收地址不匹配")

	// 签名交易
	chainID := big.NewInt(4) // Rinkeby测试网络
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	assert.NoError(t, err, "交易签名失败")

	// 验证签名后的交易
	assert.NotNil(t, signedTx, "签名后的交易不应为空")

	// 获取RLP编码
	rawTxBytes, err := signedTx.MarshalBinary()
	assert.NoError(t, err, "RLP编码失败")
	rawTxHex := hex.EncodeToString(rawTxBytes)
	assert.NotEmpty(t, rawTxHex, "RLP编码不应为空")

	// 验证交易发送者
	sender, err := signer.Sender(signedTx)
	assert.NoError(t, err, "获取交易发送者失败")
	assert.Equal(t, fromAddress, sender, "交易发送者不匹配")
}

func TestCreateRawTransactionWithInvalidPrivateKey(t *testing.T) {
	// 测试无效私钥
	_, err := crypto.HexToECDSA("invalid_private_key")
	assert.Error(t, err, "应该检测到无效的私钥")
}

func TestCreateRawTransactionWithZeroValue(t *testing.T) {
	// 创建私钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	assert.NoError(t, err)

	// 设置交易参数（零值转账）
	nonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000000)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte

	// 创建并签名交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID := big.NewInt(4)
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	assert.NoError(t, err)

	// 验证零值交易
	assert.Equal(t, big.NewInt(0), signedTx.Value(), "交易值应为0")
}

func TestCreateRawTransactionWithCustomData(t *testing.T) {
	// 创建私钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	assert.NoError(t, err)

	// 设置交易参数（包含自定义数据）
	nonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000000)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	data := []byte("Hello, Ethereum!")

	// 创建并签名交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID := big.NewInt(4)
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	assert.NoError(t, err)

	// 验证自定义数据
	assert.Equal(t, data, signedTx.Data(), "交易数据不匹配")
}
