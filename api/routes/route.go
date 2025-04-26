package routes

import (
	"github.com/gin-gonic/gin"
	"pledge-backend/api/controllers"
	"pledge-backend/api/middlewares"
	"pledge-backend/config"
)

func InitRoute(e *gin.Engine) *gin.Engine {

	// version group
	v2Group := e.Group("/api/v" + config.Config.Env.Version)

	// pledge-defi backend
	poolController := controllers.PoolController{}
	v2Group.GET("/poolBaseInfo", poolController.PoolBaseInfo)                                   //pool base information
	v2Group.GET("/poolDataInfo", poolController.PoolDataInfo)                                   //pool data information
	v2Group.GET("/token", poolController.TokenList)                                             //pool token information
	v2Group.POST("/pool/debtTokenList", middlewares.CheckToken(), poolController.DebtTokenList) //pool debtTokenList
	v2Group.POST("/pool/search", middlewares.CheckToken(), poolController.Search)               //pool search

	// 新需求
	ethController := controllers.EthController{}
	// 用于测试限流 limit.go 每个ip每秒超过2个请求报错
	v2Group.GET("/eth/test", ethController.Test)
	// 需求a
	v2Group.GET("/eth/block/:block_num", ethController.GetBlock)
	// 需求b
	v2Group.GET("/eth/tx/:tx_hash", ethController.GetTxByHash)
	// 需求c （需求d eventLoop.go）
	v2Group.GET("/etc/tx_receipt/:tx_hash", ethController.GetReceiptByHash)

	plgrController := controllers.PlgrController{}
	// 需求f 尝试通过API，让后端处理后调用ethclient来调用合约的功能
	v2Group.GET("/plgr/usdt/get_fee", plgrController.GetFee)

	// plgr-usdt price
	priceController := controllers.PriceController{}
	v2Group.GET("/price", priceController.NewPrice) //new price on ku-coin-exchange

	// pledge-defi admin backend
	multiSignPoolController := controllers.MultiSignPoolController{}
	v2Group.POST("/pool/setMultiSign", middlewares.CheckToken(), multiSignPoolController.SetMultiSign) //multi-sign set
	v2Group.POST("/pool/getMultiSign", middlewares.CheckToken(), multiSignPoolController.GetMultiSign) //multi-sign get

	userController := controllers.UserController{}
	v2Group.POST("/user/login", userController.Login)                             // login
	v2Group.POST("/user/logout", middlewares.CheckToken(), userController.Logout) // logout

	return e
}
