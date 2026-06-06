package route

import (
	"wallet/internal/handler"
	"wallet/internal/service"
	"wallet/internal/store"

	"github.com/gin-gonic/gin"
	"wallet/internal/middleware"
)

func LoadRouter(store *store.Initializer) *gin.Engine {
	router := gin.Default()

	userServ := service.ExportUserService(*store)
	transactionServ := service.ExportTransactionService(*store)
	walletServ := service.ExportWalletService(*store)

	handler.RegisterTransactionRoutes(router.Group("/transactions"), transactionServ)
	handler.RegisterUserRoutes(router.Group("/users"), userServ)
	handler.RegisterWalletRoutes(router.Group("/wallets"), walletServ)
	
	router.GET("/ping", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
