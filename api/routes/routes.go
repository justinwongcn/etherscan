package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/justinwongcn/etherscan/api/handler"
)

// RegisterRoutes 注册所有HTTP路由
func RegisterRoutes(r *gin.Engine, blockHandler *handler.BlockHandler) {
	// 区块相关路由
	r.GET("/block/height", blockHandler.GetBlockHeight)
	r.GET("/block/:number", blockHandler.GetBlockByNumberOrHash)
}
