// Package valueobject 定义了领域层的值对象
package valueobject

import (
	"github.com/justinwongcn/go-ethlibs/eth"
)

// validateWithEthLibs 使用go-ethlibs包验证哈希
// 这个函数专门用于调用go-ethlibs包，确保导入不被移除
func validateWithEthLibs(hashStr string) (string, error) {
	// 使用go-ethlibs包进行验证
	ethHash, err := eth.NewHash(hashStr)
	if err != nil {
		return "", err
	}

	// 返回标准化的哈希字符串
	return ethHash.String(), nil
}
